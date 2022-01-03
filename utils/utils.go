package utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/mdvsh/coeby/structs"
)

func GetKey(deptLink string) string {
	foo := strings.Split(deptLink, "/")[4]
	// take care of outliers
	if strings.Contains(foo, "-courses") {
		foo = strings.Split(foo, "-courses")[0]
	}
	return foo
}

func ParseKeyAliases(s string) (string, []string) {
	// var of arry of s
	var aliases []string
	var courseKey string
	// go regex doesn't support lookarounds to extract text in multiple groups of parentheses
	// hacky way ahead
	re := regexp.MustCompile(`\(([^\)]+)\)`)
	aliases = re.FindAllString(s, -1)
	if aliases != nil {
		for i, alias := range aliases {
			aliases[i] = strings.Trim(alias, "()")
			s = strings.Replace(s, alias, "", -1)
		}
	} else {
		aliases = []string{}
	}
	courseKey = strings.Trim(strings.Split(s, ".")[0], " ")
	return courseKey, aliases
}

func ParseCredits(s []string) int {
	var credits int
	re := regexp.MustCompile("[0-9]+")
	// find for creds in last element of s
	credArr := re.FindAllString(s[len(s)-1], -1)
	if len(credArr) == 0 {
		credits = 0
	} else if len(credArr) == 1 {
		credits, _ = strconv.Atoi(credArr[0])
	} else {
		var maxCredit int
		for _, cred := range credArr {
			credInt, _ := strconv.Atoi(cred)
			if credInt > maxCredit {
				maxCredit = credInt
			}
		}
		credits = maxCredit
	}
	return credits
}

func CheckElemExistence(e soup.Root) bool {
	return e.Error != nil && e.Error.(soup.Error).Type == soup.ErrElementNotFound
}

func CleanInvisText(desc string) string {
	return strings.Replace(desc, "\u00a0", " ", -1)
}

func CleanFromCredits(raw string) string {
	// get all substrings of parantheses containing the word credit
	re := regexp.MustCompile(`\(([^\)]+)\)`)
	creds := re.FindAllString(raw, -1)

	for _, cred := range creds {
		if strings.Contains(cred, "credit") {
			raw = strings.Replace(raw, cred, "", -1)
		}
	}

	return raw

}

func ParseReqs(raw string) structs.RequisiteProps {
	var reqProps structs.RequisiteProps
	rawlower := strings.ToLower(raw)
	// c1: check for no pre req
	if strings.Contains(raw, "None") || strings.Contains(raw, "none") {
		reqProps.None = true
		return reqProps
	}

	// c2: check for permission of instructor
	if strings.Contains(rawlower, "permission") {
		reqProps.InstructorPerms = true
	}

	// c3: prereqs
	check := strings.Split(raw, ": ")
	if len(check) == 0 {
		return reqProps
	} else if len(check) == 1 {
		reqProps.Notes = CleanInvisText(check[0])
		return reqProps
	}

	// when splits in key and valey for prereq
	if strings.Contains(strings.ToLower(check[0]), "prerequisite") {
		handlePrereqs(check[1], &reqProps)
	}

	reqProps.Raw = raw

	// print temp
	return reqProps
}

func handlePrereqs(raw string, reqProps *structs.RequisiteProps) {
	// split by commas

}
