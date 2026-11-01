// Package persistentstore provides the hook driven common handlers for the Persistent Store (0x18). The protocol is implemented on Nintendo game servers but never used by games, so every method is hook driven and no storage is provided protocol
package persistentstore

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	persistent_store "github.com/PretendoNetwork/nex-protocols-go/v2/persistent-store"
)

// CommonProtocol implements the PersistentStore protocol on top of hooks. A method whose hook is nil answers Core::NotImplemented
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol persistent_store.Interface

	// FindByGroup is called for PersistentStore::FindByGroup
	FindByGroup func(packet nex.PacketInterface, uiGroup types.UInt32) (types.List[types.String], *nex.Error)
	// InsertItem is called for PersistentStore::InsertItem
	InsertItem func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String, bufData types.Buffer, bReplace types.Bool) (types.Bool, *nex.Error)
	// RemoveItem is called for PersistentStore::RemoveItem
	RemoveItem func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String) (types.Bool, *nex.Error)
	// GetItem is called for PersistentStore::GetItem
	GetItem func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String) (types.Buffer, types.Bool, *nex.Error)
	// InsertCustomItem is called for PersistentStore::InsertCustomItem
	InsertCustomItem func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String, hData types.DataHolder, bReplace types.Bool) *nex.Error
	// GetCustomItem is called for PersistentStore::GetCustomItem
	GetCustomItem func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String) (types.DataHolder, *nex.Error)
	// FindItemsBySQLQuery is called for PersistentStore::FindItemsBySQLQuery
	FindItemsBySQLQuery func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String, strQuery types.String) (types.List[types.DataHolder], *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol persistent_store.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerFindByGroup(commonProtocol.findByGroup)
	protocol.SetHandlerInsertItem(commonProtocol.insertItem)
	protocol.SetHandlerRemoveItem(commonProtocol.removeItem)
	protocol.SetHandlerGetItem(commonProtocol.getItem)
	protocol.SetHandlerInsertCustomItem(commonProtocol.insertCustomItem)
	protocol.SetHandlerGetCustomItem(commonProtocol.getCustomItem)
	protocol.SetHandlerFindItemsBySQLQuery(commonProtocol.findItemsBySQLQuery)

	return commonProtocol
}
