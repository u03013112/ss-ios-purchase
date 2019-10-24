package main

import (
	"log"
	"net"

	"github.com/u03013112/ss-ios-purchase/ios"
	pb "github.com/u03013112/ss-pb/ios"
	"google.golang.org/grpc"
)

const (
	port = ":50002"
)

func main() {
	ios.InitDB()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listen %s", port)
	s := grpc.NewServer()
	pb.RegisterIOSServer(s, &ios.Srv{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
