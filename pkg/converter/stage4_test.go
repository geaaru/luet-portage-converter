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

package converter_test

import (
	"fmt"

	luet_pkg "github.com/mudler/luet/pkg/package"

	. "github.com/Luet-lab/luet-portage-converter/pkg/converter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func NewPackage(name, cat, version string, deps []*luet_pkg.DefaultPackage) *luet_pkg.DefaultPackage {
	return &luet_pkg.DefaultPackage{
		Name:            name,
		Category:        cat,
		Version:         version,
		PackageRequires: deps,
	}
}

var _ = Describe("Converter", func() {
	Context("Stage4", func() {
		fmt.Println("TEST")
		/*
			p1 := NewPackage("A", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("B", "cat1", ">=0", nil),
					NewPackage("C", "cat1", ">=0", nil),
				},
			)

			p2 := NewPackage("B", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("D", "cat1", ">=0", nil),
					NewPackage("C", "cat1", ">=0", nil),
				},
			)

			p3 := NewPackage("C", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("D", "cat1", ">=0", nil),
				},
			)

			p4 := NewPackage("D", "cat1", "1.0", nil)
			It("Drop dependency1", func() {

				p := NewPackage("C", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("D", "cat1", ">=0", nil),
						NewPackage("A", "cat1", ">=0", nil),
					},
				)

				err := RemoveDependencyFromLuetPackage(p, p1)
				Expect(err).Should(BeNil())
				Expect(len(p.GetRequires())).Should(Equal(1))
				Expect(p.GetRequires()[0].GetName()).Should(Equal(p4.GetName()))
			})
			It("Drop dependency2", func() {

				p := NewPackage("C", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("D", "cat1", ">=0", nil),
						NewPackage("A", "cat1", ">=0", nil),
					},
				)
				err := RemoveDependencyFromLuetPackage(p, p4)
				Expect(err).Should(BeNil())
				Expect(len(p.GetRequires())).Should(Equal(1))
				Expect(p.GetRequires()[0].GetName()).Should(Equal(p1.GetName()))
			})

			It("Drop dependency3", func() {

				p := NewPackage("C", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("D", "cat1", ">=0", nil),
						NewPackage("A", "cat1", ">=0", nil),
					},
				)
				err := RemoveDependencyFromLuetPackage(p, p2)
				Expect(err).ShouldNot(BeNil())
			})

			It("Drop leaf in tail", func() {
				levels := NewStage4Levels()
				tree1 := NewStage4Tree(1)
				levels.AddTree(tree1)

				err1 := tree1.AddDependency(p2, p1)
				err2 := tree1.AddDependency(p3, p1)

				err3 := tree1.DropDependency(p3)

				Expect(err1).Should(BeNil())
				Expect(err2).Should(BeNil())
				Expect(err3).Should(BeNil())

				Expect(len((*tree1.GetDeps()))).Should(Equal(1))
				Expect((*tree1.GetDeps())[0]).Should(Equal(p2))
			})

			It("Check error on Drop dependency", func() {
				levels := NewStage4Levels()
				tree1 := NewStage4Tree(1)
				levels.AddTree(tree1)

				err1 := tree1.AddDependency(p2, p1)
				err2 := tree1.AddDependency(p3, p1)
				err3 := tree1.DropDependency(p4)

				Expect(err1).Should(BeNil())
				Expect(err2).Should(BeNil())
				Expect(err3).ShouldNot(BeNil())
			})

			It("Drop leaf in head", func() {
				levels := NewStage4Levels()
				tree1 := NewStage4Tree(1)
				levels.AddTree(tree1)

				err1 := tree1.AddDependency(p2, p1)
				err2 := tree1.AddDependency(p3, p1)

				err3 := tree1.DropDependency(p2)

				Expect(err1).Should(BeNil())
				Expect(err2).Should(BeNil())
				Expect(err3).Should(BeNil())

				Expect(len((*tree1.GetDeps()))).Should(Equal(1))
				Expect((*tree1.GetDeps())[0]).Should(Equal(p3))
			})

			It("Drop leaf in middle", func() {
				levels := NewStage4Levels()
				tree1 := NewStage4Tree(1)
				levels.AddTree(tree1)

				err1 := tree1.AddDependency(p2, p1)
				err2 := tree1.AddDependency(p3, p1)
				err3 := tree1.AddDependency(p4, p1)

				err4 := tree1.DropDependency(p3)

				Expect(err1).Should(BeNil())
				Expect(err2).Should(BeNil())
				Expect(err3).Should(BeNil())
				Expect(err4).Should(BeNil())

				Expect(len((*tree1.GetDeps()))).Should(Equal(2))
				Expect((*tree1.GetDeps())[0]).Should(Equal(p2))
				Expect((*tree1.GetDeps())[1]).Should(Equal(p4))
			})

			It("Load Tree", func() {
				levels := NewStage4Levels()
				tree1 := NewStage4Tree(1)
				tree2 := NewStage4Tree(2)
				tree3 := NewStage4Tree(3)
				levels.AddTree(tree1)
				levels.AddTree(tree2)
				levels.AddTree(tree3)

				err1 := tree1.AddDependency(p2, p1)
				err2 := tree1.AddDependency(p3, p1)

				err3 := tree2.AddDependency(p4, p2)
				err4 := tree2.AddDependency(p3, p2)

				err5 := tree3.AddDependency(p4, p3)

				Expect(err1).Should(BeNil())
				Expect(err2).Should(BeNil())
				Expect(err3).Should(BeNil())
				Expect(err4).Should(BeNil())
				Expect(err5).Should(BeNil())

				Expect(tree1.GetDependency(0)).Should(Equal(p2))
				Expect(tree1.GetDependency(1)).Should(Equal(p3))
				Expect(tree2.GetDependency(0)).Should(Equal(p4))
				Expect(tree2.GetDependency(1)).Should(Equal(p3))
				Expect(tree3.GetDependency(0)).Should(Equal(p4))
			})
		*/
		It("Check clean of the father", func() {

			p1 := NewPackage("A", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("B", "cat1", ">=0", nil),
					NewPackage("C", "cat1", ">=0", nil),
				},
			)

			p2 := NewPackage("B", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("C", "cat1", ">=0", nil),
				},
			)

			p3 := NewPackage("C", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("D", "cat1", ">=0", nil),
				},
			)

			p4 := NewPackage("D", "cat1", "1.0", nil)

			levels := NewStage4Levels()
			tree1 := NewStage4Tree(1)
			tree2 := NewStage4Tree(2)
			tree3 := NewStage4Tree(3)
			tree4 := NewStage4Tree(4)
			levels.AddTree(tree1)
			levels.AddTree(tree2)
			levels.AddTree(tree3)
			levels.AddTree(tree4)

			errs := make([]error, 0)

			// A (deps B, C)
			errs = append(errs, tree1.AddDependency(p1, nil))
			// B (deps C)
			errs = append(errs, tree1.AddDependency(p2, nil))
			// C (deps D)
			errs = append(errs, tree1.AddDependency(p3, nil))
			// D
			errs = append(errs, tree1.AddDependency(p4, nil))

			// Tree2
			// B (father A)
			errs = append(errs, tree2.AddDependency(p2, p1))
			// C (father A)
			errs = append(errs, tree2.AddDependency(p3, p1))
			// C (father B)
			errs = append(errs, tree2.AddDependency(p3, p2))
			// D (father C)
			errs = append(errs, tree2.AddDependency(p4, p3))

			// Tree3
			// C (father B)
			errs = append(errs, tree3.AddDependency(p3, p2))
			// D (father C)
			errs = append(errs, tree3.AddDependency(p4, p3))

			// Tree4
			// D (father C)
			errs = append(errs, tree4.AddDependency(p4, p3))

			for i, _ := range errs {
				Expect(errs[i]).Should(BeNil())
			}

			for i, _ := range levels.Levels {
				Expect(levels.Levels[i].Id).Should(Equal(i + 1))
			}

			for i, _ := range levels.Levels {
				Expect(len(levels.Levels[i].Map)).Should(Equal(4 - i))
			}

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3, p4},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2, p3, p4},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1, p2, p3, p4},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 1,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p1, p2},
						Position: 1,
						Counter:  2,
					},
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 2,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 1,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 2,
						Counter:  1,
					},
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 3,
						Counter:  1,
					},
				},
			))

			fmt.Println("LEVELS\n", levels.Dump())

			err := levels.AnalyzeLeaf(3, levels.Levels[3],
				levels.Levels[3].Map["cat1/D"],
			)
			Expect(err).Should(BeNil())

			fmt.Println("LEVELS RESOLVED\n", levels.Dump())

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2, p3},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1, p2, p3},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p1, p2},
						Position: 1,
						Counter:  2,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 1,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 2,
						Counter:  1,
					},
				},
			))

			fmt.Println("====================================")
			fmt.Println("ANALYZE cat1/C level 3")
			fmt.Println("====================================")
			err = levels.AnalyzeLeaf(2, levels.Levels[2],
				levels.Levels[2].Map["cat1/C"],
			)
			Expect(err).Should(BeNil())
			fmt.Println("LEVELS RESOLVED\n", levels.Dump())

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1, p2},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 1,
						Counter:  1,
					},
				},
			))

			fmt.Println("====================================")
			fmt.Println("ANALYZE cat1/B level 2")
			fmt.Println("====================================")
			err = levels.AnalyzeLeaf(1, levels.Levels[1],
				levels.Levels[1].Map["cat1/B"],
			)
			Expect(err).Should(BeNil())
			fmt.Println("LEVELS RESOLVED\n", levels.Dump())

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
				},
			))

			fmt.Println("====================================")
			fmt.Println("ANALYZE cat1/A level 1")
			fmt.Println("====================================")
			err = levels.AnalyzeLeaf(0, levels.Levels[0],
				levels.Levels[0].Map["cat1/A"],
			)
			Expect(err).Should(BeNil())
			fmt.Println("LEVELS RESOLVED\n", levels.Dump())

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
				},
			))

		})

		It("Resolve Levels test1", func() {

			errs := make([]error, 0)

			p1 := NewPackage("A", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("B", "cat1", ">=0", nil),
					NewPackage("C", "cat1", ">=0", nil),
				},
			)

			p2 := NewPackage("B", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("D", "cat1", ">=0", nil),
					NewPackage("C", "cat1", ">=0", nil),
				},
			)

			p3 := NewPackage("C", "cat1", "1.0",
				[]*luet_pkg.DefaultPackage{
					NewPackage("D", "cat1", ">=0", nil),
				},
			)

			p4 := NewPackage("D", "cat1", "1.0", nil)

			levels := NewStage4Levels()
			tree1 := NewStage4Tree(1)
			tree2 := NewStage4Tree(2)
			tree3 := NewStage4Tree(3)
			tree4 := NewStage4Tree(4)
			levels.AddTree(tree1)
			levels.AddTree(tree2)
			levels.AddTree(tree3)
			levels.AddTree(tree4)

			// Level1
			// A (deps B, C)
			errs = append(errs, tree1.AddDependency(p1, nil))
			// B (deps D, C)
			errs = append(errs, tree1.AddDependency(p2, nil))
			// C (dep D)
			errs = append(errs, tree1.AddDependency(p3, nil))
			// D
			errs = append(errs, tree1.AddDependency(p4, nil))

			// Level2
			// B (father A)
			errs = append(errs, tree2.AddDependency(p2, p1))
			// C (father A)
			errs = append(errs, tree2.AddDependency(p3, p1))
			// D (father B)
			errs = append(errs, tree2.AddDependency(p4, p2))
			// C (father B)
			errs = append(errs, tree2.AddDependency(p3, p2))
			// D (father C)
			errs = append(errs, tree2.AddDependency(p4, p3))

			// Level3
			// D (father B)
			errs = append(errs, tree3.AddDependency(p4, p2))
			// C (father B)
			errs = append(errs, tree3.AddDependency(p3, p2))
			// D (father C)
			errs = append(errs, tree3.AddDependency(p4, p3))

			// Level4
			errs = append(errs, tree4.AddDependency(p4, p3))

			for i, _ := range errs {
				Expect(errs[i]).Should(BeNil())
			}

			for i, _ := range levels.Levels {
				Expect(levels.Levels[i].Id).Should(Equal(i + 1))
			}

			for i, _ := range levels.Levels {
				Expect(len(levels.Levels[i].Map)).Should(Equal(4 - i))
			}

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4, p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2, p3, p4},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1, p2, p3, p4},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))

			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 1,
						Counter:  1,
					},
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p2, p3},
						Position: 0,
						Counter:  2,
					},
				},
			))

			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p1, p2},
						Position: 1,
						Counter:  2,
					},
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p2, p3},
						Position: 2,
						Counter:  2,
					},
				},
			))

			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 1,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 2,
						Counter:  1,
					},
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 3,
						Counter:  1,
					},
				},
			))

			fmt.Println("====================================")
			fmt.Println("ANALYZE cat1/D level 4")
			fmt.Println("====================================")

			fmt.Println("LEVELS\n", levels.Dump())

			err := levels.AnalyzeLeaf(3, levels.Levels[3],
				levels.Levels[3].Map["cat1/D"],
			)
			Expect(err).Should(BeNil())

			fmt.Println("LEVELS RESOLVED\n", levels.Dump())

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2, p3},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1, p2, p3},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
				},
			))

			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p1, p2},
						Position: 1,
						Counter:  2,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 1,
						Counter:  1,
					},
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 2,
						Counter:  1,
					},
				},
			))

			fmt.Println("====================================")
			fmt.Println("ANALYZE cat1/C level 3")
			fmt.Println("====================================")

			fmt.Println("LEVELS\n", levels.Dump())

			err = levels.AnalyzeLeaf(2, levels.Levels[2],
				levels.Levels[2].Map["cat1/C"],
			)
			Expect(err).Should(BeNil())

			fmt.Println("LEVELS RESOLVED\n", levels.Dump())

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1, p2},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
				},
			))

			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 1,
						Counter:  1,
					},
				},
			))

			fmt.Println("====================================")
			fmt.Println("ANALYZE cat1/B level 2")
			fmt.Println("====================================")

			fmt.Println("LEVELS\n", levels.Dump())

			err = levels.AnalyzeLeaf(1, levels.Levels[1],
				levels.Levels[1].Map["cat1/B"],
			)
			Expect(err).Should(BeNil())

			fmt.Println("LEVELS RESOLVED\n", levels.Dump())

			// Check Deps
			Expect(levels.Levels[3].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p4},
			))
			Expect(levels.Levels[2].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p3},
			))
			Expect(levels.Levels[1].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p2},
			))
			Expect(levels.Levels[0].Deps).Should(Equal(
				[]*luet_pkg.DefaultPackage{p1},
			))

			// Check Maps
			Expect(levels.Levels[3].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/D": &Stage4Leaf{
						Package:  p4,
						Father:   []*luet_pkg.DefaultPackage{p3},
						Position: 0,
						Counter:  1,
					}},
			))
			Expect(levels.Levels[2].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/C": &Stage4Leaf{
						Package:  p3,
						Father:   []*luet_pkg.DefaultPackage{p2},
						Position: 0,
						Counter:  1,
					},
				},
			))

			Expect(levels.Levels[1].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/B": &Stage4Leaf{
						Package:  p2,
						Father:   []*luet_pkg.DefaultPackage{p1},
						Position: 0,
						Counter:  1,
					},
				},
			))
			Expect(levels.Levels[0].Map).Should(Equal(
				map[string]*Stage4Leaf{
					"cat1/A": &Stage4Leaf{
						Package:  p1,
						Father:   []*luet_pkg.DefaultPackage{},
						Position: 0,
						Counter:  1,
					},
				},
			))
		})
		/*
			It("Resolve Levels test2", func() {

				p1 := NewPackage("A", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("C", "cat1", ">=0", nil),
					},
				)

				p2 := NewPackage("B", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("C", "cat1", ">=0", nil),
					},
				)

				p3 := NewPackage("C", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("D", "cat1", ">=0", nil),
					},
				)

				p4 := NewPackage("D", "cat1", "1.0", nil)

				levels := NewStage4Levels()
				tree1 := NewStage4Tree(1)
				tree2 := NewStage4Tree(2)
				tree3 := NewStage4Tree(3)
				levels.AddTree(tree1)
				levels.AddTree(tree2)
				levels.AddTree(tree3)

				err1 := tree1.AddDependency(p1, nil)
				err2 := tree1.AddDependency(p2, nil)

				err3 := tree2.AddDependency(p3, p1)
				err4 := tree2.AddDependency(p3, p2)

				err5 := tree3.AddDependency(p4, p3)

				Expect(err1).Should(BeNil())
				Expect(err2).Should(BeNil())
				Expect(err3).Should(BeNil())
				Expect(err4).Should(BeNil())
				Expect(err5).Should(BeNil())

				fmt.Println("LEVELS\n", levels.Dump())
				err6 := levels.Resolve()
				fmt.Println("LEVELS RESOLVED\n", levels.Dump())
				Expect(err6).Should(BeNil())
				Expect(len(p1.GetRequires())).Should(Equal(1))
				Expect(p1).Should(Equal(
					NewPackage("A", "cat1", "1.0",
						[]*luet_pkg.DefaultPackage{
							NewPackage("C", "cat1", ">=0", nil),
						},
					)),
				)

				Expect(p2).Should(Equal(
					NewPackage("B", "cat1", "1.0",
						[]*luet_pkg.DefaultPackage{
							NewPackage("A", "cat1", ">=0", nil),
						},
					)),
				)

				Expect(p3).Should(Equal(
					NewPackage("C", "cat1", "1.0",
						[]*luet_pkg.DefaultPackage{
							NewPackage("D", "cat1", ">=0", nil),
						},
					)),
				)

			})
			It("Resolve Levels test3", func() {

				levels := NewStage4Levels()
				tree1 := NewStage4Tree(1)
				tree2 := NewStage4Tree(2)
				tree3 := NewStage4Tree(3)
				levels.AddTree(tree1)
				levels.AddTree(tree2)
				levels.AddTree(tree3)

				p1 := NewPackage("A", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("E", "cat1", ">=0", nil),
						NewPackage("C", "cat1", ">=0", nil),
					},
				)
				err1 := tree1.AddDependency(p1, nil)

				p2 := NewPackage("B", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("C", "cat1", ">=0", nil),
					},
				)
				err2 := tree1.AddDependency(p2, nil)

				p3 := NewPackage("C", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("D", "cat1", ">=0", nil),
					},
				)
				err3 := tree1.AddDependency(p3, nil)

				p4 := NewPackage("D", "cat1", "1.0", nil)
				err4 := tree1.AddDependency(p4, nil)

				p5 := NewPackage("E", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("D", "cat1", ">=0", nil),
					},
				)
				err5 := tree1.AddDependency(p5, nil)

				err6 := tree2.AddDependency(p5, p1)
				err7 := tree2.AddDependency(p3, p1)
				err8 := tree2.AddDependency(p3, p2)
				err9 := tree2.AddDependency(p4, p3)
				err10 := tree2.AddDependency(p4, p5)

				err11 := tree3.AddDependency(p4, p3)
				err12 := tree3.AddDependency(p4, p5)

				Expect(err1).Should(BeNil())
				Expect(err2).Should(BeNil())
				Expect(err3).Should(BeNil())
				Expect(err4).Should(BeNil())
				Expect(err5).Should(BeNil())
				Expect(err6).Should(BeNil())
				Expect(err7).Should(BeNil())
				Expect(err8).Should(BeNil())
				Expect(err9).Should(BeNil())
				Expect(err10).Should(BeNil())
				Expect(err11).Should(BeNil())
				Expect(err12).Should(BeNil())

				fmt.Println("LEVELS\n", levels.Dump())

				err13 := levels.Resolve()
				fmt.Println("LEVELS RESOLVED\n", levels.Dump())
				Expect(err13).Should(BeNil())
				//			Expect(len(p1.GetRequires())).Should(Equal(1))
				Expect(p1).Should(Equal(
					NewPackage("A", "cat1", "1.0",
						[]*luet_pkg.DefaultPackage{
							NewPackage("E", "cat1", ">=0", nil),
						},
					)),
				)
				Expect(p2).Should(Equal(
					NewPackage("B", "cat1", "1.0",
						[]*luet_pkg.DefaultPackage{
							NewPackage("A", "cat1", ">=0", nil),
						},
					)),
				)

				Expect(p3).Should(Equal(
					NewPackage("C", "cat1", "1.0",
						[]*luet_pkg.DefaultPackage{
							NewPackage("D", "cat1", ">=0", nil),
						},
					)),
				)

				Expect(p4).Should(Equal(
					NewPackage("D", "cat1", "1.0", nil),
				))

				Expect(p5).Should(Equal(
					NewPackage("E", "cat1", "1.0",
						[]*luet_pkg.DefaultPackage{
							NewPackage("C", "cat1", ">=0", nil),
						},
					)),
				)
			})

		*/
	})

})
