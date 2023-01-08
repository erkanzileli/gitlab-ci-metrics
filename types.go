package main

import (
	"github.com/xanzy/go-gitlab"
	"strconv"
	"time"
)

func newJob(gitlabJob *gitlab.Job) Job {
	return Job{
		ID:             strconv.Itoa(gitlabJob.ID),
		Name:           gitlabJob.Name,
		Stage:          gitlabJob.Stage,
		ProjectID:      strconv.Itoa(gitlabJob.Pipeline.ProjectID),
		Status:         gitlabJob.Status,
		Duration:       gitlabJob.Duration,
		QueuedDuration: gitlabJob.QueuedDuration,
		WebURL:         gitlabJob.WebURL,
		CreatedAt:      *gitlabJob.CreatedAt,
		StartedAt:      *gitlabJob.StartedAt,
		FinishedAt:     *gitlabJob.FinishedAt,
	}
}

type Job struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Stage          string    `json:"stage"`
	ProjectID      string    `json:"projectID"`
	Status         string    `json:"status"`
	Duration       float64   `json:"duration"`
	QueuedDuration float64   `json:"queued_duration"`
	WebURL         string    `json:"webURL"`
	JobTrace       JobTrace  `json:"jobTrace,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	StartedAt      time.Time `json:"started_at"`
	FinishedAt     time.Time `json:"finished_at"`
}

type JobTrace struct {
	Sections []JobSection `json:"sections,omitempty"`
}

type JobSection struct {
	Name       string    `json:"name,omitempty"`
	DurationMS int64     `json:"durationMs"`
	Start      time.Time `json:"start,omitempty"`
	End        time.Time `json:"end,omitempty"`
}
