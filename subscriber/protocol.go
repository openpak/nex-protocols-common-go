// Package subscriber provides the hook driven common handlers for the Subscriber (0x79). Every method is hook driven. The Follow family and GetTimeline have no documented layout, their hooks get the raw request parameters and return the raw response body protocol
package subscriber

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	subscriber "github.com/PretendoNetwork/nex-protocols-go/v2/subscriber"
	subscriber_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscriber/types"
)

// CommonProtocol implements the Subscriber protocol on top of hooks. A method whose hook is nil answers Core::NotImplemented
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol subscriber.Interface

	// Hello is called for Subscriber::Hello
	Hello func(packet nex.PacketInterface, unknown types.String) *nex.Error
	// PostContent is called for Subscriber::PostContent
	PostContent func(packet nex.PacketInterface, param subscriber_types.SubscriberPostContentParam) (types.UInt64, *nex.Error)
	// GetContent is called for Subscriber::GetContent
	GetContent func(packet nex.PacketInterface, param subscriber_types.SubscriberGetContentParam) (types.List[subscriber_types.SubscriberContent], *nex.Error)
	// Follow is called for Subscriber::Follow
	Follow func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
	// UnfollowAllAndFollow is called for Subscriber::UnfollowAllAndFollow
	UnfollowAllAndFollow func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
	// Unfollow is called for Subscriber::Unfollow
	Unfollow func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
	// GetFollowing is called for Subscriber::GetFollowing
	GetFollowing func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
	// GetFollower is called for Subscriber::GetFollower
	GetFollower func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
	// GetNumFollowers is called for Subscriber::GetNumFollowers
	GetNumFollowers func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
	// GetTimeline is called for Subscriber::GetTimeline
	GetTimeline func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
	// DeleteContent is called for Subscriber::DeleteContent
	DeleteContent func(packet nex.PacketInterface, topics types.List[types.String], contentID types.UInt64) *nex.Error
	// GetContentMulti is called for Subscriber::GetContentMulti
	GetContentMulti func(packet nex.PacketInterface, params types.List[subscriber_types.SubscriberGetContentParam]) (types.List[types.List[subscriber_types.SubscriberContent]], *nex.Error)
	// UpdateUserStatus is called for Subscriber::UpdateUserStatus
	UpdateUserStatus func(packet nex.PacketInterface, param types.List[subscriber_types.SubscriberUserStatusParam], isNotify types.List[types.UInt8]) *nex.Error
	// GetFriendUserStatuses is called for Subscriber::GetFriendUserStatuses
	GetFriendUserStatuses func(packet nex.PacketInterface, keys types.List[types.UInt8]) (types.List[subscriber_types.SubscriberUserStatusInfo], *nex.Error)
	// GetUserStatuses is called for Subscriber::GetUserStatuses
	GetUserStatuses func(packet nex.PacketInterface, pids types.List[types.PID], keys types.List[types.UInt8]) (types.List[subscriber_types.SubscriberUserStatusInfo], *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol subscriber.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerHello(commonProtocol.hello)
	protocol.SetHandlerPostContent(commonProtocol.postContent)
	protocol.SetHandlerGetContent(commonProtocol.getContent)
	protocol.SetHandlerFollow(commonProtocol.follow)
	protocol.SetHandlerUnfollowAllAndFollow(commonProtocol.unfollowAllAndFollow)
	protocol.SetHandlerUnfollow(commonProtocol.unfollow)
	protocol.SetHandlerGetFollowing(commonProtocol.getFollowing)
	protocol.SetHandlerGetFollower(commonProtocol.getFollower)
	protocol.SetHandlerGetNumFollowers(commonProtocol.getNumFollowers)
	protocol.SetHandlerGetTimeline(commonProtocol.getTimeline)
	protocol.SetHandlerDeleteContent(commonProtocol.deleteContent)
	protocol.SetHandlerGetContentMulti(commonProtocol.getContentMulti)
	protocol.SetHandlerUpdateUserStatus(commonProtocol.updateUserStatus)
	protocol.SetHandlerGetFriendUserStatuses(commonProtocol.getFriendUserStatuses)
	protocol.SetHandlerGetUserStatuses(commonProtocol.getUserStatuses)

	return commonProtocol
}
