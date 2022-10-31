package service

import (
	"bytes"
	"encoding/csv"
	"mime/multipart"
	"strconv"

	errs "github.com/Drozd0f/csv-app/errors"
	mock_iservices "github.com/Drozd0f/csv-app/interfaces/services/mock"
	"github.com/Drozd0f/csv-app/pkg/comparator"
	"github.com/Drozd0f/csv-app/schemes"
	"github.com/Drozd0f/csv-app/test/fixtures"
	"github.com/golang/mock/gomock"
)

func (ts *ServiceTestSuite) TestUploadCsvFile() {
	ts.Run("empty csv", func() {
		f := &multipart.FileHeader{}

		err := ts.service.UploadCsvFile(ts.ctx, f)
		ts.Require().ErrorIs(err, ErrOpenFile)
	})

	ts.Run("parsing error", func() {
		ctrl := gomock.NewController(ts.T())
		defer ctrl.Finish()

		f, err := fixtures.Fixtures.Open("parsing_error_transactions_test.csv")
		ts.Require().NoError(err)

		mfh := mock_iservices.NewMockIFileHeader(ctrl)
		mf := mock_iservices.NewMockIFile(ctrl)

		mfh.EXPECT().Open().Return(mf, nil).Times(1)
		mf.EXPECT().Close().Do(f.Close).Times(1)
		mf.EXPECT().Read(gomock.Any()).DoAndReturn(f.Read).AnyTimes()

		err = ts.service.UploadCsvFile(ts.ctx, mfh)

		var ewm *errs.ErrorWithMessage
		ts.Require().ErrorAs(err, &ewm)
		ts.Require().ErrorIs(ewm.Err, ErrParsing)
		ts.Require().Contains(ewm.Msg, "impossible parse <DateInput: some_time> to time <TransactionId: 2>")
	})

	ts.Run("transaction exist error", func() {
		ts.SeedDB()

		ctrl := gomock.NewController(ts.T())
		defer ctrl.Finish()

		f, err := fixtures.Fixtures.Open("transactions_test.csv")
		ts.Require().NoError(err)

		mfh := mock_iservices.NewMockIFileHeader(ctrl)
		mf := mock_iservices.NewMockIFile(ctrl)

		mfh.EXPECT().Open().Return(mf, nil).Times(1)
		mf.EXPECT().Close().Do(f.Close).Times(1)
		mf.EXPECT().Read(gomock.Any()).DoAndReturn(f.Read).AnyTimes()

		err = ts.service.UploadCsvFile(ts.ctx, mfh)

		var ewm *errs.ErrorWithMessage
		ts.Require().ErrorAs(err, &ewm)
		ts.Require().ErrorIs(ewm.Err, ErrTransactionExist)
		ts.Require().Contains(ewm.Msg, "transaction insert error key (transaction_id)=(1) already exists.")

		ts.CleanupDB()
	})

	ts.Run("transactions is uploaded", func() {
		ctrl := gomock.NewController(ts.T())
		defer ctrl.Finish()

		f, err := fixtures.Fixtures.Open("transactions_test.csv")
		ts.Require().NoError(err)

		mfh := mock_iservices.NewMockIFileHeader(ctrl)
		mf := mock_iservices.NewMockIFile(ctrl)

		mfh.EXPECT().Open().Return(mf, nil).Times(1)
		mf.EXPECT().Close().Do(f.Close).Times(1)
		mf.EXPECT().Read(gomock.Any()).DoAndReturn(f.Read).AnyTimes()

		ts.Require().NoError(ts.service.UploadCsvFile(ts.ctx, mfh))

		storedT, err := ts.service.GetFilteredTransactions(ts.ctx, schemes.RawTransactionFilter{})
		ts.Require().NoError(err)

		ts.Assert().Len(storedT, countTestTransactions)

		ts.CleanupDB()
	})
}

func (ts *ServiceTestSuite) TestGetFilteredTransactions() {
	ts.SeedDB()

	ts.Run("not filter", func() {
		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, schemes.RawTransactionFilter{})
		ts.Require().NoError(err)

		ts.Assert().Len(sliceT, countTestTransactions)
	})

	ts.Run("get transaction by id", func() {
		rtf := schemes.RawTransactionFilter{
			ID: "1",
		}

		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		ts.Assert().Len(sliceT, 1)
		ts.Assert().Equal(int32(1), sliceT[0].TransactionID)
	})

	ts.Run("get transaction by status", func() {
		accepted, declined := schemes.Statuses[0], schemes.Statuses[1]

		rtf := schemes.RawTransactionFilter{
			Status: accepted,
		}

		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		for _, t := range sliceT {
			ts.Assert().Equal(accepted, t.Status)
		}

		rtf.Status = declined

		sliceT, err = ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		for _, t := range sliceT {
			ts.Assert().Equal(declined, t.Status)
		}
	})

	ts.Run("get transaction by payment types", func() {
		cash, card := schemes.PaymentTypes[0], schemes.PaymentTypes[1]

		rtf := schemes.RawTransactionFilter{
			PaymentType: cash,
		}

		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		for _, t := range sliceT {
			ts.Assert().Equal(cash, t.PaymentType)
		}

		rtf.PaymentType = card

		sliceT, err = ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		for _, t := range sliceT {
			ts.Assert().Equal(card, t.PaymentType)
		}
	})

	ts.Run("get transaction by terminal IDs", func() {
		terminalIDs := []int32{3506, 3507, 3508, 3509}
		strTerminalIDs := make([]string, 0, len(terminalIDs))
		for _, tID := range terminalIDs {
			strTerminalIDs = append(strTerminalIDs, strconv.FormatInt(int64(tID), 10))
		}

		rtf := schemes.RawTransactionFilter{
			TerminalIDs: strTerminalIDs,
		}

		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		for _, t := range sliceT {
			_, err := comparator.IdxSlice(t.TerminalID, terminalIDs)
			ts.Require().NoError(err)
		}
	})

	ts.Run("get transaction by date post", func() {

		rtf := schemes.RawTransactionFilter{
			DatePost: schemes.RawDatePost{
				From: "2022-08-13",
				To:   "",
			},
		}

		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		ts.Assert().Len(sliceT, countFromThirteenth)

		rtf.DatePost = schemes.RawDatePost{
			From: "",
			To:   "2022-08-13",
		}

		sliceT, err = ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		ts.Assert().Len(sliceT, countToThirteenth)

		rtf.DatePost = schemes.RawDatePost{
			From: "2022-08-11",
			To:   "2022-08-13",
		}

		sliceT, err = ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		ts.Assert().Len(sliceT, countToThirteenth)
	})

	ts.Run("get transaction by date post", func() {
		PaymentNarrative := "А11/27122 від 19.11.2020 р."

		rtf := schemes.RawTransactionFilter{
			PaymentNarrative: PaymentNarrative,
		}

		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, rtf)
		ts.Require().NoError(err)

		for _, t := range sliceT {
			ts.Assert().Contains(t.PaymentNarrative, PaymentNarrative)
		}
	})

	ts.Run("not filter", func() {
		sliceT, err := ts.service.GetFilteredTransactions(ts.ctx, schemes.RawTransactionFilter{})
		ts.Require().NoError(err)

		ts.Assert().Len(sliceT, countTestTransactions)
	})

	ts.CleanupDB()
}

func (ts *ServiceTestSuite) TestDownloadCsvFile() {
	ts.SeedDB()

	b := new(bytes.Buffer)

	err := ts.service.DownloadCsvFile(ts.ctx, b)
	ts.Require().NoError(err)
	ts.Require().NotEqual(0, len(b.Bytes()))

	reader := csv.NewReader(b)
	_, err = reader.Read()
	ts.Require().NoError(err)

	data, err := reader.ReadAll()
	ts.Require().NoError(err)

	ts.Assert().Len(data, countTestTransactions)

	ts.CleanupDB()
}
