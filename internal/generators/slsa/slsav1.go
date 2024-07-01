package slsa

import (
	"fmt"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/pkg/digestof"
	"github.com/djcass44/ci-tools/pkg/ociutil"
	"github.com/djcass44/ci-tools/pkg/purl"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	v1 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v1"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ExecuteV1(ctx *civ1.BuildContext, r *civ1.BuildRecipe, ref, digest string, predicateOnly bool) error {
	repoURL := purl.Parse(ctx.Provider, ctx.Repo.URL, ctx.Repo.CommitSha, digestSha1, ctx.Context)
	repoDigest := common.DigestSet{digestSha1: ctx.Repo.CommitSha}

	log.Printf("generating SLSA (v1.0) provenance for ref: %s", repoURL)

	auth := ctx.Auth()

	baseDigest := ociutil.GetDigest(ctx.Image.Base, auth)
	parentDigest := ociutil.GetDigest(ctx.Image.Parent, auth)

	materials := []v1.ResourceDescriptor{
		{
			URI:    repoURL,
			Digest: repoDigest,
		},
		{
			URI:    purl.Parse(purl.TypeOCI, ctx.Image.Base, baseDigest, DigestSha256, ""),
			Digest: common.DigestSet{DigestSha256: baseDigest},
		},
		{
			URI:    purl.Parse(purl.TypeOCI, ctx.Image.Parent, parentDigest, DigestSha256, ""),
			Digest: common.DigestSet{DigestSha256: parentDigest},
		},
	}

	subjects := []in_toto.Subject{
		{
			Name:   purl.Parse(purl.TypeOCI, ref, digest, DigestSha256, ""),
			Digest: common.DigestSet{DigestSha256: digest},
		},
	}

	// collect environment variables
	env := map[string]string{}
	for _, i := range os.Environ() {
		k, _, _ := strings.Cut(i, "=")
		env[k] = ""
	}

	env["entryPoint"] = "ci"
	env["source"] = repoURL

	// parse times
	buildStart, err := time.Parse(time.RFC3339, ctx.StartTime)
	if err != nil {
		return err
	}
	buildEnd := time.Now()

	buildType := DefaultBuildType
	if val := os.Getenv(civ1.EnvBuildSLSABuildType); val != "" {
		buildType = val
	}

	configHash, err := digestof.File(ctx.ConfigPath)
	if err != nil {
		return err
	}

	// write the digest file
	outPath := filepath.Join(ctx.Root, outBuild)
	if err := os.WriteFile(outPath, []byte(fmt.Sprintf("%s@sha256:%s", ref, digest)), 0644); err != nil {
		return err
	}

	outHash, err := digestof.File(outPath)
	if err != nil {
		return err
	}

	predicate := v1.ProvenancePredicate{
		BuildDefinition: v1.ProvenanceBuildDefinition{
			BuildType:          buildType,
			ExternalParameters: env,
			InternalParameters: map[string]any{
				"commands": os.Args,
				"build":    append([]string{r.Command}, r.Args...),
				"shell":    os.Getenv("SHELL"),
			},
			ResolvedDependencies: materials,
		},
		RunDetails: v1.ProvenanceRunDetails{
			Builder: v1.Builder{
				ID: ctx.Builder,
			},
			BuildMetadata: v1.BuildMetadata{
				InvocationID: ctx.BuildID,
				StartedOn:    &buildStart,
				FinishedOn:   &buildEnd,
			},
			Byproducts: []v1.ResourceDescriptor{
				{
					URI:    ctx.ConfigPath,
					Digest: common.DigestSet{DigestSha256: configHash},
				},
				{
					URI:    outPath,
					Digest: common.DigestSet{DigestSha256: outHash},
				},
			},
		},
	}

	provenance := in_toto.ProvenanceStatementSLSA1{
		StatementHeader: in_toto.StatementHeader{
			Type:          in_toto.StatementInTotoV01,
			PredicateType: v1.PredicateSLSAProvenance,
			Subject:       subjects,
		},
		Predicate: predicate,
	}

	var data any = provenance
	if predicateOnly {
		data = predicate
	}

	return output(ctx, &data, digest)
}
