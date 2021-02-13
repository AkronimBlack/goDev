package common

import (
	"log"
	"strings"
)

type NameData struct {
	Platform    string `json:"platform"`
	Maintainer  string `json:"maintainer"`
	ProjectName string `json:"project_name"`
}

//ExtractNameData will break fullname into parts by /
//Part 0 will be platform (eg. github.com)
//Part 1 will be maintainer (er. AkronimBlack)
//Part 3 and everything after it fill be joined back by / and be set as project name
func ExtractNameData(fullname string) *NameData {
	var maintainer string
	var projectName string
	var platform string

	parts := strings.Split(fullname, "/")
	if len(parts) > 1 {
		fullname = parts[len(parts)-1]
	}
	if len(parts) == 3 {
		maintainer = parts[1]
		platform = parts[0]
		projectName = strings.Join(parts[2:], "/")
	}

	return &NameData{
		Platform:    platform,
		Maintainer:  maintainer,
		ProjectName: projectName,
	}
}

//SanitizeName will check for whitespaces. if exist will break the string into componenets and merge it back
//as camelCase eg. "Some test string" => "someTestString"
func SanitizeName(value string) string {
	if strings.Contains(value, " ") {
		log.Println("Whitespace detected. Cleaning up app name ...")
		components := strings.Split(value, " ")
		for i, v := range components {
			if i == 0 {
				v = strings.ToLower(v)
				continue
			}
			v = strings.Title(strings.ToLower(v))
		}
		value = strings.Join(components, "")
		log.Printf("App name: %s", value)
	}
	return value
}
