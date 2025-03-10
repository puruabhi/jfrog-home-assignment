package filewriter

import (
	"context"
	"fmt"
	"os"
	"path"

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
	fileName := fmt.Sprintf("%s%s", uuid.New(), fileExt)
	filePath := path.Join(w.config.WriteDir, fileName)

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		w.logger.Errorf("Failed to save file: %s\n", err)
	} else {
		w.logger.Infof("Saved: %s\n", filePath)
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
