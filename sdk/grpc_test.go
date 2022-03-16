package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type mockFuncs struct {
	mock.Mock
}

func (m *mockFuncs) DummyCall1(info, option map[string]*structpb.Value) error {
	ret := m.Called(info, option)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]*structpb.Value, map[string]*structpb.Value) error); ok {
		r0 = rf(info, option)
	} else {
		if ret.Get(0) == nil {
			r0 = nil
		} else {
			r0 = ret.Get(0).(error)
		}
	}

	return r0
}

func TestRegister(t *testing.T) {
	tests := []struct {
		Funcs     []func(in, opt map[string]*structpb.Value) error
		ExpectErr error
	}{
		{
			Funcs: []func(in, opt map[string]*structpb.Value) error{
				callbackFunc,
			},
			ExpectErr: nil,
		},

		{
			Funcs: []func(in, opt map[string]*structpb.Value) error{
				callbackFunc,
				callbackFunc,
				callbackFunc,
				callbackFunc,
				callbackFunc,
				callbackFunc,
			},
			ExpectErr: ErrRegisterMaxLimit,
		},
	}

	for _, test := range tests {
		p := plugin{}

		var err error
		for _, f := range test.Funcs {
			err = p.Register(f)
		}

		if test.ExpectErr == nil {
			assert.Equal(t, len(test.Funcs), len(p.grpc.callbacks))
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.ExpectErr, err)
		}

	}

}

func callbackFunc(in, opt map[string]*structpb.Value) error {
	return nil
}
