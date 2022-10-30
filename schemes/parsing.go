package schemes

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	errs "github.com/Drozd0f/csv-app/errors"
	"github.com/Drozd0f/csv-app/pkg/comparator"
)

var ErrParsing = errors.New("parsing error")

func parsingMsgError(colName, value, parseType, tranID string) string {
	return fmt.Sprintf("impossible parse <%s: %s> to %s <TransactionId: %s>", colName, value, parseType, tranID)
}

func BindCsv(i interface{}, row, headers []string) error {
	t := reflect.TypeOf(i).Elem()
	val := reflect.ValueOf(i).Elem()

	fieldNames := setFieldNames(t)

	if err := readToStruct(val, row, headers, fieldNames); err != nil {
		return err
	}

	return nil
}

func setFieldNames(t reflect.Type) []string {
	numField := t.NumField()
	fieldNames := make([]string, 0, numField)

	for i := 0; i < numField; i++ {
		fieldNames = append(fieldNames, t.Field(i).Tag.Get("csv"))
	}

	return fieldNames
}

func readToStruct(t reflect.Value, row []string, headers, fieldNames []string) error {
	pkIdx, err := comparator.IdxSlice("TransactionId", headers)
	if err != nil {
		return &errs.ErrorWithMessage{
			Err: ErrParsing,
			Msg: fmt.Sprint("Not found column <TransactionId> in file"),
		}
	}

	for idx := range fieldNames {
		tf := t.Field(idx)
		rowIdx, err := comparator.IdxSlice(fieldNames[idx], headers)
		if err != nil {
			return &errs.ErrorWithMessage{
				Err: ErrParsing,
				Msg: fmt.Sprintf("Not found column <%s> in file", fieldNames[idx]),
			}
		}

		val := row[rowIdx]
		translationID := row[pkIdx]

		switch tf.Interface().(type) {
		case int32:
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return &errs.ErrorWithMessage{
					Err: ErrParsing,
					Msg: parsingMsgError(headers[idx], val, "int", translationID),
				}
			}
			tf.SetInt(i)
		case float32:
			i, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return &errs.ErrorWithMessage{
					Err: ErrParsing,
					Msg: parsingMsgError(headers[idx], val, "float", translationID),
				}
			}
			tf.SetFloat(i)
		case time.Time:
			layout := "2006-01-02 15:04:05"
			i, err := time.Parse(layout, val)
			if err != nil {
				return &errs.ErrorWithMessage{
					Err: ErrParsing,
					Msg: parsingMsgError(headers[idx], val, "time", translationID),
				}
			}
			tf.Set(reflect.ValueOf(i))
		default:
			tf.SetString(val)
		}
	}

	return nil
}
