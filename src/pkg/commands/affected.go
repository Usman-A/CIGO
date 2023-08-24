package commands

import (
	"4zp6/cigo/pkg/algorithms"
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"os/exec"
	"strings"
)

func GetAffected(branch string, head string) (projects []string, err error) {

	cmd := exec.Command("git", "--no-pager", "diff", "--name-only", branch, head)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	wsPath, err := misc.GetWorkspacePath()
	if err != nil {
		return nil, err
	}
	ws, err := parser.DecodeWorkspace(wsPath, parser.JSON)
	if err != nil {
		return nil, err
	}
	paths := strings.Split(string(out), "\n")

	// get projects affected explicitly
	for k, p := range ws.Projects {
		for _, f := range paths {
			if strings.HasPrefix(f, p) {
				projects = append(projects, k)
				break
			}
		}
	}

	// get projects affected via dependencies
	graph := algorithms.GetGraph()
	toAdd := make([]string, 0)
	for _, p := range projects {
		deps, err := graph.GetProjectDependencies(p)
		if err != nil {
			return nil, err
		}
		toAdd = append(toAdd, deps...)
	}
	projects = append(projects, toAdd...)

	// remove duplicates
	projects = removeDuplicates(projects)

	return projects, nil
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
