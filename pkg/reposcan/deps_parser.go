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

Based on old code of luet simpleparser.
*/

package reposcan

import (
	"fmt"
	"strings"

	_gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
)

type GentooDependency struct {
	Use          string
	UseCondition _gentoo.PackageCond
	SubDeps      []*GentooDependency
	Dep          *_gentoo.GentooPackage
}

type EbuildDependencies struct {
	Dependencies []*GentooDependency
}

func NewGentooDependency(pkg, use string) (*GentooDependency, error) {
	var err error
	ans := &GentooDependency{
		Use:     use,
		SubDeps: make([]*GentooDependency, 0),
	}

	if strings.HasPrefix(use, "!") {
		ans.Use = ans.Use[1:]
		ans.UseCondition = _gentoo.PkgCondNot
	}

	if pkg != "" {
		ans.Dep, err = _gentoo.ParsePackageStr(pkg)
		if err != nil {
			return nil, err
		}

		// TODO: Fix this on parsing phase for handle correctly ${PV}
		if strings.HasSuffix(ans.Dep.Name, "-") {
			ans.Dep.Name = ans.Dep.Name[:len(ans.Dep.Name)-1]
		}

	}

	return ans, nil
}

func (d *GentooDependency) String() string {
	if d.Dep != nil {
		return fmt.Sprintf("%s", d.Dep)
	} else {
		return fmt.Sprintf("%s %d %s", d.Use, d.UseCondition, d.SubDeps)
	}
}

func (d *GentooDependency) GetDepsList() []*GentooDependency {
	ans := make([]*GentooDependency, 0)

	if len(d.SubDeps) > 0 {
		for _, d2 := range d.SubDeps {
			list := d2.GetDepsList()
			ans = append(ans, list...)
		}
	}

	if d.Dep != nil {
		ans = append(ans, d)
	}

	return ans
}

func (d *GentooDependency) AddSubDependency(pkg, use string) (*GentooDependency, error) {
	ans, err := NewGentooDependency(pkg, use)
	if err != nil {
		return nil, err
	}

	d.SubDeps = append(d.SubDeps, ans)

	return ans, nil
}

func (r *EbuildDependencies) GetDependencies() []*GentooDependency {
	ans := make([]*GentooDependency, 0)

	for _, d := range r.Dependencies {
		list := d.GetDepsList()
		ans = append(ans, list...)
	}

	// the same dependency could be available in multiple use flags.
	// It's needed avoid duplicate.
	m := make(map[string]*GentooDependency, 0)

	for _, p := range ans {
		m[p.String()] = p
	}

	ans = make([]*GentooDependency, 0)
	for _, p := range m {
		ans = append(ans, p)
	}

	return ans
}

func ParseDependencies(rdepend string) (*EbuildDependencies, error) {
	var lastdep []*GentooDependency = make([]*GentooDependency, 0)
	var pendingDep = false
	var orDep = false
	var dep *GentooDependency
	var err error

	ans := &EbuildDependencies{
		Dependencies: make([]*GentooDependency, 0),
	}

	if rdepend != "" {
		rdepends := strings.Split(rdepend, "\n")
		for _, rr := range rdepends {
			rr = strings.TrimSpace(rr)
			if rr == "" {
				continue
			}

			if strings.HasPrefix(rr, "|| (") {
				orDep = true
				continue
			}

			if orDep {
				rr = strings.TrimSpace(rr)
				if rr == ")" {
					orDep = false
				}
				continue
			}

			if strings.Index(rr, "?") > 0 {
				// use flag present

				if pendingDep {
					dep, err = lastdep[len(lastdep)-1].AddSubDependency("", rr[:strings.Index(rr, "?")])
					if err != nil {
						// Debug
						fmt.Println("Ignoring subdependency ", rr[:strings.Index(rr, "?")])
					}
				} else {
					dep, err = NewGentooDependency("", rr[:strings.Index(rr, "?")])
					if err != nil {
						// Debug
						fmt.Println("Ignoring dep", rr)
					} else {
						ans.Dependencies = append(ans.Dependencies, dep)
					}
				}

				if strings.Index(rr, ")") < 0 {
					pendingDep = true
					lastdep = append(lastdep, dep)
				}

				if strings.Index(rr, "|| (") >= 0 {
					// Ignore dep in or
					continue
				}

				fields := strings.Split(rr[strings.Index(rr, "?")+1:], " ")
				for _, f := range fields {
					f = strings.TrimSpace(f)
					if f == ")" || f == "(" || f == "" {
						continue
					}

					_, err = dep.AddSubDependency(f, "")
					if err != nil {
						// Debug
						fmt.Println("Ignoring subdependency ", f)
					}
				}

			} else if pendingDep {
				fields := strings.Split(rr, " ")
				for _, f := range fields {
					f = strings.TrimSpace(f)
					if f == ")" || f == "(" || f == "" {
						continue
					}
					_, err = lastdep[len(lastdep)-1].AddSubDependency(f, "")
					if err != nil {
						return nil, err
					}
				}

				if strings.Index(rr, ")") >= 0 {
					lastdep = lastdep[:len(lastdep)-1]
					if len(lastdep) == 0 {
						pendingDep = false
					}
				}

			} else {
				rr = strings.TrimSpace(rr)
				// Check if there multiple deps in single row

				fields := strings.Split(rr, " ")
				if len(fields) > 1 {
					for _, rrr := range fields {
						rrr = strings.TrimSpace(rrr)
						if rrr == "" {
							continue
						}
						dep, err := NewGentooDependency(rrr, "")
						if err != nil {
							// Debug
							fmt.Println("Ignoring dep", rr)
						} else {
							ans.Dependencies = append(ans.Dependencies, dep)
						}
					}
				} else {
					dep, err := NewGentooDependency(rr, "")
					if err != nil {
						// Debug
						fmt.Println("Ignoring dep", rr)
					} else {
						ans.Dependencies = append(ans.Dependencies, dep)
					}
				}
			}

		}

	}

	return ans, nil
}
