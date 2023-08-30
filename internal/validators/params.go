package validators

import (
	"github.com/in-toto/in-toto-golang/in_toto"
	"golang.org/x/exp/slices"
)

var allowedInternalParameters = []string{"shell", "commands", "build"}

type InternalParameterValidator struct{}

func (*InternalParameterValidator) Check1(statement *in_toto.ProvenanceStatementSLSA1) bool {
	intParams := (statement.Predicate.BuildDefinition.InternalParameters).(map[string]any)
	for k := range intParams {
		if !slices.Contains(allowedInternalParameters, k) {
			return false
		}
	}
	return true
}

func (*InternalParameterValidator) Check02(statement *in_toto.ProvenanceStatementSLSA02) bool {
	intParams := (statement.Predicate.Invocation.Environment).(map[string]any)
	for k := range intParams {
		if !slices.Contains(allowedInternalParameters, k) {
			return false
		}
	}
	return true
}
