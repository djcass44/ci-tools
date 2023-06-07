package vsa

import (
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"time"
)

const SlsaVersion1 = "1.0"
const PredicateVSA = "https://slsa.dev/verification_summary/v1"

type Verifier struct {
	ID string `json:"id"`
}

type Predicate struct {
	Verifier           Verifier                    `json:"verifier"`
	TimeVerified       time.Time                   `json:"timeVerified"`
	ResourceURI        string                      `json:"resourceURI"`
	Policy             common.ProvenanceMaterial   `json:"policy"`
	InputAttestations  []common.ProvenanceMaterial `json:"inputAttestations"`
	VerificationResult string                      `json:"verificationResult"`
	VerifiedLabels     []string                    `json:"verifiedLabels"`
	DependencyLevels   map[string]int              `json:"dependencyLevels"`
	SlsaVersion        string                      `json:"slsaVersion"`
}
