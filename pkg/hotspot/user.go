package hotspot

import (
	"fmt"
	"strings"
	"time"
)

// CreateUser creates new hotspot user on RouterOS
func (c *Client) CreateUser(user *User) error {
	// Validate
	if user.Name == "" || user.Password == "" {
		return NewError("create user", fmt.Errorf("username and password required"))
	}

	if user.Profile == "" {
		return NewError("create user", fmt.Errorf("profile required"))
	}

	// Build comment
	if user.Comment == "" {
		user.Comment = GetUserMode(user.Name, user.Password) + extractPrefix(user.Name)
	}

	if user.Server == "" {
		user.Server = "all"
	}

	_, err := c.execute("/ip/hotspot/user/add",
		"=server="+user.Server,
		"=name="+user.Name,
		"=password="+user.Password,
		"=profile="+user.Profile,
		"=disabled="+boolToString(user.Disabled),
		"=limit-uptime="+fmt.Sprintf("%d", user.LimitUptime),
		"=limit-bytes-total="+fmt.Sprintf("%d", user.LimitBytesTotal),
		"=comment="+user.Comment,
	)

	if err != nil {
		return NewError("create user", err)
	}

	return nil
}

// UpdateUser updates existing user on RouterOS
func (c *Client) UpdateUser(username string, updates map[string]interface{}) error {
	// Find user
	reply, err := c.execute("/ip/hotspot/user/print", "?name="+username)
	if err != nil {
		return NewError("update user", err)
	}

	if len(reply.Re) == 0 {
		return ErrUserNotFound
	}

	userID := reply.Re[0].Map[".id"]

	// Build update arguments
	args := []string{"=.id=" + userID}

	if profile, ok := updates["profile"].(string); ok && profile != "" {
		args = append(args, "=profile="+profile)
	}
	if disabled, ok := updates["disabled"].(bool); ok {
		args = append(args, "=disabled="+boolToString(disabled))
	}
	if comment, ok := updates["comment"].(string); ok && comment != "" {
		args = append(args, "=comment="+comment)
	}
	if uptime, ok := updates["limit_uptime"].(int64); ok {
		args = append(args, "=limit-uptime="+fmt.Sprintf("%d", uptime))
	}
	if bytesTotal, ok := updates["limit_bytes_total"].(int64); ok {
		args = append(args, "=limit-bytes-total="+fmt.Sprintf("%d", bytesTotal))
	}

	_, err = c.execute("/ip/hotspot/user/set", args...)
	if err != nil {
		return NewError("update user", err)
	}

	return nil
}

// DeleteUser deletes user from RouterOS
func (c *Client) DeleteUser(username string) error {
	reply, err := c.execute("/ip/hotspot/user/print", "?name="+username)
	if err != nil {
		return NewError("delete user", err)
	}

	if len(reply.Re) == 0 {
		return ErrUserNotFound
	}

	userID := reply.Re[0].Map[".id"]

	_, err = c.execute("/ip/hotspot/user/remove", "=.id="+userID)
	if err != nil {
		return NewError("delete user", err)
	}

	return nil
}

// GetUser retrieves user from RouterOS
func (c *Client) GetUser(username string) (*User, error) {
	reply, err := c.execute("/ip/hotspot/user/print", "?name="+username)
	if err != nil {
		return nil, NewError("get user", err)
	}

	if len(reply.Re) == 0 {
		return nil, ErrUserNotFound
	}

	re := reply.Re[0]
	user := &User{
		Name:            re.Map["name"],
		Password:        re.Map["password"],
		Profile:         re.Map["profile"],
		Comment:         re.Map["comment"],
		LimitUptime:     parseInt64(re.Map["limit-uptime"]),
		LimitBytesTotal: parseInt64(re.Map["limit-bytes-total"]),
		LimitBytesIn:    parseInt64(re.Map["limit-bytes-in"]),
		LimitBytesOut:   parseInt64(re.Map["limit-bytes-out"]),
		Disabled:        re.Map["disabled"] == "true",
		Uptime:          re.Map["uptime"],
		BytesIn:         re.Map["bytes-in"],
		BytesOut:        re.Map["bytes-out"],
	}

	return user, nil
}

// GetAllUsers retrieves all users from RouterOS
func (c *Client) GetAllUsers(filter *UserFilter) ([]User, error) {
	var args []string

	if filter != nil {
		if filter.Profile != "" {
			args = append(args, "?profile="+filter.Profile)
		}
		if filter.Comment != "" {
			args = append(args, "?comment="+filter.Comment)
		}
		if filter.Disabled != nil {
			disabled := boolToString(*filter.Disabled)
			args = append(args, "?disabled="+disabled)
		}
	}

	reply, err := c.execute("/ip/hotspot/user/print", args...)
	if err != nil {
		return nil, NewError("get all users", err)
	}

	users := make([]User, 0, len(reply.Re))
	for _, re := range reply.Re {
		user := User{
			Name:            re.Map["name"],
			Password:        re.Map["password"],
			Profile:         re.Map["profile"],
			Comment:         re.Map["comment"],
			LimitUptime:     parseInt64(re.Map["limit-uptime"]),
			LimitBytesTotal: parseInt64(re.Map["limit-bytes-total"]),
			LimitBytesIn:    parseInt64(re.Map["limit-bytes-in"]),
			LimitBytesOut:   parseInt64(re.Map["limit-bytes-out"]),
			Disabled:        re.Map["disabled"] == "true",
			Uptime:          re.Map["uptime"],
			BytesIn:         re.Map["bytes-in"],
			BytesOut:        re.Map["bytes-out"],
		}
		users = append(users, user)
	}

	// Apply pagination
	if filter != nil && filter.Limit > 0 {
		start := filter.Offset
		if start > len(users) {
			start = len(users)
		}
		end := start + filter.Limit
		if end > len(users) {
			end = len(users)
		}
		users = users[start:end]
	}

	return users, nil
}

// GetUsersByProfile retrieves users filtered by profile
func (c *Client) GetUsersByProfile(profile string) ([]User, error) {
	return c.GetAllUsers(&UserFilter{Profile: profile})
}

// GetUsersByComment retrieves users filtered by comment
func (c *Client) GetUsersByComment(comment string) ([]User, error) {
	return c.GetAllUsers(&UserFilter{Comment: comment})
}

// DisableUser disables user on RouterOS
func (c *Client) DisableUser(username string) error {
	return c.UpdateUser(username, map[string]interface{}{
		"disabled": true,
	})
}

// EnableUser enables user on RouterOS
func (c *Client) EnableUser(username string) error {
	return c.UpdateUser(username, map[string]interface{}{
		"disabled": false,
	})
}

// RemoveExpiredUsers removes users with expired dates in comment
func (c *Client) RemoveExpiredUsers(profile string) (int, error) {
	users, err := c.GetAllUsers(&UserFilter{Profile: profile})
	if err != nil {
		return 0, err
	}

	removedCount := 0

	for _, user := range users {
		// Check if comment contains expiry date
		if strings.Contains(user.Comment, " ") {
			parts := strings.Split(user.Comment, " ")
			expiryDate := parts[0]

			// Parse date
			expiryTime, err := time.Parse("Jan/02/2006", expiryDate)
			if err == nil && expiryTime.Before(time.Now()) {
				// Remove user
				err := c.DeleteUser(user.Name)
				if err == nil {
					removedCount++
				}
			}
		}
	}

	return removedCount, nil
}

// RemoveUnusedVouchers removes vouchers with uptime=0s
func (c *Client) RemoveUnusedVouchers(profile string) (int, error) {
	users, err := c.GetAllUsers(&UserFilter{Profile: profile})
	if err != nil {
		return 0, err
	}

	removedCount := 0
	for _, user := range users {
		if user.Uptime == "0s" {
			err := c.DeleteUser(user.Name)
			if err == nil {
				removedCount++
			}
		}
	}

	return removedCount, nil
}

// BatchCreateUsers creates multiple users
func (c *Client) BatchCreateUsers(users []User) (*VoucherResult, error) {
	result := &VoucherResult{
		Vouchers: make([]User, 0),
		Errors:   make([]string, 0),
	}

	for _, user := range users {
		err := c.CreateUser(&user)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("%s: %v", user.Name, err))
		} else {
			result.Success++
			result.Vouchers = append(result.Vouchers, user)
		}
	}

	return result, nil
}

// BatchRemoveUsers removes multiple users
func (c *Client) BatchRemoveUsers(usernames []string) (int, error) {
	removedCount := 0
	for _, username := range usernames {
		err := c.DeleteUser(username)
		if err == nil {
			removedCount++
		}
	}
	return removedCount, nil
}

// Helper functions
func extractPrefix(username string) string {
	parts := strings.Split(username, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
