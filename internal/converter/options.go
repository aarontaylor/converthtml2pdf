package converter

import (
	"strconv"
	"strings"
)

type Options struct {
	Format          string
	Landscape       bool
	Margin          string
	Scale           float64
	WaitTime        int
	Timeout         int
	PrintBackground bool
	HeaderTemplate  string
	FooterTemplate  string
	WaitSelector    string
	CookiesFile     string
	UserAgent       string
	ViewportWidth   int
	ViewportHeight  int
	NoImages        bool
	Grayscale       bool
}

func DefaultOptions() Options {
	return Options{
		Format:          "A4",
		Landscape:       false,
		Margin:          "0.4",
		Scale:           1.0,
		WaitTime:        2,
		Timeout:         30,
		PrintBackground: true,
		ViewportWidth:   1920,
		ViewportHeight:  1080,
		NoImages:        false,
		Grayscale:       false,
	}
}

type MarginOptions struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func (o *Options) ParseMargins() MarginOptions {
	margins := MarginOptions{
		Top:    0.4,
		Right:  0.4,
		Bottom: 0.4,
		Left:   0.4,
	}

	if o.Margin == "" {
		return margins
	}

	parts := splitMargins(o.Margin)
	
	switch len(parts) {
	case 1:
		if n, _ := parseFloat(parts[0]); n > 0 {
			margins.Top = n
			margins.Right = n
			margins.Bottom = n
			margins.Left = n
		}
	case 4:
		if n, _ := parseFloat(parts[0]); n >= 0 {
			margins.Top = n
		}
		if n, _ := parseFloat(parts[1]); n >= 0 {
			margins.Right = n
		}
		if n, _ := parseFloat(parts[2]); n >= 0 {
			margins.Bottom = n
		}
		if n, _ := parseFloat(parts[3]); n >= 0 {
			margins.Left = n
		}
	}

	return margins
}

func splitMargins(margin string) []string {
	var parts []string
	for _, part := range strings.Split(margin, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}