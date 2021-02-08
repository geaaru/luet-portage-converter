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
package converter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Luet-lab/luet-portage-converter/pkg/qdepends"
	"github.com/Luet-lab/luet-portage-converter/pkg/reposcan"
	"github.com/Luet-lab/luet-portage-converter/pkg/specs"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	luet_config "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
	luet_tree "github.com/mudler/luet/pkg/tree"
)

type PortageConverter struct {
	Config         *luet_config.LuetConfig
	Cache          map[string]*specs.PortageSolution
	ReciperBuild   luet_tree.Builder
	ReciperRuntime luet_tree.Builder
	Specs          *specs.PortageConverterSpecs
	TargetDir      string
	Solutions      []*specs.PortageSolution
	Backend        string
	Resolver       specs.PortageResolver

	Override          bool
	IgnoreMissingDeps bool
	DisabledUseFlags  []string
	TreePaths         []string
}

func NewPortageConverter(targetDir, backend string) *PortageConverter {
	return &PortageConverter{
		// TODO: we use it as singleton
		Config:            luet_config.LuetCfg,
		Cache:             make(map[string]*specs.PortageSolution, 0),
		ReciperBuild:      luet_tree.NewCompilerRecipe(luet_pkg.NewInMemoryDatabase(false)),
		ReciperRuntime:    luet_tree.NewInstallerRecipe(luet_pkg.NewInMemoryDatabase(false)),
		TargetDir:         targetDir,
		Backend:           backend,
		Override:          false,
		IgnoreMissingDeps: false,
		DisabledUseFlags:  []string{},
	}
}

func (pc *PortageConverter) LoadTrees(treePath []string) error {

	// Load trees
	for _, t := range treePath {
		InfoC(fmt.Sprintf(":evergreen_tree: Loading tree %s...", t))
		err := pc.ReciperBuild.Load(t)
		if err != nil {
			return errors.New("Error on load tree" + err.Error())
		}
		err = pc.ReciperRuntime.Load(t)
		if err != nil {
			return errors.New("Error on load tree" + err.Error())
		}
	}

	return nil
}

func (pc *PortageConverter) LoadRules(file string) error {
	spec, err := specs.LoadSpecsFile(file)
	if err != nil {
		return err
	}

	pc.Specs = spec

	if spec.BuildTmplFile == "" {
		return errors.New("No build template file defined")
	}

	return nil
}

func (pc *PortageConverter) GetSpecs() *specs.PortageConverterSpecs {
	return pc.Specs
}

func (pc *PortageConverter) IsDep2Skip(pkg *gentoo.GentooPackage) bool {

	for _, skipPkg := range pc.Specs.SkippedResolutions.Packages {
		if skipPkg.Name == pkg.Name && skipPkg.Category == pkg.Category {
			return true
		}
	}

	for _, cat := range pc.Specs.SkippedResolutions.Categories {
		if cat == pkg.Category {
			return true
		}
	}

	return false
}

func (pc *PortageConverter) IsInStack(stack []string, pkg string) bool {
	ans := false
	for _, p := range stack {
		if p == pkg {
			ans = true
			break
		}
	}
	return ans
}

func (pc *PortageConverter) AppendIfNotPresent(list []gentoo.GentooPackage, pkg gentoo.GentooPackage) []gentoo.GentooPackage {
	ans := list
	isPresent := false
	for _, p := range list {
		if p.Name == pkg.Name && p.Category == pkg.Category {
			isPresent = true
			break
		}
	}
	if !isPresent {
		ans = append(ans, pkg)
	}
	return ans
}

func (pc *PortageConverter) createSolution(pkg, treePath string, stack []string, artefact specs.PortageConverterArtefact) error {

	// Check if it's present artefact from map
	art, err := pc.Specs.GetArtefactByPackage(pkg)
	if err == nil {
		// POST: use artefact from map.
		artefact = *art
	}

	opts := specs.PortageResolverOpts{
		EnableUseFlags:   artefact.Uses.Enabled,
		DisabledUseFlags: artefact.Uses.Disabled,
	}

	if pc.IsInStack(stack, pkg) {
		DebugC(fmt.Sprintf("Intercepted cycle dep for %s: %s", pkg, stack))
		DebugC(fmt.Sprintf("[%s] I skip cycle.", pkg))
		// TODO: Is this correct?
		return nil
	}

	gp, err := gentoo.ParsePackageStr(pkg)
	// Avoid to resolve it if it's skipped. Workaround to qdepends problems.
	if err != nil {
		return err
	}

	if pc.IsDep2Skip(gp) {
		DebugC(fmt.Sprintf("[%s] Skipped dependency %s", stack[len(stack)-1], pkg))
		return nil
	}

	solution, err := pc.Resolver.Resolve(pkg, opts)
	if err != nil {
		return errors.New(fmt.Sprintf("Error on resolve %s: %s", pkg, err.Error()))
	}

	stack = append(stack, pkg)

	cacheKey := fmt.Sprintf("%s/%s",
		specs.SanitizeCategory(solution.Package.Category, solution.Package.Slot),
		solution.Package.Name)

	if _, ok := pc.Cache[cacheKey]; ok {
		DebugC(fmt.Sprintf("Package %s already in cache.", pkg))
		// Nothing to do
		return nil
	}

	InfoC(GetAurora().Bold(fmt.Sprintf(":pizza: Creating solution for %s (%s)...", pkg, treePath)))

	pkgDir := fmt.Sprintf("%s/%s/%s/",
		filepath.Join(pc.TargetDir, treePath),
		solution.Package.Category, solution.Package.Name)

	if solution.Package.Slot != "0" {
		slot := solution.Package.Slot
		// Ignore sub-slot
		if strings.Contains(solution.Package.Slot, "/") {
			slot = solution.Package.Slot[0:strings.Index(slot, "/")]
		}

		pkgDir = fmt.Sprintf("%s/%s-%s/%s",
			filepath.Join(pc.TargetDir, treePath),
			solution.Package.Category, slot, solution.Package.Name)
	}

	// Check if specs is already present. I don't check definition.yaml
	// because with collection packages could be inside collection file.
	pTarget := luet_pkg.NewPackage(solution.Package.Name, ">=0",
		[]*luet_pkg.DefaultPackage{},
		[]*luet_pkg.DefaultPackage{})
	pTarget.Category = specs.SanitizeCategory(solution.Package.Category, solution.Package.Slot)

	p, _ := pc.ReciperRuntime.GetDatabase().FindPackages(pTarget)
	// TODO: at the moment we ignore version. Do We want to handle this with Marvin?
	if p != nil && !pc.Override {
		// Nothing to do
		InfoC(fmt.Sprintf("Package %s already in tree.", pkg))
		return nil
	}

	// TODO: atm I handle build-dep and runtime-dep at the same
	//       way. I'm not sure if this correct.

	// Check every build dependency
	var bdeps []gentoo.GentooPackage = make([]gentoo.GentooPackage, 0)
	for _, bdep := range solution.BuildDeps {

		DebugC(fmt.Sprintf("[%s] Analyzing buildtime dep %s...", pkg, bdep.GetPackageName()))

		if pc.IsDep2Skip(&bdep) {
			DebugC(fmt.Sprintf("[%s] Skipped dependency %s", pkg, bdep.GetPackageName()))
			continue
		}

		dep := luet_pkg.NewPackage(bdep.Name, ">=0",
			[]*luet_pkg.DefaultPackage{},
			[]*luet_pkg.DefaultPackage{})
		dep.Category = specs.SanitizeCategory(bdep.Category, bdep.Slot)

		// Check if it's present the build dep on the tree
		p, _ := pc.ReciperBuild.GetDatabase().FindPackage(dep)
		if p == nil {

			// Check if there is a runtime deps/provide for this
			p, _ := pc.ReciperRuntime.GetDatabase().FindPackage(dep)
			if p == nil {
				dep_str := fmt.Sprintf("%s/%s", bdep.Category, bdep.Name)
				if bdep.Slot != "0" {
					dep_str += ":" + bdep.Slot
				}
				// Now we use the same treePath.
				err := pc.createSolution(dep_str, treePath, stack, artefact)
				if err != nil {
					return err
				}

				bdeps = pc.AppendIfNotPresent(bdeps, bdep)
			} else {

				DebugC(fmt.Sprintf("[%s] For buildtime dep %s is used package %s",
					pkg, bdep.GetPackageName(), p.HumanReadableString()))

				gp := gentoo.GentooPackage{
					Name:     p.GetName(),
					Category: p.GetCategory(),
					Version:  ">=0",
					Slot:     "0",
				}
				bdeps = pc.AppendIfNotPresent(bdeps, gp)
			}
		} else {
			DebugC(fmt.Sprintf("[%s] For build-time dep %s is used package %s",
				pkg, bdep.GetPackageName(), p.HumanReadableString()))

			gp := gentoo.GentooPackage{
				Name:     p.GetName(),
				Category: p.GetCategory(),
				Version:  ">=0",
				Slot:     "0",
			}
			bdeps = pc.AppendIfNotPresent(bdeps, gp)
		}

	}
	solution.BuildDeps = bdeps

	// Check buildtime conflicts
	var bconflicts []gentoo.GentooPackage = make([]gentoo.GentooPackage, 0)
	for _, bconflict := range solution.BuildConflicts {

		DebugC(fmt.Sprintf("[%s] Analyzing buildtime conflict %s...",
			pkg, bconflict.GetPackageName()))

		if pc.IsDep2Skip(&bconflict) {
			DebugC(fmt.Sprintf("[%s] Skipped dependency %s", pkg, bconflict.GetPackageName()))
			continue
		}

		dep := luet_pkg.NewPackage(bconflict.Name, ">=0",
			[]*luet_pkg.DefaultPackage{},
			[]*luet_pkg.DefaultPackage{})
		dep.Category = specs.SanitizeCategory(bconflict.Category, bconflict.Slot)

		// Check if it's present the build dep on the tree
		p, _ := pc.ReciperRuntime.GetDatabase().FindPackage(dep)
		if p == nil {
			DebugC(fmt.Sprintf("[%s] Package in conflict %s not in tree. I ignore it.",
				pkg, p))
			continue
		}

		gp := gentoo.GentooPackage{
			Name:     p.GetName(),
			Category: p.GetCategory(),
			Version:  ">=0",
			Slot:     "0",
		}
		bconflicts = pc.AppendIfNotPresent(bconflicts, gp)
	}
	solution.BuildConflicts = bconflicts

	// Check every runtime deps
	var rdeps []gentoo.GentooPackage = make([]gentoo.GentooPackage, 0)
	for _, rdep := range solution.RuntimeDeps {

		DebugC(fmt.Sprintf("[%s] Analyzing runtime dep %s...", pkg, rdep.GetPackageName()))

		if pc.IsDep2Skip(&rdep) {
			DebugC(fmt.Sprintf("[%s] Skipped dependency %s", pkg, rdep.GetPackageName()))
			continue
		}

		dep := luet_pkg.NewPackage(rdep.Name, ">=0",
			[]*luet_pkg.DefaultPackage{},
			[]*luet_pkg.DefaultPackage{})
		dep.Category = specs.SanitizeCategory(rdep.Category, rdep.Slot)

		// Check if it's present the build dep on the tree
		p, _ := pc.ReciperRuntime.GetDatabase().FindPackages(dep)
		if p == nil {
			dep_str := fmt.Sprintf("%s/%s", rdep.Category, rdep.Name)
			if rdep.Slot != "0" {
				dep_str += ":" + rdep.Slot
			}
			// Now we use the same treePath.
			err := pc.createSolution(dep_str, treePath, stack, artefact)
			if err != nil {
				return err
			}

			rdeps = pc.AppendIfNotPresent(rdeps, rdep)

		} else {
			// TODO: handle package list in a better way
			DebugC(fmt.Sprintf("[%s] For runtime dep %s is used package %s",
				pkg, rdep.GetPackageName(), p[0].HumanReadableString()))

			gp := gentoo.GentooPackage{
				Name:     p[0].GetName(),
				Category: p[0].GetCategory(),
				Version:  ">=0",
				Slot:     "0",
			}
			rdeps = pc.AppendIfNotPresent(rdeps, gp)
		}
	}
	solution.RuntimeDeps = rdeps

	// Check runtime conflicts
	var rconflicts []gentoo.GentooPackage = make([]gentoo.GentooPackage, 0)
	for _, rconflict := range solution.RuntimeConflicts {

		DebugC(fmt.Sprintf("[%s] Analyzing runtime conflict %s...",
			pkg, rconflict.GetPackageName()))

		if pc.IsDep2Skip(&rconflict) {
			DebugC(fmt.Sprintf("[%s] Skipped dependency %s", pkg, rconflict.GetPackageName()))
			continue
		}

		dep := luet_pkg.NewPackage(rconflict.Name, ">=0",
			[]*luet_pkg.DefaultPackage{},
			[]*luet_pkg.DefaultPackage{})
		dep.Category = specs.SanitizeCategory(rconflict.Category, rconflict.Slot)

		// Check if it's present the build dep on the tree
		p, _ := pc.ReciperRuntime.GetDatabase().FindPackage(dep)
		if p == nil {
			DebugC(fmt.Sprintf(
				"[%s] Runtime package dep in conflict %s not in tree. I ignore it.",
				pkg))
			continue
		}

		gp := gentoo.GentooPackage{
			Name:     p.GetName(),
			Category: p.GetCategory(),
			Version:  ">=0",
			Slot:     "0",
		}
		rconflicts = pc.AppendIfNotPresent(rconflicts, gp)
	}
	solution.RuntimeConflicts = rconflicts

	solution.PackageDir = pkgDir

	pc.Cache[cacheKey] = solution

	pc.Solutions = append(pc.Solutions, solution)

	return nil
}

func (pc *PortageConverter) Generate() error {
	// Load Build template file
	buildTmpl, err := NewLuetCompilationSpecSanitizedFromFile(pc.Specs.BuildTmplFile)
	if err != nil {
		return errors.New("Error on load template: " + err.Error())
	}

	// Initialize resolver
	if pc.Backend == "reposcan" {
		if len(pc.Specs.ReposcanSources) == 0 {
			return errors.New("No reposcan sources defined!")
		}

		resolver := reposcan.NewRepoScanResolver()
		resolver.JsonSources = pc.Specs.ReposcanSources
		resolver.SetIgnoreMissingDeps(pc.IgnoreMissingDeps)
		resolver.SetDepsWithSlot(pc.Specs.ReposcanRequiresWithSlot)
		resolver.SetDisabledUseFlags(pc.Specs.ReposcanDisabledUseFlags)
		resolver.SetDisabledKeywords(pc.Specs.ReposcanDisabledKeywords)
		InfoC(fmt.Sprintf("Using dependency with slot on category: %v",
			resolver.GetDepsWithSlot()))
		InfoC(fmt.Sprintf("Disabled keywords: %s",
			GetAurora().Bold(resolver.GetDisabledKeywords())))
		InfoC(fmt.Sprintf("Disabled USE: %s",
			GetAurora().Bold(resolver.GetDisabledUseFlags())))
		err = resolver.LoadJsonFiles()
		if err != nil {
			return err
		}

		if len(pc.Specs.ReposcanConstraints.Packages) > 0 {
			resolver.Constraints = pc.Specs.ReposcanConstraints.Packages
		}

		err = resolver.BuildMap()
		if err != nil {
			return err
		}

		pc.Resolver = resolver
	} else {
		pc.Resolver = qdepends.NewQDependsResolver()
	}

	// Resolve all packages
	for _, artefact := range pc.Specs.GetArtefacts() {
		for _, pkg := range artefact.GetPackages() {

			DebugC(fmt.Sprintf("Analyzing package %s...", pkg))
			err := pc.createSolution(pkg, artefact.GetTree(), []string{}, artefact)
			if err != nil {
				return err
			}
		}
	}

	// Stage1: Write new specs without analyzing requires for build / runtime.
	for _, pkg := range pc.Solutions {

		InfoC(fmt.Sprintf(
			":cake: Processing package %s-%s", pkg.Package.GetPackageName(), pkg.Package.GetPVR()))

		err := os.MkdirAll(pkg.PackageDir, 0755)
		if err != nil {
			return err
		}

		defFile := filepath.Join(pkg.PackageDir, "definition.yaml")
		buildFile := filepath.Join(pkg.PackageDir, "build.yaml")

		// Convert solution to luet package
		pack := pkg.ToPack(true)

		// Write definition.yaml
		err = luet_tree.WriteDefinitionFile(pack, defFile)
		if err != nil {
			return err
		}

		// Check if artefact is in map
		ignoreBuildDeps := false
		artefact, err := pc.Specs.GetArtefactByPackage(pkg.Package.GetPackageNameWithSlot())
		if err != nil {
			if pkg.Package.Slot != "" {
				artefact, err = pc.Specs.GetArtefactByPackage(pkg.Package.GetPackageName())
			}
		}
		if artefact != nil {
			ignoreBuildDeps = artefact.IgnoreBuildDeps
		}

		// create build.yaml
		bPack := pkg.ToPack(false)
		buildPack, _ := buildTmpl.Clone()
		if !ignoreBuildDeps {
			buildPack.AddRequires(bPack.PackageRequires)
		}
		buildPack.AddConflicts(bPack.PackageConflicts)

		err = buildPack.WriteBuildDefinition(buildFile)
		if err != nil {
			return err
		}
	}

	InfoC(GetAurora().Bold(fmt.Sprintf(
		"Stage1 Completed: generated %d packages.", len(pc.Solutions))))

	// Stage2: Reload tree and drop redundant dependencies
	err = pc.Stage2()
	if err != nil {
		return err
	}

	return nil
}

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
		resolvedRuntimeDeps := []*luet_pkg.DefaultPackage{}
		resolvedBuildtimeDeps := []*luet_pkg.DefaultPackage{}

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

							InfoC(fmt.Sprintf("[%s/%s-%s] Dropping buildtime dep %s/%s available in %s/%s",
								pack.GetCategory(), pack.GetName(), pack.GetVersion(),
								dep.GetCategory(), dep.GetName(),
								d2.GetCategory(), d2.GetName(),
							))
							goto next_dep
						}
					}

				}

			next_dep:

				if !alreadyInjected {
					resolvedBuildtimeDeps = append(resolvedBuildtimeDeps, dep)
				}

			} // end for idx, dep

			if len(resolvedRuntimeDeps) != len(deps) {
				updateBuildDeps = true
			}

		} else {

			DebugC(fmt.Sprintf("[%s/%s-%s] Only one buildtime dep present. Nothing to do.",
				pack.GetCategory(), pack.GetName(), pack.GetVersion()))

		}

		// Check runtime requires
		pReciper, err = pc.ReciperRuntime.GetDatabase().FindPackage(luetPkg)
		if err != nil {
			return err
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

							InfoC(fmt.Sprintf("[%s/%s-%s] Dropping runtime dep %s/%s available in %s/%s",
								pack.GetCategory(), pack.GetName(), pack.GetVersion(),
								dep.GetCategory(), dep.GetName(),
								d2.GetCategory(), d2.GetName(),
							))
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

		}

		// Write definition
		if updateRuntimeDeps {

			defFile := filepath.Join(pkg.PackageDir, "definition.yaml")
			// Convert solution to luet package
			pack := pkg.ToPack(true)
			pack.Requires(resolvedRuntimeDeps)

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
			bPack := pkg.ToPack(false)
			buildPack, _ := buildTmpl.Clone()
			buildPack.Requires(resolvedBuildtimeDeps)
			buildPack.AddConflicts(bPack.PackageConflicts)

			err = buildPack.WriteBuildDefinition(buildFile)
			if err != nil {
				return err
			}

		}

	}

	return nil
}
