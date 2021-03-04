/*
Copyright (C) 2020-2021  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/
package specs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"

	"gopkg.in/yaml.v2"
)

type PortageConverterSpecs struct {
	SkippedResolutions PortageConverterSkips `json:"skipped_resolutions,omitempty" yaml:"skipped_resolutions,omitempty"`

	IncludeFiles         []string                   `json:"include_files,omitempty" yaml:"include_files,omitempty"`
	Artefacts            []PortageConverterArtefact `json:"artefacts,omitempty" yaml:"artefacts,omitempty"`
	BuildTmplFile        string                     `json:"build_template_file" yaml:"build_template_file"`
	BuildPortageTmplFile string                     `json:"build_portage_template_file,omitempty" yaml:"build_portage_template_file,omitempty"`

	// Reposcan options
	ReposcanRequiresWithSlot bool                                `json:"reposcan_requires_slot,omitempty" yaml:"reposcan_requires_slot,omitempty"`
	ReposcanSources          []string                            `json:"reposcan_sources,omitempty" yaml:"reposcan_sources,omitempty"`
	ReposcanConstraints      PortageConverterReposcanConstraints `json:"reposcan_contraints,omitempty" yaml:"reposcan_contraints,omitempty"`
	ReposcanDisabledUseFlags []string                            `json:"reposcan_disabled_use_flags,omitempty" yaml:"reposcan_disabled_use_flags,omitempty"`
	ReposcanDisabledKeywords []string                            `json:"reposcan_disabled_keywords,omitempty" yaml:"reposcan_disabled_keywords,omitempty"`

	Replacements PortageConverterReplacements `json:"replacements,omitempty" yaml:"replacements,omitempty"`

	MapArtefacts             map[string]*PortageConverterArtefact       `json:"-" yaml:"-"`
	MapReplacementsRuntime   map[string]*PortageConverterReplacePackage `json:"-" yaml:"-"`
	MapReplacementsBuildtime map[string]*PortageConverterReplacePackage `json:"-" yaml:"-"`
}

type PortageConverterReposcanConstraints struct {
	Packages []string `json:"packages,omitempty" yaml:"packages,omitempty"`
}

type PortageConverterSkips struct {
	Packages   []PortageConverterPkg `json:"packages,omitempty" yaml:"packages,omitempty"`
	Categories []string              `json:"categories,omitempty" yaml:"categories,omitempty"`
}

type PortageConverterPkg struct {
	Name     string `json:"name" yaml:"name"`
	Category string `json:"category" yaml:"category"`
}

type PortageConverterArtefact struct {
	Tree            string                   `json:"tree" yaml:"tree"`
	Uses            PortageConverterUseFlags `json:"uses,omitempty" yaml:"uses,omitempty"`
	IgnoreBuildDeps bool                     `json:"ignore_build_deps,omitempty" yaml:"ignore_build_deps,omitempty"`
	Packages        []string                 `json:"packages" yaml:"packages"`
	OverrideVersion string                   `json:"override_version,omitempty" yaml:"override_version,omitempty"`

	Replacements PortageConverterReplacements `json:"replacements,omitempty" yaml:"replacements,omitempty"`
	Mutations    PortageConverterMutations    `json:"mutations,omitempty" yaml:"mutations,omitempty"`

	MapReplacementsRuntime   map[string]*PortageConverterReplacePackage `json:"-" yaml:"-"`
	MapReplacementsBuildtime map[string]*PortageConverterReplacePackage `json:"-" yaml:"-"`
	MapIgnoreRuntime         map[string]bool                            `json:"-" yaml:"-"`
	MapIgnoreBuildtime       map[string]bool                            `json:"-" yaml:"-"`
}

type PortageConverterMutations struct {
	RuntimeDeps   PortageConverterMutationDeps `json:"runtime_deps,omitempty" yaml:"runtime_deps,omitempty"`
	BuildTimeDeps PortageConverterMutationDeps `json:"buildtime_deps,omitempty" yaml:"buildtime_deps,omitempty"`
	Uses          []string                     `json:"uses,omitempty" yaml:"uses,omitempty"`
}

type PortageConverterMutationDeps struct {
	Packages []PortageConverterPkg `json:"packages,omitempty" yaml:"packages,omitempty"`
}

type PortageConverterReplacements struct {
	RuntimeDeps  PortageConverterDepReplacements `json:"runtime_deps,omitempty" yaml:"runtime_deps,omitempty"`
	BuiltimeDeps PortageConverterDepReplacements `json:"buildtime_deps,omitempty" yaml:"buildtime_deps,omitempty"`
}

type PortageConverterDepReplacements struct {
	Packages []PortageConverterReplacePackage `json:"packages,omitempty" yaml:"packages,omitempty"`
	Ignore   []PortageConverterPkg            `json:"ignore,omitempty" yaml:"ignore,omitempty"`
}

type PortageConverterReplacePackage struct {
	From PortageConverterPkg `json:"from" yaml:"from"`
	To   PortageConverterPkg `json:"to" yaml:"to"`
}

type PortageConverterUseFlags struct {
	Disabled []string `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Enabled  []string `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type PortageConverterInclude struct {
	SkippedResolutions PortageConverterSkips      `json:"skipped_resolutions,omitempty" yaml:"skipped_resolutions,omitempty"`
	Artefacts          []PortageConverterArtefact `json:"artefacts,omitempty" yaml:"artefacts,omitempty"`
}

func SpecsFromYaml(data []byte) (*PortageConverterSpecs, error) {
	ans := &PortageConverterSpecs{}
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}
	return ans, nil
}

func IncludeFromYaml(data []byte) (*PortageConverterInclude, error) {
	ans := &PortageConverterInclude{}
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}
	return ans, nil
}

func LoadSpecsFile(file string) (*PortageConverterSpecs, error) {

	if file == "" {
		return nil, errors.New("Invalid file path")
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	ans, err := SpecsFromYaml(content)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(path.Dir(file))
	if err != nil {
		return nil, err
	}

	if len(ans.IncludeFiles) > 0 {

		for _, include := range ans.IncludeFiles {

			if include[0:1] != "/" {
				include = filepath.Join(absPath, include)
			}
			content, err := ioutil.ReadFile(include)
			if err != nil {
				return nil, err
			}

			data, err := IncludeFromYaml(content)
			if err != nil {
				return nil, err
			}

			if len(data.SkippedResolutions.Packages) > 0 {
				ans.SkippedResolutions.Packages = append(ans.SkippedResolutions.Packages,
					data.SkippedResolutions.Packages...)
			}

			if len(data.SkippedResolutions.Categories) > 0 {
				ans.SkippedResolutions.Categories = append(ans.SkippedResolutions.Categories,
					data.SkippedResolutions.Categories...)
			}

			if len(data.Artefacts) > 0 {
				ans.Artefacts = append(ans.Artefacts, data.Artefacts...)
			}

		}
	}

	if ans.BuildTmplFile != "" && ans.BuildTmplFile[0:1] != "/" {
		// Convert in abs path
		ans.BuildTmplFile = filepath.Join(absPath, ans.BuildTmplFile)
	}

	if ans.BuildPortageTmplFile != "" && ans.BuildPortageTmplFile[0:1] != "/" {
		// Convert in abs path
		ans.BuildPortageTmplFile = filepath.Join(absPath, ans.BuildPortageTmplFile)
	}
	return ans, nil
}

func (s *PortageConverterSpecs) GenerateArtefactsMap() {

	s.MapArtefacts = make(map[string]*PortageConverterArtefact, 0)

	for idx, _ := range s.Artefacts {
		for _, pkg := range s.Artefacts[idx].Packages {
			s.MapArtefacts[pkg] = &s.Artefacts[idx]
		}
	}
}

func (s *PortageConverterSpecs) HasRuntimeReplacement(pkg string) bool {
	_, ans := s.MapReplacementsRuntime[pkg]
	return ans
}

func (s *PortageConverterSpecs) HasBuildtimeReplacement(pkg string) bool {
	_, ans := s.MapReplacementsBuildtime[pkg]
	return ans
}

func (s *PortageConverterSpecs) GetBuildtimeReplacement(pkg string) (*PortageConverterReplacePackage, error) {
	ans, ok := s.MapReplacementsBuildtime[pkg]
	if ok {
		return ans, nil
	}
	return ans, errors.New("No replacement found for key " + pkg)
}

func (s *PortageConverterSpecs) GetRuntimeReplacement(pkg string) (*PortageConverterReplacePackage, error) {
	ans, ok := s.MapReplacementsRuntime[pkg]
	if ok {
		return ans, nil
	}
	return ans, errors.New("No replacement found for key " + pkg)
}

func (s *PortageConverterSpecs) GenerateReplacementsMap() {
	s.MapReplacementsRuntime = make(map[string]*PortageConverterReplacePackage, 0)
	s.MapReplacementsBuildtime = make(map[string]*PortageConverterReplacePackage, 0)

	// Add key from global map
	if len(s.Replacements.RuntimeDeps.Packages) > 0 {
		for ridx, r := range s.Replacements.RuntimeDeps.Packages {
			s.MapReplacementsRuntime[fmt.Sprintf("%s/%s", r.From.Category, r.From.Name)] =
				&s.Replacements.RuntimeDeps.Packages[ridx]
		}
	}

	if len(s.Replacements.BuiltimeDeps.Packages) > 0 {
		for ridx, r := range s.Replacements.BuiltimeDeps.Packages {
			s.MapReplacementsBuildtime[fmt.Sprintf("%s/%s", r.From.Category, r.From.Name)] =
				&s.Replacements.BuiltimeDeps.Packages[ridx]
		}
	}

	// Create artefact maps
	for idx, _ := range s.Artefacts {
		s.Artefacts[idx].MapReplacementsBuildtime = make(map[string]*PortageConverterReplacePackage, 0)
		s.Artefacts[idx].MapReplacementsRuntime = make(map[string]*PortageConverterReplacePackage, 0)
		s.Artefacts[idx].MapIgnoreBuildtime = make(map[string]bool, 0)
		s.Artefacts[idx].MapIgnoreRuntime = make(map[string]bool, 0)

		if len(s.Artefacts[idx].Replacements.RuntimeDeps.Packages) > 0 {
			for ridx, r := range s.Artefacts[idx].Replacements.RuntimeDeps.Packages {
				s.Artefacts[idx].MapReplacementsRuntime[fmt.Sprintf(
					"%s/%s", r.From.Category, r.From.Name)] =
					&s.Artefacts[idx].Replacements.RuntimeDeps.Packages[ridx]
			}
		}

		if len(s.Artefacts[idx].Replacements.RuntimeDeps.Ignore) > 0 {
			for _, r := range s.Artefacts[idx].Replacements.RuntimeDeps.Ignore {
				s.Artefacts[idx].MapIgnoreRuntime[fmt.Sprintf(
					"%s/%s", r.Category, r.Name)] = true
			}
		}

		if len(s.Artefacts[idx].Replacements.BuiltimeDeps.Packages) > 0 {
			for ridx, r := range s.Artefacts[idx].Replacements.BuiltimeDeps.Packages {
				s.Artefacts[idx].MapReplacementsBuildtime[fmt.Sprintf(
					"%s/%s", r.From.Category, r.From.Name)] =
					&s.Artefacts[idx].Replacements.BuiltimeDeps.Packages[ridx]
			}
		}

		if len(s.Artefacts[idx].Replacements.BuiltimeDeps.Ignore) > 0 {
			for _, r := range s.Artefacts[idx].Replacements.BuiltimeDeps.Ignore {
				s.Artefacts[idx].MapIgnoreBuildtime[fmt.Sprintf(
					"%s/%s", r.Category, r.Name)] = true
			}
		}

	}
}

func (s *PortageConverterSpecs) GetArtefactByPackage(pkg string) (*PortageConverterArtefact, error) {
	if a, ok := s.MapArtefacts[pkg]; ok {
		return a, nil
	}
	return nil, errors.New("Package not found")
}

func (s *PortageConverterSpecs) GetArtefacts() []PortageConverterArtefact {
	return s.Artefacts
}

func (s *PortageConverterSpecs) AddReposcanSource(source string) {
	s.ReposcanSources = append(s.ReposcanSources, source)
}

func (s *PortageConverterSpecs) AddReposcanDisabledUseFlags(uses []string) {
	s.ReposcanDisabledUseFlags = append(s.ReposcanDisabledUseFlags, uses...)
}

type PortageResolver interface {
	Resolve(pkg string, opts PortageResolverOpts) (*PortageSolution, error)
}

type PortageResolverOpts struct {
	EnableUseFlags   []string
	DisabledUseFlags []string
}

type PortageSolution struct {
	Package          gentoo.GentooPackage   `json:"package"`
	PackageDir       string                 `json:"package_dir"`
	BuildDeps        []gentoo.GentooPackage `json:"build-deps,omitempty"`
	RuntimeDeps      []gentoo.GentooPackage `json:"runtime-deps,omitempty"`
	RuntimeConflicts []gentoo.GentooPackage `json:"runtime-conflicts,omitempty"`
	BuildConflicts   []gentoo.GentooPackage `json:"build-conflicts,omitempty"`

	Description string            `json:"description,omitempty"`
	Uri         []string          `json:"uri,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`

	OverrideVersion string `json:"override_version,omitempty"`
}
