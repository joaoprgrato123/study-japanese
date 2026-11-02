package role

import "time"

type Role struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}
