package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/bartoszmajsak/template-golang/version"
)

// NewCmd creates version cmd which prints version and Build details of the executed binary.
func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "version",
		Short:        "Prints the version number of tool",
		Long:         "All software has versions. This is ours",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			PrintVersion()

			return nil
		},
	}
}

func PrintVersion() {
	fmt.Printf("Binary Version: %s\n", version.Version)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Go OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Build Commit: %v\n", version.Commit)
	fmt.Printf("Build Time: %v\n", version.BuildTime)
}
