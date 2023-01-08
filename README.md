# gitlab-ci-metrics

Experimental PoC that shows us the bottlenecks on our GitLab CI.

It uses GitLab's Job Trace endpoint to get the complete job logs and parses the predefined pattern.

It runs for only the last 20 successful jobs and exists. You can extend it by copying the code.

## Usage

Just run command below. It'll print your job traces to stdout.

```shell
go run . --project=<your gitlab project id> --gitlab-base-url=<gitlab base url of your instance> --gitlab-token=<gitlab token of your instance> 
```

## Data

Here is a JSON array and many details.

```json
[
  {
    "id": "16628865",
    "name": "build container image",
    "stage": "Build 📦",
    "projectID": "9052",
    "status": "success",
    "duration": 62.357635,
    "queued_duration": 2.239513,
    "webURL": "https://gitlab.company.com/path/to/project/16628865",
    "jobTrace": {
      "sections": [
        {
          "name": "resolve_secrets",
          "durationMs": 0,
          "start": "2023-01-05T21:22:09+03:00",
          "end": "2023-01-05T21:22:09+03:00"
        },
        {
          "name": "prepare_executor",
          "durationMs": 0,
          "start": "2023-01-05T21:22:09+03:00",
          "end": "2023-01-05T21:22:09+03:00"
        },
        {
          "name": "prepare_script",
          "durationMs": 6000,
          "start": "2023-01-05T21:22:09+03:00",
          "end": "2023-01-05T21:22:15+03:00"
        },
        {
          "name": "get_sources",
          "durationMs": 2000,
          "start": "2023-01-05T21:22:15+03:00",
          "end": "2023-01-05T21:22:17+03:00"
        },
        {
          "name": "restore_cache",
          "durationMs": 10000,
          "start": "2023-01-05T21:22:17+03:00",
          "end": "2023-01-05T21:22:27+03:00"
        },
        {
          "name": "download_artifacts",
          "durationMs": 1000,
          "start": "2023-01-05T21:22:27+03:00",
          "end": "2023-01-05T21:22:28+03:00"
        },
        {
          "name": "step_script",
          "durationMs": 41000,
          "start": "2023-01-05T21:22:28+03:00",
          "end": "2023-01-05T21:23:09+03:00"
        },
        {
          "name": "archive_cache",
          "durationMs": 1000,
          "start": "2023-01-05T21:23:09+03:00",
          "end": "2023-01-05T21:23:10+03:00"
        },
        {
          "name": "cleanup_file_variables",
          "durationMs": 0,
          "start": "2023-01-05T21:23:10+03:00",
          "end": "2023-01-05T21:23:10+03:00"
        }
      ]
    },
    "created_at": "2023-01-05T18:21:21.319Z",
    "started_at": "2023-01-05T18:22:08.684Z",
    "finished_at": "2023-01-05T18:23:11.042Z"
  }
]
```