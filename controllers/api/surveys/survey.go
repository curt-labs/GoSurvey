package surveys

import (
	"github.com/curt-labs/GoSurvey/models/survey"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

type SurveyError struct {
	Message string `json:"error"`
}

func All(rw http.ResponseWriter, req *http.Request, r render.Render) {
	svs, err := survey.GetSurveys()
	if err != nil {
		panic(err)
	}

	r.JSON(200, svs)
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

// func Add(rw http.ResponseWriter, req *http.Request, r render.Render, params martini.Params, w warranties.Warranty) {
// 	if err := w.Add(); err != nil {
// 		r.JSON(500, WarrantyError{err.Error()})
// 		return
// 	}

// 	r.JSON(200, w)
// }

// func Delete(rw http.ResponseWriter, req *http.Request, r render.Render, params martini.Params) {

// 	var w warranties.Warranty
// 	var err error
// 	if w.ID, err = strconv.Atoi(params["id"]); err != nil {
// 		r.JSON(500, WarrantyError{err.Error()})
// 		return
// 	}

// 	if err = w.Delete(); err != nil {
// 		r.JSON(500, WarrantyError{err.Error()})
// 		return
// 	}
// 	r.Status(200)
// }
