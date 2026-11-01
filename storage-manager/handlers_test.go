package storagemanager

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/fakenet"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/v2/storage-manager"
)

const callID = 0x55

func TestRegistersEveryMethod(t *testing.T) {
	p := storage_manager.NewProtocol()
	NewCommonProtocol(p)

	if p.AcquireNexUniqueID == nil || p.NexUniqueIDToPrincipalID == nil || p.Unk3 == nil || p.AcquireCardID == nil || p.ActivateWithCardID == nil {
		t.Fatal("a StorageManager method was left without a handler")
	}
}

func TestNotImplementedWithoutHooks(t *testing.T) {
	c := &CommonProtocol{}
	p := fakenet.NewPacket(3)
	code := nex.ResultCodes.Core.NotImplemented

	msg, err := c.acquireNexUniqueID(nil, p, callID, 0)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.nexUniqueIDToPrincipalID(nil, p, callID, 0)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.unk3(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.acquireCardID(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.activateWithCardID(nil, p, callID, 0, 0)
	fakenet.RequireError(t, msg, err, code)
}

func TestAcquireCardIDLayout(t *testing.T) {
	c := &CommonProtocol{AcquireCardID: func(types.PID) (types.UInt64, error) { return 0x1122334455667788, nil }}

	msg, err := c.acquireCardID(nil, fakenet.NewPacket(3), callID)
	fakenet.RequireBody(t, msg, err, storage_manager.ProtocolID, storage_manager.MethodAcquireCardID, callID,
		[]byte{0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11})
}

func TestActivateWithCardIDLayoutAndArguments(t *testing.T) {
	var gotPID types.PID
	var gotSlot types.UInt8
	var gotCard types.UInt64

	c := &CommonProtocol{ActivateWithCardID: func(pid types.PID, slot types.UInt8, card types.UInt64) (types.UInt32, types.Bool, error) {
		gotPID, gotSlot, gotCard = pid, slot, card
		return 0x01020304, true, nil
	}}

	msg, err := c.activateWithCardID(nil, fakenet.NewPacket(3), callID, 2, 77)
	fakenet.RequireBody(t, msg, err, storage_manager.ProtocolID, storage_manager.MethodActivateWithCardID, callID,
		[]byte{4, 3, 2, 1, 1})

	if gotPID != 3 || gotSlot != 2 || gotCard != 77 {
		t.Fatalf("hook got %d %d %d", gotPID, gotSlot, gotCard)
	}
}

func TestUniqueIDMethods(t *testing.T) {
	c := &CommonProtocol{
		AcquireNexUniqueID:       func(types.PID, types.UInt8) (types.UInt32, error) { return 9, nil },
		NexUniqueIDToPrincipalID: func(types.PID, types.UInt32) (types.PID, error) { return 0x0A0B0C0D, nil },
	}
	p := fakenet.NewPacket(3)

	msg, err := c.acquireNexUniqueID(nil, p, callID, 1)
	fakenet.RequireBody(t, msg, err, storage_manager.ProtocolID, storage_manager.MethodAcquireNexUniqueID, callID, []byte{9, 0, 0, 0})

	msg, err = c.nexUniqueIDToPrincipalID(nil, p, callID, 9)
	fakenet.RequireBody(t, msg, err, storage_manager.ProtocolID, storage_manager.MethodNexUniqueIDToPrincipalID, callID, []byte{0x0D, 0x0C, 0x0B, 0x0A})
}

func TestUnk3PassesRawBody(t *testing.T) {
	c := &CommonProtocol{Unk3: func(types.PID) ([]byte, error) { return []byte{1, 2, 3}, nil }}

	msg, err := c.unk3(nil, fakenet.NewPacket(3), callID)
	fakenet.RequireBody(t, msg, err, storage_manager.ProtocolID, storage_manager.MethodUnk3, callID, []byte{1, 2, 3})
}

func TestErrors(t *testing.T) {
	c := &CommonProtocol{AcquireCardID: func(types.PID) (types.UInt64, error) { return 0, errors.New("x") }}
	p := fakenet.NewPacket(3)

	msg, err := c.acquireCardID(nil, p, callID)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.Unknown)
	msg, err = c.acquireCardID(errors.New("short"), p, callID)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.InvalidArgument)
}
