package paint

import (
	"github.com/jedib0t/go-pretty/table"
	"io"
	"os"
)

type Paint struct {
	writer table.Writer
	output io.Writer
}

func New() *Paint {
	t := &Paint{
		writer: table.NewWriter(),
		output: nil,
	}
	t.defaultTableStyle()

	return t
}

func (t *Paint)defaultTableStyle()  {
	t.SetOutput(os.Stdout)
	t.writer.SetStyle(table.StyleLight)
	t.writer.SetAutoIndex(false)
}

func (t *Paint) SetOutput(w io.Writer)  {
	t.writer.SetOutputMirror(w)
}

func (t *Paint) SetHeader(r table.Row) {
	t.writer.AppendHeader(r)
}

func (t *Paint) SetFooter(r table.Row) {
	t.writer.AppendFooter(r)
}

func (t *Paint) SetRow(r table.Row) {
	t.writer.AppendRow(r)
}

func (t *Paint) SetRows(r []table.Row) {
	t.writer.AppendRows(r)
}

func (t *Paint) Render() {
	t.writer.Render()
}