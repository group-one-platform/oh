package cmd

import (
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/cache"
	"github.com/edvin/oh/ui"
	"github.com/spf13/cobra"
)

var vpsProductCmd = &cobra.Command{
	Use:   "product",
	Short: "Manage VPS Products",
}

var listVpsProductsCmd = &cobra.Command{
	Use:               "list",
	Short:             "List Available VPS Products",
	SilenceUsage:      true,
	ValidArgsFunction: NoArgs,
	Long: `This endpoint can be used only if it is included in your subscription (your support representative can provide more information regarding how to include it with your subscription).

Provides a list of products as array each containing array(s) of product plans.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		images, err := cache.Call(cache.KeyVpsProducts, cache.DefaultTTL, func() ([]api.Product, error) {
			return api.ListVpsProducts()
		})
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(images, cmd); printed {
			return err
		}

		return ui.RenderTable(images, productColumns()...)
	},
}

func init() {
	vpsProductCmd.AddCommand(listVpsProductsCmd)
	vpsCmd.AddCommand(vpsProductCmd)
}

func productColumns() []ui.TableColumn[api.Product] {
	return []ui.TableColumn[api.Product]{
		ui.Column("Id", 11, func(i api.Product) int { return i.Id }),
		ui.Column("Name", 35, func(i api.Product) string { return i.Name }),
		ui.Column("Plans", 100, func(i api.Product) string { return i.PlansString() }),
	}
}
