package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"time"
	pb "web-toy/pkg/helloworld"
)

type grpcResult struct {
	len       int
	timeTaken float64
	succeed   bool
}

func makeGRPC(address string, duration int, ch chan<- grpcResult) {
	time.Sleep(time.Duration(rand.Intn(duration)) * time.Second)
	var failed bool
	var result grpcResult
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	start := time.Now()
	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "JACKagag"})
	if err != nil {
		failed = true
	} else {
		failed = false
	}
	secs := time.Since(start).Seconds()
	if failed {
		result = grpcResult{
			len:       0,
			timeTaken: secs,
			succeed:   false,
		}
	} else {
		result = grpcResult{
			len:       len(r.Message),
			timeTaken: secs,
			succeed:   true,
		}
	}
	ch <- result
}

func main() {
	host := flag.String("host", "localhost", "hostname or ip of grpc server")
	port := flag.Int("port", 8080, "port of grpc server")
	num_calls := flag.Int("num", 10, "number of grpc calls")
	duration := flag.Int("duration", 30, "time period to run test")
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Using GRPC Server: %s\n", address)
	total_success := 0
	total_failure := 0
	ch := make(chan grpcResult)
	for i := 1; i <= *num_calls; i++ {
		go makeGRPC(address, *duration, ch)
	}
	for i := 1; i <= *num_calls; i++ {
		result := <-ch
		if result.succeed {
			total_success += 1
		} else {
			total_failure += 1
		}
		// fmt.Printf("success: %t, len: %d in: %f\n", result.succeed, result.len, result.timeTaken)
	}
	fmt.Printf("Total succeed: %d | total failure: %d \n", total_success, total_failure)
	suc_per := (float64(total_success) / float64(*num_calls)) * 100.0
	fmt.Printf("Succeed: %.2f%s\n", suc_per, "%")
	fail_per := (float64(total_failure) / float64(*num_calls)) * 100.0
	fmt.Printf("Failure: %.2f%s\n", fail_per, "%")
}
