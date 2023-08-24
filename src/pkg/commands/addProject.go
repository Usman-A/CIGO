package commands

import (
	"4zp6/cigo/pkg/customErrors"
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"errors"
	"fmt"
	"os"
	p "path"
	"strings"

	input "github.com/tcnksm/go-input"
)

var workspace *data.Workspace

func targetPrompt(user *input.UI, workspace data.Workspace) (map[string]data.Target, error) {
	targetMap := make(map[string]data.Target)

	create, err := user.Select("Do you want to add a target?", []string{"Yes", "No"}, &input.Options{
		Default:   "No",
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return nil, err
	}

	for create == "Yes" {
		builder := data.TargetBuilder{}

		// Get the name of the target
		name, err := user.Ask("What's the name of the target?", &input.Options{
			Default:   "target",
			Required:  true,
			Loop:      true,
			HideOrder: true,
		})
		if err != nil {
			return nil, err
		}

		// Get the commands the project has
		commands, err := user.Ask("What command(s) does the target have? Separate commands with a ','", &input.Options{
			Default:   "test build run",
			Required:  true,
			Loop:      true,
			HideOrder: true,
		})
		if err != nil {
			return nil, err
		}
		commandList := strings.Split(commands, ",")
		builder.SetCmds(commandList)

		artifacts, err := user.Ask("What artifact(s) does the target have? Separate items with a space.", &input.Options{
			Default:   "",
			Required:  true,
			Loop:      true,
			HideOrder: true,
		})
		if err != nil {
			return nil, err
		}
		artifactList := strings.Split(artifacts, " ")
		builder.SetArtifacts(artifactList)

		envMap := make(map[string]string)
		_, err = user.Ask("What environment variables does the target have? Use <key>>:<value> pairs separated by a colon", &input.Options{
			Default:   "",
			Required:  false,
			Loop:      true,
			HideOrder: true,
			ValidateFunc: func(metadata string) error {
				envList := strings.Split(metadata, " ")
				if len(envList) == 1 && envList[0] == "" {
					return nil
				}
				for _, item := range envList {
					split := strings.Split(item, ":")
					if len(split) != 2 {
						return fmt.Errorf("Invalid environment variable format, should be '<key>:<value>'")
					}
					envMap[split[0]] = split[1]
				}
				return nil
			},
		})
		if err != nil {
			return nil, err
		}
		builder.SetEnv(envMap)

		// get dependencies for the target
		d, err := user.Select("Does the target have any dependencies?", []string{"Yes", "No"}, &input.Options{
			Default:   "No",
			Required:  true,
			Loop:      true,
			HideOrder: true,
		})
		if err != nil {
			return nil, err
		}
		if d == "Yes" {
			dependencies := []data.DependsTarget{}
			for d == "Yes" {
				project, err := user.Ask("Which project does the dependant target belong to?", &input.Options{
					Default:   "",
					Required:  true,
					Loop:      true,
					HideOrder: true,
					ValidateFunc: func(project string) error {
						if _, ok := workspace.Projects[project]; !ok {
							return fmt.Errorf("Project '%s' cannot be found in the workspace", project)
						}
						return nil
					},
				})
				if err != nil {
					return nil, err
				}
				target, err := user.Ask("What is target does it depend on? (Enter the name)", &input.Options{
					Default:   "",
					Required:  true,
					Loop:      true,
					HideOrder: true,
					ValidateFunc: func(target string) error {
						pPath, err := misc.GetRelativePath(workspace.Projects[project] + "/project.json")
						if err != nil {
							return err
						}

						proj, err := parser.DecodeProjectDef(pPath, parser.JSON)
						if err != nil {
							return err
						}

						if _, ok := proj.Targets[target]; !ok {
							return fmt.Errorf("Target '%s' does not exist in project %s.", target, project)
						}
						return nil
					},
				})
				if err != nil {
					return nil, err
				}
				dependencies = append(dependencies, data.DependsTarget{Project: project, Target: target})

				d, err = user.Select("Do you want to add another dependency?", []string{"Yes", "No"}, &input.Options{
					Default:   "No",
					Required:  true,
					Loop:      true,
					HideOrder: true,
				})
				if err != nil {
					return nil, err
				}
			}
			builder.SetDependsOn(dependencies)
		}

		create, err = user.Select("Would you like to add another target?", []string{"Yes", "No"}, &input.Options{
			Default:   "No",
			Required:  true,
			Loop:      true,
			HideOrder: true,
		})
		if err != nil {
			return nil, err
		}
		targetMap[name] = builder.Build()
	}
	return targetMap, nil
}

func tagValidator(tag string) error {
	// if the tag is empty, then we don't need to validate it as it is an optional field
	if tag == "" {
		return nil
	}
	tagList := strings.Split(tag, " ")
	for _, t := range tagList {
		if contains(workspace.Tags, t) {
			continue
		} else {
			return fmt.Errorf("Tag '%s' does not exist", t)
		}
	}
	return nil
}

// Wizard to add a new project to the workspace
func AddProject(user *input.UI) error {
	projBuilder := data.ProjectDefinitionBuilder{}
	workspacePath, err := misc.GetWorkspacePath()
	if err != nil {
		return err
	}
	workspace, err = parser.DecodeWorkspace(workspacePath, parser.JSON)
	if err != nil {
		return err
	}
	skip := "\nIf you would like to skip this step press enter, otherwise enter a value"

	// Get the name of the project
	name, err := user.Ask("What's the name of the project?", &input.Options{
		Default:   "my_project",
		Required:  true,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(name string) error {
			// name can't end with or contain a space
			if strings.HasPrefix(name, " ") || strings.HasSuffix(name, " ") {
				return fmt.Errorf("Project name cannot start or end with a space")
			}
			if strings.Contains(name, " ") {
				return fmt.Errorf("Project name cannot contain a space")
			}
			if _, ok := workspace.Projects[name]; ok {
				return fmt.Errorf("Project with name %s already exists", name)
			}
			return nil
		},
	})
	if err != nil {
		return err
	}
	projBuilder.SetName(name)

	//Need the path of the project
	path, err := user.
		Ask("What's the path of where the project is stored?\n"+
			"The path would be the folder containing the project folder, e.g. 'apps/' "+
			"and within this folder you would have your project folder.",
			&input.Options{
				Default:   "apps/",
				Required:  true,
				Loop:      true,
				HideOrder: true,
				ValidateFunc: func(path string) error {
					//make sure it doesnt start with a /, this is because we will be appending a / to the path
					if strings.HasPrefix(path, "/") {
						return fmt.Errorf("Path cannot start with a '/'")
					}
					return nil
				},
			})
	if err != nil {
		return err
	}

	// Get Project version
	version, err := user.Ask("What's the version of the project?", &input.Options{
		Default:   "0.1.0",
		Required:  false,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return err
	}
	projBuilder.SetVersion(version)

	// Get the language name
	language, err := user.Ask("What's the language of the project? This information would help build the project later on.", &input.Options{
		Default:   "go",
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return err
	}
	projBuilder.SetMainLanguage(language)

	// Get language version
	langVersion, err := user.Ask("What's the version of the language?", &input.Options{
		Default:   "LATEST",
		Required:  false,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return err
	}
	projBuilder.SetLangVersion(langVersion)

	// Get Targets
	targets, err := targetPrompt(user, *workspace)
	if err != nil {
		return err
	}
	projBuilder.SetTargets(targets)

	fmt.Println("Note: the following items are lists, separate each item with a space")
	// get owners
	owners, err := user.Ask("Who are the owner(s) of this project? (Eg. 'elon_musk kobe')", &input.Options{
		Default:   "",
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return err
	}
	ownersList := strings.Split(owners, " ")
	projBuilder.SetOwners(ownersList)

	var dependsOnList []string
	// get the dependencies
	_, err = user.Ask("Which project(s) does this project depend on? (Eg. 'proj_a proj_d')", &input.Options{
		Default:   "",
		Required:  false,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(dep string) error {
			if dep == "" {
				dependsOnList = []string{}
				return nil
			}
			dependsOnList = strings.Split(dep, " ")
			for _, d := range dependsOnList {
				if d == name {
					dependsOnList = []string{}
					return fmt.Errorf("Project cannot depend on itself")
				}
				if _, ok := workspace.Projects[d]; !ok {
					dependsOnList = []string{}
					return fmt.Errorf("Project with name '%s' does not exist", d)
				}
			}
			return nil
		},
	})
	if err != nil {
		return err
	}
	projBuilder.SetDependsOn(dependsOnList)

	// get affects tags
	affects, err := user.Ask("Which tag(s) does this project affect? (Eg. 'client server db')"+skip, &input.Options{
		Default:      "",
		Required:     false,
		Loop:         true,
		HideOrder:    true,
		ValidateFunc: tagValidator,
	})
	if err != nil {
		return err
	}
	affectsList := strings.Split(affects, " ")
	projBuilder.SetAffectsTags(affectsList)

	// get affected by tags
	affectedBy, err := user.Ask("Which tag(s) is this project affected by? (Eg, 'node db')"+skip, &input.Options{
		Default:      "",
		Required:     false,
		Loop:         true,
		HideOrder:    true,
		ValidateFunc: tagValidator,
	})
	if err != nil {
		return err
	}
	affectedByList := strings.Split(affectedBy, " ")
	projBuilder.SetAffectedByTags(affectedByList)

	// get custom metadata
	metadataMap := make(map[string]string)
	_, err = user.Ask("What custom metadata do you want to add? (Eg.'<key1>:<val1> <key2>:<val2> ...')"+skip, &input.Options{
		Default:   "",
		Required:  false,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(metadata string) error {
			metadataList := strings.Split(metadata, " ")
			if len(metadataList) == 1 && metadataList[0] == "" {
				return nil
			}
			for _, item := range metadataList {
				split := strings.Split(item, ":")
				if len(split) != 2 {
					return fmt.Errorf("Invalid metadata format, should be '<key>:<val>'")
				}
				// generating metadata to reduce computation, if this passes validation then we won't need
				// to run the function again to generate the metadata map
				metadataMap[split[0]] = split[1]
			}
			return nil
		},
	})
	if err != nil {
		return err
	}
	projBuilder.SetMetadata(metadataMap)

	proj := projBuilder.Build()

	// add project to the workspace
	workspace.Projects[name] = p.Join(path, name)

	// write the workspace
	err = parser.EncodeWorkspace(*workspace, workspacePath, parser.JSON)
	if err != nil {
		return err
	}
	fmt.Println("Workspace updated successfully")

	path, err = misc.GetRelativePath(workspace.Projects[name])
	if err != nil {
		return err
	}

	err = parser.EncodeProjectDef(proj, path+"/project.json", parser.JSON)
	if err != nil {
		// if illegal path error
		pathErr := new(customErrors.IllegalPathError)
		if errors.As(err, &pathErr) {
			//create the path and try again
			fmt.Println("Path `" + workspace.Projects[name] + "` doesn't exist, creating it...")
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("Failed to create path: %s", err)
			}
			err = parser.EncodeProjectDef(proj, path+"/project.json", parser.JSON)
			if err != nil {
				return fmt.Errorf("Failed to parse project.json: %s", err)
			}
		}
	}
	fmt.Println("Project created successfully")

	return nil
}
