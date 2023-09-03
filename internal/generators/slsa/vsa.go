package slsa

import (
	"encoding/json"
	"github.com/djcass44/ci-tools/pkg/in_toto/vsa"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"time"
)

func VSA[T in_toto.ProvenanceStatementSLSA1 | in_toto.ProvenanceStatementSLSA02](ok bool, provenance *T, provenanceMeta common.ProvenanceMaterial) (string, error) {
	verifiedLabel := vsa.BuildLevel3
	result := vsa.ResultSuccess
	if !ok {
		result = string(vsa.BuildFailed)
		verifiedLabel = vsa.BuildFailed
	}

	// extract information from the provenance
	// statement
	var subject []in_toto.Subject
	var slsaVersion string
	switch v := any(provenance).(type) {
	case *in_toto.ProvenanceStatementSLSA1:
		subject = v.Subject[:1]
		slsaVersion = vsa.SlsaVersion1
	case *in_toto.ProvenanceStatementSLSA02:
		subject = v.Subject[:1]
		slsaVersion = vsa.SlsaVersion02
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
			InputAttestations: []common.ProvenanceMaterial{provenanceMeta},
			// PASSED | FAILED
			VerificationResult: result,
			VerifiedLabels: []vsa.SlsaResult{
				verifiedLabel,
			},
			DependencyLevels: map[string]int{},
			SlsaVersion:      slsaVersion,
		},
	}

	data, err := json.MarshalIndent(&statement, "", "\t")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
