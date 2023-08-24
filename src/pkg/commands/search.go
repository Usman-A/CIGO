package commands

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"fmt"
	"sync"
)

type search struct {
	Name           *string
	MainLanguage   *string
	Target         *string
	Version        *string
	Owners         *string
	DependsOn      *string
	AffectsTags    *string
	AffectedByTags *string
	Others         map[string]string
}

func mapToSearch(items map[string]string) search {
	var search search

	if name, ok := items["name"]; ok {
		search.Name = &name
		delete(items, "name")
	}
	if mainLanguage, ok := items["mainLanguage"]; ok {
		search.MainLanguage = &mainLanguage
		delete(items, "mainLanguage")
	}
	if version, ok := items["version"]; ok {
		search.Version = &version
		delete(items, "version")
	}
	if owners, ok := items["owners"]; ok {
		search.Owners = &owners
		delete(items, "owners")
	}
	if target, ok := items["target"]; ok {
		search.Target = &target
		delete(items, "target")
	}
	if dependsOn, ok := items["dependsOn"]; ok {
		search.DependsOn = &dependsOn
		delete(items, "dependsOn")
	}
	if affectsTags, ok := items["affectsTags"]; ok {
		search.AffectsTags = &affectsTags
		delete(items, "affectsTags")
	}
	if affectedByTags, ok := items["affectedByTags"]; ok {
		search.AffectedByTags = &affectedByTags
		delete(items, "affectedByTags")
	}
	search.Others = items

	return search
}

func Search(items map[string]string) ([]data.ProjectDefinition, error) {
	var projects []data.ProjectDefinition
	var wg sync.WaitGroup
	var mut sync.Mutex

	root, err := misc.GetRoot()
	if err != nil {
		return nil, err
	}
	workspace, err := parser.DecodeWorkspace(root+"/workspace.json", parser.JSON)
	if err != nil {
		return nil, err
	}
	for _, path := range workspace.Projects {
		wg.Add(1)
		p := path
		go func() {
			defer wg.Done()

			absPath, err := misc.GetRelativePath(p + "/project.json")
			if err != nil {
				fmt.Println(err)
				return
			}
			proj, err := parser.DecodeProjectDef(absPath, parser.JSON)

			if err != nil {
				fmt.Println(err)
				return
			}

			mut.Lock()
			projects = append(projects, *proj)
			mut.Unlock()
		}()
	}
	wg.Wait()

	search := mapToSearch(items)
	var res []data.ProjectDefinition
Loop:
	for _, p := range projects {
		if search.Name != nil && *search.Name != p.Name {
			continue
		}
		if search.MainLanguage != nil && *search.MainLanguage != p.MainLanguage {
			continue
		}
		if search.Version != nil && *search.Version != p.Version {
			continue
		}
		if search.Target != nil {
			if _, ok := p.Targets[*search.Target]; !ok {
				continue
			}
		}
		if search.Owners != nil && !contains(p.Owners, *search.Owners) {
			continue
		}
		if search.DependsOn != nil && !contains(p.DependsOn, *search.DependsOn) {
			continue
		}
		if search.AffectsTags != nil && !contains(p.AffectsTags, *search.AffectsTags) {
			continue
		}
		if search.AffectedByTags != nil && !contains(p.AffectedByTags, *search.AffectedByTags) {
			continue
		}

		for k, v := range search.Others {
			if v1, ok := p.Metadata[k]; !ok || v1 != v {
				// Breaks out of this loop and continues with the outer loop.
				// Very similar to a GOTO, but it continues the loop
				continue Loop
			}
		}

		res = append(res, p)
	}

	return res, nil
}

func contains(l []string, v string) bool {

	for _, i := range l {
		if i == v {
			return true
		}
	}

	return false
}
