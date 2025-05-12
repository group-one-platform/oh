package cmd

import (
	"fmt"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/cache"
	"github.com/edvin/oh/ui"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var (
	flavourId int
)
var vpsFlavourCmd = &cobra.Command{
	Use:   "flavour",
	Short: "Manage VPS Flavours",
}

var listFlavoursCmd = &cobra.Command{
	Use:   "list [server-id]",
	Short: "List Possible Flavours for VPS",
	Long: `This endpoint can be used only if it is included in your subscription (your support representative can provide more information regarding how to include it with your subscription).

Returns a list of compatible flavours for a specific VPS. These can be used to change the server’s configuration.`,
	Args:              validateSingleVpsIdArg,
	ValidArgsFunction: completeVpsIds,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid server Id %q: %w", args[0], err)
		}

		flavours, err := cache.Call(cache.KeyFlavours.WithArg(serverId), cache.DefaultTTL, func() ([]api.CloudServerFlavour, error) {
			return api.ListVpsFlavours(serverId)
		})
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(flavours, cmd); printed {
			return err
		}

		return ui.RenderTable(flavours, flavourColumns()...)
	},
}

var changeFlavourCmd = &cobra.Command{
	Use:   "set [server-id]",
	Short: "Change Flavour of a VPS",
	Long: `This endpoint can be used only if it is included in your subscription (your support representative can provide more information regarding how to include it with your subscription).

Changes the current flavour (configuration) of a specified VPS instance.`,
	Args:              validateSingleVpsIdArg,
	ValidArgsFunction: completeVpsIds,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid server Id %q: %w", args[0], err)
		}

		response, err := api.ChangeVpsFlavour(serverId, flavourId)
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(response, cmd); printed {
			return err
		}

		return ui.RenderForm(response, changeFlavourColumns()...)
	},
}

func init() {
	changeFlavourCmd.Flags().IntVarP(&flavourId, "flavour", "f", 0, "The new Flavour ID")
	changeFlavourCmd.RegisterFlagCompletionFunc("flavour", completeFlavoursForServer)
	vpsFlavourCmd.AddCommand(listFlavoursCmd, changeFlavourCmd)
	vpsCmd.AddCommand(vpsFlavourCmd)
}

func flavourColumns() []ui.TableColumn[api.CloudServerFlavour] {
	return []ui.TableColumn[api.CloudServerFlavour]{
		ui.Column("Id", 11, func(f api.CloudServerFlavour) int { return f.Id }),
		ui.Column("Name", 35, func(f api.CloudServerFlavour) string { return f.Name }),
		ui.Column("Cores", 10, func(f api.CloudServerFlavour) int { return f.Cores }),
		ui.Column("Ram Size", 10, func(f api.CloudServerFlavour) int { return f.RamSize }),
		ui.Column("Storage Type", 15, func(f api.CloudServerFlavour) string { return f.StorageType }),
		ui.Column("Storage Size", 15, func(f api.CloudServerFlavour) int { return f.StorageSize }),
	}
}

func changeFlavourColumns() []ui.TableColumn[api.ChangeFlavourResponse] {
	return []ui.TableColumn[api.ChangeFlavourResponse]{
		ui.Column("ServerId", 11, func(f api.ChangeFlavourResponse) int { return f.ServerId }),
		ui.Column("FlavourId", 11, func(f api.ChangeFlavourResponse) int { return f.FlavourId }),
		ui.Column("Message", 40, func(f api.ChangeFlavourResponse) string { return f.Message }),
	}
}

func completeFlavoursForServer(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) < 1 {
		// no server ID yet, bail out
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	serverId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	flavours, err := cache.Call(cache.KeyFlavours, cache.DefaultTTL, func() ([]api.CloudServerFlavour, error) {
		return api.ListVpsFlavours(serverId)
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var comps []string
	for _, f := range flavours {
		id := strconv.Itoa(f.Id)
		if strings.HasPrefix(id, toComplete) {
			comps = append(comps, fmt.Sprintf("%s\t%s", id, f.Name))
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}
