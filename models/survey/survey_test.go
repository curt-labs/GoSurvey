package survey

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetSurveys(t *testing.T) {
	Convey("Testing the GetSurveys", t, func() {
		Convey("Testing with GetSurveys(0,0)", func() {
			surveys, err := GetSurveys(0, 0)
			So(err, ShouldEqual, nil)
			So(surveys, ShouldNotEqual, nil)
		})

		Convey("Testing with GetSurveys(10000000,10000)", func() {
			surveys, err := GetSurveys(1000000, 10000)
			So(err, ShouldEqual, nil)
			So(surveys, ShouldNotEqual, nil)
		})
	})
}

func TestSurveyCount(t *testing.T) {
	Convey("Get the total number of surveys", t, func() {
		count := SurveyCount()
		So(count, ShouldNotEqual, nil)
	})
}

func TestGet(t *testing.T) {
	var s Survey
	Convey("Get a specific survey", t, func() {
		Convey("Testing with zero for ID", func() {
			err := s.Get()
			So(err, ShouldNotEqual, nil)
			So(s.Name, ShouldEqual, "")
		})
		Convey("Testing with 1 for ID", func() {
			s.ID = 10000000000
			err := s.Get()
			So(err, ShouldNotEqual, nil)
		})
	})
}

func TestSave(t *testing.T) {
	var s Survey

	Convey("Insert a survey", t, func() {
		Convey("Inserting with empty name", func() {
			err := s.Save()
			So(err, ShouldNotEqual, nil)
		})

		Convey("Inserting with a name and no user reference", func() {
			s.Name = "Test Survey"
			err := s.Save()
			So(err, ShouldNotEqual, nil)
			So(s.Name, ShouldEqual, "Test Survey")
			So(s.ID, ShouldEqual, 0)
		})

		Convey("Deleting a survey with no id reference", func() {
			err := s.Delete()
			So(err, ShouldNotEqual, nil)
		})

		Convey("Inserting with a name and user reference", func() {
			s.Name = "Test Survey"
			s.UserID = 1
			err := s.Save()
			So(err, ShouldEqual, nil)
			So(s.Name, ShouldEqual, "Test Survey")
			So(s.ID, ShouldBeGreaterThan, 0)

			Convey("Testing get now that we've inserted", func() {
				var s2 Survey
				s2.ID = s.ID
				err = s2.Get()
				So(err, ShouldEqual, nil)
				So(s2.Name, ShouldEqual, s.Name)

				Convey("Testing update", func() {
					s2.Name = "Updated Name"
					err = s2.Save()
					So(err, ShouldEqual, nil)
				})

				Convey("Testing delete", func() {
					err := s.Delete()
					So(err, ShouldEqual, nil)

					err = s2.Delete()
					So(err, ShouldEqual, nil)
				})
			})
		})
	})
}
