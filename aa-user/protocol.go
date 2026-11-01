// Package aauser provides the hook driven common handlers for the AA User (0x7B). Every method is hook driven protocol
package aauser

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	aa_user "github.com/PretendoNetwork/nex-protocols-go/v2/aa-user"
	aauser_types "github.com/PretendoNetwork/nex-protocols-go/v2/aa-user/types"
)

// CommonProtocol implements the AAUser protocol on top of hooks. A method whose hook is nil answers Core::NotImplemented
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol aa_user.Interface

	// RegisterApplication is called for AAUser::RegisterApplication
	RegisterApplication func(packet nex.PacketInterface, titleID types.UInt64) *nex.Error
	// UnregisterApplication is called for AAUser::UnregisterApplication
	UnregisterApplication func(packet nex.PacketInterface, titleID types.UInt64) *nex.Error
	// SetApplicationInfo is called for AAUser::SetApplicationInfo
	SetApplicationInfo func(packet nex.PacketInterface, applicationInfo types.List[aauser_types.ApplicationInfo]) *nex.Error
	// GetApplicationInfo is called for AAUser::GetApplicationInfo
	GetApplicationInfo func(packet nex.PacketInterface) (types.List[aauser_types.ApplicationInfo], *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol aa_user.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerRegisterApplication(commonProtocol.registerApplication)
	protocol.SetHandlerUnregisterApplication(commonProtocol.unregisterApplication)
	protocol.SetHandlerSetApplicationInfo(commonProtocol.setApplicationInfo)
	protocol.SetHandlerGetApplicationInfo(commonProtocol.getApplicationInfo)

	return commonProtocol
}
