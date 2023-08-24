package algorithms

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"errors"
	"fmt"
	"sync"

	"github.com/dominikbraun/graph"
)

func GetGraph() IGragh {
	return &graphInternal{}
}

// need a way to store the Project name and the target name in the struct. Best
// option was creating a new struct.
type graphTarget struct {
	Project string
	Name    string
	Target  data.Target
}

type IGragh interface {
	GraphProjects() (projectGraph, error)
	GraphTargets(string, string) (targetGraph, error)
	GetProjectDependencies(string) ([]string, error)
}
type graphInternal struct {
	// Used for frequent access to the projects.
	projects map[string]data.ProjectDefinition
}

type (
	targetGraph  graph.Graph[string, graphTarget]
	projectGraph graph.Graph[string, data.ProjectDefinition]
)

// Check if the internal state has been initialized or not
func (g *graphInternal) isInit() (bool, error) {
	workspacePath, err := misc.GetWorkspacePath()
	if err != nil {
		return false, err
	}
	workspace, err := parser.DecodeWorkspace(workspacePath, parser.JSON)
	if err != nil {
		return false, err
	}

	if len(workspace.Projects) != len(g.projects) {
		return false, nil
	}

	for proj := range workspace.Projects {
		if _, ok := g.projects[proj]; !ok {
			return false, nil
		}
	}

	return true, nil
}

// Initialize the internal state
func (g *graphInternal) init() error {
	if ok, _ := g.isInit(); ok {
		return nil
	}
	if g.projects == nil {
		g.projects = make(map[string]data.ProjectDefinition)
	}
	workspacePath, err := misc.GetWorkspacePath()
	if err != nil {
		return err
	}
	workspace, err := parser.DecodeWorkspace(workspacePath, parser.JSON)
	if err != nil {
		return err
	}

	var (
		wg      sync.WaitGroup
		mutex   sync.Mutex
		projErr error
	)
	for name, path := range workspace.Projects {
		// Get the project from the file
		pPath, err := misc.GetRelativePath(path + "/project.json")
		if err != nil {
			return err
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
			g.projects[pName] = *proj
			mutex.Unlock()

			if err != nil {
				projErr = errors.New(projErr.Error() + "\n" + err.Error())
			}
		}()
	}
	wg.Wait()
	if projErr != nil {
		return projErr
	}

	return nil
}

func TargetHash(t graphTarget) string {
	return fmt.Sprintf("%s/%s", t.Project, t.Name)
}

// Abstracting this away
// It is left here due to the use of the struct
func (g graphInternal) getTargetDependencies(t graphTarget) ([]graphTarget, error) {
	var targets []graphTarget

	for _, tDeps := range t.Target.DependsOn {
		var project string
		var (
			p  data.ProjectDefinition
			ok bool
		)
		if tDeps.Project == "self" {
			project = t.Project
		} else {
			project = tDeps.Project

		}
		p, ok = g.projects[project]
		if !ok {
			return nil, fmt.Errorf("Project not found: %s", g.projects)
		}
		target, ok := p.Targets[tDeps.Target]
		if !ok {
			return nil, fmt.Errorf("Target not found: %s", p)
		}
		targets = append(targets, graphTarget{
			Project: project,
			Name:    tDeps.Target,
			Target:  target,
		})
	}

	return targets, nil
}

// Uses the internal state
func (g *graphInternal) getTarget(project string, target string) (*data.Target, error) {

	p, ok := g.projects[project]
	if !ok {
		return nil, fmt.Errorf("Project not found: %s", g.projects)
	}
	t, ok := p.Targets[target]
	if !ok {
		return nil, fmt.Errorf("Target not found: %s", p)
	}
	return &t, nil
}

func (gi *graphInternal) GraphTargets(project string, target string) (targetGraph, error) {
	err := gi.init()
	if err != nil {
		return nil, err
	}
	g := graph.New(TargetHash, graph.PreventCycles(), graph.Directed())

	// get the target
	t, err := gi.getTarget(project, target)
	if err != nil {
		return nil, err
	}
	// create the root node
	rootTarget := graphTarget{
		Project: project,
		Name:    target,
		Target:  *t,
	}

	// Stack to keep track of all the targets that needs their dependencies added
	stack := []graphTarget{rootTarget}
	// Add the root target to the Graph
	// The target must be already in the graph before its dependencies are added
	err = g.AddVertex(rootTarget)
	if err != nil {
		return nil, err
	}

	for len(stack) > 0 {
		node := stack[0]
		// 1. Get the target dependencies
		deps, err := gi.getTargetDependencies(node)
		if err != nil {
			return nil, err
		}

		// 2. Add them to the stack, for later processing
		stack = append(stack[1:], deps...)

		// 3. Add them to the graph
		for _, gt := range deps {
			// 3.1. Add the vertex
			err = g.AddVertex(gt)
			if err != nil {
				return nil, err
			}

			// 3.2. Add the edge
			// This is why we need the dependent target already in the graph
			// t1 -> t2 = t1 depends on t2
			err = g.AddEdge(TargetHash(gt), TargetHash(node))
			if err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

// functions required by graph library in order to create graphs of new types
func projectHash(p data.ProjectDefinition) string {
	return p.Name
}

// NOTE: Think about the graph size if we store the projects in the graph
func (self *graphInternal) GraphProjects() (projectGraph, error) {

	// initialize the project list
	err := self.init()

	if err != nil {
		return nil, err
	}

	//initialize graph
	g := graph.New(projectHash, graph.Directed(), graph.PreventCycles())

	// Add vertices to the graph
	for _, project := range self.projects {
		err = g.AddVertex(project)
		if err != nil {
			return nil, err
		}
	}

	//add directed edges from a project to it's dependencies
	for _, project := range self.projects {
		for _, dependency := range project.DependsOn {
			err := g.AddEdge(project.Name, dependency)
			if err != nil {
				return g, err
			}
		}
	}

	return g, nil
}

func (self *graphInternal) GetProjectDependencies(project string) (projs []string, err error) {
	// initialize the graph
	projGraph, err := self.GraphProjects()
	if err != nil {
		return nil, err
	}
	if _, ok := self.projects[project]; !ok {
		return nil, fmt.Errorf("Project not found: %s", project)
	}
	var toProcess []string = []string{project}
	deps, err := projGraph.PredecessorMap()
	for len(toProcess) > 0 {
		// pop the first element
		toCheck := toProcess[0]
		toProcess = toProcess[1:]

		// on the dependencies
		for p := range deps[toCheck] {
			toProcess = append(toProcess, p)
			projs = append(projs, p)
		}

	}

	return projs, err
}
