package table

type TableInformation interface {
	SetData(data [][]string) TableInformation
	CancelHeader()
	SetTitles(titles []string) TableInformation
	Output()
}

func NewTableInformation(style string) TableInformation {
	switch style {
	case "kube":
		return NewKubeFormat()
	case "table":
		return NewTableFormat()
	}
	return NewTableFormat()
}
