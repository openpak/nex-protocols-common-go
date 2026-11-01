package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
	uplay_win_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win/types"
)

// buildGetSectionsByKeyResponse writes the response of Uplay Win::GetSectionsByKey to the given stream
func buildGetSectionsByKeyResponse(stream *nex.ByteStreamOut, sectionList types.List[uplay_win_types.UplaySection]) {
	sectionList.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getSectionsByKey(err error, packet nex.PacketInterface, callID uint32, cultureName types.String, sectionKey types.String, platformCode types.String, gameCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetSectionsByKey == nil {
		common_globals.Logger.Warning("UplayWin::GetSectionsByKey not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respSectionList, nexError := commonProtocol.GetSectionsByKey(packet, cultureName, sectionKey, platformCode, gameCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetSectionsByKeyResponse(rmcResponseStream, respSectionList)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = uplay_win.ProtocolID
	rmcResponse.MethodID = uplay_win.MethodGetSectionsByKey
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
