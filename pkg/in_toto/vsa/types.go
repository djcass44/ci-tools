package vsa

import (
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"time"
)

const SlsaVersion1 = "1.0"
const SlsaVersion02 = "0.2"
const PredicateVSA = "https://slsa.dev/verification_summary/v1"

type SlsaResult string

const (
	ResultSuccess = "PASSED"

	BuildFailed SlsaResult = "FAILED"
	BuildLevel1 SlsaResult = "SLSA_BUILD_LEVEL_1"
	BuildLevel2 SlsaResult = "SLSA_BUILD_LEVEL_2"
	BuildLevel3 SlsaResult = "SLSA_BUILD_LEVEL_3"
)

type Verifier struct {
	ID string `json:"id"`
}

type Predicate struct {
	Verifier           Verifier                    `json:"verifier"`
	TimeVerified       time.Time                   `json:"timeVerified"`
	ResourceURI        string                      `json:"resourceURI"`
	Policy             common.ProvenanceMaterial   `json:"policy"`
	InputAttestations  []common.ProvenanceMaterial `json:"inputAttestations,omitempty"`
	VerificationResult string                      `json:"verificationResult"`
	VerifiedLabels     []SlsaResult                `json:"verifiedLabels"`
	DependencyLevels   map[string]int              `json:"dependencyLevels,omitempty"`
	SlsaVersion        string                      `json:"slsaVersion,omitempty"`
}
