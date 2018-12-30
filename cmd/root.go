package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"github.com/gofunct/common/bingen"
	"github.com/gofunct/common/errors"
)


const (
	cliName = "bingen"
	version = "v0.1.0"
)

var (
	pkgsToBeAdded []string
	flagBuild     bool
	flagVersion   bool
	flagVerbose   bool
	flagHelp      bool
)

func init() {
	BinCmd.PersistentFlags().SetInterspersed(false)
	BinCmd.PersistentFlags().StringArrayVar(&pkgsToBeAdded, "add", []string{}, "Add new tools")
	BinCmd.PersistentFlags().BoolVar(&flagBuild, "build", false, "Build all tools")
	BinCmd.PersistentFlags().BoolVar(&flagVersion, "version", false, "Print the CLI version")
	BinCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Verbose level output")
	BinCmd.PersistentFlags().BoolVarP(&flagHelp, "help", "h", false, "Help for the CLI")
	BinCmd.AddCommand(genCmd)
}


var BinCmd = &cobra.Command{
	Use: "bingen",
	Short: "bingen opts",
}

var genCmd = &cobra.Command{
	Use: "gen",
	Short: "generate needed executables for your app",
	Run: func(cmd *cobra.Command, args []string) {
		var exitCode int

		if err := run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}

		os.Exit(exitCode)
	},
}

func Execute() {
	if err := BinCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	pflag.Parse()
	args := pflag.Args()

	var cfg bingen.Config
	if flagVerbose {
		cfg.Verbose = true
		cfg.Logger = log.New(os.Stderr, "", 0)
	}
	toolRepo, err := cfg.Create()
	if err != nil {
		return errors.WithStack(err)
	}

	ctx := context.TODO()

	switch {
	case len(pkgsToBeAdded) > 0:
		err = toolRepo.Add(ctx, pkgsToBeAdded...)
	case flagVersion:
		fmt.Fprintf(os.Stdout, "%s %s\n", cliName, version)
	case flagHelp:
		printHelp(os.Stdout)
	case flagBuild:
		err = toolRepo.BuildAll(ctx)
	case len(args) > 0:
		err = toolRepo.Run(ctx, args[0], args[1:]...)
	default:
		printHelp(os.Stdout)
	}

	return errors.WithStack(err)
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, helpText)
	pflag.PrintDefaults()
}

var (
	helpText = `The implementation of clarify best practice for tool dependencies.

See https://github.com/golang/go/issues/25922#issuecomment-412992431

Usage:
  bingen [command] [args]
  bingen --add [packages...]

Flags:`
)
