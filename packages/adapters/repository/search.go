package repository

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/repository"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"golang.org/x/exp/maps"
)

type searchRepository struct {
	carCards *[]models.CarCard
}

func NewSearchRepository() repository.SearchRepository {
	var carCards *[]models.CarCard
	return &searchRepository{carCards}
}

// ограничение количества автомобилей
var limit_value int = 10

// ограничение, необходимое при поиске названий и цен автомобиля
var _limit_value int = 2 * limit_value

// Функция получения ссылок на изображения, названий, цен и характеристик автомобилей
func (sr *searchRepository) GetCarsUsingScraping(search *models.Search) *[]models.CarCard {
	// Подготовка параметров
	car_make, model, params := Prepare_parameters(search.Mark, search.Model, search.LowPriceLimit, search.HighPriceLimit, search.EarliestYear, search.LatestYear, search.Gearbox, search.Fuel, search.Drive)

	// ссылка на страницу, содержащую записи о автомобилях конкретной марки, их названиях и ценах
	var link string
	// если коробка передач является "АКПП"
	if search.Gearbox == "АКПП" {
		link = fmt.Sprintf("https://auto.drom.ru/%s/%s/?minprice=%d&maxprice=%d&minyear=%d&maxyear=%d&transmission[]=2&transmission[]=5&fueltype=%d&privod=%d&ph=1", car_make, model, params["min_price"], params["max_price"], params["earliest_year"], params["latest_year"], params["fuel"], params["drive"])
	} else {
		link = fmt.Sprintf("https://auto.drom.ru/%s/%s/?minprice=%d&maxprice=%d&minyear=%d&maxyear=%d&transmission[]=%d&fueltype=%d&privod=%d&ph=1", car_make, model, params["min_price"], params["max_price"], params["earliest_year"], params["latest_year"], params["gearbox"], params["fuel"], params["drive"])
	}

	// шаблон (регулярного выражения) страницы, представляющей один автомобиль
	pattern := fmt.Sprintf(`https://(\w+|\w+\-\w+|\w+\-\w+\-\w+)\.drom.ru/%s/(\w+\-\w+|\w+)/\w+\.html`, car_make)

	// названия автомобилей
	var names []string
	// цены автомобилей
	var prices []string
	// ссылки на страницы, представляющие свой автомобиль
	var links []string

	// получение названий, цен и ссылок на страницы автомобилей
	links, names, prices = Scrape_pages_names_prices(link, pattern, "pages")

	// ссылки на изображения автомобилей
	var images = [][]string{}
	// характеристики автомобилей
	var characteristics []map[string]string

	// шаблон (регулярного выражения) ссылки на изображение автомобиля
	pattern = `https://(s|s1).auto.drom.ru/.+.jpg`

	// цикл по ссылок на страницы, представляющие свой автомобиль
	for _, link := range links {
		// получение изображений автомобиля
		images_per_page := Scrape_links(nil, link, pattern, "images")
		// добавление изображений
		images = append(images, images_per_page)
		// получение характеристик и их проверка
		characteristics_per_one := check_characteristics(Scrape_characteristics(link))
		// добавление характеристик
		characteristics = append(characteristics, characteristics_per_one)
	}

	// автомобили
	var cars []models.CarCard
	// Заполнение данного среза
	for i, img := range images {
		// создание нового автомобиля
		car := new(models.CarCard)
		// присваивание ему названия
		car.Name = names[i]
		// присваивание ему ссылок на изображения
		car.Images = img
		// присваивание ему цены
		car.Price = Clear_the_string(prices[i])
		// присваивание ему характеристик
		car.Characteristics = characteristics[i]
		// добавление нового автомобиля в срез
		cars = append(cars, *car)
	}

	return &cars
}

// Подготовка некоторых параметров для вставки в ссылку
func Prepare_parameters(car_make, model, min_price_, max_price_, earliest_year_, latest_year_, gearbox_, fuel_, drive_ string) (string, string, map[string]int) {

	a_few_parameters := make(map[string]int)

	// перевод нижней границы цены из типа string в тип int
	a_few_parameters["min_price"], _ = strconv.Atoi(min_price_)
	// перевод верхней границы цены из типа string в тип int
	a_few_parameters["max_price"], _ = strconv.Atoi(max_price_)
	// перевод самого раннего года выпуска из типа string в тип int
	a_few_parameters["earliest_year"], _ = strconv.Atoi(earliest_year_)
	// перевод самого недавнего года выпуска из типа string в тип int
	a_few_parameters["latest_year"], _ = strconv.Atoi(latest_year_)

	switch gearbox_ {
	case "Вариатор":
		a_few_parameters["gearbox"] = 3
	case "Робот":
		a_few_parameters["gearbox"] = 4
	case "Механическая":
		a_few_parameters["gearbox"] = 1
	}

	switch fuel_ {
	case "Бензин":
		a_few_parameters["fuel"] = 1
	case "Дизель":
		a_few_parameters["fuel"] = 2
	case "Электричество":
		a_few_parameters["fuel"] = 4
	case "Гибрид":
		a_few_parameters["fuel"] = 5
	}

	switch drive_ {
	case "Передний":
		a_few_parameters["drive"] = 1
	case "Задний":
		a_few_parameters["drive"] = 2
	case "Полный":
		a_few_parameters["drive"] = 3
	}

	switch car_make {
	case "Лада":
		car_make = "lada"
	case "Mitsubishi":
		car_make = "mitsubishi"
	case "Volkswagen":
		car_make = "volkswagen"
	}

	switch model {
	case "Гранта":
		model = "granta"
	case "Приора":
		model = "priora"
	case "Веста":
		model = "vesta"
	case "Pajero":
		model = "pajero"
	case "Outlander":
		model = "outlander"
	case "L200":
		model = "l200"
	case "ASX":
		model = "asx"
	case "Pajero Sport":
		model = "pajero_sport"
	}

	return car_make, model, a_few_parameters
}

// удаление из цен знака замены � с кодом 65533
func Clear_the_string(str string) string {
	var new_str []rune
	for _, s := range str {
		if int(s) == 65533 {
			new_str = append(new_str, ' ')
			continue
		}
		new_str = append(new_str, s)
	}
	str = string(new_str)
	str += "₽"
	return str
}

// Функция, проверяющая все ли ключи, соответствующие характеристикам, содержатся в хэш-таблице
func check_characteristics(characteristics map[string]string) map[string]string {

	// набор возможных ключей
	var features []string = []string{"Двигатель", "Мощность", "Коробка передач", "Привод", "Тип кузова", "Цвет", "Пробег, км", "Руль", "Число мест", "Поколение", "Название комплектации", "Время разгона 0-100 км/ч, с", "Максимальная скорость, км/ч", "Используемое топливо", "Расход топлива в городском цикле, л/100 км", "Расход топлива за городом, л/100 км", "Расход топлива в смешанном цикле, л/100 км", "Максимальный крутящий момент, Н*м (кг*м) при об./мин.", "Максимальная мощность, л.с. (кВт) при об./мин.", "Кондиционер", "Климат-контроль", "Объем багажника, л", "Габариты кузова (Д x Ш x В), мм", "Клиренс (высота дорожного просвета), мм", "Масса, кг"}

	// если какой-либо ключ отсутствует, этот ключ добавляется со значением "Не указано"
	for _, feature := range features {
		if _, ok := characteristics[feature]; !ok {
			characteristics[feature] = "Не указано"
		}
	}
	return characteristics
}

// Функция получения ссылок на страницы, содержащие сведения о автомобилях, а также
// их названий и цен
func Scrape_pages_names_prices(link, pattern, object string) ([]string, []string, []string) {

	// объект структуры, представляющей запрос
	var response *http.Response

	// получение данной web-страницы. Установление соединения.
	document, response := Get_web_page(link)

	// прекращение соединения
	response.Body.Close()

	// получение ссылок на страницы, содержащие сведения о автомобилях
	links := Scrape_links(document, link, pattern, "pages")

	fmt.Println(links)

	// названия автомобилей
	var names = []string{}

	// цены автомобилей
	var prices = []string{}

	// ограничение количества автомобилей
	var limit int = 0

	// поиск элементов с тегом "span"
	document.Find("span").EachWithBreak(func(index int, element *goquery.Selection) bool {
		// получение значения атрибута "data-ftid"
		href, exists := element.Attr("data-ftid")
		// если значение существует
		if exists {
			if href == "bull_title" {
				// добвление нового названия
				names = append(names, element.Text())

				// реализация ограничения числа автомобилей
				limit++
				if limit == _limit_value {
					return false
				}
			}
			if href == "bull_price" {
				// добвление новой цены
				prices = append(prices, element.Text())

				// реализация ограничения числа автомобилей
				limit++
				if limit == _limit_value {
					return false
				}
			}

		}

		return true
	})

	return links, names, prices
}

// Функция получения ссылок на страницы, содержащие сведения о автомобилях,
// а также получения ссылок на фотографии одного автомобиля
func Scrape_links(document *goquery.Document, link, pattern, object string) []string {

	// объект структуры, представляющей запрос
	var response *http.Response
	// если нужно найти ссылки на фотографии одного автомобиля
	if object == "images" {
		// получение данной web-страницы. Установление соединения.
		document, response = Get_web_page(link)
		// прекращение соединения
		response.Body.Close()
	}

	// ограничение на количество автомобилей
	var limit int

	// ссылки на страницы или фотографии
	var links = []string{}

	// комиляция шаблона регулярного выражения
	re, _ := regexp.Compile(pattern)

	// поиск элементов с тегом "а"
	document.Find("a").EachWithBreak(func(index int, element *goquery.Selection) bool {
		// получение значения атрибута "href"
		href, exists := element.Attr("href")
		// если значение существует
		if exists {
			// найти все строки, удовлетворяющие шаблону
			if re.FindAllString(href, -1) != nil {
				// добавить найденную строку
				links = append(links, href)
				// реализация ограничения числа автомобилей
				if object == "pages" {
					limit++
					if limit == limit_value {
						return false
					}
				}
			}
		}
		return true
	})

	return links
}

// Функция, получения характеристик автомобиля с web-страницы, предоставляющей краткий перечень характеристик автомобиля
func Scrape_characteristics(link string) map[string]string {

	// объект структуры, представляющей запрос
	var response *http.Response

	// получение данной web-страницы. Установление соединения.
	document, response := Get_web_page(link)

	// прекращение соединения
	response.Body.Close()

	characteristics := make(map[string]string)
	// почему-то регулярные выражения плохо работают с кириллицей в go
	// поиск в строке шаблона типа "бензин, 1.5 л", поскольку иногда приходит строка вида "бензин, 1.5 лтребуется ремонт или не на ходу".
	// это связано с тем, что иногда на странице(на странице неисправной машины) встречаются два тега span с классом "css-1jygg09".
	// первый содержит "бензин, 1.5 л", а второй "требуется ремонт или не на ходу". Функция Find их склеивает.
	// шаблон (регулярного выражения), которому должна соответствовать строка, содержащая сведения о двигателе.
	re, _ := regexp.Compile(`(\W+\D\d.\d\D\W)|(\W+)`)

	// поиск элемента с тегом "span" и классом "css-1jygg09", который удовлетворяет шаблону
	characteristics["Двигатель"] = re.FindString(document.Find("span.css-1jygg09").Text())

	// обратный слэш s, обозначающий пробельный символ, не работает
	// 3 или 2 или 4 символа обозначают количество цифр в числе лошадиных сил Л.с может быть меньше ста - например 90 л.с, может быть больше 1000 - 4 цифры
	// шаблон (регулярного выражения), которому должна соответствовать строка, содержащая сведения о мощности.
	re, _ = regexp.Compile(`(\d{3}|\d{2}|\d{4})\Dл\.с\.`)

	// поиск элемента с тегом "span" и классом "css-9g0qum.e162wx9x0", который удовлетворяет шаблону
	characteristics["Мощность"] = re.FindString(document.Find("span.css-9g0qum.e162wx9x0").Text())

	// поиск элементов с тегом "td" и классом "css-9xodgi.ezjvm5n0"
	document.Find("td.css-9xodgi.ezjvm5n0").EachWithBreak(func(index int, element *goquery.Selection) bool {

		// проверка текста элемента, который находится перед текущим элементом
		switch element.Prev().Text() {
		case "Коробка передач":
			// присваивание текста текущего элемента
			characteristics["Коробка передач"] = element.Text()
		case "Привод":
			characteristics["Привод"] = element.Text()
		case "Тип кузова":
			characteristics["Тип кузова"] = element.Text()
		case "Цвет":
			characteristics["Цвет"] = element.Text()
		case "Руль":
			characteristics["Руль"] = element.Text()
		}

		return true
	})

	// поиск элемента с тегом "span" и классом "css-9g0qum.e162wx9x0"
	// css-1ni3lw9 e162wx9x0 - это новый авто, нужно добавить. что с cherry tiggo?
	characteristics["Пробег, км"] = document.Find("span.css-1osyw3j.ei6iaw00").Text()
	// если автомобиль новый
	if characteristics["Пробег, км"] == "" {
		characteristics["Пробег, км"] = document.Find("span.css-1ni3lw9 e162wx9x0").Text()
	}

	// ссылка на web-страницу, которая содержит ссылки на страницы комплектаций,
	// одна из которых подходит данному автомобилю.
	var generation_link string

	// ссылка на страницу комплектации
	var complectation_link string

	// поиск элементов с тегом "а"
	document.Find("a").EachWithBreak(func(index int, element *goquery.Selection) bool {

		// получение значения атрибута "data-ga-stats-name"
		item, exists := element.Attr("data-ga-stats-name")
		// если значение существует
		if exists {
			if item == "generation_link" {
				// присваивание значение атрибута "href"
				generation_link, _ = element.Attr("href")
				// присваивание текста текущего элемента
				characteristics["Поколение"] = element.Text()
			}
			if item == "complectation_link" {
				// присваивание значение атрибута "href"
				complectation_link, _ = element.Attr("href")
				// присваивание текста текущего элемента
				characteristics["Комплектация"] = element.Text()
			}
		}

		return true
	})

	if complectation_link != "" {
		// Копирование характеристик, полученных функцией Scrape_page_of_complectation_link в исходную хэш-таблицу
		maps.Copy(characteristics, Scrape_page_of_complectation_link(&complectation_link))
	} else if complectation_link == "" && generation_link != "" {
		// Копирование характеристик, полученных функцией Scrape_page_of_generation_link в исходную хэш-таблицу
		maps.Copy(characteristics, Scrape_page_of_generation_link(characteristics, generation_link))
	}

	return characteristics
}

// Функция получения характеристик автомобиля, использующая web-страницу, которая содержит ссылки на страницы комплектаций,
// одна из которых подходит данному автомобилю. С подходящей страницы комплектации собирается информация о характеристиках.
func Scrape_page_of_generation_link(characteristics map[string]string, generation_link string) map[string]string {

	// объект структуры, представляющей запрос
	var response *http.Response

	// получение данной web-страницы. Установление соединения.
	document, response := Get_web_page(generation_link)

	// прекращение соединения
	response.Body.Close()

	// ссылки на страницы, которые содержат ссылки на страницы конкретных комплектаций.
	var links []string

	// Добавление ссылки на текущую страницу
	links = append(links, generation_link)

	// шаблон(регулярного выражения) ссылок на последующие за первой страницы, которые содержат ссылки на
	// страницы конкретных комплектаций (на вторую, третью, и.т.д.)
	re, _ := regexp.Compile(`https://www.drom.ru/catalog/\w+/\w+/\w+/page\d\.html`)

	// поиск элементов с тегом "а"
	document.Find("a").EachWithBreak(func(index int, element *goquery.Selection) bool {
		// Получение значения атрибута "href"
		item, exists := element.Attr("href")
		// если значение существует
		if exists {
			// если ссылка соответствует шаблону регулярного выражения
			if re.MatchString(item) {
				// Добавление ссылки на следующую страницу
				links = append(links, item)
			}
		}
		return true
	})

	// удаление последней ссылки, посколку она является дубликатом ссылки внутри этого среза, если ссылок больше одной
	if len(links) != 1 {
		links = links[:len(links)-1]
	}

	// характеристики автомобиля, которые нужно для сопоставления с кратким перечнем характеристик, присущих конкретной(-ым) комплектации(-ям)
	specific_characteristics := make(map[string]string)

	// ссылка на страницу комплектации
	var complectation_link string

	for _, link := range links {

		// получение данной web-страницы. Установление соединения.
		var response *http.Response
		document, response := Get_web_page(link)

		// прекращение соединения
		response.Body.Close()

		// копирование только нужных характеристик для сопоставления с кратким перечнем характеристик, присущих конкретной(-ым) комплектации(-ям)
		// исходный хэщ-таблица должна быть неизменной. Хэш-таблица передается по ссылке и является reference type в языке go
		specific_characteristics["Двигатель"] = characteristics["Двигатель"]
		specific_characteristics["Мощность"] = characteristics["Мощность"]
		specific_characteristics["Коробка передач"] = characteristics["Коробка передач"]
		specific_characteristics["Привод"] = characteristics["Привод"]

		// удаление из строки вида "249 л.с." подстроки "л.с", поскольку regexp.MatchString() почему-то не находит
		// совпадение строки вида "249 л.с" c строкой вида ""3.5 л, 249 л.с., бензин, АКПП, полный привод (4WD)
		specific_characteristics["Мощность"] = strings.Replace(specific_characteristics["Мощность"], "л.с.", "", 1)

		// удаление пробела из строки вида "249 ", потому что функция replace выше почему-то не может убрать подстроку " л.с" из "249 л.с."
		specific_characteristics["Мощность"] = strings.TrimSpace(specific_characteristics["Мощность"])

		if specific_characteristics["Двигатель"] == "электро" {
			// замена "электро" на "электричество", поскольку только последний вариант присутствует на странице комплектаций
			specific_characteristics["Двигатель"] = "электричество"
			// удаление "АКПП", поскольку её нет на странице комплектации
			delete(specific_characteristics, "Коробка передач")
		} else {

			// удаление пробелов из таких, например, строк "бензин, 2.3 л" и "бензин, 2.0 л, гибрид"
			specific_characteristics["Двигатель"] = strings.Replace(specific_characteristics["Двигатель"], " ", "", -1)

			// разделение строки вида "бензин, 2.3 л" на две строки "бензин" и " 2.3 л"
			// может встретиться "бензин, 2.0 л, гибрид", которая будет разбита на "бензин", "2.0 л", "гибрид"
			// если встретиться "электро", то эта строка не будет разделена
			fuel := strings.Split(specific_characteristics["Двигатель"], ",")

			// удаление "л" из "2.3 л", поскольку команда выше удаляет пробел, и получается "2.3л", которая не равна строке "2.3 л", находящейся на странице комплектаций
			fuel[1] = strings.Replace(fuel[1], "л", "", 1)

			// удаление строки вида "бензин, 2.3 л"
			delete(specific_characteristics, "Двигатель")

			// добавление этих двух строк "бензин", " 2.3 л" или трех строк "бензин", "2.0 л", "гибрид"
			for _, f := range fuel {
				specific_characteristics[f] = f
			}

			// если вместо "МКПП" И "АКПП", в строке characteristics присутствуют "механика" и "автомат" соответственно, то нужно заменить их на "МКПП" и "АКПП",
			// потому что на странице комплектаций всегда встречается "МКПП" и "АКПП", и нужно чтобы characteristics[3] совпал с "МКПП" или "АКПП" на странице комплектаций
			if specific_characteristics["Коробка передач"] == "механика" {
				specific_characteristics["Коробка передач"] = "МКПП"
			} else if specific_characteristics["Коробка передач"] == "автомат" {
				specific_characteristics["Коробка передач"] = "АКПП"
			}
		}

		// поиск элементов с тегом "th"
		document.Find("th").EachWithBreak(func(index int, element *goquery.Selection) bool {

			// Получение значения атрибута "colspan"
			item, exists := element.Attr("colspan")
			if exists {
				// 7 и 6 проверяются, потому что так была составлена HTML-страница
				if item == "7" || item == "6" {

					// счетчик совпадений краткого перечня характеристик, присущих конкретной(-ым) комплектации(-ям) с хэш-таблицей specific_characteristics
					match := 0
					for _, value := range specific_characteristics {
						// поиск в строке(тексте элемента с тегом "th"и значением атрибута "colspan" = 7 или 6)) подстроки, удовлетворяющей шаблону регулярного выражения
						res, _ := regexp.MatchString(value, element.Text())
						// если совпадение найдено
						if res {
							match++
						}
					}

					// если совпали все характеристики
					if match == len(specific_characteristics) {
						// Найти ссылку на страницу комплектацию, которой присущ данный краткий перечень характеристик
						// ссылка - это значение атрибута "href" элемента с тегом "a", который является потомком элемента,
						// следующего за родителем текущего элемента
						complectation_link, _ = element.Parent().Next().Children().Find("a").Attr("href")
						// выйти из поиска элементов с тегом "th"
						return false
					}
				}
			}
			return true
		})

		// если ссылка на комплектацию найдена, нужно выйти из цикла
		if complectation_link != "" {
			break
		}
	}

	complectation_link = "https://www.drom.ru" + complectation_link

	// вызов функции, определенной ниже
	return Scrape_page_of_complectation_link(&complectation_link)
}

// Функция получения характеристик автомобиля с web-страницы, содержащей сведения о комплектации автомобиля
func Scrape_page_of_complectation_link(complectation_link *string) map[string]string {

	// объект структуры, представляющей запрос
	var response *http.Response

	// получение данной web-страницы. Установление соединения.
	document, response := Get_web_page(*complectation_link)

	// прекращение соединения
	response.Body.Close()

	// характеристики автомобиля, которые будут получены с данной web-страницы
	characteristics_from_complectation_link := make(map[string]string)

	// Поиск элементов с тегом "td"
	document.Find("td").EachWithBreak(func(index int, element *goquery.Selection) bool {

		// проверка текста элемента с тегом "td"
		switch element.Text() {
		case "Название комплектации":
			// присваивание текста элемента, следующего за элементом с тегом "td"
			// удаление из текста пробелов и символов новой строки
			characteristics_from_complectation_link["Название комплектации"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Время разгона 0-100 км/ч, с":
			characteristics_from_complectation_link["Время разгона 0-100 км/ч, с"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Максимальная скорость, км/ч":
			characteristics_from_complectation_link["Максимальная скорость, км/ч"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Используемое топливо":
			characteristics_from_complectation_link["Используемое топливо"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Расход топлива в городском цикле, л/100 км":
			characteristics_from_complectation_link["Расход топлива в городском цикле, л/100 км"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Расход топлива за городом, л/100 км":
			characteristics_from_complectation_link["Расход топлива за городом, л/100 км"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Расход топлива в смешанном цикле, л/100 км":
			characteristics_from_complectation_link["Расход топлива в смешанном цикле, л/100 км"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Максимальный крутящий момент, Н*м (кг*м) при об./мин.":
			characteristics_from_complectation_link["Максимальный крутящий момент, Н*м (кг*м) при об./мин."] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Максимальная мощность, л.с. (кВт) при об./мин.":
			characteristics_from_complectation_link["Максимальная мощность, л.с. (кВт) при об./мин."] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Кондиционер":
			// проверка значения аттрибута "href" потомков потомков элемента,следующего за элементом с тегом "td"
			if s, _ := element.Next().Children().Children().Attr("href"); s == "#yes" {
				characteristics_from_complectation_link["Кондиционер"] = "есть"

			} else if s == "#option" {
				characteristics_from_complectation_link["Кондиционер"] = "опция"

			} else {
				characteristics_from_complectation_link["Кондиционер"] = "нет"
			}

		case "Климат-контроль":
			// проверка значения аттрибута "href" потомков потомков элемента,следующего за элементом с тегом "td"
			if s, _ := element.Next().Children().Children().Attr("href"); s == "#yes" {
				characteristics_from_complectation_link["Климат-контроль"] = "есть"
			} else if s == "#option" {
				characteristics_from_complectation_link["Климат-контроль"] = "опция"
			} else {
				characteristics_from_complectation_link["Климат-контроль"] = "нет"
			}

		case "Объем багажника, л":
			// присваивание текста элемента, следующего за элементом с тегом "td"
			// удаление из текста пробелов и символов новой строки
			characteristics_from_complectation_link["Объем багажника, л"] = strings.Replace(strings.Replace(element.Next().Text(), " ", "", -1), "\n", "", -1)

		case "Габариты кузова (Д x Ш x В), мм":
			characteristics_from_complectation_link["Габариты кузова (Д x Ш x В), мм"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Клиренс (высота дорожного просвета), мм":
			characteristics_from_complectation_link["Клиренс (высота дорожного просвета), мм"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)

		case "Масса, кг":
			characteristics_from_complectation_link["Масса, кг"] = strings.Replace(strings.Trim(element.Next().Text(), " "), "\n", "", -1)
		}

		return true
	})

	return characteristics_from_complectation_link
}

// Функция получения web-страницы
func Get_web_page(link string) (*goquery.Document, *http.Response) {

	// получение страницы. Установление соединения. Посылается GET-запрос.
	response, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}

	// смена кодировки страницы с windows-1251 на utf-8
	utfBody, err := iconv.NewReader(response.Body, "windows-1251", "utf-8")
	if err != nil {
		log.Fatal("Error converting charset from windows-1251 to utf-8. ", err)
	}

	// создание объекта структуры, представляющего HTML документ
	document, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	return document, response
}

func (sr *searchRepository) GetFileName() string {
	return "offer_for_search.html"
}

func (sr *searchRepository) GetQuestion() string {
	return "Как Вы думаете, расход топлива в смешанном цикле <u>5 л/100 км</u>:"
}

// GetFileName() string
// 	GetQuestion() string
