package pages

import (
	"io"
	"io/fs"

	"fmt"
	"html/template"
	"math"
	"time"

	"embed"
)

//go:embed themes/*
var Themes embed.FS

type RenderOptionsInstanceState struct {
	Name            string
	CurrentReplicas int
	DesiredReplicas int
	Status          string
	Error           error
}

type RenderOptions struct {
	DisplayName      string
	ShowDetails      bool
	InstanceStates   []RenderOptionsInstanceState
	SessionDuration  time.Duration
	RefreshFrequency time.Duration
	Theme            string
	CustomThemes     fs.FS
	// If custom theme is loaded through os.DirFS, nothing prevents you
	// from escaping the prefix with relative path such as ..
	// The `AllowedCustomThemes` are the themes that were scanned during initilization
	AllowedCustomThemes map[string]bool
	Version             string
}

type TemplateValues struct {
	DisplayName      string
	InstanceStates   []RenderOptionsInstanceState
	SessionDuration  string
	RefreshFrequency string
	Version          string
}

func Render(options RenderOptions, writer io.Writer) error {
	var tpl *template.Template
	var err error

	// Load custom theme if provided
	if options.CustomThemes != nil && options.AllowedCustomThemes[options.Theme] {
		tpl, err = template.ParseFS(options.CustomThemes, fmt.Sprintf("%s.html", options.Theme))
	} else {
		// Load from the embedded FS
		tpl, err = template.ParseFS(Themes, fmt.Sprintf("themes/%s.html", options.Theme))
	}

	if err != nil {
		return err
	}

	instanceStates := []RenderOptionsInstanceState{}

	if options.ShowDetails {
		instanceStates = options.InstanceStates
	}

	return tpl.Execute(writer, TemplateValues{
		DisplayName:      options.DisplayName,
		InstanceStates:   instanceStates,
		SessionDuration:  humanizeDuration(options.SessionDuration),
		RefreshFrequency: fmt.Sprintf("%d", int64(options.RefreshFrequency.Seconds())),
		Version:          options.Version,
	})
}

// humanizeDuration humanizes time.Duration output to a meaningful value,
// golang's default “time.Duration“ output is badly formatted and unreadable.
func humanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		if remainingSeconds > 0 {
			return fmt.Sprintf("%d minutes %d seconds", int64(duration.Minutes()), int64(remainingSeconds))
		}
		return fmt.Sprintf("%d minutes", int64(duration.Minutes()))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)

		if remainingMinutes > 0 {
			if remainingSeconds > 0 {
				return fmt.Sprintf("%d hours %d minutes %d seconds", int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
			}
			return fmt.Sprintf("%d hours %d minutes", int64(duration.Hours()), int64(remainingMinutes))
		}
		return fmt.Sprintf("%d hours", int64(duration.Hours()))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}
