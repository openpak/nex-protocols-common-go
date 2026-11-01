package accountmanagement

import (
	"bytes"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	account_management "github.com/PretendoNetwork/nex-protocols-go/v2/account-management"
	account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/account-management/types"
)

func TestCreateAccountEncodesResult(t *testing.T) {
	commonProtocol := &CommonProtocol{
		CreateAccount: func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String) (types.QResult, *nex.Error) {
			if strPrincipalName != "name" || uiGroups != 2 {
				t.Errorf("bad params %q %d", strPrincipalName, uiGroups)
			}

			return types.NewQResultError(0x00010001), nil
		},
	}

	response, nexError := commonProtocol.createAccount(nil, rmctest.NewPacket(1), 2, "name", "key", 2, "mail")
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * QResult is a little endian u32 with the error bit (0x80000000) set
	if want := []byte{0x01, 0x00, 0x01, 0x80}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}

	if response.ProtocolID != account_management.ProtocolID || response.MethodID != account_management.MethodCreateAccount {
		t.Fatalf("bad header: %+v", response)
	}
}

func TestNintendoCreateAccountEncodesPIDAndHMAC(t *testing.T) {
	commonProtocol := &CommonProtocol{
		NintendoCreateAccount: func(packet nex.PacketInterface, strPrincipalName types.String, strKey types.String, uiGroups types.UInt32, strEmail types.String, oAuthData types.DataHolder) (types.PID, types.String, *nex.Error) {
			return 0x11223344, "x", nil
		},
	}

	response, nexError := commonProtocol.nintendoCreateAccount(nil, rmctest.NewPacket(1), 1, "", "", 0, "", types.NewDataHolder())
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * 4 byte PID, then String: u16 length (including the NUL) + bytes + NUL
	if want := []byte{0x44, 0x33, 0x22, 0x11, 2, 0, 'x', 0}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestGetLastConnectionStatsEncodesThreeDateTimes(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetLastConnectionStats: func(packet nex.PacketInterface, idPrincipal types.PID) (types.DateTime, types.DateTime, types.DateTime, *nex.Error) {
			return types.NewDateTime(1), types.NewDateTime(2), types.NewDateTime(3), nil
		},
	}

	response, nexError := commonProtocol.getLastConnectionStats(nil, rmctest.NewPacket(1), 1, 5)
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * Three little endian u64
	want := []byte{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestGetAccountDataEncoding(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetAccountData: func(packet nex.PacketInterface) (types.UInt32, account_management_types.AccountData, *nex.Error) {
			data := account_management_types.NewAccountData()
			data.PID = 7
			data.StrName = "me"

			return 1, data, nil
		},
	}

	response, nexError := commonProtocol.getAccountData(nil, rmctest.NewPacket(1), 1)
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * The Uint32 return value comes first, then the AccountData
	if !bytes.HasPrefix(response.Parameters, []byte{1, 0, 0, 0, 7, 0, 0, 0, 3, 0, 'm', 'e', 0}) {
		t.Fatalf("parameters = %x", response.Parameters)
	}
}

func TestDeleteAccountHasNoBody(t *testing.T) {
	var got types.PID

	commonProtocol := &CommonProtocol{
		DeleteAccount: func(packet nex.PacketInterface, idPrincipal types.PID) *nex.Error {
			got = idPrincipal
			return nil
		},
	}

	response, nexError := commonProtocol.deleteAccount(nil, rmctest.NewPacket(1), 1, 42)
	if nexError != nil || got != 42 || len(response.Parameters) != 0 {
		t.Fatalf("got=%d err=%v response=%+v", got, nexError, response)
	}
}

func TestHookErrorIsReturned(t *testing.T) {
	commonProtocol := &CommonProtocol{
		DisconnectAllPrincipals: func(packet nex.PacketInterface) (types.Bool, *nex.Error) {
			return false, nex.NewError(nex.ResultCodes.Core.AccessDenied, "no")
		},
	}

	response, nexError := commonProtocol.disconnectAllPrincipals(nil, rmctest.NewPacket(1), 1)
	if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.AccessDenied) {
		t.Fatalf("expected AccessDenied, got %v", nexError)
	}
}
