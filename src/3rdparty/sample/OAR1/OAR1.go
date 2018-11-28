/**
* Base64
 */
package main

import (
	//"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"3rdparty/src/proto"
	"3rdparty/src/protorpc"

	// can not import xxx.pb with rpc stub here,
	// because it will cause import cycle.
	OAR1 "3rdparty/src/protorpc/OAR1"
)

type AddressCheck int

func (t *AddressCheck) Check(args *OAR1.AddressBook, reply *OAR1.CheckResponse) error {
	log.Printf("server get from client address book!============================")
	persons := args.GetPerson()
	for _, person := range persons {
		fmt.Println("name:", person.GetName())
		fmt.Println("ID:", person.GetId())
		fmt.Println("email:", person.GetEmail())
		for _, phone := range person.GetPhone() {
			fmt.Println("phone number:", phone.GetNumber(), ", type:", phone.GetType())
		}
	}
	reply.Msg = proto.String("OK, got you!")
	return nil
}

func init() {
	//client rpc server
	err := listenAndServeClickService("tcp", ":1234")
	if err != nil {
		log.Fatalf("listenAndServeArithAndEchoService: %v", err)
	}
}

func main() {
	time.Sleep(time.Hour)
}

func listenAndServeClickService(network, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	listener, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("CheckService", new(AddressCheck)); err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("clients.Accept(): %v\n", err)
				continue
			}
			fmt.Println("!!!!!!!!!!server receive connect from client address: ", conn.RemoteAddr())
			go srv.ServeCodec(protorpc.NewServerCodec(conn, nil, false))
			client := rpc.NewClientWithCodec(protorpc.NewClientCodec(conn, false))
			go pushToClient(client)
			//go pushToClient(client)
		}
	}()
	return nil
}

func pushToClient(client *rpc.Client) {

	//defer client.Close()

	for {
		var args OAR1.AddressBook
		var reply OAR1.CheckResponse
		person := new(OAR1.Person)
		str1 := "zhouzhifeng"
		arr1 := make([]byte, 0, 100)
		arr1 = Append(arr1, []byte(str1))
		arr1 = append(arr1, 0)
		arr1 = Append(arr1, []byte("after zreo"))
		person.Name = arr1 //proto.String(string(arr1))
		fmt.Println("person.Name:", person.Name)
		person.Id = proto.Int32(1111)
		//person.Email = proto.String("lionzhou@hhhh")
		phoneNumber := new(OAR1.Person_PhoneNumber)
		phoneNumber.Number = proto.String("13645789524")
		phoneNumber.Type = OAR1.Person_MOBILE.Enum()
		person.Phone = append(person.Phone, phoneNumber)
		args.Person = append(args.Person, person)

		person = new(OAR1.Person)
		person.Name = []byte("zhouzhifeng2")
		person.Id = proto.Int32(2222)
		person.Email = proto.String("lionzhou@hhhh2")
		phoneNumber = new(OAR1.Person_PhoneNumber)
		phoneNumber.Number = proto.String("13745789524")
		phoneNumber.Type = OAR1.Person_MOBILE.Enum()
		person.Phone = append(person.Phone, phoneNumber)
		args.Person = append(args.Person, person)

		//fmt.Println("@@@@@@@@@@@@@ in pushToClient, before client.Call")
		err := client.Call("OAR1.CheckService.Check", &args, &reply)
		if err != nil {
			log.Println("@@@@@@@@@", err)
			break
		}
		//fmt.Println("@@@@@@@@@@@@@ in pushToClient, after client.Call")
		//log.Println("server push message:", args.GetMsg(), ", and get reply:", reply.GetMsg())
		time.Sleep(3 * time.Second)
	}
}

func Append(slice, data []byte) []byte {
	l := len(slice)
	if l+len(data) > cap(slice) { // reallocate
		// Allocate double what's needed, for future growth.
		newSlice := make([]byte, (l+len(data))*2)
		// The copy function is predeclared and works for any slice type.
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : l+len(data)]
	for i, c := range data {
		slice[l+i] = c
	}
	return slice
}
