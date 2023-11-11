package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/abitofhelp/bazel_grpc_pbuf_go/cmd"
	v1 "github.com/abitofhelp/bazel_grpc_pbuf_go/proto/helloworld/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

// server implements the GreeterServer.
type server struct {
	v1.UnimplementedGreeterServer
}

// SayHello implements the RPC function to say hello to a requestor.
func (s *server) SayHello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloResponse, error) {
	requestArrivedUtc := time.Now().UTC()
	log.Printf("Received: %v", req)
	return buildResponse(req.Name, requestArrivedUtc), nil
}

func timeoutInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var err error
	var result interface{}

	done := make(chan struct{})

	go func() {
		result, err = handler(ctx, req)
		done <- struct{}{}
	}()

	// Has the client set a deadline in the request?
	dt, ok := ctx.Deadline()
	if !ok {
		msg := fmt.Sprintf("client's context did not have a deadline: %v", ctx.Err())
		fmt.Println(msg)
		return nil, status.New(codes.FailedPrecondition, msg).Err()
	}
	// Determine the remaining time before the client want's a timeout.
	timeRemainingBeforeTimeout := dt.Sub(time.Now())

	select {
	case <-ctx.Done():
		msg := fmt.Sprintf("client cancelled the request: %v", ctx.Err())
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Println(msg)
			return nil, status.New(codes.Canceled, msg).Err()
		}
	case <-time.After(timeRemainingBeforeTimeout):
		msg := fmt.Sprintf("the client's deadline timeout has expired: %v", ctx.Err())
		fmt.Println(msg)
		return nil, status.New(codes.DeadlineExceeded, msg).Err()
	case <-done:
	}
	return result, err
}

func main() {
	// Create the listener for the gRPC server.
	if lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cmd.Port)); err == nil {
		// Configure the gRPC server.
		s := grpc.NewServer(grpc.ChainUnaryInterceptor(timeoutInterceptor))
		v1.RegisterGreeterServer(s, &server{})

		// Launch the server.
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve a connection from the listener: %v", err)
		}
	} else {
		log.Fatalf("failed to listen for our network traffic: %v", err)
	}
}

func buildResponse(name string, requestArrivedUtc time.Time) *v1.HelloResponse {
	messageArrivedTimestamp := timestamppb.Timestamp{
		Seconds: requestArrivedUtc.Unix(),
		Nanos:   int32(requestArrivedUtc.Nanosecond()),
	}

	repliedUtc := time.Now().UTC()
	repliedTimestamp := timestamppb.Timestamp{
		Seconds: repliedUtc.Unix(),
		Nanos:   int32(repliedUtc.Nanosecond()),
	}

	return &v1.HelloResponse{
		Message:           fmt.Sprintf("Howdy %s!", name),
		MessageArrivedUtc: &messageArrivedTimestamp,
		RepliedUtc:        &repliedTimestamp,
	}
}
