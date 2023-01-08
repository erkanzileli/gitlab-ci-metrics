package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	flagGitlabProjectID string
	flagGitlabBaseURL   string
	flagGitlabToken     string
)

const (
	sectionStateStart = "section_start"
	sectionStateEnd   = "section_end"
)

func init() {
	flag.StringVar(&flagGitlabProjectID, "project", "", "your gitlab project id")
	flag.StringVar(&flagGitlabBaseURL, "gitlab-base-url", "", "gitlab base url of your instance")
	flag.StringVar(&flagGitlabToken, "gitlab-token", "", "gitlab token of your instance")
	flag.Parse()
}

func main() {
	gitlabClient, err := gitlab.NewClient(flagGitlabToken, gitlab.WithBaseURL(flagGitlabBaseURL))
	if err != nil {
		logrus.WithError(err).Fatalf("failed to create gitlab client")
	}

	buildStateValues := []gitlab.BuildStateValue{"success"}
	jobs, _, err := gitlabClient.Jobs.ListProjectJobs(flagGitlabProjectID, &gitlab.ListJobsOptions{
		Scope: &buildStateValues,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("failed to list project jobs")
	}

	var jobAnalysis []Job
	for _, job := range jobs {
		traceFile, _, err := gitlabClient.Jobs.GetTraceFile(flagGitlabProjectID, job.ID)
		if err != nil {
			logrus.WithError(err).Fatalf("failed to get job trace")
		}

		jobLogs, err := io.ReadAll(traceFile)
		if err != nil {
			logrus.WithError(err).Fatalf("failed to read job trace")
		}

		analyzedJob := newJob(job)
		analyzedJob.JobTrace = parseJobTrace(jobLogs)
		jobAnalysis = append(jobAnalysis, analyzedJob)
	}

	jsonContent, err := json.MarshalIndent(jobAnalysis, "", "  ")
	if err != nil {
		logrus.WithError(err).Fatalf("failed to marshall job analysis")
	}

	fmt.Print(string(jsonContent))
}

func parseJobTrace(jobLogs []byte) JobTrace {
	sectionsByName := make(map[string]JobSection)
	regex := regexp.MustCompile(`(.*)(section_(start|end):[0-9]+:[a-z_]+)`)

	preProcessedLogs := strings.ReplaceAll(string(jobLogs), "\u001B[0K", "\n")

	allString := regex.FindAllString(preProcessedLogs, -1)
	if len(allString) == 0 {
		return JobTrace{}
	}

	for _, sectionLine := range allString {
		sectionParts := strings.Split(sectionLine, ":")
		if len(sectionParts) != 3 {
			logrus.Warnf("could not parse section line %q", sectionLine)
			continue
		}

		sectionState, timestampStr, sectionName := sectionParts[0], sectionParts[1], sectionParts[2]

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			logrus.Warnf("could not parse timestamp into integer, sectionLine: %q", sectionLine)
			continue
		}

		if sectionState == sectionStateStart {
			sectionsByName[sectionName] = JobSection{
				Name:  sectionName,
				Start: time.Unix(timestamp, 0),
			}
		} else if sectionState == sectionStateEnd {
			currentSection, exists := sectionsByName[sectionName]
			if !exists {
				logrus.Warnf("section ended before started, sectionLine: %q", sectionLine)
				continue
			}
			currentSection.End = time.Unix(timestamp, 0)
			currentSection.DurationMS = currentSection.End.Sub(currentSection.Start).Milliseconds()
			sectionsByName[sectionName] = currentSection
		} else {
			logrus.Warnf("unknown section state, sectionLine: %q", sectionLine)
			continue
		}
	}

	sectionList := make([]JobSection, 0, len(sectionsByName))
	for _, section := range sectionsByName {
		sectionList = append(sectionList, section)
	}

	sort.Slice(sectionList, func(i, j int) bool {
		return sectionList[i].Start.Unix() < sectionList[j].Start.Unix()
	})

	return JobTrace{
		Sections: sectionList,
	}
}
