package types

//go:generate mockgen -destination=./mocks/mock_reader.go -source=reader.go -package=mocks . Readable

type Readable interface {
	Close() error
	GetReadURLs() int32
}
