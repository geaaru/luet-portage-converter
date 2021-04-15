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
	"strings"
	"testing"

	. "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func initConfig() error {
	LuetCfg.Viper.SetEnvPrefix("LUET")
	LuetCfg.Viper.AutomaticEnv() // read in environment variables that match

	// Create EnvKey Replacer for handle complex structure
	replacer := strings.NewReplacer(".", "__")
	LuetCfg.Viper.SetEnvKeyReplacer(replacer)
	LuetCfg.Viper.SetTypeByDefaultValue(true)

	err := LuetCfg.Viper.Unmarshal(&LuetCfg)
	if err != nil {
		return err
	}

	LuetCfg.GetGeneral().Debug = true

	InitAurora()
	NewSpinner()

	return nil
}

func TestSolver(t *testing.T) {
	RegisterFailHandler(Fail)
	initConfig()
	RunSpecs(t, "Converter Suite")
}
