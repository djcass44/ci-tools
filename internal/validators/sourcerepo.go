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

func (v *SourceRepoValidator) Validate(statement *in_toto.ProvenanceStatementSLSA1) bool {
	purl := v.Expected
	// if the expected URL is not a package-url,
	// then we need to convert it into one
	if !strings.HasPrefix(purl, "pkg:") {
		domain, version, _ := strings.Cut(purl, "@")
		uri, err := url.Parse(domain)
		if err != nil {
			return false
		}
		// extract the repo name without the .git extension
		name := filepath.Base(strings.TrimSuffix(uri.Path, ".git"))
		// extract the namespace
		namespace := url.PathEscape(strings.TrimPrefix(filepath.Dir(uri.Path), "/"))
		// generate the package-url
		p := packageurl.NewPackageURL(strings.TrimSuffix(uri.Host, filepath.Ext(uri.Host)), namespace, name, version, nil, "")
		purl = p.ToString()
		log.Printf("purl: %s", purl)
	}
	extParams := (statement.Predicate.BuildDefinition.ExternalParameters).(map[string]any)
	return extParams["source"] == purl
}
