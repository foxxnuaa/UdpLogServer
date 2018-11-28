/**
* Base64
 */
package main

// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"testing"
	"time"

	"3rdparty/src/proto"
	"3rdparty/src/protorpc"

	// can not import xxx.pb with rpc stub here,
	// because it will cause import cycle.
	msg "3rdparty/src/protorpc/service.pb"
)

type Arith int

func (t *Arith) Add(args *msg.ArithRequest, reply *msg.ArithResponse) error {
	reply.C = proto.Int32(args.GetA() + args.GetB() + args.GetMarlon())
	log.Printf("Arith.Add(%v, %v, %v): %v", args.GetA(), args.GetB(), args.GetMarlon(), reply.GetC())
	return nil
}

func (t *Arith) Mul(args *msg.ArithRequest, reply *msg.ArithResponse) error {
	reply.C = proto.Int32(args.GetA() * args.GetB() * args.GetMarlon())
	log.Printf("Arith.Mul(%v, %v, %v): %v", args.GetA(), args.GetB(), args.GetMarlon(), reply.GetC())
	return nil
}

func (t *Arith) Div(args *msg.ArithRequest, reply *msg.ArithResponse) error {
	if args.GetB() == 0 {
		return errors.New("divide by zero")
	}
	reply.C = proto.Int32(args.GetA() / args.GetB())
	log.Printf("Arith.Div(%v, %v): %v", args.GetA(), args.GetB(), reply.GetC())
	return nil
}

func (t *Arith) Error(args *msg.ArithRequest, reply *msg.ArithResponse) error {
	return errors.New("ArithError")
}

func init() {
	err := listenAndServeArithAndEchoService("tcp", ":1234")
	if err != nil {
		log.Fatalf("listenAndServeArithAndEchoService: %v", err)
	}
}

func main() {
	time.Sleep(time.Hour)
}

func listenAndServeArithAndEchoService(network, addr string) error {
	clients, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("ArithService", new(Arith)); err != nil {
		return err
	}

	go func() {
		for {
			conn, err := clients.Accept()
			if err != nil {
				log.Printf("clients.Accept(): %v\n", err)
				continue
			}
			fmt.Println("conn:", conn)
			go srv.ServeCodec(protorpc.NewServerCodec(conn))
		}
	}()
	return nil
}

func testArithClient(t *testing.T, client *rpc.Client) {
	var args msg.ArithRequest
	var reply msg.ArithResponse
	var err error

	// Add
	args.A = proto.Int32(1)
	args.B = proto.Int32(2)
	if err = client.Call("ArithService.Add", &args, &reply); err != nil {
		t.Fatalf(`arith.Add: %v`, err)
	}
	if reply.GetC() != 3 {
		t.Fatalf(`arith.Add: expected = %d, got = %d`, 3, reply.GetC())
	}

	// Mul
	args.A = proto.Int32(2)
	args.B = proto.Int32(3)
	if err = client.Call("ArithService.Mul", &args, &reply); err != nil {
		t.Fatalf(`arith.Mul: %v`, err)
	}
	if reply.GetC() != 6 {
		t.Fatalf(`arith.Mul: expected = %d, got = %d`, 6, reply.GetC())
	}

	// Div
	args.A = proto.Int32(13)
	args.B = proto.Int32(5)
	if err = client.Call("ArithService.Div", &args, &reply); err != nil {
		t.Fatalf(`arith.Div: %v`, err)
	}
	if reply.GetC() != 2 {
		t.Fatalf(`arith.Div: expected = %d, got = %d`, 2, reply.GetC())
	}

	// Div zero
	args.A = proto.Int32(1)
	args.B = proto.Int32(0)
	if err = client.Call("ArithService.Div", &args, &reply); err.Error() != "divide by zero" {
		t.Fatalf(`arith.Error: expected = "%s", got = "%s"`, "divide by zero", err.Error())
	}

	// Error
	args.A = proto.Int32(1)
	args.B = proto.Int32(2)
	if err = client.Call("ArithService.Error", &args, &reply); err.Error() != "ArithError" {
		t.Fatalf(`arith.Error: expected = "%s", got = "%s"`, "ArithError", err.Error())
	}
}

func testArithClientAsync(t *testing.T, client *rpc.Client) {
	done := make(chan *rpc.Call, 16)
	callInfoList := []struct {
		method string
		args   *msg.ArithRequest
		reply  *msg.ArithResponse
		err    error
	}{
		{
			"ArithService.Add",
			&msg.ArithRequest{A: proto.Int32(1), B: proto.Int32(2)},
			&msg.ArithResponse{C: proto.Int32(3)},
			nil,
		},
		{
			"ArithService.Mul",
			&msg.ArithRequest{A: proto.Int32(2), B: proto.Int32(3)},
			&msg.ArithResponse{C: proto.Int32(6)},
			nil,
		},
		{
			"ArithService.Div",
			&msg.ArithRequest{A: proto.Int32(13), B: proto.Int32(5)},
			&msg.ArithResponse{C: proto.Int32(2)},
			nil,
		},
		{
			"ArithService.Div",
			&msg.ArithRequest{A: proto.Int32(1), B: proto.Int32(0)},
			&msg.ArithResponse{},
			errors.New("divide by zero"),
		},
		{
			"ArithService.Error",
			&msg.ArithRequest{A: proto.Int32(1), B: proto.Int32(2)},
			&msg.ArithResponse{},
			errors.New("ArithError"),
		},
	}

	// GoCall list
	calls := make([]*rpc.Call, len(callInfoList))
	for i := 0; i < len(calls); i++ {
		calls[i] = client.Go(callInfoList[i].method,
			callInfoList[i].args, callInfoList[i].reply,
			done,
		)
	}
	for i := 0; i < len(calls); i++ {
		<-calls[i].Done
	}

	// check result
	for i := 0; i < len(calls); i++ {
		if callInfoList[i].err != nil {
			if calls[i].Error.Error() != callInfoList[i].err.Error() {
				t.Fatalf(`%s: expected %v, Got = %v`,
					callInfoList[i].method,
					callInfoList[i].err,
					calls[i].Error,
				)
			}
			continue
		}

		got := calls[i].Reply.(*msg.ArithResponse).GetC()
		expected := callInfoList[i].reply.GetC()
		if got != expected {
			t.Fatalf(`%d: expected %v, Got = %v`,
				callInfoList[i].method, got, expected,
			)
		}
	}
}
