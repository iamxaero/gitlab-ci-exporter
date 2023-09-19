package controller

type PJob struct {
	Stage    string  `json:"stage"`
	Status   string  `json:"status"`
	Duration float64 `json:"duration"`
	Runner   struct {
		Tags []string `json:"tags"`
	} `json:"runner"`
	User struct {
		Username string `json:"username"`
	} `json:"user"`
}

type Pipeline struct {
	Attributes struct {
		ID       int     `json:"id"`
		Ref      string  `json:"ref"`
		Status   string  `json:"status"`
		Duration float64 `json:"duration"`
	} `json:"object_attributes"`
	Project struct {
		ID   int    `json:"id"`
		Path string `json:"path_with_namespace"`
		Name string `json:"name"`
		URL  string `json:"web_url"`
	} `json:"project"`
	User struct {
		Username string `json:"username"`
	} `json:"user"`
	Jobs []PJob `json:"builds"`
}

type Job struct {
	Ref            string  `json:"ref"`
	ID             int     `json:"build_id"`
	Stage          string  `json:"build_stage"`
	Status         string  `json:"build_status"`
	AllowFailure   bool    `json:"build_allow_failure"`
	PipelineID     int     `json:"pipeline_id"`
	ProjectID      int     `json:"project_id"`
	Project_path   string  `json:"project_name"`
	Build_duration float64 `json:"build_duration"`
	Runner         struct {
		Tags []string `json:"tags"`
	} `json:"runner"`
	User struct {
		Username string `json:"username"`
	} `json:"user"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repository"`
}
