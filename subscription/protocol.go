// Package subscription implements a hook-driven common handler set for the
// Subscription protocol (0x75). The game server supplies storage through the
// function fields on CommonProtocol; this package decodes requests, calls the
// hooks and encodes the responses.
//
// A method whose hook is nil answers Core::NotImplemented. Hooks that return an
// error answer Core::Unknown unless the error is a *nex.Error, which is sent
// as is. In every hook pid is the caller.
package subscription

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	subscription "github.com/PretendoNetwork/nex-protocols-go/v2/subscription"
	subscription_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscription/types"
)

// CommonProtocol holds the hooks for the Subscription protocol.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol subscription.Interface

	// CreateMySubscriptionData creates the caller's subscription data.
	CreateMySubscriptionData func(pid types.PID, unknown1 types.UInt32, param subscription_types.SubscriptionData, unknown2 types.Bool) error
	// UpdateMySubscriptionData replaces the caller's subscription data.
	UpdateMySubscriptionData func(pid types.PID, param subscription_types.SubscriptionData) error
	// ClearMySubscriptionData removes the caller's subscription data.
	ClearMySubscriptionData func(pid types.PID) error
	// AddTarget adds users to the caller's target list.
	AddTarget func(pid types.PID, targets types.List[types.PID]) error
	// DeleteTarget removes users from the caller's target list.
	DeleteTarget func(pid types.PID, targets types.List[types.PID]) error
	// ClearTarget empties the caller's target list.
	ClearTarget func(pid types.PID) error
	// GetFriendSubscriptionData returns the subscription data of the caller's friends.
	GetFriendSubscriptionData func(pid types.PID) (types.List[subscription_types.SubscriptionData], error)
	// GetTargetSubscriptionData returns the subscription data of the caller's targets.
	GetTargetSubscriptionData func(pid types.PID) (types.List[subscription_types.SubscriptionData], error)
	// GetActivePlayerSubscriptionData returns subscription data of recently active players.
	GetActivePlayerSubscriptionData func(pid types.PID, unknown1 types.UInt32, unknown2 types.UInt32, unknown3 types.UInt32) (types.List[subscription_types.ActivePlayerSubscriptionData], error)
	// GetSubscriptionData returns the subscription data of the given users.
	GetSubscriptionData func(pid types.PID, pids types.List[types.PID]) (types.List[subscription_types.SubscriptionData], error)
	// ReplaceTargetAndGetSubscriptionData replaces the target list and returns the new targets' data.
	ReplaceTargetAndGetSubscriptionData func(pid types.PID, newTargets types.List[types.PID]) (types.List[subscription_types.SubscriptionData], error)
	// SetPrivacyLevel stores the caller's privacy level.
	SetPrivacyLevel func(pid types.PID, privacyLevel types.UInt32) error
	// GetPrivacyLevel returns the caller's privacy level.
	GetPrivacyLevel func(pid types.PID) (types.UInt32, error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol subscription.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerCreateMySubscriptionData(commonProtocol.createMySubscriptionData)
	protocol.SetHandlerUpdateMySubscriptionData(commonProtocol.updateMySubscriptionData)
	protocol.SetHandlerClearMySubscriptionData(commonProtocol.clearMySubscriptionData)
	protocol.SetHandlerAddTarget(commonProtocol.addTarget)
	protocol.SetHandlerDeleteTarget(commonProtocol.deleteTarget)
	protocol.SetHandlerClearTarget(commonProtocol.clearTarget)
	protocol.SetHandlerGetFriendSubscriptionData(commonProtocol.getFriendSubscriptionData)
	protocol.SetHandlerGetTargetSubscriptionData(commonProtocol.getTargetSubscriptionData)
	protocol.SetHandlerGetActivePlayerSubscriptionData(commonProtocol.getActivePlayerSubscriptionData)
	protocol.SetHandlerGetSubscriptionData(commonProtocol.getSubscriptionData)
	protocol.SetHandlerReplaceTargetAndGetSubscriptionData(commonProtocol.replaceTargetAndGetSubscriptionData)
	protocol.SetHandlerSetPrivacyLevel(commonProtocol.setPrivacyLevel)
	protocol.SetHandlerGetPrivacyLevel(commonProtocol.getPrivacyLevel)

	return commonProtocol
}
