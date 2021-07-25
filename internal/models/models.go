package models

type Report struct {
	DUID         string `json:"du_id"`
	SQCCheck     string `json:"sqc_check"`
	Status       string `json:"status"`
	TemplateName string `json:"template_name"`
	Collected    int    `json:"collected"`
	Passed       int    `json:"passed"`
	Failed       int    `json:"failed"`
}
