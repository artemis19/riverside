package database

import (
	"github.com/artemis19/viz/rpc/sqltime"
)

// "Internal" hosts that have agents installed on them
type Host struct {
	OS           string             `json:"os"`
	Hostname     string             `json:"hostname"`
	Architecture string             `json:"arch"`
	MachineID    string             `json:"machine_id" gorm:"primaryKey;column:machine_id"`
	Interfaces   []NetworkInterface `json:"interfaces" gorm:"foreignKey:HostID;references:MachineID"`
	CreatedAt    sqltime.Time       `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt    sqltime.Time       `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt    *sqltime.Time      `json:"deleted_at" gorm:"type:timestamp"`
}

// Network interfaces for hosts w/agents
type NetworkInterface struct {
	ID        uint          `gorm:"primaryKey"`
	HostID    string        `json:"host_id"`
	Name      string        `json:"name"`
	IPAddress string        `json:"ip_address"`
	CreatedAt sqltime.Time  `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt sqltime.Time  `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt *sqltime.Time `json:"deleted_at" gorm:"type:timestamp"`
}

// "External" hosts or hosts that communicate with agent-installed hosts
type RemoteHost struct {
	RemoteHostID string        `json:"remote_host_id" gorm:"primaryKey;column:remote_host_id"`
	IPAddress    string        `json:"ipaddress"`
	DomainName   string        `json:"domainname"`
	CreatedAt    sqltime.Time  `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt    sqltime.Time  `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt    *sqltime.Time `json:"deleted_at" gorm:"type:timestamp"`
}

// Communication b/w hosts or "edges" on the visualization
type NetFlow struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	CreatedAt     sqltime.Time  `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt     sqltime.Time  `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt     *sqltime.Time `json:"deleted_at" gorm:"type:timestamp"`
	HostID        string        `json:"host_id"`        // maps to Host MachineID for primary key
	RemoteHostID  string        `json:"remote_host_id"` // maps to RemotHost ID for primary key
	OriginAddress string        `json:"origin_address"`
	SrcAddress    string        `json:"src_address"`
	DstAddress    string        `json:"dst_address"`
	SrcPort       int           `json:"src_port"`
	DstPort       int           `json:"dst_port"`
	Throughput    int           `json:"throughput"`
	Direction     string        `json:"direction"`
	Protocol      string        `json:"protocol"`
	StartTime     sqltime.Time  `json:"start_time"`
	EndTime       sqltime.Time  `json:"end_time"`
}

// Overrides default GORM table naming schema
type Tabler interface {
	TableName() string
}

func (NetFlow) TableName() string {
	return "net_flow"
}

//------------------------------------------------------------------

type User struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	CreatedAt    sqltime.Time  `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt    sqltime.Time  `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt    *sqltime.Time `json:"deleted_at" gorm:"type:timestamp"`
	Username     string        `json:"username"`
	PasswordHash string        `json:"password_hash"`
}
