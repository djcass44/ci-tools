package cmd

import (
	"encoding/json"
	"github.com/djcass44/ci-tools/internal/generators/slsa"
	"github.com/djcass44/ci-tools/internal/validators"
	"github.com/djcass44/ci-tools/pkg/in_toto/vsa"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var verifyCmd = &cobra.Command{
	Use:    "verify [provenance file] [outputfile]",
	Short:  "verify SLSA provenance",
	RunE:   verifyFunc,
	Args:   cobra.ExactArgs(2),
	Hidden: true,
}

const (
	flagExpectedBuildType        = "expected-build-type"
	flagExpectedSourceRepository = "expected-source-repo"

	flagProvenancePermalink = "provenance-perma-link"
)

func init() {
	verifyCmd.Flags().String(flagExpectedBuildType, slsa.DefaultBuildType, "expected value for 'buildType'")
	verifyCmd.Flags().String(flagExpectedSourceRepository, "", "expected url (package-url or http) for the source repository")

	verifyCmd.Flags().String(flagSLSAVersion, vsa.SlsaVersion02, "slsa version (1.0 or 0.2)")
	verifyCmd.Flags().String(flagProvenancePermalink, "", "permanent reference (e.g. OCI image reference) that can be used to locate the SLSA provenance being verified")

	_ = verifyCmd.MarkFlagRequired(flagExpectedSourceRepository)
	_ = verifyCmd.MarkFlagRequired(flagProvenancePermalink)
}

func verifyFunc(cmd *cobra.Command, args []string) error {
	filename := args[0]
	output := args[1]

	// flags
	expectedBuildType, _ := cmd.Flags().GetString(flagExpectedBuildType)
	expectedSourceRepo, _ := cmd.Flags().GetString(flagExpectedSourceRepository)

	slsaVersion, _ := cmd.Flags().GetString(flagSLSAVersion)
	permaLink, _ := cmd.Flags().GetString(flagProvenancePermalink)

	// read the statement
	var statement1 *in_toto.ProvenanceStatementSLSA1
	var statement02 *in_toto.ProvenanceStatementSLSA02
	var err error

	if slsaVersion == vsa.SlsaVersion1 {
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
		if slsaVersion == vsa.SlsaVersion1 {
			ok = v.Check1(statement1)
		} else {
			ok = v.Check02(statement02)
		}
		if !ok {
			log.Printf("%s... FAILED", k)
			break
		}
		log.Printf("%s... SUCCESS", k)
	}

	// todo generate a digest of the file
	provenanceMeta := common.ProvenanceMaterial{
		URI:    permaLink,
		Digest: map[string]string{},
	}

	var data string
	if slsaVersion == vsa.SlsaVersion1 {
		data, err = slsa.VSA(ok, statement1, provenanceMeta)
	} else {
		data, err = slsa.VSA(ok, statement02, provenanceMeta)
	}
	if err != nil {
		log.Printf("failed to generate VSA")
		return err
	}

	return os.WriteFile(output, []byte(data), 0644)
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
