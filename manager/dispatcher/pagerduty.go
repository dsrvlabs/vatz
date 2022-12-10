package dispatcher

import (
	"context"
	"fmt"
	pd "github.com/PagerDuty/go-pagerduty"
	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

const SUCCESS = "success"

type pagerdutyMSGEvent struct {
	flag  string
	deKey string
}

type pagerduty struct {
	host       string
	channel    tp.Channel
	secret     string
	pagerEntry sync.Map
}

func (p *pagerduty) SetDispatcher(firstRunMsg bool, preStat tp.StateFlag, notifyInfo tp.NotifyInfo) error {
	reqToNotify, _, deliverMessage := messageHandler(firstRunMsg, preStat, notifyInfo)
	if reqToNotify {
		p.SendNotification(deliverMessage)
	}
	return nil
}

func (p *pagerduty) SendNotification(msg tp.ReqMsg) error {

	var (
		pagerdutySeverity string
		emoji             string
		methodName        = msg.Option["pUnique"].(string)
	)
	/*
		Pagerduty severity Allowed values:
		- critical
		- warning
		- error
		- info
	*/
	switch {
	case msg.Severity == pb.SEVERITY_INFO:
		pagerdutySeverity = "info"
		emoji = emojiCheck
	case msg.Severity == pb.SEVERITY_CRITICAL:
		pagerdutySeverity = "critical"
		emoji = "‼️"
	case msg.Severity == pb.SEVERITY_WARNING:
		pagerdutySeverity = "warning"
		emoji = "❗"
	default:
		emoji = emojiER
		pagerdutySeverity = "error"
	}

	v2EventPayload := &pd.V2Payload{
		Source:    fmt.Sprintf(`(%s)`, p.host),
		Component: msg.ResourceType,
		Severity:  pagerdutySeverity,
		Summary: fmt.Sprintf(`
%s %s %s 
(%s)
Plugin Name: %s 
%s`, emoji, msg.Severity.String(), emoji, p.host, msg.ResourceType, msg.Msg)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if pb.STATE_SUCCESS != msg.State || pb.SEVERITY_INFO != msg.Severity {
		resp, err := pd.ManageEventWithContext(ctx, pd.V2Event{RoutingKey: p.secret, Action: "trigger", Payload: v2EventPayload})
		if err != nil {
			log.Error().Str("module", "dispatcher").Msgf("Channel(Pagerduty): Connection failed due to Error: %s", err)
			return err
		}
		if resp.Status == SUCCESS {
			preps := make([]pagerdutyMSGEvent, 0)
			if prepTriggers, ok := p.pagerEntry.Load(methodName); ok {
				preps = prepTriggers.([]pagerdutyMSGEvent)
			}
			preps = append(preps, pagerdutyMSGEvent{flag: pagerdutySeverity, deKey: resp.DedupKey})
			p.pagerEntry.Store(methodName, preps)
		}
	} else if pdResolver, ok := p.pagerEntry.Load(methodName); ok {
		resolver := pdResolver.([]pagerdutyMSGEvent)
		for _, prepFlags := range resolver {
			_, err := pd.ManageEventWithContext(ctx, pd.V2Event{
				RoutingKey: p.secret,
				Action:     "resolve",
				DedupKey:   prepFlags.deKey,
				Payload:    v2EventPayload,
			})
			if err != nil {
				log.Error().Str("module", "dispatcher").Msgf("Channel(Pagerduty): Connection failed due to Error: %s", err)
				return err
			}
		}
		p.pagerEntry.Store(methodName, make([]pagerdutyMSGEvent, 0))
	}

	return nil
}
