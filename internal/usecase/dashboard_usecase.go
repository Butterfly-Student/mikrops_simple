package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

type DashboardUsecase interface {
	GetDashboardStats() (*dto.DashboardResponse, error)
}

type dashboardUsecase struct {
	customerRepo repositories.CustomerRepository
	invoiceRepo  repositories.InvoiceRepository
	packageRepo  repositories.PackageRepository
}

func NewDashboardUsecase(
	customerRepo repositories.CustomerRepository,
	invoiceRepo repositories.InvoiceRepository,
	packageRepo repositories.PackageRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		customerRepo: customerRepo,
		invoiceRepo:  invoiceRepo,
		packageRepo:  packageRepo,
	}
}

func (u *dashboardUsecase) GetDashboardStats() (*dto.DashboardResponse, error) {
	stats := dto.DashboardStats{}

	totalCustomers, _, err := u.customerRepo.FindAll(1, 1, "")
	if err != nil {
		return nil, err
	}
	stats.TotalCustomers = int64(len(totalCustomers))

	activeCustomers, _, err := u.customerRepo.FindByStatus("active", 1, 1)
	if err != nil {
		return nil, err
	}
	stats.ActiveCustomers = int64(len(activeCustomers))

	isolatedCustomers, _, err := u.customerRepo.FindByStatus("isolated", 1, 1)
	if err != nil {
		return nil, err
	}
	stats.IsolatedCustomers = int64(len(isolatedCustomers))

	packages, err := u.packageRepo.FindAll()
	if err != nil {
		return nil, err
	}
	stats.TotalPackages = int64(len(packages))

	totalInvoices, _, err := u.invoiceRepo.FindAll(1, 1)
	if err != nil {
		return nil, err
	}
	stats.TotalInvoices = int64(len(totalInvoices))

	paidInvoices, _, err := u.invoiceRepo.FindByStatus("paid", 1, 100)
	if err != nil {
		return nil, err
	}
	stats.PaidInvoices = int64(len(paidInvoices))

	pendingInvoices, _, err := u.invoiceRepo.FindByStatus("unpaid", 1, 100)
	if err != nil {
		return nil, err
	}
	stats.PendingInvoices = int64(len(pendingInvoices))

	totalRevenue := 0.0
	for _, invoice := range paidInvoices {
		totalRevenue += invoice.Amount
	}
	stats.TotalRevenue = totalRevenue

	recentInvoices, _, err := u.invoiceRepo.FindAll(1, 10)
	if err != nil {
		return nil, err
	}

	recentInvoicesDTO := make([]dto.InvoiceSummary, 0, len(recentInvoices))
	for _, invoice := range recentInvoices {
		customerName := ""
		if invoice.Customer != nil {
			customerName = invoice.Customer.Name
		}
		recentInvoicesDTO = append(recentInvoicesDTO, dto.InvoiceSummary{
			ID:           invoice.ID,
			CustomerID:   invoice.CustomerID,
			CustomerName: customerName,
			Number:       invoice.Number,
			Amount:       invoice.Amount,
			Status:       invoice.Status,
			CreatedAt:    invoice.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	recentCustomers, _, err := u.customerRepo.FindAll(1, 5, "")
	if err != nil {
		return nil, err
	}

	recentCustomersDTO := make([]dto.CustomerSummary, 0, len(recentCustomers))
	for _, customer := range recentCustomers {
		packageName := ""
		if customer.Package != nil {
			packageName = customer.Package.Name
		}
		recentCustomersDTO = append(recentCustomersDTO, dto.CustomerSummary{
			ID:          customer.ID,
			Name:        customer.Name,
			Phone:       customer.Phone,
			PackageID:   customer.PackageID,
			PackageName: packageName,
			Status:      customer.Status,
			CreatedAt:   customer.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.DashboardResponse{
		Stats:           stats,
		RecentInvoices:  recentInvoicesDTO,
		RecentCustomers: recentCustomersDTO,
	}, nil
}
