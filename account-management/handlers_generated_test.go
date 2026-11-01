package accountmanagement

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
)

// TestHooksAbsentAnswerNotImplemented checks that every method without a hook answers Core::NotImplemented
func TestHooksAbsentAnswerNotImplemented(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"CreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.createAccount(nil, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String))
		},
		"DeleteAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.deleteAccount(nil, packet, 1, *new(types.PID))
		},
		"DisableAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.disableAccount(nil, packet, 1, *new(types.PID), *new(types.DateTime), *new(types.String))
		},
		"ChangePassword": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.changePassword(nil, packet, 1, *new(types.String))
		},
		"TestCapability": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.testCapability(nil, packet, 1, *new(types.UInt32))
		},
		"GetName":        func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getName(nil, packet, 1, *new(types.PID)) },
		"GetAccountData": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getAccountData(nil, packet, 1) },
		"GetPrivateData": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getPrivateData(nil, packet, 1) },
		"GetPublicData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getPublicData(nil, packet, 1, *new(types.PID))
		},
		"GetMultiplePublicData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getMultiplePublicData(nil, packet, 1, *new(types.List[types.PID]))
		},
		"UpdateAccountName": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountName(nil, packet, 1, *new(types.String))
		},
		"UpdateAccountEmail": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountEmail(nil, packet, 1, *new(types.String))
		},
		"UpdateCustomData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateCustomData(nil, packet, 1, *new(types.DataHolder), *new(types.DataHolder))
		},
		"FindByNameRegex": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findByNameRegex(nil, packet, 1, *new(types.UInt32), *new(types.String), *new(types.ResultRange))
		},
		"UpdateAccountExpiryDate": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountExpiryDate(nil, packet, 1, *new(types.PID), *new(types.DateTime), *new(types.String))
		},
		"UpdateAccountEffectiveDate": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountEffectiveDate(nil, packet, 1, *new(types.PID), *new(types.DateTime), *new(types.String))
		},
		"UpdateStatus": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateStatus(nil, packet, 1, *new(types.String))
		},
		"GetStatus": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getStatus(nil, packet, 1, *new(types.PID)) },
		"GetLastConnectionStats": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getLastConnectionStats(nil, packet, 1, *new(types.PID))
		},
		"ResetPassword": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.resetPassword(nil, packet, 1) },
		"CreateAccountWithCustomData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.createAccountWithCustomData(nil, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder), *new(types.DataHolder))
		},
		"RetrieveAccount": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.retrieveAccount(nil, packet, 1) },
		"UpdateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccount(nil, packet, 1, *new(types.String), *new(types.String), *new(types.DataHolder), *new(types.DataHolder))
		},
		"ChangePasswordByGuest": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.changePasswordByGuest(nil, packet, 1, *new(types.String), *new(types.String), *new(types.String))
		},
		"FindByNameLike": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findByNameLike(nil, packet, 1, *new(types.UInt32), *new(types.String), *new(types.ResultRange))
		},
		"CustomCreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.customCreateAccount(nil, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder))
		},
		"NintendoCreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.nintendoCreateAccount(nil, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder))
		},
		"LookupOrCreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.lookupOrCreateAccount(nil, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder))
		},
		"DisconnectPrincipal": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.disconnectPrincipal(nil, packet, 1, *new(types.PID))
		},
		"DisconnectAllPrincipals": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.disconnectAllPrincipals(nil, packet, 1) },
	}

	if len(calls) != 30 {
		t.Fatalf("expected 30 methods, got %d", len(calls))
	}

	for name, call := range calls {
		response, nexError := call()
		if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.NotImplemented) {
			t.Errorf("%s: expected Core::NotImplemented, got %v", name, nexError)
		}
	}
}

// TestDecodeErrorsAreInvalidArgument checks that a request that failed to decode never reaches a hook
func TestDecodeErrorsAreInvalidArgument(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"CreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.createAccount(errBadRequest, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String))
		},
		"DeleteAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.deleteAccount(errBadRequest, packet, 1, *new(types.PID))
		},
		"DisableAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.disableAccount(errBadRequest, packet, 1, *new(types.PID), *new(types.DateTime), *new(types.String))
		},
		"ChangePassword": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.changePassword(errBadRequest, packet, 1, *new(types.String))
		},
		"TestCapability": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.testCapability(errBadRequest, packet, 1, *new(types.UInt32))
		},
		"GetName": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getName(errBadRequest, packet, 1, *new(types.PID))
		},
		"GetAccountData": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getAccountData(errBadRequest, packet, 1) },
		"GetPrivateData": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getPrivateData(errBadRequest, packet, 1) },
		"GetPublicData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getPublicData(errBadRequest, packet, 1, *new(types.PID))
		},
		"GetMultiplePublicData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getMultiplePublicData(errBadRequest, packet, 1, *new(types.List[types.PID]))
		},
		"UpdateAccountName": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountName(errBadRequest, packet, 1, *new(types.String))
		},
		"UpdateAccountEmail": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountEmail(errBadRequest, packet, 1, *new(types.String))
		},
		"UpdateCustomData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateCustomData(errBadRequest, packet, 1, *new(types.DataHolder), *new(types.DataHolder))
		},
		"FindByNameRegex": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findByNameRegex(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String), *new(types.ResultRange))
		},
		"UpdateAccountExpiryDate": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountExpiryDate(errBadRequest, packet, 1, *new(types.PID), *new(types.DateTime), *new(types.String))
		},
		"UpdateAccountEffectiveDate": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccountEffectiveDate(errBadRequest, packet, 1, *new(types.PID), *new(types.DateTime), *new(types.String))
		},
		"UpdateStatus": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateStatus(errBadRequest, packet, 1, *new(types.String))
		},
		"GetStatus": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getStatus(errBadRequest, packet, 1, *new(types.PID))
		},
		"GetLastConnectionStats": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getLastConnectionStats(errBadRequest, packet, 1, *new(types.PID))
		},
		"ResetPassword": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.resetPassword(errBadRequest, packet, 1) },
		"CreateAccountWithCustomData": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.createAccountWithCustomData(errBadRequest, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder), *new(types.DataHolder))
		},
		"RetrieveAccount": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.retrieveAccount(errBadRequest, packet, 1) },
		"UpdateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.updateAccount(errBadRequest, packet, 1, *new(types.String), *new(types.String), *new(types.DataHolder), *new(types.DataHolder))
		},
		"ChangePasswordByGuest": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.changePasswordByGuest(errBadRequest, packet, 1, *new(types.String), *new(types.String), *new(types.String))
		},
		"FindByNameLike": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findByNameLike(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String), *new(types.ResultRange))
		},
		"CustomCreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.customCreateAccount(errBadRequest, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder))
		},
		"NintendoCreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.nintendoCreateAccount(errBadRequest, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder))
		},
		"LookupOrCreateAccount": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.lookupOrCreateAccount(errBadRequest, packet, 1, *new(types.String), *new(types.String), *new(types.UInt32), *new(types.String), *new(types.DataHolder))
		},
		"DisconnectPrincipal": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.disconnectPrincipal(errBadRequest, packet, 1, *new(types.PID))
		},
		"DisconnectAllPrincipals": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.disconnectAllPrincipals(errBadRequest, packet, 1)
		},
	}

	for name, call := range calls {
		response, nexError := call()
		if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
			t.Errorf("%s: expected Core::InvalidArgument, got %v", name, nexError)
		}
	}
}

var errBadRequest = errors.New("bad request")
