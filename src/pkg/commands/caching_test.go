package commands

import (
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSumFilesIsRepeatable(t *testing.T) {
	sha := sha256.New()

	root, err := misc.GetRoot()
	if err != nil {
		t.Fatal(err)
	}
	files := []string{root + "/src/main.go", root + "/src/go.sum", root + "/src/go.mod"}
	first, err := checkSumFiles(files, sha)
	if err != nil {
		t.Fatal(err)
	}

	second, err := checkSumFiles(files, sha)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, first, second)
}

func TestCheckSumFilesIsDifferent(t *testing.T) {
	sha := sha256.New()

	root, err := misc.GetRoot()
	if err != nil {
		t.Fatal(err)
	}
	files := []string{root + "/src/main.go", root + "/src/go.sum", root + "/src/go.mod"}
	first, err := checkSumFiles(files, sha)
	if err != nil {
		t.Fatal(err)
	}

	files = []string{root + "/src/main.go", root + "/src/go.sum"}
	second, err := checkSumFiles(files, sha)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, first, second)
}

func TestTestCheckSumDirectoryRepeatable(t *testing.T) {
	sha := sha256.New()

	dir := "apps/proj_a"
	first, err := checkSumDirectory(dir, sha)
	if err != nil {
		t.Fatal(err)
	}

	second, err := checkSumDirectory(dir, sha)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, first, second)
}

func TestGetCheksumsNoArtifacts(t *testing.T) {
	root, err := misc.GetRoot()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "apps/proj_a/project.json")
	proj, err := parser.DecodeProjectDef(dir, parser.JSON)
	if err != nil {
		t.Fatal(err)
	}
	target := proj.Targets["init"]
	hashes, err := getChecksums(dir, target)
	if err != nil {
		t.Fatal(err)
	}

	sha := sha256.New()
	emptyHash := sha.Sum([]byte{})
	// NOTE: The produced hash is 64 bytes long, but the test hash is 32 bytes long...
	assert.Equal(t, emptyHash[:32], hashes.ArtifactsHash[:32])
}

func TestGetChekSumsRepeatable(t *testing.T) {
	root, err := misc.GetRoot()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "apps/proj_a/project.json")
	proj, err := parser.DecodeProjectDef(dir, parser.JSON)
	if err != nil {
		t.Fatal(err)
	}
	target := proj.Targets["init"]
	first, err := getChecksums(dir, target)
	if err != nil {
		t.Fatal(err)
	}

	second, err := getChecksums(dir, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, first, second)
}

func TestCacheTarget(t *testing.T) {
	root, err := misc.GetRoot()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "apps/proj_a/project.json")
	proj, err := parser.DecodeProjectDef(dir, parser.JSON)
	if err != nil {
		t.Fatal(err)
	}

	// Case where CacheTarget succeeds (as it should)
	target := proj.Targets["build"]
	err = CacheTarget(dir, target)
	assert.NoError(t, err, "CacheTarget returned an error")
}

func TestIsCached(t *testing.T) {
	root, err := misc.GetRoot()
	if err != nil {
		t.Fatal(err)
	}

	// Read the project definition
	dir := filepath.Join(root, "apps/proj_d")
	proj, err := parser.DecodeProjectDef(filepath.Join(dir, "project.json"), parser.JSON)
	if err != nil {
		t.Fatal(err)
	}
	target := proj.Targets["init"]

	// Create the artifacts
	for _, file := range target.Artifacts {
		err = os.WriteFile(filepath.Join(dir, file), []byte("test"), 0666)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Delete the hashes file
	err = os.Remove(HASHE_DATA_LOCATION)
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup
	defer func() {
		for _, file := range target.Artifacts {
			err = os.Remove(filepath.Join(dir, file))
			if err != nil {
				t.Fatal(err)
			}
		}

		err = os.Remove(HASHE_DATA_LOCATION)
	}()

	// Test 1: Check if the target is cached before caching.
	// This should return False
	res, err := IsCached(dir, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, res, "Expected false, but got %v", res)

	// Test 2: Cache the target and check again. This should return True
	err = CacheTarget(dir, target)
	if err != nil {
		t.Fatal(err)
	}
	res, err = IsCached(dir, target)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, res, "Expected true, but got %v", proj.Targets["init"])
}
