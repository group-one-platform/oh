package cmd

import (
	"fmt"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/ui"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var (
	resetImageId  int
	resetName     string
	resetPassword string
)

var validVpsActions = []string{
	string(api.VirtualServerSoftReboot),
	string(api.VirtualServerHardReboot),
	string(api.VirtualServerPowerOff),
	string(api.VirtualServerPowerOn),
	string(api.VirtualServerReset),
}
var validVpsActionSet = func() map[string]struct{} {
	m := make(map[string]struct{}, len(validVpsActions))
	for _, a := range validVpsActions {
		m[a] = struct{}{}
	}
	return m
}()

var vpsActionCmd = &cobra.Command{
	Use:               "execute <vps-id> <action>",
	Short:             "Execute an action on a VPS (soft-reboot, hard-reboot, power-off, power-on, reset)",
	Long:              `Run one of the VirtualServerAction (soft-reboot, hard-reboot, power-off, power-on, reset) against a given VPS ID.`,
	SilenceUsage:      true,
	Args:              validateVpsExecuteArgs,
	ValidArgsFunction: completeVpsExecuteArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		vpsId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid VPS Id %q: %w", args[0], err)
		}
		action := api.VirtualServerAction(args[1])

		var request any

		if action == api.VirtualServerReset {
			err = validateResetCommand()
			if err != nil {
				return err
			}

			request = api.ResetCloudServerRequest{
				ImageId:  resetImageId,
				Name:     resetName,
				Password: resetPassword,
			}
		}

		resp, err := api.ExecuteVirtualServerAction(vpsId, action, request)
		if err != nil {
			return err
		}
		if printed, err := PrintJSON(resp, cmd); printed {
			return err
		}

		return ui.RenderForm(resp, actionResponseColumns()...)
	},
}

func validateResetCommand() error {
	if resetImageId == 0 {
		return fmt.Errorf("please supply the image id")
	}
	if resetName == "" {
		return fmt.Errorf("please supply the name for the VPS")
	}
	if resetPassword == "" {
		return fmt.Errorf("please supply the password for the VPS")
	}
	return nil
}

func validateVpsExecuteArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("vps-id is required")
	}

	if len(args) < 2 {
		return fmt.Errorf(
			"action is required; must be one of [%s]",
			strings.Join(validVpsActions, ", "),
		)
	}

	if _, ok := validVpsActionSet[args[1]]; !ok {
		return fmt.Errorf(
			"invalid action %q; must be one of [%s]",
			args[1],
			strings.Join(validVpsActions, ", "),
		)
	}
	return nil
}

func completeVpsExecuteArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	// Complete VPS ids by dynamic API lookup
	case 0:
		return completeVpsIds(cmd, args, toComplete)

	case 1:
		// Validate vps actions
		var comps []string
		for _, a := range validVpsActions {
			if strings.HasPrefix(a, toComplete) {
				comps = append(comps, a)
			}
		}
		return comps, cobra.ShellCompDirectiveNoFileComp

	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

func actionResponseColumns() []ui.TableColumn[api.VirtualServerActionResponse] {
	return []ui.TableColumn[api.VirtualServerActionResponse]{
		ui.Column("Id", 11, func(i api.VirtualServerActionResponse) int { return i.Id }),
		ui.Column("Message", 35, func(i api.VirtualServerActionResponse) string { return i.Message }),
	}
}

func init() {
	vpsActionCmd.Flags().IntVarP(&resetImageId, "image-id", "i", 0, "ID of the image to reset")
	vpsActionCmd.Flags().StringVarP(&resetName, "name", "n", "", "Name of the virtual server")
	vpsActionCmd.Flags().StringVarP(&resetPassword, "password", "p", "", "Password of the virtual server")

	vpsCmd.AddCommand(vpsActionCmd)
}
