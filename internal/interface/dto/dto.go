package dto

type DashboardStats struct {
	TotalCustomers    int64   `json:"totalCustomers"`
	ActiveCustomers   int64   `json:"activeCustomers"`
	IsolatedCustomers int64   `json:"isolatedCustomers"`
	TotalPackages     int64   `json:"totalPackages"`
	TotalInvoices     int64   `json:"totalInvoices"`
	PaidInvoices      int64   `json:"paidInvoices"`
	PendingInvoices   int64   `json:"pendingInvoices"`
	TotalRevenue      float64 `json:"totalRevenue"`
}

type DashboardResponse struct {
	Stats           DashboardStats    `json:"stats"`
	RecentInvoices  []InvoiceSummary  `json:"recentInvoices"`
	RecentCustomers []CustomerSummary `json:"recentCustomers"`
}

type CustomerSummary struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	PackageID   uint   `json:"package_id"`
	PackageName string `json:"package_name"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type InvoiceSummary struct {
	ID           uint    `json:"id"`
	CustomerID   uint    `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Number       string  `json:"number"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

type CustomerListResponse struct {
	Customers  []CustomerDetail `json:"customers"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"perPage"`
	TotalPages int              `json:"totalPages"`
}

type CustomerDetail struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	Address        string  `json:"address"`
	PackageID      uint    `json:"package_id"`
	PackageName    string  `json:"package_name"`
	PackagePrice   float64 `json:"package_price"`
	PPPoEUsername  string  `json:"pppoe_username"`
	PPPoEPassword  string  `json:"pppoe_password,omitempty"`
	Status         string  `json:"status"`
	RouterID       uint    `json:"router_id"`
	ONUID          string  `json:"onu_id"`
	ONUSerial      string  `json:"onu_serial"`
	ONUMacAddress  string  `json:"onu_mac_address"`
	ONUIPAddress   string  `json:"onu_ip_address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	IsolationDate  *string `json:"isolation_date,omitempty"`
	ActivationDate *string `json:"activation_date,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type InvoiceListResponse struct {
	Invoices   []InvoiceDetail `json:"invoices"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
	TotalPages int             `json:"totalPages"`
}

type InvoiceDetail struct {
	ID               uint    `json:"id"`
	CustomerID       uint    `json:"customer_id"`
	CustomerName     string  `json:"customer_name"`
	CustomerPhone    string  `json:"customer_phone"`
	Number           string  `json:"number"`
	Amount           float64 `json:"amount"`
	Period           string  `json:"period"`
	DueDate          string  `json:"due_date"`
	Status           string  `json:"status"`
	PaidAt           *string `json:"paid_at,omitempty"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentReference string  `json:"payment_reference"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type PackageResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Speed       string  `json:"speed"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
