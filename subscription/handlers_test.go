package subscription

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/fakenet"
	subscription "github.com/PretendoNetwork/nex-protocols-go/v2/subscription"
	subscription_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscription/types"
)

const callID = 0x11223344

func data(pid types.PID, payload ...byte) subscription_types.SubscriptionData {
	d := subscription_types.NewSubscriptionData()
	d.PrincipalID = pid
	d.Unknown = payload
	return d
}

func TestRegistersEveryMethod(t *testing.T) {
	p := subscription.NewProtocol()
	NewCommonProtocol(p)

	if p.CreateMySubscriptionData == nil || p.UpdateMySubscriptionData == nil || p.ClearMySubscriptionData == nil ||
		p.AddTarget == nil || p.DeleteTarget == nil || p.ClearTarget == nil || p.GetFriendSubscriptionData == nil ||
		p.GetTargetSubscriptionData == nil || p.GetActivePlayerSubscriptionData == nil || p.GetSubscriptionData == nil ||
		p.ReplaceTargetAndGetSubscriptionData == nil || p.SetPrivacyLevel == nil || p.GetPrivacyLevel == nil {
		t.Fatal("a Subscription method was left without a handler")
	}
}

func TestNotImplementedWithoutHooks(t *testing.T) {
	c := &CommonProtocol{}
	p := fakenet.NewPacket(7)
	code := nex.ResultCodes.Core.NotImplemented
	sd := subscription_types.NewSubscriptionData()

	msg, err := c.createMySubscriptionData(nil, p, callID, 0, sd, false)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.updateMySubscriptionData(nil, p, callID, sd)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.clearMySubscriptionData(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.addTarget(nil, p, callID, nil)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.deleteTarget(nil, p, callID, nil)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.clearTarget(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getFriendSubscriptionData(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getTargetSubscriptionData(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getActivePlayerSubscriptionData(nil, p, callID, 1, 2, 3)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getSubscriptionData(nil, p, callID, nil)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.replaceTargetAndGetSubscriptionData(nil, p, callID, nil)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.setPrivacyLevel(nil, p, callID, 1)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getPrivacyLevel(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
}

func TestActionsAnswerEmptyAndPassArguments(t *testing.T) {
	var gotPID types.PID
	var gotTargets types.List[types.PID]
	var gotUnknown1 types.UInt32
	var gotUnknown2 types.Bool

	c := &CommonProtocol{
		AddTarget: func(pid types.PID, targets types.List[types.PID]) error {
			gotPID, gotTargets = pid, targets
			return nil
		},
		CreateMySubscriptionData: func(pid types.PID, u1 types.UInt32, param subscription_types.SubscriptionData, u2 types.Bool) error {
			gotUnknown1, gotUnknown2 = u1, u2
			return nil
		},
	}

	msg, err := c.addTarget(nil, fakenet.NewPacket(9), callID, types.List[types.PID]{4, 5})
	fakenet.RequireBody(t, msg, err, subscription.ProtocolID, subscription.MethodAddTarget, callID, nil)
	if gotPID != 9 || len(gotTargets) != 2 {
		t.Fatalf("pid=%d targets=%v", gotPID, gotTargets)
	}

	msg, err = c.createMySubscriptionData(nil, fakenet.NewPacket(9), callID, 6, data(9), true)
	fakenet.RequireBody(t, msg, err, subscription.ProtocolID, subscription.MethodCreateMySubscriptionData, callID, nil)
	if gotUnknown1 != 6 || !bool(gotUnknown2) {
		t.Fatalf("unknowns %d %v", gotUnknown1, gotUnknown2)
	}
}

func TestSubscriptionListLayout(t *testing.T) {
	c := &CommonProtocol{GetSubscriptionData: func(pid types.PID, pids types.List[types.PID]) (types.List[subscription_types.SubscriptionData], error) {
		return types.List[subscription_types.SubscriptionData]{data(0x01020304, 0xAA, 0xBB)}, nil
	}}

	msg, err := c.getSubscriptionData(nil, fakenet.NewPacket(9), callID, types.List[types.PID]{1})
	fakenet.RequireBody(t, msg, err, subscription.ProtocolID, subscription.MethodGetSubscriptionData, callID, []byte{
		1, 0, 0, 0, // list length
		4, 3, 2, 1, // pid
		2, 0, 0xAA, 0xBB, // qBuffer
	})
}

func TestEmptyListsStillAnswer(t *testing.T) {
	c := &CommonProtocol{GetFriendSubscriptionData: func(types.PID) (types.List[subscription_types.SubscriptionData], error) {
		return nil, nil
	}}

	msg, err := c.getFriendSubscriptionData(nil, fakenet.NewPacket(9), callID)
	fakenet.RequireBody(t, msg, err, subscription.ProtocolID, subscription.MethodGetFriendSubscriptionData, callID, []byte{0, 0, 0, 0})
}

func TestActivePlayerLayout(t *testing.T) {
	c := &CommonProtocol{GetActivePlayerSubscriptionData: func(pid types.PID, a, b, d types.UInt32) (types.List[subscription_types.ActivePlayerSubscriptionData], error) {
		ap := subscription_types.NewActivePlayerSubscriptionData()
		ap.SubscriptionData = data(2, 0x07)
		ap.Unknown = true
		return types.List[subscription_types.ActivePlayerSubscriptionData]{ap}, nil
	}}

	msg, err := c.getActivePlayerSubscriptionData(nil, fakenet.NewPacket(9), callID, 1, 2, 3)
	fakenet.RequireBody(t, msg, err, subscription.ProtocolID, subscription.MethodGetActivePlayerSubscriptionData, callID, []byte{
		1, 0, 0, 0,
		2, 0, 0, 0, 1, 0, 0x07, // subscription data
		1, // unknown bool
	})
}

func TestPrivacyLevelLayout(t *testing.T) {
	c := &CommonProtocol{GetPrivacyLevel: func(types.PID) (types.UInt32, error) { return 1, nil }}

	msg, err := c.getPrivacyLevel(nil, fakenet.NewPacket(9), callID)
	fakenet.RequireBody(t, msg, err, subscription.ProtocolID, subscription.MethodGetPrivacyLevel, callID, []byte{1, 0, 0, 0})
}

func TestErrors(t *testing.T) {
	c := &CommonProtocol{ClearTarget: func(types.PID) error { return errors.New("x") }}
	p := fakenet.NewPacket(9)

	msg, err := c.clearTarget(nil, p, callID)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.Unknown)
	msg, err = c.clearTarget(errors.New("short"), p, callID)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.InvalidArgument)
}
