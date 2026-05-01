package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

// Job defines the interface for a background task
type Job interface {
	Name() string
	Run()
}

// Scheduler manages the cron jobs
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a new cron scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()), // Enable second-level precision if needed
	}
}

// RegisterJob adds a job to the scheduler with a given cron expression
func (s *Scheduler) RegisterJob(spec string, job Job) (cron.EntryID, error) {
	entryID, err := s.cron.AddFunc(spec, func() {
		log.Printf("[Cron] Running job: %s", job.Name())
		job.Run()
		log.Printf("[Cron] Job finished: %s", job.Name())
	})

	if err != nil {
		return 0, err
	}

	log.Printf("[Cron] Registered job: %s with schedule: %s", job.Name(), spec)
	return entryID, nil
}

// Start starts the cron scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("[Cron] Scheduler started")
}

// Stop stops the cron scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("[Cron] Scheduler stopped")
}
