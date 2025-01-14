package logger

import (
	"io"
	"os"
	"testing"

	"github.com/aripalo/vegas-credentials/internal/config"
	"github.com/aripalo/vegas-credentials/internal/profile"
	"github.com/aripalo/vegas-credentials/internal/vegastestapp"
	"github.com/gookit/color"
)

type formatTestCase struct {
	description string
	flags       config.Flags
	emoji       string
	prefix      string
	message     string
	want        string
}

func TestFormat(t *testing.T) {

	tests := []formatTestCase{
		{
			"with all formatting",
			config.Flags{
				NoColor: false,
			},
			"🚧",
			"Test",
			"Message",
			"🚧 \x1b[90m\x1b[1mTest:\x1b[0m\x1b[0m \x1b[90mMessage\x1b[0m",
		},
		{
			"witout emoji",
			config.Flags{
				NoColor: false,
			},
			"",
			"Test",
			"Message",
			"\x1b[90m\x1b[1mTest:\x1b[0m\x1b[0m \x1b[90mMessage\x1b[0m",
		},
		{
			"witout prefix",
			config.Flags{
				NoColor: false,
			},
			"🚧",
			"",
			"Message",
			"🚧 \x1b[90mMessage\x1b[0m",
		},
		{
			"without color",
			config.Flags{
				NoColor: true,
			},
			"🚧",
			"Test",
			"Message",
			"Test: Message",
		},
	}

	for _, tc := range tests {

		// Handle terminal env (i.e. in CI)
		nocolor := os.Getenv("NO_COLOR")
		term := os.Getenv("TERM")
		level := color.TermColorLevel()
		os.Unsetenv("NO_COLOR")
		os.Setenv("TERM", "xterm-256color")
		os.Setenv("FORCE_COLOR", "on")
		_ = color.ForceSetColorLevel(color.Level256)

		defer func() {
			os.Setenv("NO_COLOR", nocolor)
			os.Setenv("TERM", term)
			os.Unsetenv("FORCE_COLOR")
			color.ForceSetColorLevel(level)
		}()

		t.Run(tc.description, func(t *testing.T) {

			a := &vegastestapp.AssumeAppForTesting{
				Flags:       tc.flags,
				Profile:     profile.Profile{},
				Destination: io.Discard,
			}

			got := format(a, textColorDebug, tc.emoji, tc.prefix, tc.message)

			if got != tc.want {
				t.Fatalf(`Got %q, want %q`, got, tc.want)
			}
		})
	}

}
