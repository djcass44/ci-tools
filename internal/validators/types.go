package validators

import "github.com/in-toto/in-toto-golang/in_toto"

type ValidateFunc = func(statement *in_toto.ProvenanceStatementSLSA1) bool
