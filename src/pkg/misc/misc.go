package misc

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/parser"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// get  git root dir

func GetRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func GetWorkspacePath() (string, error) {
	root, err := GetRoot()
	if err != nil {
		return "", err
	}

	return root + "/workspace.json", nil
}

func GetRelativePath(path string) (string, error) {
	root, err := GetRoot()
	if err != nil {
		return "", err
	}

	return root + "/" + path, nil
}

func Contains[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func GetProjects() (map[string]data.ProjectDefinition, error) {
	workspacePath, err := GetWorkspacePath()
	if err != nil {
		return nil, err
	}
	workspace, err := parser.DecodeWorkspace(workspacePath, parser.JSON)
	if err != nil {
		return nil, err
	}

	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		projErr  error
		projects map[string]data.ProjectDefinition = make(map[string]data.ProjectDefinition)
	)
	for name, path := range workspace.Projects {
		// Get the project from the file
		pPath, err := GetRelativePath(path + "/project.json")
		if err != nil {
			return nil, err
		}
		pName := name
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Read project data
			proj, err := parser.DecodeProjectDef(pPath, parser.JSON)
			if err != nil {
				if projErr == nil {
					projErr = fmt.Errorf("Failed to parse projects:\n\t")
				}
				projErr = errors.New(projErr.Error() + "\n" + err.Error())
				return
			}

			// Store it in map
			mutex.Lock()
			projects[pName] = *proj
			mutex.Unlock()

			if err != nil {
				projErr = errors.New(projErr.Error() + "\n" + err.Error())
			}
		}()
	}
	wg.Wait()
	if projErr != nil {
		return nil, projErr
	}

	return projects, nil

}
