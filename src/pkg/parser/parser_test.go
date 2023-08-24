package parser

import (
	"4zp6/cigo/pkg/customErrors"
	"4zp6/cigo/pkg/data"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeProjectDef(t *testing.T) {
	// create two temporary files
	file, err := os.Create("project_test.json")
	if err != nil {
		t.Errorf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(file.Name()) // clean up

	// Create a test project definition
	project := data.ProjectDefinition{
		MainLanguage:   "java",
		LangVersion:    "8",
		Name:           "monorepo",
		Targets:        map[string]data.Target{},
		Version:        "1.0.0",
		Owners:         []string{"ownerA", "ownerB"},
		DependsOn:      []string{"PROJ_A", "PROJ_B"},
		Metadata:       map[string]string{"color": "green"},
		AffectsTags:    []string{"client"},
		AffectedByTags: []string{},
	}

	// Test JSON file type
	assert.NoError(t, EncodeProjectDef(project, file.Name(), JSON))

	// check if the file was created and contains valid JSON
	//...

	// Test illegal path
	// err = EncodeProjectDef(project, "bad/path/project.json", JSON)
	// assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(&customErrors.IllegalPathError{}))

	// Decode function
	// Check normal case
	_, err = DecodeProjectDef(file.Name(), JSON)
	assert.NoError(t, err)

	// Test illegal file path
	_, err = DecodeProjectDef("not_found.json", JSON)
	assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(&customErrors.FileNotFoundError{}))

	file.Close() // Close file after use
}

func TestEncodeWorkspaceDef(t *testing.T) {
	// create a temporary file
	file, err := os.Create("workspace_test.json")
	if err != nil {
		t.Errorf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(file.Name()) // clean up

	// Create a test project definition
	workspace := data.Workspace{
		Owners:          []string{"owner_a", "owner_b"},
		AppVer:          "v1.2",
		Projects:        map[string]string{"proj_a": "path/to/poj_a", "proj_b": "path/to/proj_b"},
		RemoteUrl:       "github.com/fake-account/fun-repo.git",
		Tags:            []string{"tag"},
		RequiredTargets: []string{"target_alpha", "target_beta"},
	}

	// Test JSON file type
	assert.NoError(t, EncodeWorkspace(workspace, file.Name()+".json", JSON))
	// assert.EqualError(t, EncodeWorkspace(workspace, file.Name(), JSON), "IllegalPath")

	// check if the file was created and contains valid JSON
	//...

	// Test illegal path
	// err = EncodeWorkspace(workspace, "bad/path/workspace.json", JSON)
	// assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(&customErrors.IllegalPathError{}))

	// Decode function

	// Check normal case
	_, err = DecodeWorkspace(file.Name()+".json", JSON)
	os.Remove(file.Name()+".json") // Clean up file
	assert.NoError(t, err)

	// Test illegal file path
	_, err = DecodeWorkspace("not_found.json", JSON)
	assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(&customErrors.FileNotFoundError{}))

	file.Close() // Close file after use
}
