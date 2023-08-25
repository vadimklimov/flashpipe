package cmd

import (
	"fmt"
	"github.com/engswee/flashpipe/internal/config"
	"github.com/engswee/flashpipe/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func NewCmdRoot() *cobra.Command {
	var version = "3.0.0" // FLASHPIPE_VERSION

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     "flashpipe",
		Version: version,
		Short:   "FlashPipe - The CI/CD Companion for SAP Integration Suite",
		Long: `FlashPipe - The CI/CD Companion for SAP Integration Suite

FlashPipe is a CLI that is used to simplify the Build-To-Deploy cycle
for SAP Integration Suite by providing CI/CD capabilities for 
automating time-consuming manual tasks like:
- synchronising integration artifacts to Git
- creating/updating integration artifacts to SAP Integration Suite
- deploying integration artifacts on SAP Integration Suite`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
	}

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/flashpipe.yaml)")

	// Define cobra flags, the default value has the lowest (least significant) precedence
	rootCmd.PersistentFlags().String("tmn-host", "", "Host for tenant management node of Cloud Integration excluding https:// [or set environment HOST_TMN]")
	rootCmd.PersistentFlags().String("tmn-userid", "", "User ID for Basic Auth [or set environment BASIC_USERID]")
	rootCmd.PersistentFlags().String("tmn-password", "", "Password for Basic Auth [or set environment BASIC_PASSWORD]")
	rootCmd.PersistentFlags().String("oauth-host", "", "Host for OAuth token server excluding https:// [or set environment HOST_OAUTH]")
	rootCmd.PersistentFlags().String("oauth-clientid", "", "Client ID for using OAuth [or set environment OAUTH_CLIENTID]")
	rootCmd.PersistentFlags().String("oauth-clientsecret", "", "Client Secret for using OAuth [or set environment OAUTH_CLIENTSECRET]")
	rootCmd.PersistentFlags().String("oauth-path", "/oauth/token", "Path for OAuth token server, e.g /oauth2/api/v1/token for Neo [or set environment HOST_OAUTH_PATH]")

	rootCmd.PersistentFlags().Bool("debug", false, "Show debug logs")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	rootCmd := NewCmdRoot()
	rootCmd.AddCommand(NewDeployCommand())
	rootCmd.AddCommand(NewSyncCommand())
	updateCmd := NewUpdateCommand()
	updateCmd.AddCommand(NewArtifactCommand())
	updateCmd.AddCommand(NewPackageCommand())
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(NewSnapshotCommand())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initializeConfig(cmd *cobra.Command) error {
	cfgFile := config.GetString(cmd, "config")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name "flashpipe.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("flashpipe")
	}

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	viper.SetEnvPrefix("FLASHPIPE")

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, e.g. --artifact-id to FLASHPIPE_ARTIFACT_ID
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables
	viper.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd)

	// Set debug flag from command line to viper
	if !viper.IsSet("debug") {
		viper.Set("debug", config.GetBool(cmd, "debug"))
	}

	logger.InitConsoleLogger(viper.GetBool("debug"))

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)
			cmd.Flags().Set(configName, fmt.Sprintf("%v", val))
		}
	})
}