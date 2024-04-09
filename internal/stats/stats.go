package stats

import (
	"fmt"
	"sort"
	"strings"

	"github.com/themeelanoid/report-sitory/internal/git"
)

type StatLine struct {
	Author  string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func CalculateStats(files []string, order string) ([]StatLine, error) {
	authorCommits := make(map[string]map[string]struct{})
	authorLines := make(map[string]int)
	authorFiles := make(map[string]int)

	for _, file := range files {
		if file == "" {
			continue
		}
		fileStats, err := git.CalculateFileStats(file)
		if err != nil {
			return nil, err
		}
		for author, commits := range fileStats.AuthorCommits {
			if _, ok := authorCommits[author]; !ok {
				authorCommits[author] = make(map[string]struct{})
			}
			authorFiles[author]++
			for commit := range commits {
				authorCommits[author][commit] = struct{}{}
			}
		}
		for author, Lines := range fileStats.AuthorLineCnt {
			authorLines[author] += Lines
		}
	}

	var stats []StatLine
	for author, commits := range authorCommits {
		line := StatLine{Author: author, Lines: authorLines[author],
			Commits: len(commits), Files: authorFiles[author]}
		stats = append(stats, line)
	}

	err := orderStats(stats, order)
	return stats, err
}

func orderStats(stats []StatLine, order string) error {
	var comparator func(i, j int) bool
	switch order {
	case "lines":
		comparator = func(i, j int) bool {
			if stats[i].Lines != stats[j].Lines {
				return stats[i].Lines > stats[j].Lines
			}
			if stats[i].Commits != stats[j].Commits {
				return stats[i].Commits > stats[j].Commits
			}
			if stats[i].Files != stats[j].Files {
				return stats[i].Files > stats[j].Files
			}
			return strings.Compare(stats[i].Author, stats[j].Author) == -1
		}
	case "commits":
		comparator = func(i, j int) bool {
			if stats[i].Commits != stats[j].Commits {
				return stats[i].Commits > stats[j].Commits
			}
			if stats[i].Lines != stats[j].Lines {
				return stats[i].Lines > stats[j].Lines
			}
			if stats[i].Files != stats[j].Files {
				return stats[i].Files > stats[j].Files
			}
			return strings.Compare(stats[i].Author, stats[j].Author) == -1
		}
	case "files":
		comparator = func(i, j int) bool {
			if stats[i].Files != stats[j].Files {
				return stats[i].Files > stats[j].Files
			}
			if stats[i].Lines != stats[j].Lines {
				return stats[i].Lines > stats[j].Lines
			}
			if stats[i].Commits != stats[j].Commits {
				return stats[i].Commits > stats[j].Commits
			}
			return strings.Compare(stats[i].Author, stats[j].Author) == -1
		}
	default:
		return fmt.Errorf("unsupported ouput ordering: %s", order)
	}
	sort.Slice(stats, comparator)
	return nil
}
