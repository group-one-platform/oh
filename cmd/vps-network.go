package cmd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/cache"
	"github.com/edvin/oh/ui"
	"github.com/spf13/cobra"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	detachNetId string
	attachNetId string
	attachIPv4  string
	attachIPv6  string
)

var vpsNetworkCommand = &cobra.Command{
	Use:   "network",
	Short: "Manage VPS Networks",
}

var listAvailableNetworksCmd = &cobra.Command{
	Use:               "list-available",
	Short:             "List Available Virtual Networks",
	SilenceUsage:      true,
	ValidArgsFunction: NoArgs,
	Long: `This endpoint can be used only if it is included in your subscription (your support representative can provide more information regarding how to include it with your subscription).

Lists all available Virtual Networks as array each containing array(s) of subnet related information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		networks, err := cache.Call(cache.KeyVirtualNetworks, cache.DefaultTTL, func() ([]api.VirtualNetwork, error) {
			return api.ListVirtualNetworks()
		})
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(networks, cmd); printed {
			return err
		}

		return ui.RenderTable(networks, networkColumns()...)
	},
}

var listAttachedNetworksCmd = &cobra.Command{
	Use:               "list [server-id]",
	Short:             "List Attached Virtual Networks on VPS",
	SilenceUsage:      true,
	ValidArgsFunction: completeVpsIds,
	Args: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			return fmt.Errorf("you must specify the VPS ID, e.g.:\n  oh vps network list 42")
		case 1:
			return nil
		default:
			return fmt.Errorf("only one positional argument expected (the VPS ID), got %d", len(args))
		}
	},

	Long: `This endpoint can be used only if it is included in your subscription (your support representative can provide more information regarding how to include it with your subscription).

Return list of all attached networks on specified VPS.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid server Id %q: %w", args[0], err)
		}

		networks, err := cache.Call(cache.KeyFlavours.WithArg(serverId), time.Minute, func() ([]api.AttachedNetwork, error) {
			return api.ListAttachedVirtualNetworks(serverId)
		})
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(networks, cmd); printed {
			return err
		}

		return ui.RenderTable(networks, attachedNetworkColumns()...)
	},
}

var detachNetworksCmd = &cobra.Command{
	Use:               "detach [server-id]",
	Short:             "Detach virtual network from server instance",
	SilenceUsage:      true,
	Args:              validateSingleVpsIdArg,
	ValidArgsFunction: completeVpsIds,
	Long:              `Detach virtual network from server instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid server Id %q: %w", args[0], err)
		}

		response, err := api.DetachVirtualNetwork(serverId, detachNetId)
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(response, cmd); printed {
			return err
		}

		return ui.RenderForm(response, detachNetworkResponseColumns()...)
	},
}

var attachNetworksCmd = &cobra.Command{
	Use:               "attach [server-id]",
	Short:             "Attach virtual network to server instance",
	SilenceUsage:      true,
	Args:              validateSingleVpsIdArg,
	ValidArgsFunction: completeVpsIds,
	Long:              `Attach virtual network to server instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid server Id %q: %w", args[0], err)
		}

		response, err := api.AttachVirtualNetwork(serverId, attachNetId, attachIPv4, attachIPv6)
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(response, cmd); printed {
			return err
		}

		return ui.RenderForm(response, attachNetworkResponseColumns()...)
	},
}

func init() {
	detachNetworksCmd.Flags().StringVarP(&detachNetId, "network-id", "n", "", "Network Id to detach")
	detachNetworksCmd.RegisterFlagCompletionFunc("network-id", completeAttachedNetworkIdsForServer)

	attachNetworksCmd.Flags().StringVarP(&attachNetId, "network-id", "n", "", "Network Id to attach")
	attachNetworksCmd.RegisterFlagCompletionFunc("network-id", completeAvailableNetworkIds)

	attachNetworksCmd.Flags().StringVarP(&attachIPv4, "ipv4", "4", "", "IPv4 address")
	attachNetworksCmd.RegisterFlagCompletionFunc("ipv4", completeAvailableIpv4Addresses)

	attachNetworksCmd.Flags().StringVarP(&attachIPv6, "ipv6", "6", "", "IPv6 address")

	vpsNetworkCommand.AddCommand(listAttachedNetworksCmd, detachNetworksCmd, attachNetworksCmd, listAvailableNetworksCmd)
	vpsCmd.AddCommand(vpsNetworkCommand)
}

func networkColumns() []ui.TableColumn[api.VirtualNetwork] {
	return []ui.TableColumn[api.VirtualNetwork]{
		ui.Column("Id", 40, func(i api.VirtualNetwork) string { return i.Id }),
		ui.Column("Name", 30, func(i api.VirtualNetwork) string { return i.Name }),
		ui.Column("Subnets", 50, func(i api.VirtualNetwork) string { return "[see json output]" }),
	}
}

func attachedNetworkColumns() []ui.TableColumn[api.AttachedNetwork] {
	return []ui.TableColumn[api.AttachedNetwork]{
		ui.Column("Id", 40, func(i api.AttachedNetwork) string { return i.Id }),
		ui.Column("Name", 30, func(i api.AttachedNetwork) string { return i.Name }),
		ui.Column("IPv4", 20, func(i api.AttachedNetwork) string { return i.IPv4 }),
		ui.Column("IPv6", 26, func(i api.AttachedNetwork) string { return i.IPv6 }),
	}
}

func detachNetworkResponseColumns() []ui.TableColumn[api.DetachVirtualNetworkResponse] {
	return []ui.TableColumn[api.DetachVirtualNetworkResponse]{
		ui.Column("ServerId", 40, func(i api.DetachVirtualNetworkResponse) int { return i.ServerId }),
		ui.Column("Message", 50, func(i api.DetachVirtualNetworkResponse) string { return i.Message }),
	}
}

func attachNetworkResponseColumns() []ui.TableColumn[api.AttachVirtualNetworkResponse] {
	return []ui.TableColumn[api.AttachVirtualNetworkResponse]{
		ui.Column("ServerId", 40, func(i api.AttachVirtualNetworkResponse) int { return i.ServerId }),
		ui.Column("NetworkId", 40, func(i api.AttachVirtualNetworkResponse) string { return i.NetworkId }),
		ui.Column("Message", 50, func(i api.AttachVirtualNetworkResponse) string { return i.Message }),
	}
}

func completeAttachedNetworkIdsForServer(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) < 1 {
		// no server ID yet, bail out
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	serverId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	key := cache.KeyAttachedNetworks.WithArg(serverId)
	networks, err := cache.Call(key, time.Minute, func() ([]api.AttachedNetwork, error) {
		return api.ListAttachedVirtualNetworks(serverId)
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var comps []string
	for _, n := range networks {
		if strings.HasPrefix(n.Id, toComplete) {
			comps = append(comps, fmt.Sprintf("%s\t%s", n.Id, n.Name))
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func completeAvailableNetworkIds(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	networks, err := cache.Call(cache.KeyVirtualNetworks, time.Minute, func() ([]api.VirtualNetwork, error) {
		return api.ListVirtualNetworks()
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var comps []string
	for _, n := range networks {
		if strings.HasPrefix(n.Id, toComplete) {
			comps = append(comps, fmt.Sprintf("%s\t%s", n.Id, n.Name))
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func completeAvailableIpv4Addresses(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	netID, err := cmd.Flags().GetString("network-id")
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	if netID == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	networks, err := cache.Call(cache.KeyVirtualNetworks, time.Minute, func() ([]api.VirtualNetwork, error) {
		return api.ListVirtualNetworks()
	})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	for _, n := range networks {
		if n.Id == netID {
			var comps []string

			for _, subnet := range n.Subnets {
				if subnet.IpVersion == 4 {
					for _, pool := range subnet.AllocationPools {
						ips, err := ListIPsInRange(pool.Start, pool.End)
						if err == nil {
							for _, ip := range ips {
								if strings.HasPrefix(ip, toComplete) {
									comps = append(comps, ip)
								}
							}
						}
					}
				}
			}

			if len(comps) > 0 {
				return comps, cobra.ShellCompDirectiveNoFileComp
			}
		}
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func ListIPsInRange(start, end string) ([]string, error) {
	sIP := net.ParseIP(start).To4()
	eIP := net.ParseIP(end).To4()
	if sIP == nil || eIP == nil {
		return nil, errors.New("start or end is not a valid IPv4 address")
	}
	s := binary.BigEndian.Uint32(sIP)
	e := binary.BigEndian.Uint32(eIP)

	if s > e {
		return nil, fmt.Errorf("start IP %s is greater than end IP %s", start, end)
	}

	size := e - s + 1
	ips := make([]string, 0, size)

	for x := s; x <= e; x++ {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], x)
		ips = append(ips, net.IP(buf[:]).String())
	}

	return ips, nil
}
