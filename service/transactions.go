package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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

func (s *Service) UploadCsvFile(ctx context.Context, f *multipart.FileHeader) error {
	file, err := f.Open()
	if err != nil {
		return ErrOpenFile
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	l, _, err := reader.ReadLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		return err
	}
	headers := strings.Split(string(l), ",")
	chunks := 10
	chRows := make(chan [][]string)
	chError := make(chan error)

	go s.r.InsertToTransactions(ctx, chunks, headers, chRows, chError)

	var errLoop error

mainLoop:
	for {
		records := make([][]string, 0, chunks)
		select {
		case <-ctx.Done():
			break mainLoop
		default:
			for len(records) != chunks {
				l, _, err := reader.ReadLine()
				if err != nil {
					if errors.Is(err, io.EOF) {
						chRows <- records
						errLoop = <-chError
						break mainLoop
					}

					errLoop = err
					break mainLoop
				}
				records = append(records, strings.Split(string(l), ","))
			}

			chRows <- records
			if errLoop = <-chError; errLoop != nil {
				break mainLoop
			}
		}
	}

	if errLoop != nil {
		switch {
		case errors.Is(errLoop, repository.ErrParsing):
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

func (s *Service) GetSliceTransactions(ctx context.Context, rrt schemes.RawRequestTransaction) (schemes.SliceResponseTransactions, error) {
	storedSliceT, err := s.r.GetSliceTransactions(ctx, schemes.NewRequestTransactionFromRaw(rrt))
	if err != nil {
		return schemes.SliceResponseTransactions{}, fmt.Errorf("repository get slice transactions: %w", err)
	}

	return schemes.NewSliceResponseTransactionFromDB(storedSliceT), nil
}
