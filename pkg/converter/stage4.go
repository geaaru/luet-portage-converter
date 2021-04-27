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
	"errors"
	"fmt"
	"path/filepath"

	//cfg "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
	luet_tree "github.com/mudler/luet/pkg/tree"
)

type Stage4Worker struct {
	Map    map[string]*luet_pkg.DefaultPackage
	Levels *Stage4Levels
}

func (pc *PortageConverter) Stage4() error {

	InfoC(GetAurora().Bold("Stage4 Starting..."))

	if len(pc.Solutions) == 0 {
		InfoC(GetAurora().Bold("Stage4: No solutions to elaborate. Nothing to do."))
		return nil
	}

	pc.ReciperBuild = luet_tree.NewCompilerRecipe(luet_pkg.NewInMemoryDatabase(false))
	pc.ReciperRuntime = luet_tree.NewInstallerRecipe(luet_pkg.NewInMemoryDatabase(false))

	err := pc.LoadTrees(pc.TreePaths)
	if err != nil {
		return err
	}

	// Create stage4 stuff
	var levels *Stage4Levels = nil
	worker := &Stage4Worker{
		Map: make(map[string]*luet_pkg.DefaultPackage, 0),
	}

	for _, pkg := range pc.Solutions {

		levels = NewStage4LevelsWithSize(1)
		worker.Levels = levels

		pack := pkg.ToPack(true)

		// Check buildtime requires
		InfoC(GetAurora().Bold(fmt.Sprintf("[%s/%s-%s]",
			pack.GetCategory(), pack.GetName(), pack.GetVersion())),
			"Preparing stage4 levels struct...")

		luetPkg := &luet_pkg.DefaultPackage{
			Name:     pack.GetName(),
			Category: pack.GetCategory(),
			Version:  pack.GetVersion(),
		}

		err := pc.stage4AddDeps2Levels(luetPkg, nil, worker, 1)
		if err != nil {
			return err
		}
		// Setup level1 with all packages
		err = pc.stage4AlignLevel1(worker)
		if err != nil {
			return err
		}

		DebugC(fmt.Sprintf(
			"Stage4: Created levels structs of %d trees for %d packages.",
			len(levels.Levels), len(levels.Map)))

		pc.stage4LevelsDumpWrapper(levels, "Starting structure")

		err = levels.Resolve()
		if err != nil {
			return errors.New("Error on resolve stage4 levels: " + err.Error())
		}

		pc.stage4LevelsDumpWrapper(levels, "Resolved structure")
	}

	err = pc.stage4UpdateBuildFiles(worker)
	if err != nil {
		return errors.New("Error on update build.yaml files: " + err.Error())
	}

	InfoC(GetAurora().Bold(
		fmt.Sprintf("Stage4 Completed. Updates: %d.", len(levels.Changed))))
	return nil
}

func (pc *PortageConverter) stage4UpdateBuildFiles(worker *Stage4Worker) error {

	if len(worker.Levels.Changed) == 0 {
		return nil
	}

	for _, pkg := range worker.Levels.Changed {

		ppp, err := pc.ReciperBuild.GetDatabase().FindPackages(pkg)
		if err != nil {
			return errors.New(
				fmt.Sprintf("Error on retrieve data of the package %s/%s: %s",
					pkg.GetCategory(), pkg.GetName(), err,
				))
		}

		buildFile := filepath.Join(ppp[0].GetPath(), "build.yaml")
		// Load Build Template file
		buildPack, err := NewLuetCompilationSpecSanitizedFromFile(buildFile)
		if err != nil {
			return err
		}

		prevReqs := len(buildPack.GetRequires())

		// Prepare requires
		reqs := []*luet_pkg.DefaultPackage{}

		for _, dep := range pkg.GetRequires() {
			reqs = append(reqs, &luet_pkg.DefaultPackage{
				Category: dep.GetCategory(),
				Name:     dep.GetName(),
				Version:  ">=0",
			})
		}

		buildPack.Requires(reqs)

		err = buildPack.WriteBuildDefinition(buildFile)
		if err != nil {
			return err
		}

		InfoC(fmt.Sprintf("[%s/%s-%s] Update requires (%d -> %d).",
			pkg.GetCategory(), pkg.GetName(), pkg.GetVersion(),
			prevReqs, len(reqs)))

	}

	return nil
}

func (pc *PortageConverter) stage4LevelsDumpWrapper(levels *Stage4Levels, msg string) {
	if len(levels.Levels) > 10 {
		InfoC(fmt.Sprintf(
			"Stage4: %s:\n", msg))
		for idx, _ := range levels.Levels {
			InfoC(levels.Levels[idx].Dump())
		}

	} else {
		DebugC(fmt.Sprintf(
			"Stage4: %s:\n%s\n", msg, levels.Dump(),
		))
	}
}

func (pc *PortageConverter) stage4AddDeps2Levels(pkg *luet_pkg.DefaultPackage,
	father *luet_pkg.DefaultPackage,
	w *Stage4Worker, level int) error {

	key := fmt.Sprintf("%s/%s", pkg.GetCategory(), pkg.GetName())

	// Check if level is already available
	if len(w.Levels.Levels) < level {
		tree := NewStage4Tree(level)
		w.Levels.AddTree(tree)
	}

	v, ok := w.Map[key]
	if ok {
		// Package already in map. I will use the same reference.
		w.Levels.AddDependency(v, father, level-1)
		pkg = v

	} else {

		pkg_search := &luet_pkg.DefaultPackage{
			Category: pkg.GetCategory(),
			Name:     pkg.GetName(),
			Version:  ">=0",
		}

		ppp, err := pc.ReciperBuild.GetDatabase().FindPackages(pkg_search)
		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"Error on retrieve dependencies of the package %s/%s: %s",
					pkg.GetCategory(), pkg.GetName(), err.Error()))
		}

		pkg.Requires(ppp[0].GetRequires())

		// Add package to first level
		w.Levels.AddDependency(pkg, father, level-1)
		w.Map[key] = pkg

	}

	if len(pkg.GetRequires()) > 0 {

		// Add requires
		for _, dep := range pkg.GetRequires() {
			err := pc.stage4AddDeps2Levels(dep, pkg, w, level+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pc *PortageConverter) stage4AlignLevel1(w *Stage4Worker) error {

	for pkg, v := range w.Levels.Map {

		if _, ok := w.Levels.Levels[0].Map[pkg]; !ok {

			DebugC(fmt.Sprintf("Adding package %s..", pkg))
			_, err := w.Levels.AddDependencyRecursive(v, nil, []string{}, 0)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
