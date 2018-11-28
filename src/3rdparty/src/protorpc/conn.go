// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protorpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
	"strings"
)

func sendFrame(w io.Writer, data []byte) (err error) {
	// Allocate enough space for the biggest uvarint
	var size [4]byte

	if data == nil || len(data) == 0 {
		//n := binary.PutUvarint(size[:], uint64(0))
		binary.BigEndian.PutUint32(size[:],uint32(0))
		if err = write(w, size[:4], false); err != nil {
			return
		}
		return
	}

	//TODO: encrypt data. marlon 20130820
	//added huyanglian 20141125

	//strResponseDest := kCommon.ProtocolEncrypt(data)
	//data = []byte(strResponseDest)
	//fmt.Println("send data", data)

	// Write the size and data
	//n := binary.PutUvarint(size[:], uint64(len(data)))
	binary.BigEndian.PutUint32(size[:], uint32(len(data)))
	fmt.Println("len(data)1==",len(data))
	if err = write(w, size[:4], false); err != nil {
		return
	}
	//send 2k data per time
	iTotalLen := len(data)
	var iSendedLen int = 0
	for iSendedLen < iTotalLen {
		var iEndPos int = iSendedLen + 1024*2
		if iEndPos > iTotalLen {
			iEndPos = iTotalLen
		}
		if err = write(w, data[iSendedLen:iEndPos], false); err != nil {
			return
		}
		iSendedLen = iEndPos
	}
	/*if err = write(w, data, false); err != nil {
		return
	}*/
	return
}

func recvFrame(r io.Reader) (data []byte, err error) {
	fmt.Println("reveive data", data)
	size, err := readUint32(r)
	//binary.BigEndian.Uint32(r)
	if err != nil {
		return nil, err
	}
	if size > 65535 {
		return nil, errors.New("recvFrame too large")
	}
	if size != 0 {
		data = make([]byte, size)
		if err = read(r, data); err != nil {
			return nil, err
		}
	}

	//TODO: decrypt data. marlon 20130820
	//decrypt:huyanglian added 20141120

	//strRequestDest := kCommon.ProtocolDecrypt(data)
	//data = []byte(strRequestDest)
	//fmt.Println("reveive data", data)

	return data, nil
}

func RecvFrame_Out(r io.Reader) (data []byte, err error) {
	//fmt.Println("reveive data", data)
	//size, err := readUvarint(r)
	size, err := readUint32(r)
	if err != nil {
		return nil, err
	}
	if size > 65535 {
		return nil, errors.New("recvFrame too large")
	}
	if size != 0 {
		data = make([]byte, size)
		if err = read(r, data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func readUint32(r io.Reader) (uint32, error) {
	buff := make([]byte, 4)
	for i := 0; i < 4; i++ {
		var b byte
		b, err := readByte(r)

		if err != nil{
			fmt.Println("readUint32 error",err.Error())
			return 0,err
		}
		buff[i] = b
	}
	return binary.BigEndian.Uint32(buff),nil
}

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func readUvarint(r io.Reader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		var b byte
		b, err := readByte(r)
		if err != nil {
			return x, err
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, errors.New("varint overflows a 64-bit integer")
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

func write(w io.Writer, data []byte, onePacket bool) error {
	if onePacket {
		if _, err := w.Write(data); err != nil {
			return err
		}
		return nil
	}
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return err
			}
		}
		index += n
	}
	return nil
}

func read(r io.Reader, data []byte) error {
	for index := 0; index < len(data); {
		n, err := r.Read(data[index:])
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return err
			}
			if strings.Contains(err.Error(),"closed"){
				fmt.Println("read error",err.Error())
				return err
			}
		}
		if n ==0{
			time.Sleep(time.Second/1000)
		}
		index += n
	}
	return nil
}

func readByte(r io.Reader) (c byte, err error) {
	data := make([]byte, 1)
	if err = read(r, data); err != nil {
		fmt.Println("readByte error",err.Error())
		return 0, err
	}
	c = data[0]
	return
}
