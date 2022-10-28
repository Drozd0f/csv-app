package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Drozd0f/csv-app/db"
	errs "github.com/Drozd0f/csv-app/errors"
	"github.com/Drozd0f/csv-app/schemes"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"
)

var (
	ErrParsing          = errors.New("parsing error")
	ErrUniqueConstraint = errors.New("unique constraint")
)

func (r *Repository) InsertToTransactions(ctx context.Context, chunkSize int, headers []string, chRows chan [][]string, chError chan error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		chError <- err
		return
	}

mainLoop:
	for {
		select {
		case <-ctx.Done():
			chError <- nil
			return
		case rows := <-chRows:
			tranParam, err := parseToTransactionsParams(headers, rows)
			if err != nil {
				chError <- err
				return
			}

			_, err = r.q.CreateTransactions(ctx, tranParam)
			if err != nil {
				var pgError *pgconn.PgError
				if errors.As(err, &pgError) {
					if pgError.Code == uniqueConstraintCode {
						err = &errs.ErrorWithMessage{
							Err: ErrUniqueConstraint,
							Msg: fmt.Sprintf("transaction insert error: %s", strings.ToLower(pgError.Detail)),
						}
						chError <- multierr.Append(err, tx.Rollback(ctx))
						return
					}
				}

				chError <- multierr.Append(err, tx.Rollback(ctx))
				return
			}

			if len(rows) < chunkSize {
				break mainLoop
			}

			chError <- nil
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		chError <- err
		return
	}

	chError <- nil
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
