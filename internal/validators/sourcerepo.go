package validators

import (
	"github.com/anchore/packageurl-go"
	"github.com/in-toto/in-toto-golang/in_toto"
	"log"
	"net/url"
	"path/filepath"
	"strings"
)

type SourceRepoValidator struct {
	Expected string
}

func (v *SourceRepoValidator) Check1(statement *in_toto.ProvenanceStatementSLSA1) bool {
	purl := v.packageUrl(v.Expected)
	if purl == "" {
		return false
	}
	// todo fix panics here
	extParams := (statement.Predicate.BuildDefinition.ExternalParameters).(map[string]any)
	val, _, _ := strings.Cut(extParams["source"].(string), "@")
	log.Printf("%s == %s", purl, val)
	return val == purl
}

func (v *SourceRepoValidator) Check02(statement *in_toto.ProvenanceStatementSLSA02) bool {
	purl := v.packageUrl(v.Expected)
	if purl == "" {
		return false
	}
	extParams := (statement.Predicate.Invocation.Parameters).(map[string]any)
	val, _, _ := strings.Cut(extParams["source"].(string), "@")
	log.Printf("%s == %s", purl, val)
	return val == purl
}

func (*SourceRepoValidator) packageUrl(s string) string {
	if !strings.HasPrefix(s, "pkg:") {
		domain, version, _ := strings.Cut(s, "@")
		uri, err := url.Parse(domain)
		if err != nil {
			return ""
		}
		// extract the repo name without the .git extension
		name := filepath.Base(strings.TrimSuffix(uri.Path, ".git"))
		// extract the namespace
		namespace := url.PathEscape(strings.TrimPrefix(filepath.Dir(uri.Path), "/"))
		// generate the package-url
		p := packageurl.NewPackageURL(strings.TrimSuffix(uri.Host, filepath.Ext(uri.Host)), namespace, name, version, nil, "")
		return p.ToString()
	}
	return s
}
