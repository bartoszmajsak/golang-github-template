package main

import (
	"fmt"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bartoszmajsak/template-golang/pkg/cmd/version"
	"github.com/bartoszmajsak/template-golang/pkg/config"
	"github.com/bartoszmajsak/template-golang/pkg/format"
	v "github.com/bartoszmajsak/template-golang/version"
)

func main() {
	rootCmd := newCmd()

	rootCmd.AddCommand(version.NewCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func newCmd() *cobra.Command {
	var configFile string
	releaseInfo := make(chan string, 1)

	rootCmd := &cobra.Command{
		Use: "cmd",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if v.Released() {
				go func() {
					latestRelease, _ := version.LatestRelease()
					if !version.IsLatestRelease(latestRelease) {
						releaseInfo <- fmt.Sprintf("WARN: you are using %s which is not the latest release (newest is %s).\n"+
							"Follow release notes for update info YOUR REPO", v.Version, latestRelease)
					} else {
						releaseInfo <- ""
					}
				}()
			}

			return errors.Wrap(config.SetupConfigSources(loadConfigFileName(cmd)), "failed setting config sources")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			shouldPrintVersion, _ := cmd.Flags().GetBool("version")
			if shouldPrintVersion {
				version.PrintVersion()
			} else {
				fmt.Print(cmd.UsageString())
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if v.Released() {
				timer := time.NewTimer(2 * time.Second)
				select {
				case release := <-releaseInfo:
					fmt.Println(release)
				case <-timer.C:
					// do nothing, just timeout
				}
			}
			close(releaseInfo)

			return nil
		},
	}

	rootCmd.PersistentFlags().
		StringVarP(&configFile, "config", "c", ".ike.config.yaml",
			fmt.Sprintf("config file (supported formats: %s)", strings.Join(config.SupportedExtensions(), ", ")))
	rootCmd.Flags().Bool("version", false, "prints the version number of ike cli")
	rootCmd.PersistentFlags().String("help-format", "standard", "prints help in asciidoc table")
	if err := rootCmd.PersistentFlags().MarkHidden("help-format"); err != nil {
		fmt.Printf("failed while trying to hide a flag: %s\n", err)
	}

	format.EnhanceHelper(rootCmd)
	format.RegisterTemplateFuncs()

	return rootCmd
}

func loadConfigFileName(cmd *cobra.Command) (configFileName string, defaultConfigSource bool) {
	configFlag := cmd.Flag("config")
	configFileName = viper.GetString("config")
	if configFileName == "" {
		if configFlag.Changed {
			configFileName = configFlag.Value.String()
		} else {
			configFileName = configFlag.DefValue
		}
	}
	defaultConfigSource = configFlag.DefValue == configFileName

	return
}
