package simpleauthentication

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	simple_authentication "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/simple-authentication"
)

// buildLoginWithTokenCafeExResponse writes the response of Simple Authentication::LoginWithTokenCafeEx to the given stream
func buildLoginWithTokenCafeExResponse(stream *nex.ByteStreamOut, retval types.QResult, pidPrincipal types.UInt32, pConnectionData types.RVConnectionData, strReturnMsg types.String) {
	retval.WriteTo(stream)
	pidPrincipal.WriteTo(stream)
	pConnectionData.WriteTo(stream)
	strReturnMsg.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) loginWithTokenCafeEx(err error, packet nex.PacketInterface, callID uint32, strNintendoToken types.String, oAnyData types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LoginWithTokenCafeEx == nil {
		common_globals.Logger.Warning("SimpleAuthentication::LoginWithTokenCafeEx not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respRetval, respPidPrincipal, respPConnectionData, respStrReturnMsg, nexError := commonProtocol.LoginWithTokenCafeEx(packet, strNintendoToken, oAnyData)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLoginWithTokenCafeExResponse(rmcResponseStream, respRetval, respPidPrincipal, respPConnectionData, respStrReturnMsg)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = simple_authentication.ProtocolID
	rmcResponse.MethodID = simple_authentication.MethodLoginWithTokenCafeEx
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
