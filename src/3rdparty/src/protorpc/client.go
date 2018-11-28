// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protorpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"sync"
	"time"

	"3rdparty/src/proto"

	wire "3rdparty/src/protorpc/wire.pb"
)

type clientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	// temporary work space
	respHeader wire.ResponseHeader

	mutex sync.Mutex

	//to implement duplex rpc, server need not reponse from client  marlon 20130805
	bNeedResponse  bool
	requestSeq     uint64
	signalToRead   chan int
	signalReadOver chan int
	signalError    chan error
}

// NewClientCodec returns a new rpc.ClientCodec using Protobuf-RPC on conn.
func NewClientCodec(conn io.ReadWriteCloser, v ...interface{}) rpc.ClientCodec {
	bNeedResponse := true
	ok := false
	if len(v) >= 1 {
		if bNeedResponse, ok = v[0].(bool); !ok {
			bNeedResponse = true
		}
	}
	return &clientCodec{
		r:              conn,
		w:              conn,
		c:              conn,
		bNeedResponse:  bNeedResponse,
		signalToRead:   make(chan int, 1),
		signalReadOver: make(chan int, 1),
		signalError:    make(chan error, 1),
	}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) (err error) {
	//fmt.Println("@@@@@@before clientCodec.WriteRequest")
	err = nil
	var request proto.Message
	if param != nil {
		var ok bool
		if request, ok = param.(proto.Message); !ok {
			return fmt.Errorf(
				"ClientCodec.WriteRequest: %T does not implement proto.Message",
				param,
			)
		}
	}

	c.mutex.Lock()
	c.requestSeq = r.Seq

	defer func() {
		if !c.bNeedResponse {
			//在不需要回复的情况下，如果WriteRequest在多个goroutine中被调用，为保证WriteRequest和ReadResponseHeader的同步，
			//即两者合起来作为一个原子操作，这里用channel来同步
			c.signalToRead <- 1
			c.signalError <- err
			<-c.signalReadOver
		}
		c.mutex.Unlock()
	}()

	//fmt.Println("@@@@@@before writeRequest")
	err = writeRequest(c.w, r.Seq, r.ServiceMethod, request)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	header := wire.ResponseHeader{}
	if c.bNeedResponse {
		err := readResponseHeader(c.r, &header)
		if err != nil {
			return err
		}
	} else {
		<-c.signalToRead //在不需要response时，通过channel模拟read的阻塞
		err := <-c.signalError
		if err != nil {
			r.Seq = c.requestSeq
			c.signalReadOver <- 1
			return err
		}
	}

	if c.bNeedResponse {
		r.Seq = header.GetId()
	} else {
		r.Seq = c.requestSeq
		c.signalReadOver <- 1
	}
	r.Error = header.GetError()
	c.respHeader = header
	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	if !c.bNeedResponse {
		return nil //add by marlon, to implement duplex rpc
	}
	var response proto.Message
	if x != nil {
		var ok bool
		response, ok = x.(proto.Message)
		if !ok {
			return fmt.Errorf(
				"ClientCodec.ReadResponseBody: %T does not implement proto.Message",
				x,
			)
		}
	}

	err := readResponseBody(c.r, &c.respHeader, response)
	if err != nil {
		return nil
	}

	c.respHeader = wire.ResponseHeader{}
	return nil
}

// Close closes the underlying connection.
func (c *clientCodec) Close() error {
	if !c.bNeedResponse {
		c.mutex.Lock()
		select {
		case c.signalToRead <- 1: //防止ReadResponseHeader阻塞导致rpc.Client.input goroutine不能退出
			select {
			case c.signalError <- errors.New(fmt.Sprintln("rpc client close")):
				select {
				case <-c.signalReadOver:
				default:
				}
			default:
			}
		default:
		}
		c.mutex.Unlock()
	}
	return c.c.Close()
}

// NewClient returns a new rpc.Client to handle requests to the
// set of services at the other end of the connection.
func NewClient(conn io.ReadWriteCloser, v ...interface{}) *rpc.Client {
	bNeedResponse := true
	ok := false
	if len(v) >= 1 {
		if bNeedResponse, ok = v[0].(bool); !ok {
			bNeedResponse = true
		}
	}
	return rpc.NewClientWithCodec(NewClientCodec(conn, bNeedResponse))
}

// Dial connects to a Protobuf-RPC server at the specified network address.
func Dial(network, address string, v ...interface{}) (*rpc.Client, error) {
	bNeedResponse := true
	ok := false
	if len(v) >= 1 {
		if bNeedResponse, ok = v[0].(bool); !ok {
			bNeedResponse = true
		}
	}

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, bNeedResponse), err
}

// DialTimeout connects to a Protobuf-RPC server at the specified network address.
func DialTimeout(network, address string, timeout time.Duration, v ...interface{}) (*rpc.Client, error) {
	bNeedResponse := true
	ok := false
	if len(v) >= 1 {
		if bNeedResponse, ok = v[0].(bool); !ok {
			bNeedResponse = true
		}
	}

	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, bNeedResponse), err
}
