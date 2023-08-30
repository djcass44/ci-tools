package validators

import (
	"github.com/in-toto/in-toto-golang/in_toto"
	v02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	v1 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v1"
)

type PredicateTypeValidator struct{}

func (*PredicateTypeValidator) Check1(statement *in_toto.ProvenanceStatementSLSA1) bool {
	return statement.PredicateType == v1.PredicateSLSAProvenance
}

func (*PredicateTypeValidator) Check02(statement *in_toto.ProvenanceStatementSLSA02) bool {
	return statement.PredicateType == v02.PredicateSLSAProvenance
}
