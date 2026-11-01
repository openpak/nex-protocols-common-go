package storagemanager

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmcutil"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/v2/storage-manager"
)

func notImplemented(method string) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Warningf("StorageManager::%s missing hook!", method)
	return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "StorageManager::"+method+" not implemented")
}

func invalidArgument(err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Error(err.Error())
	return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, err.Error())
}

func hookFailed(method string, err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Errorf("StorageManager::%s hook failed: %s", method, err.Error())
	return nil, rmcutil.HookError(err, nex.ResultCodes.Core.Unknown)
}

func reply(packet nex.PacketInterface, methodID uint32, callID uint32, write func(*nex.ByteStreamOut)) *nex.RMCMessage {
	return rmcutil.Success(packet.Sender().Endpoint(), storage_manager.ProtocolID, methodID, callID, write)
}

func (c *CommonProtocol) acquireNexUniqueID(err error, packet nex.PacketInterface, callID uint32, slot types.UInt8) (*nex.RMCMessage, *nex.Error) {
	if c.AcquireNexUniqueID == nil {
		return notImplemented("AcquireNexUniqueID")
	}
	if err != nil {
		return invalidArgument(err)
	}
	uniqueID, err := c.AcquireNexUniqueID(packet.Sender().PID(), slot)
	if err != nil {
		return hookFailed("AcquireNexUniqueID", err)
	}

	return reply(packet, storage_manager.MethodAcquireNexUniqueID, callID, func(s *nex.ByteStreamOut) { uniqueID.WriteTo(s) }), nil
}

func (c *CommonProtocol) nexUniqueIDToPrincipalID(err error, packet nex.PacketInterface, callID uint32, uniqueID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if c.NexUniqueIDToPrincipalID == nil {
		return notImplemented("NexUniqueIDToPrincipalID")
	}
	if err != nil {
		return invalidArgument(err)
	}
	owner, err := c.NexUniqueIDToPrincipalID(packet.Sender().PID(), uniqueID)
	if err != nil {
		return hookFailed("NexUniqueIDToPrincipalID", err)
	}

	return reply(packet, storage_manager.MethodNexUniqueIDToPrincipalID, callID, func(s *nex.ByteStreamOut) { owner.WriteTo(s) }), nil
}

func (c *CommonProtocol) unk3(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if c.Unk3 == nil {
		return notImplemented("Unk3")
	}
	if err != nil {
		return invalidArgument(err)
	}
	body, err := c.Unk3(packet.Sender().PID())
	if err != nil {
		return hookFailed("Unk3", err)
	}

	return reply(packet, storage_manager.MethodUnk3, callID, func(s *nex.ByteStreamOut) { s.Write(body) }), nil
}

func (c *CommonProtocol) acquireCardID(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if c.AcquireCardID == nil {
		return notImplemented("AcquireCardID")
	}
	if err != nil {
		return invalidArgument(err)
	}
	cardID, err := c.AcquireCardID(packet.Sender().PID())
	if err != nil {
		return hookFailed("AcquireCardID", err)
	}

	return reply(packet, storage_manager.MethodAcquireCardID, callID, func(s *nex.ByteStreamOut) { cardID.WriteTo(s) }), nil
}

func (c *CommonProtocol) activateWithCardID(err error, packet nex.PacketInterface, callID uint32, slot types.UInt8, cardID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if c.ActivateWithCardID == nil {
		return notImplemented("ActivateWithCardID")
	}
	if err != nil {
		return invalidArgument(err)
	}
	uniqueID, firstTime, err := c.ActivateWithCardID(packet.Sender().PID(), slot, cardID)
	if err != nil {
		return hookFailed("ActivateWithCardID", err)
	}

	return reply(packet, storage_manager.MethodActivateWithCardID, callID, func(s *nex.ByteStreamOut) {
		uniqueID.WriteTo(s)
		firstTime.WriteTo(s)
	}), nil
}
