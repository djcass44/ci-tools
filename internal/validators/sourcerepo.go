package validators

import "github.com/in-toto/in-toto-golang/in_toto"

type SourceRepoValidator struct {
	Expected string
}

func (v *SourceRepoValidator) Validate(statement *in_toto.ProvenanceStatementSLSA1) bool {
	extParams := (statement.Predicate.BuildDefinition.ExternalParameters).(map[string]any)
	return extParams["source"] == v.Expected
}
