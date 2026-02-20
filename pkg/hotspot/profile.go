package hotspot

import "fmt"

// CreateProfile creates new hotspot user profile on RouterOS
func (c *Client) CreateProfile(profile *Profile) error {
	// Validate
	if profile.Name == "" {
		return NewError("create profile", fmt.Errorf("profile name required"))
	}

	// Build on-login script
	profile.OnLoginScript = BuildOnLoginScript(profile)

	_, err := c.execute("/ip/hotspot/user/profile/add",
		"=name="+profile.Name,
		"=shared-users="+fmt.Sprintf("%d", profile.SharedUsers),
		"=rate-limit="+profile.RateLimit,
		"=on-login="+profile.OnLoginScript,
		"=keepalive-timeout="+profile.KeepaliveTimeout,
	)

	if err != nil {
		return NewError("create profile", err)
	}

	return nil
}

// UpdateProfile updates existing profile on RouterOS
func (c *Client) UpdateProfile(profileName string, updates *Profile) error {
	// Get current profile
	current, err := c.GetProfile(profileName)
	if err != nil {
		return err
	}

	// Apply updates
	if updates.RateLimit != "" {
		current.RateLimit = updates.RateLimit
	}
	if updates.SharedUsers > 0 {
		current.SharedUsers = updates.SharedUsers
	}
	if updates.Validity != "" {
		current.Validity = updates.Validity
	}
	if updates.Price > 0 {
		current.Price = updates.Price
	}
	if updates.SellingPrice > 0 {
		current.SellingPrice = updates.SellingPrice
	}
	if updates.ExpiryMode != "" {
		current.ExpiryMode = updates.ExpiryMode
	}
	if updates.LockUser != "" {
		current.LockUser = updates.LockUser
	}
	if updates.KeepaliveTimeout != "" {
		current.KeepaliveTimeout = updates.KeepaliveTimeout
	}

	// Rebuild on-login script
	current.OnLoginScript = BuildOnLoginScript(current)

	// Update on RouterOS
	reply, err := c.execute("/ip/hotspot/user/profile/print", "?name="+profileName)
	if err != nil {
		return NewError("update profile", err)
	}

	if len(reply.Re) == 0 {
		return ErrProfileNotFound
	}

	profileID := reply.Re[0].Map[".id"]

	_, err = c.execute("/ip/hotspot/user/profile/set",
		"=.id="+profileID,
		"=rate-limit="+current.RateLimit,
		"=shared-users="+fmt.Sprintf("%d", current.SharedUsers),
		"=on-login="+current.OnLoginScript,
		"=keepalive-timeout="+current.KeepaliveTimeout,
	)

	if err != nil {
		return NewError("update profile", err)
	}

	return nil
}

// DeleteProfile deletes profile from RouterOS
func (c *Client) DeleteProfile(profileName string) error {
	reply, err := c.execute("/ip/hotspot/user/profile/print", "?name="+profileName)
	if err != nil {
		return NewError("delete profile", err)
	}

	if len(reply.Re) == 0 {
		return ErrProfileNotFound
	}

	profileID := reply.Re[0].Map[".id"]

	_, err = c.execute("/ip/hotspot/user/profile/remove", "=.id="+profileID)
	if err != nil {
		return NewError("delete profile", err)
	}

	// Remove associated scheduler if exists
	// Note: RemoveExpiryScheduler will be implemented in scheduler.go

	return nil
}

// GetProfile retrieves profile from RouterOS
func (c *Client) GetProfile(profileName string) (*Profile, error) {
	reply, err := c.execute("/ip/hotspot/user/profile/print", "?name="+profileName)
	if err != nil {
		return nil, NewError("get profile", err)
	}

	if len(reply.Re) == 0 {
		return nil, ErrProfileNotFound
	}

	re := reply.Re[0]
	profile := &Profile{
		Name:             re.Map["name"],
		SharedUsers:      parseInt(re.Map["shared-users"]),
		RateLimit:        re.Map["rate-limit"],
		KeepaliveTimeout: re.Map["keepalive-timeout"],
		OnLoginScript:    re.Map["on-login"],
	}

	// Parse on-login script
	parsed, err := ParseOnLoginScript(re.Map["on-login"])
	if err == nil {
		profile.ExpiryMode = parsed.ExpiryMode
		profile.Price = parsed.Price
		profile.Validity = parsed.Validity
		profile.SellingPrice = parsed.SellingPrice
		profile.LockUser = parsed.LockUser
	}

	return profile, nil
}

// GetAllProfiles retrieves all profiles from RouterOS
func (c *Client) GetAllProfiles() ([]Profile, error) {
	reply, err := c.execute("/ip/hotspot/user/profile/print")
	if err != nil {
		return nil, NewError("get all profiles", err)
	}

	profiles := make([]Profile, 0, len(reply.Re))
	for _, re := range reply.Re {
		profile := &Profile{
			Name:             re.Map["name"],
			SharedUsers:      parseInt(re.Map["shared-users"]),
			RateLimit:        re.Map["rate-limit"],
			KeepaliveTimeout: re.Map["keepalive-timeout"],
			OnLoginScript:    re.Map["on-login"],
		}

		// Parse on-login script
		parsed, _ := ParseOnLoginScript(re.Map["on-login"])
		if parsed != nil {
			profile.ExpiryMode = parsed.ExpiryMode
			profile.Price = parsed.Price
			profile.Validity = parsed.Validity
			profile.SellingPrice = parsed.SellingPrice
			profile.LockUser = parsed.LockUser
		}

		profiles = append(profiles, *profile)
	}

	return profiles, nil
}

// GetProfileSettings extracts price/validity from on-login script
func (c *Client) GetProfileSettings(profileName string) (*Profile, error) {
	return c.GetProfile(profileName)
}

// SyncProfileToRouter ensures profile exists on RouterOS with current settings
func (c *Client) SyncProfileToRouter(profile *Profile) error {
	_, err := c.GetProfile(profile.Name)
	if err != nil {
		if err == ErrProfileNotFound {
			// Create new profile
			return c.CreateProfile(profile)
		}
		return err
	}

	// Update existing profile
	return c.UpdateProfile(profile.Name, profile)
}
