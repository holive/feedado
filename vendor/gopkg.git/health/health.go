package health

import (
	"context"
	"sync"
)

// Checker is the interface that wraps the health check method.
type Checker interface {
	Check(context.Context) DependencyResult
}

// Application health status documented in confluence.
var (
	StatusOK      = "OK"
	StatusPartial = "PARTIAL"
	StatusFail    = "FAIL"
)

// Dependency health status documented in confluence.
var (
	DependencyOK   = "OK"
	DependencyFail = "FAIL"
)

// Result which reflects the application health status.
type Result struct {
	Status       string             `json:"status"`
	Message      string             `json:"message"`
	Dependencies []DependencyResult `json:"dependencies,omitempty"`
}

// DependencyResult attributes struct.
type DependencyResult struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Critical    bool   `json:"-"`
	Description string `json:"description,omitempty"`
	Reference   string `json:"reference,omitempty"`
}

// Health status manager.
type Health struct {
	checkers []Checker
}

// NewHealth returns a initialized Health.
func NewHealth(checker []Checker) *Health {
	return &Health{checker}
}

// Check the application health.
func (h *Health) Check(ctx context.Context) Result {
	status := StatusOK
	dependencies := h.checkDependencies(ctx)

	for _, dependency := range dependencies {
		if dependency.Status == DependencyFail {
			if dependency.Critical {
				status = StatusFail
				break
			}
			status = StatusPartial
		}
	}

	return Result{
		Message:      h.message(status),
		Status:       status,
		Dependencies: dependencies,
	}
}

func (h *Health) message(status string) string {
	var message string
	switch status {
	case StatusOK:
		message = "The application is fully functional."
	case StatusPartial:
		message = "The application is partially functional."
	case StatusFail:
		message = "The application is not functional."
	}
	return message
}

func (h *Health) checkDependencies(ctx context.Context) []DependencyResult {
	var (
		dependencies = make([]DependencyResult, 0, len(h.checkers))
		wg           = sync.WaitGroup{}
		chanResult   = make(chan DependencyResult)
		chanDone     = make(chan struct{})
	)

	go func() {
		for dependency := range chanResult {
			dependencies = append(dependencies, dependency)
		}
		close(chanDone)
	}()

	wg.Add(len(h.checkers))
	for _, checker := range h.checkers {
		go func(checker Checker) {
			defer wg.Done()

			chanResult <- checker.Check(ctx)
		}(checker)
	}

	wg.Wait()
	close(chanResult)
	<-chanDone

	return dependencies
}
