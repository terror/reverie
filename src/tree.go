package main

import (
  "fmt"
  "io/fs"
  "path/filepath"
  "sort"
  "strings"

  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  devicons "github.com/epilande/go-devicons"
)

type node struct {
  name     string
  path     string
  isDir    bool
  expanded bool
  depth    int
  parent   *node
  children []*node
}

type model struct {
  root   *node
  cursor int
  width  int
  height int
}

func (m *model) Init() tea.Cmd {
  return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c", "q":
      return m, tea.Quit
    case "j", "down":
      m.moveCursor(1)
    case "k", "up":
      m.moveCursor(-1)
    case "h", "left":
      m.collapseCurrent()
    case "l", "right":
      m.expandCurrent()
    case "enter":
      m.toggleCurrent()
    case "g", "home":
      m.cursor = 0
      m.clampCursor()
    case "G", "end":
      nodes := m.visibleNodes()
      if len(nodes) > 0 {
        m.cursor = len(nodes) - 1
      }
    }
  case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
  }

  return m, nil
}

func (m *model) View() string {
  if m.root == nil {
    return "loading..."
  }

  nodes := m.visibleNodes()

  if len(nodes) == 0 {
    return "(empty)\n\n(q to quit)"
  }

  var b strings.Builder

  for i, n := range nodes {
    cursor := " "
    if i == m.cursor {
      cursor = ">"
    }

    twisty := twistyFor(n)

    icon, iconColor := iconStyleFor(n)

    if iconColor != "" {
      icon = lipgloss.NewStyle().Foreground(lipgloss.Color(iconColor)).Render(icon)
    }

    indent := strings.Repeat("  ", n.depth)

    fmt.Fprintf(&b, "%s %s%s %s %s\n", cursor, indent, twisty, icon, n.name)
  }

  b.WriteString("\n(q to quit)")

  return b.String()
}

func buildTree(rootPath string) (*node, error) {
  root := &node{
    name:     filepath.Base(rootPath),
    path:     rootPath,
    isDir:    true,
    expanded: true,
    depth:    0,
  }

  nodes := map[string]*node{
    rootPath: root,
  }

  err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
    if err != nil {
      return err
    }

    if path == rootPath {
      return nil
    }

    if d.IsDir() && d.Name() == ".git" {
      return filepath.SkipDir
    }

    parent := nodes[filepath.Dir(path)]

    if parent == nil {
      return nil
    }

    n := &node{
      name:   d.Name(),
      path:   path,
      isDir:  d.IsDir(),
      depth:  parent.depth + 1,
      parent: parent,
    }

    parent.children = append(parent.children, n)

    if d.IsDir() {
      nodes[path] = n
    }

    return nil
  })

  if err != nil {
    return nil, err
  }

  sortTree(root)

  return root, nil
}

func sortTree(n *node) {
  sort.Slice(n.children, func(i, j int) bool {
    a := n.children[i]
    b := n.children[j]

    if a.isDir != b.isDir {
      return a.isDir
    }

    return strings.ToLower(a.name) < strings.ToLower(b.name)
  })

  for _, child := range n.children {
    if child.isDir {
      sortTree(child)
    }
  }
}

func twistyFor(n *node) string {
  if !n.isDir {
    return " "
  }

  if n.expanded {
    return "▾"
  }

  return "▸"
}

func iconStyleFor(n *node) (string, string) {
  style := devicons.IconForPath(n.path)

  if style.Icon != "" {
    return style.Icon, style.Color
  }

  if n.isDir {
    return "", ""
  }

  return "󰈙", ""
}

func (m *model) visibleNodes() []*node {
  if m.root == nil {
    return nil
  }

  nodes := make([]*node, 0, 256)

  var walk func(*node)

  walk = func(n *node) {
    nodes = append(nodes, n)

    if n.isDir && n.expanded {
      for _, child := range n.children {
        walk(child)
      }
    }
  }

  walk(m.root)

  return nodes
}

func (m *model) currentNode() *node {
  nodes := m.visibleNodes()

  if len(nodes) == 0 || m.cursor < 0 || m.cursor >= len(nodes) {
    return nil
  }

  return nodes[m.cursor]
}

func (m *model) moveCursor(delta int) {
  nodes := m.visibleNodes()

  if len(nodes) == 0 {
    m.cursor = 0
    return
  }

  m.cursor += delta

  m.clampCursor()
}

func (m *model) clampCursor() {
  nodes := m.visibleNodes()

  if len(nodes) == 0 {
    m.cursor = 0
    return
  }

  if m.cursor < 0 {
    m.cursor = 0
  } else if m.cursor >= len(nodes) {
    m.cursor = len(nodes) - 1
  }
}

func (m *model) expandCurrent() {
  n := m.currentNode()

  if n == nil || !n.isDir {
    return
  }

  n.expanded = true

  m.clampCursor()
}

func (m *model) collapseCurrent() {
  n := m.currentNode()

  if n == nil {
    return
  }

  if n.isDir && n.expanded {
    n.expanded = false
    m.clampCursor()
    return
  }

  if n.parent != nil {
    m.cursor = m.indexOfNode(n.parent)
  }
}

func (m *model) toggleCurrent() {
  n := m.currentNode()

  if n == nil || !n.isDir {
    return
  }

  n.expanded = !n.expanded

  m.clampCursor()
}

func (m *model) indexOfNode(target *node) int {
  nodes := m.visibleNodes()

  for i, n := range nodes {
    if n == target {
      return i
    }
  }

  return m.cursor
}
