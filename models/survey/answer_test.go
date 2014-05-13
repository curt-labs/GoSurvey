package survey

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAnswer(t *testing.T) {
	Convey("Testing Answer", t, func() {

		var s Survey
		s.Name = "Test Survey"
		s.Description = "Now that we know who you are, I know who I am. I'm not a mistake! It all makes sense! In a comic, you know how you can tell who the arch-villain's going to be? He's the exact opposite of the hero. And most times they're friends, like you and me! I should've known way back when... You know why, David? Because of the kids. They called me Mr Glass."
		s.UserID = 1
		err := s.Save()
		So(err, ShouldEqual, nil)

		var q Question
		Convey("Adding question with answers that have no user reference", func() {
			q = Question{
				Question: "You think water moves fast?",
				UserID:   1,
			}
			q.Answers = append(q.Answers, Answer{
				Answer: "You should see ice.",
			})
			q, err := s.AddQuestion(q)
			So(err, ShouldNotEqual, nil)
			So(len(q.Answers), ShouldEqual, 0)
		})

		Convey("Adding question with answers", func() {
			// Insert question and answer
			q := Question{
				Question: "You think water moves fast?",
				UserID:   1,
			}
			q.Answers = []Answer{
				Answer{
					Answer:   "",
					UserID:   1,
					DataType: "",
				},
			}
			q, err = s.AddQuestion(q)
			So(err, ShouldNotEqual, nil)
			So(len(q.Answers), ShouldEqual, 0)

			q.Answers = []Answer{
				Answer{
					Answer:   "",
					UserID:   1,
					DataType: "multiple",
				},
			}
			q, err = s.AddQuestion(q)
			So(err, ShouldNotEqual, nil)
			So(len(q.Answers), ShouldEqual, 0)

			q.Answers = []Answer{
				Answer{
					Answer:   "",
					UserID:   1,
					DataType: "input",
				},
			}
			q, err = s.AddQuestion(q)
			So(err, ShouldEqual, nil)
			So(len(q.Answers), ShouldNotEqual, 0)

			q.Answers[0].UserID = 1
			q.Answers[0].DataType = "multiple"
			q.Answers[0].Answer = "fail"
			q, err := s.AddQuestion(q)
			So(err, ShouldEqual, nil)
			So(len(q.Answers), ShouldNotEqual, 0)

			Convey("Test getting answer revisions now that we've added a question", func() {
				err = s.Get()
				So(err, ShouldEqual, nil)
			})
		})

		Convey("Deleting answer", func() {
			err := q.Answers[0].Delete()
			So(err, ShouldEqual, nil)
		})
	})
}
