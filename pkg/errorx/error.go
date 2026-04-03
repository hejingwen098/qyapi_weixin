package errorx

import "fmt"

type QyError struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

func (e *QyError) Error() string {
	return fmt.Sprintf("企业微信错误 [code:%d]: %s", e.Code, e.Msg)
}

var (
	ErrInvalidToken       = &QyError{Code: 40014, Msg: "invalid access_token"}
	ErrTokenExpired       = &QyError{Code: 42001, Msg: "access_token expired"}
	ErrDepartmentNotFound = &QyError{Code: 60001, Msg: "department not found"}
)
