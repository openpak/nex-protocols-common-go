package accountmanagement

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	account_management "github.com/PretendoNetwork/nex-protocols-go/v2/account-management"
)

func (commonProtocol *CommonProtocol) createAccount(err error, packet nex.PacketInterface, callID uint32, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.CreateAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::CreateAccount")
	}

	result0, nexError := commonProtocol.CreateAccount(packet, strPrincipalName, strKey, uiGroups, strEmail)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodCreateAccount, callID, result0), nil
}

func (commonProtocol *CommonProtocol) deleteAccount(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.DeleteAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::DeleteAccount")
	}

	if nexError := commonProtocol.DeleteAccount(packet, idPrincipal); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodDeleteAccount, callID), nil
}

func (commonProtocol *CommonProtocol) disableAccount(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID, dtUntil types.DateTime, strMessage types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.DisableAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::DisableAccount")
	}

	result0, nexError := commonProtocol.DisableAccount(packet, idPrincipal, dtUntil, strMessage)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodDisableAccount, callID, result0), nil
}

func (commonProtocol *CommonProtocol) changePassword(err error, packet nex.PacketInterface, callID uint32, strNewKey types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.ChangePassword == nil {
		return nil, rmc.NotImplemented("AccountManagement::ChangePassword")
	}

	result0, nexError := commonProtocol.ChangePassword(packet, strNewKey)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodChangePassword, callID, result0), nil
}

func (commonProtocol *CommonProtocol) testCapability(err error, packet nex.PacketInterface, callID uint32, uiCapability types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.TestCapability == nil {
		return nil, rmc.NotImplemented("AccountManagement::TestCapability")
	}

	result0, nexError := commonProtocol.TestCapability(packet, uiCapability)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodTestCapability, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getName(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetName == nil {
		return nil, rmc.NotImplemented("AccountManagement::GetName")
	}

	result0, nexError := commonProtocol.GetName(packet, idPrincipal)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodGetName, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getAccountData(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetAccountData == nil {
		return nil, rmc.NotImplemented("AccountManagement::GetAccountData")
	}

	result0, result1, nexError := commonProtocol.GetAccountData(packet)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodGetAccountData, callID, result0, result1), nil
}

func (commonProtocol *CommonProtocol) getPrivateData(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetPrivateData == nil {
		return nil, rmc.NotImplemented("AccountManagement::GetPrivateData")
	}

	result0, result1, nexError := commonProtocol.GetPrivateData(packet)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodGetPrivateData, callID, result0, result1), nil
}

func (commonProtocol *CommonProtocol) getPublicData(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetPublicData == nil {
		return nil, rmc.NotImplemented("AccountManagement::GetPublicData")
	}

	result0, result1, nexError := commonProtocol.GetPublicData(packet, idPrincipal)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodGetPublicData, callID, result0, result1), nil
}

func (commonProtocol *CommonProtocol) getMultiplePublicData(err error, packet nex.PacketInterface, callID uint32, lstPrincipals types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetMultiplePublicData == nil {
		return nil, rmc.NotImplemented("AccountManagement::GetMultiplePublicData")
	}

	result0, result1, nexError := commonProtocol.GetMultiplePublicData(packet, lstPrincipals)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodGetMultiplePublicData, callID, result0, result1), nil
}

func (commonProtocol *CommonProtocol) updateAccountName(err error, packet nex.PacketInterface, callID uint32, strName types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateAccountName == nil {
		return nil, rmc.NotImplemented("AccountManagement::UpdateAccountName")
	}

	result0, nexError := commonProtocol.UpdateAccountName(packet, strName)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodUpdateAccountName, callID, result0), nil
}

func (commonProtocol *CommonProtocol) updateAccountEmail(err error, packet nex.PacketInterface, callID uint32, strName types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateAccountEmail == nil {
		return nil, rmc.NotImplemented("AccountManagement::UpdateAccountEmail")
	}

	result0, nexError := commonProtocol.UpdateAccountEmail(packet, strName)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodUpdateAccountEmail, callID, result0), nil
}

func (commonProtocol *CommonProtocol) updateCustomData(err error, packet nex.PacketInterface, callID uint32, oPublicData types.DataHolder, oPrivateData types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateCustomData == nil {
		return nil, rmc.NotImplemented("AccountManagement::UpdateCustomData")
	}

	result0, nexError := commonProtocol.UpdateCustomData(packet, oPublicData, oPrivateData)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodUpdateCustomData, callID, result0), nil
}

func (commonProtocol *CommonProtocol) findByNameRegex(err error, packet nex.PacketInterface, callID uint32, uiGroups types.UInt32, strRegex types.String, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.FindByNameRegex == nil {
		return nil, rmc.NotImplemented("AccountManagement::FindByNameRegex")
	}

	result0, nexError := commonProtocol.FindByNameRegex(packet, uiGroups, strRegex, resultRange)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodFindByNameRegex, callID, result0), nil
}

func (commonProtocol *CommonProtocol) updateAccountExpiryDate(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID, dtExpiry types.DateTime, strExpiredMessage types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateAccountExpiryDate == nil {
		return nil, rmc.NotImplemented("AccountManagement::UpdateAccountExpiryDate")
	}

	if nexError := commonProtocol.UpdateAccountExpiryDate(packet, idPrincipal, dtExpiry, strExpiredMessage); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodUpdateAccountExpiryDate, callID), nil
}

func (commonProtocol *CommonProtocol) updateAccountEffectiveDate(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID, dtEffectiveFrom types.DateTime, strNotEffectiveMessage types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateAccountEffectiveDate == nil {
		return nil, rmc.NotImplemented("AccountManagement::UpdateAccountEffectiveDate")
	}

	if nexError := commonProtocol.UpdateAccountEffectiveDate(packet, idPrincipal, dtEffectiveFrom, strNotEffectiveMessage); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodUpdateAccountEffectiveDate, callID), nil
}

func (commonProtocol *CommonProtocol) updateStatus(err error, packet nex.PacketInterface, callID uint32, strStatus types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateStatus == nil {
		return nil, rmc.NotImplemented("AccountManagement::UpdateStatus")
	}

	if nexError := commonProtocol.UpdateStatus(packet, strStatus); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodUpdateStatus, callID), nil
}

func (commonProtocol *CommonProtocol) getStatus(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetStatus == nil {
		return nil, rmc.NotImplemented("AccountManagement::GetStatus")
	}

	result0, nexError := commonProtocol.GetStatus(packet, idPrincipal)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodGetStatus, callID, result0), nil
}

func (commonProtocol *CommonProtocol) getLastConnectionStats(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetLastConnectionStats == nil {
		return nil, rmc.NotImplemented("AccountManagement::GetLastConnectionStats")
	}

	result0, result1, result2, nexError := commonProtocol.GetLastConnectionStats(packet, idPrincipal)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodGetLastConnectionStats, callID, result0, result1, result2), nil
}

func (commonProtocol *CommonProtocol) resetPassword(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.ResetPassword == nil {
		return nil, rmc.NotImplemented("AccountManagement::ResetPassword")
	}

	result0, nexError := commonProtocol.ResetPassword(packet)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodResetPassword, callID, result0), nil
}

func (commonProtocol *CommonProtocol) createAccountWithCustomData(err error, packet nex.PacketInterface, callID uint32, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oPublicData types.DataHolder, oPrivateData types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.CreateAccountWithCustomData == nil {
		return nil, rmc.NotImplemented("AccountManagement::CreateAccountWithCustomData")
	}

	if nexError := commonProtocol.CreateAccountWithCustomData(packet, strPrincipalName, strKey, uiGroups, strEmail, oPublicData, oPrivateData); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodCreateAccountWithCustomData, callID), nil
}

func (commonProtocol *CommonProtocol) retrieveAccount(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.RetrieveAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::RetrieveAccount")
	}

	result0, result1, result2, nexError := commonProtocol.RetrieveAccount(packet)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodRetrieveAccount, callID, result0, result1, result2), nil
}

func (commonProtocol *CommonProtocol) updateAccount(err error, packet nex.PacketInterface, callID uint32, strKey types.String, strEmail types.String, oPublicData types.DataHolder, oPrivateData types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UpdateAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::UpdateAccount")
	}

	if nexError := commonProtocol.UpdateAccount(packet, strKey, strEmail, oPublicData, oPrivateData); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodUpdateAccount, callID), nil
}

func (commonProtocol *CommonProtocol) changePasswordByGuest(err error, packet nex.PacketInterface, callID uint32, strPrincipalName types.String, strKey types.String, strEmail types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.ChangePasswordByGuest == nil {
		return nil, rmc.NotImplemented("AccountManagement::ChangePasswordByGuest")
	}

	if nexError := commonProtocol.ChangePasswordByGuest(packet, strPrincipalName, strKey, strEmail); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodChangePasswordByGuest, callID), nil
}

func (commonProtocol *CommonProtocol) findByNameLike(err error, packet nex.PacketInterface, callID uint32, uiGroups types.UInt32, strLike types.String, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.FindByNameLike == nil {
		return nil, rmc.NotImplemented("AccountManagement::FindByNameLike")
	}

	result0, nexError := commonProtocol.FindByNameLike(packet, uiGroups, strLike, resultRange)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodFindByNameLike, callID, result0), nil
}

func (commonProtocol *CommonProtocol) customCreateAccount(err error, packet nex.PacketInterface, callID uint32, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oAuthData types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.CustomCreateAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::CustomCreateAccount")
	}

	result0, nexError := commonProtocol.CustomCreateAccount(packet, strPrincipalName, strKey, uiGroups, strEmail, oAuthData)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodCustomCreateAccount, callID, result0), nil
}

func (commonProtocol *CommonProtocol) nintendoCreateAccount(err error, packet nex.PacketInterface, callID uint32, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oAuthData types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.NintendoCreateAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::NintendoCreateAccount")
	}

	result0, result1, nexError := commonProtocol.NintendoCreateAccount(packet, strPrincipalName, strKey, uiGroups, strEmail, oAuthData)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodNintendoCreateAccount, callID, result0, result1), nil
}

func (commonProtocol *CommonProtocol) lookupOrCreateAccount(err error, packet nex.PacketInterface, callID uint32, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oAuthData types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.LookupOrCreateAccount == nil {
		return nil, rmc.NotImplemented("AccountManagement::LookupOrCreateAccount")
	}

	result0, nexError := commonProtocol.LookupOrCreateAccount(packet, strPrincipalName, strKey, uiGroups, strEmail, oAuthData)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodLookupOrCreateAccount, callID, result0), nil
}

func (commonProtocol *CommonProtocol) disconnectPrincipal(err error, packet nex.PacketInterface, callID uint32, idPrincipal types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.DisconnectPrincipal == nil {
		return nil, rmc.NotImplemented("AccountManagement::DisconnectPrincipal")
	}

	result0, nexError := commonProtocol.DisconnectPrincipal(packet, idPrincipal)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodDisconnectPrincipal, callID, result0), nil
}

func (commonProtocol *CommonProtocol) disconnectAllPrincipals(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.DisconnectAllPrincipals == nil {
		return nil, rmc.NotImplemented("AccountManagement::DisconnectAllPrincipals")
	}

	result0, nexError := commonProtocol.DisconnectAllPrincipals(packet)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, account_management.ProtocolID, account_management.MethodDisconnectAllPrincipals, callID, result0), nil
}
