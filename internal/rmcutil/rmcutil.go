// Package rmcutil holds small helpers shared by the hook-driven common protocols.
package rmcutil

import (
	"errors"

	"github.com/PretendoNetwork/nex-go/v2"
)

// HookError converts an error returned by a game supplied hook into a NEX error.
// A *nex.Error is passed through untouched so a hook can pick its own result
// code; any other error becomes fallback.
func HookError(err error, fallback uint32) *nex.Error {
	var nexErr *nex.Error
	if errors.As(err, &nexErr) {
		return nexErr
	}

	return nex.NewError(fallback, err.Error())
}

// Success builds a successful RMC response whose body is whatever write puts in
// the stream. write may be nil for an empty response.
func Success(endpoint nex.EndpointInterface, protocolID uint16, methodID uint32, callID uint32, write func(stream *nex.ByteStreamOut)) *nex.RMCMessage {
	var body []byte

	if write != nil {
		stream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())
		write(stream)
		body = stream.Bytes()
	}

	response := nex.NewRMCSuccess(endpoint, body)
	response.ProtocolID = protocolID
	response.MethodID = methodID
	response.CallID = callID

	return response
}
