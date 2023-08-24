package commands

import (
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// load workspace to use during tests
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

func TestSearchByName(t *testing.T) {

	// struct containing different searches to run and test
	tests := []struct {
		testName   string
		searchArgs map[string]string
		expected   int
	}{
		{
			testName: "Searching by name",
			searchArgs: map[string]string{
				"name": "proj_a",
			},
			expected: 1,
		},
		{
			testName: "Searching by MainLanguage",
			searchArgs: map[string]string{
				"mainLanguage": "cpp",
			},
			expected: 4,
		},
		{
			testName: "Searching by version",
			searchArgs: map[string]string{
				"version": "1.2",
			},
			expected: 1,
		},
		{
			testName: "Searching by owners",
			searchArgs: map[string]string{
				"owners": "owner1@example.com",
			},
			expected: 1,
		},
		{
			testName: "Searching by target",
			searchArgs: map[string]string{
				"target": "build",
			},
			expected: 4,
		},
		{
			testName: "Searching by dependsOn",
			searchArgs: map[string]string{
				"dependsOn": "proj_b",
			},
			expected: 2,
		},
		{
			testName: "Searching by affectsTags",
			searchArgs: map[string]string{
				"affectsTags": "db",
			},
			expected: 1,
		},
		{
			testName: "Searching by affectedByTags",
			searchArgs: map[string]string{
				"affectedByTags": "server",
			},
			expected: 1,
		},
		{
			testName: "Searching with custom metadata",
			searchArgs: map[string]string{
				"knowledge": "expert",
			},
			expected: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			response, err := Search(test.searchArgs)
			if err != nil {
				t.Errorf("Error searching: %v", err)
			}
			assert.Equal(t, len(response), test.expected, len(response))
		})
	}

}
