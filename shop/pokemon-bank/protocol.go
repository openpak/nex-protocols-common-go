// Package shop_pokemon_bank implements hook-driven handlers for the Pokemon Bank Shop protocol
package shop_pokemon_bank

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	shop_pokemon_bank "github.com/PretendoNetwork/nex-protocols-go/v2/shop/pokemon-bank"
	shop_types "github.com/PretendoNetwork/nex-protocols-go/v2/shop/types"
)

// CommonProtocol handles the Shop (Pokemon Bank) protocol. Every method calls the
// matching hook, which must be set by the game server; methods without a hook
// respond with NotImplemented.
type CommonProtocol struct {
	endpoint                   nex.EndpointInterface
	protocol                   shop_pokemon_bank.Interface
	GetItems                   func(requesterPID types.PID) (types.List[shop_types.ShopItem], *nex.Error)
	GetChallengeBlob           func(requesterPID types.PID) (types.Buffer, *nex.Error)
	GetRivToken                func(requesterPID types.PID, itemCode types.String, referenceID types.QBuffer) (types.String, *nex.Error)
	GetRivTokenByItemID        func(requesterPID types.PID, itemID types.UInt32) (types.String, *nex.Error)
	GetItemRights              func(requesterPID types.PID, ticketEnvelope types.QBuffer) (types.List[shop_types.ShopItemRights], *nex.Error)
	VerifyAndRegisterTicket    func(requesterPID types.PID, ticketEnvelope types.QBuffer, purchasedTime types.DateTime) (types.DateTime, *nex.Error)
	DebugSetExpireTime         func(requesterPID types.PID, expireTime types.DateTime) *nex.Error
	PrincipalIDToSupportNumber func(requesterPID types.PID, pid types.UInt32) (types.String, *nex.Error)
	SupportNumberToPrincipalID func(requesterPID types.PID, supportNumber types.String) (types.UInt32, *nex.Error)
	GetGameServerTime          func(requesterPID types.PID) (types.DateTime, *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol shop_pokemon_bank.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetItems(commonProtocol.getItems)
	protocol.SetHandlerGetChallengeBlob(commonProtocol.getChallengeBlob)
	protocol.SetHandlerGetRivToken(commonProtocol.getRivToken)
	protocol.SetHandlerGetRivTokenByItemID(commonProtocol.getRivTokenByItemID)
	protocol.SetHandlerGetItemRights(commonProtocol.getItemRights)
	protocol.SetHandlerVerifyAndRegisterTicket(commonProtocol.verifyAndRegisterTicket)
	protocol.SetHandlerDebugSetExpireTime(commonProtocol.debugSetExpireTime)
	protocol.SetHandlerPrincipalIDToSupportNumber(commonProtocol.principalIDToSupportNumber)
	protocol.SetHandlerSupportNumberToPrincipalID(commonProtocol.supportNumberToPrincipalID)
	protocol.SetHandlerGetGameServerTime(commonProtocol.getGameServerTime)

	return commonProtocol
}
