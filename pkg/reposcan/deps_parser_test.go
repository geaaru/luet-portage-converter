package reposcan

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
)

var _ = Describe("GentooBuilder", func() {

	Context("Parse Dependencies 1", func() {

		rdepend := `
	app-crypt/sbsigntools
	x11-themes/sabayon-artwork-grub
	sys-boot/os-prober
	app-arch/xz-utils
	>=sys-libs/ncurses-5.2-r5:0=
`
		gr, err := ParseDependencies(rdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(5))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sbsigntools",
						Category: "app-crypt",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sabayon-artwork-grub",
						Category: "x11-themes",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep5", func() {
			Expect(*gr.Dependencies[4]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:          "ncurses",
						Category:      "sys-libs",
						Slot:          "0=",
						Version:       "5.2",
						VersionSuffix: "-r5",
						Condition:     _gentoo.PkgCondGreaterEqual,
					},
				},
			))
		})

	})

	Context("Parse Dependencies 2", func() {

		rdepend := `
	app-crypt/sbsigntools
	x11-themes/sabayon-artwork-grub
	sys-boot/os-prober
	app-arch/xz-utils
	>=sys-libs/ncurses-5.2-r5:0=
	mount? ( sys-fs/fuse )
`
		gr, err := ParseDependencies(rdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(6))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sbsigntools",
						Category: "app-crypt",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sabayon-artwork-grub",
						Category: "x11-themes",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep5", func() {
			Expect(*gr.Dependencies[4]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:          "ncurses",
						Category:      "sys-libs",
						Slot:          "0=",
						Version:       "5.2",
						VersionSuffix: "-r5",
						Condition:     _gentoo.PkgCondGreaterEqual,
					},
				},
			))
		})

		It("Check dep6", func() {
			Expect(*gr.Dependencies[5]).Should(Equal(
				GentooDependency{
					Use:          "mount",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:     "fuse",
								Category: "sys-fs",
								Slot:     "0",
							},
						},
					},
					Dep: nil,
				},
			))
		})

	})

	Context("Parse Dependencies 3", func() {

		rdepend := `
	app-crypt/sbsigntools
	x11-themes/sabayon-artwork-grub
	sys-boot/os-prober
	app-arch/xz-utils
	>=sys-libs/ncurses-5.2-r5:0=
	mount? ( sys-fs/fuse =sys-apps/pmount-0.9.99_alpha-r5:= )
`
		gr, err := ParseDependencies(rdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(6))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sbsigntools",
						Category: "app-crypt",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sabayon-artwork-grub",
						Category: "x11-themes",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep5", func() {
			Expect(*gr.Dependencies[4]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:          "ncurses",
						Category:      "sys-libs",
						Slot:          "0=",
						Version:       "5.2",
						VersionSuffix: "-r5",
						Condition:     _gentoo.PkgCondGreaterEqual,
					},
				},
			))
		})

		It("Check dep6", func() {
			Expect(*gr.Dependencies[5]).Should(Equal(
				GentooDependency{
					Use:          "mount",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:     "fuse",
								Category: "sys-fs",
								Slot:     "0",
							},
						},
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:          "pmount",
								Category:      "sys-apps",
								Condition:     _gentoo.PkgCondEqual,
								Version:       "0.9.99",
								VersionSuffix: "_alpha-r5",
								Slot:          "=",
							},
						},
					},
					Dep: nil,
				},
			))
		})

	})

	Context("Parse Dependencies", func() {

		rdepend := `
	app-crypt/sbsigntools
	x11-themes/sabayon-artwork-grub
	sys-boot/os-prober
	app-arch/xz-utils
	>=sys-libs/ncurses-5.2-r5:0=
	!mount? ( sys-fs/fuse =sys-apps/pmount-0.9.99_alpha-r5:= )
`
		gr, err := ParseDependencies(rdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(6))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sbsigntools",
						Category: "app-crypt",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sabayon-artwork-grub",
						Category: "x11-themes",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep5", func() {
			Expect(*gr.Dependencies[4]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:          "ncurses",
						Category:      "sys-libs",
						Slot:          "0=",
						Version:       "5.2",
						VersionSuffix: "-r5",
						Condition:     _gentoo.PkgCondGreaterEqual,
					},
				},
			))
		})

		It("Check dep6", func() {
			Expect(*gr.Dependencies[5]).Should(Equal(
				GentooDependency{
					Use:          "mount",
					UseCondition: _gentoo.PkgCondNot,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:     "fuse",
								Category: "sys-fs",
								Slot:     "0",
							},
						},
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:          "pmount",
								Category:      "sys-apps",
								Condition:     _gentoo.PkgCondEqual,
								Version:       "0.9.99",
								VersionSuffix: "_alpha-r5",
								Slot:          "=",
							},
						},
					},
					Dep: nil,
				},
			))
		})

	})

	Context("Parse Dependencies 5", func() {

		rdepend := `
	app-crypt/sbsigntools
	>=sys-libs/ncurses-5.2-r5:0=
	mount? (
		sys-fs/fuse
		=sys-apps/pmount-0.9.99_alpha-r5:=
	)
`
		gr, err := ParseDependencies(rdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(3))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sbsigntools",
						Category: "app-crypt",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:          "ncurses",
						Category:      "sys-libs",
						Slot:          "0=",
						Version:       "5.2",
						VersionSuffix: "-r5",
						Condition:     _gentoo.PkgCondGreaterEqual,
					},
				},
			))
		})

		It("Check dep3", func() {
			Expect(*gr.Dependencies[2]).Should(Equal(
				GentooDependency{
					Use:          "mount",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:     "fuse",
								Category: "sys-fs",
								Slot:     "0",
							},
						},
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:          "pmount",
								Category:      "sys-apps",
								Condition:     _gentoo.PkgCondEqual,
								Version:       "0.9.99",
								VersionSuffix: "_alpha-r5",
								Slot:          "=",
							},
						},
					},
					Dep: nil,
				},
			))
		})

	})

	Context("Parse Dependencies 6", func() {

		rdepend := `
	app-crypt/sbsigntools
	>=sys-libs/ncurses-5.2-r5:0=
	mount? (
		sys-fs/fuse
		=sys-apps/pmount-0.9.99_alpha-r5:= )
`
		gr, err := ParseDependencies(rdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(3))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sbsigntools",
						Category: "app-crypt",
						Slot:     "0",
					},
				},
			))
		})
	})

	Context("Parse Dependencies 7", func() {

		rdepend := `
	app-crypt/sbsigntools
	>=sys-libs/ncurses-5.2-r5:0=
	mount? (
		sys-fs/fuse
		=sys-apps/pmount-0.9.99_alpha-r5:=
		ext2? (
			sys-fs/genext2fs
		)
	)
`
		gr, err := ParseDependencies(rdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(3))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "sbsigntools",
						Category: "app-crypt",
						Slot:     "0",
					},
				},
			))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:          "ncurses",
						Category:      "sys-libs",
						Slot:          "0=",
						Version:       "5.2",
						VersionSuffix: "-r5",
						Condition:     _gentoo.PkgCondGreaterEqual,
					},
				},
			))
		})

		It("Check dep3", func() {
			Expect(*gr.Dependencies[2]).Should(Equal(
				GentooDependency{
					Use:          "mount",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:     "fuse",
								Category: "sys-fs",
								Slot:     "0",
							},
						},
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:          "pmount",
								Category:      "sys-apps",
								Condition:     _gentoo.PkgCondEqual,
								Version:       "0.9.99",
								VersionSuffix: "_alpha-r5",
								Slot:          "=",
							},
						},
						&GentooDependency{
							Use:          "ext2",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps: []*GentooDependency{
								&GentooDependency{
									Use:          "",
									UseCondition: _gentoo.PkgCondInvalid,
									SubDeps:      make([]*GentooDependency, 0),
									Dep: &_gentoo.GentooPackage{
										Name:     "genext2fs",
										Category: "sys-fs",
										Slot:     "0",
									},
								},
							},
							Dep: nil,
						},
					},
				},
			))
		})

	})

})
