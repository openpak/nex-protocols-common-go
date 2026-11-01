package localization

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/localization"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

// encodeGetLocaleCodeResponse encodes the GetLocaleCode response body
func encodeGetLocaleCodeResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetLocaleCodeResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getLocaleCode(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Localization::GetLocaleCode: invalid request")
	}

	if commonProtocol.GetLocaleCode == nil {
		common_globals.Logger.Warning("Localization::GetLocaleCode not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Localization::GetLocaleCode not implemented")
	}

	r0, nexErr := commonProtocol.GetLocaleCode(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetLocaleCodeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetLocaleCodeResponse{LocalCode: r0})

	return respond(endpoint, ubi.MethodGetLocaleCode, callID, body), nil
}

func (commonProtocol *CommonProtocol) setLocaleCode(err error, packet nex.PacketInterface, callID uint32, localCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Localization::SetLocaleCode: invalid request")
	}

	if commonProtocol.SetLocaleCode == nil {
		common_globals.Logger.Warning("Localization::SetLocaleCode not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Localization::SetLocaleCode not implemented")
	}

	nexErr := commonProtocol.SetLocaleCode(packet, localCode)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodSetLocaleCode, callID, nil), nil
}
