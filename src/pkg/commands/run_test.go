package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// test setup
	outTemp, errTemp := os.Stdout, os.Stderr
	defer func() {
		os.Stdout = outTemp
		os.Stderr = errTemp
	}()

	out, err := os.CreateTemp("/tmp", "out.XXX.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(out.Name())
	errFile, err := os.CreateTemp("/tmp", "err.XXX.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(errFile.Name())

	os.Stdout, os.Stderr = out, errFile

	err = Run("proj_a", "test", false)
	if err != nil {
		t.Fatal(err)
	}
	outData, err := os.ReadFile(out.Name())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `initializing proj_a
building application
done building proj_a
testing proj_a
`, string(outData))
	errData, err := os.ReadFile(errFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	assert.Empty(t, errData, "There should no error messages")
}
