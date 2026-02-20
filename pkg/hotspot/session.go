package hotspot

// GetActiveSessions retrieves all active hotspot sessions
func (c *Client) GetActiveSessions() ([]Session, error) {
	reply, err := c.execute("/ip/hotspot/active/print")
	if err != nil {
		return nil, NewError("get active sessions", err)
	}

	sessions := make([]Session, 0, len(reply.Re))
	for _, re := range reply.Re {
		session := Session{
			Name:            re.Map["user"],
			Address:         re.Map["address"],
			MacAddress:      re.Map["mac-address"],
			Uptime:          re.Map["uptime"],
			SessionTimeLeft: re.Map["session-time-left"],
			BytesIn:         re.Map["bytes-in"],
			BytesOut:        re.Map["bytes-out"],
			LoginBy:         re.Map["login-by"],
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetSessionsByServer retrieves active sessions for specific server
func (c *Client) GetSessionsByServer(server string) ([]Session, error) {
	reply, err := c.execute("/ip/hotspot/active/print", "?server="+server)
	if err != nil {
		return nil, NewError("get sessions by server", err)
	}

	sessions := make([]Session, 0, len(reply.Re))
	for _, re := range reply.Re {
		session := Session{
			Name:            re.Map["user"],
			Address:         re.Map["address"],
			MacAddress:      re.Map["mac-address"],
			Uptime:          re.Map["uptime"],
			SessionTimeLeft: re.Map["session-time-left"],
			BytesIn:         re.Map["bytes-in"],
			BytesOut:        re.Map["bytes-out"],
			LoginBy:         re.Map["login-by"],
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetSessionByUsername retrieves active session by username
func (c *Client) GetSessionByUsername(username string) (*Session, error) {
	reply, err := c.execute("/ip/hotspot/active/print", "?user="+username)
	if err != nil {
		return nil, NewError("get session by username", err)
	}

	if len(reply.Re) == 0 {
		return nil, ErrSessionNotFound
	}

	re := reply.Re[0]
	session := &Session{
		Name:            re.Map["user"],
		Address:         re.Map["address"],
		MacAddress:      re.Map["mac-address"],
		Uptime:          re.Map["uptime"],
		SessionTimeLeft: re.Map["session-time-left"],
		BytesIn:         re.Map["bytes-in"],
		BytesOut:        re.Map["bytes-out"],
		LoginBy:         re.Map["login-by"],
	}

	return session, nil
}

// DisconnectUser disconnects active user session
func (c *Client) DisconnectUser(username string) error {
	// Find user in active sessions
	reply, err := c.execute("/ip/hotspot/active/print", "?user="+username)
	if err != nil {
		return NewError("disconnect user", err)
	}

	if len(reply.Re) == 0 {
		return ErrSessionNotFound
	}

	sessionID := reply.Re[0].Map[".id"]

	_, err = c.execute("/ip/hotspot/active/remove", "=.id="+sessionID)
	if err != nil {
		return NewError("disconnect user", err)
	}

	return nil
}

// GetSessionStats retrieves session statistics
func (c *Client) GetSessionStats() (*SessionStats, error) {
	// Get active sessions
	reply, err := c.execute("/ip/hotspot/active/print")
	if err != nil {
		return nil, NewError("get session stats", err)
	}

	totalBytesIn := int64(0)
	totalBytesOut := int64(0)

	for _, re := range reply.Re {
		totalBytesIn += parseInt64(re.Map["bytes-in"])
		totalBytesOut += parseInt64(re.Map["bytes-out"])
	}

	return &SessionStats{
		TotalUsers:    len(reply.Re),
		ActiveUsers:   len(reply.Re),
		TotalBytesIn:  FormatBytes(totalBytesIn),
		TotalBytesOut: FormatBytes(totalBytesOut),
	}, nil
}
