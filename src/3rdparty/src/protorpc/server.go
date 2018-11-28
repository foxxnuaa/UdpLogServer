// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protorpc

import (
	"3rdparty/src/proto"
	wire "3rdparty/src/protorpc/wire.pb"
	//"errors"
	"fmt"
	"io"
	_ "log"
	// "math/rand"
	"net/rpc"
	"reflect"
	//"sync"
	"strings"
	"sync"
	_ "unsafe"
)

var sConnID uint32

var gConnAuto uint32
var lockConnID sync.RWMutex

type ClosingConnFunction func(connID uint32, conn io.Closer) error

type ServerCodec struct {
	r         io.Reader
	w         io.Writer
	c         io.Closer
	closingFn reflect.Value
	// temporary work space
	reqHeader       wire.RequestHeader
	bNeedToResponse bool //to implement duplex rpc, server need not reponse from client  marlon 20130805
	connID          uint32
	IP              uint64
}

// NewServerCodec returns a ServerCodec that communicates with the ClientCodec
// on the other end of the given conn.
func NewServerCodec(conn io.ReadWriteCloser, uIP uint64, v ...interface{}) (codec rpc.ServerCodec) {
	//fmt.Println("NewServerCodec", *(*uint32)(unsafe.Pointer(&conn)))
	var closeFn reflect.Value
	ok := false

	// 第一个v参数为回调函数
	if len(v) >= 1 {
		if closeFn, ok = v[0].(reflect.Value); !ok {
			closeFn = reflect.ValueOf(nil)
		}
	}
	bNeedToResponse := true

	// 第二个是一个布尔值
	if len(v) >= 2 {
		if bNeedToResponse, ok = v[1].(bool); !ok {
			bNeedToResponse = true
		}
	}
	sCodec := &ServerCodec{
		r:               conn,
		w:               conn,
		c:               conn,
		bNeedToResponse: bNeedToResponse,
		closingFn:       closeFn,
	}

	lockConnID.Lock()

	gConnAuto = gConnAuto + 1

	sCodec.connID = gConnAuto
	sCodec.IP = uIP

	lockConnID.Unlock()

	//connPtr := (uintptr)(unsafe.Pointer(&sCodec.r))
	//valConnPtr := reflect.ValueOf(connPtr).Convert(reflect.TypeOf(sCodec.connID))
	//connID, ok := valConnPtr.Interface().(uint32)
	//if ok {
	//sCodec.connID = connID
	//sCodec.IP = uIP
	//} else {
	//	log.Println("get connID failed!")
	//}
	return sCodec
}

func (c *ServerCodec) GetRConn() *io.Reader {
	return &c.r
}

func (c *ServerCodec) GetConnID() uint32 {
	return c.connID
}

func (c *ServerCodec) ReadRequestHeader(r *rpc.Request) error {
	header := wire.RequestHeader{}
	if err := readRequestHeader(c.r, &header); err != nil {
		return err
	}
	r.ServiceMethod = header.GetMethod()
	strs := strings.Split(r.ServiceMethod, ".")
	if len(strs) == 3 {
		r.ServiceMethod = strs[1] + "." + strs[2]
	}
	r.Seq = header.GetId()
	c.reqHeader = header
	return nil
}

func (c *ServerCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	request, ok := x.(proto.Message)
	if !ok {
		return fmt.Errorf("ServerCodec.ReadRequestBody: %T does not implement proto.Message", x)
	}
	if err := readRequestBody(c.r, &c.reqHeader, request); err != nil {
		return err
	}

	c.writeValueToInterface("ConnId", reflect.ValueOf(&c.connID), &request)
	c.writeValueToInterface("IP", reflect.ValueOf(&c.IP), &request)

	c.reqHeader = wire.RequestHeader{}
	return nil
}

func (c *ServerCodec) writeValueToInterface(szFieldName string, interValue reflect.Value, request *proto.Message) error {

	valueRequest := reflect.ValueOf(*request)
	if valueRequest.Kind() == reflect.Ptr {
		valueRequest = valueRequest.Elem()
	}

	//print("writeValueToInterface111111111111111111111", szFieldName, interValue, "\n")

	connField := valueRequest.FieldByName(szFieldName)
	if !connField.IsValid() {
		return nil
	}
	/*if connField.IsNil() {
		if connField.Kind() == reflect.Ptr {
			z := reflect.New(connField.Type().Elem())
			connField.Set(z)
		} else {
			z := reflect.New(connField.Type())
			connField.Set(z)
		}
	}*/

	connField.Set(interValue)

	//fmt.Println("connField ", connField.Elem().Interface(), connID)
	return nil
}

// A value sent as a placeholder for the server's response value when the server
// receives an invalid request. It is never decoded by the client since the Response
// contains an error when it is used.
var invalidRequest = struct{}{}

func (c *ServerCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	//fmt.Println("$$$$$$$$$$$$$$$$$$$ServerCodec.WriteResponse")
	var response proto.Message
	if x != nil {
		var ok bool
		if response, ok = x.(proto.Message); !ok {
			if _, ok = x.(struct{}); !ok {
				return fmt.Errorf(
					"ServerCodec.WriteResponse: %T does not implement proto.Message",
					x,
				)
			}
		}
	}
	if c.bNeedToResponse {
		if err := writeResponse(c.w, r.Seq, r.Error, response); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerCodec) Close() error {
	//fmt.Println("server.Close")
	if s.closingFn.IsValid() {
		//fmt.Println("server.Close 1")
		s.closingFn.Call([]reflect.Value{reflect.ValueOf(s.connID), reflect.ValueOf(s.c)})
	}
	return s.c.Close()
}

// ServeConn runs the Protobuf-RPC server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
func ServeConn(conn io.ReadWriteCloser, closeFn ClosingConnFunction) {
	sCodec := NewServerCodec(conn, 0, closeFn)
	rpc.ServeCodec(sCodec)
}
