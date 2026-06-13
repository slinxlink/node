package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var Version = "dev"

type updateResultMsg string

type item struct {
	label       string
	labelFunc   func() string
	action      func() string
	asyncAction func() tea.Cmd
	quit        bool
	divider     bool
	selectable  bool
}

func (it item) getLabel() string {
	if it.labelFunc != nil {
		return it.labelFunc()
	}
	return it.label
}

type model struct {
	items       []item
	cursor      int
	startCursor int
	output      string
	loading     bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateResultMsg:
		m.loading = false
		m.output = string(msg)
		return m, nil
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			for m.cursor > m.startCursor {
				m.cursor--
				if m.items[m.cursor].selectable {
					break
				}
			}
		case "down", "j":
			for m.cursor < len(m.items)-1 {
				m.cursor++
				if m.items[m.cursor].selectable {
					break
				}
			}
		case "enter":
			if !m.items[m.cursor].selectable {
				return m, nil
			}
			if m.items[m.cursor].quit {
				return m, tea.Quit
			}
			if m.items[m.cursor].asyncAction != nil {
				m.loading = true
				m.output = ""
				return m, m.items[m.cursor].asyncAction()
			}
			m.output = m.items[m.cursor].action()
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	type divInfo struct {
		lineIdx int
	}
	var dividers []divInfo
	lineIdx := 0

	maxWidth := 0
	for _, it := range m.items {
		if !it.divider {
			w := lipgloss.Width("  " + it.getLabel())
			if w > maxWidth {
				maxWidth = w
			}
		}
	}

	lines := ""
	for i, it := range m.items {
		if it.divider {
			lines += strings.Repeat(" ", maxWidth)
			dividers = append(dividers, divInfo{lineIdx: lineIdx + 1})
		} else if it.selectable && i == m.cursor {
			lines += lipgloss.NewStyle().Bold(true).Foreground(primaryColor).Render("● " + it.getLabel())
		} else {
			lines += "  " + it.getLabel()
		}
		if i < len(m.items)-1 {
			lines += "\n"
			lineIdx++
		}
	}

	rendered := box(lines)
	renderedLines := strings.Split(rendered, "\n")

	for _, div := range dividers {
		if div.lineIdx < len(renderedLines) {
			line := renderedLines[div.lineIdx]
			w := lipgloss.Width(line)
			renderedLines[div.lineIdx] = "├" + strings.Repeat("─", w-2) + "┤"
		}
	}

	result := strings.Join(renderedLines, "\n") + "\n"
	if m.loading {
		result += renderStatus("更新", "下载中，请稍候...", true) + "\n"
	} else if m.output != "" {
		result += m.output + "\n"
	}
	return result
}

func Start(version string) {
	Version = version
	items := []item{
		{labelFunc: func() string { return "版本                " + Version }, selectable: false},
		{labelFunc: func() string { return "面板                " + panelStatus() }, selectable: false},
		{labelFunc: func() string { return "核心                " + coreStatus() }, selectable: false},
		{labelFunc: func() string { return "自动重启            " + autoRestartStatus() }, selectable: false},
		{divider: true},
		{label: "退出脚本", selectable: true, quit: true},
		{divider: true},
		{label: "启动面板", selectable: true, action: func() string { return startPanel() }},
		{label: "停止面板", selectable: true, action: func() string { return stopPanel() }},
		{label: "重启面板", selectable: true, action: func() string { return restartPanel() }},
		{labelFunc: func() string {
			if autoRestartStatus() == lipgloss.NewStyle().Foreground(successColor).Render("● 已开启") {
				return "关闭自动重启"
			}
			return "开启自动重启"
		}, selectable: true, action: func() string { return toggleAutoRestart() }},
		{divider: true},
		{label: "重置面板地址", selectable: true, action: func() string { return resetUrl() }},
		{label: "重置用户名&密码", selectable: true, action: func() string { return resetCredentials() }},
		{label: "查看登录信息", selectable: true, action: func() string { return showLoginInfo() }},
		{divider: true},
		{label: "更新", selectable: true, asyncAction: update()},
		{label: "卸载", selectable: true, action: func() string { return uninstall() }},
	}

	m := model{items: items}
	for i, it := range items {
		if it.selectable {
			m.cursor = i
			m.startCursor = i
			break
		}
	}
	m.output = firstRun()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
