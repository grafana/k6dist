// Package main contains the main function for k6dist.
package main

import (
	"log/slog"
	"os"

	"github.com/grafana/k6dist/cmd"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	appname = "k6dist" //nolint:unused
	version = "dev"
)

func initLogging() *slog.LevelVar {
	levelVar := new(slog.LevelVar)

	logrus.SetLevel(logrus.DebugLevel)

	logger := slog.New(sloglogrus.Option{Level: levelVar}.NewLogrusHandler())

	slog.SetDefault(logger)

	return levelVar
}

func main() {
	runCmd(newCmd(getArgs(), initLogging()))
}

func newCmd(args []string, levelVar *slog.LevelVar) *cobra.Command {
	cmd := cmd.New(levelVar)

	cmd.Version = version
	cmd.SetArgs(args)

	return cmd
}

func runCmd(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1) //nolint:forbidigo
	}
}

func getArgs() []string {
	return cmd.AddGitHubArgs(os.Args[1:]) //nolint:forbidigo
}
