package service

import "mime/multipart"

type IFile interface {
	multipart.File
}

type IFileHeader interface {
	Open() (multipart.File, error)
}

//go:generate mockgen -destination=./mock/transactions.go . IFileHeader,IFile
