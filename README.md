# k6dist

**k6 distro builder**

k6dist is a CLI tool and a GitHub Action that allows you to build k6 distributions with extensions.

```bash
k6dist https://registry.k6.io/tier/official.json
```

The extensions to be embedded can be specified with an extension registry. The location of the extension registry can be specified with HTTP(S) URL or filesystem path.

k6dist builds k6 for the specified platforms with the extensions listed in the registry. After that, it also creates archive files corresponding to the platform (`.zip` for Windows, otherwise `.tar.gz` format).

In addition to the executable, the readme and the license file can be included in the archive.

A release notes file is also prepared, in which the extensions are listed with their current and previous versions and their description. The content of the release notes file can be customized, using a go template file.

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6dist/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6dist/cmd/k6dist@latest
```

## GitHub Action

`grafana/k6dist` is a GitHub Action that enables the building of k6 distibutions with extensions.

The extensions to be embedded can be specified with an extension registry. The location of the extension registry can be specified with HTTP(S) URL or filesystem path.

**Inputs**

See the CLI [Flags](#flags) section for the default values ​​for these parameters.

name           | reqired | description
---------------|---------|-------------
in             |    yes  | extension registry location (URL or file path)
distro_name    |    no   | distro name
distro_version |    no   | distro version
platform       |    no   | list of target platforms
executable     |    no   | executable file name template
archive        |    no   | archive file name template
docker         |    no   | generated Dockerfile name template
docker_template|    no   | template for the generated Dockerfile
notes          |    no   | generated release notes file name template
notes_template |    no   | template for the generated release notes
notes_latest   |    no   | latest release notes file for change detection
readme         |    no   | readme file to be added to the archive
license        |    no   | license file to be added to the archive
verbose        |    no   | enable verbose logging
quiet          |    no   | disable normal logging

The change can be indicated by comparing the generated release notes to the latest release notes. The location of the latest release notes path can be passed in the `notes_latest` action parameter. The `changed` output variable will be `true` or `false` depending on whether the release has changed or not compared to the latest release notes. If the release has not changed, the build step is skipped.

**Outputs**

name    | description
--------|------------
changed | `true` if the release has changed compared to `notes_latest`, otherwise `false`
version | the distro release version

**Example usage**

```yaml
- name: Build k6 distro
  uses: grafana/k6dist@v0.1.0
  with:
    in: "https://registry.k6.io/product/cloud.json"
    notes_latest: "latest-release-notes.md"
```

## CLI

<!-- #region cli -->
## k6dist

k6 distro builder

### Synopsis

Build k6 distribution with extensions.

k6dist enables building k6 distributions with extensions. The extensions to be embedded can be specified with an extension registry. The location of the registry can be specified with http(s) URL or filesystem path.

The location of the registry can be specified as a command line argument. If omitted, the URL of the registry subset containing the official k6 extensions is the default.

Each flag has a meaningful default value.

k6dist builds k6 for the specified platforms with the extensions listed in the registry. After that, it also creates archive files corresponding to the platform (`.zip` for Windows, otherwise `.tar.gz` format).

In addition to the executable, the readme and the license file can be included in the archive. Their names are detected, but can also be specified using the `--readme` and `--license` flags.

The locations of executables and archives can be specified with the go template expression using the `--executable` and `--archive` flags.

A release notes file is also prepared, in which the extensions are listed with their current and previous versions and their description. The release notes file can be customized, using the `--notes-template` flag you can specify a go template file to generate the release notes.


```
k6dist [flags] [registry-location]
```

### Flags

```
      --distro-name string       distro name (default detect)
      --distro-version string    distro version (default current date in YY.MM.DD format)
      --platform strings         target platforms (default [linux/amd64,darwin/amd64,windows/amd64])
      --executable string        executable file name (default "dist/{{.Name}}_{{.OS}}_{{.Arch}}/k6{{.ExeExt}}")
      --archive string           archive file name (default "dist/{{.Name}}_{{.Version}}_{{.OS}}_{{.Arch}}{{.ZipExt}}")
      --docker string            generated Dockerfile name (default Dockerfile next to the exe)
      --docker-template string   template for the generated Dockerfile (default embedded)
      --notes string             generated release notes file name (default "dist/{{.Name}}_{{.Version}}.md")
      --notes-template string    template for the generated release notes (default embedded)
      --notes-latest string      latest release notes file for change detection
      --readme string            readme file to be added to the archive (default detect)
      --license string           license file to be added to the archive (default detect)
      --single-target            build only for the current runtime platform
  -v, --verbose                  enable verbose logging
  -q, --quiet                    disable normal logging
  -V, --version                  print version
  -h, --help                     help for k6dist
```

<!-- #endregion cli -->

## Contribure 

If you want to contribute, start by reading [CONTRIBUTING.md](CONTRIBUTING.md).
