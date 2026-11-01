package matchmaking

import "testing"

type fakeProto struct {
	Set   func()
	Unset func()
}

func TestAlreadyHandled(t *testing.T) {
	p := &fakeProto{Set: func() {}}
	if !alreadyHandled(p, "Set") {
		t.Fatal("Set should be reported as handled")
	}
	if alreadyHandled(p, "Unset") || alreadyHandled(p, "Missing") || alreadyHandled(nil, "Set") {
		t.Fatal("unset, missing and nil must not be reported as handled")
	}
}
