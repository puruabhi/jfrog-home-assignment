package process

import (
	"context"
	"time"

	"github.com/puruabhi/jfrog/home-assignment/internal/config"
	csvreader "github.com/puruabhi/jfrog/home-assignment/internal/csv-reader"
	"github.com/puruabhi/jfrog/home-assignment/internal/downloader"
	filewriter "github.com/puruabhi/jfrog/home-assignment/internal/file-writer"
	"github.com/puruabhi/jfrog/home-assignment/internal/types"
)

type process struct {
	finish     chan struct{}
	logger     types.Logger
	csvReader  types.Readable
	downloader types.Downloadable
	writer     types.Writable
	config     *config.Config
}

// Setup initializes the process and starts the setup routine.
func Setup() (chan struct{}, error) {
	finished := make(chan struct{})
	logger := types.NewFmtLogger()
	prc := &process{
		finish: finished,
		logger: logger,
	}

	go prc.setup()
	return finished, nil
}

// setup configures the process components and starts the waitAndFinish routine.
func (prc *process) setup() error {
	defer prc.waitAndFinish()

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	prc.config = cfg

	ctx := context.Background()

	prc.writer = filewriter.NewFileWriter(ctx, prc.config.Write, prc.logger)
	prc.downloader = downloader.NewDownloader(ctx, prc.logger, prc.csvReader, prc.writer)

	csvReader, err := csvreader.NewCSVReader(prc.config.Read, prc.logger, prc.downloader.GetURLsChan())
	if err != nil {
		return err
	}
	prc.csvReader = csvReader

	return nil
}

// waitAndFinish waits for the downloader to finish and then closes the csvReader and finish channel.
func (prc *process) waitAndFinish() {
	<-prc.downloader.GetFinishChan()

	if prc.csvReader != nil {
		prc.csvReader.Close()
	}
	prc.finish <- struct{}{}

	time.AfterFunc(2*time.Second, func() {
		close(prc.finish)
	})
}
