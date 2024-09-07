# Metrics

The application can optional export Prometheus metrics.
Due to being a CLI application it requires a Prometheus Push Gateway to be setup.

## Getting started

In your CI environment, set the following environment variables:

* `PROMETHEUS_PUSH_URL` - HTTP(S) URL to the Prometheus Push Gateway
* `PROMETHEUS_JOB_NAME` - name of the Job to be recorded by Prometheus.

Once these variables have been set, the application will push metrics when it exits.

## Available metrics

The full list of available metrics can be found in [`metrics.go`](../internal/metrics/metrics.go)
