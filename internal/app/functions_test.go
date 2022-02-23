package app

import (
	"pet-app/internal/data"
	"reflect"
	"testing"
)

func TestApiVK(t *testing.T) {
	var input = []data.CastomVacancy{
		{
			Name: "Go разработчик",
			Area: data.Area{
				Name: "Москва",
			},
			Salary: data.Salary{
				From: 140000,
			},
			Url: "example.com",
			Snippet: data.Snippet{
				Requirement: "Отличное знание Golang и опыт работы с ним от 2 лет. Уверенное знание принципов ООП и основных шаблонов проектирования. Умение оценивать сроки и тщательно планировать свою работу. Непреодолимое желание покрывать код unit/интеграционными тестами",
			},
		},
		{
			Name: "Junior Golang developer",
			Area: data.Area{
				Name: "Москва",
			},
			Url: "example2.com",
			Snippet: data.Snippet{
				Requirement:    "Не менее 1 года работы на GO. Понимание работы и способность пользоваться через CLI Docker",
				Responsibility: "Проектирование и дизайн архитектуры. Разработка микро-сервисов",
			},
		},
	}
	apiVK(input)
}
func TestPrepareVacanciesURL(t *testing.T) {
	var input = []data.CastomVacancy{
		{
			Url: "https://api.hh.ru/vacancies/111111?host=hh.ru",
		},
		{
			Url: "https://api.hh.ru/vacancies/222222?host=hh.ru",
		},
	}
	var expected = []data.CastomVacancy{
		{
			Url: "https://hh.ru/vacancy/111111",
		},
		{
			Url: "https://hh.ru/vacancy/222222",
		},
	}
	gotted := prepareVacanciesURL(input)
	ok := reflect.DeepEqual(expected, gotted)
	if !ok {
		t.Errorf("want: %v got: %v", expected, gotted)
	}
}
