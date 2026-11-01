package shop_pokemon_bank

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	shop_pokemon_bank "github.com/PretendoNetwork/nex-protocols-go/v2/shop/pokemon-bank"
)

func (commonProtocol *CommonProtocol) debugSetExpireTime(err error, packet nex.PacketInterface, callID uint32, expireTime types.DateTime) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.DebugSetExpireTime == nil {
		common_globals.Logger.Warning("ShopPokemonBank::DebugSetExpireTime missing DebugSetExpireTime hook!")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	errCode := commonProtocol.DebugSetExpireTime(connection.PID(), expireTime)
	if errCode != nil {
		return nil, errCode
	}

	var rmcResponseBody []byte

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = shop_pokemon_bank.ProtocolID
	rmcResponse.MethodID = shop_pokemon_bank.MethodDebugSetExpireTime
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
