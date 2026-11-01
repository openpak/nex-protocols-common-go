package shop_pokemon_bank

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	shop_pokemon_bank "github.com/PretendoNetwork/nex-protocols-go/v2/shop/pokemon-bank"
)

func (commonProtocol *CommonProtocol) supportNumberToPrincipalID(err error, packet nex.PacketInterface, callID uint32, supportNumber types.String) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.SupportNumberToPrincipalID == nil {
		common_globals.Logger.Warning("ShopPokemonBank::SupportNumberToPrincipalID missing SupportNumberToPrincipalID hook!")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	response, errCode := commonProtocol.SupportNumberToPrincipalID(connection.PID(), supportNumber)
	if errCode != nil {
		return nil, errCode
	}

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	response.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = shop_pokemon_bank.ProtocolID
	rmcResponse.MethodID = shop_pokemon_bank.MethodSupportNumberToPrincipalID
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
