package cmd

import (
	"fmt"
	"github.com/edvin/oh/config"
	tokenui "github.com/edvin/oh/ui/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Supply your credentials to login.",
	Long:  `Log into oneHome to generate an API token and supply it here to store it in the configuration.`,
	Example: `  # prompt for token
  oh token

  # pass token directly
  oh token abc123

  # read from a file or pipe
  oh token < mytoken.txt
  cat mytoken.txt | oh token`,
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var token string
		// positional argument
		if len(args) == 1 {
			token = args[0]
		} else {
			// piped or redirected stdin
			fi, _ := os.Stdin.Stat()
			if (fi.Mode() & os.ModeCharDevice) == 0 {
				// token from a pipe or file
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("reading token from stdin: %w", err)
				}
				token = strings.TrimSpace(string(data))
			}
			// fallback to interactive UI
			if token == "" {
				t, err := tokenui.RequestToken()
				if err != nil {
					return err
				}
				token = t
			}
		}

		viper.Set("token", token)
		if err := config.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		cmd.Println("🔐 Token stored!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
