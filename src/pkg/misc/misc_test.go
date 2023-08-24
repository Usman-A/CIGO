package misc

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRoot(t *testing.T) {
	root, err := GetRoot()
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, root)
	_, err = os.Stat(root)
	assert.Nil(t, err)
}

func TestGetWorkspacePath(t *testing.T) {
	workspacePath, err := GetWorkspacePath()
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, workspacePath)
	//make sure the suffix is workspace.json
	assert.Regexp(t, `workspace\.json$`, workspacePath)
}

func TestGetRelativePath(t *testing.T) {
	relativePath, err := GetRelativePath("src/pkg/misc/misc.go")
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, relativePath)
	fmt.Println(relativePath)
	//make sure the suffix is monorepo/src/pkg/misc/misc.go
	assert.Regexp(t, `monorepo\/src\/pkg\/misc\/misc\.go$`, relativePath)
}

func TestContains(t *testing.T) {
	list := []string{"a", "b", "c"}
	assert.True(t, Contains(list, "a"))
	assert.False(t, Contains(list, "d"))
}

func TestContainsNotIn(t *testing.T) {
	list := []string{"a", "b", "c"}
	assert.False(t, Contains(list, "d"))
}

func TestGetProjects(t *testing.T) {
	projects, err := GetProjects()
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, projects)
	assert.NotEmpty(t, projects)

	// Checking if projects are in the map
	_, ok := projects["proj_a"]
	assert.True(t, ok)
	_, ok = projects["proj_b"]
	assert.True(t, ok)
	_, ok = projects["proj_c"]
	assert.True(t, ok)
	_, ok = projects["proj_d"]
	assert.True(t, ok)

	_, ok = projects["proj_e"]
	assert.False(t, ok)
}
