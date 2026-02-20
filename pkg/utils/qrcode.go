package utils

import (
	"fmt"
)

// GenerateQRCodeURL generates QR code URL using external API service
func GenerateQRCodeURL(data string, size int) (string, error) {
	return "https://api.qrserver.com/v1/create-qr-code/?data=" + data + "&size=" + fmt.Sprintf("%dx%d", size, size), nil
}

// FormatBytes formats bytes to human readable string
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatSeconds formats seconds to human readable duration
func FormatSeconds(seconds int64) string {
	if seconds == 0 {
		return "unlimited"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
