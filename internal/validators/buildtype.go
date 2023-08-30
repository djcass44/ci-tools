package validators

import "github.com/in-toto/in-toto-golang/in_toto"

type BuildTypeValidator struct {
	Expected string
}

func (v *BuildTypeValidator) Check1(statement *in_toto.ProvenanceStatementSLSA1) bool {
	return statement.Predicate.BuildDefinition.BuildType == v.Expected
}

func (v *BuildTypeValidator) Check02(statement *in_toto.ProvenanceStatementSLSA02) bool {
	return statement.Predicate.BuildType == v.Expected
}
