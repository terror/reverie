package main

import (
  "fmt"
  "os"

  tea "github.com/charmbracelet/bubbletea"
)

func main() {
  cwd, err := os.Getwd()

  if err != nil {
    fmt.Println("error:", err)
    os.Exit(1)
  }

  root, err := buildTree(cwd)

  if err != nil {
    fmt.Println("error:", err)
    os.Exit(1)
  }

  p := tea.NewProgram(
    &model{root: root},
    tea.WithAltScreen(),
    tea.WithMouseCellMotion(),
  )

  if _, err := p.Run(); err != nil {
    fmt.Println("error:", err)
    os.Exit(1)
  }
}
