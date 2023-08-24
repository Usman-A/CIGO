package algorithms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduleOrder(t *testing.T) {
	g := graphInternal{}
	taskGraph, err := g.GraphTargets("proj_a", "deploy")
	if err != nil {
		t.Fatal(err)
	}
	s, err := CreateScheduler(taskGraph)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, s.Size(), 4, "There should be 4 tasks")

	assert.Equal(t, []string{"proj_a/init", "proj_a/build", "proj_a/test", "proj_a/deploy"}, s.pending)

}

func TestNext(t *testing.T) {
	g := graphInternal{}
	taskGraph, err := g.GraphTargets("proj_a", "deploy")
	if err != nil {
		t.Fatal(err)
	}
	s, err := CreateScheduler(taskGraph)
	if err != nil {
		t.Fatal(err)
	}

	next, err := s.Next()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "proj_a/init", next, "The first task should be thee init")
	next, err = s.Next()
	assert.Error(t, err, "Should return an error when there are no available tasks.")
	assert.Equal(t, "", next, "The next task should be empty when there are no tasks to do.")
	s.Done("proj_a", "init")
	next, err = s.Next()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "proj_a/build", next, "The next task should be build.")
}
