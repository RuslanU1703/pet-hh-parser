package data

type ResponseApp struct {
	Found     int
	Vacancies map[string]CastomVacancy
}
type CastomResponseHHApi struct {
	PerPage int             `json:"per_page"`
	Items   []CastomVacancy `json:"items"`
	Page    int             `json:"page"`
	Pages   int             `json:"pages"`
	Found   int             `json:"found"`
}
type CastomVacancy struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Area    Area    `json:"area"`
	Salary  Salary  `json:"salary"`
	Type    Type    `json:"type"`
	Url     string  `json:"url"`
	Snippet Snippet `json:"snippet"`
}
type Salary struct {
	To       int    `json:"to"`
	From     int    `json:"from"`
	Currency string `json:"currency"`
	Gross    bool   `json:"gross"`
}
type Area struct {
	Url  string `json:"url"`
	Id   string `json:"id"`
	Name string `json:"name"`
}
type Type struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type Snippet struct {
	Requirement    string `json:"requirement"`
	Responsibility string `json:"responsibility"`
}

// for handlers
type SearchURL struct {
	Id  int
	Url string `json:"url"`
}
type DeleteDTOSearchURL struct {
	Id int `json:"id"`
}
