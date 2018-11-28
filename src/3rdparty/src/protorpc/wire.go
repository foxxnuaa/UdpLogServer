// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protorpc

import (
	"fmt"
	"hash/crc32"
	"io"

	"3rdparty/src/proto"
	"3rdparty/src/snappy"

	wire "3rdparty/src/protorpc/wire.pb"
	//"kCommon"
)

func WriteRequest_Out(w io.Writer, id uint64, method string, request proto.Message) error {
	// marshal request
	//fmt.Println("!!!!!!!!!!!!in wire, writeRequest")
	pbRequest := []byte{}
	if request != nil {
		var err error
		pbRequest, err = proto.Marshal(request)
		if err != nil {
			return err
		}
	}

	// compress serialized proto data
	compressedPbRequest, err := snappy.Encode(nil, pbRequest)
	if err != nil {
		return err
	}

	// generate header
	header := &wire.RequestHeader{
		Id:                         proto.Uint64(id),
		Method:                     proto.String(method),
		RawRequestLen:              proto.Uint32(uint32(len(pbRequest))),
		SnappyCompressedRequestLen: proto.Uint32(uint32(len(compressedPbRequest))),
		Checksum:                   proto.Uint32(crc32.ChecksumIEEE(compressedPbRequest)),
	}

	// check header size
	pbHeader, err := proto.Marshal(header)
	if err != err {
		return err
	}
	if uint32(len(pbHeader)) > wire.Default_Const_MaxHeaderLen {
		return fmt.Errorf("writeRequest: header larger than max_header_len: %d.", len(pbHeader))
	}
	if err := sendFrame(w, pbHeader); err != nil {
		//fmt.Println("!!!!!!!!!!!!in wire, writeRequest, after sendFrame, err", err)
		return err
	}
	if err := sendFrame(w, compressedPbRequest); err != nil {
		return err
	}
	/*pSendBuffer := append(pbHeader, compressedPbRequest...)
	if err := sendFrame(w, pSendBuffer); err != nil {
		return err
	}*/

	return nil
}

func writeRequest(w io.Writer, id uint64, method string, request proto.Message) error {
	// marshal request
	//fmt.Println("!!!!!!!!!!!!in wire, writeRequest")
	pbRequest := []byte{}
	if request != nil {
		var err error
		pbRequest, err = proto.Marshal(request)
		if err != nil {
			return err
		}
	}

	// compress serialized proto data
	compressedPbRequest, err := snappy.Encode(nil, pbRequest)
	if err != nil {
		return err
	}

	// generate header
	header := &wire.RequestHeader{
		Id:                         proto.Uint64(id),
		Method:                     proto.String(method),
		RawRequestLen:              proto.Uint32(uint32(len(pbRequest))),
		SnappyCompressedRequestLen: proto.Uint32(uint32(len(compressedPbRequest))),
		Checksum:                   proto.Uint32(crc32.ChecksumIEEE(compressedPbRequest)),
	}

	// check header size
	pbHeader, err := proto.Marshal(header)
	if err != err {
		return err
	}
	if uint32(len(pbHeader)) > wire.Default_Const_MaxHeaderLen {
		return fmt.Errorf("writeRequest: header larger than max_header_len: %d.", len(pbHeader))
	}
	if err := sendFrame(w, pbHeader); err != nil {
		//fmt.Println("!!!!!!!!!!!!in wire, writeRequest, after sendFrame, err", err)
		return err
	}
	if err := sendFrame(w, compressedPbRequest); err != nil {
		return err
	}
	/*pSendBuffer := append(pbHeader, compressedPbRequest...)
	if err := sendFrame(w, pSendBuffer); err != nil {
		return err
	}*/

	return nil
}

func readRequestHeader(r io.Reader, header *wire.RequestHeader) (err error) {
	// recv header (more)
	//fmt.Println("!!!!!!!!!!!!in wire, readRequestHeader, before recvFrame")
	pbHeader, err := recvFrame(r)
	fmt.Println("readRequestHeader ",pbHeader)
	if err != nil {
		fmt.Println("readRequestHeader1 ",err.Error())
		return err
	}

	// Marshal Header
	err = proto.Unmarshal(pbHeader, header)
	if err != nil {
		fmt.Println("readRequestHeader proto.Unmarshal error",err.Error(),string(pbHeader))
		return err
	}

	return nil
}

func readRequestBody(r io.Reader, header *wire.RequestHeader, request proto.Message) error {
	// recv body (end)
	compressedPbRequest, err := recvFrame(r)
	if err != nil {
		fmt.Println("readRequestBody error1",err.Error())
		return err
	}

	// checksum
	if crc32.ChecksumIEEE(compressedPbRequest) != header.GetChecksum() {
		fmt.Println("readRequestBody: unexpected checksum")
		return fmt.Errorf("readRequestBody: unexpected checksum.")
	}

	// decode the compressed data
	pbRequest, err := snappy.Decode(nil, compressedPbRequest)
	if err != nil {
		fmt.Println("readRequestBody snappy.Decode error",err.Error())
		return err
	}

	// check wire header: rawMsgLen
	if uint32(len(pbRequest)) != header.GetRawRequestLen() {
		return fmt.Errorf("readRequestBody: Unexcpeted header.RawRequestLen.")
	}

	//decrypt:huyanglian added 20141120
	//strRequestDest := kCommon.ProtocolDecrypt(pbRequest)
	//pbRequestDest := []byte(strRequestDest)
	// Unmarshal to proto message
	if request != nil {
		err = proto.Unmarshal(pbRequest, request)
		if err != nil {
			fmt.Println("readRequestBody proto.Unmarshal error",err.Error(),string(pbRequest))
			return err
		}
	}

	return nil
}

func writeResponse(w io.Writer, id uint64, serr string, response proto.Message) (err error) {
	// clean response if error
	if serr != "" {
		response = nil
	}

	// marshal response
	pbResponse := []byte{}
	if response != nil {
		pbResponse, err = proto.Marshal(response)
		if err != nil {
			return err
		}
	}

	// compress serialized proto data
	compressedPbResponse, err := snappy.Encode(nil, pbResponse)
	if err != nil {
		return err
	}

	// generate header
	header := &wire.ResponseHeader{
		Id:                          proto.Uint64(id),
		Error:                       proto.String(serr),
		RawResponseLen:              proto.Uint32(uint32(len(pbResponse))),
		SnappyCompressedResponseLen: proto.Uint32(uint32(len(compressedPbResponse))),
		Checksum:                    proto.Uint32(crc32.ChecksumIEEE(compressedPbResponse)),
	}

	// check header size
	pbHeader, err := proto.Marshal(header)
	if err != err {
		return
	}
	if uint32(len(pbHeader)) > wire.Default_Const_MaxHeaderLen {
		return fmt.Errorf("writeResponse: header larger than max_header_len: %d.",
			len(pbHeader),
		)
	}

	// send header (more)
	if err = sendFrame(w, pbHeader); err != nil {
		return
	}

	// send body (end)
	if err = sendFrame(w, compressedPbResponse); err != nil {
		return
	}

	return nil
}

func readResponseHeader(r io.Reader, header *wire.ResponseHeader) error {
	// recv header (more)
	pbHeader, err := recvFrame(r)
	if err != nil {
		return err
	}

	// Marshal Header
	err = proto.Unmarshal(pbHeader, header)
	if err != nil {
		return err
	}

	return nil
}

func readResponseBody(r io.Reader, header *wire.ResponseHeader, response proto.Message) error {
	// recv body (end)
	compressedPbResponse, err := recvFrame(r)
	if err != nil {
		return err
	}

	// checksum
	if crc32.ChecksumIEEE(compressedPbResponse) != header.GetChecksum() {
		return fmt.Errorf("readResponseBody: unexpected checksum.")
	}

	// decode the compressed data
	pbResponse, err := snappy.Decode(nil, compressedPbResponse)
	if err != nil {
		return err
	}
	// check wire header: rawMsgLen
	if uint32(len(pbResponse)) != header.GetRawResponseLen() {
		return fmt.Errorf("readResponseBody: Unexcpeted header.RawResponseLen.")
	}

	// Unmarshal to proto message
	if response != nil {
		err = proto.Unmarshal(pbResponse, response)
		if err != nil {
			return err
		}
	}

	return nil
}
