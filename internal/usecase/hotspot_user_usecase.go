package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/pkg/hotspot"
)

type HotspotUserUsecase interface {
	CreateUser(routerID uint, req *dto.CreateHotspotUser) error
	UpdateUser(routerID uint, username string, req *dto.UpdateHotspotUser) error
	DeleteUser(routerID uint, username string) error
	GetUser(routerID uint, username string) (*dto.HotspotUserDetail, error)
	GetAllUsers(routerID uint, filter *dto.HotspotUserFilter) ([]*dto.HotspotUserDetail, int, error)
	DisableUser(routerID uint, username string) error
	EnableUser(routerID uint, username string) error
	RemoveExpiredUsers(routerID uint, profile string) (int, error)
	RemoveUnusedVouchers(routerID uint, profile string) (int, error)
}

type hotspotUserUsecase struct {
	wrapper *mikrotik.HotspotClientWrapper
}

func NewHotspotUserUsecase(wrapper *mikrotik.HotspotClientWrapper) HotspotUserUsecase {
	return &hotspotUserUsecase{
		wrapper: wrapper,
	}
}

func (u *hotspotUserUsecase) CreateUser(routerID uint, req *dto.CreateHotspotUser) error {
	user := &hotspot.User{
		Name:            req.Name,
		Password:        req.Password,
		Profile:         req.Profile,
		Comment:         req.Comment,
		LimitUptime:     req.LimitUptime,
		LimitBytesTotal: req.LimitBytesTotal,
		LimitBytesIn:    req.LimitBytesIn,
		LimitBytesOut:   req.LimitBytesOut,
		Disabled:        req.Disabled,
		Server:          req.Server,
	}

	return u.wrapper.CreateUser(routerID, user)
}

func (u *hotspotUserUsecase) UpdateUser(routerID uint, username string, req *dto.UpdateHotspotUser) error {
	updates := make(map[string]interface{})

	if req.Profile != nil {
		updates["profile"] = *req.Profile
	}
	if req.Disabled != nil {
		updates["disabled"] = *req.Disabled
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}
	if req.LimitUptime != nil {
		updates["limit_uptime"] = *req.LimitUptime
	}
	if req.LimitBytesTotal != nil {
		updates["limit_bytes_total"] = *req.LimitBytesTotal
	}

	return u.wrapper.UpdateUser(routerID, username, updates)
}

func (u *hotspotUserUsecase) DeleteUser(routerID uint, username string) error {
	return u.wrapper.DeleteUser(routerID, username)
}

func (u *hotspotUserUsecase) GetUser(routerID uint, username string) (*dto.HotspotUserDetail, error) {
	user, err := u.wrapper.GetUser(routerID, username)
	if err != nil {
		return nil, err
	}

	return &dto.HotspotUserDetail{
		Name:            user.Name,
		Password:        user.Password,
		Profile:         user.Profile,
		Comment:         user.Comment,
		LimitUptime:     user.LimitUptime,
		LimitBytesTotal: user.LimitBytesTotal,
		Disabled:        user.Disabled,
		Uptime:          user.Uptime,
		BytesIn:         user.BytesIn,
		BytesOut:        user.BytesOut,
	}, nil
}

func (u *hotspotUserUsecase) GetAllUsers(routerID uint, filter *dto.HotspotUserFilter) ([]*dto.HotspotUserDetail, int, error) {
	users, err := u.wrapper.GetAllUsers(routerID, &hotspot.UserFilter{
		Profile:  filter.Profile,
		Comment:  filter.Comment,
		Disabled: filter.Disabled,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*dto.HotspotUserDetail, len(users))
	for i, u := range users {
		dtos[i] = &dto.HotspotUserDetail{
			Name:            u.Name,
			Password:        u.Password,
			Profile:         u.Profile,
			Comment:         u.Comment,
			LimitUptime:     u.LimitUptime,
			LimitBytesTotal: u.LimitBytesTotal,
			Disabled:        u.Disabled,
			Uptime:          u.Uptime,
			BytesIn:         u.BytesIn,
			BytesOut:        u.BytesOut,
		}
	}

	total := len(users)
	return dtos, total, nil
}

func (u *hotspotUserUsecase) DisableUser(routerID uint, username string) error {
	return u.wrapper.DisableUser(routerID, username)
}

func (u *hotspotUserUsecase) EnableUser(routerID uint, username string) error {
	return u.wrapper.EnableUser(routerID, username)
}

func (u *hotspotUserUsecase) RemoveExpiredUsers(routerID uint, profile string) (int, error) {
	return u.wrapper.RemoveExpiredUsers(routerID, profile)
}

func (u *hotspotUserUsecase) RemoveUnusedVouchers(routerID uint, profile string) (int, error) {
	return u.wrapper.RemoveUnusedVouchers(routerID, profile)
}
