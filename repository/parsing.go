package repository

import (
	"github.com/Drozd0f/csv-app/db"
	"github.com/Drozd0f/csv-app/schemes"
)

func parseToTransactionsParams(rows []schemes.Transaction) []db.CreateTransactionsParams {
	dbParam := make([]db.CreateTransactionsParams, 0, len(rows))
	for _, row := range rows {
		dbParam = append(dbParam, db.CreateTransactionsParams{
			TransactionID:      row.TransactionID,
			RequestID:          row.RequestID,
			TerminalID:         row.TerminalID,
			PartnerObjectID:    row.PartnerObjectID,
			AmountTotal:        int32(row.AmountTotal * 100),
			AmountOriginal:     int32(row.AmountOriginal * 100),
			CommissionPs:       int32(row.CommissionPs * 100),
			CommissionClient:   int32(row.CommissionClient * 100),
			CommissionProvider: int32(row.CommissionProvider * 100),
			DateInput:          row.DateInput,
			DatePost:           row.DatePost,
			Status:             row.Status,
			PaymentType:        row.PaymentType,
			PaymentNumber:      row.PaymentNumber,
			ServiceID:          row.ServiceID,
			Service:            row.Service,
			PayeeID:            row.PayeeID,
			PayeeName:          row.PayeeName,
			PayeeBankMfo:       row.PayeeBankMfo,
			PayeeBankAccount:   row.PayeeBankAccount,
			PaymentNarrative:   row.PaymentNarrative,
		})
	}

	return dbParam
}
