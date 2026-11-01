// Package ubiaccountmanagement provides the common handlers for the Ubi Account Management protocol (Rendez-Vous/Quazal, protocol ID 29).
//
// Every method decodes its request and encodes its response as documented on the NintendoClients wiki
// page "Ubi Account Management Protocol". The actual behaviour is supplied through the hook fields below; a method
// whose hook is nil answers Core::NotImplemented.
package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
)

// CommonProtocol handles the Ubi Account Management protocol on an endpoint.
// Each hook receives the decoded request and returns the values of the response in wiki order.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi_account_management.Interface

	// CreateAccount is called for Ubi Account Management::CreateAccount
	CreateAccount func(packet nex.PacketInterface, ubiAccount ubi_account_management_types.UbiAccount) (respUbiAccount ubi_account_management_types.UbiAccount, respFailedReasons types.List[ubi_account_management_types.ValidationFailureReason], nexError *nex.Error)

	// UpdateAccount is called for Ubi Account Management::UpdateAccount
	UpdateAccount func(packet nex.PacketInterface, ubiAccount ubi_account_management_types.UbiAccount) (respUbiAccount ubi_account_management_types.UbiAccount, respFailedReasons types.List[ubi_account_management_types.ValidationFailureReason], nexError *nex.Error)

	// GetAccount is called for Ubi Account Management::GetAccount
	GetAccount func(packet nex.PacketInterface) (respUbiAccount ubi_account_management_types.UbiAccount, respExists types.Bool, nexError *nex.Error)

	// LinkAccount is called for Ubi Account Management::LinkAccount
	LinkAccount func(packet nex.PacketInterface, ubiAccountUsername types.String, ubiAccountPassword types.String) (nexError *nex.Error)

	// GetTOS is called for Ubi Account Management::GetTOS
	GetTOS func(packet nex.PacketInterface, countryCode types.String, languageCode types.String, htmlVersion types.Bool) (respTos ubi_account_management_types.TOS, nexError *nex.Error)

	// ValidateUsername is called for Ubi Account Management::ValidateUsername
	ValidateUsername func(packet nex.PacketInterface, username types.String) (respUsernameValidation ubi_account_management_types.UsernameValidation, nexError *nex.Error)

	// ValidatePassword is called for Ubi Account Management::ValidatePassword
	ValidatePassword func(packet nex.PacketInterface, password types.String, username types.String) (respFailedReasons types.List[ubi_account_management_types.ValidationFailureReason], nexError *nex.Error)

	// ValidateEmail is called for Ubi Account Management::ValidateEmail
	ValidateEmail func(packet nex.PacketInterface, email types.String) (respFailedReasons types.List[ubi_account_management_types.ValidationFailureReason], nexError *nex.Error)

	// GetCountryList is called for Ubi Account Management::GetCountryList
	GetCountryList func(packet nex.PacketInterface, languageCode types.String) (respCountries types.List[ubi_account_management_types.Country], nexError *nex.Error)

	// ForgetPassword is called for Ubi Account Management::ForgetPassword
	ForgetPassword func(packet nex.PacketInterface, usernameOrEmail types.String) (nexError *nex.Error)

	// LookupPrincipalIds is called for Ubi Account Management::LookupPrincipalIds
	LookupPrincipalIds func(packet nex.PacketInterface, ubiAccountIds types.List[types.String]) (respPids types.Map[types.String, types.UInt32], nexError *nex.Error)

	// LookupUbiAccountIDsByPids is called for Ubi Account Management::LookupUbiAccountIDsByPids
	LookupUbiAccountIDsByPids func(packet nex.PacketInterface, pids types.List[types.UInt32]) (respUbiaccountIDs types.Map[types.UInt32, types.String], nexError *nex.Error)

	// LookupUbiAccountIDsByUsernames is called for Ubi Account Management::LookupUbiAccountIDsByUsernames
	LookupUbiAccountIDsByUsernames func(packet nex.PacketInterface, usernames types.List[types.String]) (respUbiAccountIDs types.Map[types.String, types.String], nexError *nex.Error)

	// LookupUsernamesByUbiAccountIDs is called for Ubi Account Management::LookupUsernamesByUbiAccountIDs
	LookupUsernamesByUbiAccountIDs func(packet nex.PacketInterface, ubiAccountIds types.List[types.String]) (respUsernames types.Map[types.String, types.String], nexError *nex.Error)

	// LookupUbiAccountIDsByUsernameSubString is called for Ubi Account Management::LookupUbiAccountIDsByUsernameSubString
	LookupUbiAccountIDsByUsernameSubString func(packet nex.PacketInterface, usernameSubString types.String) (respUbiAccountIDs types.Map[types.String, types.String], nexError *nex.Error)

	// UserHasPlayed is called for Ubi Account Management::UserHasPlayed
	UserHasPlayed func(packet nex.PacketInterface, ubiAccountIDs types.List[types.String]) (respUserPresence types.Map[types.String, types.Bool], nexError *nex.Error)

	// IsUserPlaying is called for Ubi Account Management::IsUserPlaying
	IsUserPlaying func(packet nex.PacketInterface, ubiAccountIDs types.List[types.String]) (respUserPresence types.Map[types.String, types.Bool], nexError *nex.Error)

	// LookupUbiAccountIDsByUsernamesGlobal is called for Ubi Account Management::LookupUbiAccountIDsByUsernamesGlobal
	LookupUbiAccountIDsByUsernamesGlobal func(packet nex.PacketInterface, usernames types.List[types.String]) (respUbiAccountIDs types.Map[types.String, types.String], nexError *nex.Error)

	// LookupUbiAccountIDsByEmailsGlobal is called for Ubi Account Management::LookupUbiAccountIDsByEmailsGlobal
	LookupUbiAccountIDsByEmailsGlobal func(packet nex.PacketInterface, emails types.List[types.String]) (respUbiAccountIDs types.Map[types.String, types.String], nexError *nex.Error)

	// LookupUsernamesByUbiAccountIDsGlobal is called for Ubi Account Management::LookupUsernamesByUbiAccountIDsGlobal
	LookupUsernamesByUbiAccountIDsGlobal func(packet nex.PacketInterface, ubiAccountIds types.List[types.String]) (respUsernames types.Map[types.String, types.String], nexError *nex.Error)

	// GetTOSEx is called for Ubi Account Management::GetTOSEx
	GetTOSEx func(packet nex.PacketInterface, countryCode types.String, languageCode types.String, htmlVersion types.Bool) (respTosex ubi_account_management_types.TOSEx, nexError *nex.Error)

	// HasAcceptedLatestTOS is called for Ubi Account Management::HasAcceptedLatestTOS
	HasAcceptedLatestTOS func(packet nex.PacketInterface) (respHasAccepted types.Bool, respFailedReasons types.List[ubi_account_management_types.ValidationFailureReason], nexError *nex.Error)

	// AcceptLatestTOS is called for Ubi Account Management::AcceptLatestTOS
	AcceptLatestTOS func(packet nex.PacketInterface) (respFailedReasons types.List[ubi_account_management_types.ValidationFailureReason], nexError *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi_account_management.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerCreateAccount(commonProtocol.createAccount)
	protocol.SetHandlerUpdateAccount(commonProtocol.updateAccount)
	protocol.SetHandlerGetAccount(commonProtocol.getAccount)
	protocol.SetHandlerLinkAccount(commonProtocol.linkAccount)
	protocol.SetHandlerGetTOS(commonProtocol.getTOS)
	protocol.SetHandlerValidateUsername(commonProtocol.validateUsername)
	protocol.SetHandlerValidatePassword(commonProtocol.validatePassword)
	protocol.SetHandlerValidateEmail(commonProtocol.validateEmail)
	protocol.SetHandlerGetCountryList(commonProtocol.getCountryList)
	protocol.SetHandlerForgetPassword(commonProtocol.forgetPassword)
	protocol.SetHandlerLookupPrincipalIds(commonProtocol.lookupPrincipalIds)
	protocol.SetHandlerLookupUbiAccountIDsByPids(commonProtocol.lookupUbiAccountIDsByPids)
	protocol.SetHandlerLookupUbiAccountIDsByUsernames(commonProtocol.lookupUbiAccountIDsByUsernames)
	protocol.SetHandlerLookupUsernamesByUbiAccountIDs(commonProtocol.lookupUsernamesByUbiAccountIDs)
	protocol.SetHandlerLookupUbiAccountIDsByUsernameSubString(commonProtocol.lookupUbiAccountIDsByUsernameSubString)
	protocol.SetHandlerUserHasPlayed(commonProtocol.userHasPlayed)
	protocol.SetHandlerIsUserPlaying(commonProtocol.isUserPlaying)
	protocol.SetHandlerLookupUbiAccountIDsByUsernamesGlobal(commonProtocol.lookupUbiAccountIDsByUsernamesGlobal)
	protocol.SetHandlerLookupUbiAccountIDsByEmailsGlobal(commonProtocol.lookupUbiAccountIDsByEmailsGlobal)
	protocol.SetHandlerLookupUsernamesByUbiAccountIDsGlobal(commonProtocol.lookupUsernamesByUbiAccountIDsGlobal)
	protocol.SetHandlerGetTOSEx(commonProtocol.getTOSEx)
	protocol.SetHandlerHasAcceptedLatestTOS(commonProtocol.hasAcceptedLatestTOS)
	protocol.SetHandlerAcceptLatestTOS(commonProtocol.acceptLatestTOS)

	return commonProtocol
}
