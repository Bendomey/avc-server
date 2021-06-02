package utils

import (
	"context"
	"fmt"

	"github.com/kehindesalaam/go-paystack/paystack"
)

var apiKey string
var client *paystack.Client

func init() {
	apiKey = MustGet("PAYMENT_TEST_API_KEY")

	client = paystack.NewClient(nil, paystack.SecretKey(apiKey))
}

func InitializePayment(context context.Context, req paystack.TransactionRequest) (*paystack.TransactionAuthorization, error) {
	fmt.Print(apiKey, client)
	auth, _, err := client.Transaction.Initialize(context, &req)
	return auth, err
}

func VerifyPayment(reference string) (*paystack.Response, error) {
	verify, res, err := client.Transaction.Verify(context.TODO(), reference)
	fmt.Print("transaction authorization", verify)
	return res, err
}
