package commands

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"os"

	"github.com/invopop/jsonschema"
)

func CreateSchema(t string) error {
	var schema jsonschema.Schema
	switch t {
	case "project":
		schema = *jsonschema.Reflect(&data.ProjectDefinition{})
	case "workspace":
		schema = *jsonschema.Reflect(&data.Workspace{})
	}
	path, err := misc.GetRelativePath("schema/" + t + ".json")
	if err != nil {
		return err
	}

	data, err := schema.MarshalJSON()
	if err != nil {
		return err
	}

	dir, err := misc.GetRelativePath("schema/")
	if err != nil {
		return err
	}
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
