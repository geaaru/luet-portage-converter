/*
Copyright © 2021-2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd_portage

import (
	"github.com/macaroni-os/anise-portage-converter/pkg/portage"

	. "github.com/geaaru/luet/pkg/config"
	. "github.com/geaaru/luet/pkg/logger"

	"github.com/spf13/cobra"
)

func NewPortageCacheCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cache",
		Aliases: []string{"pc", "c"},
		Short:   "Purge Portage cache and optional packages from world file.",
		Run: func(cmd *cobra.Command, args []string) {
			debug, _ := cmd.Flags().GetBool("debug")
			purge, _ := cmd.Flags().GetBool("purge")
			pkgs2remove, _ := cmd.Flags().GetStringArray("pkg")

			jobDone := false

			if debug {
				LuetCfg.GetGeneral().Debug = debug
			}

			if purge {
				err := portage.CleanEdbCache("")
				if err != nil {
					Fatal(err)
				}
				jobDone = true
			}

			if len(pkgs2remove) > 0 {
				// Read the existing world file.
				packages, err := portage.GetWorldPackages("")
				if err != nil {
					Fatal(err)
				}

				packages, err = portage.GetWorldPackagesCleaned(
					packages, pkgs2remove,
				)
				if err != nil {
					Fatal(err)
				}

				err = portage.WriteWorldPackages("", packages)
				if err != nil {
					Fatal(err)
				}

				jobDone = true
			}

			if jobDone {
				Info(":panda: All done!")
			} else {
				Warning(":parachute: Nothing done.")
			}
		},
	}

	cmd.Flags().Bool("purge", false, "Purge Portage edb cache.")
	cmd.Flags().StringArrayP("pkg", "p", []string{},
		"Define the package to analyze.")

	return cmd
}
