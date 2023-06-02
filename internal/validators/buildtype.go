package validators

import "github.com/in-toto/in-toto-golang/in_toto"

type BuildTypeValidator struct {
	Expected string
}

func (v *BuildTypeValidator) Validate(statement *in_toto.ProvenanceStatementSLSA1) bool {
	return statement.Predicate.BuildDefinition.BuildType == v.Expected
}
