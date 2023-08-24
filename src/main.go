package main

import (
	"4zp6/cigo/pkg/commandParser"
	"4zp6/cigo/pkg/commands"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/tcnksm/go-input"
)

func main() {
	args, err := commandParser.Parse(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	switch commandParser.GetCommand() {
	case commandParser.LIST:
		err := commands.List()
		if err != nil {
			fmt.Printf("Faced error while running List: %v\n", err)
		}
	case commandParser.SEARCH:
		projs, err := commands.Search(args.Search.SearchItems)
		if err != nil {
			fmt.Printf("Faced error: %v\n", err)
			return
		}
		for _, p := range projs {
			fmt.Println(p.Name)
		}
	case commandParser.CREATESCHEMA:
		err := commands.CreateSchema(args.CreateSchema.Type)
		if err != nil {
			fmt.Printf("Failed to create schema for %v: %v\n", args.CreateSchema.Type, err)
		} else {
			fmt.Printf("Created schema for %v\n", args.CreateSchema.Type)
		}
	case commandParser.RUN:
		err := commands.Run(args.Run.Positional.Project, args.Run.Positional.Target, args.DryRun)
		if err != nil {
			fmt.Println(err)
		}
	case commandParser.GETCHANGES:
		projs, err := commands.GetAffected(args.GetChanged.Positional.BaseRef, args.GetChanged.Positional.HeadRef)
		if err != nil {
			color.Errorln("Failed to get affected projects", err)
		} else if len(projs) == 0 {
			color.Infoln("No projects affected...")
		} else {
			color.Infoln("Affected projects:")
			for _, p := range projs {
				color.Infoln("\t", p)
			}
		}
	case commandParser.ADDPROJECT:
		ui := &input.UI{}
		err := commands.AddProject(ui)
		if err != nil {
			fmt.Println(err)
		}
	}
}
