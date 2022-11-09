package constant

import "errors"

var (
	ErrWrongConnAssert   = errors.New("conn assert failure")
	ErrWrongMessageAssert = errors.New("message assert failure")
	ErrWrongAuthFailure = errors.New("session auth failure")
	ErrWrongProtobufAssert = errors.New("protobuf assert failure")
)