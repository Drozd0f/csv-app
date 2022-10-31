package schemes

import (
	"reflect"
	"strconv"
	"time"

	"github.com/Drozd0f/csv-app/db"
)

var (
	Statuses      = [2]string{"accepted", "declined"}
	PaymentTypes  = [2]string{"cash", "card"}
	DefaultString = "default"
)

type RawDatePost struct {
	From string `form:"from"`
	To   string `form:"to"`
}

type DatePost struct {
	From time.Time
	To   time.Time
}

type SliceTransactions []Transaction

func NewSliceTransactionsFromDB(storedTransactions []db.Transaction) SliceTransactions {
	srt := make(SliceTransactions, 0, len(storedTransactions))
	for _, st := range storedTransactions {
		srt = append(srt, NewTransactionFromDB(st))
	}

	return srt
}

func (st SliceTransactions) GetCsvNames() []string {
	if len(st) > 0 {
		return st[0].GetCsvNames()
	}

	return nil
}

func (st SliceTransactions) ToString() [][]string {
	sliceS := make([][]string, 0, len(st))
	for _, t := range st {
		sliceS = append(sliceS, t.ToString())
	}

	return sliceS
}

type RawTransactionFilter struct {
	ID               string   `form:"transaction_id"`
	Status           string   `form:"status"`
	TerminalIDs      []string `form:"terminal_id"`
	PaymentType      string   `form:"payment_type"`
	DatePost         RawDatePost
	PaymentNarrative string `form:"payment_narrative"`
}

type TransactionFilter struct {
	ID               int32
	Status           string
	TerminalIDs      []int32
	PaymentType      string
	DatePost         DatePost
	PaymentNarrative string
}

func NewTransactionFilterFromRaw(r RawTransactionFilter) TransactionFilter {
	return TransactionFilter{
		ID:               stringToInt32(r.ID),
		Status:           stringToExpected(r.Status, Statuses),
		TerminalIDs:      sliceStringToInt32(r.TerminalIDs),
		PaymentType:      stringToExpected(r.PaymentType, PaymentTypes),
		DatePost:         newDatePostFromRaw(r.DatePost),
		PaymentNarrative: r.PaymentNarrative,
	}
}

type Transaction struct {
	TransactionID      int32     `json:"transaction_id" csv:"TransactionId"`
	RequestID          int32     `json:"request_id" csv:"RequestId"`
	TerminalID         int32     `json:"terminal_id" csv:"TerminalId"`
	PartnerObjectID    int32     `json:"partner_object_id" csv:"PartnerObjectId"`
	AmountTotal        float32   `json:"amount_total" csv:"AmountTotal"`
	AmountOriginal     float32   `json:"amount_original" csv:"AmountOriginal"`
	CommissionPs       float32   `json:"commission_ps" csv:"CommissionPS"`
	CommissionClient   float32   `json:"commission_client" csv:"CommissionClient"`
	CommissionProvider float32   `json:"commission_provider" csv:"CommissionProvider"`
	DateInput          time.Time `json:"date_input" csv:"DateInput"`
	DatePost           time.Time `json:"date_post" csv:"DatePost"`
	Status             string    `json:"status" csv:"Status"`
	PaymentType        string    `json:"payment_type" csv:"PaymentType"`
	PaymentNumber      string    `json:"payment_number" csv:"PaymentNumber"`
	ServiceID          int32     `json:"service_id" csv:"ServiceId"`
	Service            string    `json:"service" csv:"Service"`
	PayeeID            int32     `json:"payee_id" csv:"PayeeId"`
	PayeeName          string    `json:"payee_name" csv:"PayeeName"`
	PayeeBankMfo       int32     `json:"payee_bank_mfo" csv:"PayeeBankMfo"`
	PayeeBankAccount   string    `json:"payee_bank_account" csv:"PayeeBankAccount"`
	PaymentNarrative   string    `json:"payment_narrative" csv:"PaymentNarrative"`
}

func NewTransactionFromDB(storedT db.Transaction) Transaction {
	return Transaction{
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

func (t Transaction) GetCsvNames() []string {
	return getFieldNames(reflect.TypeOf(t))
}

func (t Transaction) ToString() []string {
	return BindToCsv(t)
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
