package slsa

import (
	"encoding/json"
	"fmt"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	v02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getDigest(target string) (string, error) {
	sha, err := crane.Digest(target, crane.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(sha, "sha256:"), nil
}

func Execute(ctx *civ1.BuildContext) error {
	repoURL := fmt.Sprintf("%s@%s", ctx.Repo.URL, ctx.Repo.Ref)
	repoDigest := common.DigestSet{digestSha1: ctx.Repo.CommitSha}

	imageDigest, err := getDigest(ctx.Image.Name)
	if err != nil {
		return err
	}
	baseDigest, err := getDigest(ctx.Image.Base)
	if err != nil {
		return err
	}
	parentDigest, err := getDigest(ctx.Image.Parent)
	if err != nil {
		return err
	}

	materials := []common.ProvenanceMaterial{
		{
			URI:    repoURL,
			Digest: repoDigest,
		},
		{
			URI:    ctx.Image.Base,
			Digest: common.DigestSet{digestSha256: baseDigest},
		},
		{
			URI:    ctx.Image.Parent,
			Digest: common.DigestSet{digestSha256: parentDigest},
		},
	}

	subjects := []in_toto.Subject{
		{
			Name:   ctx.Image.Name,
			Digest: common.DigestSet{digestSha256: imageDigest},
		},
	}

	// collect environment variables
	env := map[string]string{}
	for _, i := range os.Environ() {
		k, _, _ := strings.Cut(i, "=")
		env[k] = ""
	}

	// parse times
	buildStart, err := time.Parse(time.RFC3339, ctx.StartTime)
	if err != nil {
		return err
	}
	buildEnd := time.Now()

	provenance := in_toto.ProvenanceStatementSLSA02{
		StatementHeader: in_toto.StatementHeader{
			Type:          in_toto.StatementInTotoV01,
			PredicateType: v02.PredicateSLSAProvenance,
			Subject:       subjects,
		},
		Predicate: v02.ProvenancePredicate{
			BuildType: "https://github.com/djcass44/ci-tools@v1",
			Builder: common.ProvenanceBuilder{
				ID: ctx.Builder,
			},
			Invocation: v02.ProvenanceInvocation{
				ConfigSource: v02.ConfigSource{
					URI:        repoURL,
					Digest:     repoDigest,
					EntryPoint: "ci",
				},
				Environment: map[string]string{},
				Parameters:  env,
			},
			BuildConfig: map[string]any{
				"commands": os.Args,
				"shell":    os.Getenv("SHELL"),
			},
			Metadata: &v02.ProvenanceMetadata{
				BuildInvocationID: ctx.BuildID,
				BuildStartedOn:    &buildStart,
				BuildFinishedOn:   &buildEnd,
				Completeness: v02.ProvenanceComplete{
					Parameters:  true,
					Environment: true,
					Materials:   true,
				},
				Reproducible: true,
			},
			Materials: materials,
		},
	}

	data, err := json.MarshalIndent(&provenance, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(ctx.Root, "provenance.slsa.json"), data, 0644); err != nil {
		return err
	}

	return nil
}
