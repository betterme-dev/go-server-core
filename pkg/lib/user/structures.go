package user

type (
	User struct {
		ID             int   `db:"id"`
		AuthKeyExpires int64 `db:"auth_key_expires"`
	}
	Session struct {
		ID        int   `db:"id"`
		UserID    int   `db:"user_id"`
		ExpiresAt int64 `db:"expires_at"`
	}
)
