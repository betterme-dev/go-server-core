package user

type (
	User struct {
		ID             int `db:"id"`
		AuthKeyExpires int `db:"auth_key_expires"`
	}
	Session struct {
		ID        int `db:"id"`
		UserID    int `db:"user_id"`
		ExpiresAt int `db:"expires_at"`
	}
)
