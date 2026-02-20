package hotspot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	DefaultCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	MinUsername    = 4
	MaxUsername    = 12
	MinPassword    = 4
	MaxPassword    = 12
)

// GenerateRandomString generates random string from charset
func GenerateRandomString(length int, charset string) string {
	if charset == "" {
		charset = DefaultCharset
	}

	if length < MinUsername {
		length = MinUsername
	}

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// BuildOnLoginScript creates on-login script from profile parameters
func BuildOnLoginScript(profile *Profile) string {
	return fmt.Sprintf(
		`:local expmode "%s";:local price "%.2f";:local validity "%s";:local selling "%.2f";:local lock "%s";`,
		profile.ExpiryMode,
		profile.Price,
		profile.Validity,
		profile.SellingPrice,
		profile.LockUser,
	)
}

// ParseOnLoginScript extracts profile settings from on-login script
func ParseOnLoginScript(script string) (*Profile, error) {
	profile := &Profile{}

	parts := strings.Split(script, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "expmode") {
			profile.ExpiryMode = extractQuotedValue(part)
		} else if strings.Contains(part, "price") {
			profile.Price = parseFloatFromQuoted(part)
		} else if strings.Contains(part, "validity") {
			profile.Validity = extractQuotedValue(part)
		} else if strings.Contains(part, "selling") {
			profile.SellingPrice = parseFloatFromQuoted(part)
		} else if strings.Contains(part, "lock") {
			profile.LockUser = extractQuotedValue(part)
		}
	}

	return profile, nil
}

// GetUserMode determines if user is voucher (vc) or user-pass (up)
func GetUserMode(username, password string) string {
	if username == password {
		return "vc-"
	}
	return "up-"
}

// BuildSaleScriptName creates script name for sales record
func BuildSaleScriptName(sale *Sale) string {
	// Format: date-|-time-|-username-|-price-|-address-|-mac-|-validity
	scriptName := fmt.Sprintf("%s-|-%s-|-%s-|-%.2f-|-",
		sale.Date,
		sale.Time,
		sale.Username,
		sale.Price,
	)

	if sale.Address != "" {
		scriptName += sale.Address + "-|"
	}

	if sale.Mac != "" {
		scriptName += sale.Mac + "-|"
	}

	if sale.Validity != "" {
		scriptName += sale.Validity
	}

	return scriptName
}

// ParseSaleScriptName parses script name back to Sale struct
func ParseSaleScriptName(scriptName string) (*Sale, error) {
	if !strings.Contains(scriptName, "-|-") {
		return nil, fmt.Errorf("invalid sales script format")
	}

	parts := strings.Split(scriptName, "-|-")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid sales script format: insufficient parts")
	}

	sale := &Sale{
		Date:     parts[0],
		Time:     parts[1],
		Username: parts[2],
		Price:    parseFloat(parts[3]),
	}

	if len(parts) > 4 && parts[4] != "" {
		sale.Address = parts[4]
	}
	if len(parts) > 5 && parts[5] != "" {
		sale.Mac = parts[5]
	}
	if len(parts) > 6 && parts[6] != "" {
		sale.Validity = parts[6]
	}

	return sale, nil
}

// BuildExpirySchedulerScript creates scheduler script for expiry monitoring
func BuildExpirySchedulerScript(profileName string) string {
	return fmt.Sprintf(
		`:local profile "%s"; /ip hotspot user { :foreach i in=[find profile="$profile"] do={ :local cmt [get $i comment]; :if ([:len $cmt] > 0 && $cmt ~ "^[a-z]{3}/[0-9]{2}/[0-9]{4}") do={ :local expdate [:pick $cmt 0 [:find $cmt " "]]; :local now [/system clock get date]; :if ($expdate < $now) do={ remove $i; }}}}`,
		profileName,
	)
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

	duration := time.Duration(seconds) * time.Second
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// Helper functions
func extractQuotedValue(s string) string {
	start := strings.Index(s, `"`)
	end := strings.LastIndex(s, `"`)
	if start >= 0 && end > start {
		return s[start+1 : end]
	}
	return ""
}

func parseFloatFromQuoted(s string) float64 {
	val := extractQuotedValue(s)
	result := 0.0
	fmt.Sscanf(val, "%f", &result)
	return result
}

func parseFloat(s string) float64 {
	result := 0.0
	fmt.Sscanf(s, "%f", &result)
	return result
}

func boolToString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
