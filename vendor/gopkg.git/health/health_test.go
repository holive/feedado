package health

import (
	"context"
	"sort"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type TestChecker struct {
	Result DependencyResult
}

func (c TestChecker) Check(ctx context.Context) DependencyResult {
	time.Sleep(40 * time.Millisecond)
	return c.Result
}

// ByStatus is used to ensure the order during the assertion.
type ByStatus []DependencyResult

func (s ByStatus) Len() int { return len(s) }

func (s ByStatus) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByStatus) Less(i, j int) bool { return s[i].Status < s[j].Status }

func TestCheck(t *testing.T) {
	Convey("Given a list of dependencies", t, func() {
		tests := []struct {
			name     string
			checker  []Checker
			expected Result
		}{
			{
				"It should report 'OK' when it has no dependencies",
				[]Checker{},
				Result{
					Status:       StatusOK,
					Message:      "The application is fully functional.",
					Dependencies: []DependencyResult{},
				},
			},
			{
				"It should report 'OK' when all dependencies are ok",
				[]Checker{
					TestChecker{DependencyResult{Status: DependencyOK}},
				},
				Result{
					Status:  StatusOK,
					Message: "The application is fully functional.",
					Dependencies: []DependencyResult{
						{Status: StatusOK},
					},
				},
			},
			{
				"It should report 'PARTIAL' when one of dependencies are not ok",
				[]Checker{
					TestChecker{DependencyResult{Status: DependencyOK}},
					TestChecker{DependencyResult{Status: DependencyFail}},
					TestChecker{DependencyResult{Status: DependencyOK}},
				},
				Result{
					Status:  StatusPartial,
					Message: "The application is partially functional.",
					Dependencies: []DependencyResult{
						{Status: StatusFail},
						{Status: StatusOK},
						{Status: StatusOK},
					},
				},
			},
			{
				"It should report 'FAIL' when one of critical dependencies are not ok",
				[]Checker{
					TestChecker{DependencyResult{Status: DependencyOK}},
					TestChecker{DependencyResult{Status: DependencyFail, Critical: true}},
					TestChecker{DependencyResult{Status: DependencyOK}},
				},
				Result{
					Status:  StatusFail,
					Message: "The application is not functional.",
					Dependencies: []DependencyResult{
						{Status: StatusFail, Critical: true},
						{Status: StatusOK},
						{Status: StatusOK},
					},
				},
			},
		}

		for _, tt := range tests {
			Convey(tt.name, func() {
				h := NewHealth(tt.checker)
				actual := h.Check(context.Background())
				sort.Sort(ByStatus(actual.Dependencies))
				So(actual, ShouldResemble, tt.expected)
			})
		}
	})
}
