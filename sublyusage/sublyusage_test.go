package sublyusage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dc0d/subly"
	nats "github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func TestMain(m *testing.M) {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	code := m.Run()

	os.Exit(code)
}

type TimeRequest struct {
	From string `json:"from"`
}

type TimeResponse struct {
	From string    `json:"from"`
	T    time.Time `json:"time"`
}

type timeService struct {
	econn *nats.EncodedConn
}

var requests = make(chan TimeRequest, 10)

func (ts *timeService) ShowMessage(tr *TimeRequest) {
	requests <- *tr
	// fmt.Printf("%v asked for time @%v", tr.From, time.Now())
}

func (ts *timeService) TellMessage(subject, reply string, tr *TimeRequest) {
	ts.econn.Publish(reply, &TimeResponse{From: tr.From, T: time.Now()})
}

func (ts *timeService) WaitMessageQueue(subject, reply string, tr *TimeRequest) {
	ts.econn.Publish(reply, &TimeResponse{From: tr.From, T: time.Now()})
}

func TestSubscriberSubscribe(t *testing.T) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	econn, err := nats.NewEncodedConn(conn, "json")
	if err != nil {
		t.Fatal(err)
	}
	defer econn.Close()

	s := subly.NewSubscriber(ctx, econn)
	s.Subscribe(&timeService{econn})

	send := &TimeRequest{From: "dc0d"}
	rply := &TimeResponse{}
	err = econn.Publish("timeservice.show", send)
	if !assert.NoError(t, err) {
		return
	}
	select {
	case <-requests:
	case <-time.After(time.Second * 3):
		t.Fail()
	}

	send = &TimeRequest{From: "dc0d"}
	rply = &TimeResponse{}
	err = econn.Request("timeservice.tell", send, rply, time.Second)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "dc0d", rply.From)
	assert.Condition(t, func() bool {
		if time.Now().Sub(rply.T) > time.Second*2 {
			return false
		}
		return true
	})

	send = &TimeRequest{From: "dc0d"}
	rply = &TimeResponse{}
	err = econn.Request("timeservice.wait", send, rply, time.Second)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "dc0d", rply.From)
	assert.Condition(t, func() bool {
		if time.Now().Sub(rply.T) > time.Second*2 {
			return false
		}
		return true
	})
}

func TestSubscriberSubscribeFunc(t *testing.T) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	econn, err := nats.NewEncodedConn(conn, "json")
	if err != nil {
		t.Fatal(err)
	}
	defer econn.Close()

	s := subly.NewSubscriber(ctx, econn)
	{
		srv := &timeService{econn}
		s.SubscribeFunc(
			map[string]interface{}{
				"timeservice.show": srv.ShowMessage,
				"timeservice.tell": srv.TellMessage,
			})
		s.SubscribeFunc(
			map[string]interface{}{
				"timeservice.wait": srv.WaitMessageQueue,
			}, "timeservice_wait")
	}

	send := &TimeRequest{From: "dc0d"}
	rply := &TimeResponse{}
	err = econn.Publish("timeservice.show", send)
	if !assert.NoError(t, err) {
		return
	}
	select {
	case <-requests:
	case <-time.After(time.Second * 3):
		t.Fail()
	}

	send = &TimeRequest{From: "dc0d"}
	rply = &TimeResponse{}
	err = econn.Request("timeservice.tell", send, rply, time.Second)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "dc0d", rply.From)
	assert.Condition(t, func() bool {
		if time.Now().Sub(rply.T) > time.Second*2 {
			return false
		}
		return true
	})

	send = &TimeRequest{From: "dc0d"}
	rply = &TimeResponse{}
	err = econn.Request("timeservice.wait", send, rply, time.Second)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "dc0d", rply.From)
	assert.Condition(t, func() bool {
		if time.Now().Sub(rply.T) > time.Second*2 {
			return false
		}
		return true
	})
}
