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

		It("Resolve Levels test1", func() {

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

			fmt.Println("LEVELS\n", levels.Dump())
			err6 := levels.Resolve()
			fmt.Println("LEVELS RESOLVED\n", levels.Dump())
			Expect(err6).Should(BeNil())
			Expect(len(p1.GetRequires())).Should(Equal(1))
			Expect(p1).Should(Equal(
				NewPackage("A", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("B", "cat1", ">=0", nil),
					},
				)),
			)

			Expect(p2).Should(Equal(
				NewPackage("B", "cat1", "1.0",
					[]*luet_pkg.DefaultPackage{
						NewPackage("C", "cat1", ">=0", nil),
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

	})
})
