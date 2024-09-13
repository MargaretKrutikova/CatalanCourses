package main

import (
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CourseInfo struct {
	Code                       string `json:"code"`
	Schedule                   string `json:"schedule"`
	StartDate                  string `json:"startDate"`
	EndDate                    string `json:"endDate"`
	PlacesAvailable            string `json:"placesAvailable"`
	RegistrationDeadline       string `json:"registrationDeadline"`
	PlacesLeft                 string `json:"placesLeft"`
	IsPreferentialRegistration string `json:"isPreferentialRegistration"`
	Address                    string `json:"address"`
	Metros                     string `json:"metros"`
	Email                      string `json:"email"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getText(s *goquery.Selection) string {
	text := s.Text()
	trimmedText := strings.Trim(text, " \n")
	return trimmedText
}

func readFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	check(err)

	return string(data)
}

func containsMultiple(text string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(text, sub) {
			return true
		}
	}
	return false
}

func parseCourseInfo(htmlReader io.Reader) *CourseInfo {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	check(err)

	courseInfo := new(CourseInfo)
	doc.Find("div.container").Find("div.row").Children().Each(func(i int, s *goquery.Selection) {
		s.Find("dl.dl-horizontal").Children().Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Course code" || s.Text() == "Codi del curs" {
				courseInfo.Code = getText(s.Next())
			}
			if containsMultiple(s.Text(), []string{"Days and times", "Dies i horari"}) {
				courseInfo.Schedule = getText(s.Next())
			}
			if containsMultiple(s.Text(), []string{"Start and end of the course", "Inici i final del curs"}) {
				dates := strings.Split(getText(s.Next()), " - ")
				courseInfo.StartDate = dates[0]
				courseInfo.EndDate = dates[1]
			}
			if containsMultiple(s.Text(), []string{"Places"}) {
				courseInfo.PlacesAvailable = getText(s.Next())
			}
			if containsMultiple(s.Text(), []string{"Registration deadline", "Inscripció general"}) {
				courseInfo.RegistrationDeadline = getText(s.Next())
			}
			if containsMultiple(s.Text(), []string{"Preferential registration", "Inscripció preferent"}) {
				courseInfo.IsPreferentialRegistration = getText(s.Next())
			}
		})

		s.Find("div.panel").Children().Each(func(i int, child *goquery.Selection) {
			if containsMultiple(child.Find("h2").Text(), []string{"Place availability", "Disponibilitat de places"}) {
				courseInfo.PlacesLeft = getText(child.Next().Children().Find("h3"))
			}

			if containsMultiple(child.Find("h2").Text(), []string{"Classroom and Center data", "Dades de l'aula i del Centre"}) {

				child.Parent().Find("table").Find("tr").Each(func(i int, tr *goquery.Selection) {
					if containsMultiple(tr.Find("th").Text(), []string{"Address", "Adreça"}) {
						courseInfo.Address = getText(tr.Find("td"))
					}
					if strings.Contains(tr.Find("th").Text(), "Metro") {
						tr.Find("td").Children().Each(func(i int, metro *goquery.Selection) {
							courseInfo.Metros = getText(metro)
						})
					}
					if tr.Find("th").Text() == "EmailCentre" {
						courseInfo.Email = getText(tr.Find("td"))
					}
				})
			}
		})
	})

	return courseInfo
}
