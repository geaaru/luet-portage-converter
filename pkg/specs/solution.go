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
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	luet_pkg "github.com/mudler/luet/pkg/package"
)

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

		sort.Strings(useFlags)
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

			// Skip itself. Maybe we need handle this case in a better way.
			if dep.Name == s.Package.Name && dep.Category == SanitizeCategory(s.Package.Category, s.Package.Slot) {
				continue
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
