package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// OrderInput represents the inputfor the checkout order call.
type OrderInput struct {
	Vendor                string    `json:"vendor"`
	ID                    uuid.UUID `json:"order_id"`
	BuyerEmail            string    `json:"buyer_email"`
	BuyerName             string    `json:"buyer_name"`
	BuyerUserID           string    `json:"buyer_userid,omitempty"`
	BuyerPhone            string    `json:"buyer_phone"`
	GatewayBuyerUUID      string    `json:"gateway_buyer_uuid"`
	Amount                int       `json:"amount"`
	Currency              string    `json:"currency"`
	PaymentMethods        string    `json:"payment_methods"`
	RedirectURL           string    `json:"redirect_url,omitempty"`
	CancelURL             string    `json:"cancel_url,omitempty"`
	Webhook               string    `json:"webhook,omitempty"`
	BillingFirstName      string    `json:"billing.firstname"`
	BillingLastName       string    `json:"billing.lastname"`
	BillingAddress1       string    `json:"billing.address_1"`
	BillingAddress2       string    `json:"billing.address_2,omitempty"`
	BillingCity           string    `json:"billing.city"`
	BillingStateRegion    string    `json:"billing.state_or_region"`
	BillingPostCodePOBox  string    `json:"billing.postcode_or_pobox"`
	BillingCountry        string    `json:"billing.country"`
	BillingPhone          string    `json:"billing.phone"`
	ShippingFirstName     string    `json:"shipping.firstname,omitempty"`
	ShippingLastName      string    `json:"shipping.lastname,omitempty"`
	ShippingAddress1      string    `json:"shipping.address_1,omitempty"`
	ShippingAddress2      string    `json:"shipping.address_2,omitempty"`
	ShippingCity          string    `json:"shipping.city,omitempty"`
	ShippingStateRegion   string    `json:"shipping.state_or_region,omitempty"`
	ShippingPostCodePOBox string    `json:"shipping.postcode_or_pobox,omitempty"`
	ShippingCountry       string    `json:"shipping.country,omitempty"`
	ShippingPhone         string    `json:"shipping.phone,omitempty"`
	BuyerRemarks          string    `json:"buyer_remarks,omitempty"`
	MerchantRemarks       string    `json:"merchant_remarks,omitempty"`
	NumberItems           int       `json:"no_of_items,omitempty"`
	HeaderColour          string    `json:"header_colour,omitempty"`
	LinkColour            string    `json:"link_colour,omitempty"`
	ButtonColour          string    `json:"button_colour,omitempty"`
	Expiry                int       `json:"expiry,omitempty"`
}

// CreateOrder creates a payment order request to the selcom payment gateway.
// Responds with the payment url, buyer details.
func (cln *Client) CreateOrder(ctx context.Context, order OrderInput) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/create-order", cln.host, version)

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, order, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

type OrderInputMinimal struct {
	Vendor          string    `json:"vendor"`
	ID              uuid.UUID `json:"order_id"`
	BuyerEmail      string    `json:"buyer_email"`
	BuyerName       string    `json:"buyer_name"`
	BuyerPhone      string    `json:"buyer_phone"`
	Amount          int       `json:"amount"`
	Currency        string    `json:"currency"`
	RedirectURL     string    `json:"redirect_url,omitempty"`
	CancelURL       string    `json:"cancel_url,omitempty"`
	Webhook         string    `json:"webhook,omitempty"`
	BuyerRemarks    string    `json:"buyer_remarks,omitempty"`
	MerchantRemarks string    `json:"merchant_remarks,omitempty"`
	NumberItems     int       `json:"no_of_items,omitempty"`
	HeaderColour    string    `json:"header_colour,omitempty"`
	LinkColour      string    `json:"link_colour,omitempty"`
	ButtonColour    string    `json:"button_colour,omitempty"`
	Expiry          int       `json:"expiry,omitempty"`
}

// CreateOrderMinimal creates a payment order request to the selcom payment gateway.
// This is for non-card payments. Ideal for mobile wallet push payments and manual payments.
func (cln *Client) CreateOrderMinimal(ctx context.Context, order OrderInputMinimal) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/create-order-minimal", cln.host, version)

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, order, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// CancelOrder Cancels an order before customer completes the payment.
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

// CheckOrder returns status of the order.
func (cln *Client) CheckOrder(ctx context.Context, orderID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/order-status?order_id=%s", cln.host, version, orderID)

	var resp Response
	if err := cln.do(ctx, http.MethodGet, url, nil, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// Orders returns a list of orders, satisfying the startDate and endDate params.
func (cln *Client) Orders(ctx context.Context, startDate string, endDate string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/list-orders?fromdate=%s&todate=%s", cln.host, version, startDate, endDate)

	var resp Response
	if err := cln.do(ctx, http.MethodGet, url, nil, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// FetchStoredCards retursn the stored billing cards for the provided buyer.
// The gatewayBuyerUUID is generated for each user on their first order creation.
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

// DeleteStoredCard deletes the provided billing card informations.
func (cln *Client) DeleteStoredCard(ctx context.Context, cardResourceID string, gatewayBuyerUUID string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/delete-card?id=%s&gateway_buyer_uuid=%s", cln.host, version, cardResourceID, gatewayBuyerUUID)

	var resp Response
	if err := cln.do(ctx, http.MethodDelete, url, nil, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

type CardPaymentInput struct {
	TransactionID    string `json:"transid"`
	Vendor           string `json:"vendor"`
	OrderID          string `json:"order_id"`
	CardToken        string `json:"card_token"`
	BuyerUserID      string `json:"buyer_userid"`
	GatewayBuyerUUID string `json:"gateway_buyer_uuid	"`
}

// ProcessCardPayment allows the ecommerce website to process an order using stored cards directly
// without redirecting the user to payment gateway page.
func (cln *Client) ProcessCardPayment(ctx context.Context, cardPaymentInput CardPaymentInput) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/card-payment", cln.host, version)

	var resp Response
	if err := cln.do(ctx, http.MethodPost, url, cardPaymentInput, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}

// ProcessWalletPayment  api allows the ecommerce website to process an order using mobile wallets directly.
// trigger this api call to reiceve a PUSH ussd.
func (cln *Client) ProcessWalletPayment(ctx context.Context, transactionID uuid.UUID, orderID uuid.UUID, phone string) (Response, error) {
	url := fmt.Sprintf("%s/%s/checkout/wallet-payment", cln.host, version)

	body := struct {
		TransactionID uuid.UUID `json:"transid"`
		OrderID       uuid.UUID `json:"order_id"`
		MSISDN        string    `json:"msisdn"`
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