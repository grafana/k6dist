// Package cmd contains run cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/grafana/k6dist"
	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

type options struct {
	k6dist.Options
	quiet     bool
	verbose   bool
	single    bool
	platforms []string
}

var defaultPlatforms = []string{ //nolint:gochecknoglobals
	"linux/amd64",
	"darwin/amd64",
	"windows/amd64",
}

// New creates new cobra command for exec command.
func New(levelVar *slog.LevelVar) *cobra.Command {
	opts := new(options)

	root := &cobra.Command{
		Use:               "k6dist [flags] [registry-location]",
		Short:             "k6 distro builder",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Args:              cobra.MaximumNArgs(1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			return preRun(args, levelVar, opts)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			changed, err := run(cmd, opts)
			if err != nil {
				return err
			}

			return emitOutput(changed, opts.Version)
		},
	}

	root.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s\n" .Version}}`)

	flags := root.Flags()

	flags.SortFlags = false

	flags.StringVar(&opts.Name, "distro-name", "", "distro name (default detect)")
	flags.StringVar(&opts.Version, "distro-version", "", "distro version (default current date in YY.MM.DD format)")
	flags.StringSliceVar(&opts.platforms, "platform", defaultPlatforms, "target platforms")
	flags.StringVar(&opts.Executable, "executable", k6dist.DefaultExecutableTemplate, "executable file name")
	flags.StringVar(&opts.Archive, "archive", k6dist.DefaultArchiveTemplate, "archive file name")
	flags.StringVar(&opts.Dockerfile, "docker", "", "generated Dockerfile name (default Dockerfile next to the exe)")
	flags.StringVar(&opts.DockerfileTemplate, "docker-template", "",
		"template for the generated Dockerfile (default embedded)")
	flags.StringVar(&opts.Notes, "notes", k6dist.DefaultNotesTemplate, "generated release notes file name")
	flags.StringVar(&opts.NotesTemplate, "notes-template", "",
		"template for the generated release notes (default embedded)")
	flags.StringVar(&opts.NotesLatest, "notes-latest", "", "latest release notes file for change detection")
	flags.StringVar(&opts.Readme, "readme", "", "readme file to be added to the archive (default detect)")
	flags.StringVar(&opts.License, "license", "", "license file to be added to the archive (default detect)")
	flags.BoolVar(&opts.single, "single-target", false, "build only for the current runtime platform")
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "enable verbose logging")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "disable normal logging")

	flags.BoolP("version", "V", false, "print version")

	return root
}

func preRun(args []string, levelVar *slog.LevelVar, opts *options) error {
	if len(args) > 0 {
		opts.Registry = args[0]
	} else {
		opts.Registry = k6dist.DefaultRegistryURL
	}

	if opts.verbose && levelVar != nil {
		levelVar.Set(slog.LevelDebug)
	}

	if opts.quiet && levelVar != nil {
		levelVar.Set(slog.LevelWarn)
	}

	if len(opts.Name) == 0 {
		opts.Name = guessName(opts.Registry)
	}

	if len(opts.Version) == 0 {
		opts.Version = defaultDate()
	}

	if len(opts.Readme) == 0 {
		opts.Readme = findReadme()
	}

	if len(opts.License) == 0 {
		opts.License = findLicense()
	}

	if opts.single {
		opts.platforms = []string{runtime.GOOS + "/" + runtime.GOARCH}
	} else {
		var err error

		opts.Platforms, err = parsePlatforms(opts.platforms)
		if err != nil {
			return err
		}
	}

	return nil
}

func run(_ *cobra.Command, opts *options) (bool, error) {
	if len(opts.Name) == 0 {
		cwd, err := os.Getwd() //nolint:forbidigo
		if err != nil {
			return false, err
		}

		opts.Name = filepath.Base(cwd)
	}

	changed, err := k6dist.Build(context.TODO(), &opts.Options)
	if err != nil {
		return false, err
	}

	return changed, nil
}

func guessName(source string) string {
	basename := func(filename string) string {
		base := path.Base(filename)

		if ext := path.Ext(base); len(ext) != 0 {
			base = strings.TrimSuffix(base, ext)
		}

		return base
	}

	base := basename(filepath.ToSlash(source))

	if base != "registry" {
		return "k6-" + base
	}

	cwd, err := os.Getwd() //nolint:forbidigo
	if err != nil {
		return ""
	}

	return basename(filepath.ToSlash(cwd))
}

func findTextFile(basename string) string {
	for _, ext := range []string{".md", "", ".txt"} {
		filename := basename + ext

		if _, err := os.Stat(filename); err == nil { //nolint:forbidigo
			return filename
		}
	}

	return ""
}

func findReadme() string {
	return findTextFile("README")
}

func findLicense() string {
	return findTextFile("LICENSE")
}

func defaultDate() string {
	now := time.Now()

	return now.Format("06.01.02")
}

func parsePlatforms(values []string) ([]*k6dist.Platform, error) {
	platforms := make([]*k6dist.Platform, 0, len(values))

	for _, value := range values {
		platform, err := k6dist.ParsePlatform(value)
		if err != nil {
			return nil, err
		}

		platforms = append(platforms, platform)
	}

	return platforms, nil
}
