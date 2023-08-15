package action

import (
	"errors"
	"github.com/spf13/cobra"
	"io"
)

type errorShouldDisplayUsage struct {
	error
}

func CommandAction(handler func(args []string, stdout io.Writer) error) func(cmd *cobra.Command, args []string) error {
	return func(c *cobra.Command, args []string) error {
		err := handler(args, c.OutOrStdout())
		var shouldDisplayUsage errorShouldDisplayUsage
		if errors.As(err, &shouldDisplayUsage) {
			return c.Help()
		}
		return err
	}
}
