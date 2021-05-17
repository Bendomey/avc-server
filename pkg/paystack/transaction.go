package paystack

import "fmt"

// TransactionService handles operations related to transactions
// For more details see https://developers.paystack.co/v1.0/reference#create-transaction
type TransactionService service

// TransactionList is a list object for transactions.
type TransactionList struct {
	Meta   ListMeta
	Values []Transaction `json:"data"`
}

// TransactionRequest represents a request to start a transaction.
type TransactionRequest struct {
	CallbackURL       string   `json:"callback_url,omitempty"`
	Reference         string   `json:"reference,omitempty"`
	AuthorizationCode string   `json:"authorization_code,omitempty"`
	Currency          string   `json:"currency,omitempty"`
	Amount            float32  `json:"amount,omitempty"`
	Email             string   `json:"email,omitempty"`
	Plan              string   `json:"plan,omitempty"`
	InvoiceLimit      int      `json:"invoice_limit,omitempty"`
	Metadata          Metadata `json:"metadata,omitempty"`
	SubAccount        string   `json:"subaccount,omitempty"`
	TransactionCharge int      `json:"transaction_charge,omitempty"`
	Bearer            string   `json:"bearer,omitempty"`
	Channels          []string `json:"channels,omitempty"`
}

// AuthorizationRequest represents a request to enable/revoke an authorization
type AuthorizationRequest struct {
	Reference         string   `json:"reference,omitempty"`
	AuthorizationCode string   `json:"authorization_code,omitempty"`
	Amount            int      `json:"amount,omitempty"`
	Currency          string   `json:"currency,omitempty"`
	Email             string   `json:"email,omitempty"`
	Metadata          Metadata `json:"metadata,omitempty"`
}

// Transaction is the resource representing your Paystack transaction.
// For more details see https://developers.paystack.co/v1.0/reference#initialize-a-transaction
type Transaction struct {
	ID              int                    `json:"id,omitempty"`
	CreatedAt       string                 `json:"createdAt,omitempty"`
	Domain          string                 `json:"domain,omitempty"`
	Metadata        string                 `json:"metadata,omitempty"` //TODO: why is transaction metadata a string?
	Status          string                 `json:"status,omitempty"`
	Reference       string                 `json:"reference,omitempty"`
	Amount          float32                `json:"amount,omitempty"`
	Message         string                 `json:"message,omitempty"`
	GatewayResponse string                 `json:"gateway_response,omitempty"`
	PaidAt          string                 `json:"piad_at,omitempty"`
	Channel         string                 `json:"channel,omitempty"`
	Currency        string                 `json:"currency,omitempty"`
	IPAddress       string                 `json:"ip_address,omitempty"`
	Log             map[string]interface{} `json:"log,omitempty"` // TODO: same as timeline?
	Fees            int                    `json:"int,omitempty"`
	FeesSplit       string                 `json:"fees_split,omitempty"` // TODO: confirm data type
	Authorization   Authorization          `json:"authorization,omitempty"`
}

// Authorization represents Paystack authorization object
type Authorization struct {
	AuthorizationCode string `json:"authorization_code,omitempty"`
	Bin               string `json:"bin,omitempty"`
	Last4             string `json:"last4,omitempty"`
	ExpMonth          string `json:"exp_month,omitempty"`
	ExpYear           string `json:"exp_year,omitempty"`
	Channel           string `json:"channel,omitempty"`
	CardType          string `json:"card_type,omitempty"`
	Bank              string `json:"bank,omitempty"`
	CountryCode       string `json:"country_code,omitempty"`
	Brand             string `json:"brand,omitempty"`
	Resusable         bool   `json:"reusable,omitempty"`
	Signature         string `json:"signature,omitempty"`
}

// TransactionTimeline represents a timeline of events in a transaction session
type TransactionTimeline struct {
	TimeSpent      int                      `json:"time_spent,omitempty"`
	Attempts       int                      `json:"attempts,omitempty"`
	Authentication string                   `json:"authentication,omitempty"` // TODO: confirm type
	Errors         int                      `json:"errors,omitempty"`
	Success        bool                     `json:"success,omitempty"`
	Mobile         bool                     `json:"mobile,omitempty"`
	Input          []string                 `json:"input,omitempty"` // TODO: confirm type
	Channel        string                   `json:"channel,omitempty"`
	History        []map[string]interface{} `json:"history,omitempty"`
}

// Initialize initiates a transaction process
// For more details see https://developers.paystack.co/v1.0/reference#initialize-a-transaction
func (s *TransactionService) Initialize(txn *TransactionRequest) (Response, error) {
	u := fmt.Sprintf("/transaction/initialize")
	resp := Response{}
	err := s.client.Call("POST", u, txn, &resp)
	return resp, err
}

// Verify checks that transaction with the given reference exists
// For more details see https://api.paystack.co/transaction/verify/reference
func (s *TransactionService) Verify(reference string) (*Transaction, error) {
	u := fmt.Sprintf("/transaction/verify/%s", reference)
	txn := &Transaction{}
	err := s.client.Call("GET", u, nil, txn)
	return txn, err
}

// Get returns the details of a transaction.
// For more details see https://developers.paystack.co/v1.0/reference#fetch-transaction
func (s *TransactionService) Get(id int) (*Transaction, error) {
	u := fmt.Sprintf("/transaction/%d", id)
	txn := &Transaction{}
	err := s.client.Call("GET", u, nil, txn)
	return txn, err
}
