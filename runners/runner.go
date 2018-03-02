package runners

import (
	"log"
	"time"
)

type Runner struct {
	Name string
	run  func() error
}

func NewRunner(name string, run func() error) *Runner {
	return &Runner{
		Name: name,
		run:  run,
	}
}

func (r *Runner) Run() error {
	return r.run()
}

func (r *Runner) RunEvery(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			log.Printf("[INFO] [%s] Starting run", r.Name)
			if err := r.Run(); err != nil {
				log.Printf("[ERROR] [%s] %v", r.Name, err)
			}

			log.Printf("[INFO] [%s] Run complete", r.Name)
		}
	}()

	return ticker
}
