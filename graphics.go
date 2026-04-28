package main

// Our application can run in one of four graphic modes:
// - ASCII
// - Unicode (default)
// - Iconic
// - Emoji
//
// For each, we'll use/allow a different set of characters.
// This should be transparent and the rest of our app doesn't even need to know about it.
// I hope...

// Here's a simple example showing what I want to do.
// It's implementation will certainly change as we add more stuff.
// Maybe some features will only be available if we're in certain graphics mode.
// (i.e. we can have an "add emoji" dialog only if emoji is supported, etc).

// TODO: Rename this or `Glyph` (see conifg.go).
// Glyphs will hold the allowed glyphs.
type Glyphs struct {
	// TODO: strings or runes?
	Cursor string
	More   string
	Less   string
	Bullet string
}

func GetGlyphs(mode Glyph) Glyphs {
	// Set `glyphs` with only ASCII characters (start with the most restrictive set).
	// Update on the swith according to the selected option `mode`.
	// TODO: Select appropriate symbols for each set.
	glyphs := Glyphs{
		// TODO: Consider removing the extra space at the end.
		// (It's there for full-width symbols to overflow).
		// (It also aligns better with full-width characters).
		Cursor: "> ",
		More:   "+ ",
		Less:   "- ",
		Bullet: "* ",
	}

	switch mode {

	case Unicode:
		glyphs.Cursor = "🠊 "
		glyphs.More = "＋"
		glyphs.Less = "－"
		glyphs.Bullet = " •"

	// TODO: Set the glyphs for each set.  BTW since Iconic and Emoji both imply the use of
	// Unicode, we can safely skip this for now.
	case Iconic:
	case Emoji:
	case ASCII:
		// Don't need to change anything.  `glyphs` is already set with ASCII only characters.
	}

	return glyphs
}
