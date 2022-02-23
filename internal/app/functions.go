package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"pet-app/internal/data"
	"pet-app/internal/helper"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
)

const (
	dataFile = "data.json"
)

func getData(searchURLs map[int]data.SearchURL) map[string]data.CastomVacancy {
	vacanciesMap := make(map[string]data.CastomVacancy)
	result := data.ResponseApp{}
	for key := range searchURLs {
		client := &http.Client{}
		url := helper.PrepareRequestURL(searchURLs[key].Url)
		addDiff := 0
		addURL := ""

		log.Println("[hh]Getting data from searchingURL:", searchURLs[key].Id, "..")
	loop:
		for {
			req, err := http.NewRequest("GET", url+addURL, nil)
			helper.Check(err)
			req.Header.Set("User-Agent", userAgent)

			resp, err := client.Do(req)
			helper.Check(err)
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			helper.Check(err)
			var apiResponse data.CastomResponseHHApi
			err = json.Unmarshal(body, &apiResponse)
			helper.Check(err)
			for i := range apiResponse.Items {
				// check the element in the map and dosage as needed
				vacanciesMap[apiResponse.Items[i].Id] = apiResponse.Items[i]
			}
			diff := apiResponse.PerPage - apiResponse.Found + addDiff
			switch {
			case diff >= 0:
				// all vacancies load
				break loop
			default:
				// need to load more pages
				addURL = "&page=" + strconv.Itoa(apiResponse.Page+1)
				addDiff += apiResponse.PerPage
			}
		}
	}
	log.Println("[hh]Gotted")
	result.Found = len(vacanciesMap)
	result.Vacancies = vacanciesMap
	err := helper.ChangeFile(result, dataFile)
	helper.Check(err)
	return vacanciesMap
}
func fixChanges(searchResultMap, baseSearchResultMap map[string]data.CastomVacancy) {
	var changes []data.CastomVacancy

	for key := range searchResultMap {
		_, ok := baseSearchResultMap[key]
		if !ok {
			changes = append(changes, searchResultMap[key])
		}
	}
	var resultChanges []data.CastomVacancy
	if len(changes) != 0 {
		// fixing api lag
		for i := range changes {
			_, ok := allChanges[changes[i].Id]
			if !ok {
				allChanges[changes[i].Id] = changes[i]
				resultChanges = append(resultChanges, changes[i])
			}
		}
	}
	if len(resultChanges) != 0 {
		log.Println("[app]Got changes")
		err := apiVK(prepareVacanciesURL(changes))
		if err != nil {
			log.Println("[vk]Message sending error: ", err.Error())
		}
		log.Println("[vk]Messge send")
	}
	// no real changes
}
func apiVK(vacancies []data.CastomVacancy) error {
	vk := api.NewVK(config.TokenVK)

	var reqParams = api.Params{}
	reqParams["group_id"] = config.GroupID
	reqParams["user_id"] = config.ProfileID
	reqParams["random_id"] = 0
	reqParams["message"] = prepareMessageVK(vacancies)

	_, err := vk.MessagesSend(reqParams)
	return err
}
func prepareMessageVK(data []data.CastomVacancy) string {
	message := ""
	for i := range data {
		message += "🎯 " + data[i].Name + "(" + data[i].Area.Name + ")\n"
		switch {
		case data[i].Salary.From != 0 && data[i].Salary.To != 0:
			message += "💵 от " + strconv.Itoa(data[i].Salary.From) + " до " + strconv.Itoa(data[i].Salary.To) + "\n"
		case data[i].Salary.From != 0:
			message += "💵 от " + strconv.Itoa(data[i].Salary.From) + "\n"
		case data[i].Salary.To != 0:
			message += "💵 до " + strconv.Itoa(data[i].Salary.To) + "\n"
		default:
			message += "💵 зп не указана" + "\n"
		}

		if data[i].Snippet.Requirement != "" {
			message += "📌 Требования:\n" + data[i].Snippet.Requirement + "\n"
		}
		if data[i].Snippet.Responsibility != "" {
			message += "📌 Обязанности:\n" + data[i].Snippet.Responsibility + "\n"
		}
		message += "url: " + data[i].Url + "\n" + "\n"
	}
	return message
}
func prepareVacanciesURL(changes []data.CastomVacancy) []data.CastomVacancy {
	for i := range changes {
		parts := strings.Split(changes[i].Url, "?")
		changes[i].Url = "https://hh.ru/vacancy/"
		changes[i].Url += strings.Replace(parts[0], "https://api.hh.ru/vacancies/", "", 1)
	}
	return changes
}
