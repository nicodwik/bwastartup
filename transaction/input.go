package transaction

import "bwastartup/user"

type GetTransactionsDetailInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
