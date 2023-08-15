/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package download

import (
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/spf13/cobra"
	"io"
)

type downloadOptions struct {
	global *root.GlobalOptions
	Dname  string
	Dpath  string
}

func DownloadCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &downloadOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "download",
		Short: "download chart or image from harbor",
		Long: `this is a simple down load chart to local
explain:
download chart:
 harborctl download chart -n hello
download image:
 harborctl download image -n hello
`,
		RunE: action.CommandAction(opts.run),
	}

	command.Flags().StringVarP(&opts.Dname, "name", "n", "", "download chart by name required!")
	command.Flags().StringVarP(&opts.Dpath, "path", "p", ".", "download chart to directory default . ")
	// Here you will define your flags and configuration settings.
	return command
}

func (d *downloadOptions) run(args []string, stdout io.Writer) error {
	name := args[0]
	if d.Dname == "" {
		return fmt.Errorf("download chart name must bu required")
	}

	switch name {
	case "chart":
		client.NewChart(d.global).DownloadChart(d.Dname, d.Dpath)
	case "image":
		client.NewImage(d.global).DownloadImage(d.Dname, d.Dpath)
	}
	return nil
}
