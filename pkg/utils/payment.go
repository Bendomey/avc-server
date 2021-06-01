package utils

import (
	"fmt"

	"github.com/Bendomey/avc-server/pkg/paystack"
)

var apiKey string
var client *paystack.Client

func init() {
	apiKey = MustGet("PAYMENT_TEST_API_KEY")

	client = paystack.NewClient(apiKey, nil)
}

func InitializePayment(req paystack.TransactionRequest) (paystack.Response, error) {
	fmt.Print(apiKey, client)
	transaction, err := client.Transaction.Initialize(&req)
	return transaction, err
}

func VerifyPayment(reference string) (*paystack.Transaction, error) {
	verify, err := client.Transaction.Verify(reference)
	return verify, err
}
