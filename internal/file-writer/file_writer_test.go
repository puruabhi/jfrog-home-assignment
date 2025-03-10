package filewriter

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/puruabhi/jfrog/home-assignment/internal/config"
	"github.com/puruabhi/jfrog/home-assignment/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewFileWriter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := config.WriteConfig{
		WriteDir: "test-dir",
	}
	logger := types.NewLoggerStub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writer := NewFileWriter(ctx, mockConfig, logger)
	assert.NotNil(t, writer)
	assert.Equal(t, mockConfig, writer.config)
	assert.Equal(t, logger, writer.logger)
	assert.Equal(t, ctx, writer.ctx)
}

func TestWriter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := config.WriteConfig{
		WriteDir: "test-dir",
	}
	logger := types.NewLoggerStub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writer := NewFileWriter(ctx, mockConfig, logger)

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	writer.config.WriteDir = tempDir

	// Send some data to the write
	data := []byte("test data")
	writer.PushForWrite(data)

	// Allow some time for the writer goroutine to process the data
	time.Sleep(10 * time.Millisecond)

	// Check that the file was written
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	// Read the file content
	filePath := filepath.Join(tempDir, files[0].Name())
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, data, content)
}

func TestWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := config.WriteConfig{
		WriteDir: "test-dir",
	}
	logger := types.NewLoggerStub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writer := NewFileWriter(ctx, mockConfig, logger)

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	writer.config.WriteDir = tempDir

	// Write some data
	data := []byte("test data")
	writer.write(data)

	// Check that the file was written
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	// Read the file content
	filePath := filepath.Join(tempDir, files[0].Name())
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, data, content)
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := config.WriteConfig{
		WriteDir: "test-dir",
	}
	logger := types.NewLoggerStub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writer := NewFileWriter(ctx, mockConfig, logger)

	// Create a channel to signal when the writer goroutine has finished
	done := make(chan struct{})

	go func() {
		writer.writer()
		close(done)
	}()

	// Cancel the context to stop the writer goroutine
	cancel()

	// Wait for the writer goroutine to finish
	<-done

	// Check that the writeChan is closed
	_, ok := <-writer.writeChan
	assert.False(t, ok)
}
