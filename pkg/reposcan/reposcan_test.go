/*
Copyright (C) 2021  Daniele Rondina <geaaru@sabayonlinux.org>

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
package reposcan_test

import (
	"errors"
	//	"fmt"

	. "github.com/Luet-lab/luet-portage-converter/pkg/reposcan"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reposcan", func() {
	Context("Json parser", func() {
		It("Load raw JSON", func() {

			source := `
{
  "cache_data_version": "1.0.5",
  "atoms": {
    "x11-misc/services_notificator-svn-0.1": {
      "relations": [],
      "relations_by_kind": {},
      "category": "x11-misc",
      "revision": "0",
      "package": "services_notificator-svn",
      "catpkg": "x11-misc/services_notificator-svn",
      "atom": "x11-misc/services_notificator-svn-0.1",
      "eclasses": [],
      "kit": "geaaru_overlay",
      "branch": "master",
      "metadata": null,
      "md5": "1750e323fe2cf4a6171da167afde6157",
      "metadata_out": "",
      "manifest_md5": "0a8c57df4540f66abba8f10d048fae38"
    }
  }
}
`
			resolver := NewRepoScanResolver()
			err := resolver.LoadRawJson(source, "test")
			Expect(err).Should(BeNil())

			err = resolver.BuildMap()
			Expect(err).Should(BeNil())

			val, ok := resolver.GetMap()["x11-misc/services_notificator-svn"]
			atom := val[0]
			Expect(ok).Should(Equal(true))
			Expect(atom.Package).Should(Equal("services_notificator-svn"))
			Expect(atom.Category).Should(Equal("x11-misc"))
			Expect(atom.Kit).Should(Equal("geaaru_overlay"))
		})

	})

	Context("Json parser2", func() {

		It("Load Json file with deps1", func() {

			resolver := NewRepoScanResolver()
			err := resolver.LoadJson("../../tests/fixtures/reposcan-01.json")
			Expect(err).Should(BeNil())
			resolver.SetIgnoreMissingDeps(true)
			err = resolver.BuildMap()
			Expect(err).Should(BeNil())

			val, ok := resolver.GetMap()["net-dns/noip-updater"]
			atom := val[0]
			Expect(ok).Should(Equal(true))
			Expect(atom.Package).Should(Equal("noip-updater"))
			Expect(atom.Category).Should(Equal("net-dns"))
			Expect(atom.Kit).Should(Equal("geaaru_overlay"))

			solution, err := resolver.Resolve("net-dns/noip-updater")
			Expect(err).Should(BeNil())

			Expect(solution.Package.GetPackageName()).Should(Equal("net-dns/noip-updater"))
			Expect(solution.Package.Slot).Should(Equal("0"))
			Expect(solution.Package.License).Should(Equal("GPL-2"))
			Expect(len(solution.BuildDeps)).Should(Equal(1))
			Expect(solution.BuildDeps[0]).Should(Equal(gentoo.GentooPackage{
				Name:          "gcc",
				Category:      "sys-devel",
				Version:       "10.2.0",
				VersionSuffix: "-r3",
				Slot:          "10",
				License:       "GPL-3+ LGPL-3+ || ( GPL-3+ libgcc libstdc++ gcc-runtime-library-exception-3.1 ) FDL-1.3+",
				Repository:    "gentoo",
				Condition:     gentoo.PkgCondEqual,
			}))
			Expect(len(solution.RuntimeDeps)).Should(Equal(0))

		})

		It("Load Json file with deps2", func() {

			resolver := NewRepoScanResolver()
			err := resolver.LoadJson("../../tests/fixtures/reposcan-01.json")
			Expect(err).Should(BeNil())
			resolver.SetIgnoreMissingDeps(true)
			err = resolver.BuildMap()
			Expect(err).Should(BeNil())

			val, ok := resolver.GetMap()["net-dns/noip-updater"]
			atom := val[0]
			Expect(ok).Should(Equal(true))
			Expect(atom.Package).Should(Equal("noip-updater"))
			Expect(atom.Category).Should(Equal("net-dns"))
			Expect(atom.Kit).Should(Equal("geaaru_overlay"))

			_, err = resolver.Resolve(">=net-dns/noip-updater-2.2.0")
			Expect(err).Should(Equal(errors.New("No package found matching >=net-dns/noip-updater-2.2.0")))

		})

		It("Load Json file with deps3", func() {

			resolver := NewRepoScanResolver()
			err := resolver.LoadJson("../../tests/fixtures/reposcan-01.json")
			Expect(err).Should(BeNil())
			resolver.SetIgnoreMissingDeps(true)
			err = resolver.BuildMap()
			Expect(err).Should(BeNil())

			val, ok := resolver.GetMap()["net-dns/noip-updater"]
			atom := val[0]
			Expect(ok).Should(Equal(true))
			Expect(atom.Package).Should(Equal("noip-updater"))
			Expect(atom.Category).Should(Equal("net-dns"))
			Expect(atom.Kit).Should(Equal("geaaru_overlay"))

			_, err = resolver.Resolve("net-dns/noip-updater:2")
			Expect(err).Should(Equal(errors.New("No package found matching net-dns/noip-updater:2")))

		})

		It("Load Json file with deps4 (with runtime deps)", func() {

			resolver := NewRepoScanResolver()
			err := resolver.LoadJson("../../tests/fixtures/reposcan-01.json")
			Expect(err).Should(BeNil())
			resolver.SetIgnoreMissingDeps(true)
			err = resolver.BuildMap()
			Expect(err).Should(BeNil())

			val, ok := resolver.GetMap()["sys-devel/gcc"]
			atom := val[0]
			Expect(ok).Should(Equal(true))
			Expect(atom.Package).Should(Equal("gcc"))
			Expect(atom.Category).Should(Equal("sys-devel"))
			Expect(atom.Kit).Should(Equal("gentoo"))
			rdeps, err := atom.GetRuntimeDeps()
			Expect(err).Should(BeNil())
			Expect(len(rdeps)).Should(Equal(8))

			solution, err := resolver.Resolve("sys-devel/gcc:9.3.0")
			Expect(err).Should(BeNil())
			// We have only one runtime deps because atom aren't
			// available on json file.
			Expect(len(solution.RuntimeDeps)).Should(Equal(1))
		})
	})

})
