package table

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

type TableFormat struct {
	Titles      []string
	Data        [][]string
	TableWriter *tablewriter.Table
}

func NewTableFormat() TableInformation {

	return &TableFormat{
		Titles:      make([]string, 0, 100),
		Data:        make([][]string, 0, 100),
		TableWriter: tablewriter.NewWriter(os.Stdout),
	}
}

func (t *TableFormat) SetData(data [][]string) TableInformation {
	t.Data = data
	return t
}

func (t *TableFormat) SetTitles(titles []string) TableInformation {
	t.Titles = titles
	return t
}
func (t *TableFormat) CancelHeader() {
	t.TableWriter.SetHeader([]string{})
}

func (t *TableFormat) Output() {
	t.TableWriter.SetHeader(t.Titles) // Enable row line
	t.TableWriter.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.TableWriter.SetAlignment(tablewriter.ALIGN_LEFT)
	t.TableWriter.AppendBulk(t.Data)
	t.TableWriter.Render()
	t.TableWriter.ClearRows()
}
