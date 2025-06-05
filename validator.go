package main

import (
	"fmt"
	"log"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/yaml"
)

func validateWithCUE(context *cue.Context, yamlPath string, schemaPath string, fSchemaName string) error {
	yamlFile, err := yaml.Extract(yamlPath, nil)
	if err != nil {
		log.Fatal(err)
	}

	yamlAsCUE := context.BuildFile(yamlFile)

	// load schema file
	cueSchema := loadCUESchema(context, schemaPath)

	var errInvalidSchema = fmt.Errorf("invalid schema name")

	switch fSchemaName {
	case "cpinit":
		cueSchema = cueSchema.LookupPath(cue.ParsePath("#CpInitRoot"))
	case "cp":
		cueSchema = cueSchema.LookupPath(cue.ParsePath("#CpRoot"))
	case "map":
		cueSchema = cueSchema.LookupPath(cue.ParsePath("#MapRoot"))
	case "exp":
		cueSchema = cueSchema.LookupPath(cue.ParsePath("#ExpRoot"))
	case "exec":
		cueSchema = cueSchema.LookupPath(cue.ParsePath("#ExecRoot"))
	default:
		return errInvalidSchema
	}

	if err := cueSchema.Err(); err != nil {
		print(err.Error())
		panic(err)
	}

	unifiedCue := cueSchema.Unify(yamlAsCUE)
	if err != nil {
		fmt.Println(err)
	}

	err = strict(unifiedCue)
	if err != nil {
		msg := errors.Details(err, nil)
		fmt.Printf("Validate Error:\n%s\n", msg)
	}

	return err
}

func strict(v cue.Value) error {
	opt := []cue.Option{
		cue.Attributes(true),
		cue.Definitions(true),
		cue.Hidden(true),
	}
	return v.Validate(append(opt, cue.Concrete(true))...)
}

func loadCUESchema(context *cue.Context, schemaPath string) cue.Value {
	valBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Println(err)
	}

	val := context.CompileBytes(valBytes, cue.Filename(schemaPath))
	if err := val.Err(); err != nil {
		fmt.Println(err)
	}

	return val
}

func main() {
	c := cuecontext.New()
	cpyamlPath := "test.yaml"
	schemaPath := "schema.cue"
	if err := validateWithCUE(c, cpyamlPath, schemaPath, "cpinit"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Syntax and schema validated successfully")
	}
}
