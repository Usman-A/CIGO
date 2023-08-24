package commands

import (
	"4zp6/cigo/pkg/misc"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomString() (str string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		n := rand.Intn(26)
		str += string(rune('a' + n))
	}
	return str
}

func setup(path string) error {
	// write to file
	if err := os.WriteFile(path, []byte(randomString()), 0666); err != nil {
		return fmt.Errorf("Failed to write the file: %w", err)
	}
	root, err := misc.GetRoot()
	if err != nil {
		return err
	}

	// stage changes
	cmd := exec.Command("git", "add", path)
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to add the file: %w", err)
	}

	// Set commit author
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to set user email: %w", err)
	}
	cmd = exec.Command("git", "config", "user.name", "test")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to set user name: %w", err)
	}

	// commit
	cmd = exec.Command("git", "commit", "-m", "testing affected, remove this commit if you see it")
	cmd.Dir = root
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("Failed to commit the file: %w", err)
	}

	return nil
}

func tearDown() error {
	// reset git
	root, err := misc.GetRoot()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "reset", "--hard", "HEAD~1")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to undo testing commit. Require manual intervention... %w", err)
	}
	return nil
}

func TestGetAffectedSimple(t *testing.T) {
	// setup
	path, err := misc.GetRelativePath("apps/proj_a/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	if err := setup(path); err != nil {
		t.Fatal(err)
	}

	// teardown
	defer func() {
		if err := tearDown(); err != nil {
			t.Fatal(err)
		}
	}()

	projs, err := GetAffected("HEAD~1", "HEAD")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(projs))
	assert.Equal(t, "proj_a", projs[0])

}

func TestGetAffectedDependency(t *testing.T) {
	// setup
	path, err := misc.GetRelativePath("apps/proj_b/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	if err := setup(path); err != nil {
		t.Fatal(err)
	}

	// teardown
	defer func() {
		if err := tearDown(); err != nil {
			t.Fatal(err)
		}
	}()

	// testing
	projs, err := GetAffected("HEAD~1", "HEAD")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(projs))
	assert.Equal(t, "proj_b", projs[0])
	assert.Contains(t, projs, "proj_a")
	assert.Contains(t, projs, "proj_c")
}
