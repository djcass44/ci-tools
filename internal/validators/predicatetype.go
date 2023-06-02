package validators

import (
	"github.com/in-toto/in-toto-golang/in_toto"
	v1 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v1"
)

func PredicateTypeValidator(statement *in_toto.ProvenanceStatementSLSA1) bool {
	return statement.PredicateType == v1.PredicateSLSAProvenance
}
