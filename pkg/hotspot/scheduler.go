package hotspot

import "fmt"

// CreateExpiryScheduler creates scheduler for expiry monitoring
func (c *Client) CreateExpiryScheduler(profileName string) error {
	schedulerName := "monitor-" + profileName

	// Build scheduler script
	schedulerScript := BuildExpirySchedulerScript(profileName)

	_, err := c.execute("/system/scheduler/add",
		"=name="+schedulerName,
		"=interval=5m",
		"=start-time=startup",
		"=policy=read,write,policy,test",
		"=on-event="+schedulerScript,
	)

	if err != nil {
		return NewError("create scheduler", err)
	}

	return nil
}

// RemoveExpiryScheduler removes expiry scheduler
func (c *Client) RemoveExpiryScheduler(profileName string) error {
	schedulerName := "monitor-" + profileName

	// Find scheduler
	reply, err := c.execute("/system/scheduler/print", "?name="+schedulerName)
	if err != nil {
		return NewError("remove scheduler", err)
	}

	if len(reply.Re) == 0 {
		// Scheduler not found, ignore
		return nil
	}

	schedulerID := reply.Re[0].Map[".id"]

	// Remove scheduler
	_, err = c.execute("/system/scheduler/remove", "=.id="+schedulerID)
	if err != nil {
		return NewError("remove scheduler", err)
	}

	return nil
}

// GetAllSchedulers retrieves all schedulers
func (c *Client) GetAllSchedulers() ([]Scheduler, error) {
	reply, err := c.execute("/system/scheduler/print")
	if err != nil {
		return nil, NewError("get all schedulers", err)
	}

	schedulers := make([]Scheduler, 0, len(reply.Re))
	for _, re := range reply.Re {
		scheduler := Scheduler{
			Name:      re.Map["name"],
			Interval:  re.Map["interval"],
			StartTime: re.Map["start-time"],
			Policy:    re.Map["policy"],
			OnEvent:   re.Map["on-event"],
			Enabled:   re.Map["disabled"] != "true",
		}
		schedulers = append(schedulers, scheduler)
	}

	return schedulers, nil
}

// GetSchedulerByName retrieves scheduler by name
func (c *Client) GetSchedulerByName(name string) (*Scheduler, error) {
	reply, err := c.execute("/system/scheduler/print", "?name="+name)
	if err != nil {
		return nil, NewError("get scheduler", err)
	}

	if len(reply.Re) == 0 {
		return nil, fmt.Errorf("scheduler not found: %s", name)
	}

	re := reply.Re[0]
	scheduler := &Scheduler{
		Name:      re.Map["name"],
		Interval:  re.Map["interval"],
		StartTime: re.Map["start-time"],
		Policy:    re.Map["policy"],
		OnEvent:   re.Map["on-event"],
		Enabled:   re.Map["disabled"] != "true",
	}

	return scheduler, nil
}
