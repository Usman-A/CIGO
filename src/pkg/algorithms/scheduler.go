package algorithms

import (
	"4zp6/cigo/pkg/data"
	"4zp6/cigo/pkg/misc"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/dominikbraun/graph"
)

type IScheduler interface {
	Next() (string, error)
	Done(string, string)
	Size() int
}

type Scheduler struct {
	pending   []string
	done      []string
	doneMutex sync.Mutex
	projects  map[string]data.ProjectDefinition
}

func CreateScheduler(targetGraph targetGraph) (*Scheduler, error) {
	// Type aliases in golang are bad. Need explicit cast here
	scheduledJobs, err := graph.TopologicalSort(graph.Graph[string, graphTarget](targetGraph))
	if err != nil {
		return nil, err
	}

	projects, err := misc.GetProjects()
	if err != nil {
		return nil, err
	}

	s := Scheduler{
		pending:  scheduledJobs,
		done:     []string{},
		projects: projects,
	}

	return &s, nil
}

func (s *Scheduler) Next() (string, error) {
	for i, v := range s.pending {
		ready := true

		t := strings.Split(v, "/")
		target := s.projects[t[0]].Targets[t[1]]
		for _, dt := range target.DependsOn {
			p := dt.Project
			if dt.Project == "self" {
				p = t[0]
			}
			if !misc.Contains(s.done, fmt.Sprintf("%s/%s", p, dt.Target)) {
				ready = false
				break
			}
		}

		if ready {
			// golang is bad :(
			// removing the target from pending
			s.pending = append(s.pending[:i], s.pending[i+1:]...)
			return v, nil
		}

		for _, dt := range target.DependsOn {
			if misc.Contains(s.pending, fmt.Sprintf("%s/%s", dt.Project, dt.Target)) {
				return "", errors.New("No ready tasks")
			}
		}
	}

	return "", errors.New("No ready tasks")
}

func (s *Scheduler) Done(project string, target string) {
	s.doneMutex.Lock()
	s.done = append(s.done, fmt.Sprintf("%s/%s", project, target))
	s.doneMutex.Unlock()
}

func (s *Scheduler) Size() int {
	return len(s.pending)
}
