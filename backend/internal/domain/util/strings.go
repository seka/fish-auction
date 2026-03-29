package util

// IsPrintableASCII returns true if the rune is a printable ASCII character (0x20-0x7E).
func IsPrintableASCII(r rune) bool {
	return r >= 32 && r <= 126
}
