package slsa

import (
	"encoding/json"
	"fmt"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/pkg/ociutil"
	"github.com/djcass44/ci-tools/pkg/purl"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	v02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Execute(ctx *civ1.BuildContext, r *civ1.BuildRecipe, digest string) error {
	repoURL := purl.Parse(ctx.Provider, ctx.Repo.URL, ctx.Repo.CommitSha, digestSha1, ctx.Context)
	repoDigest := common.DigestSet{digestSha1: ctx.Repo.CommitSha}

	log.Printf("generating SLSA provenance for ref: %s", repoURL)

	baseDigest := ociutil.GetDigest(ctx.Image.Base)
	parentDigest := ociutil.GetDigest(ctx.Image.Parent)

	materials := []common.ProvenanceMaterial{
		{
			URI:    repoURL,
			Digest: repoDigest,
		},
		{
			URI:    purl.Parse(purl.TypeOCI, ctx.Image.Base, baseDigest, digestSha256, ""),
			Digest: common.DigestSet{digestSha256: baseDigest},
		},
		{
			URI:    purl.Parse(purl.TypeOCI, ctx.Image.Parent, parentDigest, digestSha256, ""),
			Digest: common.DigestSet{digestSha256: parentDigest},
		},
	}

	subjects := []in_toto.Subject{
		{
			Name:   purl.Parse(purl.TypeOCI, ctx.Image.Name, digest, digestSha256, ""),
			Digest: common.DigestSet{digestSha256: digest},
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

	buildType := defaultBuildType
	if val := os.Getenv(civ1.EnvBuildSLSABuildType); val != "" {
		buildType = val
	}

	provenance := in_toto.ProvenanceStatementSLSA02{
		StatementHeader: in_toto.StatementHeader{
			Type:          in_toto.StatementInTotoV01,
			PredicateType: v02.PredicateSLSAProvenance,
			Subject:       subjects,
		},
		Predicate: v02.ProvenancePredicate{
			BuildType: buildType,
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
				"build":    append([]string{r.Command}, r.Args...),
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

	// write the provenance file
	if err := os.WriteFile(filepath.Join(ctx.Root, outProvenance), data, 0644); err != nil {
		return err
	}
	// write the digest file
	if err := os.WriteFile(filepath.Join(ctx.Root, outBuild), []byte(fmt.Sprintf("%s:%s@sha256:%s", ctx.Image.Name, ctx.Repo.CommitSha, digest)), 0644); err != nil {
		return err
	}

	return nil
}
