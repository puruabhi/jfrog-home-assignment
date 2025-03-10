package csvreader

import (
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/puruabhi/jfrog/home-assignment/internal/csv-reader/mocks"
	"github.com/puruabhi/jfrog/home-assignment/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestFetchURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVReader := mocks.NewMockCSVReadable(ctrl)
	mockFileReader := mocks.NewMockFileReadable(ctrl)
	logger := types.NewLoggerStub()
	urlChan := make(chan string, 10)

	csv := &csvReader{
		reader:     mockCSVReader,
		fileReader: mockFileReader,
		logger:     logger,
		urls:       urlChan,
	}

	// Mock the read method to return a header and then URLs
	gomock.InOrder(
		mockCSVReader.EXPECT().Read().Return([]string{"Urls"}, nil),
		mockCSVReader.EXPECT().Read().Return([]string{"www.example.com"}, nil),
		mockCSVReader.EXPECT().Read().Return([]string{"www.someotherurl.com/api/v1"}, nil),
		mockCSVReader.EXPECT().Read().Return([]string{"www.anotherone.com"}, nil),
		mockCSVReader.EXPECT().Read().Return(nil, io.EOF),
	)

	go csv.fetchURLs()

	// Collect URLs from the channel
	var urls []string
	for url := range urlChan {
		urls = append(urls, url)
	}

	expectedURLs := []string{
		"www.example.com",
		"www.someotherurl.com/api/v1",
		"www.anotherone.com",
	}

	assert.Equal(t, expectedURLs, urls)
}

func TestFetchURLs_ErrorReadingHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVReader := mocks.NewMockCSVReadable(ctrl)
	mockFileReader := mocks.NewMockFileReadable(ctrl)
	logger := types.NewLoggerStub()
	urlChan := make(chan string, 10)

	csv := &csvReader{
		reader:     mockCSVReader,
		fileReader: mockFileReader,
		logger:     logger,
		urls:       urlChan,
	}

	mockCSVReader.EXPECT().Read().Return(nil, errors.New("error reading header"))

	go csv.fetchURLs()

	// Ensure the channel is closed
	_, ok := <-urlChan
	assert.False(t, ok)
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileReader := mocks.NewMockFileReadable(ctrl)
	logger := types.NewLoggerStub()

	csv := &csvReader{
		fileReader: mockFileReader,
		logger:     logger,
	}

	mockFileReader.EXPECT().Close().Return(nil)

	err := csv.Close()
	assert.NoError(t, err)
}

func TestClose_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileReader := mocks.NewMockFileReadable(ctrl)
	logger := types.NewLoggerStub()

	csv := &csvReader{
		fileReader: mockFileReader,
		logger:     logger,
	}

	mockFileReader.EXPECT().Close().Return(errors.New("error closing file"))

	err := csv.Close()
	assert.Error(t, err)
	assert.Equal(t, "error closing file", err.Error())
}
