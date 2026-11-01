package subscriber

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	subscriber_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscriber/types"
)

// TestHooksAbsentAnswerNotImplemented checks that every method without a hook answers Core::NotImplemented
func TestHooksAbsentAnswerNotImplemented(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"Hello": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.hello(nil, packet, 1, *new(types.String)) },
		"PostContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.postContent(nil, packet, 1, *new(subscriber_types.SubscriberPostContentParam))
		},
		"GetContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getContent(nil, packet, 1, *new(subscriber_types.SubscriberGetContentParam))
		},
		"Follow": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.follow(nil, packet, 1, *new([]byte)) },
		"UnfollowAllAndFollow": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.unfollowAllAndFollow(nil, packet, 1, *new([]byte))
		},
		"Unfollow":     func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.unfollow(nil, packet, 1, *new([]byte)) },
		"GetFollowing": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getFollowing(nil, packet, 1, *new([]byte)) },
		"GetFollower":  func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getFollower(nil, packet, 1, *new([]byte)) },
		"GetNumFollowers": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getNumFollowers(nil, packet, 1, *new([]byte))
		},
		"GetTimeline": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getTimeline(nil, packet, 1, *new([]byte)) },
		"DeleteContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.deleteContent(nil, packet, 1, *new(types.List[types.String]), *new(types.UInt64))
		},
		"GetContentMulti": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getContentMulti(nil, packet, 1, *new(types.List[subscriber_types.SubscriberGetContentParam]))
		},
		"UpdateUserStatus": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateUserStatus(nil, packet, 1, *new(types.List[subscriber_types.SubscriberUserStatusParam]), *new(types.List[types.UInt8]))
		},
		"GetFriendUserStatuses": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getFriendUserStatuses(nil, packet, 1, *new(types.List[types.UInt8]))
		},
		"GetUserStatuses": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getUserStatuses(nil, packet, 1, *new(types.List[types.PID]), *new(types.List[types.UInt8]))
		},
	}

	if len(calls) != 15 {
		t.Fatalf("expected 15 methods, got %d", len(calls))
	}

	for name, call := range calls {
		response, nexError := call()
		if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.NotImplemented) {
			t.Errorf("%s: expected Core::NotImplemented, got %v", name, nexError)
		}
	}
}

// TestDecodeErrorsAreInvalidArgument checks that a request that failed to decode never reaches a hook
func TestDecodeErrorsAreInvalidArgument(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"Hello": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.hello(errBadRequest, packet, 1, *new(types.String))
		},
		"PostContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.postContent(errBadRequest, packet, 1, *new(subscriber_types.SubscriberPostContentParam))
		},
		"GetContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getContent(errBadRequest, packet, 1, *new(subscriber_types.SubscriberGetContentParam))
		},
		"Follow": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.follow(errBadRequest, packet, 1, *new([]byte))
		},
		"UnfollowAllAndFollow": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.unfollowAllAndFollow(errBadRequest, packet, 1, *new([]byte))
		},
		"Unfollow": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.unfollow(errBadRequest, packet, 1, *new([]byte))
		},
		"GetFollowing": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getFollowing(errBadRequest, packet, 1, *new([]byte))
		},
		"GetFollower": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getFollower(errBadRequest, packet, 1, *new([]byte))
		},
		"GetNumFollowers": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getNumFollowers(errBadRequest, packet, 1, *new([]byte))
		},
		"GetTimeline": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getTimeline(errBadRequest, packet, 1, *new([]byte))
		},
		"DeleteContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.deleteContent(errBadRequest, packet, 1, *new(types.List[types.String]), *new(types.UInt64))
		},
		"GetContentMulti": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getContentMulti(errBadRequest, packet, 1, *new(types.List[subscriber_types.SubscriberGetContentParam]))
		},
		"UpdateUserStatus": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateUserStatus(errBadRequest, packet, 1, *new(types.List[subscriber_types.SubscriberUserStatusParam]), *new(types.List[types.UInt8]))
		},
		"GetFriendUserStatuses": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getFriendUserStatuses(errBadRequest, packet, 1, *new(types.List[types.UInt8]))
		},
		"GetUserStatuses": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getUserStatuses(errBadRequest, packet, 1, *new(types.List[types.PID]), *new(types.List[types.UInt8]))
		},
	}

	for name, call := range calls {
		response, nexError := call()
		if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
			t.Errorf("%s: expected Core::InvalidArgument, got %v", name, nexError)
		}
	}
}

var errBadRequest = errors.New("bad request")
