package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/pkg/hotspot"
)

type VoucherUsecase interface {
	GenerateVouchers(routerID uint, req *dto.GenerateVouchers) (*dto.VoucherResult, error)
	GetVouchersByPrefix(routerID uint, prefix string) ([]*dto.HotspotUserDetail, error)
}

type voucherUsecase struct {
	wrapper *mikrotik.HotspotClientWrapper
}

func NewVoucherUsecase(wrapper *mikrotik.HotspotClientWrapper) VoucherUsecase {
	return &voucherUsecase{
		wrapper: wrapper,
	}
}

func (u *voucherUsecase) GenerateVouchers(routerID uint, req *dto.GenerateVouchers) (*dto.VoucherResult, error) {
	gen := &hotspot.VoucherGenerator{
		Profile:        req.Profile,
		Prefix:         req.Prefix,
		Quantity:       req.Quantity,
		LengthUsername: req.LengthUsername,
		LengthPassword: req.LengthPassword,
		Charset:        req.Charset,
		TimeLimit:      req.TimeLimit,
		DataLimit:      req.DataLimit,
	}

	result, err := u.wrapper.GenerateVouchers(routerID, gen)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.VoucherDetail, len(result.Vouchers))
	for i, v := range result.Vouchers {
		dtos[i] = dto.VoucherDetail{
			Username: v.Name,
			Password: v.Password,
			Profile:  v.Profile,
		}
	}

	return &dto.VoucherResult{
		Success:  result.Success,
		Failed:   result.Failed,
		Vouchers: dtos,
		Errors:   result.Errors,
	}, nil
}

func (u *voucherUsecase) GetVouchersByPrefix(routerID uint, prefix string) ([]*dto.HotspotUserDetail, error) {
	users, err := u.wrapper.GetAllUsers(routerID, nil)
	if err != nil {
		return nil, err
	}

	var filteredUsers []hotspot.User
	for _, user := range users {
		if prefix == "" || len(user.Name) >= len(prefix) && user.Name[:len(prefix)] == prefix {
			filteredUsers = append(filteredUsers, user)
		}
	}

	dtos := make([]*dto.HotspotUserDetail, len(filteredUsers))
	for i, u := range filteredUsers {
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

	return dtos, nil
}
