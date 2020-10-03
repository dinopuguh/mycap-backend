package group

// CreateGroup is a data transfer object for create group
type CreateGroup struct {
	Type string `json:"type"`
}

// JoinGroup is a data transfer object for joining group
type JoinGroup struct {
	AdminUsername string `json:"admin_username"`
}

// LeaveGroup is a data transfer object for leaving group
type LeaveGroup struct {
	AdminUsername string `json:"admin_username"`
	RemainingTime int64  `json:"remaining_time"`
}
