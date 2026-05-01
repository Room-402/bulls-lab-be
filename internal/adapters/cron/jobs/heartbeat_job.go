package jobs

import (
	"log"
)

type HeartbeatJob struct{}

func NewHeartbeatJob() *HeartbeatJob {
	return &HeartbeatJob{}
}

func (j *HeartbeatJob) Name() string {
	return "Heartbeat"
}

func (j *HeartbeatJob) Run() {
	log.Println("💓 Heartbeat: Service is alive and running background tasks")
}
