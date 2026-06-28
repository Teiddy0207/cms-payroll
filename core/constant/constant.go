package constant

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusDeleted  = "deleted"
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

var ValidStatuses = map[string]bool{
	StatusActive:   true,
	StatusInactive: true,
	StatusDeleted:  true,
	StatusPending:  true,
	StatusApproved: true,
	StatusRejected: true,
}

func IsValidStatus(s string) bool {
	return ValidStatuses[s]
}