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
	ErrUniqueConstraint = errors.New("unique constraint")
)

func (r *Repository) InsertToTransactions(ctx context.Context, chunkSize int32, chRows chan []schemes.Transaction, chError chan error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		chError <- err
		return
	}
	chError <- nil

mainLoop:
	for {
		select {
		case <-ctx.Done():
			return
		case rows := <-chRows:
			tranParam := parseToTransactionsParams(rows)
			_, err = r.q.CreateTransactions(ctx, tranParam)
			if err != nil {
				var pgError *pgconn.PgError
				if errors.As(err, &pgError) {
					if pgError.Code == uniqueConstraintCode {
						err = &errs.ErrorWithMessage{
							Err: ErrUniqueConstraint,
							Msg: fmt.Sprintf("transaction insert error %s", strings.ToLower(pgError.Detail)),
						}
						chError <- multierr.Append(err, tx.Rollback(ctx))
						return
					}
				}

				chError <- multierr.Append(err, tx.Rollback(ctx))
				return
			}

			if int32(len(rows)) < chunkSize {
				break mainLoop
			}

			chError <- nil
		}
	}

	if err = tx.Commit(ctx); err != nil {
		chError <- err
		return
	}

	chError <- nil
}

func (r *Repository) GetSliceTransactions(ctx context.Context, rt schemes.TransactionFilter) ([]db.Transaction, error) {
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

func (r *Repository) GetTransactions(ctx context.Context, p schemes.Paginator) ([]db.Transaction, error) {
	storedTrans, err := r.q.GetTransactions(ctx, db.GetTransactionsParams{
		Offset: p.Offset(),
		Limit:  p.Limit,
	})
	if err != nil {
		return nil, err
	}

	return storedTrans, nil
}
