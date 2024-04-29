package domain

type User struct {
	ID        int64    `json:"id"`
	TguID     int64    `json:"tgu_id"`
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Subs      []string `json:"subs"`
	Created   int64    `json:"created"`
	Edited    int64    `json:"edited"`
	Deleted   int64    `json:"deleted"`
}
