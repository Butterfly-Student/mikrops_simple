package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/pkg/hotspot"
)

type HotspotSaleUsecase interface {
	RecordSale(routerID uint, req *dto.RecordHotspotSale) error
	GetAllSales(routerID uint, filter *dto.HotspotSaleFilter) ([]*dto.HotspotSaleDetail, int, error)
	GetTotalRevenue(routerID uint, startDate, endDate string) (float64, error)
}

type hotspotSaleUsecase struct {
	wrapper *mikrotik.HotspotClientWrapper
}

func NewHotspotSaleUsecase(wrapper *mikrotik.HotspotClientWrapper) HotspotSaleUsecase {
	return &hotspotSaleUsecase{
		wrapper: wrapper,
	}
}

func (u *hotspotSaleUsecase) RecordSale(routerID uint, req *dto.RecordHotspotSale) error {
	sale := &hotspot.Sale{
		Username: req.Username,
		Price:    req.Price,
		Address:  req.Address,
		Mac:      req.Mac,
		Validity: req.Validity,
	}

	return u.wrapper.RecordSale(routerID, sale)
}

func (u *hotspotSaleUsecase) GetAllSales(routerID uint, filter *dto.HotspotSaleFilter) ([]*dto.HotspotSaleDetail, int, error) {
	sales, err := u.wrapper.GetAllSales(routerID, &hotspot.SaleFilter{
		StartDate: filter.StartDate,
		EndDate:   filter.EndDate,
		Prefix:    filter.Prefix,
		Limit:     filter.Limit,
		Offset:    filter.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*dto.HotspotSaleDetail, len(sales))
	for i, s := range sales {
		dtos[i] = &dto.HotspotSaleDetail{
			Date:     s.Date,
			Time:     s.Time,
			Username: s.Username,
			Price:    s.Price,
			Address:  s.Address,
			Mac:      s.Mac,
			Validity: s.Validity,
		}
	}

	total := len(sales)
	return dtos, total, nil
}

func (u *hotspotSaleUsecase) GetTotalRevenue(routerID uint, startDate, endDate string) (float64, error) {
	return u.wrapper.GetTotalRevenue(routerID, startDate, endDate)
}
