package main

import (
	"encoding/json"
	"fmt"
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

func parseCourseInfo(htmlLines string) *CourseInfo {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlLines))
	check(err)

	courseInfo := new(CourseInfo)
	doc.Find("div.container").Find("div.row").Children().Each(func(i int, s *goquery.Selection) {
		s.Find("dl.dl-horizontal").Children().Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Course code" {
				courseInfo.Code = getText(s.Next())
			}
			if strings.Contains(s.Text(), "Days and times") {
				courseInfo.Schedule = getText(s.Next())
			}
			if strings.Contains(s.Text(), "Start and end of the course") {
				dates := strings.Split(getText(s.Next()), " - ")
				courseInfo.StartDate = dates[0]
				courseInfo.EndDate = dates[1]
			}
			if strings.Contains(s.Text(), "Places") {
				courseInfo.PlacesAvailable = getText(s.Next())
			}
			if strings.Contains(s.Text(), "Registration deadline") {
				courseInfo.RegistrationDeadline = getText(s.Next())
			}
			if strings.Contains(s.Text(), "Preferential registration") {
				courseInfo.IsPreferentialRegistration = getText(s.Next())
			}
		})

		s.Find("div.panel").Children().Each(func(i int, child *goquery.Selection) {
			if strings.Contains(child.Find("h2").Text(), "Place availability") {
				courseInfo.PlacesLeft = getText(child.Next().Children().Find("h3"))
			}

			if strings.Contains(child.Find("h2").Text(), "Classroom and Center data") {

				child.Parent().Find("table").Find("tr").Each(func(i int, tr *goquery.Selection) {
					if strings.Contains(tr.Find("th").Text(), "Address") {
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

func main() {
	info := parseCourseInfo(readFile("data/course_example.html"))

	jsonBytes, err := json.Marshal(info)
	check(err)

	fmt.Println(string(jsonBytes[:]))
}
