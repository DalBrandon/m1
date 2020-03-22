// --------------------------------------------------------------------
// tablemaker.go -- Formats printed output into a table.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package util

import (
	"bytes"
	"fmt"
)

type Table struct {
	ncols     int
	header    []string
	rows      [][]string
	colwidths []int
	fixwidths []int
}

// NewTable starts a table.  Input are the names of the columns.
func NewTable(ColumnNames ...string) *Table {
	tbl := new(Table)
	tbl.ncols = len(ColumnNames)
	tbl.header = ColumnNames
	tbl.rows = make([][]string, 0, 1000)
	tbl.colwidths = make([]int, tbl.ncols)
	tbl.fixwidths = make([]int, tbl.ncols)
	for i, v := range tbl.header {
		tbl.colwidths[i] = len(v)
	}
	return tbl
}

// AddRow adds one row to the table.  Input are the columns.
func (tbl *Table) AddRow(values ...string) {
	newrow := make([]string, tbl.ncols)
	for i, s := range values {
		if i >= tbl.ncols {
			break
		}
		if len(s) > tbl.colwidths[i] {
			tbl.colwidths[i] = len(s)
		}
		newrow[i] = s
	}
	tbl.rows = append(tbl.rows, newrow)
}

// SetColumnWidths fixes the width (in chars) of each column.
// Use zero as a width to indicate automatic sizing.
func (tbl *Table) SetColumnWidths(Widths ...int) {
	for i, w := range Widths {
		if i >= len(tbl.fixwidths) {
			return
		}
		tbl.fixwidths[i] = w
	}
}

// GetSimple() is depreciated. It is same as Text().
func (tbl *Table) GetSimple() string {
	return tbl.Text()
}

// Returns the complete table in text, using graphical characters to mark
// the columns and header.
func (tbl *Table) Text() string {
	colw := make([]int, tbl.ncols)
	for i, v := range tbl.colwidths {
		colw[i] = v
		if colw[i] > tbl.fixwidths[i] && tbl.fixwidths[i] != 0 {
			colw[i] = tbl.fixwidths[i]
		}
	}

	var buf bytes.Buffer
	divider := make_divider(colw, "+-", "-+-", "-+", "-")
	fmt.Fprintf(&buf, "%s\n", divider)
	fmt.Fprintf(&buf, "| %s |\n", make_row(tbl.header, colw, " | "))
	fmt.Fprintf(&buf, "%s\n", divider)
	for _, v := range tbl.rows {
		fmt.Fprintf(&buf, "| %s |\n", make_row(v, colw, " | "))
	}
	fmt.Fprintf(&buf, "%s\n", divider)
	return string(buf.Bytes())
}

// +----------------------+-----------------+
// | abcdefghijklmnopqrst |                 |
// +----------------------+-----------------+
// |                      |                 |

func make_divider(widths []int, front, mid, back, fill string) string {
	dst := front
	for i, w := range widths {
		for i := 0; i < w; i++ {
			dst += fill
		}
		if i != len(widths)-1 {
			dst += mid
		}
	}
	dst += back
	return dst
}

// makes a row, with the separator between each column.
func make_row(src []string, widths []int, sep string) string {
	dst := ""
	for i, w := range widths {
		v := ""
		if i < len(src) {
			v = src[i]
		}
		dst += FixStrLen(v, w, "...")
		if i != len(widths)-1 {
			dst += sep
		}
	}
	return dst
}
