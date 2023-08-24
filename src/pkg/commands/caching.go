package commands

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const HASHE_DATA_LOCATION = "/tmp/cigo/hashes.json"

type Caching struct {
	Hashes map[string]Hashes
}

type Hashes struct {
	ProjectHash   []byte `json:"project_hash"`
	ArtifactsHash []byte `json:"artifacts_hash"`
}

// Path: src/pkg/commands/caching.go
// Name: IsCached
// Description: checks if the project is cached
// Parameters:
//   - projectPath: the path to the project
//   - target: the execution target to check
//
// Returns:
//   - bool: true if the project is cached, false otherwise
//   - error: any errors that occur
func IsCached(projectPath string, target data.Target) (bool, error) {

	// read the hashes file and check if the project is cached
	var cache Caching
	cacheData, err := os.ReadFile(HASHE_DATA_LOCATION)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	err = json.Unmarshal(cacheData, &cache)
	if err != nil {
		return false, err
	}
	if _, ok := cache.Hashes[projectPath]; !ok {
		return false, nil

	}

	// get the checksums of the project and the target artifacts
	hashes, err := getChecksums(projectPath, target)
	if err != nil {
		return false, err
	}

	// compare the checksum
	if bytes.Equal(cache.Hashes[projectPath].ProjectHash, hashes.ProjectHash) &&
		bytes.Equal(cache.Hashes[projectPath].ArtifactsHash, hashes.ArtifactsHash) {
		return true, nil
	}

	return false, nil
}

// Path: src/pkg/commands/caching.go
// Name: getChecksums
// Description: gets the checksums of the project and the target artifacts
// Parameters:
//   - projectPath: the path to the project
//   - target: the target to get the checksums for
//
// Returns:
//   - *Hashes: the checksums of the project and the target artifacts
//   - error: any errors that occur
func getChecksums(projectPath string, target data.Target) (*Hashes, error) {

	var hashes Hashes
	// get files tracked by git
	// NOTE: this is not the best way to do this, but it is the easiest
	//       consider using git2go or go-git
	cmd := exec.Command("git", "ls-tree", "-r", "HEAD", "--name-only")
	root, err := misc.GetRoot()
	if err != nil {
		return nil, fmt.Errorf("Failed to get repo root: %v\n", err)
	}
	cmd.Dir = root
	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("Failed to get tracked files: %v\n", err)
	}

	// split the output into paths
	paths := strings.Split(string(out), "\n")

	projectPaths := []string{}
	for _, path := range paths {
		if strings.HasPrefix(path, projectPath) {
			projectPaths = append(projectPaths, path)
		}
	}

	sha256 := sha256.New()

	projectCheckSum, err := checkSumFiles(projectPaths, sha256)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the project files checksum: %v\n", err)
	}

	var artifactsCheckSum []byte = []byte{}
	artifactFiles := []string{}
	// sort the artifacts so that the checksum is consistent
	sort.StringSlice(target.Artifacts).Sort()
	for _, artifact := range target.Artifacts {
		// check if it is a file or a directory
		fileInfo, err := os.Stat(artifact)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to get the artifact directory stats: %v\n", err)
		}

		if fileInfo.IsDir() {
			sum, err := checkSumDirectory(artifact, sha256)
			if err != nil {
				return nil, fmt.Errorf("Failed to get the artifact directory checksum: %v\n", err)
			}
			artifactsCheckSum = append(artifactsCheckSum, sum...)
		} else {
			root, err := misc.GetRoot()
			if err != nil {
				return nil, fmt.Errorf("Failed to get repo root: %v\n", err)
			}
			artifactFiles = append(artifactFiles, filepath.Join(root, artifact))
		}
	}

	// get the checksum of the artifacts
	sum, err := checkSumFiles(artifactFiles, sha256)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the artifact files stats: %v\n", err)
	}

	artifactsCheckSum = append(artifactsCheckSum, sum...)
	artifactsCheckSum = sha256.Sum(artifactsCheckSum)

	hashes.ProjectHash = projectCheckSum
	hashes.ArtifactsHash = artifactsCheckSum

	return &hashes, nil
}

// Path: src/pkg/commands/caching.go
// Name: checkSumDirectory
// Description: gets the checksum of all files in a directory
// Parameters:
//   - dirPath: the path to the directory, relative to the repository
//
// Returns:
//   - []byte: the checksum of all files in the directory
//   - error: any errors that occur
func checkSumDirectory(dirPath string, sha256 hash.Hash) ([]byte, error) {
	relPath, err := misc.GetRelativePath(dirPath)
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(relPath)
	if err != nil {
		return nil, err
	}

	filePaths := make([]string, len(files))
	for i, file := range files {
		filePaths[i] = filepath.Join(relPath, file.Name())
	}

	return checkSumFiles(filePaths, sha256)
}

// Path: src/pkg/commands/caching.go
// Name: checkSumFiles
// Description: gets the checksum of all files in a list of paths
// Parameters:
//   - filePath: the paths to the files, absolute
//
// Returns:
//   - []byte: the checksum of all files in the list of paths
//   - error: any errors that occur
func checkSumFiles(filePaths []string, sha256 hash.Hash) ([]byte, error) {

	sort.StringSlice(filePaths).Sort()
	var runningCheckSum []byte = []byte{}
	for _, file := range filePaths {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		sum := sha256.Sum(data)
		runningCheckSum = append(runningCheckSum, sum...)
	}

	return sha256.Sum(runningCheckSum), nil
}

// Path: src/pkg/commands/caching.go
// Name: CacheTarget
// Description: caches the target artifacts
// Parameters:
//   - projectPath: the path to the project
//   - target: the target to cache
//
// Returns:
//   - error: any errors that occur
func CacheTarget(projectPath string, target data.Target) error {
	hashes, err := getChecksums(projectPath, target)
	if err != nil {
		return err
	}

	var cache Caching

	cachedData, err := os.ReadFile(HASHE_DATA_LOCATION)
	if err == nil {
		err = json.Unmarshal(cachedData, &cache)
		if err != nil {
			return err
		}

	} else if errors.Is(err, os.ErrNotExist) {
		path := filepath.Dir(HASHE_DATA_LOCATION)
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		cache = Caching{
			Hashes: map[string]Hashes{},
		}
	} else {
		return err
	}

	cache.Hashes[projectPath] = *hashes

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = os.WriteFile(HASHE_DATA_LOCATION, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
