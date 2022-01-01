package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/mdvsh/coeby/utils"
)

type DepartmentCourseMap struct {
	Key           string
	DeptName      string
	CourseListURL string
}

type Course struct {
	DeptKey     string
	Key         string
	Name        string
	Desc        string
	ProfileLink string
	Credits     int
	CoLists     []string
	Prereqs     []Requisite
	Coreqs      []Requisite
}

/*
* type can be:
	enforced
	advisory
	permission of instructor
*/
type Requisite struct {
	CourseKey string
	Type      string
}

func parseCourse(courseElem soup.Root, dept DepartmentCourseMap) Course {
	var c Course
	var unparsedKeyName string
	c.DeptKey = dept.Key
	unparsedKeyName = courseElem.Find("strong").FullText()

	if len(strings.Split(unparsedKeyName, ".")) != 2 {
		unparsedKeyName += courseElem.Find("strong").FindNextSibling().FullText()
	}

	c.Key = strings.Split(unparsedKeyName, ".")[0]
	c.Name = strings.Trim(strings.Split(unparsedKeyName, ".")[1], " ")

	fmt.Println("key: ", c.Key, " name: ", c.Name)

	reqsPlus := courseElem.Find("em")
	// catch error if no cannot find element having information for reqs and credits
	if reqsPlus.Error != nil && reqsPlus.Error.(soup.Error).Type == soup.ErrElementNotFound {
		c.Prereqs = nil
		c.Credits = 0
		return c
	}

	unparsedReqCreds := strings.Split(reqsPlus.FullText(), ".")
	c.Credits = utils.ParseCredits(unparsedReqCreds)
	fmt.Println("credits: ", c.Credits)
	return c
}

func seedDeptCourses(dept DepartmentCourseMap) {
	// var courses []Course

	resp, err := soup.Get(dept.CourseListURL)
	if err != nil && err.(soup.Error).Type == soup.ErrInGetRequest {
		log.Printf("Error in fetching course list\n %s", err)
		panic(err)
	}
	doc := soup.HTMLParse(resp)
	courseList := doc.Find("div", "class", "entry-content").FindAll("p")
	for _, courseElem := range courseList[1:] {
		parseCourse(courseElem, dept)
	}
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
	fmt.Println("saved to database")
}

func seed() []DepartmentCourseMap {
	var depts []DepartmentCourseMap

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
			Key:           utils.GetKey(link),
			DeptName:      department.Text(),
			CourseListURL: link,
		})
	}
	return depts

}

func main() {
	fmt.Println("booting up coeby")
	depts := seed()
	save(depts, "depts.json")
	seedDeptCourses(depts[0])
}
