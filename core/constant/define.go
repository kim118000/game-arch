package constant

import "errors"

const (
	CONN_CTX_KEY        = "CONN_CTX"
	HAND_SHAKE_ARR_LEN  = 20
	HAND_SHAKE_SIGN_LEN = 4
	HAND_SHAKE_SUC      = "HAND_SHAKE_SUC"
)

var (
	ErrWrongConnAssert   = errors.New("conn assert failure")
	ErrWrongMessageAssert = errors.New("message assert failure")
	ErrWrongAuthFailure = errors.New("session auth failure")
	ErrWrongProtobufAssert = errors.New("protobuf assert failure")

	ErrWrongConnShake     = errors.New("err wrong conn shake")
	ErrWrongDatapackLength = errors.New("err wrong data pack lenght")
	ErrWrongValueType = errors.New("protobuf: convert on wrong type value")
)
