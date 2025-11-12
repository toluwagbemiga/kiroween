package internal

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	portalsession "github.com/stripe/stripe-go/v76/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeClient wraps all Stripe SDK operations
type StripeClient struct {
	apiKey string
}

// NewStripeClient creates a new Stripe client
func NewStripeClient(apiKey string) *StripeClient {
	stripe.Key = apiKey
	return &StripeClient{
		apiKey: apiKey,
	}
}

// Product Operations

// CreateProduct creates a Stripe product
func (c *StripeClient) CreateProduct(name string, metadata map[string]string) (*stripe.Product, error) {
	params := &stripe.ProductParams{
		Name: stripe.String(name),
	}
	
	if metadata != nil {
		params.Metadata = metadata
	}
	
	return product.New(params)
}

// UpdateProduct updates a Stripe product
func (c *StripeClient) UpdateProduct(productID, name string, metadata map[string]string) (*stripe.Product, error) {
	params := &stripe.ProductParams{
		Name: stripe.String(name),
	}
	
	if metadata != nil {
		params.Metadata = metadata
	}
	
	return product.Update(productID, params)
}

// Price Operations

// CreatePrice creates a Stripe price
func (c *StripeClient) CreatePrice(productID string, amountCents int64, currency, interval string) (*stripe.Price, error) {
	params := &stripe.PriceParams{
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(amountCents),
		Currency:   stripe.String(currency),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(interval),
		},
	}
	
	return price.New(params)
}

// Customer Operations

// CreateCustomer creates a Stripe customer
func (c *StripeClient) CreateCustomer(email, teamID string, metadata map[string]string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Metadata: map[string]string{
			"team_id": teamID,
		},
	}
	
	// Merge additional metadata
	if metadata != nil {
		for k, v := range metadata {
			params.Metadata[k] = v
		}
	}
	
	return customer.New(params)
}

// GetCustomer retrieves a Stripe customer
func (c *StripeClient) GetCustomer(customerID string) (*stripe.Customer, error) {
	return customer.Get(customerID, nil)
}

// UpdateCustomer updates a Stripe customer
func (c *StripeClient) UpdateCustomer(customerID, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	
	return customer.Update(customerID, params)
}

// Checkout Session Operations

// CreateCheckoutSession creates a Stripe Checkout session
func (c *StripeClient) CreateCheckoutSession(priceID, customerID, successURL, cancelURL string, metadata map[string]string, trialDays int32) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
	}
	
	if customerID != "" {
		params.Customer = stripe.String(customerID)
	}
	
	if metadata != nil {
		params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: metadata,
		}
	}
	
	// Add trial period if specified
	if trialDays > 0 {
		if params.SubscriptionData == nil {
			params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{}
		}
		params.SubscriptionData.TrialPeriodDays = stripe.Int64(int64(trialDays))
	}
	
	return checkoutsession.New(params)
}

// GetCheckoutSession retrieves a Stripe Checkout session
func (c *StripeClient) GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("subscription")
	params.AddExpand("customer")
	
	return checkoutsession.Get(sessionID, params)
}

// Subscription Operations

// GetSubscription retrieves a Stripe subscription
func (c *StripeClient) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionID, nil)
}

// CancelSubscription cancels a Stripe subscription
func (c *StripeClient) CancelSubscription(subscriptionID string, cancelAtPeriodEnd bool) (*stripe.Subscription, error) {
	if cancelAtPeriodEnd {
		// Cancel at period end
		params := &stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(true),
		}
		return subscription.Update(subscriptionID, params)
	}
	
	// Cancel immediately
	params := &stripe.SubscriptionCancelParams{}
	return subscription.Cancel(subscriptionID, params)
}

// UpdateSubscription updates a Stripe subscription (e.g., change plan)
func (c *StripeClient) UpdateSubscription(subscriptionID, newPriceID string, prorationBehavior string) (*stripe.Subscription, error) {
	// Get current subscription to find the subscription item ID
	currentSub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get current subscription: %w", err)
	}
	
	if len(currentSub.Items.Data) == 0 {
		return nil, fmt.Errorf("subscription has no items")
	}
	
	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(currentSub.Items.Data[0].ID),
				Price: stripe.String(newPriceID),
			},
		},
		ProrationBehavior: stripe.String(prorationBehavior),
	}
	
	return subscription.Update(subscriptionID, params)
}

// ReactivateSubscription reactivates a canceled subscription
func (c *StripeClient) ReactivateSubscription(subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}
	return subscription.Update(subscriptionID, params)
}

// Customer Portal Operations

// CreateCustomerPortalSession creates a Stripe Customer Portal session
func (c *StripeClient) CreateCustomerPortalSession(customerID, returnURL string) (*stripe.BillingPortalSession, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(returnURL),
	}
	
	return portalsession.New(params)
}

// Invoice Operations

// GetUpcomingInvoice retrieves the upcoming invoice for a customer
func (c *StripeClient) GetUpcomingInvoice(customerID string) (*stripe.Invoice, error) {
	params := &stripe.InvoiceUpcomingParams{
		Customer: stripe.String(customerID),
	}
	
	return invoice.Upcoming(params)
}

// ListInvoices lists invoices for a customer
func (c *StripeClient) ListInvoices(customerID string, limit int64) ([]*stripe.Invoice, error) {
	params := &stripe.InvoiceListParams{
		Customer: stripe.String(customerID),
	}
	params.Limit = stripe.Int64(limit)
	
	iter := invoice.List(params)
	var invoices []*stripe.Invoice
	
	for iter.Next() {
		invoices = append(invoices, iter.Invoice())
	}
	
	if err := iter.Err(); err != nil {
		return nil, err
	}
	
	return invoices, nil
}

// GetInvoice retrieves a specific invoice
func (c *StripeClient) GetInvoice(invoiceID string) (*stripe.Invoice, error) {
	return invoice.Get(invoiceID, nil)
}

// Webhook Operations

// ConstructEvent constructs a Stripe event from webhook payload and signature
func (c *StripeClient) ConstructEvent(payload []byte, signature, webhookSecret string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signature, webhookSecret)
}
