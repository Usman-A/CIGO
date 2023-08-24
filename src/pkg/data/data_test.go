package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetBuilder(t *testing.T) {
	testDependsTarget := DependsTarget{
		Project: "project",
		Target:  "target",
	}

	testTarget := Target{
		DependsOn: []DependsTarget{testDependsTarget},
		Artifacts: []string{"artifact_1", "artifact_2"},
		Cmds:      []string{"rm -r ~", "ls -a"},
		Env:       map[string]string{"API": "dr_franek_is_the_best"},
	}

	builder := TargetBuilder{}
	builder.SetArtifacts([]string{"artifact_1", "artifact_2"})
	builder.SetCmds([]string{"rm -r ~", "ls -a"})
	builder.SetDependsOn([]DependsTarget{testDependsTarget})
	builder.SetEnv(map[string]string{"API": "dr_franek_is_the_best"})
	builtTarget := builder.Build()

	assert.Equal(t, testTarget, builtTarget)
}

func TestBuilderValue(t *testing.T) {
	builder := ProjectDefinitionBuilder{}
	builder.SetName("old name")
	created := builder.Build()
	builder.SetName("new name")

	assert.NotEqual(t, created.Name, "new name")
	assert.NotEqual(t, created.Name, builder.projectDefinition.Name)
	assert.NotEqual(t, &created, &builder.projectDefinition)
	assert.Equal(t, "old name", created.Name)
}

func TestAddDependsOn(t *testing.T) {
	testTarget := Target{DependsOn: []DependsTarget{
		{
			Project: "proj_a",
			Target:  "build",
		},
		{
			Project: "proj_b",
			Target:  "build",
		},
	}}
	target := Target{}

	target.AddDependsOn("proj_a", "build")
	target.AddDependsOn("proj_b", "build")

	assert.Equal(t, testTarget, target)
}

func TestAddCmd(t *testing.T) {
	testTarget := Target{Cmds: []string{"rm -r ~", "ls -a"}}
	target := Target{}

	target.AddCmd("rm -r ~")
	target.AddCmd("ls -a")

	assert.Equal(t, testTarget, target)
}

func TestAddArtifact(t *testing.T) {
	testTarget := Target{Artifacts: []string{"artifact_1", "artifact_2"}}
	target := Target{}

	target.AddArtifact("artifact_1")
	target.AddArtifact("artifact_2")

	assert.Equal(t, testTarget, target)
}

func TestAddEnv(t *testing.T) {
	testTarget := Target{Env: map[string]string{"API": "dr_franek_is_the_best"}}
	target := Target{Env: make(map[string]string)}

	target.AddEnv("API", "dr_franek_is_the_best")

	assert.Equal(t, testTarget, target)
}

func TestWorkspaceBuilder(t *testing.T) {
	testWorkspace := Workspace{
		Owners: []string{"owner_a", "owner_b"},
		AppVer: "v1.2",
		Projects: map[string]string{
			"proj_a": "path/to/proj_a",
			"proj_b": "path/to/proj_b",
		},
		RemoteUrl:       "github.com/fake-account/fun-repo.git",
		Tags:            []string{"tag"},
		RequiredTargets: []string{"target_alpha", "target_beta"},
	}

	builder := WorkspaceBuilder{}
	builder.SetOwners([]string{"owner_a", "owner_b"})
	builder.SetAppVer("v1.2")
	builder.SetProjects(map[string]string{"proj_a": "path/to/proj_a", "proj_b": "path/to/proj_b"})
	builder.SetRemoteUrl("github.com/fake-account/fun-repo.git")
	builder.SetTags([]string{"tag"})
	builder.SetRequiredTargets([]string{"target_alpha", "target_beta"})
	builtWorkspace := builder.Build()

	assert.Equal(t, testWorkspace, builtWorkspace)
}

func TestAddOwner(t *testing.T) {
	testWorkspace := Workspace{Owners: []string{"owner_a", "owner_b"}}
	workspace := Workspace{}

	workspace.AddOwner("owner_a")
	workspace.AddOwner("owner_b")

	assert.Equal(t, testWorkspace, workspace)
}

func TestAddProject(t *testing.T) {
	testWorkspace := Workspace{Projects: map[string]string{"proj_a": "path/to/proj_a", "proj_b": "path/to/proj_b"}}
	workspace := Workspace{}

	workspace.AddProject("proj_a", "path/to/proj_a")
	workspace.AddProject("proj_b", "path/to/proj_b")

	assert.Equal(t, testWorkspace, workspace)

}

func TestAddTag(t *testing.T) {
	testWorkspace := Workspace{Tags: []string{"tag"}}
	workspace := Workspace{}

	workspace.AddTag("tag")

	assert.Equal(t, testWorkspace, workspace)
}

func TestAddRequiredTarget(t *testing.T) {
	testWorkspace := Workspace{RequiredTargets: []string{"target_alpha", "target_beta"}}
	workspace := Workspace{}

	workspace.AddRequiredTarget("target_alpha")
	workspace.AddRequiredTarget("target_beta")

	assert.Equal(t, testWorkspace, workspace)
}

func TestProjectDefinitionBuilder(t *testing.T) {
	testProjDef := ProjectDefinition{MainLanguage: "java", LangVersion: "8", Name: "monorepo", Version: "1.0.0", Targets: map[string]Target{}, Owners: []string{"ownerA", "ownerB"}, DependsOn: []string{"PROJ_A", "PROJ_B"}, AffectsTags: []string{"client"}, AffectedByTags: []string{}, Metadata: map[string]string{"color": "green"}}

	builder := ProjectDefinitionBuilder{}
	builder.SetMainLanguage("java")
	builder.SetLangVersion("8")
	builder.SetName("monorepo")
	builder.SetVersion("1.0.0")
	builder.SetTargets(map[string]Target{})
	builder.SetOwners([]string{"ownerA", "ownerB"})
	builder.SetDependsOn([]string{"PROJ_A", "PROJ_B"})
	builder.SetAffectsTags([]string{"client"})
	builder.SetAffectedByTags([]string{})
	builder.SetMetadata(map[string]string{"color": "green"})
	builtProjectDef := builder.Build()

	assert.Equal(t, testProjDef, builtProjectDef)
}

func TestAddTarget(t *testing.T) {
	testProjDef := ProjectDefinition{Targets: map[string]Target{"target_a": {}}}
	projDef := ProjectDefinition{Targets: make(map[string]Target)}

	projDef.AddTarget("target_a", Target{})

	assert.Equal(t, testProjDef, projDef)
}

func TestProjAddOwner(t *testing.T) {
	testProjDef := ProjectDefinition{Owners: []string{"ownerA", "ownerB"}}
	projDef := ProjectDefinition{}

	projDef.AddOwner("ownerA")
	projDef.AddOwner("ownerB")

	assert.Equal(t, testProjDef, projDef)
}

func TestProjAddDependsOn(t *testing.T) {
	testProjDef := ProjectDefinition{DependsOn: []string{"PROJ_A", "PROJ_B"}}
	projDef := ProjectDefinition{}

	projDef.AddDependsOn("PROJ_A")
	projDef.AddDependsOn("PROJ_B")

	assert.Equal(t, testProjDef, projDef)
}

func TestAddMetadata(t *testing.T) {
	testProjDef := ProjectDefinition{Metadata: map[string]string{"color": "green"}}
	projDef := ProjectDefinition{Metadata: make(map[string]string)}

	projDef.AddMetadata("color", "green")

	assert.Equal(t, testProjDef, projDef)
}

func TestProjAddAffectsTag(t *testing.T) {
	testProjDef := ProjectDefinition{AffectsTags: []string{"client"}}
	projDef := ProjectDefinition{}

	projDef.AddAffectsTag("client")

	assert.Equal(t, testProjDef, projDef)
}

func TestProjAddAffectedByTag(t *testing.T) {
	testProjDef := ProjectDefinition{AffectedByTags: []string{"client"}}
	projDef := ProjectDefinition{}

	projDef.AddAffectedByTag("client")

	assert.Equal(t, testProjDef, projDef)
}
