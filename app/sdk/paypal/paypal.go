package paypal

import (
	"context"
	"log"

	"github.com/plutov/paypal/v4"
)

// PayPalClient holds the PayPal SDK client
type PayPalClient struct {
	Client    *paypal.Client
	ClientURL string
}

// NewPayPalClient initializes and returns a new PayPal client
func NewPayPalClient(clientID string, clientSecret string, url string) (*PayPalClient, error) {
	//clientID := os.Getenv("PAYPAL_CLIENT_ID")
	//clientSecret := os.Getenv("PAYPAL_SECRET_ID")

	if clientID == "" || clientSecret == "" {
		log.Fatal("PayPal credentials are missing")
	}

	client, err := paypal.NewClient(clientID, clientSecret, paypal.APIBaseSandBox)
	if err != nil {
		return nil, err
	}

	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}

	return &PayPalClient{Client: client, ClientURL: url}, nil
}

func (p *PayPalClient) CreateOrder(ctx context.Context, amount string, currency string) (*paypal.Order, error) {
	order, err := p.Client.CreateOrder(ctx, paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{{
		Amount: &paypal.PurchaseUnitAmount{
			Currency: currency,
			Value:    amount,
		},
	},
	},
		nil, // No PaymentSource
		&paypal.ApplicationContext{
			ReturnURL: p.ClientURL + "/payment-return",
			CancelURL: p.ClientURL + "/payment-cancel",
		},
	)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// CaptureOrder captures a PayPal order
func (p *PayPalClient) CaptureOrder(orderID string) (*paypal.CaptureOrderResponse, error) {
	return p.Client.CaptureOrder(context.Background(), orderID, paypal.CaptureOrderRequest{})
}
