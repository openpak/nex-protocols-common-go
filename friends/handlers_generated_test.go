package friends

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
)

// TestHooksAbsentAnswerNotImplemented checks that every method without a hook answers Core::NotImplemented
func TestHooksAbsentAnswerNotImplemented(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"AddFriend": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriend(nil, packet, 1, *new(types.UInt32), *new(types.UInt32), *new(types.String))
		},
		"AddFriendByName": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriendByName(nil, packet, 1, *new(types.String), *new(types.UInt32), *new(types.String))
		},
		"AddFriendWithDetails": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriendWithDetails(nil, packet, 1, *new(types.UInt32), *new(types.UInt32), *new(types.String))
		},
		"AddFriendByNameWithDetails": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriendByNameWithDetails(nil, packet, 1, *new(types.UInt32), *new(types.UInt32), *new(types.String))
		},
		"AcceptFriendship": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.acceptFriendship(nil, packet, 1, *new(types.UInt32))
		},
		"DeclineFriendship": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.declineFriendship(nil, packet, 1, *new(types.UInt32))
		},
		"BlackList": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.blackList(nil, packet, 1, *new(types.UInt32), *new(types.UInt32))
		},
		"BlackListByName": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.blackListByName(nil, packet, 1, *new(types.String), *new(types.UInt32))
		},
		"ClearRelationship": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.clearRelationship(nil, packet, 1, *new(types.UInt32))
		},
		"UpdateDetails": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateDetails(nil, packet, 1, *new(types.UInt32), *new(types.UInt32))
		},
		"GetList": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getList(nil, packet, 1, *new(types.UInt8), *new(types.Bool))
		},
		"GetDetailedList": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getDetailedList(nil, packet, 1, *new(types.UInt8), *new(types.Bool))
		},
		"GetRelationships": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getRelationships(nil, packet, 1, *new(types.ResultRange))
		},
	}

	if len(calls) != 13 {
		t.Fatalf("expected 13 methods, got %d", len(calls))
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
		"AddFriend": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriend(errBadRequest, packet, 1, *new(types.UInt32), *new(types.UInt32), *new(types.String))
		},
		"AddFriendByName": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriendByName(errBadRequest, packet, 1, *new(types.String), *new(types.UInt32), *new(types.String))
		},
		"AddFriendWithDetails": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriendWithDetails(errBadRequest, packet, 1, *new(types.UInt32), *new(types.UInt32), *new(types.String))
		},
		"AddFriendByNameWithDetails": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.addFriendByNameWithDetails(errBadRequest, packet, 1, *new(types.UInt32), *new(types.UInt32), *new(types.String))
		},
		"AcceptFriendship": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.acceptFriendship(errBadRequest, packet, 1, *new(types.UInt32))
		},
		"DeclineFriendship": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.declineFriendship(errBadRequest, packet, 1, *new(types.UInt32))
		},
		"BlackList": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.blackList(errBadRequest, packet, 1, *new(types.UInt32), *new(types.UInt32))
		},
		"BlackListByName": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.blackListByName(errBadRequest, packet, 1, *new(types.String), *new(types.UInt32))
		},
		"ClearRelationship": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.clearRelationship(errBadRequest, packet, 1, *new(types.UInt32))
		},
		"UpdateDetails": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateDetails(errBadRequest, packet, 1, *new(types.UInt32), *new(types.UInt32))
		},
		"GetList": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getList(errBadRequest, packet, 1, *new(types.UInt8), *new(types.Bool))
		},
		"GetDetailedList": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getDetailedList(errBadRequest, packet, 1, *new(types.UInt8), *new(types.Bool))
		},
		"GetRelationships": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getRelationships(errBadRequest, packet, 1, *new(types.ResultRange))
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
