# How to contribute to k6registry

Thank you for your interest in contributing to **k6registry**!

Before you begin, make sure to familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md). If you've previously contributed to other open source project, you may recognize it as the classic [Contributor Covenant](https://contributor-covenant.org/).


## Prerequisites

The tools listed in the [tools] section should be installed before contributing. It is advisable to first install the [cdo] tool, which can be used to easily perform the tasks described here. The [cdo] tool can most conveniently be installed using the [eget] tool.

```bash
eget szkiba/cdo
```

The [cdo] tool can then be used to perform the tasks described in the following sections.

Help about tasks:

```
cdo
```

[cdo]: (https://github.com/szkiba/cdo)
[eget]: https://github.com/zyedidia/eget

## tools - Install the required tools

Contributing will require the use of some tools, which can be installed most easily with a well-configured [eget] tool.

```bash
eget mikefarah/yq
eget atombender/go-jsonschema
eget szkiba/mdcode
eget golangci/golangci-lint
eget goreleaser/goreleaser
```

## schema - Update schema and generate code

After modifying registry schema ([registry.schema.json]), the [registry_gen.go] file must be regenerated.

```bash
curl -s -o internal/registry/registry.schema.json https://raw.githubusercontent.com/grafana/k6-extension-registry/main/registry.schema.json
go-jsonschema --capitalization URL --capitalization OSS -p registry --only-models -o internal/registry/registry_gen.go internal/registry/registry.schema.json
```

[registry.schema.json]: internal/registry/registry.schema.json
[registry_gen.go]: internal/registry/registry_gen.go

## readme - Update README.md

After changing the CLI tool or example registry, the documentation must be updated in README.md.

```bash
go run ./tools/gendoc README.md
```

### lint - Run the linter

The `golangci-lint` tool is used for static analysis of the source code.
It is advisable to run it before committing the changes.

```bash
golangci-lint run
```

### test - Run the tests

```bash
go test -count 1 -race -coverprofile=coverage.txt ./...
```

[test]: <#test---run-the-tests>

### coverage - View the test coverage report

Requires
: [test]

```bash
go tool cover -html=coverage.txt
```

### build - Build the executable binary

This is the easiest way to create an executable binary (although the release process uses the `goreleaser` tool to create release versions).

```bash
go build -ldflags="-w -s" -o k6dist ./cmd/k6dist
```

[build]: <#build---build-the-executable-binary>

### snapshot - Creating an executable binary with a snapshot version

The goreleaser command-line tool is used during the release process. During development, it is advisable to create binaries with the same tool from time to time.

```bash
rm -f k6registry
goreleaser build --snapshot --clean --single-target -o k6dist
```

[snapshot]: <#snapshot---creating-an-executable-binary-with-a-snapshot-version>

### docker - Build docker image

Building a Docker image. Before building the image, it is advisable to perform a snapshot build using goreleaser. To build the image, it is advisable to use the same `Docker.goreleaser` file that `goreleaser` uses during release.

Requires
: snapshot

```bash
docker build -t k6dist -f Dockerfile.goreleaser .
```

### clean - Delete the build directory

```bash
rm -rf build
```

### update - Update everything

The most robust thing is to update everything (both the schema and the example) after modifying the source.

Requires
: schema, readme

