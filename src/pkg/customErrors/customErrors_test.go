package customErrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIllegalPath(t *testing.T) {
	err := IllegalPathError{Message: "Illegal path"}
	assert.Equal(t, "Illegal path", err.Error())
}

func TestFileNotFound(t *testing.T) {
	err := FileNotFoundError{Message: "File not found"}
	assert.Equal(t, "File not found", err.Error())
}

func TestInvalidFormat(t *testing.T) {
	err := InvalidFormatError{Message: "Invalid format"}
	assert.Equal(t, "Invalid format", err.Error())
}

func TestEncodingError(t *testing.T) {
	err := EncodingError{Message: "Encoding error"}
	assert.Equal(t, "Encoding error", err.Error())
}

func TestIOError(t *testing.T) {
	err := IOError{Message: "IO error"}
	assert.Equal(t, "IO error", err.Error())
}
