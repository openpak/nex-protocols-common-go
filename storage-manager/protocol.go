// Package storagemanager implements a hook-driven common handler set for the
// StorageManager protocol. The game server supplies the unique id and card id
// bookkeeping through the function fields on CommonProtocol.
//
// A method whose hook is nil answers Core::NotImplemented. Hooks that return an
// error answer Core::Unknown unless the error is a *nex.Error, which is sent
// as is. In every hook pid is the caller.
//
// Response layouts: AcquireNexUniqueID and ActivateWithCardID answer a Uint32
// unique id (ActivateWithCardID adds a Bool that is true the first time the
// user is seen), AcquireCardID answers a Uint64 and NexUniqueIDToPrincipalID
// answers a PID. Method 3 has no known name or layout; its hook returns the raw
// response body.
package storagemanager

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/v2/storage-manager"
)

// CommonProtocol holds the hooks for the StorageManager protocol.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol storage_manager.Interface

	// AcquireNexUniqueID allocates a unique id for the caller in a slot.
	AcquireNexUniqueID func(pid types.PID, slot types.UInt8) (types.UInt32, error)
	// NexUniqueIDToPrincipalID resolves a unique id to its owner.
	NexUniqueIDToPrincipalID func(pid types.PID, uniqueID types.UInt32) (types.PID, error)
	// Unk3 handles method 3 (name unknown) and returns the raw response body.
	Unk3 func(pid types.PID) ([]byte, error)
	// AcquireCardID allocates a card id for the caller.
	AcquireCardID func(pid types.PID) (types.UInt64, error)
	// ActivateWithCardID binds a card id to a slot and returns the unique id and
	// whether this is the first time the user is seen.
	ActivateWithCardID func(pid types.PID, slot types.UInt8, cardID types.UInt64) (uniqueID types.UInt32, firstTime types.Bool, err error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol storage_manager.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerAcquireNexUniqueID(commonProtocol.acquireNexUniqueID)
	protocol.SetHandlerNexUniqueIDToPrincipalID(commonProtocol.nexUniqueIDToPrincipalID)
	protocol.SetHandlerUnk3(commonProtocol.unk3)
	protocol.SetHandlerAcquireCardID(commonProtocol.acquireCardID)
	protocol.SetHandlerActivateWithCardID(commonProtocol.activateWithCardID)

	return commonProtocol
}
