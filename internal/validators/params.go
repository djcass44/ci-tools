package validators

import (
	"github.com/in-toto/in-toto-golang/in_toto"
	"golang.org/x/exp/slices"
)

var allowedInternalParameters = []string{"shell", "commands", "build"}

func InternalParameterValidator(statement *in_toto.ProvenanceStatementSLSA1) bool {
	intParams := (statement.Predicate.BuildDefinition.InternalParameters).(map[string]any)
	for k := range intParams {
		if !slices.Contains(allowedInternalParameters, k) {
			return false
		}
	}
	return true
}
