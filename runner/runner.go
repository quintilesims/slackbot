package runner

import (
	"log"
	"time"
)

// Runner objects manage running operations asynchronously, similar to a daemon
type Runner struct {
	Name string
	run  func() error
}

// NewRunner will create a new Runner object with the specified name and run function
func NewRunner(name string, run func() error) *Runner {
	return &Runner{
		Name: name,
		run:  run,
	}
}

// Run executes the runner's function
func (r *Runner) Run() error {
	log.Printf("[INFO] [%s] Starting run", r.Name)
	defer log.Printf("[INFO] [%s] Run complete", r.Name)
	return r.run()
}

// RunEvery will execute the runner's function at the specified interval
func (r *Runner) RunEvery(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			if err := r.Run(); err != nil {
				log.Printf("[ERROR] [%s] %v", r.Name, err)
			}
		}
	}()

	return ticker
}
