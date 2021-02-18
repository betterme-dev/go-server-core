package user

type User struct {
	ID             int   `db:"id"`
	AuthKeyExpires int64 `db:"auth_key_expires"`
}
