package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID           int       `json:"user_id"`
	UUID             uuid.UUID `json:"uuid"`
	Username         string    `json:"username"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	PhoneNumber      string    `json:"phone_number"`
	RoleId           int64     `json:"role_id"`
	MasterAgentId    *int64    `json:"master_agent_id,omitempty"`
	AgentId          *int64    `json:"agent_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatedBy        uuid.UUID `json:"created_by"`
	UpdatedBy        uuid.UUID `json:"updated_by"`
	LastBluetoothMac string    `json:"last_bluetooth_mac,omitempty"`
	RoleName         string    `json:"role_name,omitempty"`
	ApiPassKey       string    `json:"api_pass_key,omitempty"`
}

type Role struct {
	RoleID      int       `json:"role_id"`
	RoleName    string    `json:"role_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
