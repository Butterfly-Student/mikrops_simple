package entities

import "time"

type AdminUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Email     string    `json:"email"`
	Role      string    `gorm:"default:'admin'" json:"role"`
	Status    string    `gorm:"default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Customer struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Name           string     `gorm:"not null" json:"name"`
	Phone          string     `gorm:"uniqueIndex;not null" json:"phone"`
	Email          string     `json:"email"`
	Address        string     `gorm:"type:text" json:"address"`
	PackageID      uint       `json:"package_id"`
	Package        *Package   `gorm:"foreignKey:PackageID" json:"package,omitempty"`
	PPPoEUsername  string     `gorm:"uniqueIndex" json:"pppoe_username"`
	PPPoEPassword  string     `gorm:"not null" json:"-"`
	Status         string     `gorm:"default:'active'" json:"status"`
	RouterID       uint       `json:"router_id"`
	ONUID          string     `json:"onu_id"`
	ONUSerial      string     `json:"onu_serial"`
	ONUMacAddress  string     `json:"onu_mac_address"`
	ONUIPAddress   string     `json:"onu_ip_address"`
	Latitude       float64    `json:"latitude"`
	Longitude      float64    `json:"longitude"`
	IsolationDate  *time.Time `json:"isolation_date,omitempty"`
	ActivationDate *time.Time `json:"activation_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Package struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"uniqueIndex;not null" json:"name"`
	Price         float64   `gorm:"not null" json:"price"`
	Speed         string    `json:"speed"`
	Description   string    `gorm:"type:text" json:"description"`
	ProfileNormal string    `gorm:"column:profile_normal" json:"profile_normal"`
	ProfileIsolir string    `gorm:"column:profile_isolir" json:"profile_isolir"`
	Status        string    `gorm:"default:'active'" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Invoice struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	CustomerID       uint       `gorm:"not null" json:"customer_id"`
	Customer         *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Number           string     `gorm:"uniqueIndex;not null" json:"number"`
	Amount           float64    `gorm:"not null" json:"amount"`
	Period           string     `gorm:"not null" json:"period"`
	DueDate          time.Time  `json:"due_date"`
	Status           string     `gorm:"default:'unpaid'" json:"status"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	PaymentMethod    string     `json:"payment_method"`
	PaymentReference string     `json:"payment_reference"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type Router struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Host      string    `gorm:"not null" json:"host"`
	Username  string    `gorm:"not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Port      int       `gorm:"default:8728" json:"port"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ONULocation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CustomerID     uint      `gorm:"not null;index" json:"customer_id"`
	Customer       *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	RouterID       uint      `json:"router_id"`
	ONUID          string    `gorm:"uniqueIndex" json:"onu_id"`
	SerialNumber   string    `json:"serial_number"`
	MacAddress     string    `json:"mac_address"`
	IPAddress      string    `json:"ip_address"`
	PortNumber     string    `json:"port_number"`
	SignalStrength string    `json:"signal_strength"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	Address        string    `gorm:"type:text" json:"address"`
	Notes          string    `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TroubleTicket struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CustomerID  uint       `gorm:"not null;index" json:"customer_id"`
	Customer    *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Subject     string     `gorm:"not null" json:"subject"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Priority    string     `gorm:"default:'medium'" json:"priority"`
	Status      string     `gorm:"default:'open'" json:"status"`
	AssignedTo  string     `json:"assigned_to"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Setting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SettingKey   string    `gorm:"uniqueIndex;not null" json:"setting_key"`
	SettingValue string    `gorm:"type:text" json:"setting_value"`
	Description  string    `gorm:"type:text" json:"description"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CronSchedule struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"not null" json:"name"`
	TaskType     string     `gorm:"not null" json:"task_type"`
	ScheduleTime string     `json:"schedule_time"`
	ScheduleDays string     `json:"schedule_days"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastRunAt    *time.Time `json:"last_run_at,omitempty"`
	NextRunAt    *time.Time `json:"next_run_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CronLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScheduleID uint      `gorm:"index" json:"schedule_id"`
	TaskType   string    `json:"task_type"`
	Status     string    `json:"status"`
	Output     string    `gorm:"type:text" json:"output"`
	Error      string    `gorm:"type:text" json:"error"`
	CreatedAt  time.Time `json:"created_at"`
}

type WebhookLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Event      string    `gorm:"not null;index" json:"event"`
	URL        string    `gorm:"not null" json:"url"`
	Payload    string    `gorm:"type:longtext" json:"payload"`
	Response   string    `gorm:"type:longtext" json:"response"`
	StatusCode int       `json:"status_code"`
	Duration   int       `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}
