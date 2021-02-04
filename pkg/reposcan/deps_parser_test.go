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
		gr, err := ParseDependenciesMultiline(rdepend)
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
		gr, err := ParseDependenciesMultiline(rdepend)
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
		gr, err := ParseDependenciesMultiline(rdepend)
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
		gr, err := ParseDependenciesMultiline(rdepend)
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
		gr, err := ParseDependenciesMultiline(rdepend)
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
		gr, err := ParseDependenciesMultiline(rdepend)
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
		gr, err := ParseDependenciesMultiline(rdepend)
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

	Context("Parse Dependencies 8", func() {

		rdepend := `sys-libs/zlib  virtual/libintl`
		gr, err := ParseDependencies(rdepend)

		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(2))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "zlib",
						Category: "sys-libs",
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
						Name:     "libintl",
						Category: "virtual",
						Slot:     "0",
					},
				},
			))
		})
	})

	Context("Parse Dependencies 9", func() {

		rdepend := `sys-libs/zlib nls? ( virtual/libintl )`
		gr, err := ParseDependencies(rdepend)

		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(2))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "zlib",
						Category: "sys-libs",
						Slot:     "0",
					},
				},
			))
		})
		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "nls",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:     "libintl",
								Category: "virtual",
								Slot:     "0",
							},
						},
					},
				},
			))
		})
	})

	Context("Parse Dependencies 10", func() {

		rdepend := `sys-libs/zlib nls? ( virtual/libintl ) virtual/libiconv >=dev-libs/gmp-4.3.2:0= >=dev-libs/mpfr-2.4.2:0= >=dev-libs/mpc-0.8.1:0= objc-gc? ( >=dev-libs/boehm-gc-7.4.2 ) graphite? ( >=dev-libs/isl-0.14:0= )`
		gr, err := ParseDependencies(rdepend)

		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(8))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "zlib",
						Category: "sys-libs",
						Slot:     "0",
					},
				},
			))
		})
		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "nls",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:     "libintl",
								Category: "virtual",
								Slot:     "0",
							},
						},
					},
				},
			))
		})
		It("Check dep3", func() {
			Expect(*gr.Dependencies[2]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:     "libiconv",
						Category: "virtual",
						Slot:     "0",
					},
				},
			))
		})
		It("Check dep4", func() {
			Expect(*gr.Dependencies[3]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:      "gmp",
						Category:  "dev-libs",
						Version:   "4.3.2",
						Condition: _gentoo.PkgCondGreaterEqual,
						Slot:      "0=",
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
						Name:      "mpfr",
						Category:  "dev-libs",
						Version:   "2.4.2",
						Condition: _gentoo.PkgCondGreaterEqual,
						Slot:      "0=",
					},
				},
			))
		})

		It("Check dep6", func() {
			Expect(*gr.Dependencies[5]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:      "mpc",
						Category:  "dev-libs",
						Version:   "0.8.1",
						Condition: _gentoo.PkgCondGreaterEqual,
						Slot:      "0=",
					},
				},
			))
		})

		It("Check dep7", func() {
			Expect(*gr.Dependencies[6]).Should(Equal(
				GentooDependency{
					Use:          "objc-gc",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:      "boehm-gc",
								Category:  "dev-libs",
								Version:   "7.4.2",
								Condition: _gentoo.PkgCondGreaterEqual,
								Slot:      "0",
							},
						},
					},
				},
			))
		})

		It("Check dep8", func() {
			Expect(*gr.Dependencies[7]).Should(Equal(
				GentooDependency{
					Use:          "graphite",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:      "isl",
								Category:  "dev-libs",
								Version:   "0.14",
								Condition: _gentoo.PkgCondGreaterEqual,
								Slot:      "0=",
							},
						},
					},
				},
			))
		})

	})

	//"BDEPEND": "minizip? ( || ( >=sys-devel/automake-1.16.1:1.16 >=sys-devel/automake-1.15.1:1.15 ) >=sys-devel/autoconf-2.69 >=sys-devel/libtool-2.4 ) >=app-portage/elt-patches-20170815",

	Context("Parse Dependencies 11", func() {

		rdepend := `minizip? ( || ( >=sys-devel/automake-1.16.1:1.16 >=sys-devel/automake-1.15.1:1.15 ) >=sys-devel/autoconf-2.69 >=sys-devel/libtool-2.4 ) >=app-portage/elt-patches-20170815`

		gr, err := ParseDependencies(rdepend)

		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(2))
		})

		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(
				GentooDependency{
					Use:          "minizip",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							DepInOr:      true,
							SubDeps: []*GentooDependency{
								&GentooDependency{
									Use:          "",
									UseCondition: _gentoo.PkgCondInvalid,
									SubDeps:      make([]*GentooDependency, 0),
									Dep: &_gentoo.GentooPackage{
										Name:      "automake",
										Category:  "sys-devel",
										Version:   "1.16.1",
										Condition: _gentoo.PkgCondGreaterEqual,
										Slot:      "1.16",
									},
								},

								&GentooDependency{
									Use:          "",
									UseCondition: _gentoo.PkgCondInvalid,
									SubDeps:      make([]*GentooDependency, 0),
									Dep: &_gentoo.GentooPackage{
										Name:      "automake",
										Category:  "sys-devel",
										Version:   "1.15.1",
										Condition: _gentoo.PkgCondGreaterEqual,
										Slot:      "1.15",
									},
								},
							},
						},

						// sys-devel/autoconf
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:      "autoconf",
								Category:  "sys-devel",
								Version:   "2.69",
								Condition: _gentoo.PkgCondGreaterEqual,
								Slot:      "0",
							},
						},
						// sys-devel/libtool
						&GentooDependency{
							Use:          "",
							UseCondition: _gentoo.PkgCondInvalid,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:      "libtool",
								Category:  "sys-devel",
								Version:   "2.4",
								Condition: _gentoo.PkgCondGreaterEqual,
								Slot:      "0",
							},
						},
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
						Name:      "elt-patches",
						Category:  "app-portage",
						Version:   "20170815",
						Condition: _gentoo.PkgCondGreaterEqual,
						Slot:      "0",
					},
				},
			))
		})

	})

})
