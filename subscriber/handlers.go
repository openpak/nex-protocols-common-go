package subscriber

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	subscriber "github.com/PretendoNetwork/nex-protocols-go/v2/subscriber"
	subscriber_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscriber/types"
)

func (commonProtocol *CommonProtocol) hello(err error, packet nex.PacketInterface, callID uint32, unknown types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.Hello == nil {
		return nil, rmc.NotImplemented("Subscriber::Hello")
	}

	if nexError := commonProtocol.Hello(packet, unknown); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodHello, callID), nil
}

func (commonProtocol *CommonProtocol) postContent(err error, packet nex.PacketInterface, callID uint32, param subscriber_types.SubscriberPostContentParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.PostContent == nil {
		return nil, rmc.NotImplemented("Subscriber::PostContent")
	}

	result0, nexError := commonProtocol.PostContent(packet, param)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodPostContent, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getContent(err error, packet nex.PacketInterface, callID uint32, param subscriber_types.SubscriberGetContentParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetContent == nil {
		return nil, rmc.NotImplemented("Subscriber::GetContent")
	}

	result0, nexError := commonProtocol.GetContent(packet, param)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodGetContent, callID, result0), nil
}

func (commonProtocol *CommonProtocol) follow(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.Follow == nil {
		return nil, rmc.NotImplemented("Subscriber::Follow")
	}

	body, nexError := commonProtocol.Follow(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, subscriber.ProtocolID, subscriber.MethodFollow, callID, body), nil
}

func (commonProtocol *CommonProtocol) unfollowAllAndFollow(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UnfollowAllAndFollow == nil {
		return nil, rmc.NotImplemented("Subscriber::UnfollowAllAndFollow")
	}

	body, nexError := commonProtocol.UnfollowAllAndFollow(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, subscriber.ProtocolID, subscriber.MethodUnfollowAllAndFollow, callID, body), nil
}

func (commonProtocol *CommonProtocol) unfollow(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.Unfollow == nil {
		return nil, rmc.NotImplemented("Subscriber::Unfollow")
	}

	body, nexError := commonProtocol.Unfollow(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, subscriber.ProtocolID, subscriber.MethodUnfollow, callID, body), nil
}

func (commonProtocol *CommonProtocol) getFollowing(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetFollowing == nil {
		return nil, rmc.NotImplemented("Subscriber::GetFollowing")
	}

	body, nexError := commonProtocol.GetFollowing(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, subscriber.ProtocolID, subscriber.MethodGetFollowing, callID, body), nil
}

func (commonProtocol *CommonProtocol) getFollower(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetFollower == nil {
		return nil, rmc.NotImplemented("Subscriber::GetFollower")
	}

	body, nexError := commonProtocol.GetFollower(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, subscriber.ProtocolID, subscriber.MethodGetFollower, callID, body), nil
}

func (commonProtocol *CommonProtocol) getNumFollowers(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetNumFollowers == nil {
		return nil, rmc.NotImplemented("Subscriber::GetNumFollowers")
	}

	body, nexError := commonProtocol.GetNumFollowers(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, subscriber.ProtocolID, subscriber.MethodGetNumFollowers, callID, body), nil
}

func (commonProtocol *CommonProtocol) getTimeline(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetTimeline == nil {
		return nil, rmc.NotImplemented("Subscriber::GetTimeline")
	}

	body, nexError := commonProtocol.GetTimeline(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, subscriber.ProtocolID, subscriber.MethodGetTimeline, callID, body), nil
}

func (commonProtocol *CommonProtocol) deleteContent(err error, packet nex.PacketInterface, callID uint32, topics types.List[types.String], contentID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.DeleteContent == nil {
		return nil, rmc.NotImplemented("Subscriber::DeleteContent")
	}

	if nexError := commonProtocol.DeleteContent(packet, topics, contentID); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodDeleteContent, callID), nil
}

func (commonProtocol *CommonProtocol) getContentMulti(err error, packet nex.PacketInterface, callID uint32, params types.List[subscriber_types.SubscriberGetContentParam]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetContentMulti == nil {
		return nil, rmc.NotImplemented("Subscriber::GetContentMulti")
	}

	result0, nexError := commonProtocol.GetContentMulti(packet, params)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodGetContentMulti, callID, result0), nil
}

func (commonProtocol *CommonProtocol) updateUserStatus(err error, packet nex.PacketInterface, callID uint32, param types.List[subscriber_types.SubscriberUserStatusParam], isNotify types.List[types.UInt8]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateUserStatus == nil {
		return nil, rmc.NotImplemented("Subscriber::UpdateUserStatus")
	}

	if nexError := commonProtocol.UpdateUserStatus(packet, param, isNotify); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodUpdateUserStatus, callID), nil
}

func (commonProtocol *CommonProtocol) getFriendUserStatuses(err error, packet nex.PacketInterface, callID uint32, keys types.List[types.UInt8]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetFriendUserStatuses == nil {
		return nil, rmc.NotImplemented("Subscriber::GetFriendUserStatuses")
	}

	result0, nexError := commonProtocol.GetFriendUserStatuses(packet, keys)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodGetFriendUserStatuses, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getUserStatuses(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.PID], keys types.List[types.UInt8]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetUserStatuses == nil {
		return nil, rmc.NotImplemented("Subscriber::GetUserStatuses")
	}

	result0, nexError := commonProtocol.GetUserStatuses(packet, pids, keys)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, subscriber.ProtocolID, subscriber.MethodGetUserStatuses, callID, result0), nil
}
