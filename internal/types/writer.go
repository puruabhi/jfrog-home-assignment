package types

//go:generate mockgen -destination=./mocks/mock_writer.go -source=writer.go -package=mocks . Writable

type Writable interface {
	PushForWrite([]byte)
	GetStats() any
}
