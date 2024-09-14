package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type CourseCode struct {
	ApiCode    string
	CourseCode string
}

const findCourseUrl = "https://inscripcions.cpnl.cat/preregistration/index/"

func loadEnv() {
	err := godotenv.Load()
	check(err)
}

func getDetailedCourseInfo(courseCode string) string {
	url := findCourseUrl + courseCode

	resp, err := http.NewRequest(http.MethodGet, url, nil)
	cookie := os.Getenv("SITE_COOKIE")
	resp.Header.Add("Cookie", cookie)

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		println(err)
		println("Error for course: " + courseCode)
		return ""
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	check(err)

	return string(body)
}

func readCourseCodes() []CourseCode {
	byteValue, err := os.ReadFile("./data/course_list.json")
	check(err)

	var result map[string]any
	if err := json.Unmarshal(byteValue, &result); err != nil {
		panic(err)
	}

	data := result["data"].([]map[string]any)

	var codes []CourseCode
	for _, element := range data {
		if code, ok := element["codiPlain"].(string); ok {
			if apiCode, ok := element["codi"].(string); ok {
				codes = append(codes, CourseCode{ApiCode: apiCode, CourseCode: code})
			}
		}
	}

	return codes
}

func getCourseCodeFromFileName(fileName string) string {
	return strings.Replace(strings.Replace(fileName, "course_", "", 1), ".html", "", 1)
}

func getReadCourseCodes() []string {
	entries, err := os.ReadDir("./data/courses")
	check(err)

	var codes []string
	for _, entry := range entries {
		code := getCourseCodeFromFileName(entry.Name())
		codes = append(codes, code)
	}
	return codes
}

func saveDetailedCourseInfo(code CourseCode) {
	courseHtml := getDetailedCourseInfo(code.ApiCode)

	writeToFile("./data/courses/course_"+code.CourseCode+".html", courseHtml)
}

func loadAllCourseHtmls() {
	codes := readCourseCodes()
	readCodes := getReadCourseCodes()

	for _, code := range codes {
		if slices.Contains(readCodes, code.CourseCode) {
			println("Skipping course: " + code.CourseCode)
			continue
		}

		saveDetailedCourseInfo(code)
		time.Sleep(2 * time.Second)
	}
}

func getCourseApiCodeByCodes(courseCodes []string) []CourseCode {
	allCodes := readCourseCodes()
	apiCodes := []CourseCode{}

	for _, code := range allCodes {
		if slices.Contains(courseCodes, code.CourseCode) {
			apiCodes = append(apiCodes, code)
		}
	}

	return apiCodes
}

func main() {
	loadEnv()

	entries, err := os.ReadDir("./data/courses")
	check(err)

	courseInfos := []*CourseInfo{}

	emptyCourses := []string{}
	for _, entry := range entries {
		path := "./data/courses/" + entry.Name()
		reader, err := os.Open(path)
		check(err)

		parsedInfo := parseCourseInfo(reader)
		if parsedInfo.Code == "" {
			courseCode := getCourseCodeFromFileName(entry.Name())
			println(courseCode)
			emptyCourses = append(emptyCourses, courseCode)
		} else {
			courseInfos = append(courseInfos, parsedInfo)
		}
	}

	writeToFile("./data/empty_courses.json", strings.Join(emptyCourses, ","))
	jsonBytes, err := json.Marshal(courseInfos)
	check(err)

	writeToFile("./data/complete_course_info.json", string(jsonBytes[:]))
}
