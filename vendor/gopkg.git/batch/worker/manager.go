package worker

import (
	"context"
	"runtime/debug"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"gocloud.dev/pubsub"
	"golang.org/x/sync/errgroup"
)

// Manager monitors workers.
type Manager struct {
	// name of the worker, used for metric and logging purposes.
	name        string
	logger      log.Logger
	concurrency int
	receiver    Receiver
	worker      Worker
}

// Worker process a single message.
type Worker interface {
	Work(ctx context.Context, message []byte) error
}

// Receiver of pubsub messages.
type Receiver interface {
	Receive(ctx context.Context) (*pubsub.Message, error)
}

// ManagerOptions passed during initialization.
type ManagerOptions struct {
	Name        string
	Logger      log.Logger
	Concurrency int
	Receiver    Receiver
	Worker      Worker
}

// Start receiving messages.
func (m *Manager) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < m.concurrency; i++ {
		g.Go(func() error {
			for {
				if err := m.receive(ctx); err != nil {
					// level.Error(m.logger).Log(
					// 	"error", err,
					// 	"name", m.name,
					// )
					// Workaround for connection reset by peer.
					// TODO: Integrate APM.
					return err
				}

				if err := ctx.Err(); err != nil {
					return err
				}
			}
		})
	}

	return g.Wait()
}

func (m *Manager) receive(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			level.Error(m.logger).Log(
				"error", err,
				"name", m.name,
				"stack", string(debug.Stack()),
				"message", "Panic",
			)
		}
	}()

	message, err := m.receiver.Receive(ctx)
	if err != nil {
		return errors.Wrap(err, "message receive")
	}

	err = m.worker.Work(ctx, message.Body)
	if err != nil {
		if message.Nackable() {
			message.Nack()
		}
		return errors.Wrap(err, "work")
	}
	message.Ack()

	return nil
}

// NewManager returns a initialized Manager.
func NewManager(opt ManagerOptions) *Manager {
	return &Manager{
		name:        opt.Name,
		logger:      opt.Logger,
		concurrency: opt.Concurrency,
		receiver:    opt.Receiver,
		worker:      opt.Worker,
	}
}
