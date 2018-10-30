package wtf

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	//"sync"
)

const SimpleDateFormat = "Jan 2"
const SimpleTimeFormat = "15:04 MST"
const MinimumTimeFormat = "15:04"
const FullDateFormat = "Monday, Jan 2"
const FriendlyDateFormat = "Mon, Jan 2"

//const FriendlyDateTimeFormat = "Mon, Jan 2, 15:04"
const TimestampFormat = "2006-01-02T15:04:05-0700"

func ExecuteCommand(cmd *exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Sprintf("%v\n", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("%v\n", err)
	}

	var str string
	if b, err := ioutil.ReadAll(stdout); err == nil {
		str += string(b)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Sprintf("%v\n", err)
	}

	return str
}

func Exclude(strs []string, val string) bool {
	for _, str := range strs {
		if val == str {
			return false
		}
	}
	return true
}

func FindMatch(pattern string, data string) [][]string {
	r := regexp.MustCompile(pattern)
	return r.FindAllStringSubmatch(data, -1)
}

func NameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	return strings.Title(strings.Replace(parts[0], ".", " ", -1))
}

func NamesFromEmails(emails []string) []string {
	var names []string

	for _, email := range emails {
		names = append(names, NameFromEmail(email))
	}

	return names
}

/* -------------------- Slice Conversion -------------------- */

func ToInts(slice []interface{}) []int {
	var results []int

	for _, val := range slice {
		results = append(results, val.(int))
	}

	return results
}

func ToStrs(slice []interface{}) []string {
	var results []string

	for _, val := range slice {
		results = append(results, val.(string))
	}

	return results
}
