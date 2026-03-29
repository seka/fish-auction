package stringutil

// IsPrintableASCII returns true if the rune is a printable ASCII character (0x20-0x7E).
func IsPrintableASCII(r rune) bool {
	return r >= 32 && r <= 126
}

// IsNonPrintableASCII returns true if the rune is NOT a printable ASCII character.
// This is useful for strings.ContainsFunc to find invalid characters.
func IsNonPrintableASCII(r rune) bool {
	return !IsPrintableASCII(r)
}
