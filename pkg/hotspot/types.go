package hotspot

// Profile represents MikroTik hotspot user profile
type Profile struct {
	Name             string  `json:"name"`
	SharedUsers      int     `json:"shared_users"`
	RateLimit        string  `json:"rate_limit"`
	Validity         string  `json:"validity"`
	Price            float64 `json:"price"`
	SellingPrice     float64 `json:"selling_price"`
	ExpiryMode       string  `json:"expiry_mode"` // rem, ntf, remc, ntfc
	LockUser         string  `json:"lock_user"`
	KeepaliveTimeout string  `json:"keepalive_timeout"`
	OnLoginScript    string  `json:"on_login_script"`
}

// User represents MikroTik hotspot user
type User struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	Profile         string `json:"profile"`
	Comment         string `json:"comment"`
	LimitUptime     int64  `json:"limit_uptime"`      // seconds, 0 = unlimited
	LimitBytesTotal int64  `json:"limit_bytes_total"` // bytes, 0 = unlimited
	LimitBytesIn    int64  `json:"limit_bytes_in"`
	LimitBytesOut   int64  `json:"limit_bytes_out"`
	Disabled        bool   `json:"disabled"`
	Server          string `json:"server"`
	Uptime          string `json:"uptime"`
	BytesIn         string `json:"bytes_in"`
	BytesOut        string `json:"bytes_out"`
}

// Session represents active hotspot session
type Session struct {
	Name            string `json:"name"`
	Address         string `json:"address"`
	MacAddress      string `json:"mac_address"`
	Uptime          string `json:"uptime"`
	SessionTimeLeft string `json:"session_time_left"`
	BytesIn         string `json:"bytes_in"`
	BytesOut        string `json:"bytes_out"`
	LoginBy         string `json:"login_by"`
}

// Sale represents sales record stored in RouterOS scripts
type Sale struct {
	Date     string  `json:"date"`
	Time     string  `json:"time"`
	Username string  `json:"username"`
	Price    float64 `json:"price"`
	Address  string  `json:"address"`
	Mac      string  `json:"mac"`
	Validity string  `json:"validity"`
	ScriptID string  `json:"script_id"`
}

// Scheduler represents RouterOS scheduler for expiry monitoring
type Scheduler struct {
	Name      string `json:"name"`
	Interval  string `json:"interval"`
	StartTime string `json:"start_time"`
	Policy    string `json:"policy"`
	OnEvent   string `json:"on_event"`
	Enabled   bool   `json:"enabled"`
}

// VoucherGenerator configuration for batch voucher generation
type VoucherGenerator struct {
	Profile        string
	Prefix         string
	Charset        string
	LengthUsername int
	LengthPassword int
	Quantity       int
	TimeLimit      int64 // seconds
	DataLimit      int64 // bytes
}

// VoucherResult contains generated voucher information
type VoucherResult struct {
	Success  int      `json:"success"`
	Failed   int      `json:"failed"`
	Vouchers []User   `json:"vouchers"`
	Errors   []string `json:"errors,omitempty"`
}

// SessionStats contains session statistics
type SessionStats struct {
	TotalUsers    int    `json:"total_users"`
	ActiveUsers   int    `json:"active_users"`
	TotalBytesIn  string `json:"total_bytes_in"`
	TotalBytesOut string `json:"total_bytes_out"`
}

// Filter options for user queries
type UserFilter struct {
	Profile  string
	Comment  string
	Disabled *bool
	Limit    int
	Offset   int
}

// Filter options for sales queries
type SaleFilter struct {
	StartDate string // format: "jan/01/2024"
	EndDate   string // format: "jan/31/2024"
	Prefix    string
	Limit     int
	Offset    int
}
