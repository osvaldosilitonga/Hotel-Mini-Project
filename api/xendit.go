package api

import (
	"context"
	"encoding/json"
	"fmt"
	"hotel/dto"
	"os"
	"time"

	xendit "github.com/xendit/xendit-go/v3"
	invoice "github.com/xendit/xendit-go/v3/invoice"
)

func CreteInvoice(d *dto.CreateInvoiceData) *dto.InvoiceResponse {
	desc := d.Description
	payerEmail := d.Email

	createInvoiceRequest := invoice.CreateInvoiceRequest{
		ExternalId:  d.ExternalID,
		Amount:      d.Amount,
		PayerEmail:  &payerEmail,
		Description: &desc,
	}

	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	resp, r, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).CreateInvoiceRequest(createInvoiceRequest).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

		b, _ := json.Marshal(err.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return nil
	}
	// response from `CreateInvoice`: Invoice
	// fmt.Fprintf(os.Stdout, "Response from `InvoiceApi.CreateInvoice`: %v\n", resp)

	data := dto.InvoiceResponse{
		InvoiceUrl:  resp.InvoiceUrl,
		Status:      string(resp.Status),
		Description: *resp.Description,
		CreateAt:    resp.Created,
		ExpairyDate: resp.ExpiryDate.In(time.Local),
		ExternalId:  resp.ExternalId,
		PayerEmail:  payerEmail,
	}

	return &data
}
