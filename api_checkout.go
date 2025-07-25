package client

import (
	"context"
	"fmt"
	"net/http"
)

// OrderInput represents the input for the checkout order call.
type OrderInput struct {
	Vendor                string `json:"vendor"`
	ID                    string `json:"order_id"`
	BuyerEmail            string `json:"buyer_email"`
	BuyerName             string `json:"buyer_name"`
	BuyerUserID           string `json:"buyer_userid,omitempty"`
	BuyerPhone            string `json:"buyer_phone"`
	GatewayBuyerUUID      string `json:"gateway_buyer_uuid"`
	Amount                int    `json:"amount"`
	Currency              string `json:"currency"`
	PaymentMethods        string `json:"payment_methods"`
	RedirectURL           string `json:"redirect_url,omitempty"`
	CancelURL             string `json:"cancel_url,omitempty"`
	Webhook               string `json:"webhook,omitempty"`
	BillingFirstName      string `json:"billing.firstname"`
	BillingLastName       string `json:"billing.lastname"`
	BillingAddress1       string `json:"billing.address_1"`
	BillingAddress2       string `json:"billing.address_2,omitempty"`
	BillingCity           string `json:"billing.city"`
	BillingStateRegion    string `json:"billing.state_or_region"`
	BillingPostCodePOBox  string `json:"billing.postcode_or_pobox"`
	BillingCountry        string `json:"billing.country"`
	BillingPhone          string `json:"billing.phone"`
	ShippingFirstName     string `json:"shipping.firstname,omitempty"`
	ShippingLastName      string `json:"shipping.lastname,omitempty"`
	ShippingAddress1      string `json:"shipping.address_1,omitempty"`
	ShippingAddress2      string `json:"shipping.address_2,omitempty"`
	ShippingCity          string `json:"shipping.city,omitempty"`
	ShippingStateRegion   string `json:"shipping.state_or_region,omitempty"`
	ShippingPostCodePOBox string `json:"shipping.postcode_or_pobox,omitempty"`
	ShippingCountry       string `json:"shipping.country,omitempty"`
	ShippingPhone         string `json:"shipping.phone,omitempty"`
	BuyerRemarks          string `json:"buyer_remarks,omitempty"`
	MerchantRemarks       string `json:"merchant_remarks,omitempty"`
	NumberItems           int    `json:"no_of_items,omitempty"`
	HeaderColour          string `json:"header_colour,omitempty"`
	LinkColour            string `json:"link_colour,omitempty"`
	ButtonColour          string `json:"button_colour,omitempty"`
	Expiry                int    `json:"expiry,omitempty"`
}

// CreateOrder creates a payment order request to the selcom payment gateway.
// Responds with the payment url, buyer details.
func (cln *Client) CreateOrder(ctx context.Context, order OrderInput) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/create-order", cln.host, version)

	// Encoding the webhook string as required.
	if order.Webhook != "" {
		order.Webhook = base64Encode([]byte(order.Webhook))
	}

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, order, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

type OrderInputMinimal struct {
	Vendor          string `json:"vendor"`
	ID              string `json:"order_id"`
	BuyerEmail      string `json:"buyer_email"`
	BuyerName       string `json:"buyer_name"`
	BuyerPhone      string `json:"buyer_phone"`
	Amount          int    `json:"amount"`
	Currency        string `json:"currency"`
	RedirectURL     string `json:"redirect_url,omitempty"`
	CancelURL       string `json:"cancel_url,omitempty"`
	Webhook         string `json:"webhook,omitempty"`
	BuyerRemarks    string `json:"buyer_remarks,omitempty"`
	MerchantRemarks string `json:"merchant_remarks,omitempty"`
	NumberItems     int    `json:"no_of_items,omitempty"`
	HeaderColour    string `json:"header_colour,omitempty"`
	LinkColour      string `json:"link_colour,omitempty"`
	ButtonColour    string `json:"button_colour,omitempty"`
	Expiry          int    `json:"expiry,omitempty"`
}

// CreateOrderMinimal initiates a payment order request to the Selcom payment gateway.
// This method is intended for non-card payments such as mobile wallet push payments
// and manual payment methods. If a webhook URL is provided, it is base64-encoded
// as required by the API.
func (cln *Client) CreateOrderMinimal(ctx context.Context, order OrderInputMinimal) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/create-order-minimal", cln.host, version)

	// Ecoding the webhook string as required.
	if order.Webhook != "" {
		order.Webhook = base64Encode([]byte(order.Webhook))
	}

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, order, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// CancelOrder cancels a payment order before the customer completes the payment.
// It sends a DELETE request to the Selcom checkout API using the provided order ID.
func (cln *Client) CancelOrder(ctx context.Context, orderID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/cancel-order", cln.host, version)

	body := struct {
		OrderID string `json:"order_id"`
	}{
		OrderID: orderID,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodDelete, url, body, &resp); err != nil {
		return Response{}, nil
	}

	return resp, nil
}

// CheckOrder retrieves the current status of a payment order from the Selcom checkout API.
// It sends a GET request using the provided order ID.
func (cln *Client) CheckOrder(ctx context.Context, orderID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/order-status?order_id=%s", cln.host, version, orderID)

	body := struct {
		OrderID string `json:"order_id"`
	}{
		OrderID: orderID,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodGet, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// Orders retrieves a list of orders from the Selcom checkout API within the specified date range.
// It uses the provided start and end dates to filter the results.
func (cln *Client) Orders(ctx context.Context, startDate string, endDate string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/list-orders?fromdate=%s&todate=%s", cln.host, version, startDate, endDate)

	body := struct {
		From string `json:"fromdate"`
		To   string `json:"todate"`
	}{
		From: startDate,
		To:   endDate,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodGet, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// FetchStoredCards retrieves the stored billing cards for the specified buyer.
// The gatewayBuyerUUID is a unique identifier generated for the user upon their first order creation.
func (cln *Client) FetchStoredCards(ctx context.Context, buyerUserID string, gatewayBuyerUUID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/stored-cards", cln.host, version)

	body := struct {
		BuyerUserID      string `json:"buyer_userid"`
		GatewayBuyerUUID string `json:"gateway_buyer_uuid"`
	}{
		BuyerUserID:      buyerUserID,
		GatewayBuyerUUID: gatewayBuyerUUID,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, body, &resp); err != nil {
		return Response{}, nil
	}

	return resp, nil
}

// DeleteStoredCard removes a previously saved billing card from the Selcom gateway.
// It uses the card's resource ID and the associated gatewayBuyerUUID to identify the card.
func (cln *Client) DeleteStoredCard(ctx context.Context, cardResourceID string, gatewayBuyerUUID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/delete-card?id=%s&gateway_buyer_uuid=%s", cln.host, version, cardResourceID, gatewayBuyerUUID)

	body := struct {
		ID               string `json:"id"`
		GatewayBuyerUUID string `json:"gateway_buyer_uuid"`
	}{
		ID:               cardResourceID,
		GatewayBuyerUUID: gatewayBuyerUUID,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodDelete, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// CardPaymentInput represents the payload required to initiate a card payment
// using a stored card through the Selcom payment gateway.
type CardPaymentInput struct {
	TransactionID    string `json:"transid"`
	Vendor           string `json:"vendor"`
	OrderID          string `json:"order_id"`
	CardToken        string `json:"card_token"`
	BuyerUserID      string `json:"buyer_userid"`
	GatewayBuyerUUID string `json:"gateway_buyer_uuid"`
}

// CardPayment processes an order using a stored billing card without redirecting
// the user to the Selcom payment gateway page.
// This is intended for seamless in-app or on-site card payments.
func (cln *Client) CardPayment(ctx context.Context, cardPaymentInput CardPaymentInput) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/card-payment", cln.host, version)

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, cardPaymentInput, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// WalletPayment initiates a mobile wallet payment for the specified order.
// This API triggers a USSD push (STK push) to the customer's phone to authorize the payment.
// It is intended for seamless integration of mobile wallet payments (e.g., M-Pesa, Tigo Pesa).

func (cln *Client) WalletPayment(ctx context.Context, transactionID string, orderID string, phone string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/wallet-payment", cln.host, version)

	body := struct {
		TransactionID string `json:"transid"`
		OrderID       string `json:"order_id"`
		MSISDN        string `json:"msisdn"`
	}{
		TransactionID: transactionID,
		OrderID:       orderID,
		MSISDN:        phone,
	}

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, body, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}
