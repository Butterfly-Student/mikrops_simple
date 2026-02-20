package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

type HotspotSessionUsecase interface {
	GetActiveSessions(routerID uint) ([]*dto.HotspotSessionDetail, error)
	GetSessionByUsername(routerID uint, username string) (*dto.HotspotSessionDetail, error)
	DisconnectUser(routerID uint, username string) error
	GetSessionStats(routerID uint) (*dto.SessionStats, error)
}

type hotspotSessionUsecase struct {
	wrapper *mikrotik.HotspotClientWrapper
}

func NewHotspotSessionUsecase(wrapper *mikrotik.HotspotClientWrapper) HotspotSessionUsecase {
	return &hotspotSessionUsecase{
		wrapper: wrapper,
	}
}

func (u *hotspotSessionUsecase) GetActiveSessions(routerID uint) ([]*dto.HotspotSessionDetail, error) {
	sessions, err := u.wrapper.GetActiveSessions(routerID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.HotspotSessionDetail, len(sessions))
	for i, s := range sessions {
		dtos[i] = &dto.HotspotSessionDetail{
			Name:            s.Name,
			Address:         s.Address,
			MacAddress:      s.MacAddress,
			Uptime:          s.Uptime,
			SessionTimeLeft: s.SessionTimeLeft,
			BytesIn:         s.BytesIn,
			BytesOut:        s.BytesOut,
			LoginBy:         s.LoginBy,
		}
	}

	return dtos, nil
}

func (u *hotspotSessionUsecase) GetSessionByUsername(routerID uint, username string) (*dto.HotspotSessionDetail, error) {
	session, err := u.wrapper.GetSessionByUsername(routerID, username)
	if err != nil {
		return nil, err
	}

	return &dto.HotspotSessionDetail{
		Name:            session.Name,
		Address:         session.Address,
		MacAddress:      session.MacAddress,
		Uptime:          session.Uptime,
		SessionTimeLeft: session.SessionTimeLeft,
		BytesIn:         session.BytesIn,
		BytesOut:        session.BytesOut,
		LoginBy:         session.LoginBy,
	}, nil
}

func (u *hotspotSessionUsecase) DisconnectUser(routerID uint, username string) error {
	return u.wrapper.DisconnectUser(routerID, username)
}

func (u *hotspotSessionUsecase) GetSessionStats(routerID uint) (*dto.SessionStats, error) {
	stats, err := u.wrapper.GetSessionStats(routerID)
	if err != nil {
		return nil, err
	}

	return &dto.SessionStats{
		TotalUsers:    stats.TotalUsers,
		ActiveUsers:   stats.ActiveUsers,
		TotalBytesIn:  stats.TotalBytesIn,
		TotalBytesOut: stats.TotalBytesOut,
	}, nil
}
