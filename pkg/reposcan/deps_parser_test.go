package reposcan

import (
	"fmt"

	_gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Context("Parse Dependencies 11", func() {
		bdepend := `minizip? ( || ( >=sys-devel/automake-1.16.1:1.16 >=sys-devel/automake-1.15.1:1.15 ) >=sys-devel/autoconf-2.69 >=sys-devel/libtool-2.4 ) >=app-portage/elt-patches-20170815`

		gr, err := ParseDependencies(bdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(2))
		})

		d1, e1 := NewGentooDependency(">=sys-devel/autoconf-2.69", "")
		d2, e2 := NewGentooDependency(">=sys-devel/libtool-2.4", "")
		d3, e3 := NewGentooDependency(">=app-portage/elt-patches-20170815", "")
		d5, e4 := NewGentooDependency(">=sys-devel/automake-1.16.1:1.16", "")
		d6, e5 := NewGentooDependency(">=sys-devel/automake-1.15.1:1.15", "")
		It("Check error d1", func() {
			Expect(e1).Should(BeNil())
			Expect(e2).Should(BeNil())
			Expect(e3).Should(BeNil())
			Expect(e4).Should(BeNil())
			Expect(e5).Should(BeNil())
		})
		aggr1, err_aggr := NewGentooDependencyWithSubdeps("", "", []*GentooDependency{d5, d6})
		It("Check error dor", func() {
			Expect(err_aggr).Should(BeNil())
		})
		dor, err_dor := NewGentooDependencyWithSubdeps("", "", []*GentooDependency{aggr1})
		It("Check error dor", func() {
			Expect(err_dor).Should(BeNil())
		})
		dor.DepInOr = true
		minizip, err := NewGentooDependencyWithSubdeps("", "minizip",
			[]*GentooDependency{dor, d1, d2})
		It("Check error minizip", func() {
			Expect(err).Should(BeNil())
		})
		It("Check dep1", func() {
			Expect(*gr.Dependencies[0]).Should(Equal(*minizip))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(*d3))
		})

	})
	Context("Parse Dependencies 12", func() {

		rdepend := `aqua? ( x11-libs/gtk+:2[aqua=,abi_x86_32(-)?,abi_x86_64(-)?,abi_x86_x32(-)?,abi_mips_n32(-)?,abi_mips_n64(-)?,abi_mips_o32(-)?,abi_s390_32(-)?,abi_s390_64(-)?] virtual/jpeg:0=[abi_x86_32(-)?,abi_x86_64(-)?,abi_x86_x32(-)?,abi_mips_n32(-)?,abi_mips_n64(-)?,abi_mips_o32(-)?,abi_s390_32(-)?,abi_s390_64(-)?] tiff? ( media-libs/tiff:0[abi_x86_32(-)?,abi_x86_64(-)?,abi_x86_x32(-)?,abi_mips_n32(-)?,abi_mips_n64(-)?,abi_mips_o32(-)?,abi_s390_32(-)?,abi_s390_64(-)?] ) )`
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
					Use:          "aqua",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							UseCondition: 0,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:      "gtk+",
								Category:  "x11-libs",
								Slot:      "2",
								Condition: 0,
							},
							DepInOr: false,
						},

						&GentooDependency{
							UseCondition: 0,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:      "jpeg",
								Category:  "virtual",
								Slot:      "0=",
								Condition: 0,
							},
							DepInOr: false,
						},
					},

					Dep: nil,
				},
			))
		})

		It("Check dep2", func() {
			Expect(*gr.Dependencies[1]).Should(Equal(
				GentooDependency{
					Use:          "tiff",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps: []*GentooDependency{
						&GentooDependency{
							UseCondition: 0,
							SubDeps:      make([]*GentooDependency, 0),
							Dep: &_gentoo.GentooPackage{
								Name:      "tiff",
								Category:  "media-libs",
								Slot:      "0",
								Condition: 0,
							},
							DepInOr: false,
						},
					},

					Dep: nil,
				},
			))
		})
	})
	Context("Parse Dependencies 14", func() {

		bdepend := `gnutls? ( net-libs/gnutls[tools] ) !gnutls? ( !libressl? ( dev-libs/openssl:0= ) libressl? ( dev-libs/libressl:0= ) ) >=net-libs/courier-authlib-0.66.4 >=net-libs/courier-unicode-2 >=net-mail/mailbase-0.00-r8 net-dns/libidn:= berkdb? ( sys-libs/db:= ) fam? ( virtual/fam ) gdbm? ( >=sys-libs/gdbm-1.8.0 ) selinux? ( sec-policy/selinux-courier ) !mail-mta/courier !net-mail/bincimap !net-mail/cyrus-imapd !net-mail/uw-imap`

		gr, err := ParseDependencies(bdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(16))
		})

		It("Check dep15", func() {
			Expect(*gr.Dependencies[15]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:      "uw-imap",
						Category:  "net-mail",
						Slot:      "0",
						Condition: _gentoo.PkgCondNot,
					},
					DepInOr: false,
				},
			))
		})
		fmt.Println(gr)
	})

	Context("Parse Dependencies 15", func() {

		bdepend := `|| ( dev-lang/python:3.9 dev-lang/python:3.8 dev-lang/python:3.7 dev-lang/python:3.6 ) opencl? ( >=sys-devel/gcc-4.6 ) sys-devel/bison sys-devel/flex sys-devel/gettext virtual/pkgconfig || ( ( dev-lang/python:3.9 >=dev-python/mako-0.8.0[python_targets_python3_9(-),python_single_target_python3_9(+)] ) ( dev-lang/python:3.8 >=dev-python/mako-0.8.0[python_targets_python3_8(-),python_single_target_python3_8(+)] ) ( dev-lang/python:3.7 >=dev-python/mako-0.8.0[python_targets_python3_7(-),python_single_target_python3_7(+)] ) ( dev-lang/python:3.6 >=dev-python/mako-0.8.0[python_targets_python3_6(-),python_single_target_python3_6(+)] ) ) >=dev-util/meson-0.54.0 >=dev-util/ninja-1.8.2`

		gr, err := ParseDependencies(bdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(9))
		})

		It("Check dep15", func() {
			Expect(*gr.Dependencies[8]).Should(Equal(
				GentooDependency{
					Use:          "",
					UseCondition: _gentoo.PkgCondInvalid,
					SubDeps:      make([]*GentooDependency, 0),
					Dep: &_gentoo.GentooPackage{
						Name:      "ninja",
						Category:  "dev-util",
						Version:   "1.8.2",
						Slot:      "0",
						Condition: _gentoo.PkgCondGreaterEqual,
					},
					DepInOr: false,
				},
			))
		})
		fmt.Println(gr)
	})
	Context("Parse Dependencies 13", func() {

		bdepend := `virtual/pkgconfig test? ( || ( dev-lang/python:3.8 dev-lang/python:3.7 dev-lang/python:3.6 ) || ( ( dev-lang/python:3.8 dev-python/pytest[python_targets_python3_8(-),python_single_target_python3_8(+)] ) ( dev-lang/python:3.7 dev-python/pytest[python_targets_python3_7(-),python_single_target_python3_7(+)] ) ( dev-lang/python:3.6 dev-python/pytest[python_targets_python3_6(-),python_single_target_python3_6(+)] ) ) ) >=dev-util/meson-0.54.0 >=dev-util/ninja-1.8.2 virtual/pkgconfig`

		gr, err := ParseDependencies(bdepend)
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check gr", func() {
			Expect(gr).ShouldNot(BeNil())
		})

		It("Check deps #", func() {
			Expect(len(gr.Dependencies)).Should(Equal(5))
		})

		fmt.Println(gr)
	})
})
