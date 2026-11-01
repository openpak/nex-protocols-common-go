// Package useraccountmanagement provides the common handlers for the User Account Management protocol (Rendez-Vous/Quazal, protocol ID 72).
//
// Every method decodes its request and encodes its response as documented on the NintendoClients wiki
// page "User Account Management Protocol". The actual behaviour is supplied through the hook fields below; a method
// whose hook is nil answers Core::NotImplemented.
package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	user_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-account-management"
)

// CommonProtocol handles the User Account Management protocol on an endpoint.
// Each hook receives the decoded request and returns the values of the response in wiki order.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol user_account_management.Interface

	// LookupSceNpIds is called for User Account Management::LookupSceNpIds
	LookupSceNpIds func(packet nex.PacketInterface, pids types.List[types.UInt32]) (respNpids types.Map[types.UInt32, types.QBuffer], nexError *nex.Error)

	// LookupPrincipalIDs is called for User Account Management::LookupPrincipalIDs
	LookupPrincipalIDs func(packet nex.PacketInterface, firstPartyIds types.List[types.String], platformId types.UInt32) (respPids types.Map[types.String, types.UInt32], nexError *nex.Error)

	// LookupFirstPartyIds is called for User Account Management::LookupFirstPartyIds
	LookupFirstPartyIds func(packet nex.PacketInterface, pids types.List[types.UInt32], platformId types.UInt32) (respFirstPartyIds types.Map[types.UInt32, types.String], nexError *nex.Error)

	// UserHasPlayed is called for User Account Management::UserHasPlayed
	UserHasPlayed func(packet nex.PacketInterface, firstPartyIds types.List[types.String], platformId types.UInt32) (respUserPresence types.Map[types.String, types.Bool], nexError *nex.Error)

	// IsUserPlaying is called for User Account Management::IsUserPlaying
	IsUserPlaying func(packet nex.PacketInterface, firstPartyIds types.List[types.String], platformId types.UInt32) (respUserPresence types.Map[types.String, types.Bool], nexError *nex.Error)

	// UpdateSonyAccountInfo is called for User Account Management::UpdateSonyAccountInfo
	UpdateSonyAccountInfo func(packet nex.PacketInterface, ticketData types.QBuffer, ticketSize types.UInt32) (nexError *nex.Error)

	// LookupUsernames is called for User Account Management::LookupUsernames
	LookupUsernames func(packet nex.PacketInterface, pids types.List[types.UInt32], platformId types.UInt32) (respUserNames types.Map[types.UInt32, types.String], nexError *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol user_account_management.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerLookupSceNpIds(commonProtocol.lookupSceNpIds)
	protocol.SetHandlerLookupPrincipalIDs(commonProtocol.lookupPrincipalIDs)
	protocol.SetHandlerLookupFirstPartyIds(commonProtocol.lookupFirstPartyIds)
	protocol.SetHandlerUserHasPlayed(commonProtocol.userHasPlayed)
	protocol.SetHandlerIsUserPlaying(commonProtocol.isUserPlaying)
	protocol.SetHandlerUpdateSonyAccountInfo(commonProtocol.updateSonyAccountInfo)
	protocol.SetHandlerLookupUsernames(commonProtocol.lookupUsernames)

	return commonProtocol
}
