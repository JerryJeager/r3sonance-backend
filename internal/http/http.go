package http

import "context"

type UserIdPP struct {
	UserID string `uri:"user_id" binding:"required,uuid_rfc4122"`
}
type PublicIDPP struct {
	PublicID string `uri:"public_id" binding:"required"`
}

func GetUserID(ctx context.Context) (string, error) {
	return ctx.Value("user_id").(string), nil
}
