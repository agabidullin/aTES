package validator

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/**
var schemasFs embed.FS

func ValidateFromString(source string, topic string, key string, version int) ([]byte, error) {
	path := fmt.Sprintf("schemas/%v/%v/%v.json", topic, key, version)
	file, _ := schemasFs.ReadFile(path)
	schemaLoader := gojsonschema.NewStringLoader(string(file))
	documentLoader := gojsonschema.NewBytesLoader([]byte(source))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	if result.Valid() {
		return []byte(source), nil
	} else {
		res := "Errors: :\n"
		for _, desc := range result.Errors() {
			res += fmt.Sprintf("- %s\n", desc)
		}

		return nil, errors.New(res)
	}
}

func Validate(source interface{}, topic string, key string, version int) ([]byte, error) {
	jsonSource, err := json.Marshal(source)

	if err != nil {
		return nil, err
	}

	return ValidateFromString(string(jsonSource), topic, key, version)

}
