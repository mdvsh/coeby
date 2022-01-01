package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
)

type DepartmentCourseMap struct {
	Key           string
	DeptName      string
	CourseListURL string
}

func getKey(deptLink string) string {
	foo := strings.Split(deptLink, "/")[4]
	// take care of outliers
	if strings.Contains(foo, "-courses") {
		foo = strings.Split(foo, "-courses")[0]
	}
	return foo
}

func save(depts []DepartmentCourseMap, fname string) {
	jsonFile, err := os.Create("data/" + fname)
	if err != nil {
		log.Printf("Error creating json file\n %s", err)
		panic(err)
	}
	defer jsonFile.Close()
	jsonWriter := bufio.NewWriter(jsonFile)
	defer jsonWriter.Flush()
	enc := json.NewEncoder(jsonWriter)
	enc.SetIndent("", "  ")
	enc.Encode(depts)
}

func seed() {
	var depts []DepartmentCourseMap

	var Headers map[string]string
	Headers = make(map[string]string)
	Headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"
	Headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9"
	Headers["Accept-Language"] = "en-US,"
	Headers["Accept-Encoding"] = "gzip, deflate"
	Headers["Connection"] = "keep-alive"

	resp, err := soup.Get("https://bulletin.engin.umich.edu/courses/")
	if err != nil && err.(soup.Error).Type == soup.ErrInGetRequest {
		log.Printf("Error in fetching bulletin\n %s", err)
		panic(err)
	}
	doc := soup.HTMLParse(resp)
	departments := doc.Find("ul", "aria-labelledby", "menu-item-dropdown-897").FindAll("a")
	for _, department := range departments[2:] {
		link := department.Attrs()["href"]
		depts = append(depts, DepartmentCourseMap{
			Key:           getKey(link),
			DeptName:      department.Text(),
			CourseListURL: link,
		})
	}
	save(depts, "departments.json")

}

func main() {
	fmt.Println("booting up coeby")
	seed()
}
