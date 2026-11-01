// Package simpleauthentication provides the common handlers for the Simple Authentication protocol (Rendez-Vous/Quazal, protocol ID 16).
//
// Every method decodes its request and encodes its response as documented on the NintendoClients wiki
// page "Simple Authentication Protocol". The actual behaviour is supplied through the hook fields below; a method
// whose hook is nil answers Core::NotImplemented.
package simpleauthentication

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	simple_authentication "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/simple-authentication"
)

// CommonProtocol handles the Simple Authentication protocol on an endpoint.
// Each hook receives the decoded request and returns the values of the response in wiki order.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol simple_authentication.Interface

	// Authenticate is called for Simple Authentication::Authenticate
	Authenticate func(packet nex.PacketInterface, strUserName types.String) (respRetval types.QResult, respPidPrincipal types.UInt32, respPConnectionData types.RVConnectionData, respStrReturnMsg types.String, nexError *nex.Error)

	// LoginWithToken is called for Simple Authentication::LoginWithToken
	LoginWithToken func(packet nex.PacketInterface, strToken types.String) (respRetval types.QResult, respPidPrincipal types.UInt32, respPConnectionData types.RVConnectionData, respStrReturnMsg types.String, nexError *nex.Error)

	// LoginWithTokenEx is called for Simple Authentication::LoginWithTokenEx
	LoginWithTokenEx func(packet nex.PacketInterface, strToken types.String, oAnyData types.DataHolder) (respRetval types.QResult, respPidPrincipal types.UInt32, respPConnectionData types.RVConnectionData, respStrReturnMsg types.String, nexError *nex.Error)

	// Login is called for Simple Authentication::Login
	Login func(packet nex.PacketInterface, strUsername types.String, strPassword types.String) (respRetval types.QResult, respPidPrincipal types.UInt32, respPConnectionData types.RVConnectionData, respStrReturnMsg types.String, nexError *nex.Error)

	// LoginWithSubAccount is called for Simple Authentication::LoginWithSubAccount
	LoginWithSubAccount func(packet nex.PacketInterface, loginData types.DataHolder) (respRetval types.QResult, respPidPrincipal types.UInt32, respPConnectionData types.RVConnectionData, respStrReturnMsg types.String, nexError *nex.Error)

	// Register is called for Simple Authentication::Register
	Register func(packet nex.PacketInterface, vecMyURLs types.List[types.StationURL]) (respRetval types.QResult, respPidConnectionID types.UInt32, respUrlPublic types.StationURL, nexError *nex.Error)

	// RegisterEx is called for Simple Authentication::RegisterEx
	RegisterEx func(packet nex.PacketInterface, vecMyURLs types.List[types.StationURL], hCustomData types.DataHolder) (respRetval types.QResult, respPidConnectionID types.UInt32, respUrlPublic types.StationURL, nexError *nex.Error)

	// LoginWithTokenCafe is called for Simple Authentication::LoginWithTokenCafe
	LoginWithTokenCafe func(packet nex.PacketInterface, strNintendoToken types.String) (respRetval types.QResult, respPidPrincipal types.UInt32, respPConnectionData types.RVConnectionData, respStrReturnMsg types.String, nexError *nex.Error)

	// LoginWithTokenCafeEx is called for Simple Authentication::LoginWithTokenCafeEx
	LoginWithTokenCafeEx func(packet nex.PacketInterface, strNintendoToken types.String, oAnyData types.DataHolder) (respRetval types.QResult, respPidPrincipal types.UInt32, respPConnectionData types.RVConnectionData, respStrReturnMsg types.String, nexError *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol simple_authentication.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerAuthenticate(commonProtocol.authenticate)
	protocol.SetHandlerLoginWithToken(commonProtocol.loginWithToken)
	protocol.SetHandlerLoginWithTokenEx(commonProtocol.loginWithTokenEx)
	protocol.SetHandlerLogin(commonProtocol.login)
	protocol.SetHandlerLoginWithSubAccount(commonProtocol.loginWithSubAccount)
	protocol.SetHandlerRegister(commonProtocol.register)
	protocol.SetHandlerRegisterEx(commonProtocol.registerEx)
	protocol.SetHandlerLoginWithTokenCafe(commonProtocol.loginWithTokenCafe)
	protocol.SetHandlerLoginWithTokenCafeEx(commonProtocol.loginWithTokenCafeEx)

	return commonProtocol
}
