/*
Copyright (C) 2020  Daniele Rondina <geaaru@sabayonlinux.org>

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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	luet_pkg "github.com/mudler/luet/pkg/package"

	"gopkg.in/yaml.v2"
)

type PortageConverterSpecs struct {
	SkippedResolutions PortageConverterSkips `json:"skipped_resolutions,omitempty" yaml:"skipped_resolutions,omitempty"`

	IncludeFiles  []string                   `json:"include_files,omitempty" yaml:"include_files,omitempty"`
	Artefacts     []PortageConverterArtefact `json:"artefacts,omitempty" yaml:"artefacts,omitempty"`
	BuildTmplFile string                     `json:"build_template_file" yaml:"build_template_file"`

	// Reposcan options
	ReposcanRequiresWithSlot bool                                `json:"reposcan_requires_slot,omitempty" yaml:"reposcan_requires_slot,omitempty"`
	ReposcanSources          []string                            `json:"reposcan_sources,omitempty" yaml:"reposcan_sources,omitempty"`
	ReposcanConstraints      PortageConverterReposcanConstraints `json:"reposcan_contraints,omitempty" yaml:"reposcan_contraints,omitempty"`
	ReposcanDisabledUseFlags []string                            `json:"reposcan_disabled_use_flags,omitempty" yaml:"reposcan_disabled_use_flags,omitempty"`
	ReposcanDisabledKeywords []string                            `json:"reposcan_disabled_keywords,omitempty" yaml:"reposcan_disabled_keywords,omitempty"`
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
	Tree     string                   `json:"tree" yaml:"tree"`
	Uses     PortageConverterUseFlags `json:"uses,omitempty" yaml:"uses,omitempty"`
	Packages []string                 `json:"packages" yaml:"packages"`
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

	return ans, nil
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

func (a *PortageConverterArtefact) GetPackages() []string { return a.Packages }
func (a *PortageConverterArtefact) GetTree() string       { return a.Tree }

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
}

func NewPortageResolverOpts() PortageResolverOpts {
	return PortageResolverOpts{
		EnableUseFlags:   []string{},
		DisabledUseFlags: []string{},
	}
}

func (o *PortageResolverOpts) IsAdmitUseFlag(u string) bool {
	ans := true
	if len(o.EnableUseFlags) > 0 {
		for _, ue := range o.EnableUseFlags {
			if ue == u {
				return true
			}
		}

		return false
	}

	if len(o.DisabledUseFlags) > 0 {
		for _, ud := range o.DisabledUseFlags {
			if ud == u {
				ans = false
				break
			}
		}
	}

	return ans
}

func (s *PortageSolution) SetLabel(k, v string) {
	if v != "" {
		s.Labels[k] = v
	}
}

func (s *PortageSolution) ToPack(runtime bool) *luet_pkg.DefaultPackage {

	version := s.Package.Version
	// TODO: handle particular use cases
	if strings.HasPrefix(s.Package.VersionSuffix, "_pre") {
		version = version + s.Package.VersionSuffix
	}

	emergePackage := s.Package.GetPackageName()
	if s.Package.Slot != "0" {
		emergePackage = emergePackage + ":" + s.Package.Slot
	}

	labels := s.Labels
	labels["original.package.name"] = s.Package.GetPackageName()
	labels["original.package.version"] = s.Package.GetPVR()
	labels["emerge.packages"] = emergePackage
	labels["kit"] = s.Package.Repository

	useFlags := []string{}

	if len(s.Package.UseFlags) > 0 {
		// Avoid duplicated
		m := make(map[string]int, 0)
		for _, u := range s.Package.UseFlags {
			m[u] = 1
		}
		for k, _ := range m {
			useFlags = append(useFlags, k)
		}
	}

	if len(useFlags) == 0 {
		useFlags = nil
	}

	ans := &luet_pkg.DefaultPackage{
		Name:        s.Package.Name,
		Category:    SanitizeCategory(s.Package.Category, s.Package.Slot),
		Version:     version,
		UseFlags:    useFlags,
		Labels:      labels,
		License:     s.Package.License,
		Description: s.Description,
		Uri:         s.Uri,
	}

	deps := s.BuildDeps
	if runtime {
		deps = s.RuntimeDeps
	}

	for _, req := range deps {

		dep := &luet_pkg.DefaultPackage{
			Name:     req.Name,
			Category: SanitizeCategory(req.Category, req.Slot),
			UseFlags: req.UseFlags,
		}
		if req.Version != "" && req.Condition != gentoo.PkgCondNot &&
			req.Condition != gentoo.PkgCondAnyRevision &&
			req.Condition != gentoo.PkgCondMatchVersion &&
			req.Condition != gentoo.PkgCondEqual {

			// TODO: to complete
			dep.Version = fmt.Sprintf("%s%s%s",
				req.Condition.String(), req.Version, req.VersionSuffix)

		} else {
			dep.Version = ">=0"
		}

		ans.PackageRequires = append(ans.PackageRequires, dep)
	}

	if runtime && len(s.RuntimeConflicts) > 0 {

		for _, req := range s.RuntimeConflicts {

			dep := &luet_pkg.DefaultPackage{
				Name:     req.Name,
				Category: SanitizeCategory(req.Category, req.Slot),
				UseFlags: req.UseFlags,
			}
			if req.Version != "" && req.Condition == gentoo.PkgCondNot {
				// TODO: to complete
				dep.Version = fmt.Sprintf("%s%s%s",
					req.Condition.String(), req.Version, req.VersionSuffix)

			} else {
				dep.Version = ">=0"
			}

			ans.PackageConflicts = append(ans.PackageConflicts, dep)
		}

	} else if !runtime && len(s.BuildConflicts) > 0 {

		for _, req := range s.BuildConflicts {

			dep := &luet_pkg.DefaultPackage{
				Name:     req.Name,
				Category: SanitizeCategory(req.Category, req.Slot),
				UseFlags: req.UseFlags,
			}
			if req.Version != "" && req.Condition == gentoo.PkgCondNot {
				// TODO: to complete
				dep.Version = fmt.Sprintf("%s%s%s",
					req.Condition.String(), req.Version, req.VersionSuffix)

			} else {
				dep.Version = ">=0"
			}

			ans.PackageConflicts = append(ans.PackageConflicts, dep)
		}

	}

	return ans
}

func (s *PortageSolution) String() string {
	data, _ := json.Marshal(*s)
	return string(data)
}

func SanitizeCategory(cat string, slot string) string {
	ans := cat
	if slot != "0" {
		// Ignore sub-slot
		if strings.Contains(slot, "/") {
			slot = slot[0:strings.Index(slot, "/")]
		}

		if slot != "0" && slot != "" {
			ans = fmt.Sprintf("%s-%s", cat, slot)
		}
	}
	return ans
}
