package main

import (
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
	MinGrade    string
	Aliases     []string
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
	unparsedKeyName = courseElem.Find("strong").FullText()

	// aero specific
	if len(strings.Split(unparsedKeyName, ".")) != 2 {
		unparsedKeyName += courseElem.Find("strong").FindNextSibling().FullText()
	}

	var key string
	var aliases []string
	key, aliases = utils.ParseKeyAliases(unparsedKeyName)
	c.Key = key
	c.Aliases = aliases

	c.DeptKey = dept.Key
	strDatArr := strings.Split(unparsedKeyName, ".")

	c.Name = utils.CleanInvisText(strings.Trim(strDatArr[1], " "))
	c.Desc = utils.CleanInvisText(courseElem.Text())

	link := courseElem.Find("a")
	if utils.CheckElemExistence(link) {
		c.ProfileLink = ""
	} else {
		c.ProfileLink = link.Attrs()["href"]
	}

	reqsPlus := courseElem.Find("em")
	// catch error if no cannot find element having information for reqs and credits
	if utils.CheckElemExistence(reqsPlus) {
		c.Prereqs = nil
		c.Coreqs = nil
		c.Credits = 0
		return c
	}

	unparsedReqCreds := strings.Split(reqsPlus.FullText(), ".")

	c.Credits = utils.ParseCredits(unparsedReqCreds)
	return c
}

func seedDeptCourses(dept DepartmentCourseMap) {
	var courses []Course

	resp, err := soup.Get(dept.CourseListURL)
	if err != nil && err.(soup.Error).Type == soup.ErrInGetRequest {
		log.Printf("Error in fetching course list\n %s", err)
		panic(err)
	}
	doc := soup.HTMLParse(resp)
	courseList := doc.Find("div", "class", "entry-content").FindAll("p")
	for _, courseElem := range courseList[1:] {
		// get course and append to courses
		c := parseCourse(courseElem, dept)
		courses = append(courses, c)
	}
	fname := fmt.Sprintf("%s.json", dept.Key)

	utils.SaveDB(fname, courses)
}

func seedDepartments() []DepartmentCourseMap {
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

	var depts []DepartmentCourseMap
	if _, err := os.Stat("data/depts.json"); os.IsNotExist(err) {
		fmt.Println("depts.json not found, fetching from bulletin")
		depts = seedDepartments()
		utils.SaveDB("depts.json", depts)
	} else {
		fmt.Println("depts.json found, seeding database")
		utils.LoadDeptDB(&depts)
	}

	seedDeptCourses(depts[0])

	// for _, dept := range depts {
	// 	fmt.Println("seeding", dept.DeptName)
	// 	seedDeptCourses(dept)
	// }
}
