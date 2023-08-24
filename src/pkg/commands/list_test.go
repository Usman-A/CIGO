package commands

import (
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

// load workspace to use during tests
func init() {
	workspacePath, err := misc.GetWorkspacePath()
	if err != nil {
		fmt.Println(err)
	}
	workspace, err = parser.DecodeWorkspace(workspacePath, parser.JSON)
	if err != nil {
		fmt.Println("Error reading workspace", err)
	}
}

func TestList(t *testing.T) {

	// set a buffer to capture the output
	var buf bytes.Buffer

	// use pipe to capture stdout
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := List()
	if err != nil {
		t.Error("Error listing projects", err)
	}

	// read the output from the buffer
	w.Close()
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Error("Error reading from buffer", err)
	}

	parsed_output := strings.Split(strings.TrimSpace(buf.String()), "\n\t")[1:]

	if (len(parsed_output)) != len(workspace.Projects) {
		t.Error("List output should match the expected output")
	}

	for _, listed_proj := range parsed_output {
		fmt.Println(listed_proj)
		_, ok := workspace.Projects[listed_proj]
		if !ok {
			t.Error("Project listed is not in the workspace")
		}
	}
}
