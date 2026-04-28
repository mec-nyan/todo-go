package main

import "flag"

type Glyph int

const (
	ASCII Glyph = iota
	Unicode
	Iconic // Requires a nerd font or similar iconic font.
	Emoji  // Maybe this belongs somewhere else...
)

const (
	DefaultConfigFile = "~/.config/notes-go/notes.conf"
	DefaultSaveFile   = "~/.local/share/notes-go/notes.json"
)

// Options for our program.
type Options struct {
	// ConfigFile is our configuration file (usually `~/.config/notes-go/notes.conf`).
	ConfigFile string
	// SaveFile is where we're gonna save our notes.
	// Ideally this will be `~/.local/share/notes-go/notes.json`.  We can also consider saving our
	// files under `~` with a `dot` name (i.e. .notes-go.json and .notes-go.cache.json).
	SaveFile string
	// FullScreen tells us wether to run our app inline or in the alternate buffer.
	FullScreen bool
	// Graphics defines what character set we use to draw our app on the screen:
	// - ASCII (safe but old, we can do better).
	// - Unicode (supported by every modern teminal emulator).
	// - Iconic (requires an iconic font installed).
	// - Emoji (wether to use emojis or not. Implies Unicode).
	Graphics Glyph
}

func LoadConfig() Options {
	// TODO: Check the config dir.  If it exists, load configuration from there.
	panic("To be implemented!")
}

func parseArgs() Options {
	// TODO: Mayebe use a package other than `flag` for more refined argument parsing
	// i.e. support long and short options, etc.
	configFile := flag.String("c", DefaultConfigFile, "configuration directory")
	saveFile := flag.String("f", DefaultSaveFile, "file to read/write our notes")
	fullScreen := flag.Bool("F", false, "open in full screen mode")
	// These flags supersedes each other, in the following order: ascii -> unicode -> icons -> emoji.
	// If neither icons nor emoji are set, unicode glyphs will be used.  ASCII will only be used if
	// any of the other flags is unset.  By default I'll use Unicode symbols without emoji.
	// BTW:
	// - ASCII is the most restrictive, "safe" set (yet we didn't need to restrict to ASCII for a
	// long time now...).
	// - Unicode implies ASCII (or ASCII is contained in Unicode).
	// - Emojis implies the use of Unicode since it's part of it.
	// - Iconic is the only "special" one since it requires an iconic font to be installed
	// (it uses Unicode's private area).
	ascii := flag.Bool("ascii", false, "restrict to ASCII characters")
	unicode := flag.Bool("unicode", false, "use Unicode")
	icons := flag.Bool("icons", false, "show icons")
	emoji := flag.Bool("emoji", false, "show emoji")

	flag.Parse()

	// If not set use Unicode by default.
	graphics := Unicode

	// Level down to ASCII.
	if *ascii {
		graphics = ASCII
	}

	if *unicode {
		graphics = Unicode
	}

	if *icons {
		graphics = Iconic
	}

	if *emoji {
		graphics = Emoji
	}

	return Options{
		ConfigFile: *configFile,
		SaveFile:   *saveFile,
		FullScreen: *fullScreen,
		Graphics:   graphics,
	}
}
