// Markdown Table Generator
package md

import (
	"strings"
)

type Cell struct {
	Value string
}

func NewPlainTextCell(data string) Cell {
	return Cell{
		Value: data,
	}
}

func (c *Cell) String() string {
	value := c.Value
	value = strings.Replace(value, "\n", "<br/>", -1)
	value = strings.Replace(value, "\\n", "<br/>", -1)
	return value
}

type Row struct {
	Columns []Cell
}

func (r *Row) String() string {
	scells := []string{}
	for _, c := range r.Columns {
		scells = append(scells, c.String())
	}
	return strings.Join(scells, "|")
}

type Table struct {
	Headers []Cell
	Rows    []Row
}

func (t *Table) generateAlignmentLabel(length int) string {
	if length < 2 {
		length = 2
	}
	length -= 2
	base := ":-"
	for i := 0; i < length; i ++ {
		base += "-"
	}
	return base
}

func (t *Table) righpad(s string, length int) string {
	ls := len(s)
	if length < ls {
		return s
	}
	length -= ls
	for i := 0; i < length; i ++ {
		s += " "
	}
	return s
}

func (t *Table) String() string {
	maxSize := len(t.Headers)
	for _, r := range t.Rows {
		if len(r.Columns) > maxSize {
			maxSize = len(r.Columns)
		}
	}

	colPadSizes := []int{}

	for i := 0; i < maxSize; i ++ {
		colPadSizes = append(colPadSizes, 1)
	}

	for i, c := range t.Headers {
		cl := len(c.String())
		if colPadSizes[i] < cl {
			colPadSizes[i] = cl
		}
	}
	for _, r := range t.Rows {
		for i, c := range r.Columns {
			cl := len(c.String())
			if colPadSizes[i] < cl {
				colPadSizes[i] = cl
			}
		}
	}

	getHeader := func(index int) string {
		if index >= len(t.Headers) || index < 0 {
			return ""
		} else {
			return t.Headers[index].String()
		}
	}

	sheaders := []string{}
	alignments := []string{}
	for i := 0; i < maxSize; i ++ {
		cHeader := getHeader(i)
		alignment := t.generateAlignmentLabel(colPadSizes[i])
		sheaders = append(sheaders, t.righpad(cHeader, colPadSizes[i]))
		alignments = append(alignments, alignment)
	}

	sdata := []string{}

	for _, r := range t.Rows {
		srdata := []string{}
		for i, c := range r.Columns {
			srdata = append(srdata, t.righpad(c.String(), colPadSizes[i]))
		}
		// append empty cells
		for i := len(r.Columns); i < maxSize; i ++ {
			srdata = append(srdata, t.righpad("", colPadSizes[i]))
		}
		sdata = append(sdata, "|"+strings.Join(srdata, "|")+"|")
	}

	return strings.Join(append([]string{
		"|" + strings.Join(sheaders, "|") + "|",
		"|" + strings.Join(alignments, "|") + "|",
	}, sdata...), "\n")
}
