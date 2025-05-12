package vps

// Currently not in use, can be added to allow UI interaction with vps commands (needs adaptation from POC code below)

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/ui"
	"os"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type UIModel struct {
	table   table.Model
	server  api.CloudServer
	servers []api.CloudServer
	msg     string
}

func (m UIModel) Init() tea.Cmd { return nil }

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.server = m.servers[m.table.Cursor()]
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m UIModel) View() string {
	return m.msg + ":\n\n" + baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func SelectServer(servers []api.CloudServer, msg string) api.CloudServer {
	columns := []table.Column{
		{Title: "Server Id", Width: 11},
		{Title: "Name", Width: 30},
		{Title: "Image", Width: 20},
		{Title: "Status", Width: 15},
	}

	var rows []table.Row

	for _, vps := range servers {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", vps.Id),
			vps.Name,
			vps.Image.Name,
			vps.Status,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t = ui.StyleTable(t)

	m := UIModel{table: t, servers: servers, msg: msg}

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if typedModel, ok := finalModel.(UIModel); ok {
		return typedModel.server
	} else {
		fmt.Println("Error: could not type-assert final model")
		os.Exit(1)
		return api.CloudServer{}
	}
}
