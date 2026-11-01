// Package uplaywin provides the common handlers for the Uplay Win protocol (Rendez-Vous/Quazal, protocol ID 49).
//
// Every method decodes its request and encodes its response as documented on the NintendoClients wiki
// page "Uplay Win Protocol". The actual behaviour is supplied through the hook fields below; a method
// whose hook is nil answers Core::NotImplemented.
package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
	uplay_win_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win/types"
)

// CommonProtocol handles the Uplay Win protocol on an endpoint.
// Each hook receives the decoded request and returns the values of the response in wiki order.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol uplay_win.Interface

	// GetActions is called for Uplay Win::GetActions
	GetActions func(packet nex.PacketInterface, startRowIndex types.Int32, maximumRows types.Int32, sortExpression types.String, cultureName types.String, platformCode types.String, gameCode types.String) (respActionList types.List[uplay_win_types.UplayAction], nexError *nex.Error)

	// GetActionsCompleted is called for Uplay Win::GetActionsCompleted
	GetActionsCompleted func(packet nex.PacketInterface, startRowIndex types.Int32, maximumRows types.Int32, sortExpression types.String, cultureName types.String, platformCode types.String, gameCode types.String) (respActionList types.List[uplay_win_types.UplayAction], nexError *nex.Error)

	// GetActionsCount is called for Uplay Win::GetActionsCount
	GetActionsCount func(packet nex.PacketInterface, platformCode types.String, gameCode types.String) (respActionsCount types.Int32, nexError *nex.Error)

	// GetActionsCompletedCount is called for Uplay Win::GetActionsCompletedCount
	GetActionsCompletedCount func(packet nex.PacketInterface, platformCode types.String, gameCode types.String) (respActionsCount types.Int32, nexError *nex.Error)

	// GetRewards is called for Uplay Win::GetRewards
	GetRewards func(packet nex.PacketInterface, startRowIndex types.Int32, maximumRows types.Int32, sortExpression types.String, cultureName types.String, platformCode types.String, gameCode types.String) (respRewardList types.List[uplay_win_types.UplayReward], nexError *nex.Error)

	// GetRewardsPurchased is called for Uplay Win::GetRewardsPurchased
	GetRewardsPurchased func(packet nex.PacketInterface, startRowIndex types.Int32, maximumRows types.Int32, sortExpression types.String, cultureName types.String, platformCode types.String, gameCode types.String) (respRewardList types.List[uplay_win_types.UplayReward], nexError *nex.Error)

	// UplayWelcome is called for Uplay Win::UplayWelcome
	UplayWelcome func(packet nex.PacketInterface, cultureName types.String, platformCode types.String, gameCode types.String) (respActionList types.List[uplay_win_types.UplayAction], nexError *nex.Error)

	// SetActionCompleted is called for Uplay Win::SetActionCompleted
	SetActionCompleted func(packet nex.PacketInterface, actionCode types.String, cultureName types.String, platformCode types.String, gameCode types.String) (respUnlockedAction uplay_win_types.UplayAction, nexError *nex.Error)

	// SetActionsCompleted is called for Uplay Win::SetActionsCompleted
	SetActionsCompleted func(packet nex.PacketInterface, actionCodeList types.List[types.String], cultureName types.String, platformCode types.String, gameCode types.String) (respActionList types.List[uplay_win_types.UplayAction], nexError *nex.Error)

	// GetUserToken is called for Uplay Win::GetUserToken
	GetUserToken func(packet nex.PacketInterface) (respToken types.String, nexError *nex.Error)

	// GetVirtualCurrencyUserBalance is called for Uplay Win::GetVirtualCurrencyUserBalance
	GetVirtualCurrencyUserBalance func(packet nex.PacketInterface, platformCode types.String) (respVirtualCurrencyUserBalance types.Int32, nexError *nex.Error)

	// GetSectionsByKey is called for Uplay Win::GetSectionsByKey
	GetSectionsByKey func(packet nex.PacketInterface, cultureName types.String, sectionKey types.String, platformCode types.String, gameCode types.String) (respSectionList types.List[uplay_win_types.UplaySection], nexError *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol uplay_win.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetActions(commonProtocol.getActions)
	protocol.SetHandlerGetActionsCompleted(commonProtocol.getActionsCompleted)
	protocol.SetHandlerGetActionsCount(commonProtocol.getActionsCount)
	protocol.SetHandlerGetActionsCompletedCount(commonProtocol.getActionsCompletedCount)
	protocol.SetHandlerGetRewards(commonProtocol.getRewards)
	protocol.SetHandlerGetRewardsPurchased(commonProtocol.getRewardsPurchased)
	protocol.SetHandlerUplayWelcome(commonProtocol.uplayWelcome)
	protocol.SetHandlerSetActionCompleted(commonProtocol.setActionCompleted)
	protocol.SetHandlerSetActionsCompleted(commonProtocol.setActionsCompleted)
	protocol.SetHandlerGetUserToken(commonProtocol.getUserToken)
	protocol.SetHandlerGetVirtualCurrencyUserBalance(commonProtocol.getVirtualCurrencyUserBalance)
	protocol.SetHandlerGetSectionsByKey(commonProtocol.getSectionsByKey)

	return commonProtocol
}
