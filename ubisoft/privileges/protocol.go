// Package privileges provides the common handlers for the Privileges protocol (Rendez-Vous/Quazal, protocol ID 35).
//
// Every method decodes its request and encodes its response as documented on the NintendoClients wiki
// page "Privileges Protocol". The actual behaviour is supplied through the hook fields below; a method
// whose hook is nil answers Core::NotImplemented.
package privileges

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	privileges "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges"
	privileges_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges/types"
)

// CommonProtocol handles the Privileges protocol on an endpoint.
// Each hook receives the decoded request and returns the values of the response in wiki order.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol privileges.Interface

	// GetPrivileges is called for Privileges::GetPrivileges
	GetPrivileges func(packet nex.PacketInterface, localeCode types.String) (respPrivileges types.Map[types.UInt32, privileges_types.Privilege], nexError *nex.Error)

	// ActivateKey is called for Privileges::ActivateKey
	ActivateKey func(packet nex.PacketInterface, uniqueKey types.String, languageCode types.String) (respPrivilege privileges_types.PrivilegeGroup, nexError *nex.Error)

	// ActivateKeyWithExpectedPrivileges is called for Privileges::ActivateKeyWithExpectedPrivileges
	ActivateKeyWithExpectedPrivileges func(packet nex.PacketInterface, uniqueKey types.String, languageCode types.String, expectedPrivilegeIDs types.List[types.UInt32]) (respPrivilege privileges_types.PrivilegeGroup, nexError *nex.Error)

	// GetPrivilegeRemainDuration is called for Privileges::GetPrivilegeRemainDuration
	GetPrivilegeRemainDuration func(packet nex.PacketInterface, privilegeID types.UInt32) (respSeconds types.Int32, nexError *nex.Error)

	// GetExpiredPrivileges is called for Privileges::GetExpiredPrivileges
	GetExpiredPrivileges func(packet nex.PacketInterface) (respExpiredPrivileges types.List[privileges_types.PrivilegeEx], nexError *nex.Error)

	// GetPrivilegesEx is called for Privileges::GetPrivilegesEx
	GetPrivilegesEx func(packet nex.PacketInterface, localeCode types.String) (respPrivilegesEx types.List[privileges_types.PrivilegeEx], nexError *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol privileges.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetPrivileges(commonProtocol.getPrivileges)
	protocol.SetHandlerActivateKey(commonProtocol.activateKey)
	protocol.SetHandlerActivateKeyWithExpectedPrivileges(commonProtocol.activateKeyWithExpectedPrivileges)
	protocol.SetHandlerGetPrivilegeRemainDuration(commonProtocol.getPrivilegeRemainDuration)
	protocol.SetHandlerGetExpiredPrivileges(commonProtocol.getExpiredPrivileges)
	protocol.SetHandlerGetPrivilegesEx(commonProtocol.getPrivilegesEx)

	return commonProtocol
}
