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

	verifyCmd.Flags().String(flagSLSAVersion, slsaVersion02, "slsa version (1.0 or 0.2)")

	_ = verifyCmd.MarkFlagRequired(flagExpectedSourceRepository)
}

func verifyFunc(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// flags
	expectedBuildType, _ := cmd.Flags().GetString(flagExpectedBuildType)
	expectedSourceRepo, _ := cmd.Flags().GetString(flagExpectedSourceRepository)

	slsaVersion, _ := cmd.Flags().GetString(flagSLSAVersion)

	// read the statement
	var statement1 *in_toto.ProvenanceStatementSLSA1
	var statement02 *in_toto.ProvenanceStatementSLSA02
	var err error

	if slsaVersion == slsaVersion10 {
		statement1, err = loadFile[in_toto.ProvenanceStatementSLSA1](filename)
	} else {
		statement02, err = loadFile[in_toto.ProvenanceStatementSLSA02](filename)
	}
	if err != nil {
		return err
	}

	buildTypeValidator := &validators.BuildTypeValidator{Expected: expectedBuildType}
	sourceRepoValidator := &validators.SourceRepoValidator{Expected: expectedSourceRepo}
	internalParameterValidator := &validators.InternalParameterValidator{}
	predicateTypeValidator := &validators.PredicateTypeValidator{}

	assertions := map[string]validators.Validator{
		"Build type":                  buildTypeValidator,
		"Internal parameters":         internalParameterValidator,
		"Predicate type":              predicateTypeValidator,
		"Canonical source repository": sourceRepoValidator,
	}

	var ok bool
	for k, v := range assertions {
		if slsaVersion == slsaVersion10 {
			ok = v.Check1(statement1)
		} else {
			ok = v.Check02(statement02)
		}
		if !ok {
			log.Printf("%s... FAILED", k)
			return errors.New("statement validation failed")
		}
		log.Printf("%s... SUCCESS", k)
	}

	if slsaVersion == slsaVersion10 {
		err = slsa.VSA(statement1)
	} else {
		err = slsa.VSA(statement02)
	}
	if err != nil {
		log.Printf("failed to generate VSA")
		return err
	}

	return nil
}

func loadFile[T in_toto.ProvenanceStatementSLSA1 | in_toto.ProvenanceStatementSLSA02](path string) (*T, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var statement T
	if err := json.NewDecoder(f).Decode(&statement); err != nil {
		return nil, err
	}
	return &statement, nil
}
