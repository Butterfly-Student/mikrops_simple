package mikrotik

import "time"

type RouterStatus struct {
	RouterID    uint      `json:"router_id"`
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	Status      string    `json:"status"`
	LastCheck   time.Time `json:"last_check"`
	ActiveUsers int       `json:"active_users"`
	CPU         float64   `json:"cpu"`
	Memory      int       `json:"memory"`
	Uptime      string    `json:"uptime"`
	Error       string    `json:"error,omitempty"`
}

type PPPoEUser struct {
	Name      string `json:"name"`
	Service   string `json:"service"`
	Profile   string `json:"profile"`
	CallerID  string `json:"caller_id"`
	Disabled  bool   `json:"disabled"`
	LastLogin string `json:"last_login,omitempty"`
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

type Profile struct {
	Name         string `json:"name"`
	RateLimit    string `json:"rate_limit"`
	LocalAddress string `json:"local_address"`
	OnlyOne      bool   `json:"only_one"`
}

type HotspotUser struct {
	Name     string `json:"name"`
	Profile  string `json:"profile"`
	BytesIn  string `json:"bytes_in"`
	BytesOut string `json:"bytes_out"`
	Uptime   string `json:"uptime"`
}

type HotspotSession struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	MacAddress string `json:"mac_address"`
	Uptime     string `json:"uptime"`
	BytesIn    string `json:"bytes_in"`
	BytesOut   string `json:"bytes_out"`
}

type HotspotLog struct {
	Time    time.Time `json:"time"`
	Topic   string    `json:"topic"`
	Message string    `json:"message"`
}

type Voucher struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Profile    string `json:"profile"`
	ValidUntil string `json:"valid_until"`
}

type SimpleQueue struct {
	Number     int    `json:"number"`
	Username   string `json:"username"`
	Expiration string `json:"expiration"`
}
