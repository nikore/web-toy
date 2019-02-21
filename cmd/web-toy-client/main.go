package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "web-toy/pkg/helloworld"

	"gopkg.in/alecthomas/kingpin.v2"
)

func makeGRPC(address string, msg string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: msg})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate)
	host := kingpin.Flag("host", "grpc host to connect too").Short('h').Default("127.0.0.1").String()
	port := kingpin.Flag("port", "grpc port to connect too").Short('p').Default("8080").Int()
	msg := kingpin.Arg("name", "name to say hello too with grpc").Default("world").String()
	kingpin.Parse()

	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Using GRPC Server: %s\n", address)

	makeGRPC(address, *msg)
}
