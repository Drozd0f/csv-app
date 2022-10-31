package schemes

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Drozd0f/csv-app/db"
	"github.com/stretchr/testify/require"
)

const countStoredTransaction = 50

var expectedCsvName = [21]string{
	"TransactionId", "RequestId", "TerminalId", "PartnerObjectId",
	"AmountTotal", "AmountOriginal", "CommissionPS", "CommissionClient",
	"CommissionProvider", "DateInput", "DatePost", "Status",
	"PaymentType", "PaymentNumber", "ServiceId", "Service",
	"PayeeId", "PayeeName", "PayeeBankMfo", "PayeeBankAccount", "PaymentNarrative",
}

func createStoredTransaction() db.Transaction {
	return db.Transaction{
		TransactionID:      rand.Int31(),
		RequestID:          rand.Int31(),
		TerminalID:         rand.Int31(),
		PartnerObjectID:    rand.Int31(),
		AmountTotal:        rand.Int31(),
		AmountOriginal:     rand.Int31(),
		CommissionPs:       rand.Int31(),
		CommissionClient:   rand.Int31(),
		CommissionProvider: rand.Int31(),
		DateInput:          time.Now(),
		DatePost:           time.Now(),
		Status:             Statuses[rand.Intn(1)],
		PaymentType:        PaymentTypes[rand.Intn(1)],
		PaymentNumber:      "PS16698205",
		ServiceID:          rand.Int31(),
		Service:            "Поповнення карток",
		PayeeID:            rand.Int31(),
		PayeeName:          "privat",
		PayeeBankMfo:       rand.Int31(),
		PayeeBankAccount:   "UA713451373919523",
		PaymentNarrative:   "Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р.",
	}
}

func generateStoredSliceTransactions() []db.Transaction {
	storedT := make([]db.Transaction, 0, countStoredTransaction)
	for idx := 0; idx < countStoredTransaction; idx++ {
		storedT = append(storedT, createStoredTransaction())
	}

	return storedT
}

func TestNewSliceTransactionsFromDB(t *testing.T) {
	storedTs := generateStoredSliceTransactions()
	sliceTs := NewSliceTransactionsFromDB(storedTs)
	require.Equal(t, len(storedTs), len(sliceTs))

	for idx := 0; idx < countStoredTransaction; idx++ {
		storedT, sliceT := storedTs[idx], sliceTs[idx]
		require.Equal(t, float32(storedT.AmountTotal)/100, sliceT.AmountTotal)
		require.Equal(t, float32(storedT.AmountOriginal)/100, sliceT.AmountOriginal)
		require.Equal(t, float32(storedT.CommissionPs)/100, sliceT.CommissionPs)
		require.Equal(t, float32(storedT.CommissionClient)/100, sliceT.CommissionClient)
		require.Equal(t, float32(storedT.CommissionProvider)/100, sliceT.CommissionProvider)

	}
}

func TestSliceTransactions_GetCsvNames(t *testing.T) {
	st := SliceTransactions{}
	names := st.GetCsvNames()
	for idx, got := range names {
		expected := expectedCsvName[idx]
		require.Equal(t, expected, got)
	}
}

func TestSliceTransactions_ToString(t *testing.T) {
	storedTs := generateStoredSliceTransactions()
	sliceTs := NewSliceTransactionsFromDB(storedTs)
	stringSliceTs := sliceTs.ToString()
	require.Equal(t, len(sliceTs), len(stringSliceTs))
}

func TestNewTransactionFromDB(t *testing.T) {
	storedT := createStoredTransaction()
	sliceT := NewTransactionFromDB(storedT)

	require.Equal(t, float32(storedT.AmountTotal)/100, sliceT.AmountTotal)
	require.Equal(t, float32(storedT.AmountOriginal)/100, sliceT.AmountOriginal)
	require.Equal(t, float32(storedT.CommissionPs)/100, sliceT.CommissionPs)
	require.Equal(t, float32(storedT.CommissionClient)/100, sliceT.CommissionClient)
	require.Equal(t, float32(storedT.CommissionProvider)/100, sliceT.CommissionProvider)
}

func TestTransaction_GetCsvNames(t *testing.T) {
	st := Transaction{}
	names := st.GetCsvNames()
	for idx, got := range names {
		expected := expectedCsvName[idx]
		require.Equal(t, expected, got)
	}
}
