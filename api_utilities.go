package client

import (
	"context"
	"fmt"
	"net/http"
)

type UtilityPaymentInput struct {
	TransactionID    string  `json:"transid"`
	UtilityCode      string  `json:"utilitycode"`
	UtilityReference string  `json:"utilityref"`
	Amount           float64 `json:"amount"`
	Vendor           string  `json:"vendor"`
	Pin              string  `json:"pin"`
	Phone            string  `json:"msisdn"`
}

// UtilityPayment process payment for a particular payment service.
func (cln *Client) UtilityPayment(ctx context.Context, body UtilityPaymentInput) (Response, error) {
	url := fmt.Sprintf("%s/%s/utilitypayment/process", cln.host, version)

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

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

// UtilityPaymentStatus checks ths status of the utility payment.
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
