# traci

OpenTelemetry tracing for CI pipelines

Traci is a CI command wrapper that generates and uploads traces in OpenTelemetry format. It uses the OpenTelemetry SDK
under the hood and can be configured to send traces to any OpenTelemetry-compatible backend.

Traci supports CI detection for the following CI platforms:

- GitHub Actions
- GitLab CI
- CircleCI
- Travis CI
- Bitbucket

If you don't see your CI provider, open an issue or submit a PR!

## Quickstart

```bash
go get github.com/nextrevision/traci

traci exec echo "hello world"
```

Traci detects the CI environment and automatically generates trace and span IDs for each command. It will also detect
the OTel endpoint using the `OTEL_EXPORTER_OTLP_ENDPOINT` environment variable.

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT=https://jaeger.mycompany:4317
export OTEL_EXPORTER_OTLP_PROTOCOL=grpc

traci exec echo "hello jaeger"
```

## Configuration

Traci supports configuration via environment variables in both `exec` and `execf` commands. The `execf` command also
supports CLI flags. CLI flags take precedence over environment variables.

### Traci Config

| Environment Variable     | `execf` CLI Flag     | Description                                                  |
|--------------------------|----------------------|--------------------------------------------------------------|
| `TRACI_SERVICE_NAME`     | `--service-name`     | The name of the service                                      |
| `TRACI_SPAN_NAME`        | `--span-name`        | The name of the span                                         |
| `TRACI_TRACE_BOUNDARY`   | `--trace-boundary`   | The scope of the generated trace. Can be `pipeline` or `job` |
| `TRACI_TAG_COMMAND_ARGS` | `--tag-command-args` | Include command args as tags in the span                     |

### OpenTelemetry Config

Traci uses the OpenTelemetry Go SDK under the hood and can be configured using the same environment variables. See the
[OpenTelemetry General SDK documentation](https://opentelemetry.io/docs/concepts/sdk-configuration/general-sdk-configuration/)
for a full list of available configuration options. Below are some of the more common options.

| Environment Variable          | Description                                        | Example Value                                    |
|-------------------------------|----------------------------------------------------|--------------------------------------------------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | The OTel endpoint to send traces to.               | `https://jaeger.mycompany:4317`                  |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | The OTel protocol to use. Can be `grpc` or `http`. | `grpc`                                           |
| `OTEL_RESOURCE_ATTRIBUTES`    | Resource attributes to include in the trace.       | `service.namespace=tutorial,service.version=1.0` |
| `OTEL_EXPORTER_OTLP_HEADERS`  | Headers to include in the request.                 | `x-something=foo,x-something-else=bar`           |

## Propagation

### CI Deterministic Trace IDs

Traci detects CI providers by environment variables. Most CI providers offer a consistent, unique identifier per
"pipeline" or "workflow". Traci uses these identifiers to generate repeatable trace IDs without the need for manual
propagation via some other means. This means commands in different jobs/steps of a pipeline/workflow will associate
with the same trace ID automatically.

#### Trace Boundary

Traci can use the CI deterministic trace ID to determine the trace boundary. The trace boundary determines how far a
trace will propagate. For example, if the trace boundary is set to `pipeline`, the trace will include all wrapped commands
in every job in a pipeline. If the trace boundary is set to `job`, the trace will only include all commands in the job.
The default trace boundary is `pipeline`.

`Pipeline` is a generic term to describe a grouping of jobs. A pipeline includes the following CI vendor concepts:

- GitHub Actions: Workflow
- GitLab CI: Pipeline
- CircleCI: Workflow
- Travis CI: Build
- Bitbucket: Pipeline

`Job` is a generic term to describe a grouping of steps. A job includes the following CI vendor concepts:

- GitHub Actions: Job
- GitLab CI: Job
- CircleCI: Job
- Travis CI: Job
- Bitbucket: Step

### `TRACEPARENT` Environment Variable

Traci supports propagating trace context between commands using the `TRACEPARENT` environment variables. If a valid
trace context is set in the `TRACEPARENT` variable, traci will use that to determine the parent trace and span ids. The
`exec` and `execf` commands will also inject a [W3C Trace Context compatible](https://www.w3.org/TR/trace-context/)
`TRACEPARENT` variable into the environment of the command being executed. This allows traces to be propagated between
commands by traci or other tools that look for that variable. For example:

```bash
traci exec /bin/sh -c 'traci exec /bin/sh -c "echo $TRACEPARENT"'
```

## Commands

## `traci exec`

The `traci exec` commands does not support any traci CLI options opting to pass all args to the command being executed. It
can be configured via environment variables, however. For example:

```bash
export TRACI_SPAN_NAME=foo
traci exec echo "hello world"
# or
TRACI_SPAN_NAME=bar traci exec echo "hello world"
```

## `traci execf`

The `traci execf` command does parse traci CLI options and requires the use of the bash end of command-line options `--`
separator to pass args to the command being executed. This command can also be configured via environment variables. For
example:

```bash
traci execf --span-name foo -- echo "hello world"
```

## Examples

### GitLab CI

```yaml
stages:
  - test
  - build
  - deploy

variables:
  OTEL_EXPORTER_OTLP_ENDPOINT: https://jaeger:4317
  OTEL_EXPORTER_OTLP_PROTOCOL: grpc

default:
  image: alpine
  before_script:
    - wget -O /tmp/traci.tgz https://github.com/nextrevision/traci/releases/download/v0.3.0/traci_0.3.0_linux_amd64.tar.gz && tar -C /usr/local/bin -xzf /tmp/traci.tgz traci
    - traci detect

test:
  stage: test
  script:
    - traci execf --tag-command-args -- /bin/sh -c 'echo $TRACEPARENT'
    - traci exec sleep 1

build:
  stage: build
  script:
    - traci exec env
    - traci exec sleep 2

deploy:
  stage: deploy
  script:
    - trace exec /bin/sh -c 'traci /bin/sh -c "echo $TRACEPARENT"'
    - traci exec sleep 3
```

## Inspiration

Traci is heavily inspired by the following projects:

- [buildevents](https://github.com/honeycombio/buildevents)
- [otel-cli](https://github.com/equinix-labs/otel-cli)
- [Jenkins otel plugin](https://github.com/jenkinsci/opentelemetry-plugin/tree/master)
- [tracepusher](https://github.com/agardnerIT/tracepusher)
