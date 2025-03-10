package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/puruabhi/jfrog/home-assignment/internal/types"
	typeMocks "github.com/puruabhi/jfrog/home-assignment/internal/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFormatURL(t *testing.T) {
	logger := types.NewLoggerStub()
	d := &downloader{
		logger: logger,
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "https://example.com"},
		{"http://example.com", "http://example.com"},
		{"https://example.com", "https://example.com"},
	}

	for _, test := range tests {
		result := d.formatURL(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestFetchContent(t *testing.T) {
	logger := types.NewLoggerStub()
	d := &downloader{
		logger: logger,
	}

	serverMockResponse := "test content"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverMockResponse))
	}))
	defer server.Close()

	url := server.URL
	expectedContent := serverMockResponse

	content, err := d.fetchContent(url)
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, string(content))
}

func TestFetchContent_Error(t *testing.T) {
	logger := types.NewLoggerStub()
	d := &downloader{
		logger: logger,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	url := server.URL

	_, err := d.fetchContent(url)
	assert.Error(t, err)
}

func TestDownloadAndPush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serverMockResponse := "test content"

	mockWriter := typeMocks.NewMockWritable(ctrl)
	logger := types.NewLoggerStub()
	d := &downloader{
		logger: logger,
		writer: mockWriter,
		lock:   make(chan struct{}, 1),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverMockResponse))
	}))
	defer server.Close()

	url := server.URL
	content := []byte(serverMockResponse)

	mockWriter.EXPECT().PushForWrite(content).Times(1)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	d.lock <- struct{}{}
	go func() {
		d.downloadAndPush(url, wg)
	}()

	wg.Wait()
}

func TestDownloadWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := typeMocks.NewMockReadable(ctrl)
	mockWriter := typeMocks.NewMockWritable(ctrl)
	logger := types.NewLoggerStub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := &downloader{
		ctx:    ctx,
		logger: logger,
		reader: mockReader,
		writer: mockWriter,
		urls:   make(chan string, 1),
		lock:   make(chan struct{}, 1),
	}

	serverMockResponse := "test content"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverMockResponse))
	}))
	defer server.Close()

	url := server.URL
	content := []byte(serverMockResponse)

	mockWriter.EXPECT().PushForWrite(content).Times(1)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go d.downloadWorker(wg)

	d.urls <- url
	close(d.urls)

	wg.Wait()
}
