/*
Copyright © 2021-2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package reposcan

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	helpers "github.com/MottainaiCI/lxd-compose/pkg/helpers"
)

type ManifestFile struct {
	Md5   string         `json:"manifest_md5,omitempty" yaml:"manifest_md5,omitempty"`
	Files []RepoScanFile `json:"files,omitempty" yaml:"files,omitempty"`
}

func (m *ManifestFile) GetFiles(srcUri string) ([]RepoScanFile, error) {
	ans := []RepoScanFile{}

	srcUri = strings.TrimSpace(srcUri)

	if srcUri == "" {
		// POST: no tarballs and/or files defined.
		return ans, nil
	}

	words := strings.Split(srcUri, " ")

	toParse := len(words)
	idx := 0
	originUri := words[idx]
	for toParse > 0 {

		if words[idx] == "->" {
			idx++
			toParse--

			// Avoid to add two time the same file when the origin is equal
			// to alias
			if words[idx] == filepath.Base(originUri) {
				idx++
				toParse--
				continue
			}
		} else {
			originUri = words[idx]
		}

		baseName := filepath.Base(words[idx])
		// Check if the file is defined in the manifest
		for _, f := range m.Files {
			if f.Name == baseName {
				f.SrcUri = []string{originUri}
				ans = append(ans, f)
				break
			}
		}

		idx++
		toParse--
	}

	return ans, nil
}

func ParseManifest(f string) (*ManifestFile, error) {
	ans := &ManifestFile{
		Files: []RepoScanFile{},
	}

	if helpers.Exists(f) {
		content, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		ans.Md5 = fmt.Sprintf("%x", md5.Sum(content))

		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			words := strings.Split(line, " ")
			if len(words) <= 3 || words[0] != "DIST" {
				continue
			}

			// The src_uri is populate later on processing metadata.
			file := &RepoScanFile{
				Size:   words[2],
				Name:   words[1],
				Hashes: make(map[string]string, 0),
			}
			pos := 3
			for pos < len(words) {
				file.Hashes[strings.ToLower(words[pos])] = words[pos+1]
				pos += 2
			}

			ans.Files = append(ans.Files, *file)
		}
	}

	return ans, nil
}
