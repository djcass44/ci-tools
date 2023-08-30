package validators

import "github.com/in-toto/in-toto-golang/in_toto"

type ValidateFunc = func(statement *in_toto.ProvenanceStatementSLSA1) bool

type Validator interface {
	Check1(statement *in_toto.ProvenanceStatementSLSA1) bool
	Check02(statement *in_toto.ProvenanceStatementSLSA02) bool
}
