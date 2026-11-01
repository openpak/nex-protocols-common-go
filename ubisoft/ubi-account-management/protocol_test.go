package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
	"testing"
)

func init() {
	types.RegisterObjectHolderType(types.NewData())
}

func TestBuildCreateAccountResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "050073323132000500733231330005007332313400d810000000010005007332323000dd77665544332211de10000005007332323300010005007332323600050073323237000500733232380002000000e71000000500733233320005007332333300eb100000050073323336000500733233370002000000f010000005007332343100f310000005007332343400"}, {true, "00840000000500733231320005007332313300050073323134000007000000d810000000010005007332323000dd77665544332211de100000050073323233000100050073323236000500733232370005007332323800020000000012000000e710000005007332333200050073323333000012000000eb100000050073323336000500733233370002000000000b000000f010000005007332343100000b000000f310000005007332343400"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildCreateAccountResponse(out, func() ubi_account_management_types.UbiAccount {
			v := ubi_account_management_types.NewUbiAccount()
			v.UbiAccountID = types.NewString("s212")
			v.Username = types.NewString("s213")
			v.Password = types.NewString("s214")
			v.Status = func() ubi_account_management_types.UbiAccountStatus {
				v := ubi_account_management_types.NewUbiAccountStatus()
				v.BasicStatus = types.NewUInt32(4312)
				v.MissingRequiredInformations = types.NewBool(false)
				v.RecoveringPassword = types.NewBool(true)
				v.PendingDeactivation = types.NewBool(false)
				return v
			}()
			v.Email = types.NewString("s220")
			v.DateOfBirth = types.NewDateTime(1234605616436508637)
			v.Gender = types.NewUInt32(4318)
			v.CountryCode = types.NewString("s223")
			v.OptIn = types.NewBool(true)
			v.ThirdPartyOptIn = types.NewBool(false)
			v.FirstName = types.NewString("s226")
			v.LastName = types.NewString("s227")
			v.PreferredLanguage = types.NewString("s228")
			v.ExternalAccounts = types.List[ubi_account_management_types.ExternalAccount]{func() ubi_account_management_types.ExternalAccount {
				v := ubi_account_management_types.NewExternalAccount()
				v.AccountType = types.NewUInt32(4327)
				v.ID = types.NewString("s232")
				v.Username = types.NewString("s233")
				return v
			}(), func() ubi_account_management_types.ExternalAccount {
				v := ubi_account_management_types.NewExternalAccount()
				v.AccountType = types.NewUInt32(4331)
				v.ID = types.NewString("s236")
				v.Username = types.NewString("s237")
				return v
			}()}
			return v
		}(), types.List[ubi_account_management_types.ValidationFailureReason]{func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4336)
			v.Description = types.NewString("s241")
			return v
		}(), func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4339)
			v.Description = types.NewString("s244")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestCreateAccountHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.createAccount(nil, nil, 1, ubi_account_management_types.NewUbiAccount())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.createAccount(errTest, nil, 1, ubi_account_management_types.NewUbiAccount())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildUpdateAccountResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "050073323436000500733234370005007332343800fa10000000010005007332353400ff7766554433221100110000050073323537000100050073323630000500733236310005007332363200020000000911000005007332363600050073323637000d11000005007332373000050073323731000200000012110000050073323735001511000005007332373800"}, {true, "00840000000500733234360005007332343700050073323438000007000000fa10000000010005007332353400ff776655443322110011000005007332353700010005007332363000050073323631000500733236320002000000001200000009110000050073323636000500733236370000120000000d110000050073323730000500733237310002000000000b0000001211000005007332373500000b0000001511000005007332373800"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildUpdateAccountResponse(out, func() ubi_account_management_types.UbiAccount {
			v := ubi_account_management_types.NewUbiAccount()
			v.UbiAccountID = types.NewString("s246")
			v.Username = types.NewString("s247")
			v.Password = types.NewString("s248")
			v.Status = func() ubi_account_management_types.UbiAccountStatus {
				v := ubi_account_management_types.NewUbiAccountStatus()
				v.BasicStatus = types.NewUInt32(4346)
				v.MissingRequiredInformations = types.NewBool(false)
				v.RecoveringPassword = types.NewBool(true)
				v.PendingDeactivation = types.NewBool(false)
				return v
			}()
			v.Email = types.NewString("s254")
			v.DateOfBirth = types.NewDateTime(1234605616436508671)
			v.Gender = types.NewUInt32(4352)
			v.CountryCode = types.NewString("s257")
			v.OptIn = types.NewBool(true)
			v.ThirdPartyOptIn = types.NewBool(false)
			v.FirstName = types.NewString("s260")
			v.LastName = types.NewString("s261")
			v.PreferredLanguage = types.NewString("s262")
			v.ExternalAccounts = types.List[ubi_account_management_types.ExternalAccount]{func() ubi_account_management_types.ExternalAccount {
				v := ubi_account_management_types.NewExternalAccount()
				v.AccountType = types.NewUInt32(4361)
				v.ID = types.NewString("s266")
				v.Username = types.NewString("s267")
				return v
			}(), func() ubi_account_management_types.ExternalAccount {
				v := ubi_account_management_types.NewExternalAccount()
				v.AccountType = types.NewUInt32(4365)
				v.ID = types.NewString("s270")
				v.Username = types.NewString("s271")
				return v
			}()}
			return v
		}(), types.List[ubi_account_management_types.ValidationFailureReason]{func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4370)
			v.Description = types.NewString("s275")
			return v
		}(), func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4373)
			v.Description = types.NewString("s278")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestUpdateAccountHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.updateAccount(nil, nil, 1, ubi_account_management_types.NewUbiAccount())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.updateAccount(errTest, nil, 1, ubi_account_management_types.NewUbiAccount())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetAccountResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0500733238300005007332383100050073323832001c11000000010005007332383800217866554433221122110000050073323931000100050073323934000500733239350005007332393600020000002b11000005007333303000050073333031002f110000050073333034000500733330350001"}, {true, "008400000005007332383000050073323831000500733238320000070000001c110000000100050073323838002178665544332211221100000500733239310001000500733239340005007332393500050073323936000200000000120000002b110000050073333030000500733330310000120000002f110000050073333034000500733330350001"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetAccountResponse(out, func() ubi_account_management_types.UbiAccount {
			v := ubi_account_management_types.NewUbiAccount()
			v.UbiAccountID = types.NewString("s280")
			v.Username = types.NewString("s281")
			v.Password = types.NewString("s282")
			v.Status = func() ubi_account_management_types.UbiAccountStatus {
				v := ubi_account_management_types.NewUbiAccountStatus()
				v.BasicStatus = types.NewUInt32(4380)
				v.MissingRequiredInformations = types.NewBool(false)
				v.RecoveringPassword = types.NewBool(true)
				v.PendingDeactivation = types.NewBool(false)
				return v
			}()
			v.Email = types.NewString("s288")
			v.DateOfBirth = types.NewDateTime(1234605616436508705)
			v.Gender = types.NewUInt32(4386)
			v.CountryCode = types.NewString("s291")
			v.OptIn = types.NewBool(true)
			v.ThirdPartyOptIn = types.NewBool(false)
			v.FirstName = types.NewString("s294")
			v.LastName = types.NewString("s295")
			v.PreferredLanguage = types.NewString("s296")
			v.ExternalAccounts = types.List[ubi_account_management_types.ExternalAccount]{func() ubi_account_management_types.ExternalAccount {
				v := ubi_account_management_types.NewExternalAccount()
				v.AccountType = types.NewUInt32(4395)
				v.ID = types.NewString("s300")
				v.Username = types.NewString("s301")
				return v
			}(), func() ubi_account_management_types.ExternalAccount {
				v := ubi_account_management_types.NewExternalAccount()
				v.AccountType = types.NewUInt32(4399)
				v.ID = types.NewString("s304")
				v.Username = types.NewString("s305")
				return v
			}()}
			return v
		}(), types.NewBool(true))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetAccountHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getAccount(nil, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getAccount(errTest, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLinkAccountResponse(t *testing.T) {
	out := newOut(false)
	buildLinkAccountResponse(out)
	checkBytes(t, out.Bytes(), "")
}

func TestLinkAccountHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.linkAccount(nil, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.linkAccount(errTest, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetTOSResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "050073333038000500733330390005007333313000"}, {true, "0015000000050073333038000500733330390005007333313000"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetTOSResponse(out, func() ubi_account_management_types.TOS {
			v := ubi_account_management_types.NewTOS()
			v.LocaleCode = types.NewString("s308")
			v.Content = types.NewString("s309")
			v.StoringInfoQuestion = types.NewString("s310")
			return v
		}())
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetTOSHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getTOS(nil, nil, 1, types.NewString(""), types.NewString(""), types.NewBool(false))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getTOS(errTest, nil, 1, types.NewString(""), types.NewString(""), types.NewBool(false))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildValidateUsernameResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "020000000500733331330005007333313400020000003d110000050073333138004011000005007333323100"}, {true, "003600000002000000050073333133000500733331340002000000000b0000003d11000005007333313800000b0000004011000005007333323100"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildValidateUsernameResponse(out, func() ubi_account_management_types.UsernameValidation {
			v := ubi_account_management_types.NewUsernameValidation()
			v.Suggestions = types.List[types.String]{types.NewString("s313"), types.NewString("s314")}
			v.Reasons = types.List[ubi_account_management_types.ValidationFailureReason]{func() ubi_account_management_types.ValidationFailureReason {
				v := ubi_account_management_types.NewValidationFailureReason()
				v.ValidationID = types.NewUInt32(4413)
				v.Description = types.NewString("s318")
				return v
			}(), func() ubi_account_management_types.ValidationFailureReason {
				v := ubi_account_management_types.NewValidationFailureReason()
				v.ValidationID = types.NewUInt32(4416)
				v.Description = types.NewString("s321")
				return v
			}()}
			return v
		}())
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestValidateUsernameHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.validateUsername(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.validateUsername(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildValidatePasswordResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000044110000050073333235004711000005007333323800"}, {true, "02000000000b0000004411000005007333323500000b0000004711000005007333323800"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildValidatePasswordResponse(out, types.List[ubi_account_management_types.ValidationFailureReason]{func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4420)
			v.Description = types.NewString("s325")
			return v
		}(), func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4423)
			v.Description = types.NewString("s328")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestValidatePasswordHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.validatePassword(nil, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.validatePassword(errTest, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildValidateEmailResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "020000004b110000050073333332004e11000005007333333500"}, {true, "02000000000b0000004b11000005007333333200000b0000004e11000005007333333500"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildValidateEmailResponse(out, types.List[ubi_account_management_types.ValidationFailureReason]{func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4427)
			v.Description = types.NewString("s332")
			return v
		}(), func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4430)
			v.Description = types.NewString("s335")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestValidateEmailHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.validateEmail(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.validateEmail(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetCountryListResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000005007333333800050073333339000500733334310005007333343200"}, {true, "02000000000e0000000500733333380005007333333900000e0000000500733334310005007333343200"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetCountryListResponse(out, types.List[ubi_account_management_types.Country]{func() ubi_account_management_types.Country {
			v := ubi_account_management_types.NewCountry()
			v.Code = types.NewString("s338")
			v.Name = types.NewString("s339")
			return v
		}(), func() ubi_account_management_types.Country {
			v := ubi_account_management_types.NewCountry()
			v.Code = types.NewString("s341")
			v.Name = types.NewString("s342")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetCountryListHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getCountryList(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getCountryList(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildForgetPasswordResponse(t *testing.T) {
	out := newOut(false)
	buildForgetPasswordResponse(out)
	checkBytes(t, out.Bytes(), "")
}

func TestForgetPasswordHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.forgetPassword(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.forgetPassword(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupPrincipalIdsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733334340059110000"}, {true, "010000000500733334340059110000"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupPrincipalIdsResponse(out, types.Map[types.String, types.UInt32]{types.NewString("s344"): types.NewUInt32(4441)})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupPrincipalIdsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupPrincipalIds(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupPrincipalIds(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUbiAccountIDsByPidsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000005b11000005007333343800"}, {true, "010000005b11000005007333343800"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUbiAccountIDsByPidsResponse(out, types.Map[types.UInt32, types.String]{types.NewUInt32(4443): types.NewString("s348")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUbiAccountIDsByPidsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUbiAccountIDsByPids(nil, nil, 1, types.NewList[types.UInt32]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUbiAccountIDsByPids(errTest, nil, 1, types.NewList[types.UInt32]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUbiAccountIDsByUsernamesResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733335300005007333353100"}, {true, "010000000500733335300005007333353100"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUbiAccountIDsByUsernamesResponse(out, types.Map[types.String, types.String]{types.NewString("s350"): types.NewString("s351")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUbiAccountIDsByUsernamesHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUbiAccountIDsByUsernames(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUbiAccountIDsByUsernames(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUsernamesByUbiAccountIDsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733335330005007333353400"}, {true, "010000000500733335330005007333353400"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUsernamesByUbiAccountIDsResponse(out, types.Map[types.String, types.String]{types.NewString("s353"): types.NewString("s354")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUsernamesByUbiAccountIDsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUsernamesByUbiAccountIDs(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUsernamesByUbiAccountIDs(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUbiAccountIDsByUsernameSubStringResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733335360005007333353700"}, {true, "010000000500733335360005007333353700"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUbiAccountIDsByUsernameSubStringResponse(out, types.Map[types.String, types.String]{types.NewString("s356"): types.NewString("s357")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUbiAccountIDsByUsernameSubStringHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUbiAccountIDsByUsernameSubString(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUbiAccountIDsByUsernameSubString(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildUserHasPlayedResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733335390001"}, {true, "010000000500733335390001"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildUserHasPlayedResponse(out, types.Map[types.String, types.Bool]{types.NewString("s359"): types.NewBool(true)})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestUserHasPlayedHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.userHasPlayed(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.userHasPlayed(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildIsUserPlayingResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733336320000"}, {true, "010000000500733336320000"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildIsUserPlayingResponse(out, types.Map[types.String, types.Bool]{types.NewString("s362"): types.NewBool(false)})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestIsUserPlayingHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.isUserPlaying(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.isUserPlaying(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUbiAccountIDsByUsernamesGlobalResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733336350005007333363600"}, {true, "010000000500733336350005007333363600"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUbiAccountIDsByUsernamesGlobalResponse(out, types.Map[types.String, types.String]{types.NewString("s365"): types.NewString("s366")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUbiAccountIDsByUsernamesGlobalHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUbiAccountIDsByUsernamesGlobal(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUbiAccountIDsByUsernamesGlobal(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUbiAccountIDsByEmailsGlobalResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733336380005007333363900"}, {true, "010000000500733336380005007333363900"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUbiAccountIDsByEmailsGlobalResponse(out, types.Map[types.String, types.String]{types.NewString("s368"): types.NewString("s369")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUbiAccountIDsByEmailsGlobalHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUbiAccountIDsByEmailsGlobal(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUbiAccountIDsByEmailsGlobal(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUsernamesByUbiAccountIDsGlobalResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733337310005007333373200"}, {true, "010000000500733337310005007333373200"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUsernamesByUbiAccountIDsGlobalResponse(out, types.Map[types.String, types.String]{types.NewString("s371"): types.NewString("s372")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUsernamesByUbiAccountIDsGlobalHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUsernamesByUbiAccountIDsGlobal(nil, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUsernamesByUbiAccountIDsGlobal(errTest, nil, 1, types.NewList[types.String]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetTOSExResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0500733337340002000000050073333736000500733337370005007333373800"}, {true, "00200000000500733337340002000000050073333736000500733337370005007333373800"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetTOSExResponse(out, func() ubi_account_management_types.TOSEx {
			v := ubi_account_management_types.NewTOSEx()
			v.LocaleCode = types.NewString("s374")
			v.Contents = types.List[types.String]{types.NewString("s376"), types.NewString("s377")}
			v.StoringInfoQuestion = types.NewString("s378")
			return v
		}())
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetTOSExHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getTOSEx(nil, nil, 1, types.NewString(""), types.NewString(""), types.NewBool(false))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getTOSEx(errTest, nil, 1, types.NewString(""), types.NewString(""), types.NewBool(false))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildHasAcceptedLatestTOSResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "00020000007e110000050073333833008111000005007333383600"}, {true, "0002000000000b0000007e11000005007333383300000b0000008111000005007333383600"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildHasAcceptedLatestTOSResponse(out, types.NewBool(false), types.List[ubi_account_management_types.ValidationFailureReason]{func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4478)
			v.Description = types.NewString("s383")
			return v
		}(), func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4481)
			v.Description = types.NewString("s386")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestHasAcceptedLatestTOSHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.hasAcceptedLatestTOS(nil, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.hasAcceptedLatestTOS(errTest, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildAcceptLatestTOSResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000085110000050073333930008811000005007333393300"}, {true, "02000000000b0000008511000005007333393000000b0000008811000005007333393300"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildAcceptLatestTOSResponse(out, types.List[ubi_account_management_types.ValidationFailureReason]{func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4485)
			v.Description = types.NewString("s390")
			return v
		}(), func() ubi_account_management_types.ValidationFailureReason {
			v := ubi_account_management_types.NewValidationFailureReason()
			v.ValidationID = types.NewUInt32(4488)
			v.Description = types.NewString("s393")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestAcceptLatestTOSHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.acceptLatestTOS(nil, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.acceptLatestTOS(errTest, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}
