package cmd

import (
	"encoding/json"
	"errors"
	"github.com/djcass44/ci-tools/internal/generators/slsa"
	"github.com/djcass44/ci-tools/internal/validators"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var verifyCmd = &cobra.Command{
	Use:    "verify",
	Short:  "verify SLSA provenance",
	RunE:   verifyFunc,
	Args:   cobra.ExactArgs(1),
	Hidden: true,
}

const (
	flagExpectedBuildType        = "expected-build-type"
	flagExpectedSourceRepository = "expected-source-repo"
)

func init() {
	verifyCmd.Flags().String(flagExpectedBuildType, slsa.DefaultBuildType, "expected value for 'buildType'")
	verifyCmd.Flags().String(flagExpectedSourceRepository, "", "expected url (package-url or http) for the source repository")

	_ = verifyCmd.MarkFlagRequired(flagExpectedSourceRepository)
}

func verifyFunc(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// flags
	expectedBuildType, _ := cmd.Flags().GetString(flagExpectedBuildType)
	expectedSourceRepo, _ := cmd.Flags().GetString(flagExpectedSourceRepository)

	// read the statement
	statement, err := loadFile(filename)
	if err != nil {
		return err
	}

	buildTypeValidator := validators.BuildTypeValidator{Expected: expectedBuildType}
	sourceRepoValidator := validators.SourceRepoValidator{Expected: expectedSourceRepo}

	assertions := map[string]validators.ValidateFunc{
		"Build type":                  buildTypeValidator.Validate,
		"Internal parameters":         validators.InternalParameterValidator,
		"Predicate type":              validators.PredicateTypeValidator,
		"Canonical source repository": sourceRepoValidator.Validate,
	}

	for k, v := range assertions {
		ok := v(statement)
		if !ok {
			log.Printf("%s... FAILED", k)
			return errors.New("statement validation failed")
		}
		log.Printf("%s... SUCCESS", k)
	}

	if err := slsa.VSA(statement); err != nil {
		log.Printf("failed to generate VSA")
		return err
	}

	return nil
}

func loadFile(path string) (*in_toto.ProvenanceStatementSLSA1, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var statement in_toto.ProvenanceStatementSLSA1
	if err := json.NewDecoder(f).Decode(&statement); err != nil {
		return nil, err
	}
	return &statement, nil
}
