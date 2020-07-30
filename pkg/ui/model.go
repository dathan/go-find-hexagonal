package ui

import (
	"fmt"
	"strings"

	"github.com/dathan/go-find-hexagonal/pkg/find"
	"github.com/davecheney/errors"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SatisifyString string

func (s SatisifyString) String() string {
	return string(s)
}

func Display(f find.FindResults) error {

	if err := ui.Init(); err != nil {
		return errors.Annotate(err, "ui.Display(...) init failed")
	}

	defer ui.Close()

	nodes := []*widgets.TreeNode{}

	for _, row := range f {
		paths := strings.Split(row.Path, "\n\n")
		fmt.Printf("Fixing Row: %s\n\t%v\n", row.Name, paths)
		node := &widgets.TreeNode{}
		node.Value = SatisifyString(row.Name)
		for _, path := range paths {
			node.Nodes = append(node.Nodes, &widgets.TreeNode{Value: SatisifyString(path)})
		}
		nodes = append(nodes, node)
	}

	l := widgets.NewTree()
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	x, y := ui.TerminalDimensions()
	l.SetRect(0, 0, x, y)

	l.SetNodes(nodes)
	ui.Render(l)

	uiEvents := ui.PollEvents()
	previousKey := ""
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return nil
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "<Enter>":
			l.ToggleExpand()
		case "G", "<End>":
			l.ScrollBottom()
		case "E":
			l.ExpandAll()
		case "C":
			l.CollapseAll()
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			l.SetRect(0, 0, x, y)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(l)
	}
}
