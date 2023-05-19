package cmd

import (
	"fmt"
	"github.com/engswee/flashpipe/runner"
	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg"},
	Short:   "Upload/update integration package",
	Long: `Upload or update integration package on the
SAP Integration Suite tenant.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[INFO] Executing update package command")

		setMandatoryVariable("package.file", "PACKAGE_FILE")
		setOptionalVariable("package.override.id", "PACKAGE_ID")
		setOptionalVariable("package.override.id", "PACKAGE_NAME")

		runner.JavaCmd("io.github.engswee.flashpipe.cpi.exec.UpdateIntegrationPackage", mavenRepoLocation, flashpipeLocation, log4jFile)
	},
}

func init() {
	updateCmd.AddCommand(packageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setStringFlagAndBind(packageCmd, "package.file", "", "Path to location of package file [or set environment PACKAGE_FILE]")
	setStringFlagAndBind(packageCmd, "package.override.id", "", "Override package ID from file [or set environment PACKAGE_ID]")
	setStringFlagAndBind(packageCmd, "package.override.name", "", "Override package name from file [or set environment PACKAGE_NAME]")
}