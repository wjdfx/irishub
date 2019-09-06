//nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

// Rand errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "htlc"

	CodeInvalidAddress              sdk.CodeType = 100
	CodeInvalidAmount               sdk.CodeType = 101
	CodeInvalidSecretHashLock       sdk.CodeType = 102
	CodeSecretHashLockAlreadyExists sdk.CodeType = 103
	CodeInvalidTimeLock             sdk.CodeType = 104
	CodeInvalidSecret               sdk.CodeType = 105
	CodeStateIsNotOpen              sdk.CodeType = 106
	CodeStateIsNotExpired           sdk.CodeType = 107
)

//----------------------------------------
// HTLC error constructors

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrInvalidAmount(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAmount, msg)
}

func ErrInvalidSecretHashLock(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSecretHashLock, msg)
}

func ErrSecretHashLockAlreadyExists(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeSecretHashLockAlreadyExists, msg)
}

func ErrInvalidTimeLock(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTimeLock, msg)
}

func ErrInvalidSecret(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSecret, msg)
}

func ErrStateIsNotOpen(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeStateIsNotOpen, msg)
}

func ErrStateIsNotExpired(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeStateIsNotExpired, msg)
}