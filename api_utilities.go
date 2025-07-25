package client

import (
	"context"
	"fmt"
	"net/http"
)

// UtilityPaymentInput represents the data required to process a utility payment
// through the Selcom payment gateway.
type UtilityPaymentInput struct {
	TransactionID    string  `json:"transid"`
	UtilityCode      string  `json:"utilitycode"`
	UtilityReference string  `json:"utilityref"`
	Amount           float64 `json:"amount"`
	Vendor           string  `json:"vendor"`
	Pin              string  `json:"pin"`
	Phone            string  `json:"msisdn"`
}

// UtilityPayment process payment for a particular payment service.// UtilityPayment processes a payment for a specified utility service.
// The payment details are provided in the UtilityPaymentInput parameter.
func (cln *Client) UtilityPayment(ctx context.Context, body UtilityPaymentInput) (Response, error) {
	url := fmt.Sprintf("%s/%s/utilitypayment/process", cln.host, version)

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// UtilityLookup retrieves information related to a specific utility payment.
// It queries the utility payment service using the utility code, utility reference,
// and transaction ID to validate or fetch details before processing a payment.

func (cln *Client) UtilityLookup(ctx context.Context, utilityCode, utilityRef, transactionID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/utilitypayment/lookup?utilitycode=%s&utilityref=%s&transid=%s", cln.host, version, utilityCode, utilityRef, transactionID)

	var body = struct {
		TransactionID    string `json:"transid"`
		UtilityCode      string `json:"utilitycode"`
		UtilityReference string `json:"utilref"`
	}{
		TransactionID:    transactionID,
		UtilityCode:      utilityCode,
		UtilityReference: utilityRef,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodGet, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// UtilityPaymentStatus checks the current status of a utility payment
// by querying the utility payment service with the provided transaction ID.
func (cln *Client) UtilityPaymentStatus(ctx context.Context, trasactionID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/utilitypayment/query?transid=%s", cln.host, version, trasactionID)

	var body = struct {
		TransactionID string `json:"transid"`
	}{
		TransactionID: trasactionID,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodGet, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}
