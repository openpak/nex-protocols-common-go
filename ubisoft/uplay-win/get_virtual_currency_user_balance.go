package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
)

// buildGetVirtualCurrencyUserBalanceResponse writes the response of Uplay Win::GetVirtualCurrencyUserBalance to the given stream
func buildGetVirtualCurrencyUserBalanceResponse(stream *nex.ByteStreamOut, virtualCurrencyUserBalance types.Int32) {
	virtualCurrencyUserBalance.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getVirtualCurrencyUserBalance(err error, packet nex.PacketInterface, callID uint32, platformCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetVirtualCurrencyUserBalance == nil {
		common_globals.Logger.Warning("UplayWin::GetVirtualCurrencyUserBalance not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respVirtualCurrencyUserBalance, nexError := commonProtocol.GetVirtualCurrencyUserBalance(packet, platformCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetVirtualCurrencyUserBalanceResponse(rmcResponseStream, respVirtualCurrencyUserBalance)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = uplay_win.ProtocolID
	rmcResponse.MethodID = uplay_win.MethodGetVirtualCurrencyUserBalance
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
