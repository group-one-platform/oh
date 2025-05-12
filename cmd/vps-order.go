package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/ui"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	orderFile string
)

var orderVpsCmd = &cobra.Command{
	Use:   "order [<order-json>]",
	Short: "Order VPS",
	Long: `This endpoint can be used only if it is included in your subscription (your support representative can provide more information regarding how to include it with your subscription).

Allows ordering of a new VPS instance by specifying required information during provisioning.

Pass the JSON payload describing the new VPS either as a positional argument or via --file (use '-' for stdin).

Example usage:

# Read order from file
oh vps order -f my-order.json

# Or pipe through stdin
cat my-order.json | oh vps order -f -

Example payload:

{
  "productId": 123,
  "productPlanId": 123,
  "imageId": 123,
  "password": "asdf",
  "availabilityZone": "asdf",
  "name": "test-server",
  "sshKey": "asdf",
  "storageSize": "10",
  "networks": [
    {"network": "uuid", "fixed_ipv4": "192.168.1.1" }
  ]
}
`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var reader io.Reader
		switch {
		case orderFile == "-":
			reader = os.Stdin
		case orderFile != "":
			f, err := os.Open(orderFile)
			if err != nil {
				return fmt.Errorf("cannot open %q: %w", orderFile, err)
			}
			defer f.Close()
			reader = f
		case len(args) == 1:
			reader = strings.NewReader(args[0])
		default:
			return fmt.Errorf("you must supply JSON via a positional arg or --file")
		}

		// decode with strict checking
		var order api.CloudServerOrder
		dec := json.NewDecoder(reader)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&order); err != nil {
			return fmt.Errorf("invalid order payload: %w", err)
		}

		response, err := api.OrderVps(order)
		if err != nil {
			return err
		}

		if printed, err := PrintJSON(response, cmd); printed {
			return err
		}

		return ui.RenderForm(response, vpsOrderColumns()...)
	},
}

func init() {
	orderVpsCmd.Flags().
		StringVarP(&orderFile, "file", "f", "",
			"JSON file to read order from (`-` for stdin); if omitted you can pass raw JSON as the sole positional argument")
	vpsCmd.AddCommand(orderVpsCmd)
}

func vpsOrderColumns() []ui.TableColumn[api.CloudServerOrderResponse] {
	return []ui.TableColumn[api.CloudServerOrderResponse]{
		ui.Column("Id", 11, func(i api.CloudServerOrderResponse) int { return i.Id }),
		ui.Column("ContractId", 11, func(i api.CloudServerOrderResponse) int { return i.ContractId }),
		ui.Column("OrderId", 11, func(i api.CloudServerOrderResponse) string { return i.OrderId }),
	}
}
