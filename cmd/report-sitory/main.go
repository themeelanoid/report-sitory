//go:build !solution

package main

import (
	"log"

	"github.com/spf13/pflag"

	"github.com/themeelanoid/report-sitory/internal/files"
	"github.com/themeelanoid/report-sitory/internal/git"
	"github.com/themeelanoid/report-sitory/internal/stats"
	"github.com/themeelanoid/report-sitory/internal/writers"
)

var (
	repository   = pflag.String("repository", ".", "Path to the repository.")
	revision     = pflag.String("revision", "HEAD", "Commit hash with which to perform calculations.")
	orderBy      = pflag.String("order-by", "lines", "Parameter to order the contributors by (lines, commits, files).")
	useCommitter = pflag.Bool("use-committer", false, "Use committer instead of author in calculations.")
	format       = pflag.String("format", "tabular", "Format of the output (tabular, csv, json, json-lines).")
	extensions   = pflag.StringSlice("extensions", []string{}, "File extensions to include in calculations. Separated by comma.")
	langs        = pflag.StringSlice("languages", []string{}, "Programming languages' files to include in calculations. Separated by comma.")
	exclude      = pflag.StringSlice("exclude", []string{}, "File name patterns to exclude from calculations. Separated by comma.")
	restrict     = pflag.StringSlice("restrict-to", []string{}, "File name patterns to use in calculations. Separated by comma.")
)

func main() {
	pflag.Parse()

	git.RepositoryPath = *repository
	git.CommitHash = *revision
	git.UseCommitter = *useCommitter

	relevantFiles, err := files.ListFiles(*langs, *extensions, *exclude, *restrict)
	if err != nil {
		log.Fatal(err)
	}

	statistics, err := stats.CalculateStats(relevantFiles, *orderBy)
	if err != nil {
		log.Fatal(err)
	}

	err = writers.WriteStatistics(statistics, *format)
	if err != nil {
		log.Fatal(err)
	}
}
