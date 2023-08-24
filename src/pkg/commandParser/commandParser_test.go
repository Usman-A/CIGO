package commandParser

import (
	"testing"

	"github.com/Potterli20/go-flags-fork"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	args := []string{
		"-d",
		"-V",
		"list",
	}
	parsed, err := Parse(args)
	assert.Nil(t, err, "Should not return any errors")

	assert.True(t, parsed.DryRun)
	assert.True(t, parsed.Version)
	assert.Equal(t, LIST, command)
}

func TestSearchCommand(t *testing.T) {
	args := []string{
		"search",
		"-l", "3",
		"-s", "key1:val1",
		"-s", "key2:val2",
	}
	parsed, err := Parse(args)

	assert.Nil(t, err)
	assert.Equal(t, SEARCH, command)
	assert.Equal(t, 3, parsed.Search.Limit, "Limit should be three")
	assert.Equal(t, 2, len(parsed.Search.SearchItems), "Should have 2 search items")
	assert.Equal(t, "val1", parsed.Search.SearchItems["key1"], "`key1` should equal `val1`")
	assert.Equal(t, "val2", parsed.Search.SearchItems["key2"], "`key2` should equal `val2`")
}

func TestRunCommand(t *testing.T) {
	args := []string{
		"run",
		"-t",
		"target",
		"-p",
		"proj1",
	}
	parsed, err := Parse(args)

	assert.Nil(t, err)
	assert.Equal(t, RUN, GetCommand())
	assert.Equal(t, "proj1", parsed.Run.Positional.Project)
	assert.Equal(t, "target", parsed.Run.Positional.Target)
}

func TestChangesCommand(test *testing.T) {
	args := []string{
		"get-changed",
		"-B",
		"main",
	}

	parsed, err := Parse(args)

	assert.Nil(test, err)
	assert.Equal(test, "HEAD", parsed.GetChanged.Positional.HeadRef)
	assert.Equal(test, "main", parsed.GetChanged.Positional.BaseRef)
}

func TestAddProjectCommand(test *testing.T) {
	args := []string{
		"add-project",
	}

	parsed, err := Parse(args)

	assert.Nil(test, err)
	assert.Equal(test, ADDPROJECT, GetCommand())
	assert.NotNil(test, parsed.AddProject)
}

func TestCreateSchemaCommand(test *testing.T) {
	args := []string{
		"create-schema",
		"-t",
		"workspace",
	}

	parsed, err := Parse(args)

	assert.Nil(test, err)
	assert.Equal(test, CREATESCHEMA, GetCommand())
	assert.NotNil(test, parsed.CreateSchema)
	assert.Equal(test, "workspace", parsed.CreateSchema.Type)
}

func TestBadArgs(t *testing.T) {
	args := []string{
		"--bad-flag",
	}
	parsed, err := Parse(args)
	assert.Nil(t, parsed)

	assert.NotNil(t, err)
}

func TestPrintHelp(t *testing.T) {
	args := []string{
		"get-changed",
		"-h",
	}
	_, err := Parse(args)
	assert.NotNil(t, err.Error())
	assert.True(t, flags.WroteHelp(err))
}
