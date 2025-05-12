package cmd

import (
	"fmt"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/cache"
	"github.com/edvin/oh/ui"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

var vpsImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Manage VPS Images",
}

var listVpsImagesCmd = &cobra.Command{
	Use:               "list",
	Short:             "List Available Images",
	SilenceUsage:      true,
	Long:              `Returns an array of all image available to use when creating or redeploying a VPS.`,
	ValidArgsFunction: NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		images, err := cache.Call(cache.KeyVpsImages, cache.DefaultTTL, func() ([]api.CloudServerImage, error) {
			return api.ListVpsImages()
		})
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(images, cmd); printed {
			return err
		}

		return ui.RenderTable(images, imageColumns()...)
	},
}

var getVpsImageCmd = &cobra.Command{
	Use:               "get [id]",
	Short:             "Get Image Details",
	Args:              validateSingleVpsIdArg,
	ValidArgsFunction: completeVpsImageIds,
	SilenceUsage:      true,
	Long:              `Fetches the detailed information of the specified image.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid image Id %q: %w", args[0], err)
		}

		image, err := cache.Call(cache.KeyVpsImages.WithArg(id), 24*time.Hour, func() (api.CloudServerImage, error) {
			return api.GetVpsImage(id)
		})
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(image, cmd); printed {
			return err
		}

		return ui.RenderForm(image, imageColumns()...)
	},
}

func init() {
	vpsImageCmd.AddCommand(getVpsImageCmd, listVpsImagesCmd)
	vpsCmd.AddCommand(vpsImageCmd)
}

func imageColumns() []ui.TableColumn[api.CloudServerImage] {
	return []ui.TableColumn[api.CloudServerImage]{
		ui.Column("Id", 11, func(i api.CloudServerImage) int { return i.Id }),
		ui.Column("Name", 35, func(i api.CloudServerImage) string { return i.Name }),
		ui.Column("Distro", 10, func(i api.CloudServerImage) string { return i.OSDistro }),
		ui.Column("Version", 10, func(i api.CloudServerImage) string { return i.OSVersion }),
		ui.Column("Release Date", 25, func(i api.CloudServerImage) api.Date { return i.ReleaseDate }),
		ui.Column("Size", 15, func(i api.CloudServerImage) api.Size64 { return i.Size }),
		ui.Column("Virtual Size", 15, func(i api.CloudServerImage) api.Size64 { return i.VirtualSize }),
		ui.Column("Min RAM", 10, func(i api.CloudServerImage) int { return i.MinRAM }),
		ui.Column("Min Disk", 10, func(i api.CloudServerImage) int { return i.MinDisk }),
	}
}

func completeVpsImageIds(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	vpsList, err := cache.Call(cache.KeyVpsImages, cache.DefaultTTL, func() ([]api.CloudServerImage, error) {
		return api.ListVpsImages()
	})

	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var comps []string
	for _, i := range vpsList {
		id := strconv.Itoa(i.Id)
		if strings.HasPrefix(id, toComplete) {
			comps = append(comps, fmt.Sprintf("%s\t%s", id, i.Name))
		}
	}

	return comps, cobra.ShellCompDirectiveNoFileComp
}
