package aauser

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	aa_user "github.com/PretendoNetwork/nex-protocols-go/v2/aa-user"
	aauser_types "github.com/PretendoNetwork/nex-protocols-go/v2/aa-user/types"
)

func (commonProtocol *CommonProtocol) registerApplication(err error, packet nex.PacketInterface, callID uint32, titleID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.RegisterApplication == nil {
		return nil, rmc.NotImplemented("AAUser::RegisterApplication")
	}

	if nexError := commonProtocol.RegisterApplication(packet, titleID); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, aa_user.ProtocolID, aa_user.MethodRegisterApplication, callID), nil
}

func (commonProtocol *CommonProtocol) unregisterApplication(err error, packet nex.PacketInterface, callID uint32, titleID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.UnregisterApplication == nil {
		return nil, rmc.NotImplemented("AAUser::UnregisterApplication")
	}

	if nexError := commonProtocol.UnregisterApplication(packet, titleID); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, aa_user.ProtocolID, aa_user.MethodUnregisterApplication, callID), nil
}

func (commonProtocol *CommonProtocol) setApplicationInfo(err error, packet nex.PacketInterface, callID uint32, applicationInfo types.List[aauser_types.ApplicationInfo]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.SetApplicationInfo == nil {
		return nil, rmc.NotImplemented("AAUser::SetApplicationInfo")
	}

	if nexError := commonProtocol.SetApplicationInfo(packet, applicationInfo); nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, aa_user.ProtocolID, aa_user.MethodSetApplicationInfo, callID), nil
}

func (commonProtocol *CommonProtocol) getApplicationInfo(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.GetApplicationInfo == nil {
		return nil, rmc.NotImplemented("AAUser::GetApplicationInfo")
	}

	result0, nexError := commonProtocol.GetApplicationInfo(packet)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.Respond(packet, aa_user.ProtocolID, aa_user.MethodGetApplicationInfo, callID, result0), nil
}
