package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/spf13/pflag"
)

// AddGitHubArgs adds GitHub input parameters as arguments.
func AddGitHubArgs(args []string, flags *pflag.FlagSet) ([]string, error) {
	if !isGitHubAction() {
		return args, nil
	}

	args = make([]string, 0)

	flags.Visit(func(flag *pflag.Flag) {
		args = ghinput(args, flag)
	})

	if iargs := os.Getenv("INPUT_ARGS"); len(iargs) > 0 { //nolint:forbidigo
		items, err := shlex.Split(iargs)
		if err != nil {
			return nil, err
		}

		args = append(args, items...)
	}

	return args, nil
}

//nolint:forbidigo
func isGitHubAction() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func ghinput(args []string, flag *pflag.Flag) []string {
	name := "INPUT_" + strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_"))

	if val := os.Getenv(name); len(val) > 0 { //nolint:forbidigo
		if flag.Value.Type() == "bool" {
			if val == "true" {
				args = append(args, "--"+flag.Name)
			}
		} else {
			args = append(args, "--"+flag.Name, val)
		}
	}

	return args
}

//nolint:forbidigo
func emitOutput(changed bool, version string) error {
	ghOutput := os.Getenv("GITHUB_OUTPUT")
	if len(ghOutput) == 0 {
		return nil
	}

	file, err := os.Create(filepath.Clean(ghOutput))
	if err != nil {
		return err
	}

	slog.Debug("Emit changed", "changed", changed)

	_, err = fmt.Fprintf(file, "changed=%t\n", changed)
	if err != nil {
		return err
	}

	slog.Debug("Emit version", "version", version)

	_, err = fmt.Fprintf(file, "version=%s\n", version)
	if err != nil {
		return err
	}

	return file.Close()
}
