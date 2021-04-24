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

	//cfg "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
)

type Stage4Leaf struct {
	Package  *luet_pkg.DefaultPackage
	Father   []*luet_pkg.DefaultPackage
	Position int
	Counter  int
}

type Stage4Tree struct {
	Id   int
	Map  map[string]*Stage4Leaf
	Deps []*luet_pkg.DefaultPackage
}

type Stage4Levels struct {
	Levels  []*Stage4Tree
	Map     map[string]*luet_pkg.DefaultPackage
	Changed map[string]*luet_pkg.DefaultPackage
}

func NewStage4Levels() *Stage4Levels {
	return &Stage4Levels{
		Levels:  []*Stage4Tree{},
		Map:     make(map[string]*luet_pkg.DefaultPackage, 0),
		Changed: make(map[string]*luet_pkg.DefaultPackage, 0),
	}
}

func NewStage4LevelsWithSize(nLevels int) *Stage4Levels {
	ans := NewStage4Levels()
	for i := 0; i < nLevels; i++ {
		tree := NewStage4Tree(i + 1)
		ans.AddTree(tree)
	}

	return ans
}

func (l *Stage4Levels) GetMap() *map[string]*luet_pkg.DefaultPackage { return &l.Map }

func (l *Stage4Levels) Dump() string {
	ans := ""
	for _, t := range l.Levels {
		ans += t.Dump()
	}

	return ans
}

func (l *Stage4Levels) AddTree(t *Stage4Tree) {
	l.Levels = append(l.Levels, t)
}

func (l *Stage4Levels) AddDependency(p, father *luet_pkg.DefaultPackage, level int) error {
	if level >= len(l.Levels) {
		return errors.New("Invalid level")
	}

	key := fmt.Sprintf("%s/%s", p.GetCategory(), p.GetName())

	v, ok := l.Map[key]
	if ok {
		l.Levels[level].AddDependency(v, father)
	} else {
		l.Map[key] = p
		l.Levels[level].AddDependency(p, father)
	}

	return nil
}

func (l *Stage4Levels) AddDependencyRecursive(p, father *luet_pkg.DefaultPackage, level int) error {
	if level >= len(l.Levels) {
		return errors.New("Invalid level")
	}

	key := fmt.Sprintf("%s/%s", p.GetCategory(), p.GetName())

	DebugC(fmt.Sprintf("Adding recursive %s package to level %d.", key, level))

	_, ok := l.Map[key]
	if !ok {
		return errors.New(fmt.Sprintf("On add dependency not found package %s/%s",
			p.GetCategory(), p.GetName()))
	}

	l.Levels[level].AddDependency(p, father)

	if len(p.GetRequires()) > 0 {

		for _, d := range p.GetRequires() {

			key = fmt.Sprintf("%s/%s", d.GetCategory(), d.GetName())
			v, ok := l.Map[key]
			if !ok {
				return errors.New(fmt.Sprintf("For package %s/%s not found dependency %s",
					p.GetCategory(), p.GetName(), key))
			}

			err := l.AddDependencyRecursive(v, p, level+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Stage4Levels) AddChangedPackage(pkg *luet_pkg.DefaultPackage) {
	key := fmt.Sprintf("%s/%s", pkg.GetCategory(), pkg.GetName())
	l.Changed[key] = pkg
}

func (l *Stage4Levels) AnalyzeLeaf(pos int, tree *Stage4Tree, leaf *Stage4Leaf) (bool, error) {
	var firstFatherLeaf *luet_pkg.DefaultPackage = nil
	var lastFatherLeaf *luet_pkg.DefaultPackage = nil
	rescan := false
	isLowerLevel := false
	nextLevel := pos - 1

	DebugC(GetAurora().Bold(
		fmt.Sprintf("[P%d-L%d] Levels:\n%s", pos, tree.Id, l.Dump())))

	key := fmt.Sprintf("%s/%s", leaf.Package.GetCategory(), leaf.Package.GetName())

	if pos == len(l.Levels)-1 {
		isLowerLevel = true
	}

	DebugC(fmt.Sprintf("[P%d-L%d] Processing leaf %s - lower lever = %v",
		pos, tree.Id, leaf, isLowerLevel))

	if len(leaf.Father) > 0 {
		lastFatherLeaf = leaf.Father[0]
	}

	if leaf.Counter > 1 {
		// POST: we have different packages
		//       with the same dependency

		if len(leaf.Father) > 0 {
			toRemove := []*luet_pkg.DefaultPackage{}
			for idx, _ := range leaf.Father {
				if idx == 0 {
					firstFatherLeaf = leaf.Father[idx]
					continue
				}
				// The father must to point at the father of the last leaf.
				err := RemoveDependencyFromLuetPackage(
					leaf.Father[idx],
					leaf.Package)
				if err != nil {
					return false, err
				}

				// Add the dependency
				AddDependencyToLuetPackage(leaf.Father[idx], leaf.Father[idx-1])

				//tree.AddDependency(leaf.Father[idx-1], leaf.Father[idx])
				l.AddDependencyRecursive(leaf.Father[idx-1], leaf.Father[idx], pos)
				l.AddChangedPackage(leaf.Father[idx])

				lastFatherLeaf = leaf.Father[idx]
				toRemove = append(toRemove, leaf.Father[idx])
				rescan = true

			}

			if len(toRemove) > 0 {
				for _, f := range toRemove {
					leaf.DelFather(f)
				}
			}
		}

	} else if len(leaf.Father) > 0 {
		firstFatherLeaf = leaf.Father[0]
	}

	for nextLevel >= 0 {
		DebugC(fmt.Sprintf("[P%d-L%d] Analyze upper levels for leaf %s (%d)",
			pos, tree.Id, key, nextLevel))

		treeUpper := l.Levels[nextLevel]

		l2, ok := treeUpper.Map[key]
		if ok {

			// POST: found the package with the selected key
			if len(l2.Father) == 0 {
				DebugC(fmt.Sprintf("For %s father is nil.", key))
				if nextLevel == 0 {
					treeUpper.DropDependency(leaf.Package)
				}
			} else {

				toRemove := []*luet_pkg.DefaultPackage{}

				for idx, _ := range l2.Father {

					fmt.Println("FATHER ", l2)
					DebugC(fmt.Sprintf("[L%d] For %s analyze father %s (%v)",
						treeUpper.Id, key, l2.Father[idx], isLowerLevel))

					if firstFatherLeaf != nil &&
						l2.Father[idx].GetCategory() == firstFatherLeaf.GetCategory() &&
						l2.Father[idx].GetName() == firstFatherLeaf.GetName() {

						DebugC(fmt.Sprintf(
							"[L%d] For key %s the father %s/%s is the first father. Nothing to do.",
							treeUpper.Id, key, l2.Father[idx].GetCategory(), l2.Father[idx].GetName()))
						treeUpper.DropDependency(leaf.Package)

					} else if lastFatherLeaf != nil &&
						l2.Father[idx].GetCategory() == lastFatherLeaf.GetCategory() &&
						l2.Father[idx].GetName() == lastFatherLeaf.GetName() {

						DebugC(fmt.Sprintf(
							"[L%d] For key %s the father %s/%s is the last father. Nothing to do.",
							treeUpper.Id, key, l2.Father[idx].GetCategory(), l2.Father[idx].GetName()))

					} else {

						if len(l2.Father[idx].GetRequires()) == 1 {
							pdep := l2.Father[idx].GetRequires()[0]

							if pdep.GetCategory() != lastFatherLeaf.GetCategory() ||
								pdep.GetName() != lastFatherLeaf.GetName() {
								err := RemoveDependencyFromLuetPackage(
									l2.Father[idx], pdep)
								if err != nil {
									return false, err
								}

								treeUpper.DropDependency(pdep)
								AddDependencyToLuetPackage(l2.Father[idx], lastFatherLeaf)
								//treeUpper.AddDependency(l2.Father[idx], lastFatherLeaf)
								l.AddDependencyRecursive(l2.Father[idx], lastFatherLeaf, nextLevel)

								l.AddChangedPackage(l2.Father[idx])

								lastFatherLeaf = l2.Father[idx]
								rescan = true

							} else {
								DebugC(fmt.Sprintf("[L%d] Father %s/%s has already the right dependency.",
									treeUpper.Id, l2.Father[idx].GetCategory(), l2.Father[idx].GetName()))
							}

						} else if len(l2.Father[idx].GetRequires()) > 1 {
							// The father must to point at the father of the last leaf.
							err := RemoveDependencyFromLuetPackage(
								l2.Father[idx],
								leaf.Package)
							if err != nil {
								return false, err
							}

							AddDependencyToLuetPackage(l2.Father[idx], lastFatherLeaf)

							l.AddChangedPackage(l2.Father[idx])
							l.AddDependencyRecursive(lastFatherLeaf, l2.Father[idx], nextLevel)

							lastFatherLeaf = l2.Father[idx]
							toRemove = append(toRemove, l2.Father[idx])
							rescan = true
						} else {
							DebugC(fmt.Sprintf("For %s father %s/%s is with deps: %s",
								key, l2.Father[idx].GetCategory(), l2.Father[idx].GetName(),
								l2.Father[idx].GetRequires()))
						} ///

					}

				} // end for

				if len(toRemove) > 0 {
					for _, f := range toRemove {
						l2.DelFather(f)
					}
				}

				if nextLevel > 0 {
					// Remove the package from the tree.
					treeUpper.DropDependency(leaf.Package)
				}
			}

		}

		DebugC(GetAurora().Bold(fmt.Sprintf(
			"[P%d-L%d] Completed analysis of the level %d for leaf %s: key found: %v (lasfFather = %s/%s)",
			pos, tree.Id, treeUpper.Id, key, ok, lastFatherLeaf.GetCategory(),
			lastFatherLeaf.GetName())))

		nextLevel--
	}

	return rescan, nil
}

func (l *Stage4Levels) analyzeLevelLeafs(pos int) (bool, error) {
	tree := l.Levels[pos]
	rescan := false

	DebugC(fmt.Sprintf("[%d-%d] Tree:\n%s", tree.Id, pos, tree.Dump()))

	for _, leaf := range *tree.GetMap() {

		r, err := l.AnalyzeLeaf(pos, tree, leaf)
		if err != nil {
			return rescan, err
		}

		if r {
			rescan = true
		}
	} // end for key map

	return rescan, nil
}

func (l *Stage4Levels) Resolve() error {
	// Start from bottom
	pos := len(l.Levels)

	// Check if the levels are sufficient for serialization
	missingLevels := len(l.Levels[0].Map) - pos
	fmt.Println("MISSING ", missingLevels)
	// POST: we need to add levels.
	for i := 0; i < missingLevels; i++ {
		tree := NewStage4Tree(pos + i)
		l.AddTree(tree)
	}
	initialLevels := len(l.Levels)

	for pos > 0 {
		pos--

		rescan, err := l.analyzeLevelLeafs(pos)
		if err != nil {
			return err
		}
		if rescan {
			// Restarting analysis from begin
			pos = initialLevels
		}
	}

	return nil
}

func NewStage4Leaf(pkg, father *luet_pkg.DefaultPackage, pos int) (*Stage4Leaf, error) {
	if pkg == nil || pos < 0 {
		return nil, errors.New("Invalid parameter on create stage4 leaf")
	}

	ans := &Stage4Leaf{
		Package:  pkg,
		Father:   []*luet_pkg.DefaultPackage{},
		Position: pos,
		Counter:  1,
	}

	if father != nil {
		ans.Father = append(ans.Father, father)
	}

	return ans, nil
}

func (l *Stage4Leaf) String() string {

	ans := fmt.Sprintf("%s/%s (%d, %d) ",
		l.Package.GetCategory(), l.Package.GetName(),
		l.Position, l.Counter)

	if len(l.Father) > 0 {
		ans += " father: ["
		for _, f := range l.Father {
			ans += fmt.Sprintf("%s/%s, ", f.GetCategory(), f.GetName())
		}
		ans += "]"
	}

	return ans
}

func (l *Stage4Leaf) AddFather(father *luet_pkg.DefaultPackage) {

	// Check if the father is already present
	notFound := true

	for _, f := range l.Father {
		if f.GetCategory() == father.GetCategory() &&
			f.GetName() == father.GetName() {
			notFound = false
			break
		}
	}

	if notFound {
		l.Father = append(l.Father, father)
		l.Counter++
	}
}

func (l *Stage4Leaf) DelFather(father *luet_pkg.DefaultPackage) {
	pos := -1
	for idx, f := range l.Father {
		if f.GetCategory() == father.GetCategory() &&
			f.GetName() == father.GetName() {
			pos = idx
			break
		}
	}

	if pos >= 0 {
		l.Father = append(l.Father[:pos], l.Father[pos+1:]...)
		l.Counter--

		DebugC(fmt.Sprintf("From leaf %s/%s delete father %s/%s (%d)",
			l.Package.GetCategory(), l.Package.GetName(),
			father.GetCategory(), father.GetName(), l.Counter,
		))
	}
}

func NewStage4Tree(id int) *Stage4Tree {
	return &Stage4Tree{
		Id:   id,
		Map:  make(map[string]*Stage4Leaf, 0),
		Deps: []*luet_pkg.DefaultPackage{},
	}
}

func (t *Stage4Tree) Dump() string {
	ans := fmt.Sprintf("[%d] Map: [ ", t.Id)

	for k, v := range t.Map {
		ans += fmt.Sprintf("%s-%d ", k, v.Counter)

		if v.Counter > 0 {
			ans += "("
			for _, father := range v.Father {
				ans += fmt.Sprintf("%s/%s, ", father.GetCategory(), father.GetName())
			}
			ans += ")"
		}

		ans += ", "
	}
	ans += "]\n"
	ans += fmt.Sprintf("[%d] Deps: [ ", t.Id)
	for _, d := range t.Deps {
		ans += fmt.Sprintf("%s/%s", d.GetCategory(), d.GetName())

		if len(d.GetRequires()) > 0 {
			ans += " ("
			for _, d2 := range d.GetRequires() {
				ans += fmt.Sprintf("%s/%s, ", d2.GetCategory(), d2.GetName())
			}
			ans += " )"
		}
		ans += ", "

	}
	ans += "]\n"
	return ans
}

func (t *Stage4Tree) GetMap() *map[string]*Stage4Leaf {
	return &t.Map
}

func (t *Stage4Tree) GetDeps() *[]*luet_pkg.DefaultPackage { return &t.Deps }

func (t *Stage4Tree) GetDependency(pos int) (*luet_pkg.DefaultPackage, error) {
	if len(t.Deps) <= pos {
		return nil, errors.New("Invalid position")
	}

	return t.Deps[pos], nil
}

func (t *Stage4Tree) DropDependency(p *luet_pkg.DefaultPackage) error {
	key := fmt.Sprintf("%s/%s", p.GetCategory(), p.GetName())
	leaf, ok := t.Map[key]
	if !ok {
		return errors.New(fmt.Sprintf("Package %s is not present on tree.", key))
	}

	if leaf.Counter > 1 {
		leaf.Counter--
		DebugC(fmt.Sprintf("For %s decrement counter.", key))
		return nil
	}

	DebugC(fmt.Sprintf("[%d] Dropping dependency %s...", t.Id, key), len(t.Map))

	delete(t.Map, key)

	if len(t.Deps) == (leaf.Position - 1) {
		// POST: just drop last element
		t.Deps = t.Deps[:leaf.Position-1]

	} else if leaf.Position == 0 {
		// POST: just drop the first element
		t.Deps = t.Deps[1:]

	} else {
		// POST: drop element between other elements
		t.Deps = append(t.Deps[:leaf.Position], t.Deps[leaf.Position+1:]...)
	}

	for k, v := range t.Map {
		if v.Position > leaf.Position {
			t.Map[k].Position--
		}
	}

	DebugC(fmt.Sprintf("[%d] Dropped dependency %s...", t.Id, key), len(t.Map))

	return nil
}

func (t *Stage4Tree) AddDependency(p, father *luet_pkg.DefaultPackage) error {
	fatherKey := "-"
	if father != nil {
		fatherKey = fmt.Sprintf("%s/%s", father.GetCategory(), father.GetName())
	}
	DebugC(fmt.Sprintf("[L%d] Adding dep %s/%s with father %s",
		t.Id, p.GetCategory(), p.GetName(), fatherKey))

	if leaf, ok := t.Map[fmt.Sprintf("%s/%s", p.GetCategory(), p.GetName())]; ok {
		leaf.AddFather(father)
		return nil
	}

	leaf, err := NewStage4Leaf(p, father, len(t.Deps))
	if err != nil {
		return err
	}

	t.Map[fmt.Sprintf("%s/%s", p.GetCategory(), p.GetName())] = leaf
	t.Deps = append(t.Deps, p)

	return nil
}

// TODO: Move in another place
func RemoveDependencyFromLuetPackage(pkg, dep *luet_pkg.DefaultPackage) error {

	DebugC(GetAurora().Bold(fmt.Sprintf("Removing dep %s/%s from package %s/%s...",
		dep.GetCategory(), dep.GetName(),
		pkg.GetCategory(), pkg.GetName(),
	)))

	pos := -1
	for idx, d := range pkg.GetRequires() {
		if d.GetCategory() == dep.GetCategory() && d.GetName() == dep.GetName() {
			pos = idx
			break
		}
	}

	if pos < 0 {
		// Ignore error until we fix father cleanup
		Warning(fmt.Sprintf("Dependency %s/%s not found on package %s/%s",
			dep.GetCategory(), dep.GetName(),
			pkg.GetCategory(), pkg.GetName(),
		))
		/*
			return nil
		*/
		return errors.New(
			fmt.Sprintf("Dependency %s/%s not found on package %s/%s",
				dep.GetCategory(), dep.GetName(),
				pkg.GetCategory(), pkg.GetName(),
			))
	}

	pkg.PackageRequires = append(pkg.PackageRequires[:pos], pkg.PackageRequires[pos+1:]...)

	return nil
}

func AddDependencyToLuetPackage(pkg, dep *luet_pkg.DefaultPackage) {

	DebugC(fmt.Sprintf("Adding %s/%s dep to package %s/%s",
		dep.GetCategory(), dep.GetName(),
		pkg.GetCategory(), pkg.GetName()))

	// Check if dependency is already present.
	notFound := true
	for _, d := range pkg.GetRequires() {
		if d.GetCategory() == dep.GetCategory() && d.GetName() == dep.GetName() {
			notFound = false
			break
		}
	}
	if notFound {
		DebugC(GetAurora().Bold(fmt.Sprintf("Added %s/%s dep to package %s/%s",
			dep.GetCategory(), dep.GetName(),
			pkg.GetCategory(), pkg.GetName())))

		// TODO: check if we need set PackageRequires
		pkg.PackageRequires = append(pkg.PackageRequires,
			&luet_pkg.DefaultPackage{
				Category: dep.GetCategory(),
				Name:     dep.GetName(),
				Version:  ">=0",
			},
		)
	}

}
