package worker

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// remove
var TestQuantity = 0
var Timetest time.Time

type MessageTest struct {
	// remove
	Body []byte
}

func (m *Worker) testReceiver(ctx context.Context) (MessageTest, error) {
	time.Sleep(time.Millisecond * 500)
	// remove
	m.Lock()
	defer m.Unlock()

	if TestQuantity == 0 {
		Timetest = time.Now()
	}

	TestQuantity += 1
	if TestQuantity > 1 {
		fmt.Println("Processing time: ", time.Now().Sub(Timetest).String())
		return MessageTest{}, errors.New("finish here")
	}

	return MessageTest{Body: []byte("{\"schema_id\": \"5eeff53aea05875533bcfa75\"}")}, nil
}

func (t MessageTest) Nackable() bool {
	// remove
	return true
}

func (t MessageTest) Nack() {
	// remove
}

func (t MessageTest) Ack() {
	// remove
}
