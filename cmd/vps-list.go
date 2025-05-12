package cmd

import (
	"fmt"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/ui"
	"github.com/spf13/cobra"
	"strconv"
)

var listVpsCmd = &cobra.Command{
	Use:          "list",
	Short:        "List Available VPS instances",
	Long:         `Retrieves a list of all VPS instances`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		servers, err := api.ListCloudServers()
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(servers, cmd); printed {
			return err
		}

		return ui.RenderTable(servers, serverColumns()...)
	},
}

var getVpsCmd = &cobra.Command{
	Use:               "get [server-id]",
	Short:             "Get Image Details",
	Long:              `Fetches the detailed information of the specified image.`,
	ValidArgsFunction: completeVpsIds,
	Args:              validateSingleVpsIdArg,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid server Id %q: %w", args[0], err)
		}

		image, err := api.GetVirtualServer(serverId)
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(image, cmd); printed {
			return err
		}

		return ui.RenderForm(image, serverColumns()...)
	},
}

func init() {
	vpsCmd.AddCommand(listVpsCmd, getVpsCmd)
}

func serverColumns() []ui.TableColumn[api.CloudServer] {
	return []ui.TableColumn[api.CloudServer]{
		ui.Column("Id", 11, func(i api.CloudServer) int { return i.Id }),
		ui.Column("Name", 35, func(i api.CloudServer) string { return i.Name }),
		ui.Column("IPv4", 16, func(i api.CloudServer) string { return i.IPv4 }),
		ui.Column("IPv6", 26, func(i api.CloudServer) string { return i.IPv6 }),
		ui.Column("Status", 12, func(i api.CloudServer) string { return i.Status }),
		ui.Column("Image #", 10, func(i api.CloudServer) int { return i.Image.Id }),
		ui.Column("OS Distro", 15, func(i api.CloudServer) string { return i.Image.OSDistro }),
		ui.Column("OS Version", 15, func(i api.CloudServer) string { return i.Image.OSVersion }),
		ui.Column("Image Release Date", 20, func(i api.CloudServer) api.Date { return i.Image.ReleaseDate }),
	}
}
