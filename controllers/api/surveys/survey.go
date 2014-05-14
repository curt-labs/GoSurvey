package surveys

import (
	"encoding/json"
	"github.com/curt-labs/GoSurvey/models/survey"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

type SurveyError struct {
	Message string `json:"error"`
}

type SurveyResponse struct {
	Surveys      []survey.Survey `json:"surveys"`
	TotalSurveys int             `json:"total_surveys"`
	CurrentPage  int             `json:"current_page"`
	TotalResults int             `json:"total_results"`
}

func All(rw http.ResponseWriter, req *http.Request, r render.Render) {
	params := req.URL.Query()
	var take int
	var page int
	var err error
	total := make(chan int, 0)

	go func() {
		total <- survey.SurveyCount()
	}()

	take, err = strconv.Atoi(params.Get("count"))
	page, err = strconv.Atoi(params.Get("page"))

	skip := page * take
	if page > 0 {
		skip = (page - 1) * take
	}

	svs, err := survey.GetSurveys(skip, take)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if page == 0 {
		page = 1
	}

	sr := SurveyResponse{
		Surveys:      svs,
		CurrentPage:  page,
		TotalResults: len(svs),
		TotalSurveys: <-total,
	}

	r.JSON(200, sr)
}

func Get(rw http.ResponseWriter, req *http.Request, r render.Render, params martini.Params) {
	var sv survey.Survey
	var err error

	if sv.ID, err = strconv.Atoi(params["id"]); err != nil {
		r.JSON(500, SurveyError{err.Error()})
		return
	}

	if err := sv.Get(); err != nil {
		r.JSON(500, SurveyError{err.Error()})
		return
	}

	r.JSON(200, sv)
}

func Submit(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	dec := json.NewDecoder(req.Body)
	var s survey.SurveySubmission
	err := dec.Decode(&s)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Submit()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	success := struct {
		Success bool `json:"success"`
	}{
		true,
	}

	js, _ := json.Marshal(success)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}
