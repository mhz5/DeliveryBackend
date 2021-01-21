package main

import (
	"context"
	"crypto/tls"
	"strings"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	cpb "fda/proto/common"
	lpb "fda/proto/location"
	ppb "fda/proto/payment"
	rpb "fda/proto/restaurant"
)

const defaultName = "world"
const localHost = "127.0.0.1:50051"

var (
	addr     = flag.String("addr", localHost, "Address of grpc server.")
	key      = flag.String("api-key", "", "API key.")
	token    = flag.String("token", "", "Authentication token.")
	keyfile  = flag.String("keyfile", "", "Path to a Google service account key file.")
	audience = flag.String("audience", "", "Audience.")
)

func sendGetRestaurantsGrpc(conn *grpc.ClientConn, ctx context.Context, header metadata.MD) {
	restaurantClient := rpb.NewRestaurantServiceClient(conn)
	r, err := restaurantClient.GetRestaurants(ctx, &rpb.GetRestaurantsRequest{}, grpc.Header(&header))
	if err != nil {
		log.Fatalf("could not send GetRestaurants request: %v", err)
		return
	}
	log.Printf("GetRestaurants Response: %s", r.Restaurants)
}

func sendCreatePaymentGrpc(conn *grpc.ClientConn, ctx context.Context, header metadata.MD)  {
	paymentClient := ppb.NewPaymentServiceClient(conn)
	req := &ppb.CreatePaymentRequest{
		Nonce: "cnon:CBASEPUWKIy69prEnYsnOux00FM",
		BuyerVerificationToken: "verf:CBASEPCoIgjqrfBWGaJ7vFULqT4",
		AmountCents: 123,
		Currency: "USD",
	}

	r, err := paymentClient.CreatePayment(ctx, req)
	if err != nil {
		log.Printf("could not send CreatePayment request: %v", err)
		return
	}
	log.Printf("CreatePayment Response: %s", r)
}

func sendPublishLocationGrpc(conn *grpc.ClientConn, ctx context.Context, header metadata.MD) {
	locationClient := lpb.NewLocationServiceClient(conn)
	req := &lpb.PublishLocationRequest{
		UserToken: "snth1234",
		Location: &cpb.Point{
			Lat: 37.4219983,
			Long: -122.084,
		},
	}
	r, err := locationClient.PublishLocation(ctx, req)
	if err != nil {
		log.Fatalf("could not send PublishLocation request: %v", err)
		return
	}

	log.Printf("PublishLocation Response: %s", r)
}

func main() {
	flag.Parse()

	// Set up a connection to the server.
	tlsDialOption := grpc.WithInsecure()
	// Use fake TLS if using HTTPS.
	if strings.Contains(*addr, "443") {
		config := &tls.Config{ InsecureSkipVerify: true }
		tlsCred := credentials.NewTLS(config)
		tlsDialOption = grpc.WithTransportCredentials(tlsCred)
	}
	conn, err := grpc.Dial(*addr, tlsDialOption)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	if *key != "" {
		log.Printf("Using API key: %s", *key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", *key)
	}
	if *token != "" {
		log.Printf("Using authentication token: %s", *token)
		ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	var header metadata.MD
	log.Printf("==== GetRestaurants ====")
	sendGetRestaurantsGrpc(conn, ctx, header)
	log.Printf("==== CreatePayment ====")
	sendCreatePaymentGrpc(conn, ctx, header)
	log.Printf("==== PublishLocation ====")
	sendPublishLocationGrpc(conn, ctx, header)
}
