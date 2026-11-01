package shop_pokemon_bank

import (
	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	shop_pokemon_bank "github.com/PretendoNetwork/nex-protocols-go/v2/shop/pokemon-bank"
)

func (commonProtocol *CommonProtocol) getItems(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetItems == nil {
		common_globals.Logger.Warning("ShopPokemonBank::GetItems missing GetItems hook!")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	response, errCode := commonProtocol.GetItems(connection.PID())
	if errCode != nil {
		return nil, errCode
	}

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	response.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = shop_pokemon_bank.ProtocolID
	rmcResponse.MethodID = shop_pokemon_bank.MethodGetItems
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
