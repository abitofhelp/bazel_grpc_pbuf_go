package main

import (
	"context"
	"fmt"
	"github.com/abitofhelp/bazel_grpc_pbuf_go/cmd"
	v1 "github.com/abitofhelp/bazel_grpc_pbuf_go/proto/helloworld/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func main() {

	// Build the message to send.
	req := buildRequest("Mike")

	// Connect to the gRPC server.
	serverAddress := fmt.Sprintf("%s:%s", cmd.Host, cmd.Port)

	// Timeout for initial client connection.
	ctx, cancel := context.WithTimeout(context.Background(), cmd.ConnectionTimeOut)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to the gRPC server at '%s': %v", serverAddress, err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close the connection to the gRPC server: %v", err)
		}
	}(conn)

	reqctx, cancel := context.WithTimeout(context.Background(), cmd.RequestTimeOut)
	defer cancel()

	// Contact the server and print out the response.
	c := v1.NewGreeterClient(conn)
	if rsp, err := c.SayHello(reqctx, &req); err == nil {
		now := time.Now().UTC()
		requestDuration := rsp.MessageArrivedUtc.AsTime().Sub(req.SentUtc.AsTime())
		responseDuration := now.Sub(rsp.RepliedUtc.AsTime())
		totalDuration := now.Sub(req.SentUtc.AsTime())

		log.Printf("%s %s\n %s: %dμs\n %s: %dμs\n %s: %dμs\n",
			rsp.Message,
			"Here are some details about your request:",
			"Request Duration",
			requestDuration.Microseconds(),
			"Response Duration",
			responseDuration.Microseconds(),
			"Total Duration",
			totalDuration.Microseconds(),
		)
	} else {
		log.Fatalf("failed to say hello: %v", err)
	}
}

func buildRequest(name string) v1.HelloRequest {
	now := time.Now().UTC()
	timestamp := timestamppb.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.Nanosecond()),
	}

	return v1.HelloRequest{
		Name:    name,
		SentUtc: &timestamp,
	}
}
