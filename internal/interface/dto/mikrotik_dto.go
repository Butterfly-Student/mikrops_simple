package dto

// Router DTOs
type RouterCreate struct {
	Name     string `json:"name" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Port     int    `json:"port"`
	IsActive bool   `json:"is_active"`
}

type RouterUpdate struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     *int   `json:"port"`
	IsActive *bool  `json:"is_active"`
}

type RouterStatus struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Host        string  `json:"host"`
	Status      string  `json:"status"`
	LastCheck   string  `json:"last_check"`
	ActiveUsers int     `json:"active_users"`
	CPU         float64 `json:"cpu"`
	Memory      int     `json:"memory"`
	Uptime      string  `json:"uptime"`
	Error       string  `json:"error,omitempty"`
}

type RouterDetail struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	Username  string `json:"username"`
	Port      int    `json:"port"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ConnectionTestResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Latency      int    `json:"latency_ms"`
	RouterOS     string `json:"routeros_version"`
	ResourceName string `json:"resource_name"`
}

// PPPoE DTOs
type PPPUser struct {
	Name      string `json:"name"`
	Service   string `json:"service"`
	Profile   string `json:"profile"`
	CallerID  string `json:"caller_id"`
	Disabled  bool   `json:"disabled"`
	LastLogin string `json:"last_login,omitempty"`
}

type PPPUsersResponse struct {
	Users  []PPPUser     `json:"users"`
	Total  int           `json:"total"`
	Router *RouterDetail `json:"router,omitempty"`
}

type ActiveSession struct {
	Name     string `json:"name"`
	CallerID string `json:"caller_id"`
	Address  string `json:"address"`
	Uptime   string `json:"uptime"`
	BytesIn  string `json:"bytes_in"`
	BytesOut string `json:"bytes_out"`
	Encoding string `json:"encoding"`
}

type ActiveSessionsResponse struct {
	Sessions []ActiveSession `json:"sessions"`
	Total    int             `json:"total"`
	Router   *RouterDetail   `json:"router,omitempty"`
}

type Profile struct {
	Name         string `json:"name"`
	RateLimit    string `json:"rate_limit"`
	LocalAddress string `json:"local_address"`
	OnlyOne      bool   `json:"only_one"`
}

type ProfilesResponse struct {
	Profiles []Profile     `json:"profiles"`
	Total    int           `json:"total"`
	Router   *RouterDetail `json:"router,omitempty"`
}

type AddPPPUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Profile  string `json:"profile" binding:"required"`
	Service  string `json:"service"`
	RouterID uint   `json:"router_id"`
}

type UpdatePPPUserRequest struct {
	Profile  string `json:"profile"`
	Disabled *bool  `json:"disabled"`
}

type DisconnectUserRequest struct {
	Username string `json:"username" binding:"required"`
	RouterID uint   `json:"router_id"`
}

// Customer MikroTik Operations DTOs
type BulkOperationRequest struct {
	CustomerIDs []uint `json:"customer_ids" binding:"required"`
}

type BulkOperationResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailedIDs    []uint   `json:"failed_ids"`
	Errors       []string `json:"errors,omitempty"`
}
