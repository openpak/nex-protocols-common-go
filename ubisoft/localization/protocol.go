// Package localization implements the common Localization protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Localization Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
package localization

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/localization"
)

// CommonProtocol handles the Localization protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// GetLocaleCode is required: when nil the request is answered with Core::NotImplemented.
	GetLocaleCode func(packet nex.PacketInterface) (types.String, *nex.Error)

	// SetLocaleCode is required: when nil the request is answered with Core::NotImplemented.
	SetLocaleCode func(packet nex.PacketInterface, localCode types.String) *nex.Error
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetLocaleCode(commonProtocol.getLocaleCode)
	protocol.SetHandlerSetLocaleCode(commonProtocol.setLocaleCode)

	return commonProtocol
}
