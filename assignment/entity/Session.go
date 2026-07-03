package entity

type Session struct {
	UserID    int    `json:"userID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	TTL       uint64 `json:"ttl"`
}
