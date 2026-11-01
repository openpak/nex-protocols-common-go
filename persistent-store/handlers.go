package persistentstore

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	persistent_store "github.com/PretendoNetwork/nex-protocols-go/v2/persistent-store"
)

func (commonProtocol *CommonProtocol) findByGroup(err error, packet nex.PacketInterface, callID uint32, uiGroup types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.FindByGroup == nil {
		return nil, rmc.NotImplemented("PersistentStore::FindByGroup")
	}

	result0, nexError := commonProtocol.FindByGroup(packet, uiGroup)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, persistent_store.ProtocolID, persistent_store.MethodFindByGroup, callID, result0), nil
}

func (commonProtocol *CommonProtocol) insertItem(err error, packet nex.PacketInterface, callID uint32, uiGroup types.UInt32, strTag types.String, bufData types.Buffer, bReplace types.Bool) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.InsertItem == nil {
		return nil, rmc.NotImplemented("PersistentStore::InsertItem")
	}

	result0, nexError := commonProtocol.InsertItem(packet, uiGroup, strTag, bufData, bReplace)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, persistent_store.ProtocolID, persistent_store.MethodInsertItem, callID, result0), nil
}

func (commonProtocol *CommonProtocol) removeItem(err error, packet nex.PacketInterface, callID uint32, uiGroup types.UInt32, strTag types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.RemoveItem == nil {
		return nil, rmc.NotImplemented("PersistentStore::RemoveItem")
	}

	result0, nexError := commonProtocol.RemoveItem(packet, uiGroup, strTag)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, persistent_store.ProtocolID, persistent_store.MethodRemoveItem, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getItem(err error, packet nex.PacketInterface, callID uint32, uiGroup types.UInt32, strTag types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetItem == nil {
		return nil, rmc.NotImplemented("PersistentStore::GetItem")
	}

	result0, result1, nexError := commonProtocol.GetItem(packet, uiGroup, strTag)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, persistent_store.ProtocolID, persistent_store.MethodGetItem, callID, result0, result1), nil
}

func (commonProtocol *CommonProtocol) insertCustomItem(err error, packet nex.PacketInterface, callID uint32, uiGroup types.UInt32, strTag types.String, hData types.DataHolder, bReplace types.Bool) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.InsertCustomItem == nil {
		return nil, rmc.NotImplemented("PersistentStore::InsertCustomItem")
	}

	if nexError := commonProtocol.InsertCustomItem(packet, uiGroup, strTag, hData, bReplace); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, persistent_store.ProtocolID, persistent_store.MethodInsertCustomItem, callID), nil
}

func (commonProtocol *CommonProtocol) getCustomItem(err error, packet nex.PacketInterface, callID uint32, uiGroup types.UInt32, strTag types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetCustomItem == nil {
		return nil, rmc.NotImplemented("PersistentStore::GetCustomItem")
	}

	result0, nexError := commonProtocol.GetCustomItem(packet, uiGroup, strTag)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, persistent_store.ProtocolID, persistent_store.MethodGetCustomItem, callID, result0), nil
}

func (commonProtocol *CommonProtocol) findItemsBySQLQuery(err error, packet nex.PacketInterface, callID uint32, uiGroup types.UInt32, strTag types.String, strQuery types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.FindItemsBySQLQuery == nil {
		return nil, rmc.NotImplemented("PersistentStore::FindItemsBySQLQuery")
	}

	result0, nexError := commonProtocol.FindItemsBySQLQuery(packet, uiGroup, strTag, strQuery)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, persistent_store.ProtocolID, persistent_store.MethodFindItemsBySQLQuery, callID, result0), nil
}
