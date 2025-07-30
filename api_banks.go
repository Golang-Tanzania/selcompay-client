package client

import (
	"context"
	"fmt"
	"net/http"
)

// BankLookup retrieves the account information for a given bank and account number.
// It performs a GET request to the QwikSend lookup endpoint using the provided bank code,
// account number, and transaction ID.
func (cln *Client) BankLookup(ctx context.Context, bank, account, transactionId string) (Response, error) {
	url := fmt.Sprintf("%s/%s/qwiksend/lookup?bank=%s&account=%s&transid=%s", cln.host, version, bank, account, transactionId)

	var body = struct {
		TransactionID string `json:"transid"`
		Bank          string `json:"bank"`
		Account       string `json:"account"`
	}{
		TransactionID: transactionId,
		Bank:          bank,
		Account:       account,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodGet, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}
