package friends

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	friends "github.com/PretendoNetwork/nex-protocols-go/v2/friends"
)

func (commonProtocol *CommonProtocol) addFriend(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.AddFriend == nil {
		return nil, rmc.NotImplemented("Friends::AddFriend")
	}

	result0, nexError := commonProtocol.AddFriend(packet, uiPlayer, uiDetails, strMessage)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodAddFriend, callID, result0), nil
}

func (commonProtocol *CommonProtocol) addFriendByName(err error, packet nex.PacketInterface, callID uint32, strPlayerName types.String, uiDetails types.UInt32, strMessage types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.AddFriendByName == nil {
		return nil, rmc.NotImplemented("Friends::AddFriendByName")
	}

	result0, nexError := commonProtocol.AddFriendByName(packet, strPlayerName, uiDetails, strMessage)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodAddFriendByName, callID, result0), nil
}

func (commonProtocol *CommonProtocol) addFriendWithDetails(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.AddFriendWithDetails == nil {
		return nil, rmc.NotImplemented("Friends::AddFriendWithDetails")
	}

	result0, nexError := commonProtocol.AddFriendWithDetails(packet, uiPlayer, uiDetails, strMessage)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodAddFriendWithDetails, callID, result0), nil
}

func (commonProtocol *CommonProtocol) addFriendByNameWithDetails(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.AddFriendByNameWithDetails == nil {
		return nil, rmc.NotImplemented("Friends::AddFriendByNameWithDetails")
	}

	result0, nexError := commonProtocol.AddFriendByNameWithDetails(packet, uiPlayer, uiDetails, strMessage)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodAddFriendByNameWithDetails, callID, result0), nil
}

func (commonProtocol *CommonProtocol) acceptFriendship(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.AcceptFriendship == nil {
		return nil, rmc.NotImplemented("Friends::AcceptFriendship")
	}

	result0, nexError := commonProtocol.AcceptFriendship(packet, uiPlayer)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodAcceptFriendship, callID, result0), nil
}

func (commonProtocol *CommonProtocol) declineFriendship(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.DeclineFriendship == nil {
		return nil, rmc.NotImplemented("Friends::DeclineFriendship")
	}

	result0, nexError := commonProtocol.DeclineFriendship(packet, uiPlayer)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodDeclineFriendship, callID, result0), nil
}

func (commonProtocol *CommonProtocol) blackList(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32, uiDetails types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.BlackList == nil {
		return nil, rmc.NotImplemented("Friends::BlackList")
	}

	result0, nexError := commonProtocol.BlackList(packet, uiPlayer, uiDetails)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodBlackList, callID, result0), nil
}

func (commonProtocol *CommonProtocol) blackListByName(err error, packet nex.PacketInterface, callID uint32, strPlayerName types.String, uiDetails types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.BlackListByName == nil {
		return nil, rmc.NotImplemented("Friends::BlackListByName")
	}

	result0, nexError := commonProtocol.BlackListByName(packet, strPlayerName, uiDetails)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodBlackListByName, callID, result0), nil
}

func (commonProtocol *CommonProtocol) clearRelationship(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.ClearRelationship == nil {
		return nil, rmc.NotImplemented("Friends::ClearRelationship")
	}

	result0, nexError := commonProtocol.ClearRelationship(packet, uiPlayer)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodClearRelationship, callID, result0), nil
}

func (commonProtocol *CommonProtocol) updateDetails(err error, packet nex.PacketInterface, callID uint32, uiPlayer types.UInt32, uiDetails types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateDetails == nil {
		return nil, rmc.NotImplemented("Friends::UpdateDetails")
	}

	result0, nexError := commonProtocol.UpdateDetails(packet, uiPlayer, uiDetails)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodUpdateDetails, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getList(err error, packet nex.PacketInterface, callID uint32, byRelationship types.UInt8, bReversed types.Bool) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetList == nil {
		return nil, rmc.NotImplemented("Friends::GetList")
	}

	result0, nexError := commonProtocol.GetList(packet, byRelationship, bReversed)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodGetList, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getDetailedList(err error, packet nex.PacketInterface, callID uint32, byRelationship types.UInt8, bReversed types.Bool) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetDetailedList == nil {
		return nil, rmc.NotImplemented("Friends::GetDetailedList")
	}

	result0, nexError := commonProtocol.GetDetailedList(packet, byRelationship, bReversed)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodGetDetailedList, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getRelationships(err error, packet nex.PacketInterface, callID uint32, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetRelationships == nil {
		return nil, rmc.NotImplemented("Friends::GetRelationships")
	}

	result0, result1, nexError := commonProtocol.GetRelationships(packet, resultRange)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, friends.ProtocolID, friends.MethodGetRelationships, callID, result0, result1), nil
}
