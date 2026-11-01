// Package accountmanagement provides the hook driven common handlers for the Account Management (0x19). Every method is hook driven, the common library has no account store protocol
package accountmanagement

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	account_management "github.com/PretendoNetwork/nex-protocols-go/v2/account-management"
	account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/account-management/types"
)

// CommonProtocol implements the AccountManagement protocol on top of hooks. A method whose hook is nil answers Core::NotImplemented
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol account_management.Interface

	// CreateAccount is called for AccountManagement::CreateAccount
	CreateAccount func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String) (types.QResult, *nex.Error)
	// DeleteAccount is called for AccountManagement::DeleteAccount
	DeleteAccount func(packet nex.PacketInterface, idPrincipal types.PID) *nex.Error
	// DisableAccount is called for AccountManagement::DisableAccount
	DisableAccount func(packet nex.PacketInterface, idPrincipal types.PID, dtUntil types.DateTime, strMessage types.String) (types.QResult, *nex.Error)
	// ChangePassword is called for AccountManagement::ChangePassword
	ChangePassword func(packet nex.PacketInterface, strNewKey types.String) (types.Bool, *nex.Error)
	// TestCapability is called for AccountManagement::TestCapability
	TestCapability func(packet nex.PacketInterface, uiCapability types.UInt32) (types.Bool, *nex.Error)
	// GetName is called for AccountManagement::GetName
	GetName func(packet nex.PacketInterface, idPrincipal types.PID) (types.String, *nex.Error)
	// GetAccountData is called for AccountManagement::GetAccountData
	GetAccountData func(packet nex.PacketInterface) (types.UInt32, account_management_types.AccountData, *nex.Error)
	// GetPrivateData is called for AccountManagement::GetPrivateData
	GetPrivateData func(packet nex.PacketInterface) (types.Bool, types.DataHolder, *nex.Error)
	// GetPublicData is called for AccountManagement::GetPublicData
	GetPublicData func(packet nex.PacketInterface, idPrincipal types.PID) (types.Bool, types.DataHolder, *nex.Error)
	// GetMultiplePublicData is called for AccountManagement::GetMultiplePublicData
	GetMultiplePublicData func(packet nex.PacketInterface, lstPrincipals types.List[types.PID]) (types.Bool, types.List[types.DataHolder], *nex.Error)
	// UpdateAccountName is called for AccountManagement::UpdateAccountName
	UpdateAccountName func(packet nex.PacketInterface, strName types.String) (types.QResult, *nex.Error)
	// UpdateAccountEmail is called for AccountManagement::UpdateAccountEmail
	UpdateAccountEmail func(packet nex.PacketInterface, strName types.String) (types.QResult, *nex.Error)
	// UpdateCustomData is called for AccountManagement::UpdateCustomData
	UpdateCustomData func(packet nex.PacketInterface, oPublicData types.DataHolder, oPrivateData types.DataHolder) (types.QResult, *nex.Error)
	// FindByNameRegex is called for AccountManagement::FindByNameRegex
	FindByNameRegex func(packet nex.PacketInterface, uiGroups types.UInt32, strRegex types.String, resultRange types.ResultRange) (types.List[account_management_types.BasicAccountInfo], *nex.Error)
	// UpdateAccountExpiryDate is called for AccountManagement::UpdateAccountExpiryDate
	UpdateAccountExpiryDate func(packet nex.PacketInterface, idPrincipal types.PID, dtExpiry types.DateTime, strExpiredMessage types.String) *nex.Error
	// UpdateAccountEffectiveDate is called for AccountManagement::UpdateAccountEffectiveDate
	UpdateAccountEffectiveDate func(packet nex.PacketInterface, idPrincipal types.PID, dtEffectiveFrom types.DateTime, strNotEffectiveMessage types.String) *nex.Error
	// UpdateStatus is called for AccountManagement::UpdateStatus
	UpdateStatus func(packet nex.PacketInterface, strStatus types.String) *nex.Error
	// GetStatus is called for AccountManagement::GetStatus
	GetStatus func(packet nex.PacketInterface, idPrincipal types.PID) (types.String, *nex.Error)
	// GetLastConnectionStats is called for AccountManagement::GetLastConnectionStats
	GetLastConnectionStats func(packet nex.PacketInterface, idPrincipal types.PID) (types.DateTime, types.DateTime, types.DateTime, *nex.Error)
	// ResetPassword is called for AccountManagement::ResetPassword
	ResetPassword func(packet nex.PacketInterface) (types.Bool, *nex.Error)
	// CreateAccountWithCustomData is called for AccountManagement::CreateAccountWithCustomData
	CreateAccountWithCustomData func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oPublicData types.DataHolder, oPrivateData types.DataHolder) *nex.Error
	// RetrieveAccount is called for AccountManagement::RetrieveAccount
	RetrieveAccount func(packet nex.PacketInterface) (account_management_types.AccountData, types.DataHolder, types.DataHolder, *nex.Error)
	// UpdateAccount is called for AccountManagement::UpdateAccount
	UpdateAccount func(packet nex.PacketInterface, strKey types.String, strEmail types.String, oPublicData types.DataHolder, oPrivateData types.DataHolder) *nex.Error
	// ChangePasswordByGuest is called for AccountManagement::ChangePasswordByGuest
	ChangePasswordByGuest func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, strEmail types.String) *nex.Error
	// FindByNameLike is called for AccountManagement::FindByNameLike
	FindByNameLike func(packet nex.PacketInterface, uiGroups types.UInt32, strLike types.String, resultRange types.ResultRange) (types.List[account_management_types.BasicAccountInfo], *nex.Error)
	// CustomCreateAccount is called for AccountManagement::CustomCreateAccount
	CustomCreateAccount func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oAuthData types.DataHolder) (types.PID, *nex.Error)
	// NintendoCreateAccount is called for AccountManagement::NintendoCreateAccount
	NintendoCreateAccount func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oAuthData types.DataHolder) (types.PID, types.String, *nex.Error)
	// LookupOrCreateAccount is called for AccountManagement::LookupOrCreateAccount
	LookupOrCreateAccount func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oAuthData types.DataHolder) (types.PID, *nex.Error)
	// DisconnectPrincipal is called for AccountManagement::DisconnectPrincipal
	DisconnectPrincipal func(packet nex.PacketInterface, idPrincipal types.PID) (types.Bool, *nex.Error)
	// DisconnectAllPrincipals is called for AccountManagement::DisconnectAllPrincipals
	DisconnectAllPrincipals func(packet nex.PacketInterface) (types.Bool, *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol account_management.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerCreateAccount(commonProtocol.createAccount)
	protocol.SetHandlerDeleteAccount(commonProtocol.deleteAccount)
	protocol.SetHandlerDisableAccount(commonProtocol.disableAccount)
	protocol.SetHandlerChangePassword(commonProtocol.changePassword)
	protocol.SetHandlerTestCapability(commonProtocol.testCapability)
	protocol.SetHandlerGetName(commonProtocol.getName)
	protocol.SetHandlerGetAccountData(commonProtocol.getAccountData)
	protocol.SetHandlerGetPrivateData(commonProtocol.getPrivateData)
	protocol.SetHandlerGetPublicData(commonProtocol.getPublicData)
	protocol.SetHandlerGetMultiplePublicData(commonProtocol.getMultiplePublicData)
	protocol.SetHandlerUpdateAccountName(commonProtocol.updateAccountName)
	protocol.SetHandlerUpdateAccountEmail(commonProtocol.updateAccountEmail)
	protocol.SetHandlerUpdateCustomData(commonProtocol.updateCustomData)
	protocol.SetHandlerFindByNameRegex(commonProtocol.findByNameRegex)
	protocol.SetHandlerUpdateAccountExpiryDate(commonProtocol.updateAccountExpiryDate)
	protocol.SetHandlerUpdateAccountEffectiveDate(commonProtocol.updateAccountEffectiveDate)
	protocol.SetHandlerUpdateStatus(commonProtocol.updateStatus)
	protocol.SetHandlerGetStatus(commonProtocol.getStatus)
	protocol.SetHandlerGetLastConnectionStats(commonProtocol.getLastConnectionStats)
	protocol.SetHandlerResetPassword(commonProtocol.resetPassword)
	protocol.SetHandlerCreateAccountWithCustomData(commonProtocol.createAccountWithCustomData)
	protocol.SetHandlerRetrieveAccount(commonProtocol.retrieveAccount)
	protocol.SetHandlerUpdateAccount(commonProtocol.updateAccount)
	protocol.SetHandlerChangePasswordByGuest(commonProtocol.changePasswordByGuest)
	protocol.SetHandlerFindByNameLike(commonProtocol.findByNameLike)
	protocol.SetHandlerCustomCreateAccount(commonProtocol.customCreateAccount)
	protocol.SetHandlerNintendoCreateAccount(commonProtocol.nintendoCreateAccount)
	protocol.SetHandlerLookupOrCreateAccount(commonProtocol.lookupOrCreateAccount)
	protocol.SetHandlerDisconnectPrincipal(commonProtocol.disconnectPrincipal)
	protocol.SetHandlerDisconnectAllPrincipals(commonProtocol.disconnectAllPrincipals)

	return commonProtocol
}
