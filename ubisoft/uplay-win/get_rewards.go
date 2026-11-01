package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
	uplay_win_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win/types"
)

// buildGetRewardsResponse writes the response of Uplay Win::GetRewards to the given stream
func buildGetRewardsResponse(stream *nex.ByteStreamOut, rewardList types.List[uplay_win_types.UplayReward]) {
	rewardList.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getRewards(err error, packet nex.PacketInterface, callID uint32, startRowIndex types.Int32, maximumRows types.Int32, sortExpression types.String, cultureName types.String, platformCode types.String, gameCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetRewards == nil {
		common_globals.Logger.Warning("UplayWin::GetRewards not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respRewardList, nexError := commonProtocol.GetRewards(packet, startRowIndex, maximumRows, sortExpression, cultureName, platformCode, gameCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetRewardsResponse(rmcResponseStream, respRewardList)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = uplay_win.ProtocolID
	rmcResponse.MethodID = uplay_win.MethodGetRewards
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
