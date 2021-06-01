package utils

import (
	"context"
	"fmt"

	// "github.com/Bendomey/avc-server/pkg/paystack"
	"github.com/kehindesalaam/go-paystack/paystack"
)

var apiKey string
var client *paystack.Client

func init() {
	apiKey = MustGet("PAYMENT_TEST_API_KEY")

	client = paystack.NewClient(nil, paystack.SecretKey(apiKey))
}

func InitializePayment(context context.Context, req paystack.TransactionRequest) (*paystack.Response, error) {
	fmt.Print(apiKey, client)
	auth, transaction, err := client.Transaction.Initialize(context, &req)
	fmt.Print("transaction authorization", auth)
	return transaction, err
}

func VerifyPayment(reference string) (*paystack.Response, error) {
	verify, res, err := client.Transaction.Verify(context.TODO(), reference)
	fmt.Print("transaction authorization", verify)
	return res, err
}
