package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/pkg/hotspot"
)

type HotspotProfileUsecase interface {
	CreateProfile(routerID uint, req *dto.CreateHotspotProfile) error
	UpdateProfile(routerID uint, profileName string, req *dto.UpdateHotspotProfile) error
	DeleteProfile(routerID uint, profileName string) error
	GetAllProfiles(routerID uint) ([]*dto.HotspotProfileDetail, error)
	GetProfile(routerID uint, profileName string) (*dto.HotspotProfileDetail, error)
}

type hotspotProfileUsecase struct {
	wrapper *mikrotik.HotspotClientWrapper
}

func NewHotspotProfileUsecase(wrapper *mikrotik.HotspotClientWrapper) HotspotProfileUsecase {
	return &hotspotProfileUsecase{
		wrapper: wrapper,
	}
}

func (u *hotspotProfileUsecase) CreateProfile(routerID uint, req *dto.CreateHotspotProfile) error {
	price := req.Price
	sellingPrice := req.SellingPrice

	profile := &hotspot.Profile{
		Name:             req.Name,
		SharedUsers:      1,
		RateLimit:        req.RateLimit,
		Validity:         req.Validity,
		Price:            price,
		SellingPrice:     sellingPrice,
		ExpiryMode:       req.ExpiryMode,
		LockUser:         req.LockUser,
		KeepaliveTimeout: "5m",
	}

	return u.wrapper.CreateProfile(routerID, profile)
}

func (u *hotspotProfileUsecase) UpdateProfile(routerID uint, profileName string, req *dto.UpdateHotspotProfile) error {
	updates := &hotspot.Profile{}

	if req.RateLimit != "" {
		updates.RateLimit = req.RateLimit
	}
	if req.SharedUsers != nil && *req.SharedUsers > 0 {
		updates.SharedUsers = *req.SharedUsers
	}
	if req.Validity != "" {
		updates.Validity = req.Validity
	}
	if req.Price != nil {
		updates.Price = *req.Price
	}
	if req.SellingPrice != nil {
		updates.SellingPrice = *req.SellingPrice
	}
	if req.ExpiryMode != "" {
		updates.ExpiryMode = req.ExpiryMode
	}
	if req.LockUser != "" {
		updates.LockUser = req.LockUser
	}

	return u.wrapper.UpdateProfile(routerID, profileName, updates)
}

func (u *hotspotProfileUsecase) DeleteProfile(routerID uint, profileName string) error {
	return u.wrapper.DeleteProfile(routerID, profileName)
}

func (u *hotspotProfileUsecase) GetAllProfiles(routerID uint) ([]*dto.HotspotProfileDetail, error) {
	profiles, err := u.wrapper.GetAllProfiles(routerID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.HotspotProfileDetail, len(profiles))
	for i, p := range profiles {
		dtos[i] = &dto.HotspotProfileDetail{
			Name:         p.Name,
			SharedUsers:  p.SharedUsers,
			RateLimit:    p.RateLimit,
			Validity:     p.Validity,
			Price:        p.Price,
			SellingPrice: p.SellingPrice,
			ExpiryMode:   p.ExpiryMode,
			LockUser:     p.LockUser,
		}
	}

	return dtos, nil
}

func (u *hotspotProfileUsecase) GetProfile(routerID uint, profileName string) (*dto.HotspotProfileDetail, error) {
	profile, err := u.wrapper.GetProfile(routerID, profileName)
	if err != nil {
		return nil, err
	}

	return &dto.HotspotProfileDetail{
		Name:         profile.Name,
		SharedUsers:  profile.SharedUsers,
		RateLimit:    profile.RateLimit,
		Validity:     profile.Validity,
		Price:        profile.Price,
		SellingPrice: profile.SellingPrice,
		ExpiryMode:   profile.ExpiryMode,
		LockUser:     profile.LockUser,
	}, nil
}
