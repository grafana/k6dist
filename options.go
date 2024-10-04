// Package k6dist contains k6dists's go API.
package k6dist

import (
	"runtime"
)

// Options contains the optional parameters of the Build function.
type Options struct {
	// Name contains short name of the distribution.
	// Templating is not supported.
	Name string
	// Version contains distribution version.
	// Templating is not supported.
	Version string
	// Executable is the name of the k6 executable file to be built.
	// Templating is supported.
	// It defaults to DefaultExecutableTemplate.
	Executable string
	// Archive is the name of the archive file to be created.
	// Templating is supported.
	// It defaults to DefaultArchiveTemplate.
	Archive string
	// Readme is the name of the readme file to be added.
	// Templating is not supported.
	// No readme will be added if it is empty.
	Readme string
	// License is the name of the license file to be added.
	// Templating is not supported.
	// No license will be added if it is empty.
	License string
	// NotesTemplate is the name of the release notes (go) template file.
	// Templating is not supported in filename.
	// Default is an embedded template.
	NotesTemplate string
	// Notes is the name of the generated release notes file.
	// Templating is supported.
	// It defaults to DefaultNotesTemplate.
	Notes string
	// DockerfileTemplate is the name of the Dockerfile template.
	// Templating is not supported in filename.
	// Default is an embedded template.
	DockerfileTemplate string
	// Dockerfile is the name of the generated Dockerfile.
	// Templating is supported.
	// It defaults to Dockerfile in the same directory as the executable.
	Dockerfile string
	// NotesLatest contains the name of the latest release notes file.
	// It is used for change detection, if the release notes have not changed, no new release is made.
	// Within the release notes, an HTML comment contains the module list.
	// Change detection is done by comparing this.
	// Templating is not supported.
	NotesLatest string
	// Platforms contains the target platforms.
	// If absent, the current runtime platform will be used.
	Platforms []*Platform
	// Registry contains the location of the registry. Its value is http(s) URL or filesystem path.
	// Templating is not supported.
	// It defaults to the URL of the registry subset containing the official k6 extensions.
	Registry string
}

const (
	// DefaultExecutableTemplate is a default go template for executable file path.
	DefaultExecutableTemplate = `dist/{{.Name}}_{{.OS}}_{{.Arch}}/k6{{.ExeExt}}`
	// DefaultNotesTemplate is a default go template for release notes file path.
	DefaultNotesTemplate = `dist/{{.Name}}_{{.Version}}.md`
	// DefaultArchiveTemplate is a default go template for archive file path.
	DefaultArchiveTemplate = `dist/{{.Name}}_{{.Version}}_{{.OS}}_{{.Arch}}{{.ZipExt}}`
	// DefaultRegistryURL is a default registry URL.
	DefaultRegistryURL = `https://registry.k6.io/tier/official.json`
)

func (o *Options) setDefaults() {
	if len(o.Executable) == 0 {
		o.Executable = DefaultExecutableTemplate
	}

	if len(o.Notes) == 0 {
		o.Notes = DefaultNotesTemplate
	}

	if len(o.Archive) == 0 {
		o.Archive = DefaultArchiveTemplate
	}

	if len(o.Registry) == 0 {
		o.Registry = DefaultRegistryURL
	}

	if o.Platforms == nil {
		o.Platforms = []*Platform{{OS: runtime.GOOS, Arch: runtime.GOARCH}}
	}
}
