package request

type ReqQueryJob struct {
	JobName string `json:"job_name" form:"job_name"`
}

type ReqCopyJob struct {
	SrcJobName  string `json:"src_job_name" form:"src_job_name"`
	NewJobName  string `json:"new_job_name" form:"new_job_name"`
	ViewName    string `json:"view_name" form:"view_name"`
	GitUuid     string `json:"git_uuid" form:"git_uuid"`
	BranchUuid  string `json:"branch_uuid" form:"branch_uuid"`
	Application string `json:"application" form:"application"`
}
