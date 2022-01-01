package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func GetKey(deptLink string) string {
	foo := strings.Split(deptLink, "/")[4]
	// take care of outliers
	if strings.Contains(foo, "-courses") {
		foo = strings.Split(foo, "-courses")[0]
	}
	return foo
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
