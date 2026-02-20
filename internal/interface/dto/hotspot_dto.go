package dto

// Profile DTOs
type CreateHotspotProfile struct {
	Name         string  `json:"name" binding:"required"`
	SharedUsers  int     `json:"shared_users"`
	RateLimit    string  `json:"rate_limit"`
	Validity     string  `json:"validity"`
	Price        float64 `json:"price"`
	SellingPrice float64 `json:"selling_price"`
	ExpiryMode   string  `json:"expiry_mode"`
	LockUser     string  `json:"lock_user"`
}

type UpdateHotspotProfile struct {
	SharedUsers  *int     `json:"shared_users"`
	RateLimit    string   `json:"rate_limit"`
	Validity     string   `json:"validity"`
	Price        *float64 `json:"price"`
	SellingPrice *float64 `json:"selling_price"`
	ExpiryMode   string   `json:"expiry_mode"`
	LockUser     string   `json:"lock_user"`
}

type HotspotProfileDetail struct {
	Name         string  `json:"name"`
	SharedUsers  int     `json:"shared_users"`
	RateLimit    string  `json:"rate_limit"`
	Validity     string  `json:"validity"`
	Price        float64 `json:"price"`
	SellingPrice float64 `json:"selling_price"`
	ExpiryMode   string  `json:"expiry_mode"`
	LockUser     string  `json:"lock_user"`
}

// User DTOs
type CreateHotspotUser struct {
	Name            string `json:"name" binding:"required"`
	Password        string `json:"password" binding:"required"`
	Profile         string `json:"profile" binding:"required"`
	Comment         string `json:"comment"`
	LimitUptime     int64  `json:"limit_uptime"`
	LimitBytesTotal int64  `json:"limit_bytes_total"`
	LimitBytesIn    int64  `json:"limit_bytes_in"`
	LimitBytesOut   int64  `json:"limit_bytes_out"`
	Disabled        bool   `json:"disabled"`
	Server          string `json:"server"`
}

type UpdateHotspotUser struct {
	Profile         *string `json:"profile"`
	Disabled        *bool   `json:"disabled"`
	Comment         *string `json:"comment"`
	LimitUptime     *int64  `json:"limit_uptime"`
	LimitBytesTotal *int64  `json:"limit_bytes_total"`
}

type HotspotUserDetail struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	Profile         string `json:"profile"`
	Comment         string `json:"comment"`
	LimitUptime     int64  `json:"limit_uptime"`
	LimitBytesTotal int64  `json:"limit_bytes_total"`
	Disabled        bool   `json:"disabled"`
	Uptime          string `json:"uptime"`
	BytesIn         string `json:"bytes_in"`
	BytesOut        string `json:"bytes_out"`
}

type HotspotUserFilter struct {
	Profile  string `form:"profile"`
	Comment  string `form:"comment"`
	Disabled *bool  `form:"disabled"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
}

// Voucher DTOs
type GenerateVouchers struct {
	Profile        string `json:"profile" binding:"required"`
	Prefix         string `json:"prefix"`
	Quantity       int    `json:"quantity" binding:"required"`
	LengthUsername int    `json:"length_username"`
	LengthPassword int    `json:"length_password"`
	Charset        string `json:"charset"`
	TimeLimit      int64  `json:"time_limit"`
	DataLimit      int64  `json:"data_limit"`
}

type VoucherResult struct {
	Success  int             `json:"success"`
	Failed   int             `json:"failed"`
	Vouchers []VoucherDetail `json:"vouchers"`
	Errors   []string        `json:"errors,omitempty"`
}

type VoucherDetail struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
	Validity string `json:"validity"`
	Price    string `json:"price"`
	QrCode   string `json:"qr_code,omitempty"`
}

// Sales DTOs
type RecordHotspotSale struct {
	Username string  `json:"username" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Address  string  `json:"address"`
	Mac      string  `json:"mac"`
	Validity string  `json:"validity"`
}

type HotspotSaleDetail struct {
	Date     string  `json:"date"`
	Time     string  `json:"time"`
	Username string  `json:"username"`
	Price    float64 `json:"price"`
	Address  string  `json:"address"`
	Mac      string  `json:"mac"`
	Validity string  `json:"validity"`
}

type HotspotSaleFilter struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Prefix    string `form:"prefix"`
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
}

// Session DTOs
type HotspotSessionDetail struct {
	Name            string `json:"name"`
	Address         string `json:"address"`
	MacAddress      string `json:"mac_address"`
	Uptime          string `json:"uptime"`
	SessionTimeLeft string `json:"session_time_left"`
	BytesIn         string `json:"bytes_in"`
	BytesOut        string `json:"bytes_out"`
	LoginBy         string `json:"login_by"`
}

type SessionStats struct {
	TotalUsers    int    `json:"total_users"`
	ActiveUsers   int    `json:"active_users"`
	TotalBytesIn  string `json:"total_bytes_in"`
	TotalBytesOut string `json:"total_bytes_out"`
}
