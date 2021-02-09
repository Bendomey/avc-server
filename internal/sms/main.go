package sms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Bendomey/avc-server/pkg/utils"
)

var accountSid, authToken, from string

func init() {
	accountSid = utils.MustGet("TWILIO_ACCOUNT_SID")
	authToken = utils.MustGet("TWILIO_AUTH_TOKEN")
	from = utils.MustGet("TWILIO_FROM_NUMBER")
}

//Send text message
func Send(to string, body string) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", from)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
