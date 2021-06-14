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
	Context("Stage4 - PackageHasAncient", func() {

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

		It("HasAncient1", func() {
			levels := NewStage4LevelsWithSize(4)

			errs := make([]error, 0)

			// A (deps B, C)
			errs = append(errs, levels.AddDependency(p1, nil, 0))
			// B (deps C)
			errs = append(errs, levels.AddDependency(p2, nil, 0))
			// C (deps D)
			errs = append(errs, levels.AddDependency(p3, nil, 0))
			// D
			errs = append(errs, levels.AddDependency(p4, nil, 0))

			// Tree2
			// B (father A)
			errs = append(errs, levels.AddDependency(p2, p1, 1))
			// C (father A)
			errs = append(errs, levels.AddDependency(p3, p1, 1))
			// C (father B)
			errs = append(errs, levels.AddDependency(p3, p2, 1))
			// D (father C)
			errs = append(errs, levels.AddDependency(p4, p3, 1))

			// Tree3
			// C (father B)
			errs = append(errs, levels.AddDependency(p3, p2, 2))
			// D (father C)
			errs = append(errs, levels.AddDependency(p4, p3, 2))

			// Tree4
			// D (father C)
			errs = append(errs, levels.AddDependency(p4, p3, 3))

			for i, _ := range errs {
				Expect(errs[i]).Should(BeNil())
			}

			hasAncient, err := levels.PackageHasAncient(p1, p4, 0)
			Expect(err).Should(BeNil())
			Expect(hasAncient).Should(Equal(true))

			hasAncient, err = levels.PackageHasAncient(p3, p1, 0)
			Expect(err).Should(BeNil())
			Expect(hasAncient).Should(Equal(false))
		})
	})

})
