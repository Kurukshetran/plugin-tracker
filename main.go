package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/jroimartin/gocui"
)

func getNpmPackages() string {
	cmd := exec.Command("npm", "list", "-g", "--depth=0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return out.String()
}

func getPipPackages() string {
	cmd := exec.Command("pip", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return out.String()
}

func getMavenPlugins() string {
	cmd := exec.Command("mvn", "dependency:tree")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return out.String()
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("menu", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Menu"
		fmt.Fprintln(v, "Press n for npm packages, p for pip packages, m for Maven plugins")
	}

	if v, err := g.SetView("output", 0, 3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Output"
		v.Wrap = true
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'n', gocui.ModNone, displayNpmPackages); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'p', gocui.ModNone, displayPipPackages); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'm', gocui.ModNone, displayMavenPlugins); err != nil {
		return err
	}
	return nil
}

func displayNpmPackages(g *gocui.Gui, v *gocui.View) error {
	output := getNpmPackages()
	return displayOutput(g, output)
}

func displayPipPackages(g *gocui.Gui, v *gocui.View) error {
	output := getPipPackages()
	return displayOutput(g, output)
}

func displayMavenPlugins(g *gocui.Gui, v *gocui.View) error {
	output := getMavenPlugins()
	return displayOutput(g, output)
}

func displayOutput(g *gocui.Gui, output string) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("output")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, output)
		return nil
	})
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
