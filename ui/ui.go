package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func StyleTable(t table.Model) table.Model {
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return t
}

type TableColumn[T any] struct {
	Title string
	Width int
	Value func(T) string
}

func Column[T any, V any](title string, width int, extract func(T) V) TableColumn[T] {
	return TableColumn[T]{
		Title: title,
		Width: width,
		Value: func(t T) string {
			return fmt.Sprint(extract(t))
		},
	}
}

func RenderTable[T any](items []T, cols ...TableColumn[T]) error {
	tblCols := make([]table.Column, len(cols))
	for i, c := range cols {
		tblCols[i] = table.Column{Title: c.Title, Width: c.Width}
	}

	tblRows := make([]table.Row, len(items))
	for i, item := range items {
		row := make(table.Row, len(cols))
		for j, c := range cols {
			row[j] = c.Value(item)
		}
		tblRows[i] = row
	}

	t := table.New(
		table.WithColumns(tblCols),
		table.WithRows(tblRows),
	)

	t.SetStyles(table.Styles{
		Header:   table.Styles{}.Header,
		Cell:     table.Styles{}.Cell,
		Selected: table.Styles{}.Cell,
	})

	out := strings.TrimSpace(t.View())
	_, err := fmt.Println(out)
	return err
}

func RenderForm[T any](item T, cols ...TableColumn[T]) error {
	for _, col := range cols {
		val := col.Value(item)
		_, err := fmt.Printf("%s: %v\n", col.Title, val)
		if err != nil {
			return err
		}
	}
	return nil
}
