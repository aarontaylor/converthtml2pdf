package converter

import (
	"testing"
)

func TestParseMargins(t *testing.T) {
	tests := []struct {
		name   string
		margin string
		want   MarginOptions
	}{
		{
			name:   "single value",
			margin: "1.0",
			want: MarginOptions{
				Top:    1.0,
				Right:  1.0,
				Bottom: 1.0,
				Left:   1.0,
			},
		},
		{
			name:   "four values",
			margin: "1.0,2.0,3.0,4.0",
			want: MarginOptions{
				Top:    1.0,
				Right:  2.0,
				Bottom: 3.0,
				Left:   4.0,
			},
		},
		{
			name:   "empty margin",
			margin: "",
			want: MarginOptions{
				Top:    0.4,
				Right:  0.4,
				Bottom: 0.4,
				Left:   0.4,
			},
		},
		{
			name:   "invalid format defaults",
			margin: "1.0,2.0",
			want: MarginOptions{
				Top:    0.4,
				Right:  0.4,
				Bottom: 0.4,
				Left:   0.4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{Margin: tt.margin}
			got := opts.ParseMargins()
			
			if got.Top != tt.want.Top || got.Right != tt.want.Right ||
				got.Bottom != tt.want.Bottom || got.Left != tt.want.Left {
				t.Errorf("ParseMargins() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	
	if opts.Format != "A4" {
		t.Errorf("Default format = %s, want A4", opts.Format)
	}
	
	if opts.Timeout != 30 {
		t.Errorf("Default timeout = %d, want 30", opts.Timeout)
	}
	
	if opts.WaitTime != 2 {
		t.Errorf("Default wait time = %d, want 2", opts.WaitTime)
	}
	
	if !opts.PrintBackground {
		t.Error("Default PrintBackground should be true")
	}
	
	if opts.ViewportWidth != 1920 {
		t.Errorf("Default viewport width = %d, want 1920", opts.ViewportWidth)
	}
}