package parser

import (
	"4zp6/cigo/pkg/customErrors"
	"4zp6/cigo/pkg/data"
	"encoding/json"
	"os"
	"syscall"
)

// Create a FileType that takes in a JSON argument
type FileType uint8

const (
	JSON FileType = iota
)

// Function that, given a filePath and a fileType, decodes the path in the given file
// and returns an instance of a ProjectDefinition struct
func DecodeProjectDef(filePath string, fileType FileType) (*data.ProjectDefinition, error) {

	var out *data.ProjectDefinition

	// Opens the file in the given path, or returns an error
	openFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &customErrors.FileNotFoundError{Message: "FileNotFound: " + err.Error()}
	}

	// Decodes the file from JSON/YAML into a ProjectDefinition type
	if fileType == JSON {
		err = json.Unmarshal([]byte(openFile), &out)
		if err != nil {
			// In the case of a serialization error
			return nil, &customErrors.InvalidFormatError{Message: "InvalidFormat: " + err.Error()}
		}
	}
	return out, nil
}

// Function that, given a ProjectDefinition type, a file path, and a FileType, encodes the
// ProjectDefinition into a file of the specified type (JSON)
func EncodeProjectDef(project data.ProjectDefinition, filePath string, fileType FileType) error {
	// Encodes the file given file into a JSON/YAML format (depending on type specified) and writes it
	// on a corresponding JSON file that it creates/updates in the given path, or returns an error
	if fileType == JSON {
		// Makes sure the path specifies a JSON file
		if filePath[len(filePath)-4:] != "json" {
			// In the case of a bad path
			return &customErrors.IllegalPathError{Message: "IllegalPath, path must end in .json"}
		}
		data, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			// In the case of a deserialization error
			return &customErrors.EncodingError{Message: "EncodingError: " + err.Error()}
		}
		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			if pathError, ok := err.(*os.PathError); ok {
				if pathError.Err == syscall.ENOENT {
					// In the case of a bad path
					return &customErrors.IllegalPathError{Message: "IllegalPath, path must end in .json:" + err.Error()}
				} else {
					// If it fails to write
					return &customErrors.IOError{Message: "IOError" + err.Error()}
				}
			}
		}
	}
	return nil
}

// Function that, given a filePath and a fileType, decodes the path in the given file
// and returns an instance of a Workspace struct
func DecodeWorkspace(filePath string, fileType FileType) (*data.Workspace, error) {

	var out *data.Workspace

	// Opens the file in the given path, or returns an error
	openFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &customErrors.FileNotFoundError{Message: "FileNotFound: " + err.Error()}
	}

	// Decodes the file from JSON into a Workspace type
	if fileType == JSON {
		err = json.Unmarshal([]byte(openFile), &out)
		if err != nil {
			// In the case of a serialization error
			return nil, &customErrors.InvalidFormatError{Message: "InvalidFormat: " + err.Error()}
		}
	}
	return out, nil
}

// Function that, given a Workspace type, a file path, and a FileType, encodes the
// Workspace into a file of the specified type (JSON)
func EncodeWorkspace(workspace data.Workspace, filePath string, fileType FileType) error {
	// Encodes the file given file into a JSON format (depending on type specified) and writes it
	// on a corresponding JSON file that it creates/updates in the given path, or returns an error
	if fileType == JSON {
		// Makes sure the path specifies a JSON file
		if filePath[len(filePath)-4:] != "json" {
			// In the case of a bad path
			return &customErrors.IllegalPathError{Message: "IllegalPath, path must end in .json:"}
		}
		data, err := json.MarshalIndent(workspace, "", "  ")
		if err != nil {
			// In the case of a deserialization error
			return &customErrors.EncodingError{Message: "EncodingError: " + err.Error()}
		}
		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			if pathError, ok := err.(*os.PathError); ok {
				if pathError.Err == syscall.ENOENT {
					// In the case of a bad path
					return &customErrors.IllegalPathError{Message: "IllegalPath:" + err.Error()}
				} else {
					// If it fails to write
					return &customErrors.IOError{Message: "IOError" + err.Error()}
				}
			}
		}
	}
	return nil
}
