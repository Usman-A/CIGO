package commands

import (
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"fmt"
	"strings"
)

// Prints the list to the output directly
func List() error {
	// read the workspace directory
	root, err := misc.GetRoot()
	if err != nil {
		return err
	}

	workspace, err := parser.DecodeWorkspace(strings.TrimSpace(root)+"/workspace.json", parser.JSON)
	if err != nil {
		return err
	}

	fmt.Println("\nProjects in this workspace:")
	for name := range workspace.Projects {
		fmt.Printf("\t%s\n", name)
	}

	return nil
}
