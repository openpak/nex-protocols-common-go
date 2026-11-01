// Package rmc holds the small helpers shared by the hook-driven common protocols:
// building a success response from a list of values, and the two standard errors.
package rmc

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
)

// Response encodes values back to back as the response body of a successful call
func Response(endpoint nex.EndpointInterface, protocolID uint16, methodID uint32, callID uint32, values ...types.RVType) *nex.RMCMessage {
	var body []byte

	if len(values) > 0 {
		stream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

		for _, value := range values {
			value.WriteTo(stream)
		}

		body = stream.Bytes()
	}

	response := nex.NewRMCSuccess(endpoint, body)
	response.ProtocolID = protocolID
	response.MethodID = methodID
	response.CallID = callID

	return response
}

// Respond builds the response for the endpoint the packet arrived on
func Respond(packet nex.PacketInterface, protocolID uint16, methodID uint32, callID uint32, values ...types.RVType) *nex.RMCMessage {
	return Response(packet.Sender().Endpoint(), protocolID, methodID, callID, values...)
}

// InvalidArgument logs a request that failed to decode and returns the matching error
func InvalidArgument(err error) *nex.Error {
	common_globals.Logger.Error(err.Error())

	return nex.NewError(nex.ResultCodes.Core.InvalidArgument, err.Error())
}

// NotImplemented returns Core::NotImplemented for a method whose hook was not set
func NotImplemented(method string) *nex.Error {
	message := fmt.Sprintf("%s not implemented", method)

	common_globals.Logger.Warning(message)

	return nex.NewError(nex.ResultCodes.Core.NotImplemented, message)
}

// RespondRaw answers with an already encoded body, for methods whose response layout is not documented
func RespondRaw(packet nex.PacketInterface, protocolID uint16, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	if len(body) == 0 {
		body = nil
	}

	response := nex.NewRMCSuccess(packet.Sender().Endpoint(), body)
	response.ProtocolID = protocolID
	response.MethodID = methodID
	response.CallID = callID

	return response
}
