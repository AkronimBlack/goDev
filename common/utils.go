package common

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/bxcodec/faker"
)

//PanicOnError if err != nil log.Panic
func PanicOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type FakerOptions struct {
	Email           string `faker:"email"`
	PhoneNumber     string `faker:"phone_number"`
	URL             string `faker:"url"`
	UserName        string `faker:"username"`
	TitleMale       string `faker:"title_male"`
	TitleFemale     string `faker:"title_female"`
	FirstName       string `faker:"first_name"`
	FirstNameMale   string `faker:"first_name_male"`
	FirstNameFemale string `faker:"first_name_female"`
	LastName        string `faker:"last_name"`
	Name            string `faker:"name"`
	Date            string `faker:"date"`
	Time            string `faker:"time"`
	MonthName       string `faker:"month_name"`
	Year            string `faker:"year"`
	DayOfWeek       string `faker:"day_of_week"`
	DayOfMonth      string `faker:"day_of_month"`
	Timestamp       string `faker:"timestamp"`
	Century         string `faker:"century"`
	TimeZone        string `faker:"timezone"`
	TimePeriod      string `faker:"time_period"`
	Word            string `faker:"word"`
	Sentence        string `faker:"sentence"`
	Paragraph       string `faker:"paragraph"`
	Currency        string `faker:"currency"`
	UUID            string `faker:"uuid_digit"`
}

//ReplacePlaceholder ...
func ReplacePlaceholder(value string) string {
	if strings.Contains(value, "{{") {
		val := trimExcessFat(value, "{{", "}}")
		splitData := strings.Split(val, ".")
		if splitData[0] == "faker" {
			return newFaked(splitData[1])
		}
	}
	return value
}

func newFaked(key string) string {
	v := FakerOptions{}
	err := faker.FakeData(&v)
	if err != nil {
		log.Println("Problem faking data")
		fmt.Println(err)
	}
	rv := reflect.ValueOf(v)
	fv := rv.FieldByName(key)
	return fv.String()
}

func trimExcessFat(value string, left string, right string) string {
	stringValue := strings.TrimLeft(value, left)
	stringValue = strings.TrimRight(stringValue, right)
	return stringValue
}

func LogJson(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println(string(b))
}
