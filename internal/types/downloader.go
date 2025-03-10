package types

//go:generate mockgen -destination=./mocks/mock_downloader.go -source=downloader.go -package=mocks . Downloadable

type Downloadable interface {
	GetFinishChan() chan struct{}
	GetURLsChan() chan string
	GetStats() any
}
