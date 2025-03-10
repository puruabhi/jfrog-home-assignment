package process

import (
	"context"
	"time"

	"github.com/puruabhi/jfrog/home-assignment/internal/config"
	csvreader "github.com/puruabhi/jfrog/home-assignment/internal/csv-reader"
	"github.com/puruabhi/jfrog/home-assignment/internal/downloader"
	filewriter "github.com/puruabhi/jfrog/home-assignment/internal/file-writer"
	"github.com/puruabhi/jfrog/home-assignment/internal/logger"
	"github.com/puruabhi/jfrog/home-assignment/internal/types"
)

type process struct {
	finish     chan struct{}
	logger     types.Logger
	csvReader  types.Readable
	downloader types.Downloadable
	writer     types.Writable
	config     *config.Config
	ctx        context.Context
}

// Setup initializes the process and starts the setup routine.
func Setup() (chan struct{}, error) {
	finished := make(chan struct{})
	log := logger.NewZapLogger()
	prc := &process{
		finish: finished,
		logger: log,
		ctx:    context.Background(),
	}

	go prc.setup()
	go prc.printPeriodicStats()
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

	prc.writer = filewriter.NewFileWriter(prc.ctx, prc.config.Write, prc.logger)
	prc.downloader = downloader.NewDownloader(prc.ctx, prc.logger, prc.csvReader, prc.writer)

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

func (prc *process) printPeriodicStats() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			prc.printCSVReaderStats()
			prc.printDownloaderStats()
			prc.printWriterStats()

		case <-prc.ctx.Done():
			return
		}
	}
}

func (prc *process) printCSVReaderStats() {
	readUrls := prc.csvReader.GetReadURLs()
	prc.logger.Infof("CSV Reader: urls read: %d", readUrls)
}

func (prc *process) printDownloaderStats() {
	stats := prc.downloader.GetStats()
	prc.logger.Infof("URL Downloader: %+v", stats)
}

func (prc *process) printWriterStats() {
	stats := prc.writer.GetStats()
	prc.logger.Infof("File Writer: %+v", stats)
}
