package worker

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pkg/errors"

	"go.uber.org/zap"

	"gocloud.dev/pubsub"
)

// Processor process a single message.
type Processor interface {
	Process(ctx context.Context, message []byte) error
}

// Receiver of pubsub messages.
type Receiver interface {
	Receive(ctx context.Context) (*pubsub.Message, error)
}

// Worker monitors processors.
type Worker struct {
	// name of the worker, used for metric and logging purposes.
	name           string
	logger         *zap.SugaredLogger
	concurrency    int
	receiver       Receiver
	processor      Processor
	receiveTimeout time.Duration
	exit           bool
	sync.Mutex
}

// Options passed during initialization.
type Options struct {
	Name           string
	Concurrency    int
	ReceiveTimeout time.Duration
}

// Start receiving messages.
func (m *Worker) Start(ctx context.Context) {
	for i := 0; i < m.concurrency; i++ {
		go func(ctx context.Context, i int) {
			m.logger.Info("Starting", "name", m.name, "number", i)
			for {
				m.Lock()
				if m.exit == true {
					m.Unlock()
					break
				}
				m.Unlock()

				if err := m.receive(ctx); err != nil {
					m.logger.Errorw("Worker receive error",
						"error", err.Error(),
						"name", m.name,
					)

					// TODO: remove this break to release the worker
					break
				}

				if err := ctx.Err(); err != nil {
					m.logger.Info("Exiting", "name", m.name, "number", i)
					break
				}
			}
		}(ctx, i)
	}
}

func (m *Worker) shutdown() error {
	m.Lock()
	m.exit = true
	m.Unlock()
	return nil
}

func (m *Worker) receive(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			m.logger.Errorw(
				"error", err,
				"name", m.name,
				"stack", string(debug.Stack()),
				"message", "Panic",
			)
		}
	}()

	if m.receiveTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, m.receiveTimeout)
	}

	//message, err := m.receiver.Receive(ctx)
	message, err := m.testReceiver(ctx) // TODO: remove test
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			return m.shutdown()
		}

		return errors.Wrap(err, "at message receive method")
	}

	err = m.processor.Process(ctx, message.Body)
	if err != nil {
		if message.Nackable() {
			message.Nack()
		}
		return errors.Wrap(err, "work")
	}
	message.Ack()

	return nil
}

// NewWorker returns a initialized Worker.
func New(opt *Options, logger *zap.SugaredLogger, receiver Receiver, worker Processor) (*Worker, error) {

	return &Worker{
		name:           opt.Name,
		concurrency:    opt.Concurrency,
		receiveTimeout: opt.ReceiveTimeout,
		logger:         logger,
		receiver:       receiver,
		processor:      worker,
	}, nil
}
