package mikrotik

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	routeros "github.com/go-routeros/routeros/v3"
	"go.uber.org/zap"
)

// MikroTikClient manages RouterOS API connections per router.
type MikroTikClient struct {
	connections map[uint]*routeros.Client
	connInfo    map[uint]*MikroTikConnectionInfo
	mu          sync.RWMutex
	routerRepo  repositories.RouterRepository
}

// MikroTikConnectionInfo holds connection parameters.
type MikroTikConnectionInfo struct {
	Host     string
	Username string
	Password string
	Port     int
}

// MikroTikConnection is kept for backward compatibility with other files.
type MikroTikConnection = MikroTikConnectionInfo

func NewMikroTikClient(routerRepo repositories.RouterRepository) *MikroTikClient {
	return &MikroTikClient{
		connections: make(map[uint]*routeros.Client),
		connInfo:    make(map[uint]*MikroTikConnectionInfo),
		routerRepo:  routerRepo,
	}
}

// dial opens a new RouterOS API connection.
func (c *MikroTikClient) dial(info *MikroTikConnectionInfo) (*routeros.Client, error) {
	addr := fmt.Sprintf("%s:%d", info.Host, info.Port)
	return routeros.Dial(addr, info.Username, info.Password)
}

// Connect loads router config from DB and opens a connection.
func (c *MikroTikClient) Connect(routerID uint) (*MikroTikConnectionInfo, error) {
	router, err := c.routerRepo.FindByID(routerID)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}

	info := &MikroTikConnectionInfo{
		Host:     router.Host,
		Username: router.Username,
		Password: router.Password,
		Port:     router.Port,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.connInfo[routerID] = info

	// Try to establish connection
	client, err := c.dial(info)
	if err != nil {
		logger.Warn("MikroTik connection failed",
			zap.Uint("router_id", routerID),
			zap.String("host", info.Host),
			zap.Error(err),
		)
		return info, nil // Return info even if connection failed; will retry on use
	}

	c.connections[routerID] = client
	logger.Info("MikroTik connected",
		zap.Uint("router_id", routerID),
		zap.String("host", info.Host),
	)
	return info, nil
}

// getClient returns a live RouterOS client, re-dialing if necessary.
func (c *MikroTikClient) getClient(routerID uint) (*routeros.Client, error) {
	c.mu.RLock()
	client, exists := c.connections[routerID]
	info := c.connInfo[routerID]
	c.mu.RUnlock()

	if exists && client != nil {
		return client, nil
	}

	if info == nil {
		// Not connected yet — load from DB
		router, err := c.routerRepo.FindByID(routerID)
		if err != nil {
			return nil, fmt.Errorf("router not found: %w", err)
		}
		info = &MikroTikConnectionInfo{
			Host:     router.Host,
			Username: router.Username,
			Password: router.Password,
			Port:     router.Port,
		}
		c.mu.Lock()
		c.connInfo[routerID] = info
		c.mu.Unlock()
	}

	client, err := c.dial(info)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router %d (%s): %w", routerID, info.Host, err)
	}

	c.mu.Lock()
	c.connections[routerID] = client
	c.mu.Unlock()

	return client, nil
}

// GetClient returns the connection info (for backward compat).
func (c *MikroTikClient) GetClient(routerID uint) (*MikroTikConnectionInfo, error) {
	c.mu.RLock()
	info, exists := c.connInfo[routerID]
	c.mu.RUnlock()
	if exists {
		return info, nil
	}
	return c.Connect(routerID)
}

func (c *MikroTikClient) GetActiveRouter() (*MikroTikConnectionInfo, uint, error) {
	router, err := c.routerRepo.FindActive()
	if err != nil {
		return nil, 0, fmt.Errorf("no active router found: %w", err)
	}
	info, err := c.GetClient(router.ID)
	return info, router.ID, err
}

func (c *MikroTikClient) ConnectAll() error {
	routers, err := c.routerRepo.FindAll()
	if err != nil {
		return err
	}
	for _, router := range routers {
		if _, err := c.Connect(router.ID); err != nil {
			logger.Warn("Failed to pre-connect to router",
				zap.Uint("router_id", router.ID),
				zap.Error(err),
			)
		}
	}
	logger.Info("MikroTik connect-all finished", zap.Int("count", len(routers)))
	return nil
}

func (c *MikroTikClient) HealthCheck(routerID uint) error {
	client, err := c.getClient(routerID)
	if err != nil {
		return err
	}
	_, err = client.Run("/system/resource/print")
	return err
}

// ─── PPPoE Secrets ───────────────────────────────────────────────────────────

func (c *MikroTikClient) AddUser(routerID uint, username, password, profile string) error {
	client, err := c.getClient(routerID)
	if err != nil {
		return err
	}
	_, err = client.Run("/ppp/secret/add",
		"=name="+username,
		"=password="+password,
		"=profile="+profile,
		"=service=pppoe",
	)
	if err != nil {
		return fmt.Errorf("AddUser failed: %w", err)
	}
	logger.Info("MikroTik: AddUser ok", zap.String("username", username), zap.Uint("router_id", routerID))
	return nil
}

func (c *MikroTikClient) RemoveUser(routerID uint, username string) error {
	client, err := c.getClient(routerID)
	if err != nil {
		return err
	}
	// Find the internal ID first
	reply, err := client.Run("/ppp/secret/print", "?name="+username)
	if err != nil {
		return fmt.Errorf("RemoveUser find failed: %w", err)
	}
	if len(reply.Re) == 0 {
		return fmt.Errorf("PPPoE user not found: %s", username)
	}
	id := reply.Re[0].Map[".id"]
	_, err = client.Run("/ppp/secret/remove", "=.id="+id)
	if err != nil {
		return fmt.Errorf("RemoveUser failed: %w", err)
	}
	logger.Info("MikroTik: RemoveUser ok", zap.String("username", username))
	return nil
}

func (c *MikroTikClient) UpdateUser(routerID uint, username string, args ...string) error {
	client, err := c.getClient(routerID)
	if err != nil {
		return err
	}
	reply, err := client.Run("/ppp/secret/print", "?name="+username)
	if err != nil || len(reply.Re) == 0 {
		return fmt.Errorf("UpdateUser: user not found: %s", username)
	}
	id := reply.Re[0].Map[".id"]
	setArgs := []string{"/ppp/secret/set", "=.id=" + id}
	setArgs = append(setArgs, args...)
	_, err = client.Run(setArgs...)
	if err != nil {
		return fmt.Errorf("UpdateUser failed: %w", err)
	}
	return nil
}

func (c *MikroTikClient) GetAllUsers(routerID uint) ([]PPPoEUser, error) {
	client, err := c.getClient(routerID)
	if err != nil {
		return nil, err
	}
	reply, err := client.Run("/ppp/secret/print")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers failed: %w", err)
	}
	users := make([]PPPoEUser, 0, len(reply.Re))
	for _, re := range reply.Re {
		users = append(users, PPPoEUser{
			Name:      re.Map["name"],
			Service:   re.Map["service"],
			Profile:   re.Map["profile"],
			CallerID:  re.Map["caller-id"],
			Disabled:  re.Map["disabled"] == "true",
			LastLogin: re.Map["last-logged-out"],
		})
	}
	return users, nil
}

func (c *MikroTikClient) GetActiveSessions(routerID uint) ([]ActiveSession, error) {
	client, err := c.getClient(routerID)
	if err != nil {
		return nil, err
	}
	reply, err := client.Run("/ppp/active/print")
	if err != nil {
		return nil, fmt.Errorf("GetActiveSessions failed: %w", err)
	}
	sessions := make([]ActiveSession, 0, len(reply.Re))
	for _, re := range reply.Re {
		sessions = append(sessions, ActiveSession{
			Name:     re.Map["name"],
			CallerID: re.Map["caller-id"],
			Address:  re.Map["address"],
			Uptime:   re.Map["uptime"],
			BytesIn:  re.Map["bytes-in"],
			BytesOut: re.Map["bytes-out"],
			Encoding: re.Map["encoding"],
		})
	}
	return sessions, nil
}

func (c *MikroTikClient) GetAllProfiles(routerID uint) ([]Profile, error) {
	client, err := c.getClient(routerID)
	if err != nil {
		return nil, err
	}
	reply, err := client.Run("/ppp/profile/print")
	if err != nil {
		return nil, fmt.Errorf("GetAllProfiles failed: %w", err)
	}
	profiles := make([]Profile, 0, len(reply.Re))
	for _, re := range reply.Re {
		profiles = append(profiles, Profile{
			Name:         re.Map["name"],
			RateLimit:    re.Map["rate-limit"],
			LocalAddress: re.Map["local-address"],
			OnlyOne:      re.Map["only-one"] == "yes",
		})
	}
	return profiles, nil
}

func (c *MikroTikClient) DisconnectUser(routerID uint, username string) error {
	client, err := c.getClient(routerID)
	if err != nil {
		return err
	}
	reply, err := client.Run("/ppp/active/print", "?name="+username)
	if err != nil {
		return fmt.Errorf("DisconnectUser find failed: %w", err)
	}
	if len(reply.Re) == 0 {
		return fmt.Errorf("no active session for user: %s", username)
	}
	id := reply.Re[0].Map[".id"]
	_, err = client.Run("/ppp/active/remove", "=.id="+id)
	if err != nil {
		return fmt.Errorf("DisconnectUser failed: %w", err)
	}
	logger.Info("MikroTik: DisconnectUser ok", zap.String("username", username))
	return nil
}

func (c *MikroTikClient) SetActiveProfile(routerID uint, username, profile string) error {
	client, err := c.getClient(routerID)
	if err != nil {
		return err
	}
	reply, err := client.Run("/ppp/secret/print", "?name="+username)
	if err != nil || len(reply.Re) == 0 {
		return fmt.Errorf("SetActiveProfile: user not found: %s", username)
	}
	id := reply.Re[0].Map[".id"]
	_, err = client.Run("/ppp/secret/set", "=.id="+id, "=profile="+profile)
	if err != nil {
		return fmt.Errorf("SetActiveProfile failed: %w", err)
	}
	logger.Info("MikroTik: SetActiveProfile ok",
		zap.String("username", username),
		zap.String("profile", profile),
	)
	return nil
}

// ─── Router Status ────────────────────────────────────────────────────────────

func (c *MikroTikClient) GetRouterStatus(routerID uint) (*RouterStatus, error) {
	c.mu.RLock()
	info := c.connInfo[routerID]
	c.mu.RUnlock()

	status := &RouterStatus{
		RouterID:  routerID,
		LastCheck: time.Now(),
	}

	if info != nil {
		status.Host = info.Host
	}

	client, err := c.getClient(routerID)
	if err != nil {
		status.Status = "disconnected"
		status.Error = err.Error()
		return status, nil
	}

	reply, err := client.Run("/system/resource/print")
	if err != nil || len(reply.Re) == 0 {
		status.Status = "error"
		if err != nil {
			status.Error = err.Error()
		}
		return status, nil
	}

	res := reply.Re[0]
	status.Status = "connected"
	status.Uptime = res.Map["uptime"]

	if cpuLoad, err := strconv.ParseFloat(res.Map["cpu-load"], 64); err == nil {
		status.CPU = cpuLoad
	}
	if totalMem, err := strconv.ParseInt(res.Map["total-memory"], 10, 64); err == nil {
		if freeMem, err2 := strconv.ParseInt(res.Map["free-memory"], 10, 64); err2 == nil {
			used := totalMem - freeMem
			if totalMem > 0 {
				status.Memory = int(used * 100 / totalMem)
			}
		}
	}

	// Count active PPPoE sessions
	activeReply, err := client.Run("/ppp/active/print", "count-only=")
	if err == nil && len(activeReply.Re) > 0 {
		if count, err := strconv.Atoi(activeReply.Re[0].Map["ret"]); err == nil {
			status.ActiveUsers = count
		}
	}

	return status, nil
}

func (c *MikroTikClient) GetAllRoutersStatus() ([]RouterStatus, error) {
	routers, err := c.routerRepo.FindAll()
	if err != nil {
		return nil, err
	}

	statuses := make([]RouterStatus, 0, len(routers))
	for _, router := range routers {
		s, _ := c.GetRouterStatus(router.ID)
		if s == nil {
			s = &RouterStatus{
				RouterID:  router.ID,
				Name:      router.Name,
				Host:      router.Host,
				Status:    "unknown",
				LastCheck: time.Now(),
			}
		}
		s.Name = router.Name
		statuses = append(statuses, *s)
	}
	return statuses, nil
}
