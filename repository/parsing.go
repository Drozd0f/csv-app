package repository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Drozd0f/csv-app/db"
	errs "github.com/Drozd0f/csv-app/errors"
	"github.com/Drozd0f/csv-app/pkg/comparator"
)

var (
	headerParseToInt32 = []string{
		"TransactionId", "RequestId", "TerminalId",
		"PartnerObjectId", "ServiceId", "PayeeId",
		"PayeeBankMfo",
	}
	headerParseTofloat32 = []string{
		"AmountTotal", "AmountOriginal", "CommissionPS",
		"CommissionClient", "CommissionProvider",
	}
	headerParseToTime = []string{"DateInput", "DatePost"}
)

func parsingMsgError(colName, value, parseType, tranID string) string {
	return fmt.Sprintf("impossible parse <%s: %s> to %s <TransactionId: %s>", colName, value, parseType, tranID)
}

func parseToInt32(v string) (int32, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	return int32(i), err
}

func parseToFloat32(v string) (float32, error) {
	i, err := strconv.ParseFloat(v, 32)
	return float32(i), err
}

func parseToTime(v string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, v)
}

func parseRow(row, headers []string) (map[string]any, error) {
	newRow := make(map[string]any, len(row))
	translationID := row[comparator.IdxSlice("TransactionId", headers)]

	for idx, col := range row {
		header := headers[idx]
		if comparator.InSlice(header, headerParseToInt32) {
			v, err := parseToInt32(col)
			if err != nil {
				return nil, &errs.ErrorWithMessage{
					Err: ErrParsing,
					Msg: parsingMsgError(header, col, "int", translationID),
				}
			}

			newRow[header] = v
			continue
		}
		if comparator.InSlice(header, headerParseTofloat32) {
			v, err := parseToFloat32(col)
			if err != nil {
				return nil, &errs.ErrorWithMessage{
					Err: ErrParsing,
					Msg: parsingMsgError(header, col, "float", translationID),
				}
			}

			newRow[header] = v * 100
			continue
		}
		if comparator.InSlice(header, headerParseToTime) {
			v, err := parseToTime(col)
			if err != nil {
				return nil, &errs.ErrorWithMessage{
					Err: ErrParsing,
					Msg: parsingMsgError(header, col, "time", translationID),
				}
			}

			newRow[header] = v
			continue
		}
		newRow[header] = col
	}

	return newRow, nil
}

func parseToTransactionsParams(headers []string, rows [][]string) ([]db.CreateTransactionsParams, error) {
	dbParam := make([]db.CreateTransactionsParams, 0, len(rows))

	for _, row := range rows {
		newRow, err := parseRow(row, headers)
		if err != nil {
			return nil, err
		}

		dbParam = append(dbParam, db.CreateTransactionsParams{
			TransactionID:      newRow["TransactionId"].(int32),
			RequestID:          newRow["RequestId"].(int32),
			TerminalID:         newRow["TerminalId"].(int32),
			PartnerObjectID:    newRow["PartnerObjectId"].(int32),
			AmountTotal:        int32(newRow["AmountTotal"].(float32)),
			AmountOriginal:     int32(newRow["AmountOriginal"].(float32)),
			CommissionPs:       int32(newRow["CommissionPS"].(float32)),
			CommissionClient:   int32(newRow["CommissionClient"].(float32)),
			CommissionProvider: int32(newRow["CommissionProvider"].(float32)),
			DateInput:          newRow["DateInput"].(time.Time),
			DatePost:           newRow["DatePost"].(time.Time),
			Status:             newRow["Status"].(string),
			PaymentType:        newRow["PaymentType"].(string),
			PaymentNumber:      newRow["PaymentNumber"].(string),
			ServiceID:          newRow["ServiceId"].(int32),
			Service:            newRow["Service"].(string),
			PayeeID:            newRow["PayeeId"].(int32),
			PayeeName:          newRow["PayeeName"].(string),
			PayeeBankMfo:       newRow["PayeeBankMfo"].(int32),
			PayeeBankAccount:   newRow["PayeeBankAccount"].(string),
			PaymentNarrative:   newRow["PaymentNarrative"].(string),
		})
	}

	return dbParam, nil
}
