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
package specs

func NewPortageResolverOpts() PortageResolverOpts {
	return PortageResolverOpts{
		EnableUseFlags:   []string{},
		DisabledUseFlags: []string{},
	}
}

func (o *PortageResolverOpts) IsAdmitUseFlag(u string) bool {
	ans := true
	if len(o.EnableUseFlags) > 0 {
		for _, ue := range o.EnableUseFlags {
			if ue == u {
				return true
			}
		}

		return false
	}

	if len(o.DisabledUseFlags) > 0 {
		for _, ud := range o.DisabledUseFlags {
			if ud == u {
				ans = false
				break
			}
		}
	}

	return ans
}
