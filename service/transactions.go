package service

import (
	"bufio"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

func parseHeader(h []string) map[string]int {
	header := make(map[string]int, len(h))
	for idx, nameCol := range h {
		header[nameCol] = idx
	}
	return header
}

func (s *Service) UploadCsvFile(ctx context.Context, f *multipart.FileHeader) error {
	file, err := f.Open()
	if err != nil {
		return err
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
	headerIdx := parseHeader(strings.Split(string(l), ","))
	chunks := 10

main_loop:
	for {
		records := make([][]string, 0, chunks)
		for len(records) != chunks {
			l, _, err := reader.ReadLine()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break main_loop
				}

				return err
			}
			records = append(records, strings.Split(string(l), ","))
		}
		s.r.InsertToTransactions(ctx, headerIdx, records)
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
