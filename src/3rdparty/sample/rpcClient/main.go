/*
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Args struct {
	A, B int
}
type Quotient struct {
	Quo *int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		os.Exit(1)
	}

	service := os.Args[1]
	client, err := rpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &Args{17, 8}

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("#################get: address:%0x, ", quot.Quo)
	if quot.Quo != nil {
		fmt.Printf("value:%d", *quot.Quo)
	}
	fmt.Printf("\n")
}
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

type AddressCheck struct {
	TypeInt
}

type TypeInt int

/*
func (t *AddressCheck) Check(args *OAR1.AddressBook, reply *OAR1.CheckResponse) error {
	log.Printf("AddressCheck, client get pushed address book!============================")
	return nil
}*/

func (t *TypeInt) Check(args *OAR1.AddressBook, reply *OAR1.CheckResponse) error {
	log.Printf("client get pushed address book!============================")
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
	err := ConnectToServer("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("ConnectToServer: %v", err)
	}
}

func main() {
	time.Sleep(time.Hour)
}

func ConnectToServer(network, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	conn, err := net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		return err
	}

	srv := rpc.NewServer()
	if err := srv.Register(new(AddressCheck)); err != nil {
		return err
	}

	go srv.ServeCodec(protorpc.NewServerCodec(conn, nil, false))
	go CallServer(conn)

	return nil
}

func CallServer(conn net.Conn) {
	client := rpc.NewClientWithCodec(protorpc.NewClientCodec(conn, false))
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
		err := client.Call("CheckService.Check", &args, &reply)
		if err != nil {
			log.Println("@@@@@@@@@", err)
			break
		}
		//fmt.Println("@@@@@@@@@@@@@ in pushToClient, after client.Call")
		//log.Println("server push message:", args.GetMsg(), ", and get reply:", reply.GetMsg())
		time.Sleep(2 * time.Second)
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
