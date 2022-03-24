/*
	Copyright © 2021 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package reposcan

import (
	"fmt"
	"strings"

	"github.com/geaaru/luet-portage-converter/pkg/specs"
	gentoo "github.com/geaaru/pkgs-checker/pkg/gentoo"

	"gopkg.in/yaml.v2"
)

type RepoScanSpec struct {
	CacheDataVersion string                  `json:"cache_data_version" yaml:"cache_data_version"`
	Atoms            map[string]RepoScanAtom `json:"atoms" yaml:"atoms"`

	File string `json:"-"`
}

type RepoScanAtom struct {
	Atom string `json:"atom" yaml:"atom"`

	Category string     `json:"category" yaml:"category"`
	Package  string     `json:"package" yaml:"package"`
	Revision string     `json:"revision" yaml:"revision"`
	CatPkg   string     `json:"catpkg" yaml:"catpkg"`
	Eclasses [][]string `json:"eclasses" yaml:"eclasses"`

	Kit    string `json:"kit" yaml:"kit"`
	Branch string `json:"branch" yaml:"branch"`

	// Relations contains the list of the keys defined on
	// relations_by_kind. The values could be RDEPEND, DEPEND, BDEPEND
	Relations       []string            `json:"relations" yaml:"relations"`
	RelationsByKind map[string][]string `json:"relations_by_kind" yaml:"relations_by_kind"`

	// Metadata contains ebuild variables.
	// Ex: SLOT, SRC_URI, HOMEPAGE, etc.
	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
	MetadataOut string            `json:"metadata_out" yaml:"metadata_out"`

	ManifestMd5 string `json:"manifest_md5" yaml:"manifest_md5"`
	Md5         string `json:"md5" yaml:"md5"`

	// Fields present on failure
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
	Output string `json:"output,omitempty" yaml:"output,omitempty"`

	Files []RepoScanFile `json:"files" yaml:"files"`
}

type RepoScanFile struct {
	SrcUri []string          `json:"src_uri"`
	Size   string            `json:"size"`
	Hashes map[string]string `json:"hashes"`
	Name   string            `json:"name"`
}

func (r *RepoScanSpec) Yaml() (string, error) {
	data, err := yaml.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *RepoScanAtom) GetPackageName() string {
	return fmt.Sprintf("%s/%s", r.GetCategory(), r.Package)
}

func (r *RepoScanAtom) GetCategory() string {
	slot := "0"

	if r.HasMetadataKey("SLOT") {
		slot = r.GetMetadataValue("SLOT")
		// We ignore subslot atm.
		if strings.Contains(slot, "/") {
			slot = slot[0:strings.Index(slot, "/")]
		}

	}

	return specs.SanitizeCategory(r.Category, slot)
}

func (r *RepoScanAtom) HasMetadataKey(k string) bool {
	_, ans := r.Metadata[k]
	return ans
}

func (r *RepoScanAtom) GetMetadataValue(k string) string {
	ans, _ := r.Metadata[k]
	return ans
}

func (r *RepoScanAtom) ToGentooPackage() (*gentoo.GentooPackage, error) {
	ans, err := gentoo.ParsePackageStr(r.Atom)
	if err != nil {
		return nil, err
	}

	// Retrieve license
	if l, ok := r.Metadata["LICENSE"]; ok {
		ans.License = l
	}

	if slot, ok := r.Metadata["SLOT"]; ok {
		// TOSEE: We ignore subslot atm.
		if strings.Contains(slot, "/") {
			slot = slot[0:strings.Index(slot, "/")]
		}
		ans.Slot = slot
	}

	ans.Repository = r.Kit

	return ans, nil
}

func (r *RepoScanAtom) GetRuntimeDeps() ([]gentoo.GentooPackage, error) {
	ans := []gentoo.GentooPackage{}

	if len(r.Relations) > 0 {
		if _, ok := r.RelationsByKind["RDEPEND"]; ok {

			deps, err := r.getDepends("RDEPEND")
			if err != nil {
				return ans, err
			}
			ans = append(ans, deps...)
		}
		// TODO: Check if it's needed add PDEPEND here
	}

	return ans, nil
}

func (r *RepoScanAtom) GetBuildtimeDeps() ([]gentoo.GentooPackage, error) {
	ans := []gentoo.GentooPackage{}

	if len(r.Relations) > 0 {
		if _, ok := r.RelationsByKind["DEPEND"]; ok {
			deps, err := r.getDepends("DEPEND")
			if err != nil {
				return ans, err
			}
			ans = append(ans, deps...)
		}

		if _, ok := r.RelationsByKind["BDEPEND"]; ok {
			deps, err := r.getDepends("BDEPEND")
			if err != nil {
				return ans, err
			}
			ans = append(ans, deps...)
		}
	}

	return ans, nil
}

func (r *RepoScanAtom) getDepends(depType string) ([]gentoo.GentooPackage, error) {
	ans := []gentoo.GentooPackage{}
	if _, ok := r.RelationsByKind[depType]; ok {

		for _, pkg := range r.RelationsByKind[depType] {
			gp, err := gentoo.ParsePackageStr(pkg)
			if err != nil {
				return ans, err
			}
			gp.Slot = ""
			ans = append(ans, *gp)
		}
	}

	return ans, nil
}
