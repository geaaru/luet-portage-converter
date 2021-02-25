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
)

func (a *PortageConverterArtefact) GetPackages() []string { return a.Packages }
func (a *PortageConverterArtefact) GetTree() string       { return a.Tree }

func (a *PortageConverterArtefact) HasOverrideVersion(pkg string) bool {
	ans := false
	hasPkg := false

	for _, p := range a.Packages {
		if p == pkg {
			hasPkg = true
			break
		}
	}

	if hasPkg && a.OverrideVersion != "" {
		ans = true
	}

	return ans
}

func (a *PortageConverterArtefact) GetOverrideVersion() string {
	return a.OverrideVersion
}

func (a *PortageConverterArtefact) IgnoreBuildtime(pkg string) bool {
	_, ans := a.MapIgnoreBuildtime[pkg]
	return ans
}

func (a *PortageConverterArtefact) IgnoreRuntime(pkg string) bool {
	_, ans := a.MapIgnoreRuntime[pkg]
	return ans
}

func (s *PortageConverterArtefact) HasRuntimeReplacement(pkg string) bool {
	_, ans := s.MapReplacementsRuntime[pkg]
	return ans
}

func (s *PortageConverterArtefact) HasBuildtimeReplacement(pkg string) bool {
	_, ans := s.MapReplacementsBuildtime[pkg]
	return ans
}

func (s *PortageConverterArtefact) GetBuildtimeReplacement(pkg string) (*PortageConverterReplacePackage, error) {
	ans, ok := s.MapReplacementsBuildtime[pkg]
	if ok {
		return ans, nil
	}
	return ans, errors.New("No replacement found for key " + pkg)
}

func (s *PortageConverterArtefact) GetRuntimeReplacement(pkg string) (*PortageConverterReplacePackage, error) {
	ans, ok := s.MapReplacementsRuntime[pkg]
	if ok {
		return ans, nil
	}
	return ans, errors.New("No replacement found for key " + pkg)
}
