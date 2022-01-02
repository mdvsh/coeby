package utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
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
	courseKey = strings.Split(s, ".")[0]
	return courseKey, aliases
}

func ParseCredits(s []string) int {
	var credits int
	re := regexp.MustCompile("[0-9]+")
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

func ParseReqs(s []string) {
	// todo
}
