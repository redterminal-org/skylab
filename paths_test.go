package main

import "testing"

func TestAppPrefix(t *testing.T) {
	tests := []struct {
		prefix string
		want   string
	}{
		{prefix: "skylab", want: "/skylab"},
		{prefix: "misfin", want: "/misfin"},
		{prefix: "xyz", want: "/xyz"},
		{prefix: "/skylab/", want: "/skylab"},
		{prefix: "", want: ""},
	}

	for _, tt := range tests {
		if got := appPrefix(tt.prefix); got != tt.want {
			t.Errorf("appPrefix(%q) = %q, want %q", tt.prefix, got, tt.want)
		}
	}
}

func TestAppPath(t *testing.T) {
	tests := []struct {
		prefix string
		suffix string
		want   string
	}{
		{prefix: "skylab", suffix: "mail", want: "/skylab/mail"},
		{prefix: "misfin", suffix: "", want: "/misfin"},
		{prefix: "xyz", suffix: "compose/new", want: "/xyz/compose/new"},
		{prefix: "/skylab/", suffix: "/mail/", want: "/skylab/mail"},
		{prefix: "", suffix: "mail", want: "/mail"},
		{prefix: "", suffix: "", want: "/"},
	}

	for _, tt := range tests {
		if got := appPath(tt.prefix, tt.suffix); got != tt.want {
			t.Errorf("appPath(%q, %q) = %q, want %q", tt.prefix, tt.suffix, got, tt.want)
		}
	}
}
