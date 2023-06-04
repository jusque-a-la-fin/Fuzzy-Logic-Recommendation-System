package models

type Search struct {
	Mark           string `form:"mark"`
	Model          string `form:"model"`
	Gearbox        string `form:"gearbox"`
	LowPriceLimit  string `form:"low_price_limit"`
	HighPriceLimit string `form:"high_price_limit"`
	Drive          string `form:"drive"`
	EarliestYear   string `form:"earliest_year"`
	LatestYear     string `form:"lastest_year"`
	Fuel           string `form:"fuel"`
	IsNewCar       string `form:"new"`
}

type CarCard struct {
	// название
	Name string
	// ссылки на изображения
	Images []string
	// цена
	Price string
	// характеристики
	Characteristics map[string]string
}
