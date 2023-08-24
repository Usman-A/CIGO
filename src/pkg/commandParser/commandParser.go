package commandParser

import (
	flags "github.com/Potterli20/go-flags-fork"
)

type Command string

const (
	SEARCH       Command = "Search"
	LIST         Command = "List"
	RUN          Command = "Run"
	GETCHANGES   Command = "GetChanges"
	CREATESCHEMA Command = "CreateSchema"
	ADDPROJECT   Command = "AddProject"
)

var command Command = ""

type Args struct {
	DryRun       bool         `short:"d" long:"dry" description:"Print commands to run in order without running anything."`
	Version      bool         `short:"V" long:"version" description:"Print version"`
	List         List         `command:"list"`
	Search       Search       `command:"search"`
	Run          Run          `command:"run"`
	GetChanged   GetChanged   `command:"get-changed"`
	CreateSchema CreateSchema `command:"create-schema"`
	AddProject   AddProject   `command:"add-project"`
}

type (
	List   struct{}
	Search struct {
		Limit       int               `short:"l" description:"Search limit" default:"10"`
		SearchItems map[string]string `short:"s" description:"Search key value pair" required:"true"`
	}
	Run struct {
		Positional struct {
			Target  string `short:"t" long:"target" description:"Target you would like to run" required:"yes"`
			Project string `short:"p" long:"project" description:"The project that the target belongs to" required:"yes"`
		}
	}
	GetChanged struct {
		Positional struct {
			BaseRef string `short:"B" long:"base" description:"The base git reference/hash or branch to compare against" required:"true"`
			HeadRef string `short:"H" long:"head" description:"The head git reference/hash or branch to compare. Optional, defaults to HEAD"`
		}
	}
	// nolint:staticcheck
	CreateSchema struct {
		Type string `short:"t" long:"schema-type" choice:"workspace" choice:"project" required:"yes" positional:"yes"`
	}
	AddProject struct{}
)

// This function is called when the subcommand is called in the command line.
// See https://pkg.go.dev/github.com/Potterli20/go-flags-fork#Commander
// It is used to set the command that is active because checking the populated
// structs is not reliable or requires a lot of code.
func (l *List) Execute(args []string) error {
	command = LIST
	return nil
}
func (l *Search) Execute(args []string) error {
	command = SEARCH
	return nil
}
func (l *Run) Execute(args []string) error {
	command = RUN
	return nil
}
func (l *GetChanged) Execute(args []string) error {
	command = GETCHANGES
	if l.Positional.HeadRef == "" {
		l.Positional.HeadRef = "HEAD"
	}
	return nil
}
func (l *CreateSchema) Execute(args []string) error {
	command = CREATESCHEMA
	return nil
}

func (l *AddProject) Execute(args []string) error {
	command = ADDPROJECT
	return nil
}

func Parse(args []string) (*Args, error) {
	var parsedArgs Args
	_, err := flags.ParseArgs(&parsedArgs, args)
	if err != nil {
		return nil, err
	}
	return &parsedArgs, nil
}

func GetCommand() Command {
	return command
}
