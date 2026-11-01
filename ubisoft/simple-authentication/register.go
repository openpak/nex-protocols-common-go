package simpleauthentication

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	simple_authentication "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/simple-authentication"
)

// buildRegisterResponse writes the response of Simple Authentication::Register to the given stream
func buildRegisterResponse(stream *nex.ByteStreamOut, retval types.QResult, pidConnectionID types.UInt32, urlPublic types.StationURL) {
	retval.WriteTo(stream)
	pidConnectionID.WriteTo(stream)
	urlPublic.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) register(err error, packet nex.PacketInterface, callID uint32, vecMyURLs types.List[types.StationURL]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.Register == nil {
		common_globals.Logger.Warning("SimpleAuthentication::Register not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respRetval, respPidConnectionID, respUrlPublic, nexError := commonProtocol.Register(packet, vecMyURLs)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildRegisterResponse(rmcResponseStream, respRetval, respPidConnectionID, respUrlPublic)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = simple_authentication.ProtocolID
	rmcResponse.MethodID = simple_authentication.MethodRegister
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
