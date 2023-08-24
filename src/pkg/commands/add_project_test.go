package commands

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	input "github.com/tcnksm/go-input"
)

func init() {
	workspacePath, err := misc.GetWorkspacePath()
	if err != nil {
		fmt.Println(err)
	}
	workspace, err = parser.DecodeWorkspace(workspacePath, parser.JSON)
	if err != nil {
		fmt.Println("Error reading workspace", err)
	}
}

func TestSuccessfulTargetCreation(t *testing.T) {
	buf := bytes.NewBufferString("1\ntest_target\necho \"Hello World\"\nexe\npass:1234\n1\nproj_a\nbuild\n1\nproj_b\ndeploy\n2\n2")
	ui := &input.UI{
		Reader: buf,
	}
	targets, err := targetPrompt(ui, *workspace)
	if err != nil {
		t.Errorf("Faced error while running TargetPrompt: %v", err)
	}

	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, targets, "Should return a list of targets")
	assert.Equal(t, 1, len(targets), "Should return a list of 1 target")
	// ensure that the target_name key is present
	_, ok := targets["test_target"]
	assert.True(t, ok, "Should contain a key named `test_target`")
	assert.Equal(t, data.Target{
		DependsOn: []data.DependsTarget{{Project: "proj_a", Target: "build"}, {Project: "proj_b", Target: "deploy"}},
		Cmds:      []string{"echo \"Hello World\""},
		Env:       map[string]string{"pass": "1234"},
		Artifacts: []string{"exe"},
	}, targets["test_target"], "Target should contain these values")
	buf.Reset()
}

func TestTargetCreationWithInitialWrongInput(t *testing.T) {
	buf := bytes.NewBufferString("1\ntest_target\necho \"Hello World\"\nexe\npass:1234\n1\nproj_g\nproj_a\nsmelt\nbuild\n1\nproj_b\ndeploy\n2\n2")
	ui := &input.UI{
		Reader: buf,
	}
	targets, err := targetPrompt(ui, *workspace)
	if err != nil {
		t.Errorf("Faced error while running TargetPrompt: %v", err)
	}

	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, targets, "Should return a list of targets")
	assert.Equal(t, 1, len(targets), "Should return a list of 1 target")
	// ensure that the target_name key is present
	_, ok := targets["test_target"]
	assert.True(t, ok, "Should contain a key named `test_target`")
	assert.Equal(t, data.Target{
		DependsOn: []data.DependsTarget{{Project: "proj_a", Target: "build"}, {Project: "proj_b", Target: "deploy"}},
		Cmds:      []string{"echo \"Hello World\""},
		Env:       map[string]string{"pass": "1234"},
		Artifacts: []string{"exe"},
	}, targets["test_target"], "Target should contain these values")
	buf.Reset()
}

func TestTargetCreationWithNoDependancy(t *testing.T) {
	buf := bytes.NewBufferString("1\ntest_target\necho \"Hello World\"\nexe\n\n2\n2")
	ui := &input.UI{
		Reader: buf,
	}
	targets, err := targetPrompt(ui, *workspace)
	if err != nil {
		t.Errorf("Faced error while running TargetPrompt: %v", err)
	}

	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, targets, "Should return a list of targets")
	assert.Equal(t, 1, len(targets), "Should return a list of 1 target")
	// ensure that the target_name key is present
	_, ok := targets["test_target"]
	assert.True(t, ok, "Should contain a key named `test_target`")
	assert.Equal(t, data.Target{
		DependsOn: []data.DependsTarget(nil),
		Cmds:      []string{"echo \"Hello World\""},
		Env:       map[string]string{},
		Artifacts: []string{"exe"},
	}, targets["test_target"], "Target should contain these values")
	buf.Reset()
}

func TestTargetCreationWithWrongEnvDataFormat(t *testing.T) {
	buf := bytes.NewBufferString("1\ntest_target\necho \"Hello World\"\nexe\npass=121\npass:1234\n2\n2")
	ui := &input.UI{
		Reader: buf,
	}
	targets, err := targetPrompt(ui, *workspace)
	if err != nil {
		t.Errorf("Faced error while running TargetPrompt: %v", err)
	}

	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, targets, "Should return a list of targets")
	assert.Equal(t, 1, len(targets), "Should return a list of 1 target")
	// ensure that the target_name key is present
	_, ok := targets["test_target"]
	assert.True(t, ok, "Should contain a key named `test_target`")
	assert.Equal(t, data.Target{
		DependsOn: []data.DependsTarget(nil),
		Cmds:      []string{"echo \"Hello World\""},
		Env:       map[string]string{"pass": "1234"},
		Artifacts: []string{"exe"},
	}, targets["test_target"], "Target should contain these values")
	buf.Reset()
}

func TestNoTargetOptionWorks(t *testing.T) {
	buf := bytes.NewBufferString("2\n")
	ui := &input.UI{
		Reader: buf,
	}
	targets, err := targetPrompt(ui, *workspace)
	if err != nil {
		t.Errorf("Faced error while running TargetPrompt: %v", err)
	}

	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, targets, "Should return a list of targets")
	assert.Equal(t, 0, len(targets), "Should return a list of 0 target")
	buf.Reset()
}

func TestTagValidatorSuccessful(t *testing.T) {
	valid := tagValidator("client")
	assert.Nil(t, valid, "Should return nil for valid tag")
}

func TestInvalidTagValidator(t *testing.T) {
	valid := tagValidator("invalid_tag")
	assert.NotNil(t, valid, "Should return an error for invalid tag")
}

func TestTagValidatorWithEmptyString(t *testing.T) {
	valid := tagValidator("")
	assert.Nil(t, valid, "Should return an error for empty string")
}

func TestAddProject(t *testing.T) {
	// Spent a couple hours trying to figure out why i was having issues with the buffer,
	// it would sometimes skip the firs line and sometimes not, I could reproduce the error
	// tried different types of buffers using different packages, and nothing worked
	// so this is a hacky solution to get the mock input to work
	test_input := "need this for weird bug\ntestproj\napps/\n1.0.0\nbhailang\nLATEST\n1\ntest_target\necho \"Hello World\"\nexe\npass:1234\n1\nproj_a\nbuild\n1\nproj_b\ndeploy\n2\n2\nusman_a asadu\nproj_a\ndb\n\nsuper_secret_info:123\n"
	buf := bytes.NewBufferString(test_input)
	ui := &input.UI{
		Reader: buf,
	}
	_, _ = ui.Ask("weird_bug_for_test_case", &input.Options{
		Default:  "need this for weird bug",
		Required: true,
		Loop:     false,
	})

	err := AddProject(ui)
	fmt.Println("Successfully completed wizard")
	if err != nil {
		t.Errorf("Faced error while running AddProject: %v", err)
	}

	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, workspace.Projects, "Should return a list of projects")
	assert.Contains(t, workspace.Projects, "testproj", "Should contain a project named `testproj`")
	//file exists check
	root, err := misc.GetRoot()
	if err != nil {
		t.Errorf("Faced error while running getting root: %v", err)
	}
	_, err = os.Stat(root + "/" + workspace.Projects["testproj"] + "/project.json")
	assert.Nil(t, err, "The project definition file should exist")
	//load project.json
	proj, err := parser.DecodeProjectDef(root+"/"+workspace.Projects["testproj"]+"/project.json", parser.JSON)
	if err != nil {
		t.Errorf("Faced error while opening created project: %v", err)
	}

	// cleanup
	workspacePath, err := misc.GetWorkspacePath()
	if err != nil {
		t.Errorf("Error getting workspace path: %v", err)
	}
	os.RemoveAll(root + "/" + workspace.Projects["testproj"])
	delete(workspace.Projects, "testproj")
	err = parser.EncodeWorkspace(*workspace, workspacePath, parser.JSON)
	if err != nil {
		t.Errorf("Error ocurred while saving workspace: %v", err)
	}

	assert.Equal(t, proj, &data.ProjectDefinition{
		Name:         "testproj",
		Version:      "1.0.0",
		MainLanguage: "bhailang",
		LangVersion:  "LATEST",
		Targets: map[string]data.Target{
			"test_target": {
				DependsOn: []data.DependsTarget{{Project: "proj_a", Target: "build"}, {Project: "proj_b", Target: "deploy"}},
				Cmds:      []string{"echo \"Hello World\""},
				Env:       map[string]string{"pass": "1234"},
				Artifacts: []string{"exe"},
			}},
		Owners:         []string{"usman_a", "asadu"},
		DependsOn:      []string{"proj_a"},
		AffectsTags:    []string{"db"},
		AffectedByTags: []string{""},
		Metadata: map[string]string{
			"super_secret_info": "123",
		},
	},
	)

}

func TestAddProjectWithInitialWrongOutput(t *testing.T) {
	// Spent a couple hours trying to figure out why i was having issues with the buffer,
	// it would sometimes skip the firs line and sometimes not, I could reproduce the error
	// tried different types of buffers using different packages, and nothing worked
	// so this is a hacky solution to get the mock input to work
	test_input := "need this for weird bug\n my_proj\nmy_proj \nmy proj\nproj_a\ntestproj\n/apps/\napps/\n1.0.0\nbhailang\nLATEST\n1\ntest_target\necho \"Hello World\"\nexe\npass:1234\n1\nproj_a\nbuild\n1\nproj_b\ndeploy\n2\n2\nusman_a asadu\ntestproj\nfake\n\ndb\n\nfake=wrong\nsuper_secret_info:123\n"
	buf := bytes.NewBufferString(test_input)
	ui := &input.UI{
		Reader: buf,
	}
	_, _ = ui.Ask("weird_bug_for_test_case", &input.Options{
		Default:  "need this for weird bug",
		Required: true,
		Loop:     false,
	})

	err := AddProject(ui)
	fmt.Println("Successfully completed wizard")
	if err != nil {
		t.Errorf("Faced error while running AddProject: %v", err)
	}

	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, workspace.Projects, "Should return a list of projects")
	assert.Contains(t, workspace.Projects, "testproj", "Should contain a project named `testproj`")
	//file exists check
	root, err := misc.GetRoot()
	if err != nil {
		t.Errorf("Faced error while running getting root: %v", err)
	}
	_, err = os.Stat(root + "/" + workspace.Projects["testproj"] + "/project.json")
	assert.Nil(t, err, "The project definition file should exist")
	//load project.json
	proj, err := parser.DecodeProjectDef(root+"/"+workspace.Projects["testproj"]+"/project.json", parser.JSON)
	if err != nil {
		t.Errorf("Faced error while opening created project: %v", err)
	}

	// cleanup
	workspacePath, err := misc.GetWorkspacePath()
	if err != nil {
		t.Errorf("Error getting workspace path: %v", err)
	}
	os.RemoveAll(root + "/" + workspace.Projects["testproj"])
	delete(workspace.Projects, "testproj")
	err = parser.EncodeWorkspace(*workspace, workspacePath, parser.JSON)
	if err != nil {
		t.Errorf("Error ocurred while saving workspace: %v", err)
	}

	assert.Equal(t, proj, &data.ProjectDefinition{
		Name:         "testproj",
		Version:      "1.0.0",
		MainLanguage: "bhailang",
		LangVersion:  "LATEST",
		Targets: map[string]data.Target{
			"test_target": {
				DependsOn: []data.DependsTarget{{Project: "proj_a", Target: "build"}, {Project: "proj_b", Target: "deploy"}},
				Cmds:      []string{"echo \"Hello World\""},
				Env:       map[string]string{"pass": "1234"},
				Artifacts: []string{"exe"},
			}},
		Owners:         []string{"usman_a", "asadu"},
		DependsOn:      []string{},
		AffectsTags:    []string{"db"},
		AffectedByTags: []string{""},
		Metadata: map[string]string{
			"super_secret_info": "123",
		},
	},
	)
}
