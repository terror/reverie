package main

import (
  "fmt"
  "os"

  tea "github.com/charmbracelet/bubbletea"
)

type model struct {
  ready bool
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c", "q":
      return m, tea.Quit
    }
  case tea.WindowSizeMsg:
    m.ready = true
  }

  return m, nil
}

func (m model) View() string {
  if !m.ready {
    return "loading…"
  }

  return "hello bubbletea\n\n(q to quit)"
}

func main() {
  p := tea.NewProgram(
    model{},
    tea.WithAltScreen(),
    tea.WithMouseCellMotion(),
  )

  if _, err := p.Run(); err != nil {
    fmt.Println("error:", err)
    os.Exit(1)
  }
}
