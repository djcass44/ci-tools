package slsa

import (
	"encoding/json"
	"github.com/djcass44/ci-tools/pkg/in_toto/vsa"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"os"
	"path/filepath"
	"time"
)

func VSA(provenance *in_toto.ProvenanceStatementSLSA1) error {
	statement := in_toto.Statement{
		StatementHeader: in_toto.StatementHeader{
			Type:          in_toto.StatementInTotoV01,
			PredicateType: vsa.PredicateVSA,
			// todo limit to single subject
			Subject: provenance.Subject,
		},
		Predicate: vsa.Predicate{
			Verifier: vsa.Verifier{
				ID: DefaultVerifierType,
			},
			TimeVerified: time.Now(),
			// subject[0].uri
			ResourceURI: "",
			// need to make something up here
			Policy: common.ProvenanceMaterial{
				URI:    "",
				Digest: common.DigestSet{},
			},
			// digest of provenance data
			InputAttestations: nil,
			// PASSED | FAILED
			VerificationResult: "PASSED",
			VerifiedLabels: []string{
				"SLSA_LEVEL_3",
			},
			SlsaVersion: vsa.SlsaVersion1,
		},
	}

	data, err := json.MarshalIndent(&statement, "", "\t")
	if err != nil {
		return err
	}

	// write the file
	if err := os.WriteFile(filepath.Join(os.TempDir(), outVSA), data, 0644); err != nil {
		return err
	}
	return nil
}
