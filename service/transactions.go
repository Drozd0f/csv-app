package service

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	errs "github.com/Drozd0f/csv-app/errors"
	"github.com/Drozd0f/csv-app/repository"
	"github.com/Drozd0f/csv-app/schemes"
)

var (
	ErrOpenFile         = errors.New("invalid file")
	ErrParsing          = errors.New("invalid file signature")
	ErrTransactionExist = errors.New("transaction exist")
)

func (s *Service) UploadCsvFile(ctx context.Context, f IFileHeader) error {
	file, err := f.Open()
	if err != nil {
		return ErrOpenFile
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	// reading file headers
	l, _, err := reader.ReadLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		return err
	}
	headers := strings.Split(string(l), ",")
	chTrans := make(chan []schemes.Transaction)
	chError := make(chan error)

	// create goroutine for insert all transactions by one database transaction
	go s.r.InsertToTransactions(ctx, s.c.ChunkSize, chTrans, chError)
	// check channel error on error at goroutine initialization
	if err = <-chError; err != nil {
		return fmt.Errorf("InsertToTransactions initial goroutine: %w", err)
	}

	var errLoop error

mainLoop:
	for {
		records := make([]schemes.Transaction, 0, s.c.ChunkSize)
		select {
		case <-ctx.Done():
			// if user request is turn down connection
			break mainLoop
		default:
			// while ChunkSize to insert not equal to config value
			// fil chunk slices
			for int32(len(records)) != s.c.ChunkSize {
				row, _, err := reader.ReadLine()
				if err != nil {
					if errors.Is(err, io.EOF) {
						chTrans <- records
						errLoop = <-chError
						break mainLoop
					}

					errLoop = err
					break mainLoop
				}

				var t schemes.Transaction
				if err = schemes.BindFromCsv(&t, strings.Split(string(row), ","), headers); err != nil {
					errLoop = err
					break mainLoop
				}

				records = append(records, t)
			}

			chTrans <- records
			if errLoop = <-chError; errLoop != nil {
				break mainLoop
			}
		}
	}

	if errLoop != nil {
		switch {
		case errors.Is(errLoop, schemes.ErrParsing):
			var erw *errs.ErrorWithMessage
			if errors.As(errLoop, &erw) {
				return &errs.ErrorWithMessage{
					Err: ErrParsing,
					Msg: erw.Msg,
				}
			}
		case errors.Is(errLoop, repository.ErrUniqueConstraint):
			var erw *errs.ErrorWithMessage
			if errors.As(errLoop, &erw) {
				return &errs.ErrorWithMessage{
					Err: ErrTransactionExist,
					Msg: erw.Msg,
				}
			}
		}

		return errLoop
	}

	return nil
}

func (s *Service) GetFilteredTransactions(ctx context.Context, rrt schemes.RawTransactionFilter) (schemes.SliceTransactions, error) {
	storedSliceT, err := s.r.GetSliceTransactions(ctx, schemes.NewTransactionFilterFromRaw(rrt))
	if err != nil {
		return schemes.SliceTransactions{}, fmt.Errorf("repository get slice transactions: %w", err)
	}

	return schemes.NewSliceTransactionsFromDB(storedSliceT), nil
}

func (s *Service) DownloadCsvFile(ctx context.Context, w io.Writer) error {
	p := schemes.Paginator{
		Page:  1,
		Limit: s.c.ChunkSize, // TODO: put in config chunks
	}
	writer := csv.NewWriter(w)
	isHeaderWrote := false

	for {
		storTrans, err := s.r.GetTransactions(ctx, p)
		if err != nil {
			return err
		}

		st := schemes.NewSliceTransactionsFromDB(storTrans)
		if !isHeaderWrote {
			if err = writer.Write(st.GetCsvNames()); err != nil {
				return err
			}
			isHeaderWrote = true
		}

		if err = writer.WriteAll(st.ToString()); err != nil {
			return err
		}

		if int32(len(storTrans)) < p.Limit {
			return nil
		}

		p.Page++
	}
}
