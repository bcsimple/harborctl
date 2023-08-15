package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

const version = "v1.0"

func VersionCmd() *cobra.Command {

	command := &cobra.Command{
		Use:   "version",
		Short: "harborctl version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("harborctl version: \n", version)
		},
		Args: cobra.NoArgs,
	}
	return command
}
