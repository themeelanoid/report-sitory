package files

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/themeelanoid/report-sitory/configs"
	"github.com/themeelanoid/report-sitory/internal/git"
)

func ListFiles(langs, extensions, exclude, restrict []string) ([]string, error) {
	files, err := git.GetGitFiles()
	if err != nil {
		return files, err
	}

	files, err = filterExtensions(files, langs, extensions)
	if err != nil {
		return files, err
	}

	files, err = filterPatterns(files, exclude, restrict)
	return files, err
}

type lang struct {
	Name       string   `json:"name"`
	Extensions []string `json:"extensions"`
}

func extsFromLangs(langs []string) (map[string]struct{}, error) {
	var allLangs []lang
	err := json.Unmarshal(configs.LangExts, &allLangs)
	if err != nil {
		return nil, err
	}

	langToExts := make(map[string][]string)
	for _, l := range allLangs {
		langToExts[strings.ToLower(l.Name)] = l.Extensions
	}

	extensions := make(map[string]struct{})
	for _, l := range langs {
		lc := strings.ToLower(l)
		if exts, ok := langToExts[lc]; ok {
			for _, ext := range exts {
				extensions[ext] = struct{}{}
			}
		}
	}

	return extensions, nil
}

func filterExtensions(files, langs, extensions []string) ([]string, error) {
	if len(langs) == 0 && len(extensions) == 0 {
		return files, nil
	}

	allowedExts, err := extsFromLangs(langs)
	if err != nil {
		return nil, err
	}
	for _, ext := range extensions {
		allowedExts[ext] = struct{}{}
	}

	var filtered []string
	for _, f := range files {
		if _, ok := allowedExts[filepath.Ext(f)]; ok {
			filtered = append(filtered, f)
		}
	}

	return filtered, nil
}

func filterPatterns(files, exclude, restrict []string) ([]string, error) {
	if len(exclude) == 0 && len(restrict) == 0 {
		return files, nil
	}

	excludeSet := sliceToSet(exclude)
	restrictSet := sliceToSet(restrict)
	var filtered []string

Loop:
	for _, file := range files {
		for excl := range excludeSet {
			match, err := filepath.Match(excl, file)
			if err != nil {
				return nil, err
			}
			if match {
				continue Loop
			}
		}

		if len(restrict) == 0 {
			filtered = append(filtered, file)
		}

		for restr := range restrictSet {
			match, err := filepath.Match(restr, file)
			if err != nil {
				return nil, err
			}
			if match {
				filtered = append(filtered, file)
				continue Loop
			}
		}
	}

	return filtered, nil
}

func sliceToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, elem := range slice {
		set[elem] = struct{}{}
	}
	return set
}
