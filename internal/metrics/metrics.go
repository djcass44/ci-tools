package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var DefaultLabels = []string{"recipe", "repo_url", "provider", "context"}

var (
	MetricOpsBuilds = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ci_ops_builds_total",
	}, DefaultLabels)
	MetricOpsBuildSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ci_ops_build_success_total",
	}, DefaultLabels)
	MetricOpsBuildErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ci_ops_build_errors_total",
	}, DefaultLabels)
	MetricOpsBuildDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ci_ops_build_duration",
		Buckets: prometheus.LinearBuckets(0.1, 1, 10),
	})
	MetricOpsProvenanceDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ci_ops_provenance_duration",
		Buckets: prometheus.LinearBuckets(0.1, 1, 10),
	})
	MetricOpsDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ci_ops_duration",
		Buckets: prometheus.LinearBuckets(0.1, 1, 10),
	})
	MetricDockerCfgGenerated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ci_ops_dockercfg_total",
	}, DefaultLabels)
	MetricOciVerify = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ci_oci_verify_total",
	}, []string{"type"})
	MetricBOMGenerated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ci_sbom_generated_total",
	})
	MetricProvenanceGenerated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ci_provenance_generated_total",
	}, []string{"version"})
)
