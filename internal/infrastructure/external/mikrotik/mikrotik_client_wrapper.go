package mikrotik

import (
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/pkg/hotspot"
)

// HotspotClientWrapper wraps pkg/hotspot library
type HotspotClientWrapper struct {
	mtClient   *MikroTikClient
	routerRepo repositories.RouterRepository
}

// NewHotspotClientWrapper creates new wrapper
func NewHotspotClientWrapper(mtClient *MikroTikClient, routerRepo repositories.RouterRepository) *HotspotClientWrapper {
	return &HotspotClientWrapper{
		mtClient:   mtClient,
		routerRepo: routerRepo,
	}
}

// getHotspotClient gets hotspot client for specific router
func (w *HotspotClientWrapper) getHotspotClient(routerID uint) (*hotspot.Client, error) {
	// Get router from repository
	_, err := w.routerRepo.FindByID(routerID)
	if err != nil {
		return nil, err
	}

	// Get MikroTik client (from existing MikroTikClient)
	mtClient, err := w.mtClient.getClient(routerID)
	if err != nil {
		return nil, err
	}

	// Create hotspot client
	hotspotClient := hotspot.NewClient(routerID, mtClient)
	return hotspotClient, nil
}

// Wrapper methods delegate to hotspot.Client

// Profile Methods
func (w *HotspotClientWrapper) CreateProfile(routerID uint, profile *hotspot.Profile) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.CreateProfile(profile)
}

func (w *HotspotClientWrapper) UpdateProfile(routerID uint, profileName string, profile *hotspot.Profile) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.UpdateProfile(profileName, profile)
}

func (w *HotspotClientWrapper) DeleteProfile(routerID uint, profileName string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.DeleteProfile(profileName)
}

func (w *HotspotClientWrapper) GetProfile(routerID uint, profileName string) (*hotspot.Profile, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetProfile(profileName)
}

func (w *HotspotClientWrapper) GetAllProfiles(routerID uint) ([]hotspot.Profile, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetAllProfiles()
}

// User Methods
func (w *HotspotClientWrapper) CreateUser(routerID uint, user *hotspot.User) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.CreateUser(user)
}

func (w *HotspotClientWrapper) UpdateUser(routerID uint, username string, updates map[string]interface{}) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.UpdateUser(username, updates)
}

func (w *HotspotClientWrapper) DeleteUser(routerID uint, username string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.DeleteUser(username)
}

func (w *HotspotClientWrapper) GetUser(routerID uint, username string) (*hotspot.User, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetUser(username)
}

func (w *HotspotClientWrapper) GetAllUsers(routerID uint, filter *hotspot.UserFilter) ([]hotspot.User, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetAllUsers(filter)
}

func (w *HotspotClientWrapper) GetUsersByProfile(routerID uint, profile string) ([]hotspot.User, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetUsersByProfile(profile)
}

func (w *HotspotClientWrapper) GetUsersByComment(routerID uint, comment string) ([]hotspot.User, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetUsersByComment(comment)
}

func (w *HotspotClientWrapper) DisableUser(routerID uint, username string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.DisableUser(username)
}

func (w *HotspotClientWrapper) EnableUser(routerID uint, username string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.EnableUser(username)
}

func (w *HotspotClientWrapper) RemoveExpiredUsers(routerID uint, profile string) (int, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return 0, err
	}
	return client.RemoveExpiredUsers(profile)
}

func (w *HotspotClientWrapper) RemoveUnusedVouchers(routerID uint, profile string) (int, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return 0, err
	}
	return client.RemoveUnusedVouchers(profile)
}

func (w *HotspotClientWrapper) BatchCreateUsers(routerID uint, users []hotspot.User) (*hotspot.VoucherResult, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.BatchCreateUsers(users)
}

func (w *HotspotClientWrapper) BatchRemoveUsers(routerID uint, usernames []string) (int, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return 0, err
	}
	return client.BatchRemoveUsers(usernames)
}

// Voucher Methods
func (w *HotspotClientWrapper) GenerateVouchers(routerID uint, gen *hotspot.VoucherGenerator) (*hotspot.VoucherResult, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GenerateVouchers(gen)
}

func (w *HotspotClientWrapper) GenerateUserPasswordMode(routerID uint, gen *hotspot.VoucherGenerator) (*hotspot.VoucherResult, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GenerateUserPasswordMode(gen)
}

// Session Methods
func (w *HotspotClientWrapper) GetActiveSessions(routerID uint) ([]hotspot.Session, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetActiveSessions()
}

func (w *HotspotClientWrapper) GetSessionsByServer(routerID uint, server string) ([]hotspot.Session, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetSessionsByServer(server)
}

func (w *HotspotClientWrapper) GetSessionByUsername(routerID uint, username string) (*hotspot.Session, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetSessionByUsername(username)
}

func (w *HotspotClientWrapper) DisconnectUser(routerID uint, username string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.DisconnectUser(username)
}

func (w *HotspotClientWrapper) GetSessionStats(routerID uint) (*hotspot.SessionStats, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetSessionStats()
}

// Sales Methods
func (w *HotspotClientWrapper) RecordSale(routerID uint, sale *hotspot.Sale) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.RecordSale(sale)
}

func (w *HotspotClientWrapper) GetAllSales(routerID uint, filter *hotspot.SaleFilter) ([]hotspot.Sale, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetAllSales(filter)
}

func (w *HotspotClientWrapper) GetSalesByDateRange(routerID uint, startDate, endDate string) ([]hotspot.Sale, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetSalesByDateRange(startDate, endDate)
}

func (w *HotspotClientWrapper) GetSalesByPrefix(routerID uint, prefix string) ([]hotspot.Sale, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetSalesByPrefix(prefix)
}

func (w *HotspotClientWrapper) GetTotalRevenue(routerID uint, startDate, endDate string) (float64, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return 0, err
	}
	return client.GetTotalRevenue(startDate, endDate)
}

func (w *HotspotClientWrapper) DeleteSale(routerID uint, scriptID string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.DeleteSale(scriptID)
}

// Scheduler Methods
func (w *HotspotClientWrapper) CreateExpiryScheduler(routerID uint, profileName string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.CreateExpiryScheduler(profileName)
}

func (w *HotspotClientWrapper) RemoveExpiryScheduler(routerID uint, profileName string) error {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return err
	}
	return client.RemoveExpiryScheduler(profileName)
}

func (w *HotspotClientWrapper) GetAllSchedulers(routerID uint) ([]hotspot.Scheduler, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetAllSchedulers()
}

func (w *HotspotClientWrapper) GetSchedulerByName(routerID uint, name string) (*hotspot.Scheduler, error) {
	client, err := w.getHotspotClient(routerID)
	if err != nil {
		return nil, err
	}
	return client.GetSchedulerByName(name)
}
