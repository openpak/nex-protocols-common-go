// Package friends provides the hook driven common handlers for the Friends (0x14), the generic Rendez-Vous friends protocol (not the Wii U / 3DS friends protocols). Every method is hook driven protocol
package friends

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends "github.com/PretendoNetwork/nex-protocols-go/v2/friends"
	friends_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends/types"
)

// CommonProtocol implements the Friends protocol on top of hooks. A method whose hook is nil answers Core::NotImplemented
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol friends.Interface

	// AddFriend is called for Friends::AddFriend
	AddFriend func(packet nex.PacketInterface, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (types.Bool, *nex.Error)
	// AddFriendByName is called for Friends::AddFriendByName
	AddFriendByName func(packet nex.PacketInterface, strPlayerName types.String, uiDetails types.UInt32, strMessage types.String) (types.Bool, *nex.Error)
	// AddFriendWithDetails is called for Friends::AddFriendWithDetails
	AddFriendWithDetails func(packet nex.PacketInterface, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (friends_types.RelationshipData, *nex.Error)
	// AddFriendByNameWithDetails is called for Friends::AddFriendByNameWithDetails
	AddFriendByNameWithDetails func(packet nex.PacketInterface, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (friends_types.RelationshipData, *nex.Error)
	// AcceptFriendship is called for Friends::AcceptFriendship
	AcceptFriendship func(packet nex.PacketInterface, uiPlayer types.UInt32) (types.Bool, *nex.Error)
	// DeclineFriendship is called for Friends::DeclineFriendship
	DeclineFriendship func(packet nex.PacketInterface, uiPlayer types.UInt32) (types.Bool, *nex.Error)
	// BlackList is called for Friends::BlackList
	BlackList func(packet nex.PacketInterface, uiPlayer types.UInt32, uiDetails types.UInt32) (types.Bool, *nex.Error)
	// BlackListByName is called for Friends::BlackListByName
	BlackListByName func(packet nex.PacketInterface, strPlayerName types.String, uiDetails types.UInt32) (types.Bool, *nex.Error)
	// ClearRelationship is called for Friends::ClearRelationship
	ClearRelationship func(packet nex.PacketInterface, uiPlayer types.UInt32) (types.Bool, *nex.Error)
	// UpdateDetails is called for Friends::UpdateDetails
	UpdateDetails func(packet nex.PacketInterface, uiPlayer types.UInt32, uiDetails types.UInt32) (types.Bool, *nex.Error)
	// GetList is called for Friends::GetList
	GetList func(packet nex.PacketInterface, byRelationship types.UInt8, bReversed types.Bool) (types.List[types.UInt32], *nex.Error)
	// GetDetailedList is called for Friends::GetDetailedList
	GetDetailedList func(packet nex.PacketInterface, byRelationship types.UInt8, bReversed types.Bool) (types.List[friends_types.FriendData], *nex.Error)
	// GetRelationships is called for Friends::GetRelationships
	GetRelationships func(packet nex.PacketInterface, resultRange types.ResultRange) (types.UInt32, types.List[friends_types.RelationshipData], *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol friends.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerAddFriend(commonProtocol.addFriend)
	protocol.SetHandlerAddFriendByName(commonProtocol.addFriendByName)
	protocol.SetHandlerAddFriendWithDetails(commonProtocol.addFriendWithDetails)
	protocol.SetHandlerAddFriendByNameWithDetails(commonProtocol.addFriendByNameWithDetails)
	protocol.SetHandlerAcceptFriendship(commonProtocol.acceptFriendship)
	protocol.SetHandlerDeclineFriendship(commonProtocol.declineFriendship)
	protocol.SetHandlerBlackList(commonProtocol.blackList)
	protocol.SetHandlerBlackListByName(commonProtocol.blackListByName)
	protocol.SetHandlerClearRelationship(commonProtocol.clearRelationship)
	protocol.SetHandlerUpdateDetails(commonProtocol.updateDetails)
	protocol.SetHandlerGetList(commonProtocol.getList)
	protocol.SetHandlerGetDetailedList(commonProtocol.getDetailedList)
	protocol.SetHandlerGetRelationships(commonProtocol.getRelationships)

	return commonProtocol
}
