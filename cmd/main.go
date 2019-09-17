package main

import (
	"fmt"
	"github.com/maistra/istio-workspace/pkg/cmd/config"
	"time"

	"github.com/bartoszmajsak/template-golang/pkg/cmd/version"

	"github.com/bartoszmajsak/template-golang/pkg/format"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := newCmd()

	rootCmd.AddCommand(version.NewCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func newCmd() *cobra.Command {
	releaseInfo := make(chan string, 1)

	rootCmd := &cobra.Command{
		Use: "cmd",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { //nolint[:unparam]
			if v.Released() {
				go func() {
					latestRelease, _ := version.LatestRelease()
					if !version.IsLatestRelease(latestRelease) {
						releaseInfo <- fmt.Sprintf("WARN: you are using %s which is not the latest release (newest is %s).\n"+
							"Follow release notes for update info https://github.com/Maistra/istio-workspace/releases/latest", v.Version, latestRelease)
					} else {
						releaseInfo <- ""
					}
				}()
			}
			return config.SetupConfigSources(configFile, cmd.Flag("config").Changed)
		},
		RunE: func(cmd *cobra.Command, args []string) error { //nolint[:unparam]
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
					log.Info(release)
				case <-timer.C:
					// do nothing, just timeout
				}
			}
			close(releaseInfo)
			return nil
		},
	}

	format.EnhanceHelper(rootCmd)
	format.RegisterTemplateFuncs()

	return rootCmd
}
