package schemes

import (
	"strconv"
	"time"

	"github.com/Drozd0f/csv-app/db"
)

var (
	statuses      = [2]string{"accepted", "declined"}
	paymentTypes  = [2]string{"cash", "card"}
	DefaultString = "default"
)

type SliceResponseTransactions []ResponseTransaction

type RawDatePost struct {
	From string `form:"from"`
	To   string `form:"to"`
}

type DatePost struct {
	From time.Time
	To   time.Time
}

type RawRequestTransaction struct {
	ID               string   `form:"transaction_id"`
	Status           string   `form:"status"`
	TerminalIDs      []string `form:"terminal_id"`
	PaymentType      string   `form:"payment_type"`
	DatePost         RawDatePost
	PaymentNarrative string `form:"payment_narrative"`
}

type RequestTransaction struct {
	ID               int32
	Status           string
	TerminalIDs      []int32
	PaymentType      string
	DatePost         DatePost
	PaymentNarrative string
}

type ResponseTransaction struct {
	TransactionID      int32     `json:"transaction_id"`
	RequestID          int32     `json:"request_id"`
	TerminalID         int32     `json:"terminal_id"`
	PartnerObjectID    int32     `json:"partner_object_id"`
	AmountTotal        float32   `json:"amount_total"`
	AmountOriginal     float32   `json:"amount_original"`
	CommissionPs       float32   `json:"commission_ps"`
	CommissionClient   float32   `json:"commission_client"`
	CommissionProvider float32   `json:"commission_provider"`
	DateInput          time.Time `json:"date_input"`
	DatePost           time.Time `json:"date_post"`
	Status             string    `json:"status"`
	PaymentType        string    `json:"payment_type"`
	PaymentNumber      string    `json:"payment_number"`
	ServiceID          int32     `json:"service_id"`
	Service            string    `json:"service"`
	PayeeID            int32     `json:"payee_id"`
	PayeeName          string    `json:"payee_name"`
	PayeeBankMfo       int32     `json:"payee_bank_mfo"`
	PayeeBankAccount   string    `json:"payee_bank_account"`
	PaymentNarrative   string    `json:"payment_narrative"`
}

func NewRequestTransactionFromRaw(r RawRequestTransaction) RequestTransaction {
	return RequestTransaction{
		ID:               stringToInt32(r.ID),
		Status:           stringToExpected(r.Status, statuses),
		TerminalIDs:      sliceStringToInt32(r.TerminalIDs),
		PaymentType:      stringToExpected(r.PaymentType, paymentTypes),
		DatePost:         newDatePostFromRaw(r.DatePost),
		PaymentNarrative: r.PaymentNarrative,
	}
}

func NewSliceResponseTransactionFromDB(storedTransactions []db.Transaction) SliceResponseTransactions {
	srt := make(SliceResponseTransactions, 0, len(storedTransactions))
	for _, st := range storedTransactions {
		srt = append(srt, NewResponseTransactionFromDB(st))
	}

	return srt
}

func NewResponseTransactionFromDB(storedT db.Transaction) ResponseTransaction {
	return ResponseTransaction{
		TransactionID:      storedT.TransactionID,
		RequestID:          storedT.RequestID,
		TerminalID:         storedT.TerminalID,
		PartnerObjectID:    storedT.PartnerObjectID,
		AmountTotal:        float32(storedT.AmountTotal) / 100,
		AmountOriginal:     float32(storedT.AmountOriginal) / 100,
		CommissionPs:       float32(storedT.CommissionPs) / 100,
		CommissionClient:   float32(storedT.CommissionClient) / 100,
		CommissionProvider: float32(storedT.CommissionProvider) / 100,
		DateInput:          storedT.DateInput,
		DatePost:           storedT.DatePost,
		Status:             storedT.Status,
		PaymentType:        storedT.PaymentType,
		PaymentNumber:      storedT.PaymentNumber,
		ServiceID:          storedT.ServiceID,
		Service:            storedT.Service,
		PayeeID:            storedT.PayeeID,
		PayeeName:          storedT.PayeeName,
		PayeeBankMfo:       storedT.PayeeBankMfo,
		PayeeBankAccount:   storedT.PayeeBankAccount,
		PaymentNarrative:   storedT.PaymentNarrative,
	}
}

func stringToInt32(s string) int32 {
	// Ignore error and return 0 as default value
	intS, _ := strconv.ParseInt(s, 10, 32)
	return int32(intS)
}

func sliceStringToInt32(sliceS []string) []int32 {
	sliceI := make([]int32, 0, len(sliceS))

	for _, s := range sliceS {
		sliceI = append(sliceI, stringToInt32(s))
	}

	return sliceI
}

func stringToExpected(s string, expectedStrings [2]string) string {
	for _, str := range expectedStrings {
		if s == str {
			return s
		}
	}

	return DefaultString
}

func newDatePostFromRaw(dpr RawDatePost) DatePost {
	layout := "2006-01-02"
	timeFrom, err := time.Parse(layout, dpr.From)
	timeTo, err := time.Parse(layout, dpr.To)
	if err != nil {
		timeTo = time.Now()
	}
	return DatePost{
		From: timeFrom,
		To:   timeTo,
	}
}
