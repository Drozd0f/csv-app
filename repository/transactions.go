package repository

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/Drozd0f/csv-app/db"
	"github.com/Drozd0f/csv-app/schemes"
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

	err = tx.Commit(ctx)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (r *Repository) GetSliceTransactions(ctx context.Context, rt schemes.RequestTransaction) ([]db.Transaction, error) {
	sliceT, err := r.q.SliceTransactions(ctx, db.SliceTransactionsParams{
		TransactionID: rt.ID,
		Status:        rt.Status,
		PaymentType:   rt.PaymentType,
		DatePostFrom:  rt.DatePost.From,
		DatePostTo:    rt.DatePost.To,
		PaymentNarrative: sql.NullString{
			String: rt.PaymentNarrative,
			Valid:  true,
		},
		TerminalID: rt.TerminalIDs,
	})
	if err != nil {
		return nil, err
	}

	return sliceT, nil
}

func parseToInt32(v string) int32 {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		log.Println("parse int:", err)
	}
	return int32(i)
}

func parseToFloat32(v string) float32 {
	i, err := strconv.ParseFloat(v, 32)
	if err != nil {
		log.Println("parse int:", err)
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
			TransactionID:      parseToInt32(row[headersIdx["TransactionId"]]),
			RequestID:          parseToInt32(row[headersIdx["RequestId"]]),
			TerminalID:         parseToInt32(row[headersIdx["TerminalId"]]),
			PartnerObjectID:    parseToInt32(row[headersIdx["PartnerObjectId"]]),
			AmountTotal:        int32(parseToFloat32(row[headersIdx["AmountTotal"]]) * 100),
			AmountOriginal:     int32(parseToFloat32(row[headersIdx["AmountOriginal"]]) * 100),
			CommissionPs:       int32(parseToFloat32(row[headersIdx["CommissionPS"]]) * 100),
			CommissionClient:   int32(parseToFloat32(row[headersIdx["CommissionClient"]]) * 100),
			CommissionProvider: int32(parseToFloat32(row[headersIdx["CommissionProvider"]]) * 100),
			DateInput:          parseToTime(row[headersIdx["DateInput"]]),
			DatePost:           parseToTime(row[headersIdx["DatePost"]]),
			Status:             row[headersIdx["Status"]],
			PaymentType:        row[headersIdx["PaymentType"]],
			PaymentNumber:      row[headersIdx["PaymentNumber"]],
			ServiceID:          parseToInt32(row[headersIdx["ServiceId"]]),
			Service:            row[headersIdx["Service"]],
			PayeeID:            parseToInt32(row[headersIdx["PayeeId"]]),
			PayeeName:          row[headersIdx["PayeeName"]],
			PayeeBankMfo:       parseToInt32(row[headersIdx["PayeeBankMfo"]]),
			PayeeBankAccount:   row[headersIdx["PayeeBankAccount"]],
			PaymentNarrative:   row[headersIdx["PaymentNarrative"]],
		})
	}

	return dbParam
}
