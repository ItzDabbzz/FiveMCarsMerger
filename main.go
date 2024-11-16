package main

import (
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/config"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/flags"
	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/merger"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Flags flags.Flags

func main() {
	appFlags, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if appFlags == nil {
		// First time setup
		appFlags = &flags.Flags{}
		if err := initialSetup(appFlags); err != nil {
			log.Fatal(err)
		}
		if err := config.SaveConfig(appFlags); err != nil {
			log.Fatal(err)
		}
	}

	for {
		mainMenu := []string{"Start Merge Process", "Edit Settings", "Exit"}
		var selected string

		form := huh.NewSelect[string]().
			Title("FiveM Cars Merger").
			Options(huh.NewOptions(mainMenu...)...).
			Value(&selected)

		if err := form.Run(); err != nil {
			log.Fatal(err)
		}

		switch selected {
		case "Start Merge Process":
			ConfigureLogger()
			carsMerger := merger.New(*appFlags)
			if err := carsMerger.Merge(); err != nil {
				log.Fatal(err)
			}
		case "Edit Settings":
			if err := editSettings(appFlags); err != nil {
				log.Fatal(err)
			}
			if err := config.SaveConfig(appFlags); err != nil {
				log.Fatal(err)
			}
		case "Exit":
			return
		}
	}
}

func initialSetup(flags *flags.Flags) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input Path").
				Description("Path to all cars").
				Placeholder("cars").
				Value(&flags.InputPath),
			huh.NewInput().
				Title("Output Path").
				Description("Output path for merged cars").
				Placeholder("merged-cars").
				Value(&flags.OutputPath),
			huh.NewConfirm().
				Title("Enable Verbose Logging").
				Value(&flags.Verbose),
			huh.NewConfirm().
				Title("Clean Output Directory").
				Value(&flags.Clean),
		).Title("Initial Configuration"),
	).Run()
}

func editSettings(flags *flags.Flags) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input Path").
				Description("Path to all cars").
				Value(&flags.InputPath),
			huh.NewInput().
				Title("Output Path").
				Description("Output path for merged cars").
				Value(&flags.OutputPath),
			huh.NewConfirm().
				Title("Enable Verbose Logging").
				Value(&flags.Verbose),
			huh.NewConfirm().
				Title("Clean Output Directory").
				Value(&flags.Clean),
		).Title("Edit Settings"),
	).Run()
}

func ConfigureLogger() {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff7df9"))

	log.SetFormatter(log.TextFormatter)
	log.SetReportCaller(true)
	log.SetPrefix(style.Render("FiveMCarsMerger"))

	if Flags.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
