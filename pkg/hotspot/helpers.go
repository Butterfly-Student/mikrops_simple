package hotspot

import "fmt"

// parseInt parses string to int
func parseInt(s string) int {
	result := 0
	fmt.Sscanf(s, "%d", &result)
	return result
}

// parseInt64 parses string to int64
func parseInt64(s string) int64 {
	result := int64(0)
	fmt.Sscanf(s, "%d", &result)
	return result
}
