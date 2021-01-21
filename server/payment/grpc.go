package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	fpb "fda/proto/payment"
)

type Server struct{}

type createPaymentData struct {
	SourceId          string       `json:"source_id"`
	VerificationToken string       `json:"verification_token"`
	Autocomplete      bool         `json:"autocomplete"`
	LocationId        string       `json:"location_id"`
	IdempotencyKey    string       `json:"idempotency_key"`
	AmountMoney       *moneyAmount `json:"amount_money"`
}

type moneyAmount struct {
	AmountCents int32  `json:"amount"`
	Currency    string `json:"currency"`
}

func (s *Server) CreatePayment(ctx context.Context, in *fpb.CreatePaymentRequest) (*fpb.CreatePaymentResponse, error) {
	log.Printf("Handling CreatePayment request [%v] with context %v", in, ctx)
	createPaymentUrl := "https://connect.squareupsandbox.com/v2/payments"

	data := &createPaymentData{
		SourceId:          in.Nonce,
		VerificationToken: in.BuyerVerificationToken,
		Autocomplete:      true,
		LocationId:        "LGDEVVWVR4BF6",
		IdempotencyKey:    uuid.New().String(),
		AmountMoney:       &moneyAmount{AmountCents: in.AmountCents, Currency: "USD"},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return &fpb.CreatePaymentResponse{Success: false}, fmt.Errorf("could not serialize payload to json: %v", err)
	}

	req, err :=
		http.NewRequest("POST",
			createPaymentUrl,
			strings.NewReader(string(jsonData)))
	if err != nil {
		return &fpb.CreatePaymentResponse{Success: false}, fmt.Errorf("could not create createPayment http request: %v", err)
	}
	req.Header.Add("Authorization",
		"Bearer EAAAEOK5LqkFWsGfcWxCS81GrapCH1pw-b-9aaLxwRJhGZMbzb28JMr55nsTsiFe")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &fpb.CreatePaymentResponse{Success: false}, fmt.Errorf("could not send Square createPayment request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)
		return &fpb.CreatePaymentResponse{Success: false},
			fmt.Errorf("http error: %s", string(b))
	}

	return &fpb.CreatePaymentResponse{Success: true}, nil
}
