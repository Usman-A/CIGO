package commands

import (
	"4zp6/cigo/pkg/algorithms"
	"4zp6/cigo/pkg/misc"
	"4zp6/cigo/pkg/parser"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/gookit/color"
)

func Run(project string, target string, dryRun bool) error {

	graph, err := algorithms.GetGraph().GraphTargets(project, target)
	if err != nil {
		return err
	}
	s, err := algorithms.CreateScheduler(graph)
	if err != nil {
		return err
	}

	var (
		wait sync.WaitGroup
		done = make(chan bool, s.Size())
	)

	wsPath, err := misc.GetWorkspacePath()
	if err != nil {
		return err
	}
	ws, err := parser.DecodeWorkspace(wsPath, parser.JSON)
	if err != nil {
		return err
	}

	var runErr error = nil

	for s.Size() > 0 {
		next, err := s.Next()
		if err == nil && next != "" {
			wait.Add(1)
			go func() {
				t := strings.Split(next, "/")
				// these should be the inputs of the function
				project := t[0]
				targetStr := t[1]
				defer s.Done(project, targetStr)
				defer func() { done <- true }()
				defer wait.Done()

				target, err := graph.Vertex(next)
				if err != nil {
					if runErr == nil {
						runErr = fmt.Errorf("Errors faced when running commands: \n")
					}
					runErr = errors.New(runErr.Error() + "\t" + err.Error())
					return
				}

				// Check if the target is cached
				if len(target.Target.Artifacts) > 0 {
					isCached, err := IsCached(target.Project, target.Target)
					if err != nil {
						if runErr == nil {
							runErr = fmt.Errorf("Errors faced when running commands: \n")
						}
						runErr = errors.New(runErr.Error() + "\t" + err.Error())
						return
					}
					if isCached {
						fmt.Printf("Target %s for project %s is cached. Skipping execution.\n", target.Name, target.Project)
						return
					}
				}

				// Put the environment variables into the correct format
				env := []string{}
				for k, v := range target.Target.Env {
					env = append(env, fmt.Sprintf("%s=%s", k, v))
				}

				// run the target
				// Run all the commands in the target
				color.Infof("Executing %s for %s:\n", target.Name, target.Project)
				for _, v := range target.Target.Cmds {
					t := strings.Split(v, " ")
					cmd := exec.Command(t[0], t[1:]...)
					cmd.Env = env
					dir, err := misc.GetRelativePath(ws.Projects[project])
					if err != nil {
						if runErr == nil {
							runErr = fmt.Errorf("Errors faced when running commands: \n")
						}
						runErr = errors.New(runErr.Error() + "\t" + err.Error())
						return
					}
					// set the current directory as the project directory
					cmd.Dir = dir
					// set the output to the terminal
					cmd.Stderr = os.Stderr
					cmd.Stdout = os.Stdout

					//  run the command
					// if the dryRun flag is set, don't actually run the command
					// just print it
					if dryRun {
						color.Infof("Dry run: %s\n", cmd.String())
					} else {
						// run the command
						err = cmd.Run()
					}

					// if there is an error, print it
					if err != nil {
						if runErr == nil {
							runErr = fmt.Errorf("Errors faced when running commands: \n")
						}
						runErr = errors.New(runErr.Error() + "\t" + err.Error())
						color.Errorf("Error while executing %s for %s\n", target.Name, target.Project)
						return
					}
				}
				color.Successf("Finished execution.\n\n")
			}()

			continue
		}
		<-done
		for len(done) > 0 {
			<-done
		}
	}

	wait.Wait()

	return runErr
}
