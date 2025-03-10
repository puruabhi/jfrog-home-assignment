package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/puruabhi/jfrog/home-assignment/internal/types"
)

const (
	ParallelDownload = 50
	HTTPPrefix       = "http"
	HTTPSPrefix      = "https://"
)

type downloader struct {
	ctx           context.Context
	logger        types.Logger
	reader        types.Readable
	writer        types.Writable
	finish        chan struct{}
	urls          chan string
	lock          chan struct{}
	activeCounter int32 // Atomic counter for active downloads
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
	go down.monitorActiveDownloads()
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
		atomic.AddInt32(&d.activeCounter, -1) // Decrement the counter
	}()

	atomic.AddInt32(&d.activeCounter, 1) // Increment the counter

	content, err := d.download(url)
	if err != nil {
		d.logger.Errorf("Error downloading URL: %s - %s\n", url, err)
		return
	}
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

// GetActiveCounter returns the current count of active downloads.
func (d *downloader) GetActiveCounter() int32 {
	return atomic.LoadInt32(&d.activeCounter)
}

// monitorActiveDownloads periodically logs the number of active downloads.
func (d *downloader) monitorActiveDownloads() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			activeDownloads := d.GetActiveCounter()
			d.logger.Infof("Active downloads: %d", activeDownloads)
		}
	}
}
