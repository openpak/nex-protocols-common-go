package playerstats

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/player-stats"
	player_stats_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/player-stats/types"
)

type fakeEndpoint struct{ nex.EndpointInterface }

func (fakeEndpoint) LibraryVersions() *nex.LibraryVersions {
	versions := nex.NewLibraryVersions()
	versions.Main = nex.NewLibraryVersion(3, 0, 0)

	return versions
}

func (fakeEndpoint) ByteStreamSettings() *nex.ByteStreamSettings {
	return nex.NewByteStreamSettings()
}

type fakeConnection struct{ nex.ConnectionInterface }

func (fakeConnection) Endpoint() nex.EndpointInterface { return fakeEndpoint{} }

type fakePacket struct{ nex.PacketInterface }

func (fakePacket) Sender() nex.ConnectionInterface { return fakeConnection{} }

func checkMessage(t *testing.T, message *nex.RMCMessage, nexErr *nex.Error, methodID uint32, callID uint32, want []byte) {
	t.Helper()

	if nexErr != nil {
		t.Fatalf("unexpected error: %s", nexErr.Error())
	}

	if message == nil {
		t.Fatal("nil message")
	}

	if message.ProtocolID != ubi.ProtocolID || message.MethodID != methodID || message.CallID != callID || !message.IsSuccess {
		t.Fatalf("bad message header: protocol %d method %d call %d success %v", message.ProtocolID, message.MethodID, message.CallID, message.IsSuccess)
	}

	if !bytes.Equal(message.Parameters, want) {
		t.Fatalf("response bytes differ\n got %x\nwant %x", message.Parameters, want)
	}
}

func checkNotImplemented(t *testing.T, message *nex.RMCMessage, nexErr *nex.Error) {
	t.Helper()

	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "").ResultCode {
		t.Fatalf("want Core::NotImplemented, got message %v error %v", message, nexErr)
	}
}

func TestWriteStats_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.writeStats(nil, fakePacket{}, 7, types.List[player_stats_types.PlayerStatUpdate]{func() player_stats_types.PlayerStatUpdate {
		v := player_stats_types.NewPlayerStatUpdate()
		v.BoardID = types.NewUInt32(3032)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}
		v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(9086)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(12113)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v13")}
			return v
		}()}
		return v
	}(), func() player_stats_types.PlayerStatUpdate {
		v := player_stats_types.NewPlayerStatUpdate()
		v.BoardID = types.NewUInt32(15140)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(21194)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-27)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(24221)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v25")}
			return v
		}()}
		return v
	}()})
	checkNotImplemented(t, message, nexErr)
}

func TestWriteStats_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.WriteStats = func(packet nex.PacketInterface, _ types.List[player_stats_types.PlayerStatUpdate]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.writeStats(errors.New("bad"), fakePacket{}, 7, types.NewList[player_stats_types.PlayerStatUpdate]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestWriteStats_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.WriteStats = func(packet nex.PacketInterface, got0 types.List[player_stats_types.PlayerStatUpdate]) *nex.Error {
		called = true

		if !(got0.Equals(types.List[player_stats_types.PlayerStatUpdate]{func() player_stats_types.PlayerStatUpdate {
			v := player_stats_types.NewPlayerStatUpdate()
			v.BoardID = types.NewUInt32(3032)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}
			v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(9086)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(12113)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v13")}
				return v
			}()}
			return v
		}(), func() player_stats_types.PlayerStatUpdate {
			v := player_stats_types.NewPlayerStatUpdate()
			v.BoardID = types.NewUInt32(15140)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
			v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(21194)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-27)}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(24221)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v25")}
				return v
			}()}
			return v
		}()})) {
			t.Errorf("playerStatUpdates differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.writeStats(nil, fakePacket{}, 9, types.List[player_stats_types.PlayerStatUpdate]{func() player_stats_types.PlayerStatUpdate {
		v := player_stats_types.NewPlayerStatUpdate()
		v.BoardID = types.NewUInt32(3032)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}
		v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(9086)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(12113)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v13")}
			return v
		}()}
		return v
	}(), func() player_stats_types.PlayerStatUpdate {
		v := player_stats_types.NewPlayerStatUpdate()
		v.BoardID = types.NewUInt32(15140)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(21194)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-27)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(24221)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v25")}
			return v
		}()}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodWriteStats, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.WriteStats = func(packet nex.PacketInterface, _ types.List[player_stats_types.PlayerStatUpdate]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.writeStats(nil, fakePacket{}, 9, types.List[player_stats_types.PlayerStatUpdate]{func() player_stats_types.PlayerStatUpdate {
		v := player_stats_types.NewPlayerStatUpdate()
		v.BoardID = types.NewUInt32(3032)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}
		v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(9086)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(12113)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v13")}
			return v
		}()}
		return v
	}(), func() player_stats_types.PlayerStatUpdate {
		v := player_stats_types.NewPlayerStatUpdate()
		v.BoardID = types.NewUInt32(15140)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(21194)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-27)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(24221)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v25")}
			return v
		}()}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadStatsByPlayers_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readStatsByPlayers(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, types.List[player_stats_types.StatboardQuery]{func() player_stats_types.StatboardQuery {
		v := player_stats_types.NewStatboardQuery()
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.ResetFrequency = types.NewUInt32(10095)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(16149)
			v.SortOrder = types.NewUInt32(17158)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(19176)
			v.SortOrder = types.NewUInt32(20185)
			return v
		}()}
		return v
	}(), func() player_stats_types.StatboardQuery {
		v := player_stats_types.NewStatboardQuery()
		v.BoardID = types.NewUInt32(22203)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(24221), types.NewUInt32(25230)}
		v.ResetFrequency = types.NewUInt32(26239)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(28257), types.NewUInt32(29266)}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(32293)
			v.SortOrder = types.NewUInt32(33302)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(35320)
			v.SortOrder = types.NewUInt32(36329)
			return v
		}()}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadStatsByPlayers, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadStatsByPlayers_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadStatsByPlayers = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ types.List[player_stats_types.StatboardQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.StatboardResult](), nil
	}

	message, nexErr := commonProtocol.readStatsByPlayers(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32](), types.NewList[player_stats_types.StatboardQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadStatsByPlayers_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadStatsByPlayers = func(packet nex.PacketInterface, got0 types.List[types.UInt32], got1 types.List[player_stats_types.StatboardQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("playerPIDs differs")
		}

		if !(got1.Equals(types.List[player_stats_types.StatboardQuery]{func() player_stats_types.StatboardQuery {
			v := player_stats_types.NewStatboardQuery()
			v.BoardID = types.NewUInt32(6059)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
			v.ResetFrequency = types.NewUInt32(10095)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
			v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(16149)
				v.SortOrder = types.NewUInt32(17158)
				return v
			}(), func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(19176)
				v.SortOrder = types.NewUInt32(20185)
				return v
			}()}
			return v
		}(), func() player_stats_types.StatboardQuery {
			v := player_stats_types.NewStatboardQuery()
			v.BoardID = types.NewUInt32(22203)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(24221), types.NewUInt32(25230)}
			v.ResetFrequency = types.NewUInt32(26239)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(28257), types.NewUInt32(29266)}
			v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(32293)
				v.SortOrder = types.NewUInt32(33302)
				return v
			}(), func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(35320)
				v.SortOrder = types.NewUInt32(36329)
				return v
			}()}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.StatboardResult]{func() player_stats_types.StatboardResult {
			v := player_stats_types.NewStatboardResult()
			v.BoardID = types.NewUInt32(39356)
			v.ContextID = types.NewUInt32(40365)
			v.ResetFrequency = types.NewUInt32(41374)
			v.PlayerStatSets = types.List[player_stats_types.PlayerStatSet]{func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(44401)
				v.PlayerName = types.NewString("s45")
				v.SubmittedTime = types.NewDateTime(46000215)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(49446)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-55)}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(52473)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v53")}
					return v
				}()}
				return v
			}(), func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(55500)
				v.PlayerName = types.NewString("s56")
				v.SubmittedTime = types.NewDateTime(57000248)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(60545)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v61")}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(63572)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-69)}
					return v
				}()}
				return v
			}()}
			v.DefaultStatValues = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(67608)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-73)}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(70635)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v71")}
				return v
			}()}
			return v
		}(), func() player_stats_types.StatboardResult {
			v := player_stats_types.NewStatboardResult()
			v.BoardID = types.NewUInt32(73662)
			v.ContextID = types.NewUInt32(74671)
			v.ResetFrequency = types.NewUInt32(75680)
			v.PlayerStatSets = types.List[player_stats_types.PlayerStatSet]{func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(78707)
				v.PlayerName = types.NewString("s79")
				v.SubmittedTime = types.NewDateTime(80000317)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(83752)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-89)}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(86779)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v87")}
					return v
				}()}
				return v
			}(), func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(89806)
				v.PlayerName = types.NewString("s90")
				v.SubmittedTime = types.NewDateTime(91000350)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(94851)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v95")}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(97878)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-103)}
					return v
				}()}
				return v
			}()}
			v.DefaultStatValues = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(101914)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-107)}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(104941)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v105")}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readStatsByPlayers(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, types.List[player_stats_types.StatboardQuery]{func() player_stats_types.StatboardQuery {
		v := player_stats_types.NewStatboardQuery()
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.ResetFrequency = types.NewUInt32(10095)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(16149)
			v.SortOrder = types.NewUInt32(17158)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(19176)
			v.SortOrder = types.NewUInt32(20185)
			return v
		}()}
		return v
	}(), func() player_stats_types.StatboardQuery {
		v := player_stats_types.NewStatboardQuery()
		v.BoardID = types.NewUInt32(22203)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(24221), types.NewUInt32(25230)}
		v.ResetFrequency = types.NewUInt32(26239)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(28257), types.NewUInt32(29266)}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(32293)
			v.SortOrder = types.NewUInt32(33302)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(35320)
			v.SortOrder = types.NewUInt32(36329)
			return v
		}()}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadStatsByPlayers, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xbc, 0x99, 0x00, 0x00, 0xad, 0x9d, 0x00, 0x00, 0x9e, 0xa1, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x71, 0xad, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x35, 0x00, 0x57, 0xe8, 0xbd, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x26, 0xc1, 0x00, 0x00, 0x01, 0xc9, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xf9, 0xcc, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x33, 0x00, 0xcc, 0xd8, 0x00, 0x00, 0x04, 0x00, 0x73, 0x35, 0x36, 0x00, 0x38, 0xc1, 0x65, 0x03, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x81, 0xec, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x36, 0x31, 0x00, 0x54, 0xf8, 0x00, 0x00, 0x01, 0xbb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x02, 0x00, 0x00, 0x00, 0x18, 0x08, 0x01, 0x00, 0x01, 0xb7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xeb, 0x13, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x31, 0x00, 0xbe, 0x1f, 0x01, 0x00, 0xaf, 0x23, 0x01, 0x00, 0xa0, 0x27, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x73, 0x33, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x39, 0x00, 0x3d, 0xb5, 0xc4, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x28, 0x47, 0x01, 0x00, 0x01, 0xa7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfb, 0x52, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x37, 0x00, 0xce, 0x5e, 0x01, 0x00, 0x04, 0x00, 0x73, 0x39, 0x30, 0x00, 0x1e, 0x8e, 0x6c, 0x05, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x83, 0x72, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x35, 0x00, 0x56, 0x7e, 0x01, 0x00, 0x01, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x02, 0x00, 0x00, 0x00, 0x1a, 0x8e, 0x01, 0x00, 0x01, 0x95, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xed, 0x99, 0x01, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x30, 0x35, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadStatsByPlayers = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ types.List[player_stats_types.StatboardQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		return types.NewList[player_stats_types.StatboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readStatsByPlayers(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, types.List[player_stats_types.StatboardQuery]{func() player_stats_types.StatboardQuery {
		v := player_stats_types.NewStatboardQuery()
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.ResetFrequency = types.NewUInt32(10095)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(16149)
			v.SortOrder = types.NewUInt32(17158)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(19176)
			v.SortOrder = types.NewUInt32(20185)
			return v
		}()}
		return v
	}(), func() player_stats_types.StatboardQuery {
		v := player_stats_types.NewStatboardQuery()
		v.BoardID = types.NewUInt32(22203)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(24221), types.NewUInt32(25230)}
		v.ResetFrequency = types.NewUInt32(26239)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(28257), types.NewUInt32(29266)}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(32293)
			v.SortOrder = types.NewUInt32(33302)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(35320)
			v.SortOrder = types.NewUInt32(36329)
			return v
		}()}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsNearPlayer_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readLeaderboardsNearPlayer(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(12113)
		v.ContextID = types.NewUInt32(13122)
		v.ResetFrequency = types.NewUInt32(14131)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsNearPlayer, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadLeaderboardsNearPlayer_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadLeaderboardsNearPlayer = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.LeaderboardResult](), nil
	}

	message, nexErr := commonProtocol.readLeaderboardsNearPlayer(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0), types.NewList[player_stats_types.LeaderboardQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsNearPlayer_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadLeaderboardsNearPlayer = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32, got2 types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("playerPID differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("count differs")
		}

		if !(got2.Equals(types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
			v := player_stats_types.NewLeaderboardQuery()
			v.BoardID = types.NewUInt32(5050)
			v.ContextID = types.NewUInt32(6059)
			v.ResetFrequency = types.NewUInt32(7068)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
			return v
		}(), func() player_stats_types.LeaderboardQuery {
			v := player_stats_types.NewLeaderboardQuery()
			v.BoardID = types.NewUInt32(12113)
			v.ContextID = types.NewUInt32(13122)
			v.ResetFrequency = types.NewUInt32(14131)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.LeaderboardResult]{func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(20185)
			v.ContextID = types.NewUInt32(21194)
			v.ResetFrequency = types.NewUInt32(22203)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(23212)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(27248)
					v.PlayerName = types.NewString("s28")
					v.SubmittedTime = types.NewDateTime(29000164)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(32293)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v33")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(35320)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-41)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(37338)
				v.Rank = types.NewUInt32(38347)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v39")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(42383)
					v.PlayerName = types.NewString("s43")
					v.SubmittedTime = types.NewDateTime(44000209)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(47428)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-53)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(50455)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v51")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(52473)
				v.Rank = types.NewUInt32(53482)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
				return v
			}()}
			return v
		}(), func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(56509)
			v.ContextID = types.NewUInt32(57518)
			v.ResetFrequency = types.NewUInt32(58527)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(59536)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(63572)
					v.PlayerName = types.NewString("s64")
					v.SubmittedTime = types.NewDateTime(65000272)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(68617)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v69")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(71644)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-77)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(73662)
				v.Rank = types.NewUInt32(74671)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v75")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(78707)
					v.PlayerName = types.NewString("s79")
					v.SubmittedTime = types.NewDateTime(80000317)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(83752)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-89)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(86779)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v87")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(88797)
				v.Rank = types.NewUInt32(89806)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-95)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readLeaderboardsNearPlayer(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(12113)
		v.ContextID = types.NewUInt32(13122)
		v.ResetFrequency = types.NewUInt32(14131)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsNearPlayer, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x38, 0x00, 0xe4, 0x81, 0xba, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x25, 0x7e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x33, 0x00, 0xf8, 0x89, 0x00, 0x00, 0x01, 0xd7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xda, 0x91, 0x00, 0x00, 0xcb, 0x95, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x39, 0x00, 0x8f, 0xa5, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x33, 0x00, 0xd1, 0x63, 0x9f, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x44, 0xb9, 0x00, 0x00, 0x01, 0xcb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x17, 0xc5, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x31, 0x00, 0xf9, 0xcc, 0x00, 0x00, 0xea, 0xd0, 0x00, 0x00, 0x01, 0xc5, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbd, 0xdc, 0x00, 0x00, 0xae, 0xe0, 0x00, 0x00, 0x9f, 0xe4, 0x00, 0x00, 0x90, 0xe8, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x54, 0xf8, 0x00, 0x00, 0x04, 0x00, 0x73, 0x36, 0x34, 0x00, 0x50, 0xd3, 0xdf, 0x03, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x09, 0x0c, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x36, 0x39, 0x00, 0xdc, 0x17, 0x01, 0x00, 0x01, 0xb3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbe, 0x1f, 0x01, 0x00, 0xaf, 0x23, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x35, 0x00, 0x73, 0x33, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x39, 0x00, 0x3d, 0xb5, 0xc4, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x28, 0x47, 0x01, 0x00, 0x01, 0xa7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfb, 0x52, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x37, 0x00, 0xdd, 0x5a, 0x01, 0x00, 0xce, 0x5e, 0x01, 0x00, 0x01, 0xa1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadLeaderboardsNearPlayer = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		return types.NewList[player_stats_types.LeaderboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readLeaderboardsNearPlayer(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(12113)
		v.ContextID = types.NewUInt32(13122)
		v.ResetFrequency = types.NewUInt32(14131)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByRank_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readLeaderboardsByRank(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(12113)
		v.ContextID = types.NewUInt32(13122)
		v.ResetFrequency = types.NewUInt32(14131)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByRank, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadLeaderboardsByRank_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadLeaderboardsByRank = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.LeaderboardResult](), nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByRank(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0), types.NewList[player_stats_types.LeaderboardQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByRank_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadLeaderboardsByRank = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32, got2 types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("startingRank differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("count differs")
		}

		if !(got2.Equals(types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
			v := player_stats_types.NewLeaderboardQuery()
			v.BoardID = types.NewUInt32(5050)
			v.ContextID = types.NewUInt32(6059)
			v.ResetFrequency = types.NewUInt32(7068)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
			return v
		}(), func() player_stats_types.LeaderboardQuery {
			v := player_stats_types.NewLeaderboardQuery()
			v.BoardID = types.NewUInt32(12113)
			v.ContextID = types.NewUInt32(13122)
			v.ResetFrequency = types.NewUInt32(14131)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.LeaderboardResult]{func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(20185)
			v.ContextID = types.NewUInt32(21194)
			v.ResetFrequency = types.NewUInt32(22203)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(23212)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(27248)
					v.PlayerName = types.NewString("s28")
					v.SubmittedTime = types.NewDateTime(29000164)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(32293)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v33")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(35320)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-41)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(37338)
				v.Rank = types.NewUInt32(38347)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v39")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(42383)
					v.PlayerName = types.NewString("s43")
					v.SubmittedTime = types.NewDateTime(44000209)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(47428)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-53)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(50455)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v51")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(52473)
				v.Rank = types.NewUInt32(53482)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
				return v
			}()}
			return v
		}(), func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(56509)
			v.ContextID = types.NewUInt32(57518)
			v.ResetFrequency = types.NewUInt32(58527)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(59536)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(63572)
					v.PlayerName = types.NewString("s64")
					v.SubmittedTime = types.NewDateTime(65000272)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(68617)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v69")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(71644)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-77)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(73662)
				v.Rank = types.NewUInt32(74671)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v75")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(78707)
					v.PlayerName = types.NewString("s79")
					v.SubmittedTime = types.NewDateTime(80000317)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(83752)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-89)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(86779)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v87")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(88797)
				v.Rank = types.NewUInt32(89806)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-95)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByRank(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(12113)
		v.ContextID = types.NewUInt32(13122)
		v.ResetFrequency = types.NewUInt32(14131)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByRank, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x38, 0x00, 0xe4, 0x81, 0xba, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x25, 0x7e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x33, 0x00, 0xf8, 0x89, 0x00, 0x00, 0x01, 0xd7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xda, 0x91, 0x00, 0x00, 0xcb, 0x95, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x39, 0x00, 0x8f, 0xa5, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x33, 0x00, 0xd1, 0x63, 0x9f, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x44, 0xb9, 0x00, 0x00, 0x01, 0xcb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x17, 0xc5, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x31, 0x00, 0xf9, 0xcc, 0x00, 0x00, 0xea, 0xd0, 0x00, 0x00, 0x01, 0xc5, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbd, 0xdc, 0x00, 0x00, 0xae, 0xe0, 0x00, 0x00, 0x9f, 0xe4, 0x00, 0x00, 0x90, 0xe8, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x54, 0xf8, 0x00, 0x00, 0x04, 0x00, 0x73, 0x36, 0x34, 0x00, 0x50, 0xd3, 0xdf, 0x03, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x09, 0x0c, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x36, 0x39, 0x00, 0xdc, 0x17, 0x01, 0x00, 0x01, 0xb3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbe, 0x1f, 0x01, 0x00, 0xaf, 0x23, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x35, 0x00, 0x73, 0x33, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x39, 0x00, 0x3d, 0xb5, 0xc4, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x28, 0x47, 0x01, 0x00, 0x01, 0xa7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfb, 0x52, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x37, 0x00, 0xdd, 0x5a, 0x01, 0x00, 0xce, 0x5e, 0x01, 0x00, 0x01, 0xa1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadLeaderboardsByRank = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		return types.NewList[player_stats_types.LeaderboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readLeaderboardsByRank(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(12113)
		v.ContextID = types.NewUInt32(13122)
		v.ResetFrequency = types.NewUInt32(14131)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(16149), types.NewUInt32(17158)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByPlayers_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readLeaderboardsByPlayers(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(6059)
		v.ContextID = types.NewUInt32(7068)
		v.ResetFrequency = types.NewUInt32(8077)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(13122)
		v.ContextID = types.NewUInt32(14131)
		v.ResetFrequency = types.NewUInt32(15140)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByPlayers, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadLeaderboardsByPlayers_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadLeaderboardsByPlayers = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.LeaderboardResult](), nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByPlayers(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32](), types.NewList[player_stats_types.LeaderboardQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByPlayers_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadLeaderboardsByPlayers = func(packet nex.PacketInterface, got0 types.List[types.UInt32], got1 types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("playerPIDs differs")
		}

		if !(got1.Equals(types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
			v := player_stats_types.NewLeaderboardQuery()
			v.BoardID = types.NewUInt32(6059)
			v.ContextID = types.NewUInt32(7068)
			v.ResetFrequency = types.NewUInt32(8077)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
			return v
		}(), func() player_stats_types.LeaderboardQuery {
			v := player_stats_types.NewLeaderboardQuery()
			v.BoardID = types.NewUInt32(13122)
			v.ContextID = types.NewUInt32(14131)
			v.ResetFrequency = types.NewUInt32(15140)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.LeaderboardResult]{func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(21194)
			v.ContextID = types.NewUInt32(22203)
			v.ResetFrequency = types.NewUInt32(23212)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(24221)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(28257)
					v.PlayerName = types.NewString("s29")
					v.SubmittedTime = types.NewDateTime(30000167)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(33302)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-39)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(36329)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v37")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(38347)
				v.Rank = types.NewUInt32(39356)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-45)}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(43392)
					v.PlayerName = types.NewString("s44")
					v.SubmittedTime = types.NewDateTime(45000212)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(48437)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v49")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(51464)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-57)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(53482)
				v.Rank = types.NewUInt32(54491)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v55")}
				return v
			}()}
			return v
		}(), func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(57518)
			v.ContextID = types.NewUInt32(58527)
			v.ResetFrequency = types.NewUInt32(59536)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(60545)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(64581)
					v.PlayerName = types.NewString("s65")
					v.SubmittedTime = types.NewDateTime(66000275)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(69626)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-75)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(72653)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v73")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(74671)
				v.Rank = types.NewUInt32(75680)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-81)}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(79716)
					v.PlayerName = types.NewString("s80")
					v.SubmittedTime = types.NewDateTime(81000320)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(84761)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v85")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(87788)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-93)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(89806)
				v.Rank = types.NewUInt32(90815)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v91")}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByPlayers(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(6059)
		v.ContextID = types.NewUInt32(7068)
		v.ResetFrequency = types.NewUInt32(8077)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(13122)
		v.ContextID = types.NewUInt32(14131)
		v.ResetFrequency = types.NewUInt32(15140)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByPlayers, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x9d, 0x5e, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x61, 0x6e, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x39, 0x00, 0x27, 0xc4, 0xc9, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x16, 0x82, 0x00, 0x00, 0x01, 0xd9, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe9, 0x8d, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x37, 0x00, 0xcb, 0x95, 0x00, 0x00, 0xbc, 0x99, 0x00, 0x00, 0x01, 0xd3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x80, 0xa9, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x34, 0x00, 0x14, 0xa6, 0xae, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x35, 0xbd, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x34, 0x39, 0x00, 0x08, 0xc9, 0x00, 0x00, 0x01, 0xc7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xea, 0xd0, 0x00, 0x00, 0xdb, 0xd4, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x35, 0x00, 0xae, 0xe0, 0x00, 0x00, 0x9f, 0xe4, 0x00, 0x00, 0x90, 0xe8, 0x00, 0x00, 0x81, 0xec, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x45, 0xfc, 0x00, 0x00, 0x04, 0x00, 0x73, 0x36, 0x35, 0x00, 0x93, 0x15, 0xef, 0x03, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xfa, 0x0f, 0x01, 0x00, 0x01, 0xb5, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xcd, 0x1b, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x33, 0x00, 0xaf, 0x23, 0x01, 0x00, 0xa0, 0x27, 0x01, 0x00, 0x01, 0xaf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x64, 0x37, 0x01, 0x00, 0x04, 0x00, 0x73, 0x38, 0x30, 0x00, 0x80, 0xf7, 0xd3, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x4b, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x35, 0x00, 0xec, 0x56, 0x01, 0x00, 0x01, 0xa3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xce, 0x5e, 0x01, 0x00, 0xbf, 0x62, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x31, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadLeaderboardsByPlayers = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		return types.NewList[player_stats_types.LeaderboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readLeaderboardsByPlayers(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, types.List[player_stats_types.LeaderboardQuery]{func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(6059)
		v.ContextID = types.NewUInt32(7068)
		v.ResetFrequency = types.NewUInt32(8077)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
		return v
	}(), func() player_stats_types.LeaderboardQuery {
		v := player_stats_types.NewLeaderboardQuery()
		v.BoardID = types.NewUInt32(13122)
		v.ContextID = types.NewUInt32(14131)
		v.ResetFrequency = types.NewUInt32(15140)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadStatboardHistory_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readStatboardHistory(nil, fakePacket{}, 7, types.List[player_stats_types.StatboardHistoryQuery]{func() player_stats_types.StatboardHistoryQuery {
		v := player_stats_types.NewStatboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(16149)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(19176)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(23000146)
			v.EndingDatetime = types.NewDateTime(24000149)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(30275)
			v.SortOrder = types.NewUInt32(31284)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(33302)
			v.SortOrder = types.NewUInt32(34311)
			return v
		}()}
		return v
	}(), func() player_stats_types.StatboardHistoryQuery {
		v := player_stats_types.NewStatboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(37338), types.NewUInt32(38347)}
		v.BoardID = types.NewUInt32(39356)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(41374), types.NewUInt32(42383)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(44401), types.NewUInt32(45410)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(138), Length: types.NewUInt32(148)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(49446)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-55)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(52473)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v53")}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(56000245)
			v.EndingDatetime = types.NewDateTime(57000248)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(59000254)
			v.EndingDatetime = types.NewDateTime(60000257)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(63572)
			v.SortOrder = types.NewUInt32(64581)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(66599)
			v.SortOrder = types.NewUInt32(67608)
			return v
		}()}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadStatboardHistory, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadStatboardHistory_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadStatboardHistory = func(packet nex.PacketInterface, _ types.List[player_stats_types.StatboardHistoryQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.StatboardResult](), nil
	}

	message, nexErr := commonProtocol.readStatboardHistory(errors.New("bad"), fakePacket{}, 7, types.NewList[player_stats_types.StatboardHistoryQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadStatboardHistory_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadStatboardHistory = func(packet nex.PacketInterface, got0 types.List[player_stats_types.StatboardHistoryQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[player_stats_types.StatboardHistoryQuery]{func() player_stats_types.StatboardHistoryQuery {
			v := player_stats_types.NewStatboardHistoryQuery()
			v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
			v.BoardID = types.NewUInt32(6059)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
			v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(16149)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(19176)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
				return v
			}()}
			v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(23000146)
				v.EndingDatetime = types.NewDateTime(24000149)
				return v
			}(), func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(26000155)
				v.EndingDatetime = types.NewDateTime(27000158)
				return v
			}()}
			v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(30275)
				v.SortOrder = types.NewUInt32(31284)
				return v
			}(), func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(33302)
				v.SortOrder = types.NewUInt32(34311)
				return v
			}()}
			return v
		}(), func() player_stats_types.StatboardHistoryQuery {
			v := player_stats_types.NewStatboardHistoryQuery()
			v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(37338), types.NewUInt32(38347)}
			v.BoardID = types.NewUInt32(39356)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(41374), types.NewUInt32(42383)}
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(44401), types.NewUInt32(45410)}
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(138), Length: types.NewUInt32(148)}
			v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(49446)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-55)}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(52473)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v53")}
				return v
			}()}
			v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(56000245)
				v.EndingDatetime = types.NewDateTime(57000248)
				return v
			}(), func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(59000254)
				v.EndingDatetime = types.NewDateTime(60000257)
				return v
			}()}
			v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(63572)
				v.SortOrder = types.NewUInt32(64581)
				return v
			}(), func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(66599)
				v.SortOrder = types.NewUInt32(67608)
				return v
			}()}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.StatboardResult]{func() player_stats_types.StatboardResult {
			v := player_stats_types.NewStatboardResult()
			v.BoardID = types.NewUInt32(70635)
			v.ContextID = types.NewUInt32(71644)
			v.ResetFrequency = types.NewUInt32(72653)
			v.PlayerStatSets = types.List[player_stats_types.PlayerStatSet]{func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(75680)
				v.PlayerName = types.NewString("s76")
				v.SubmittedTime = types.NewDateTime(77000308)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(80725)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v81")}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(83752)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-89)}
					return v
				}()}
				return v
			}(), func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(86779)
				v.PlayerName = types.NewString("s87")
				v.SubmittedTime = types.NewDateTime(88000341)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(91824)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-97)}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(94851)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v95")}
					return v
				}()}
				return v
			}()}
			v.DefaultStatValues = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(98887)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v99")}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(101914)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-107)}
				return v
			}()}
			return v
		}(), func() player_stats_types.StatboardResult {
			v := player_stats_types.NewStatboardResult()
			v.BoardID = types.NewUInt32(104941)
			v.ContextID = types.NewUInt32(105950)
			v.ResetFrequency = types.NewUInt32(106959)
			v.PlayerStatSets = types.List[player_stats_types.PlayerStatSet]{func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(109986)
				v.PlayerName = types.NewString("s110")
				v.SubmittedTime = types.NewDateTime(111000410)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(115031)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v115")}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(118058)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-123)}
					return v
				}()}
				return v
			}(), func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(121085)
				v.PlayerName = types.NewString("s121")
				v.SubmittedTime = types.NewDateTime(122000443)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(126130)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-131)}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(129157)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v129")}
					return v
				}()}
				return v
			}()}
			v.DefaultStatValues = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(133193)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v133")}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(136220)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-141)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readStatboardHistory(nil, fakePacket{}, 9, types.List[player_stats_types.StatboardHistoryQuery]{func() player_stats_types.StatboardHistoryQuery {
		v := player_stats_types.NewStatboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(16149)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(19176)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(23000146)
			v.EndingDatetime = types.NewDateTime(24000149)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(30275)
			v.SortOrder = types.NewUInt32(31284)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(33302)
			v.SortOrder = types.NewUInt32(34311)
			return v
		}()}
		return v
	}(), func() player_stats_types.StatboardHistoryQuery {
		v := player_stats_types.NewStatboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(37338), types.NewUInt32(38347)}
		v.BoardID = types.NewUInt32(39356)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(41374), types.NewUInt32(42383)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(44401), types.NewUInt32(45410)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(138), Length: types.NewUInt32(148)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(49446)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-55)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(52473)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v53")}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(56000245)
			v.EndingDatetime = types.NewDateTime(57000248)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(59000254)
			v.EndingDatetime = types.NewDateTime(60000257)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(63572)
			v.SortOrder = types.NewUInt32(64581)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(66599)
			v.SortOrder = types.NewUInt32(67608)
			return v
		}()}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadStatboardHistory, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xeb, 0x13, 0x01, 0x00, 0xdc, 0x17, 0x01, 0x00, 0xcd, 0x1b, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0xa0, 0x27, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x36, 0x00, 0x74, 0xee, 0x96, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x55, 0x3b, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x31, 0x00, 0x28, 0x47, 0x01, 0x00, 0x01, 0xa7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfb, 0x52, 0x01, 0x00, 0x04, 0x00, 0x73, 0x38, 0x37, 0x00, 0x55, 0xc7, 0x3e, 0x05, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xb0, 0x66, 0x01, 0x00, 0x01, 0x9f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x83, 0x72, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x35, 0x00, 0x02, 0x00, 0x00, 0x00, 0x47, 0x82, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x39, 0x00, 0x1a, 0x8e, 0x01, 0x00, 0x01, 0x95, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xed, 0x99, 0x01, 0x00, 0xde, 0x9d, 0x01, 0x00, 0xcf, 0xa1, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0xa2, 0xad, 0x01, 0x00, 0x05, 0x00, 0x73, 0x31, 0x31, 0x30, 0x00, 0x5a, 0xbb, 0x9d, 0x06, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x57, 0xc1, 0x01, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x31, 0x35, 0x00, 0x2a, 0xcd, 0x01, 0x00, 0x01, 0x85, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfd, 0xd8, 0x01, 0x00, 0x05, 0x00, 0x73, 0x31, 0x32, 0x31, 0x00, 0x3b, 0x94, 0x45, 0x07, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xb2, 0xec, 0x01, 0x00, 0x01, 0x7d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x85, 0xf8, 0x01, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x32, 0x39, 0x00, 0x02, 0x00, 0x00, 0x00, 0x49, 0x08, 0x02, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x33, 0x33, 0x00, 0x1c, 0x14, 0x02, 0x00, 0x01, 0x73, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadStatboardHistory = func(packet nex.PacketInterface, _ types.List[player_stats_types.StatboardHistoryQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		return types.NewList[player_stats_types.StatboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readStatboardHistory(nil, fakePacket{}, 9, types.List[player_stats_types.StatboardHistoryQuery]{func() player_stats_types.StatboardHistoryQuery {
		v := player_stats_types.NewStatboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(16149)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(19176)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(23000146)
			v.EndingDatetime = types.NewDateTime(24000149)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(30275)
			v.SortOrder = types.NewUInt32(31284)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(33302)
			v.SortOrder = types.NewUInt32(34311)
			return v
		}()}
		return v
	}(), func() player_stats_types.StatboardHistoryQuery {
		v := player_stats_types.NewStatboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(37338), types.NewUInt32(38347)}
		v.BoardID = types.NewUInt32(39356)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(41374), types.NewUInt32(42383)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(44401), types.NewUInt32(45410)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(138), Length: types.NewUInt32(148)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(49446)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-55)}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(52473)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v53")}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(56000245)
			v.EndingDatetime = types.NewDateTime(57000248)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(59000254)
			v.EndingDatetime = types.NewDateTime(60000257)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(63572)
			v.SortOrder = types.NewUInt32(64581)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(66599)
			v.SortOrder = types.NewUInt32(67608)
			return v
		}()}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadLeaderboardHistory_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readLeaderboardHistory(nil, fakePacket{}, 7, types.List[player_stats_types.LeaderboardHistoryQuery]{func() player_stats_types.LeaderboardHistoryQuery {
		v := player_stats_types.NewLeaderboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextID = types.NewUInt32(7068)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(24), Length: types.NewUInt32(34)}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(11000110)
			v.EndingDatetime = types.NewDateTime(12000113)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(14000119)
			v.EndingDatetime = types.NewDateTime(15000122)
			return v
		}()}
		v.FilterOption = types.NewUInt32(16149)
		return v
	}(), func() player_stats_types.LeaderboardHistoryQuery {
		v := player_stats_types.NewLeaderboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.BoardID = types.NewUInt32(21194)
		v.ContextID = types.NewUInt32(22203)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(69), Length: types.NewUInt32(79)}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(29000164)
			v.EndingDatetime = types.NewDateTime(30000167)
			return v
		}()}
		v.FilterOption = types.NewUInt32(31284)
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardHistory, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadLeaderboardHistory_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadLeaderboardHistory = func(packet nex.PacketInterface, _ types.List[player_stats_types.LeaderboardHistoryQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.LeaderboardResult](), nil
	}

	message, nexErr := commonProtocol.readLeaderboardHistory(errors.New("bad"), fakePacket{}, 7, types.NewList[player_stats_types.LeaderboardHistoryQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadLeaderboardHistory_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadLeaderboardHistory = func(packet nex.PacketInterface, got0 types.List[player_stats_types.LeaderboardHistoryQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[player_stats_types.LeaderboardHistoryQuery]{func() player_stats_types.LeaderboardHistoryQuery {
			v := player_stats_types.NewLeaderboardHistoryQuery()
			v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
			v.BoardID = types.NewUInt32(6059)
			v.ContextID = types.NewUInt32(7068)
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(24), Length: types.NewUInt32(34)}
			v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(11000110)
				v.EndingDatetime = types.NewDateTime(12000113)
				return v
			}(), func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(14000119)
				v.EndingDatetime = types.NewDateTime(15000122)
				return v
			}()}
			v.FilterOption = types.NewUInt32(16149)
			return v
		}(), func() player_stats_types.LeaderboardHistoryQuery {
			v := player_stats_types.NewLeaderboardHistoryQuery()
			v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
			v.BoardID = types.NewUInt32(21194)
			v.ContextID = types.NewUInt32(22203)
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(69), Length: types.NewUInt32(79)}
			v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(26000155)
				v.EndingDatetime = types.NewDateTime(27000158)
				return v
			}(), func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(29000164)
				v.EndingDatetime = types.NewDateTime(30000167)
				return v
			}()}
			v.FilterOption = types.NewUInt32(31284)
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.LeaderboardResult]{func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(34311)
			v.ContextID = types.NewUInt32(35320)
			v.ResetFrequency = types.NewUInt32(36329)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(37338)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(41374)
					v.PlayerName = types.NewString("s42")
					v.SubmittedTime = types.NewDateTime(43000206)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(46419)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v47")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(49446)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-55)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(51464)
				v.Rank = types.NewUInt32(52473)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v53")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(56509)
					v.PlayerName = types.NewString("s57")
					v.SubmittedTime = types.NewDateTime(58000251)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(61554)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-67)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(64581)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v65")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(66599)
				v.Rank = types.NewUInt32(67608)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-73)}
				return v
			}()}
			return v
		}(), func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(70635)
			v.ContextID = types.NewUInt32(71644)
			v.ResetFrequency = types.NewUInt32(72653)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(73662)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(77698)
					v.PlayerName = types.NewString("s78")
					v.SubmittedTime = types.NewDateTime(79000314)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(82743)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v83")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(85770)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-91)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(87788)
				v.Rank = types.NewUInt32(88797)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v89")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(92833)
					v.PlayerName = types.NewString("s93")
					v.SubmittedTime = types.NewDateTime(94000359)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(97878)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-103)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(100905)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v101")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(102923)
				v.Rank = types.NewUInt32(103932)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-109)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readLeaderboardHistory(nil, fakePacket{}, 9, types.List[player_stats_types.LeaderboardHistoryQuery]{func() player_stats_types.LeaderboardHistoryQuery {
		v := player_stats_types.NewLeaderboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextID = types.NewUInt32(7068)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(24), Length: types.NewUInt32(34)}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(11000110)
			v.EndingDatetime = types.NewDateTime(12000113)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(14000119)
			v.EndingDatetime = types.NewDateTime(15000122)
			return v
		}()}
		v.FilterOption = types.NewUInt32(16149)
		return v
	}(), func() player_stats_types.LeaderboardHistoryQuery {
		v := player_stats_types.NewLeaderboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.BoardID = types.NewUInt32(21194)
		v.ContextID = types.NewUInt32(22203)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(69), Length: types.NewUInt32(79)}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(29000164)
			v.EndingDatetime = types.NewDateTime(30000167)
			return v
		}()}
		v.FilterOption = types.NewUInt32(31284)
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardHistory, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x07, 0x86, 0x00, 0x00, 0xf8, 0x89, 0x00, 0x00, 0xe9, 0x8d, 0x00, 0x00, 0xda, 0x91, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x9e, 0xa1, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x32, 0x00, 0x8e, 0x21, 0x90, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x53, 0xb5, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x34, 0x37, 0x00, 0x26, 0xc1, 0x00, 0x00, 0x01, 0xc9, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x08, 0xc9, 0x00, 0x00, 0xf9, 0xcc, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x33, 0x00, 0xbd, 0xdc, 0x00, 0x00, 0x04, 0x00, 0x73, 0x35, 0x37, 0x00, 0x7b, 0x03, 0x75, 0x03, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x72, 0xf0, 0x00, 0x00, 0x01, 0xbd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x45, 0xfc, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x36, 0x35, 0x00, 0x27, 0x04, 0x01, 0x00, 0x18, 0x08, 0x01, 0x00, 0x01, 0xb7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xeb, 0x13, 0x01, 0x00, 0xdc, 0x17, 0x01, 0x00, 0xcd, 0x1b, 0x01, 0x00, 0xbe, 0x1f, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x82, 0x2f, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x38, 0x00, 0xfa, 0x72, 0xb5, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x37, 0x43, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x33, 0x00, 0x0a, 0x4f, 0x01, 0x00, 0x01, 0xa5, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xec, 0x56, 0x01, 0x00, 0xdd, 0x5a, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x39, 0x00, 0xa1, 0x6a, 0x01, 0x00, 0x04, 0x00, 0x73, 0x39, 0x33, 0x00, 0xe7, 0x54, 0x9a, 0x05, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x56, 0x7e, 0x01, 0x00, 0x01, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x29, 0x8a, 0x01, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x30, 0x31, 0x00, 0x0b, 0x92, 0x01, 0x00, 0xfc, 0x95, 0x01, 0x00, 0x01, 0x93, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadLeaderboardHistory = func(packet nex.PacketInterface, _ types.List[player_stats_types.LeaderboardHistoryQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		return types.NewList[player_stats_types.LeaderboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readLeaderboardHistory(nil, fakePacket{}, 9, types.List[player_stats_types.LeaderboardHistoryQuery]{func() player_stats_types.LeaderboardHistoryQuery {
		v := player_stats_types.NewLeaderboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextID = types.NewUInt32(7068)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(24), Length: types.NewUInt32(34)}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(11000110)
			v.EndingDatetime = types.NewDateTime(12000113)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(14000119)
			v.EndingDatetime = types.NewDateTime(15000122)
			return v
		}()}
		v.FilterOption = types.NewUInt32(16149)
		return v
	}(), func() player_stats_types.LeaderboardHistoryQuery {
		v := player_stats_types.NewLeaderboardHistoryQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.BoardID = types.NewUInt32(21194)
		v.ContextID = types.NewUInt32(22203)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(69), Length: types.NewUInt32(79)}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(29000164)
			v.EndingDatetime = types.NewDateTime(30000167)
			return v
		}()}
		v.FilterOption = types.NewUInt32(31284)
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadStatboardHistoryAggregated_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readStatboardHistoryAggregated(nil, fakePacket{}, 7, types.List[player_stats_types.StatboardHistoryAggregatedQuery]{func() player_stats_types.StatboardHistoryAggregatedQuery {
		v := player_stats_types.NewStatboardHistoryAggregatedQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(16149)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(19176)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(23000146)
			v.EndingDatetime = types.NewDateTime(24000149)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(30275)
			v.SortOrder = types.NewUInt32(31284)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(33302)
			v.SortOrder = types.NewUInt32(34311)
			return v
		}()}
		v.FilterOption = types.NewUInt32(35320)
		return v
	}(), func() player_stats_types.StatboardHistoryAggregatedQuery {
		v := player_stats_types.NewStatboardHistoryAggregatedQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(38347), types.NewUInt32(39356)}
		v.BoardID = types.NewUInt32(40365)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(42383), types.NewUInt32(43392)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(45410), types.NewUInt32(46419)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(141), Length: types.NewUInt32(151)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(50455)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v51")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(53482)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(57000248)
			v.EndingDatetime = types.NewDateTime(58000251)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(60000257)
			v.EndingDatetime = types.NewDateTime(61000260)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(64581)
			v.SortOrder = types.NewUInt32(65590)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(67608)
			v.SortOrder = types.NewUInt32(68617)
			return v
		}()}
		v.FilterOption = types.NewUInt32(69626)
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadStatboardHistoryAggregated, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadStatboardHistoryAggregated_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadStatboardHistoryAggregated = func(packet nex.PacketInterface, _ types.List[player_stats_types.StatboardHistoryAggregatedQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.StatboardResult](), nil
	}

	message, nexErr := commonProtocol.readStatboardHistoryAggregated(errors.New("bad"), fakePacket{}, 7, types.NewList[player_stats_types.StatboardHistoryAggregatedQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadStatboardHistoryAggregated_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadStatboardHistoryAggregated = func(packet nex.PacketInterface, got0 types.List[player_stats_types.StatboardHistoryAggregatedQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[player_stats_types.StatboardHistoryAggregatedQuery]{func() player_stats_types.StatboardHistoryAggregatedQuery {
			v := player_stats_types.NewStatboardHistoryAggregatedQuery()
			v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
			v.BoardID = types.NewUInt32(6059)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
			v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(16149)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(19176)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
				return v
			}()}
			v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(23000146)
				v.EndingDatetime = types.NewDateTime(24000149)
				return v
			}(), func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(26000155)
				v.EndingDatetime = types.NewDateTime(27000158)
				return v
			}()}
			v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(30275)
				v.SortOrder = types.NewUInt32(31284)
				return v
			}(), func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(33302)
				v.SortOrder = types.NewUInt32(34311)
				return v
			}()}
			v.FilterOption = types.NewUInt32(35320)
			return v
		}(), func() player_stats_types.StatboardHistoryAggregatedQuery {
			v := player_stats_types.NewStatboardHistoryAggregatedQuery()
			v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(38347), types.NewUInt32(39356)}
			v.BoardID = types.NewUInt32(40365)
			v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(42383), types.NewUInt32(43392)}
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(45410), types.NewUInt32(46419)}
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(141), Length: types.NewUInt32(151)}
			v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(50455)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v51")}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(53482)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
				return v
			}()}
			v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(57000248)
				v.EndingDatetime = types.NewDateTime(58000251)
				return v
			}(), func() player_stats_types.DateRange {
				v := player_stats_types.NewDateRange()
				v.StartingDatetime = types.NewDateTime(60000257)
				v.EndingDatetime = types.NewDateTime(61000260)
				return v
			}()}
			v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(64581)
				v.SortOrder = types.NewUInt32(65590)
				return v
			}(), func() player_stats_types.PlayerStatSortCriteria {
				v := player_stats_types.NewPlayerStatSortCriteria()
				v.SortStatID = types.NewUInt32(67608)
				v.SortOrder = types.NewUInt32(68617)
				return v
			}()}
			v.FilterOption = types.NewUInt32(69626)
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.StatboardResult]{func() player_stats_types.StatboardResult {
			v := player_stats_types.NewStatboardResult()
			v.BoardID = types.NewUInt32(72653)
			v.ContextID = types.NewUInt32(73662)
			v.ResetFrequency = types.NewUInt32(74671)
			v.PlayerStatSets = types.List[player_stats_types.PlayerStatSet]{func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(77698)
				v.PlayerName = types.NewString("s78")
				v.SubmittedTime = types.NewDateTime(79000314)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(82743)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v83")}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(85770)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-91)}
					return v
				}()}
				return v
			}(), func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(88797)
				v.PlayerName = types.NewString("s89")
				v.SubmittedTime = types.NewDateTime(90000347)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(93842)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-99)}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(96869)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v97")}
					return v
				}()}
				return v
			}()}
			v.DefaultStatValues = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(100905)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v101")}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(103932)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-109)}
				return v
			}()}
			return v
		}(), func() player_stats_types.StatboardResult {
			v := player_stats_types.NewStatboardResult()
			v.BoardID = types.NewUInt32(106959)
			v.ContextID = types.NewUInt32(107968)
			v.ResetFrequency = types.NewUInt32(108977)
			v.PlayerStatSets = types.List[player_stats_types.PlayerStatSet]{func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(112004)
				v.PlayerName = types.NewString("s112")
				v.SubmittedTime = types.NewDateTime(113000416)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(117049)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v117")}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(120076)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-125)}
					return v
				}()}
				return v
			}(), func() player_stats_types.PlayerStatSet {
				v := player_stats_types.NewPlayerStatSet()
				v.PlayerPID = types.NewUInt32(123103)
				v.PlayerName = types.NewString("s123")
				v.SubmittedTime = types.NewDateTime(124000449)
				v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(128148)
					v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-133)}
					return v
				}(), func() player_stats_types.PropertyVariant {
					v := player_stats_types.NewPropertyVariant()
					v.ID = types.NewUInt32(131175)
					v.Value = types.Variant{TypeID: 4, Type: types.NewString("v131")}
					return v
				}()}
				return v
			}()}
			v.DefaultStatValues = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(135211)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v135")}
				return v
			}(), func() player_stats_types.PropertyVariant {
				v := player_stats_types.NewPropertyVariant()
				v.ID = types.NewUInt32(138238)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-143)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readStatboardHistoryAggregated(nil, fakePacket{}, 9, types.List[player_stats_types.StatboardHistoryAggregatedQuery]{func() player_stats_types.StatboardHistoryAggregatedQuery {
		v := player_stats_types.NewStatboardHistoryAggregatedQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(16149)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(19176)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(23000146)
			v.EndingDatetime = types.NewDateTime(24000149)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(30275)
			v.SortOrder = types.NewUInt32(31284)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(33302)
			v.SortOrder = types.NewUInt32(34311)
			return v
		}()}
		v.FilterOption = types.NewUInt32(35320)
		return v
	}(), func() player_stats_types.StatboardHistoryAggregatedQuery {
		v := player_stats_types.NewStatboardHistoryAggregatedQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(38347), types.NewUInt32(39356)}
		v.BoardID = types.NewUInt32(40365)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(42383), types.NewUInt32(43392)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(45410), types.NewUInt32(46419)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(141), Length: types.NewUInt32(151)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(50455)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v51")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(53482)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(57000248)
			v.EndingDatetime = types.NewDateTime(58000251)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(60000257)
			v.EndingDatetime = types.NewDateTime(61000260)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(64581)
			v.SortOrder = types.NewUInt32(65590)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(67608)
			v.SortOrder = types.NewUInt32(68617)
			return v
		}()}
		v.FilterOption = types.NewUInt32(69626)
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadStatboardHistoryAggregated, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xcd, 0x1b, 0x01, 0x00, 0xbe, 0x1f, 0x01, 0x00, 0xaf, 0x23, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x82, 0x2f, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x38, 0x00, 0xfa, 0x72, 0xb5, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x37, 0x43, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x33, 0x00, 0x0a, 0x4f, 0x01, 0x00, 0x01, 0xa5, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xdd, 0x5a, 0x01, 0x00, 0x04, 0x00, 0x73, 0x38, 0x39, 0x00, 0xdb, 0x4b, 0x5d, 0x05, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x92, 0x6e, 0x01, 0x00, 0x01, 0x9d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x65, 0x7a, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x37, 0x00, 0x02, 0x00, 0x00, 0x00, 0x29, 0x8a, 0x01, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x30, 0x31, 0x00, 0xfc, 0x95, 0x01, 0x00, 0x01, 0x93, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xcf, 0xa1, 0x01, 0x00, 0xc0, 0xa5, 0x01, 0x00, 0xb1, 0xa9, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x84, 0xb5, 0x01, 0x00, 0x05, 0x00, 0x73, 0x31, 0x31, 0x32, 0x00, 0xe0, 0x3f, 0xbc, 0x06, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x39, 0xc9, 0x01, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x31, 0x37, 0x00, 0x0c, 0xd5, 0x01, 0x00, 0x01, 0x83, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xdf, 0xe0, 0x01, 0x00, 0x05, 0x00, 0x73, 0x31, 0x32, 0x33, 0x00, 0xc1, 0x18, 0x64, 0x07, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x94, 0xf4, 0x01, 0x00, 0x01, 0x7b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x67, 0x00, 0x02, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x33, 0x31, 0x00, 0x02, 0x00, 0x00, 0x00, 0x2b, 0x10, 0x02, 0x00, 0x04, 0x05, 0x00, 0x76, 0x31, 0x33, 0x35, 0x00, 0xfe, 0x1b, 0x02, 0x00, 0x01, 0x71, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadStatboardHistoryAggregated = func(packet nex.PacketInterface, _ types.List[player_stats_types.StatboardHistoryAggregatedQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error) {
		return types.NewList[player_stats_types.StatboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readStatboardHistoryAggregated(nil, fakePacket{}, 9, types.List[player_stats_types.StatboardHistoryAggregatedQuery]{func() player_stats_types.StatboardHistoryAggregatedQuery {
		v := player_stats_types.NewStatboardHistoryAggregatedQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(4041), types.NewUInt32(5050)}
		v.BoardID = types.NewUInt32(6059)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(11104), types.NewUInt32(12113)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(39), Length: types.NewUInt32(49)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(16149)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(19176)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-25)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(23000146)
			v.EndingDatetime = types.NewDateTime(24000149)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(26000155)
			v.EndingDatetime = types.NewDateTime(27000158)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(30275)
			v.SortOrder = types.NewUInt32(31284)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(33302)
			v.SortOrder = types.NewUInt32(34311)
			return v
		}()}
		v.FilterOption = types.NewUInt32(35320)
		return v
	}(), func() player_stats_types.StatboardHistoryAggregatedQuery {
		v := player_stats_types.NewStatboardHistoryAggregatedQuery()
		v.PlayerPIDs = types.List[types.UInt32]{types.NewUInt32(38347), types.NewUInt32(39356)}
		v.BoardID = types.NewUInt32(40365)
		v.ContextIDs = types.List[types.UInt32]{types.NewUInt32(42383), types.NewUInt32(43392)}
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(45410), types.NewUInt32(46419)}
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(141), Length: types.NewUInt32(151)}
		v.PlayerStats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(50455)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v51")}
			return v
		}(), func() player_stats_types.PropertyVariant {
			v := player_stats_types.NewPropertyVariant()
			v.ID = types.NewUInt32(53482)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
			return v
		}()}
		v.DateRanges = types.List[player_stats_types.DateRange]{func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(57000248)
			v.EndingDatetime = types.NewDateTime(58000251)
			return v
		}(), func() player_stats_types.DateRange {
			v := player_stats_types.NewDateRange()
			v.StartingDatetime = types.NewDateTime(60000257)
			v.EndingDatetime = types.NewDateTime(61000260)
			return v
		}()}
		v.SortCriterias = types.List[player_stats_types.PlayerStatSortCriteria]{func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(64581)
			v.SortOrder = types.NewUInt32(65590)
			return v
		}(), func() player_stats_types.PlayerStatSortCriteria {
			v := player_stats_types.NewPlayerStatSortCriteria()
			v.SortStatID = types.NewUInt32(67608)
			v.SortOrder = types.NewUInt32(68617)
			return v
		}()}
		v.FilterOption = types.NewUInt32(69626)
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetStatboardNextPurgeDate_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getStatboardNextPurgeDate(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023))
	checkNotImplemented(t, message, nexErr)
}

func TestGetStatboardNextPurgeDate_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetStatboardNextPurgeDate = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32) (types.DateTime, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewDateTime(0), nil
	}

	message, nexErr := commonProtocol.getStatboardNextPurgeDate(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetStatboardNextPurgeDate_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetStatboardNextPurgeDate = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32) (types.DateTime, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("boardID differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("resetFrequency differs")
		}

		return types.NewDateTime(3000086), nil
	}

	message, nexErr := commonProtocol.getStatboardNextPurgeDate(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023))
	checkMessage(t, message, nexErr, ubi.MethodGetStatboardNextPurgeDate, 9, []byte{0x16, 0xc7, 0x2d, 0x00, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetStatboardNextPurgeDate = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32) (types.DateTime, *nex.Error) {
		return types.NewDateTime(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getStatboardNextPurgeDate(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsNearPlayer2_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readLeaderboardsNearPlayer2(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(15140)
		v.ContextID = types.NewUInt32(16149)
		v.ResetFrequency = types.NewUInt32(17158)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsNearPlayer2, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadLeaderboardsNearPlayer2_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadLeaderboardsNearPlayer2 = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.LeaderboardResult](), nil
	}

	message, nexErr := commonProtocol.readLeaderboardsNearPlayer2(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0), types.NewList[player_stats_types.LeaderboardQuery2]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsNearPlayer2_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadLeaderboardsNearPlayer2 = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32, got2 types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("playerPID differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("count differs")
		}

		if !(got2.Equals(types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
			v := player_stats_types.NewLeaderboardQuery2()
			v.BoardID = types.NewUInt32(5050)
			v.ContextID = types.NewUInt32(6059)
			v.ResetFrequency = types.NewUInt32(7068)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
			v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
			return v
		}(), func() player_stats_types.LeaderboardQuery2 {
			v := player_stats_types.NewLeaderboardQuery2()
			v.BoardID = types.NewUInt32(15140)
			v.ContextID = types.NewUInt32(16149)
			v.ResetFrequency = types.NewUInt32(17158)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
			v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.LeaderboardResult]{func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(26239)
			v.ContextID = types.NewUInt32(27248)
			v.ResetFrequency = types.NewUInt32(28257)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(29266)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(33302)
					v.PlayerName = types.NewString("s34")
					v.SubmittedTime = types.NewDateTime(35000182)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(38347)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v39")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(41374)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-47)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(43392)
				v.Rank = types.NewUInt32(44401)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v45")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(48437)
					v.PlayerName = types.NewString("s49")
					v.SubmittedTime = types.NewDateTime(50000227)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(53482)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(56509)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v57")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(58527)
				v.Rank = types.NewUInt32(59536)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-65)}
				return v
			}()}
			return v
		}(), func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(62563)
			v.ContextID = types.NewUInt32(63572)
			v.ResetFrequency = types.NewUInt32(64581)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(65590)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(69626)
					v.PlayerName = types.NewString("s70")
					v.SubmittedTime = types.NewDateTime(71000290)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(74671)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v75")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(77698)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-83)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(79716)
				v.Rank = types.NewUInt32(80725)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v81")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(84761)
					v.PlayerName = types.NewString("s85")
					v.SubmittedTime = types.NewDateTime(86000335)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(89806)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-95)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(92833)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v93")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(94851)
				v.Rank = types.NewUInt32(95860)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-101)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readLeaderboardsNearPlayer2(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(15140)
		v.ContextID = types.NewUInt32(16149)
		v.ResetFrequency = types.NewUInt32(17158)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsNearPlayer2, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x7f, 0x66, 0x00, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x61, 0x6e, 0x00, 0x00, 0x52, 0x72, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x16, 0x82, 0x00, 0x00, 0x04, 0x00, 0x73, 0x33, 0x34, 0x00, 0x76, 0x0f, 0x16, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xcb, 0x95, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x39, 0x00, 0x9e, 0xa1, 0x00, 0x00, 0x01, 0xd1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x80, 0xa9, 0x00, 0x00, 0x71, 0xad, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x34, 0x35, 0x00, 0x35, 0xbd, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x39, 0x00, 0x63, 0xf1, 0xfa, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xea, 0xd0, 0x00, 0x00, 0x01, 0xc5, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbd, 0xdc, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x37, 0x00, 0x9f, 0xe4, 0x00, 0x00, 0x90, 0xe8, 0x00, 0x00, 0x01, 0xbf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x63, 0xf4, 0x00, 0x00, 0x54, 0xf8, 0x00, 0x00, 0x45, 0xfc, 0x00, 0x00, 0x36, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0xfa, 0x0f, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x30, 0x00, 0xe2, 0x60, 0x3b, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xaf, 0x23, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x35, 0x00, 0x82, 0x2f, 0x01, 0x00, 0x01, 0xad, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x64, 0x37, 0x01, 0x00, 0x55, 0x3b, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x31, 0x00, 0x19, 0x4b, 0x01, 0x00, 0x04, 0x00, 0x73, 0x38, 0x35, 0x00, 0xcf, 0x42, 0x20, 0x05, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xce, 0x5e, 0x01, 0x00, 0x01, 0xa1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa1, 0x6a, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x33, 0x00, 0x83, 0x72, 0x01, 0x00, 0x74, 0x76, 0x01, 0x00, 0x01, 0x9b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadLeaderboardsNearPlayer2 = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		return types.NewList[player_stats_types.LeaderboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readLeaderboardsNearPlayer2(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(15140)
		v.ContextID = types.NewUInt32(16149)
		v.ResetFrequency = types.NewUInt32(17158)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByRank2_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readLeaderboardsByRank2(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(15140)
		v.ContextID = types.NewUInt32(16149)
		v.ResetFrequency = types.NewUInt32(17158)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByRank2, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadLeaderboardsByRank2_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadLeaderboardsByRank2 = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.LeaderboardResult](), nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByRank2(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0), types.NewList[player_stats_types.LeaderboardQuery2]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByRank2_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadLeaderboardsByRank2 = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32, got2 types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("startingRank differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("count differs")
		}

		if !(got2.Equals(types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
			v := player_stats_types.NewLeaderboardQuery2()
			v.BoardID = types.NewUInt32(5050)
			v.ContextID = types.NewUInt32(6059)
			v.ResetFrequency = types.NewUInt32(7068)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
			v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
			return v
		}(), func() player_stats_types.LeaderboardQuery2 {
			v := player_stats_types.NewLeaderboardQuery2()
			v.BoardID = types.NewUInt32(15140)
			v.ContextID = types.NewUInt32(16149)
			v.ResetFrequency = types.NewUInt32(17158)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
			v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.LeaderboardResult]{func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(26239)
			v.ContextID = types.NewUInt32(27248)
			v.ResetFrequency = types.NewUInt32(28257)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(29266)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(33302)
					v.PlayerName = types.NewString("s34")
					v.SubmittedTime = types.NewDateTime(35000182)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(38347)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v39")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(41374)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-47)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(43392)
				v.Rank = types.NewUInt32(44401)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v45")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(48437)
					v.PlayerName = types.NewString("s49")
					v.SubmittedTime = types.NewDateTime(50000227)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(53482)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-59)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(56509)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v57")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(58527)
				v.Rank = types.NewUInt32(59536)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-65)}
				return v
			}()}
			return v
		}(), func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(62563)
			v.ContextID = types.NewUInt32(63572)
			v.ResetFrequency = types.NewUInt32(64581)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(65590)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(69626)
					v.PlayerName = types.NewString("s70")
					v.SubmittedTime = types.NewDateTime(71000290)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(74671)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v75")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(77698)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-83)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(79716)
				v.Rank = types.NewUInt32(80725)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v81")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(84761)
					v.PlayerName = types.NewString("s85")
					v.SubmittedTime = types.NewDateTime(86000335)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(89806)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-95)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(92833)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v93")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(94851)
				v.Rank = types.NewUInt32(95860)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-101)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByRank2(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(15140)
		v.ContextID = types.NewUInt32(16149)
		v.ResetFrequency = types.NewUInt32(17158)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByRank2, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x7f, 0x66, 0x00, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x61, 0x6e, 0x00, 0x00, 0x52, 0x72, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x16, 0x82, 0x00, 0x00, 0x04, 0x00, 0x73, 0x33, 0x34, 0x00, 0x76, 0x0f, 0x16, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xcb, 0x95, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x39, 0x00, 0x9e, 0xa1, 0x00, 0x00, 0x01, 0xd1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x80, 0xa9, 0x00, 0x00, 0x71, 0xad, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x34, 0x35, 0x00, 0x35, 0xbd, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x39, 0x00, 0x63, 0xf1, 0xfa, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xea, 0xd0, 0x00, 0x00, 0x01, 0xc5, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbd, 0xdc, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x37, 0x00, 0x9f, 0xe4, 0x00, 0x00, 0x90, 0xe8, 0x00, 0x00, 0x01, 0xbf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x63, 0xf4, 0x00, 0x00, 0x54, 0xf8, 0x00, 0x00, 0x45, 0xfc, 0x00, 0x00, 0x36, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0xfa, 0x0f, 0x01, 0x00, 0x04, 0x00, 0x73, 0x37, 0x30, 0x00, 0xe2, 0x60, 0x3b, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xaf, 0x23, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x35, 0x00, 0x82, 0x2f, 0x01, 0x00, 0x01, 0xad, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x64, 0x37, 0x01, 0x00, 0x55, 0x3b, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x38, 0x31, 0x00, 0x19, 0x4b, 0x01, 0x00, 0x04, 0x00, 0x73, 0x38, 0x35, 0x00, 0xcf, 0x42, 0x20, 0x05, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xce, 0x5e, 0x01, 0x00, 0x01, 0xa1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa1, 0x6a, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x33, 0x00, 0x83, 0x72, 0x01, 0x00, 0x74, 0x76, 0x01, 0x00, 0x01, 0x9b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadLeaderboardsByRank2 = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		return types.NewList[player_stats_types.LeaderboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readLeaderboardsByRank2(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(5050)
		v.ContextID = types.NewUInt32(6059)
		v.ResetFrequency = types.NewUInt32(7068)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(9086), types.NewUInt32(10095)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(15140)
		v.ContextID = types.NewUInt32(16149)
		v.ResetFrequency = types.NewUInt32(17158)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(19176), types.NewUInt32(20185)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(22203), types.NewUInt32(23212)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByPlayers2_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readLeaderboardsByPlayers2(nil, fakePacket{}, 7, types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(3032)
		v.ContextID = types.NewUInt32(4041)
		v.ResetFrequency = types.NewUInt32(5050)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(7068), types.NewUInt32(8077)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(13122)
		v.ContextID = types.NewUInt32(14131)
		v.ResetFrequency = types.NewUInt32(15140)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(20185), types.NewUInt32(21194)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByPlayers2, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadLeaderboardsByPlayers2_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadLeaderboardsByPlayers2 = func(packet nex.PacketInterface, _ types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.LeaderboardResult](), nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByPlayers2(errors.New("bad"), fakePacket{}, 7, types.NewList[player_stats_types.LeaderboardQuery2]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadLeaderboardsByPlayers2_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadLeaderboardsByPlayers2 = func(packet nex.PacketInterface, got0 types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
			v := player_stats_types.NewLeaderboardQuery2()
			v.BoardID = types.NewUInt32(3032)
			v.ContextID = types.NewUInt32(4041)
			v.ResetFrequency = types.NewUInt32(5050)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(7068), types.NewUInt32(8077)}
			v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
			return v
		}(), func() player_stats_types.LeaderboardQuery2 {
			v := player_stats_types.NewLeaderboardQuery2()
			v.BoardID = types.NewUInt32(13122)
			v.ContextID = types.NewUInt32(14131)
			v.ResetFrequency = types.NewUInt32(15140)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
			v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(20185), types.NewUInt32(21194)}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.LeaderboardResult]{func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(24221)
			v.ContextID = types.NewUInt32(25230)
			v.ResetFrequency = types.NewUInt32(26239)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(27248)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(31284)
					v.PlayerName = types.NewString("s32")
					v.SubmittedTime = types.NewDateTime(33000176)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(36329)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v37")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(39356)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-45)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(41374)
				v.Rank = types.NewUInt32(42383)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v43")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(46419)
					v.PlayerName = types.NewString("s47")
					v.SubmittedTime = types.NewDateTime(48000221)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(51464)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-57)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(54491)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v55")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(56509)
				v.Rank = types.NewUInt32(57518)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-63)}
				return v
			}()}
			return v
		}(), func() player_stats_types.LeaderboardResult {
			v := player_stats_types.NewLeaderboardResult()
			v.BoardID = types.NewUInt32(60545)
			v.ContextID = types.NewUInt32(61554)
			v.ResetFrequency = types.NewUInt32(62563)
			v.LeaderboardTotalPlayerCount = types.NewUInt32(63572)
			v.PlayerRanks = types.List[player_stats_types.PlayerRank]{func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(67608)
					v.PlayerName = types.NewString("s68")
					v.SubmittedTime = types.NewDateTime(69000284)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(72653)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v73")}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(75680)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-81)}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(77698)
				v.Rank = types.NewUInt32(78707)
				v.Score = types.Variant{TypeID: 4, Type: types.NewString("v79")}
				return v
			}(), func() player_stats_types.PlayerRank {
				v := player_stats_types.NewPlayerRank()
				v.PlayerStatSet = func() player_stats_types.PlayerStatSet {
					v := player_stats_types.NewPlayerStatSet()
					v.PlayerPID = types.NewUInt32(82743)
					v.PlayerName = types.NewString("s83")
					v.SubmittedTime = types.NewDateTime(84000329)
					v.Stats = types.List[player_stats_types.PropertyVariant]{func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(87788)
						v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-93)}
						return v
					}(), func() player_stats_types.PropertyVariant {
						v := player_stats_types.NewPropertyVariant()
						v.ID = types.NewUInt32(90815)
						v.Value = types.Variant{TypeID: 4, Type: types.NewString("v91")}
						return v
					}()}
					return v
				}()
				v.RankStatus = types.NewUInt32(92833)
				v.Rank = types.NewUInt32(93842)
				v.Score = types.Variant{TypeID: 1, Type: types.NewInt64(-99)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readLeaderboardsByPlayers2(nil, fakePacket{}, 9, types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(3032)
		v.ContextID = types.NewUInt32(4041)
		v.ResetFrequency = types.NewUInt32(5050)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(7068), types.NewUInt32(8077)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(13122)
		v.ContextID = types.NewUInt32(14131)
		v.ResetFrequency = types.NewUInt32(15140)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(20185), types.NewUInt32(21194)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadLeaderboardsByPlayers2, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x9d, 0x5e, 0x00, 0x00, 0x8e, 0x62, 0x00, 0x00, 0x7f, 0x66, 0x00, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x34, 0x7a, 0x00, 0x00, 0x04, 0x00, 0x73, 0x33, 0x32, 0x00, 0xf0, 0x8a, 0xf7, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xe9, 0x8d, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x37, 0x00, 0xbc, 0x99, 0x00, 0x00, 0x01, 0xd3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x9e, 0xa1, 0x00, 0x00, 0x8f, 0xa5, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x34, 0x33, 0x00, 0x53, 0xb5, 0x00, 0x00, 0x04, 0x00, 0x73, 0x34, 0x37, 0x00, 0xdd, 0x6c, 0xdc, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0xc9, 0x00, 0x00, 0x01, 0xc7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xdb, 0xd4, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x35, 0x35, 0x00, 0xbd, 0xdc, 0x00, 0x00, 0xae, 0xe0, 0x00, 0x00, 0x01, 0xc1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x81, 0xec, 0x00, 0x00, 0x72, 0xf0, 0x00, 0x00, 0x63, 0xf4, 0x00, 0x00, 0x54, 0xf8, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x18, 0x08, 0x01, 0x00, 0x04, 0x00, 0x73, 0x36, 0x38, 0x00, 0x5c, 0xdc, 0x1c, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xcd, 0x1b, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x33, 0x00, 0xa0, 0x27, 0x01, 0x00, 0x01, 0xaf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x82, 0x2f, 0x01, 0x00, 0x73, 0x33, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x37, 0x39, 0x00, 0x37, 0x43, 0x01, 0x00, 0x04, 0x00, 0x73, 0x38, 0x33, 0x00, 0x49, 0xbe, 0x01, 0x05, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xec, 0x56, 0x01, 0x00, 0x01, 0xa3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbf, 0x62, 0x01, 0x00, 0x04, 0x04, 0x00, 0x76, 0x39, 0x31, 0x00, 0xa1, 0x6a, 0x01, 0x00, 0x92, 0x6e, 0x01, 0x00, 0x01, 0x9d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadLeaderboardsByPlayers2 = func(packet nex.PacketInterface, _ types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error) {
		return types.NewList[player_stats_types.LeaderboardResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readLeaderboardsByPlayers2(nil, fakePacket{}, 9, types.List[player_stats_types.LeaderboardQuery2]{func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(3032)
		v.ContextID = types.NewUInt32(4041)
		v.ResetFrequency = types.NewUInt32(5050)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(7068), types.NewUInt32(8077)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(10095), types.NewUInt32(11104)}
		return v
	}(), func() player_stats_types.LeaderboardQuery2 {
		v := player_stats_types.NewLeaderboardQuery2()
		v.BoardID = types.NewUInt32(13122)
		v.ContextID = types.NewUInt32(14131)
		v.ResetFrequency = types.NewUInt32(15140)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
		v.EstimatedPIDs = types.List[types.UInt32]{types.NewUInt32(20185), types.NewUInt32(21194)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReadPopulationStats_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.readPopulationStats(nil, fakePacket{}, 7, types.List[player_stats_types.PopulationStatQuery]{func() player_stats_types.PopulationStatQuery {
		v := player_stats_types.NewPopulationStatQuery()
		v.BoardID = types.NewUInt32(3032)
		v.ContextID = types.NewUInt32(4041)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
		return v
	}(), func() player_stats_types.PopulationStatQuery {
		v := player_stats_types.NewPopulationStatQuery()
		v.BoardID = types.NewUInt32(9086)
		v.ContextID = types.NewUInt32(10095)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadPopulationStats, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestReadPopulationStats_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReadPopulationStats = func(packet nex.PacketInterface, _ types.List[player_stats_types.PopulationStatQuery]) (types.List[player_stats_types.PopulationStatResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[player_stats_types.PopulationStatResult](), nil
	}

	message, nexErr := commonProtocol.readPopulationStats(errors.New("bad"), fakePacket{}, 7, types.NewList[player_stats_types.PopulationStatQuery]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReadPopulationStats_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReadPopulationStats = func(packet nex.PacketInterface, got0 types.List[player_stats_types.PopulationStatQuery]) (types.List[player_stats_types.PopulationStatResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[player_stats_types.PopulationStatQuery]{func() player_stats_types.PopulationStatQuery {
			v := player_stats_types.NewPopulationStatQuery()
			v.BoardID = types.NewUInt32(3032)
			v.ContextID = types.NewUInt32(4041)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
			return v
		}(), func() player_stats_types.PopulationStatQuery {
			v := player_stats_types.NewPopulationStatQuery()
			v.BoardID = types.NewUInt32(9086)
			v.ContextID = types.NewUInt32(10095)
			v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
			return v
		}()})) {
			t.Errorf("queries differs")
		}

		return types.List[player_stats_types.PopulationStatResult]{func() player_stats_types.PopulationStatResult {
			v := player_stats_types.NewPopulationStatResult()
			v.BoardID = types.NewUInt32(16149)
			v.ContextID = types.NewUInt32(17158)
			v.StatID = types.NewUInt32(18167)
			v.Sum = types.NewDouble(19.5)
			v.Average = types.NewDouble(20.5)
			v.StandardDeviation = types.NewDouble(21.5)
			return v
		}(), func() player_stats_types.PopulationStatResult {
			v := player_stats_types.NewPopulationStatResult()
			v.BoardID = types.NewUInt32(23212)
			v.ContextID = types.NewUInt32(24221)
			v.StatID = types.NewUInt32(25230)
			v.Sum = types.NewDouble(26.5)
			v.Average = types.NewDouble(27.5)
			v.StandardDeviation = types.NewDouble(28.5)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.readPopulationStats(nil, fakePacket{}, 9, types.List[player_stats_types.PopulationStatQuery]{func() player_stats_types.PopulationStatQuery {
		v := player_stats_types.NewPopulationStatQuery()
		v.BoardID = types.NewUInt32(3032)
		v.ContextID = types.NewUInt32(4041)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
		return v
	}(), func() player_stats_types.PopulationStatQuery {
		v := player_stats_types.NewPopulationStatQuery()
		v.BoardID = types.NewUInt32(9086)
		v.ContextID = types.NewUInt32(10095)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReadPopulationStats, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x15, 0x3f, 0x00, 0x00, 0x06, 0x43, 0x00, 0x00, 0xf7, 0x46, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x33, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x34, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x35, 0x40, 0xac, 0x5a, 0x00, 0x00, 0x9d, 0x5e, 0x00, 0x00, 0x8e, 0x62, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x3a, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x3b, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x3c, 0x40})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReadPopulationStats = func(packet nex.PacketInterface, _ types.List[player_stats_types.PopulationStatQuery]) (types.List[player_stats_types.PopulationStatResult], *nex.Error) {
		return types.NewList[player_stats_types.PopulationStatResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.readPopulationStats(nil, fakePacket{}, 9, types.List[player_stats_types.PopulationStatQuery]{func() player_stats_types.PopulationStatQuery {
		v := player_stats_types.NewPopulationStatQuery()
		v.BoardID = types.NewUInt32(3032)
		v.ContextID = types.NewUInt32(4041)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
		return v
	}(), func() player_stats_types.PopulationStatQuery {
		v := player_stats_types.NewPopulationStatQuery()
		v.BoardID = types.NewUInt32(9086)
		v.ContextID = types.NewUInt32(10095)
		v.StatIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
