package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/puruabhi/jfrog/home-assignment/internal/types"
)

const (
	ParallelDownload = 50
	HTTPPrefix       = "http"
	HTTPSPrefix      = "https://"
)

type downloader struct {
	ctx    context.Context
	logger types.Logger
	reader types.Readable
	writer types.Writable
	finish chan struct{}
	urls   chan string
	lock   chan struct{}

	stats struct {
		activeDownloads    atomic.Int32
		downloadSuccessful atomic.Int32
		downloadFailed     atomic.Int32
	}
}

// NewDownloader initializes a new downloader instance and starts processing and monitoring routines.
func NewDownloader(ctx context.Context, logger types.Logger, reader types.Readable, writer types.Writable) *downloader {
	down := &downloader{
		ctx:    ctx,
		logger: logger,
		reader: reader,
		writer: writer,
		finish: make(chan struct{}),
		urls:   make(chan string),
		lock:   make(chan struct{}, ParallelDownload),
	}

	go down.startProcessing()

	down.logger.Infof("Downloader started")
	return down
}

// formatURL ensures the URL has the correct prefix (http or https).
func (d *downloader) formatURL(url string) string {
	if !strings.HasPrefix(url, HTTPPrefix) {
		url = fmt.Sprintf("%s%s", HTTPSPrefix, url)
	}
	return url
}

// fetchContent retrieves the content from the given URL.
func (d *downloader) fetchContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response from URL %s: %s", url, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// download formats the URL and fetches its content.
func (d *downloader) download(url string) ([]byte, error) {
	formattedURL := d.formatURL(url)
	return d.fetchContent(formattedURL)
}

// downloadAndPush downloads the content from the URL and pushes it to the writer.
func (d *downloader) downloadAndPush(url string, wg *sync.WaitGroup) {
	defer func() {
		<-d.lock
		wg.Done()
		d.stats.activeDownloads.Add(-1)
	}()

	d.stats.activeDownloads.Add(1) // Increment the counter

	content, err := d.download(url)
	if err != nil {
		d.logger.Debugf("Error downloading URL: %s - %s", url, err)
		d.stats.downloadFailed.Add(1)
		return
	}

	d.stats.downloadSuccessful.Add(1)
	d.writer.PushForWrite(content)
}

// downloadWorker processes URLs from the channel and starts downloadAndPush for each URL.
func (d *downloader) downloadWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	downloadWG := sync.WaitGroup{}
	defer downloadWG.Wait()

	for {
		select {
		case <-d.ctx.Done():
			return
		case url, ok := <-d.urls:
			if !ok {
				return
			}

			d.lock <- struct{}{}
			downloadWG.Add(1)
			go d.downloadAndPush(url, &downloadWG)
		}
	}
}

// startProcessing starts the download worker and waits for it to finish.
func (d *downloader) startProcessing() {
	defer d.finishProcessing()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go d.downloadWorker(&wg)
	wg.Wait()
}

// finishProcessing signals that processing is finished.
func (d *downloader) finishProcessing() {
	d.finish <- struct{}{}
}

// GetFinishChan returns the finish channel.
func (d *downloader) GetFinishChan() chan struct{} {
	return d.finish
}

// GetURLsChan returns the URLs channel.
func (d *downloader) GetURLsChan() chan string {
	return d.urls
}

func (d *downloader) GetStats() any {
	type stats struct {
		ActiveDownloads    int32 `json:"active_downloads"`
		DownloadSuccessful int32 `json:"download_successful"`
		DownloadFailed     int32 `json:"download_failed"`
	}

	return stats{
		ActiveDownloads:    d.stats.activeDownloads.Load(),
		DownloadSuccessful: d.stats.downloadSuccessful.Load(),
		DownloadFailed:     d.stats.downloadFailed.Load(),
	}
}
