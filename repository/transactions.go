package repository

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Drozd0f/csv-app/db"
	"github.com/jackc/pgx/v4"
)

func (r *Repository) InsertToTransactions(ctx context.Context, headersIdx map[string]int, rows [][]string) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	_, err = r.q.CreateTransactions(ctx, parseToTransactionsParams(headersIdx, rows)) // i - chunks
	if err != nil {
		log.Println(err)
		err = tx.Rollback(ctx)
		if err != nil {
			log.Println(err)
		}
		return err
	}

	tx.Commit(ctx)
	return nil
}

func parseToInt32(v string) int32 {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		log.Println("parse int:", err)
		return 0
	}
	return int32(i)
}

func parseToFloat32(v string) float32 {
	i, err := strconv.ParseFloat(v, 32)
	if err != nil {
		log.Println("parse int:", err)
		return 0
	}
	return float32(i)
}

func parseToTime(v string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, v)
	if err != nil {
		log.Println("parse time:", err)
		return time.Time{}
	}
	return t
}

func parseToTransactionsParams(headersIdx map[string]int, rows [][]string) []db.CreateTransactionsParams {
	dbParam := make([]db.CreateTransactionsParams, 0, len(rows))

	for _, row := range rows {
		dbParam = append(dbParam, db.CreateTransactionsParams{
			TransactionID:      parseToInt32(row[headersIdx["TransactionID"]]),
			RequestID:          parseToInt32(row[headersIdx["RequestID"]]),
			TerminalID:         parseToInt32(row[headersIdx["TerminalID"]]),
			PartnerObjectID:    parseToInt32(row[headersIdx["PartnerObjectID"]]),
			AmountTotal:        int32(parseToFloat32(row[headersIdx["AmountTotal"]]) * 100),
			AmountOriginal:     int32(parseToFloat32(row[headersIdx["AmountOriginal"]]) * 100),
			CommissionPs:       int32(parseToFloat32(row[headersIdx["CommissionPs"]]) * 100),
			CommissionClient:   int32(parseToFloat32(row[headersIdx["CommissionClient"]]) * 100),
			CommissionProvider: int32(parseToFloat32(row[headersIdx["CommissionProvider"]]) * 100),
			DateInput:          parseToTime(row[headersIdx["DateInput"]]),
			DatePost:           parseToTime(row[headersIdx["DatePost"]]),
			Status:             row[headersIdx["Status"]],
			PaymentType:        row[headersIdx["PaymentType"]],
			PaymentNumber:      row[headersIdx["PaymentNumber"]],
			ServiceID:          parseToInt32(row[headersIdx["ServiceID"]]),
			Service:            row[headersIdx["Service"]],
			PayeeID:            parseToInt32(row[headersIdx["PayeeID"]]),
			PayeeName:          row[headersIdx["PayeeName"]],
			PayeeBankMfo:       parseToInt32(row[headersIdx["PayeeBankMfo"]]),
			PayeeBankAccount:   row[headersIdx["PayeeBankAccount"]],
			PaymentNarrative:   row[headersIdx["PaymentNarrative"]],
		})
	}

	return dbParam
}
