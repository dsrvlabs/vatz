package dispatcher

import (
	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageHandler(t *testing.T) {
	preStat := tp.StateFlag{
		State:    pb.STATE_FAILURE,
		Severity: pb.SEVERITY_CRITICAL,
	}
	notifyInfoTPOn := tp.NotifyInfo{
		Plugin:     "SamplePlugin",
		Method:     "IsSuccess",
		State:      pb.STATE_SUCCESS,
		Severity:   pb.SEVERITY_WARNING,
		ExecuteMsg: "ExecuteMsg",
	}

	notifyInfoTPOff := tp.NotifyInfo{
		Plugin:     "SamplePlugin",
		Method:     "IsSuccess",
		State:      pb.STATE_SUCCESS,
		Severity:   pb.SEVERITY_INFO,
		ExecuteMsg: "ExecuteMsg",
	}

	notifyInfoTPHang := tp.NotifyInfo{
		Plugin:     "SamplePlugin",
		Method:     "IsSuccess",
		State:      pb.STATE_FAILURE,
		Severity:   pb.SEVERITY_CRITICAL,
		ExecuteMsg: "ExecuteMsg",
	}

	no1, reminderSt1, deliverMessage1 := messageHandler(true, preStat, notifyInfoTPOn)
	no2, reminderSt2, deliverMessage2 := messageHandler(false, preStat, notifyInfoTPOff)
	no3, reminderSt3, deliverMessage3 := messageHandler(false, preStat, notifyInfoTPHang)

	assert.True(t, true == no1)
	assert.Equal(t, tp.ON, reminderSt1)
	assert.Equal(t, tp.ReqMsg{
		FuncName:     "IsSuccess",
		State:        pb.STATE_SUCCESS,
		Msg:          "ExecuteMsg",
		Severity:     pb.SEVERITY_WARNING,
		ResourceType: "SamplePlugin",
	}, deliverMessage1)

	assert.False(t, false == no2)
	assert.Equal(t, tp.OFF, reminderSt2)
	assert.Equal(t, tp.ReqMsg{
		FuncName:     "IsSuccess",
		State:        pb.STATE_SUCCESS,
		Msg:          "ExecuteMsg",
		Severity:     pb.SEVERITY_INFO,
		ResourceType: "SamplePlugin",
	}, deliverMessage2)

	assert.False(t, no3)
	assert.Equal(t, tp.HANG, reminderSt3)
	assert.Equal(t, tp.ReqMsg{
		FuncName:     "IsSuccess",
		State:        pb.STATE_FAILURE,
		Msg:          "ExecuteMsg",
		Severity:     pb.SEVERITY_CRITICAL,
		ResourceType: "SamplePlugin",
	}, deliverMessage3)
}
