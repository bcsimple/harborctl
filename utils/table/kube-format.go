package table

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

type KubeFormat struct {
	Titles      []string
	Data        [][]string
	TableWriter *tablewriter.Table
}

func NewKubeFormat() TableInformation {

	return &KubeFormat{
		Titles:      make([]string, 0, 100),
		Data:        make([][]string, 0, 100),
		TableWriter: tablewriter.NewWriter(os.Stdout),
	}
}

func (t *KubeFormat) SetData(data [][]string) TableInformation {
	t.Data = data
	return t
}

func (t *KubeFormat) SetTitles(titles []string) TableInformation {
	t.Titles = titles
	return t
}
func (t *KubeFormat) CancelHeader() {
	t.TableWriter.SetHeader([]string{})
}

func (t *KubeFormat) Output() {
	t.TableWriter.SetAutoWrapText(false)
	t.TableWriter.SetAutoFormatHeaders(true)
	t.TableWriter.SetNoWhiteSpace(true)
	t.TableWriter.SetCenterSeparator("")
	t.TableWriter.SetColumnSeparator("")
	t.TableWriter.SetRowSeparator("")
	t.TableWriter.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.TableWriter.SetAlignment(tablewriter.ALIGN_LEFT)
	t.TableWriter.SetBorder(false)
	t.TableWriter.SetHeaderLine(false)
	t.TableWriter.SetRowLine(false)
	t.TableWriter.SetTablePadding("\t")
	t.TableWriter.SetFooterAlignment(tablewriter.ALIGN_LEFT)
	t.TableWriter.SetHeader(t.Titles)
	t.TableWriter.AppendBulk(t.Data)
	t.TableWriter.Render()
}
