package dto

import "time"

type CreateInvoiceData struct {
	Description string
	Email       string
	ExternalID  string
	Amount      float32
}

type InvoiceResponse struct {
	InvoiceUrl  string    `json:"invoice_url"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"created"`
	ExpairyDate time.Time `json:"expairy_date"`
	ExternalId  string    `json:"external_id"`
	PayerEmail  string    `json:"payer_email"`
}

type XenditCallbackBody struct {
	ID                 string `json:"id"`
	ExternalID         string `json:"external_id"`
	UserID             string `json:"user_id"`
	PaymentMethod      string `json:"payment_method"`
	Status             string `json:"status"`
	MerchantName       string `json:"merchant_name"`
	Amount             int    `json:"amount"`
	PaidAmount         int    `json:"paid_amount"`
	BankCode           string `json:"bank_code"`
	PaidAt             string `json:"paid_at"`
	PayerEmail         string `json:"payer_email"`
	Description        string `json:"description"`
	Currency           string `json:"currency"`
	PaymentChannel     string `json:"payment_channel"`
	PaymentDestination string `json:"payment_destination"`
}
