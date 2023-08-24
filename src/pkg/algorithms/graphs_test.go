package algorithms

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	wPath, err := misc.GetWorkspacePath()
	if err != nil {
		t.Fatal(err)
	}
	workspace, err := parser.DecodeWorkspace(wPath, parser.JSON)
	if err != nil {
		t.Fatal(err)
	}
	g := graphInternal{}
	err = (&g).init()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(workspace.Projects), len(g.projects), "Number of projects is correct")
	for name := range workspace.Projects {
		if _, ok := g.projects[name]; !ok {
			t.Errorf("Project in the workspace and not in the graph: %s\n", name)
		}
	}

}

func TestIsInit(t *testing.T) {
	g := graphInternal{}
	// check if the graph is initialized
	ok, err := g.isInit()
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, ok, "Graph is not initialized")
	// init the graph
	err = (&g).init()
	if err != nil {
		t.Fatal(err)
	}
	// check again
	ok, err = g.isInit()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, ok, "Graph is already initialized")

	//clear g projects and check again
	for k := range g.projects {
		delete(g.projects, k)
	}
	// add a new project
	g.projects["fake_1"] = data.ProjectDefinition{}
	g.projects["fake_2"] = data.ProjectDefinition{}
	g.projects["fake_3"] = data.ProjectDefinition{}
	g.projects["fake_4"] = data.ProjectDefinition{}

	ok, err = g.isInit()
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, ok, "Graph projects don't match the workspace projects")
}

func TestProjectGraph(t *testing.T) {
	g := graphInternal{}
	graph, err := g.GraphProjects()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 4, graph.Order())
	assert.Equal(t, 2, graph.Size(), "There are two dependencies in the graph.")
	p, err := graph.Vertex("proj_a")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, p, "Proj_a exists in the graph")
	e, err := graph.Edge("proj_a", "proj_b")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, e, "proj_a depends on proj_b")
}

func TestGraphTargets(t *testing.T) {
	g := graphInternal{}
	graph, err := g.GraphTargets("proj_a", "deploy")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 4, graph.Order())
	assert.Equal(t, 3, graph.Size())
	e, err := graph.Edge("proj_a/init", "proj_a/build")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, e, "Build depends on init")
}

func TestGetGraph(t *testing.T) {
	g := GetGraph()
	assert.NotNil(t, g, "GetGraph returns a non-nil value")
}

func TestGetDependecies(t *testing.T) {
	g := GetGraph()
	dependencies, err := g.GetProjectDependencies("proj_a")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, len(dependencies), "proj_a has no dependencies")
}

func TestGetDependecies2(t *testing.T) {
	g := GetGraph()
	dependencies, err := g.GetProjectDependencies("proj_b")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(dependencies), "proj_b has one dependency")
	assert.Contains(t, dependencies, "proj_a", "proj_a depends on proj_b")
	assert.Contains(t, dependencies, "proj_c", "proj_c depends on proj_b")
}

func TestGetDependeciesNoProject(t *testing.T) {
	g := GetGraph()
	_, err := g.GetProjectDependencies("doesn't exist")
	assert.NotNil(t, err, "project does not exist")
}
