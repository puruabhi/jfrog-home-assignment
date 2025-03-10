package filewriter

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/puruabhi/jfrog/home-assignment/internal/config"
	"github.com/puruabhi/jfrog/home-assignment/internal/types"
)

const fileExt = ".txt"

type fileWriter struct {
	ctx       context.Context
	config    config.WriteConfig
	logger    types.Logger
	writeChan chan []byte

	// stat variables
	stats struct {
		writeFailed  atomic.Int32
		writing      atomic.Int32
		writeSuccess atomic.Int32
	}
}

// NewFileWriter initializes a new fileWriter instance and starts the writer goroutine.
func NewFileWriter(ctx context.Context, config config.WriteConfig, logger types.Logger) *fileWriter {
	writer := &fileWriter{
		config:    config,
		logger:    logger,
		ctx:       ctx,
		writeChan: make(chan []byte, 100),
	}

	go writer.writer()

	writer.logger.Infof("File writer started")
	return writer
}

// writer listens for bytes on the writeChan and writes them to files.
func (w *fileWriter) writer() {
	defer w.logger.Infof("File writer stopped\n")
	defer w.close()

	for {
		select {
		case bytes, ok := <-w.writeChan:
			if !ok {
				return
			}
			w.write(bytes)
		case <-w.ctx.Done():
			return
		}
	}
}

// write writes the given bytes to a new file.
func (w *fileWriter) write(bytes []byte) {
	w.stats.writing.Add(1)
	defer w.stats.writing.Add(-1)

	fileName := fmt.Sprintf("%s%s", uuid.New(), fileExt)
	filePath := path.Join(w.config.WriteDir, fileName)

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		w.logger.Errorf("Failed to save file: %s\n", err)
		w.stats.writeFailed.Add(1)
	} else {
		w.logger.Debugf("Saved: %s\n", filePath)
		w.stats.writeSuccess.Add(1)
	}
}

// PushForWrite sends bytes to the writeChan for writing to a file.
func (w *fileWriter) PushForWrite(bytes []byte) {
	w.writeChan <- bytes
}

// close closes the writeChan.
func (w *fileWriter) close() {
	close(w.writeChan)
}

func (w *fileWriter) GetStats() any {
	type stats struct {
		WriteFailed  int32 `json:"write_failed"`
		Writing      int32 `json:"writing"`
		WriteSuccess int32 `json:"write_success"`
	}
	return stats{
		WriteFailed:  w.stats.writeFailed.Load(),
		Writing:      w.stats.writing.Load(),
		WriteSuccess: w.stats.writeSuccess.Load(),
	}
}
