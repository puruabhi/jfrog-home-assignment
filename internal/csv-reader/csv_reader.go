package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/puruabhi/jfrog/home-assignment/internal/config"
	"github.com/puruabhi/jfrog/home-assignment/internal/types"
)

//go:generate mockgen -destination=./mocks/mock_csv_reader.go -source=csv_reader.go -package=mocks .
type CSVReadable interface {
	Read() (record []string, err error)
}

type FileReadable interface {
	Close() error
}

type csvReader struct {
	config     config.ReadConfig
	reader     CSVReadable
	fileReader FileReadable
	logger     types.Logger
	urls       chan string
	readUrls   int32
}

// NewCSVReader initializes a new csvReader instance and starts fetching URLs.
func NewCSVReader(config config.ReadConfig, logger types.Logger, urlChan chan string) (*csvReader, error) {
	fileReader, err := os.Open(config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("caught err while opening file: %w", err)
	}

	csvFileReader := csv.NewReader(fileReader)
	csv := &csvReader{
		config:     config,
		reader:     csvFileReader,
		fileReader: fileReader,
		logger:     logger,
		urls:       urlChan,
	}

	go csv.fetchURLs()

	csv.logger.Infof("CSV reader started")
	return csv, nil
}

// fetchURLs reads URLs from the CSV file and sends them to the urls channel.
func (r *csvReader) fetchURLs() {
	defer close(r.urls)

	header, err := r.read()
	if err != nil {
		r.logger.Errorf("Error reading csv file: %s", err)
		return
	}
	r.logger.Debugf("CSV header: %+v\n", header)

	for {
		url, err := r.read()
		if err != nil {
			if err != io.EOF {
				r.logger.Errorf("Error reading csv file: %s", err)
			}
			return
		}

		r.readUrls++
		r.logger.Debugf("URL: %s\n", url[0])
		r.urls <- url[0]
	}
}

// read reads a single record from the CSV file.
func (r *csvReader) read() ([]string, error) {
	if r.reader == nil {
		return nil, fmt.Errorf("csv reader is not initialized")
	}
	return r.reader.Read()
}

// Close closes the underlying file reader.
func (r *csvReader) Close() error {
	if r.fileReader == nil {
		return fmt.Errorf("file reader is not initialized")
	}
	return r.fileReader.Close()
}

func (r *csvReader) GetReadURLs() int32 {
	return r.readUrls
}
