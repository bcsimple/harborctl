package scan

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

type scanOptions struct {
	global               *root.GlobalOptions
	WithFile             bool
	WithCompare          bool
	WithCompareOnlyFalse bool
	path                 string
	release              string
}

func ScanCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &scanOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "scan",
		Short: "Scan image file and compare",
		Long: `For Example:
 harborctl scan -f FILENAME`,
		RunE: action.CommandAction(opts.run),
	}

	command.Flags().BoolVarP(&opts.WithFile, "file", "F", false, "check and print auto deploy script images")
	command.Flags().BoolVarP(&opts.WithCompare, "compare", "C", false, "compare images from script and harbor")
	command.Flags().BoolVarP(&opts.WithCompareOnlyFalse, "diff", "d", false, "compare images from script and harbor (only false type)")
	command.Flags().StringVarP(&opts.path, "path", "s", "", "images file with your indicate path")
	command.Flags().StringVarP(&opts.release, "release", "r", "paas_v20230101", "panji release")
	command.Flags().SortFlags = false
	_ = command.MarkFlagRequired("path")
	if !opts.WithFile && !opts.WithCompare && !opts.WithCompareOnlyFalse {
		pflag.Usage()
	}

	return command
}

func (c *scanOptions) run(args []string, stdout io.Writer) error {
	scan := client.NewScanImage(c.global, c.WithFile, c.WithCompare, c.WithCompareOnlyFalse, c.path, c.release)
	if c.WithFile {
		scan.PrintFile()
		return nil
	}
	if c.WithCompare {
		scan.PrintCompare()
		return nil
	}
	if c.WithCompareOnlyFalse {
		scan.PrintDiff()
		return nil
	}
	return nil
}
