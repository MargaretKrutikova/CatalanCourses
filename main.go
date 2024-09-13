package main

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Load env variables

const findCourseUrl = "https://inscripcions.cpnl.cat/preregistration/index/"

func loadEnv() {
	err := godotenv.Load()
	check(err)
}

func writeToFile(fileName string, data string) {
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(data)

	w.Flush()
}

func getDetailedCourseInfo(courseCode string) string {
	url := findCourseUrl + courseCode

	resp, err := http.NewRequest(http.MethodGet, url, nil)
	cookie := os.Getenv("SITE_COOKIE")
	resp.Header.Add("Cookie", cookie)

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		print(err)
		print("Error for course: " + courseCode)
		return ""
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	check(err)
	sb := string(body)
	println(response.StatusCode)

	return sb
}

type CourseCode struct {
	ApiCode    string
	CourseCode string
}

func readCourseCodes() []CourseCode {
	jsonFile, err := os.Open("./data/course_list.json")
	check(err)

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}
	if err := json.Unmarshal(byteValue, &result); err != nil {
		panic(err)
	}

	data := result["data"].([]interface{})

	var codes []CourseCode
	for _, element := range data {
		parsed := element.(map[string]interface{})

		if code, ok := parsed["codiPlain"].(string); ok {
			if apiCode, ok := parsed["codi"].(string); ok {
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

func loadCourseHtmls(codes []CourseCode) {
	for _, code := range codes {
		courseHtml := getDetailedCourseInfo(code.ApiCode)
		writeToFile("./data/courses/course_"+code.CourseCode+".html", courseHtml)
		time.Sleep(2 * time.Second)
	}
}

func loadAllCourseHtmls() {
	loadEnv()
	codes := readCourseCodes()
	readCodes := getReadCourseCodes()
	println(readCodes)

	for _, code := range codes {
		if slices.Contains(readCodes, code.CourseCode) {
			println("Skipping course: " + code.CourseCode)
			continue
		}

		courseHtml := getDetailedCourseInfo(code.ApiCode)
		writeToFile("./data/courses/course_"+code.CourseCode+".html", courseHtml)
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
