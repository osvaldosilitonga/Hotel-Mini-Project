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
	ID                 string `json:"id" validate:"required"`
	ExternalID         string `json:"external_id" validate:"required"`
	UserID             string `json:"user_id" validate:"required"`
	PaymentMethod      string `json:"payment_method" validate:"required"`
	Status             string `json:"status" validate:"required"`
	MerchantName       string `json:"merchant_name"`
	Amount             int    `json:"amount" validate:"required"`
	PaidAmount         int    `json:"paid_amount" validate:"required"`
	BankCode           string `json:"bank_code" validate:"required"`
	PaidAt             string `json:"paid_at" validate:"required"`
	PayerEmail         string `json:"payer_emal" validate:"required"`
	Description        string `json:"description" validate:"required"`
	Currency           string `json:"currency" validate:"required"`
	PaymentChannel     string `json:"payment_channel" validate:"required"`
	PaymentDestination string `json:"payment_destination"`
}
