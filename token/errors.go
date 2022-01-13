package token

import (
	"fmt"
)

type Error struct {
	ErrorDescription string `json:"error_description"`
	ErrorCode        string `json:"error"`
}

func (r *Error) Error() string {
	return fmt.Sprintf("[token] ErrCode = %s,ErrorDescription = %s", r.ErrorCode, r.ErrorDescription)
}
