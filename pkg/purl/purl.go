package purl

import (
	"fmt"
	"github.com/anchore/packageurl-go"
	"path/filepath"
	"strings"
)

const (
	TypeGitLab = "gitlab"
	TypeOCI    = "oci"

	QualifierTag     = "tag"
	QualifierRepoURL = "repository_url"
)

func Parse(purlType, target, digest, digestType, path string) string {
	name := filepath.Base(target)
	name, tag, _ := strings.Cut(name, ":")
	tag, _, _ = strings.Cut(tag, "@sha256:")
	namespace := filepath.Dir(strings.TrimPrefix(target, "https://"))
	repoURL, _, _ := strings.Cut(namespace, "/")
	namespace = strings.TrimPrefix(strings.TrimPrefix(namespace, repoURL), "/")

	qualifiers := map[string]string{}

	if purlType == TypeOCI {
		repoURL = filepath.Join(repoURL, namespace, name)
		namespace = ""
		name = strings.ToLower(name)
	}

	if purlType == TypeOCI || purlType == TypeGitLab {
		qualifiers[QualifierRepoURL] = repoURL
	}

	if purlType == TypeOCI && tag != "" {
		qualifiers[QualifierTag] = tag
	}

	return packageurl.NewPackageURL(purlType, namespace, name, fmt.Sprintf("%s:%s", digestType, digest), packageurl.QualifiersFromMap(qualifiers), path).String()
}
