# gitlab-exporter

Export data from Gitlab (Events) to Prometheus

## Data structure

Metrics and Lables:

```
gitlab_ci_job_number_total      {ci=gitlab, project_name="", branch="", user="", status="", job_stage="", job_tags} counter
gitlab_ci_job_time_total        {ci=gitlab, project_name="", branch="", user="", status="", job_stage=""} counter

gitlab_ci_pipeline_number_total {ci=gitlab, project_name="", branch="", user="", status=""} counter
gitlab_ci_pipeline_time_total   {ci=gitlab, project_name="", branch="", user="", status=""} counter
gitlab_ci_pipeline_size         {ci=gitlab, project_name="", branch="", user="", status=""} counter
```

branch = master,release,MR,other<br />
user = name of user which started pipeline<br />

status = failed,skipped,success,canceled<br />
job_stage = Group by type of jobs<br />
jobs_tags = Group by type of runner for linked it with AWS tags<br />
ci = source data system (gitlab/jenkins)<br />

# Deploy to k8s

```kubectl apply -f k8s-manifest.yaml -n microservices```

## The status of job:
```
failed
warning
pending
running
manual
scheduled
canceled
success
skipped
created
```
## The status of pipeline
```
created
waiting_for_resource
preparing
pending
running
success
failed
canceled
skipped
manual
scheduled
```
