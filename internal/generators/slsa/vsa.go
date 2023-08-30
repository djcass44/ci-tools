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

func VSA[T in_toto.ProvenanceStatementSLSA1 | in_toto.ProvenanceStatementSLSA02](provenance *T) error {
	// extract the subject
	var subject []in_toto.Subject
	switch v := any(provenance).(type) {
	case in_toto.ProvenanceStatementSLSA1:
		subject = v.Subject[:1]
	case in_toto.ProvenanceStatementSLSA02:
		subject = v.Subject[:1]
	}
	statement := in_toto.Statement{
		StatementHeader: in_toto.StatementHeader{
			Type:          in_toto.StatementInTotoV01,
			PredicateType: vsa.PredicateVSA,
			Subject:       subject,
		},
		Predicate: vsa.Predicate{
			Verifier: vsa.Verifier{
				ID: DefaultVerifierType,
			},
			TimeVerified: time.Now(),
			// subject[0].uri
			ResourceURI: subject[0].Name,
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
