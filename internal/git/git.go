package git

import (
	"os/exec"
	"strconv"
	"strings"
)

var (
	RepositoryPath string
	CommitHash     string
	UseCommitter   bool
)

func GetGitFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-tree", CommitHash, "--name-only", "-r")
	cmd.Dir = RepositoryPath
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(out), "\n"), nil
}

type FileStats struct {
	AuthorCommits map[string]map[string]struct{}
	AuthorLineCnt map[string]int
}

func CalculateFileStats(file string) (FileStats, error) {
	cmd := exec.Command("git", "blame", "--porcelain", file, CommitHash)
	cmd.Dir = RepositoryPath
	out, err := cmd.Output()
	if err != nil {
		return FileStats{}, err
	}
	// git blame does not handle empty files properly
	if len(out) == 0 {
		return handleEmptyFile(file, CommitHash)
	}
	lines := strings.Split(string(out), "\n")

	commitLineCnt := make(map[string]int)
	readAuthor, numReadLines := false, 0
	curHash := ""
	stats := FileStats{
		AuthorCommits: make(map[string]map[string]struct{}),
		AuthorLineCnt: make(map[string]int),
	}

	for _, line := range lines {
		switch {
		case line == "":
			continue
		case numReadLines == 0:
			args := strings.Split(line, " ")
			curHash = args[0]
			numReadLines, _ = strconv.Atoi(args[len(args)-1])
			commitLineCnt[curHash] += numReadLines
			readAuthor = true
		case line[0] == '\t':
			numReadLines--
		case readAuthor:
			args := strings.Split(line, " ")
			if (!UseCommitter && args[0] == "author") || (UseCommitter && args[0] == "committer") {
				readAuthor = false
				author := strings.Join(args[1:], " ")
				if _, ok := stats.AuthorCommits[author]; !ok {
					stats.AuthorCommits[author] = make(map[string]struct{})
				}
				stats.AuthorCommits[author][curHash] = struct{}{}
			}
		}
	}

	for author, commits := range stats.AuthorCommits {
		for c, num := range commitLineCnt {
			if _, ok := commits[c]; ok {
				stats.AuthorLineCnt[author] += num
			}
		}
	}

	return stats, nil
}

func handleEmptyFile(file, commit string) (FileStats, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%an%n%H", commit, "--", file)
	cmd.Dir = RepositoryPath
	out, err := cmd.Output()
	if err != nil {
		return FileStats{}, err
	}

	result := strings.Split(string(out), "\n")
	author := result[0]
	commit = result[1]
	return FileStats{
		AuthorCommits: map[string]map[string]struct{}{author: {commit: struct{}{}}},
		AuthorLineCnt: map[string]int{author: 0},
	}, nil
}
