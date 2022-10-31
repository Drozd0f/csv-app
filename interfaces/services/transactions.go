package iservices

import "mime/multipart"

type IFile interface {
	multipart.File
}

type IFileHeader interface {
	Open() (multipart.File, error)
}

//go:generate mockgen -destination=./mock/transactions.go -package=mock_iservices . IFileHeader,IFile
