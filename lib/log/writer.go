package log

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	libruntime "xcore/lib/runtime"
)

// 生成 normal log Writer
func newNormalFileWriter(filePath string, namePrefix string, logDuration int) (*os.File, error) {
	return newFileWriter(filePath, namePrefix, logDuration, normalLogFileBaseName)
}

// 生成 error log Writer
func newErrorFileWriter(filePath string, namePrefix string, logDuration int) (*os.File, error) {
	return newFileWriter(filePath, namePrefix, logDuration, errorLogFileBaseName)
}

// 生成 log Writer
func newFileWriter(filePath string, namePrefix string, logDuration int, fileBaseName string) (*os.File, error) {
	fileName := fmt.Sprintf(fileFormat, filePath, namePrefix, logDuration, fileBaseName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return nil, errors.WithMessage(err, libruntime.Location())
	}
	return file, nil
}
