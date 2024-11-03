/*
Copyright © 2021-2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd_portage

import (
	"fmt"

	. "github.com/geaaru/luet/pkg/config"
	. "github.com/geaaru/luet/pkg/logger"
	"github.com/macaroni-os/anise-portage-converter/pkg/portage"

	"github.com/spf13/cobra"
)

func NewPortageShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"show"},
		Short:   "Show data from Portage.",
		Run: func(cmd *cobra.Command, args []string) {
			debug, _ := cmd.Flags().GetBool("debug")
			world, _ := cmd.Flags().GetBool("world")

			if debug {
				LuetCfg.GetGeneral().Debug = debug
			}

			if world {
				packages, err := portage.GetWorldPackages("")
				if err != nil {
					Fatal(err)
				}

				for _, p := range packages {
					fmt.Println(p)
				}
			}

		},
	}

	cmd.Flags().Bool("world", false, "Purge Portage edb cache.")

	return cmd
}
