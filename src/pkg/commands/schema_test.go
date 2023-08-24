package commands

import (
	"4zp6/cigo/pkg/misc"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSchema(t *testing.T) {
	schemasToTest := []string{"project", "workspace"}

	for _, schemaType := range schemasToTest {
		err := CreateSchema(schemaType)

		assert.Nil(t, err, "No errors should be returned when creating a schema")
		assert.NoError(t, err)
		schemaPath, _ := misc.GetRelativePath("schema/" + schemaType + ".json")
		_, err = os.Stat(schemaPath)
		assert.False(t, os.IsNotExist(err), "Schema file should exist after running the command")

		// Clean up created test files
		err = os.Remove(schemaPath)
		assert.Nil(t, err, "Clean up should not return an error")
	}
}
