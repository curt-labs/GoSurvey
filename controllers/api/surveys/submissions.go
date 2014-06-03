package surveys

import (
	"github.com/curt-labs/GoSurvey/models/survey"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

type SubmissionResponse struct {
	Submissions      []survey.SurveySubmission `json:"submissions"`
	TotalSubmissions int                       `json:"total_submissions"`
	CurrentPage      int                       `json:"current_page"`
	TotalResults     int                       `json:"total_results"`
}

func AllSubmissions(rw http.ResponseWriter, req *http.Request, r render.Render) {
	params := req.URL.Query()
	var take int
	var page int
	var err error
	total := make(chan int, 0)

	go func() {
		total <- survey.SubmissionCount(0)
	}()

	take, err = strconv.Atoi(params.Get("count"))
	page, err = strconv.Atoi(params.Get("page"))

	submissions, err := survey.GetAllSubmissions(page, take, 0)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if page == 0 {
		page = 1
	}

	sr := SubmissionResponse{
		Submissions:      submissions,
		CurrentPage:      page,
		TotalResults:     len(submissions),
		TotalSubmissions: <-total,
	}

	r.JSON(200, sr)
}

func GetSubmissionsBySurvey(rw http.ResponseWriter, req *http.Request, params martini.Params, r render.Render) {
	qs := req.URL.Query()
	var take int
	var page int
	var surveyID int
	var err error
	total := make(chan int, 0)

	take, err = strconv.Atoi(qs.Get("count"))
	page, err = strconv.Atoi(qs.Get("page"))
	surveyID, err = strconv.Atoi(params["id"])
	if surveyID == 0 {
		http.Error(rw, "invalid survey reference", http.StatusInternalServerError)
		return
	}

	go func() {
		total <- survey.SubmissionCount(surveyID)
	}()

	submissions, err := survey.GetAllSubmissions(page, take, surveyID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if page == 0 {
		page = 1
	}

	sr := SubmissionResponse{
		Submissions:      submissions,
		CurrentPage:      page,
		TotalResults:     len(submissions),
		TotalSubmissions: <-total,
	}

	r.JSON(200, sr)
}
