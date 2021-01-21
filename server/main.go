package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	lpb "fda/proto/location"
	ppb "fda/proto/payment"
	rpb "fda/proto/restaurant"
	"fda/server/location"
	"fda/server/payment"
	"fda/server/restaurant"
)

var (
	addr = flag.String("addr", ":50051", "Network host:port to listen on for gRPC connections.")
)

// TODO: recover from panic.
func main() {
	flag.Parse()
	log.Print("Starting server at: " + *addr)

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	rpb.RegisterRestaurantServiceServer(s, &restaurant.Server{})
	ppb.RegisterPaymentServiceServer(s, &payment.Server{})
	lpb.RegisterLocationServiceServer(s, &location.Server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
