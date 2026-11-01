package subscription

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmcutil"
	subscription "github.com/PretendoNetwork/nex-protocols-go/v2/subscription"
	subscription_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscription/types"
)

type subscriptionList = types.List[subscription_types.SubscriptionData]

func notImplemented(method string) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Warningf("Subscription::%s missing hook!", method)
	return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Subscription::"+method+" not implemented")
}

func invalidArgument(err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Error(err.Error())
	return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, err.Error())
}

func hookFailed(method string, err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Errorf("Subscription::%s hook failed: %s", method, err.Error())
	return nil, rmcutil.HookError(err, nex.ResultCodes.Core.Unknown)
}

func reply(packet nex.PacketInterface, methodID uint32, callID uint32, write func(*nex.ByteStreamOut)) *nex.RMCMessage {
	return rmcutil.Success(packet.Sender().Endpoint(), subscription.ProtocolID, methodID, callID, write)
}

// runAction handles the many methods that only call a hook and answer with an empty body.
func runAction(method string, methodID uint32, hookSet bool, err error, packet nex.PacketInterface, callID uint32, call func(pid types.PID) error) (*nex.RMCMessage, *nex.Error) {
	if !hookSet {
		return notImplemented(method)
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := call(packet.Sender().PID()); err != nil {
		return hookFailed(method, err)
	}

	return reply(packet, methodID, callID, nil), nil
}

// runList handles the methods whose response is one list.
func runList[T types.RVType](method string, methodID uint32, hookSet bool, err error, packet nex.PacketInterface, callID uint32, call func(pid types.PID) (types.List[T], error)) (*nex.RMCMessage, *nex.Error) {
	if !hookSet {
		return notImplemented(method)
	}
	if err != nil {
		return invalidArgument(err)
	}
	list, err := call(packet.Sender().PID())
	if err != nil {
		return hookFailed(method, err)
	}

	return reply(packet, methodID, callID, func(s *nex.ByteStreamOut) { list.WriteTo(s) }), nil
}

func (c *CommonProtocol) createMySubscriptionData(err error, packet nex.PacketInterface, callID uint32, unknown1 types.UInt32, param subscription_types.SubscriptionData, unknown2 types.Bool) (*nex.RMCMessage, *nex.Error) {
	return runAction("CreateMySubscriptionData", subscription.MethodCreateMySubscriptionData, c.CreateMySubscriptionData != nil, err, packet, callID, func(pid types.PID) error {
		return c.CreateMySubscriptionData(pid, unknown1, param, unknown2)
	})
}

func (c *CommonProtocol) updateMySubscriptionData(err error, packet nex.PacketInterface, callID uint32, param subscription_types.SubscriptionData) (*nex.RMCMessage, *nex.Error) {
	return runAction("UpdateMySubscriptionData", subscription.MethodUpdateMySubscriptionData, c.UpdateMySubscriptionData != nil, err, packet, callID, func(pid types.PID) error {
		return c.UpdateMySubscriptionData(pid, param)
	})
}

func (c *CommonProtocol) clearMySubscriptionData(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	return runAction("ClearMySubscriptionData", subscription.MethodClearMySubscriptionData, c.ClearMySubscriptionData != nil, err, packet, callID, func(pid types.PID) error {
		return c.ClearMySubscriptionData(pid)
	})
}

func (c *CommonProtocol) addTarget(err error, packet nex.PacketInterface, callID uint32, targets types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	return runAction("AddTarget", subscription.MethodAddTarget, c.AddTarget != nil, err, packet, callID, func(pid types.PID) error {
		return c.AddTarget(pid, targets)
	})
}

func (c *CommonProtocol) deleteTarget(err error, packet nex.PacketInterface, callID uint32, targets types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	return runAction("DeleteTarget", subscription.MethodDeleteTarget, c.DeleteTarget != nil, err, packet, callID, func(pid types.PID) error {
		return c.DeleteTarget(pid, targets)
	})
}

func (c *CommonProtocol) clearTarget(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	return runAction("ClearTarget", subscription.MethodClearTarget, c.ClearTarget != nil, err, packet, callID, func(pid types.PID) error {
		return c.ClearTarget(pid)
	})
}

func (c *CommonProtocol) setPrivacyLevel(err error, packet nex.PacketInterface, callID uint32, privacyLevel types.UInt32) (*nex.RMCMessage, *nex.Error) {
	return runAction("SetPrivacyLevel", subscription.MethodSetPrivacyLevel, c.SetPrivacyLevel != nil, err, packet, callID, func(pid types.PID) error {
		return c.SetPrivacyLevel(pid, privacyLevel)
	})
}

func (c *CommonProtocol) getFriendSubscriptionData(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	return runList("GetFriendSubscriptionData", subscription.MethodGetFriendSubscriptionData, c.GetFriendSubscriptionData != nil, err, packet, callID, func(pid types.PID) (subscriptionList, error) {
		return c.GetFriendSubscriptionData(pid)
	})
}

func (c *CommonProtocol) getTargetSubscriptionData(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	return runList("GetTargetSubscriptionData", subscription.MethodGetTargetSubscriptionData, c.GetTargetSubscriptionData != nil, err, packet, callID, func(pid types.PID) (subscriptionList, error) {
		return c.GetTargetSubscriptionData(pid)
	})
}

func (c *CommonProtocol) getActivePlayerSubscriptionData(err error, packet nex.PacketInterface, callID uint32, unknown1 types.UInt32, unknown2 types.UInt32, unknown3 types.UInt32) (*nex.RMCMessage, *nex.Error) {
	return runList("GetActivePlayerSubscriptionData", subscription.MethodGetActivePlayerSubscriptionData, c.GetActivePlayerSubscriptionData != nil, err, packet, callID, func(pid types.PID) (types.List[subscription_types.ActivePlayerSubscriptionData], error) {
		return c.GetActivePlayerSubscriptionData(pid, unknown1, unknown2, unknown3)
	})
}

func (c *CommonProtocol) getSubscriptionData(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	return runList("GetSubscriptionData", subscription.MethodGetSubscriptionData, c.GetSubscriptionData != nil, err, packet, callID, func(pid types.PID) (subscriptionList, error) {
		return c.GetSubscriptionData(pid, pids)
	})
}

func (c *CommonProtocol) replaceTargetAndGetSubscriptionData(err error, packet nex.PacketInterface, callID uint32, newTargets types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	return runList("ReplaceTargetAndGetSubscriptionData", subscription.MethodReplaceTargetAndGetSubscriptionData, c.ReplaceTargetAndGetSubscriptionData != nil, err, packet, callID, func(pid types.PID) (subscriptionList, error) {
		return c.ReplaceTargetAndGetSubscriptionData(pid, newTargets)
	})
}

func (c *CommonProtocol) getPrivacyLevel(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if c.GetPrivacyLevel == nil {
		return notImplemented("GetPrivacyLevel")
	}
	if err != nil {
		return invalidArgument(err)
	}
	level, err := c.GetPrivacyLevel(packet.Sender().PID())
	if err != nil {
		return hookFailed("GetPrivacyLevel", err)
	}

	return reply(packet, subscription.MethodGetPrivacyLevel, callID, func(s *nex.ByteStreamOut) { level.WriteTo(s) }), nil
}
