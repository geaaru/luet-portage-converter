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
package converter

import (
	"fmt"
	"path/filepath"

	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
	luet_tree "github.com/mudler/luet/pkg/tree"
)

func (pc *PortageConverter) Stage2() error {

	InfoC(GetAurora().Bold("Stage2 Starting..."))

	// Reset reciper
	pc.ReciperBuild = luet_tree.NewCompilerRecipe(luet_pkg.NewInMemoryDatabase(false))
	pc.ReciperRuntime = luet_tree.NewInstallerRecipe(luet_pkg.NewInMemoryDatabase(false))

	err := pc.LoadTrees(pc.TreePaths)
	if err != nil {
		return err
	}

	for _, pkg := range pc.Solutions {

		pack := pkg.ToPack(true)
		updateBuildDeps := false
		updateRuntimeDeps := false
		runtimeDepsRemoved := 0
		runtimeConflictsRemoved := 0
		buildtimeDepsRemoved := 0
		buildtimeConflictsRemoved := 0
		resolvedRuntimeDeps := []*luet_pkg.DefaultPackage{}
		resolvedBuildtimeDeps := []*luet_pkg.DefaultPackage{}
		resolvedRuntimeConflicts := []*luet_pkg.DefaultPackage{}
		resolvedBuildConflicts := []*luet_pkg.DefaultPackage{}

		// Check buildtime requires
		DebugC(GetAurora().Bold(fmt.Sprintf("[%s/%s-%s]",
			pack.GetCategory(), pack.GetName(), pack.GetVersion())),
			"Checking buildtime dependencies...")

		luetPkg := &luet_pkg.DefaultPackage{
			Name:     pack.GetName(),
			Category: pack.GetCategory(),
			Version:  pack.GetVersion(),
		}
		pReciper, err := pc.ReciperBuild.GetDatabase().FindPackage(luetPkg)
		if err != nil {
			return err
		}

		// Drop conflicts not present on tree
		conflicts := pReciper.GetConflicts()
		if len(conflicts) > 0 {

			for _, dep := range conflicts {
				pp, _ := pc.ReciperBuild.GetDatabase().FindPackages(
					&luet_pkg.DefaultPackage{
						Name:     dep.GetName(),
						Category: dep.GetCategory(),
						Version:  ">=0",
					},
				)

				if pp == nil || len(pp) == 0 {
					pp, _ := pc.ReciperRuntime.GetDatabase().FindPackages(
						&luet_pkg.DefaultPackage{
							Name:     dep.GetName(),
							Category: dep.GetCategory(),
							Version:  ">=0",
						},
					)
					if pp == nil || len(pp) == 0 {
						InfoC(fmt.Sprintf("[%s/%s-%s] Dropping buildtime conflict %s/%s not available in tree.",
							pack.GetCategory(), pack.GetName(), pack.GetVersion(),
							dep.GetCategory(), dep.GetName(),
						))
						buildtimeConflictsRemoved++
					} else {
						resolvedBuildConflicts = append(resolvedBuildConflicts, dep)
					}
				}
			}

			if len(resolvedBuildConflicts) != len(conflicts) {
				updateBuildDeps = true
			}

		}

		deps := pReciper.GetRequires()
		if len(deps) > 1 {

			for idx, dep := range deps {
				alreadyInjected := false

				for idx2, d2 := range deps {
					if idx2 == idx {
						continue
					}

					d2pkgs, err := pc.ReciperBuild.GetDatabase().FindPackages(
						&luet_pkg.DefaultPackage{
							Name:     d2.GetName(),
							Category: d2.GetCategory(),
							Version:  ">=0",
						},
					)
					if err != nil {
						return err
					}

					for _, d3 := range d2pkgs[0].GetRequires() {
						if d3.GetName() == dep.GetName() && d3.GetCategory() == dep.GetCategory() {
							alreadyInjected = true

							DebugC(fmt.Sprintf("[%s/%s-%s] Dropping buildtime dep %s/%s available in %s/%s",
								pack.GetCategory(), pack.GetName(), pack.GetVersion(),
								dep.GetCategory(), dep.GetName(),
								d2.GetCategory(), d2.GetName(),
							))
							buildtimeDepsRemoved++
							goto next_dep
						}
					}

				}

			next_dep:

				if !alreadyInjected {
					resolvedBuildtimeDeps = append(resolvedBuildtimeDeps, dep)
				}

			} // end for idx, dep

			if len(resolvedBuildtimeDeps) != len(deps) {
				updateBuildDeps = true
			}

		} else {

			DebugC(fmt.Sprintf("[%s/%s-%s] Only one buildtime dep present. Nothing to do.",
				pack.GetCategory(), pack.GetName(), pack.GetVersion()))

			resolvedBuildtimeDeps = deps
		}

		// Check runtime requires
		pReciper, err = pc.ReciperRuntime.GetDatabase().FindPackage(luetPkg)
		if err != nil {
			return err
		}

		// Drop conflicts not present on tree
		conflicts = pReciper.GetConflicts()
		if len(conflicts) > 0 {

			for _, dep := range conflicts {
				pp, _ := pc.ReciperRuntime.GetDatabase().FindPackages(
					&luet_pkg.DefaultPackage{
						Name:     dep.GetName(),
						Category: dep.GetCategory(),
						Version:  ">=0",
					},
				)
				if pp == nil || len(pp) == 0 {

					InfoC(fmt.Sprintf("[%s/%s-%s] Dropping runtime conflict %s/%s not available in tree.",
						pack.GetCategory(), pack.GetName(), pack.GetVersion(),
						dep.GetCategory(), dep.GetName(),
					))
					runtimeConflictsRemoved++
				} else {
					resolvedRuntimeConflicts = append(resolvedRuntimeConflicts, dep)
				}
			}

			if len(resolvedRuntimeConflicts) != len(conflicts) {
				updateRuntimeDeps = true
			}

		}

		deps = pReciper.GetRequires()
		if len(deps) > 1 {

			for idx, dep := range deps {
				alreadyInjected := false

				for idx2, d2 := range deps {
					if idx2 == idx {
						continue
					}

					d2pkgs, err := pc.ReciperRuntime.GetDatabase().FindPackages(
						&luet_pkg.DefaultPackage{
							Name:     d2.GetName(),
							Category: d2.GetCategory(),
							Version:  ">=0",
						},
					)
					if err != nil {
						return err
					}

					for _, d3 := range d2pkgs[0].GetRequires() {
						if d3.GetName() == dep.GetName() && d3.GetCategory() == dep.GetCategory() {
							alreadyInjected = true

							DebugC(fmt.Sprintf("[%s/%s-%s] Dropping runtime dep %s/%s available in %s/%s",
								pack.GetCategory(), pack.GetName(), pack.GetVersion(),
								dep.GetCategory(), dep.GetName(),
								d2.GetCategory(), d2.GetName(),
							))
							runtimeDepsRemoved++
							goto next_rdep
						}
					}

				}

			next_rdep:

				if !alreadyInjected {
					resolvedRuntimeDeps = append(resolvedRuntimeDeps, dep)
				}

			} // end for idx, dep

			if len(resolvedRuntimeDeps) != len(deps) {
				updateRuntimeDeps = true
			}

		} else {

			DebugC(fmt.Sprintf("[%s/%s-%s] Only one runtime dep present. Nothing to do.",
				pack.GetCategory(), pack.GetName(), pack.GetVersion()))

			resolvedRuntimeDeps = deps

		}

		// Write definition
		if updateRuntimeDeps {

			defFile := filepath.Join(pkg.PackageDir, "definition.yaml")
			// Convert solution to luet package
			pack := pkg.ToPack(true)
			pack.Requires(resolvedRuntimeDeps)
			pack.Conflicts(resolvedRuntimeConflicts)

			// Write definition.yaml
			err = luet_tree.WriteDefinitionFile(pack, defFile)
			if err != nil {
				return err
			}
		}

		// Write build.yaml
		if updateBuildDeps {

			buildFile := filepath.Join(pkg.PackageDir, "build.yaml")
			// Load Build template file
			buildTmpl, err := NewLuetCompilationSpecSanitizedFromFile(pc.Specs.BuildTmplFile)
			if err != nil {
				return err
			}

			// create build.yaml
			buildPack, _ := buildTmpl.Clone()
			buildPack.Requires(resolvedBuildtimeDeps)
			buildPack.Conflicts(resolvedBuildConflicts)

			err = buildPack.WriteBuildDefinition(buildFile)
			if err != nil {
				return err
			}

		}

		if updateBuildDeps || updateRuntimeDeps {
			InfoC(GetAurora().Bold(
				fmt.Sprintf(
					":angel: [%s/%s-%s] removed: r.deps %d, r.conflicts %d, b.deps %d, b.conflicts %d.",
					pack.GetCategory(), pack.GetName(), pack.GetVersion(),
					runtimeDepsRemoved, runtimeConflictsRemoved,
					buildtimeDepsRemoved, buildtimeConflictsRemoved,
				)))
		}

	}

	return nil
}
