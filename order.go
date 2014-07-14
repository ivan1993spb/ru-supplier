package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

type ErrParsing uint8

const (
	_INVALID_LAW_ID ErrParsing = iota
	_INVALID_EXHIBITION_NUMBER
	_INVALID_START_ORDER_PRICE
	_UNKNOWN_CURRENCY
	_UNKNOWN_START_FILING_DATE
	_UNKNOWN_FINISH_FILING_DATE
)

func (err ErrParsing) Error() string {
	switch err {
	case _INVALID_LAW_ID:
		return "Invalid or unknown law id"
	case _INVALID_EXHIBITION_NUMBER:
		return "Invalid exhibition number"
	case _INVALID_START_ORDER_PRICE:
		return "Invalid start order price"
	case _UNKNOWN_CURRENCY:
		return "Unknown currency"
	case _UNKNOWN_START_FILING_DATE:
		return "Unknown start filing date"
	case _UNKNOWN_FINISH_FILING_DATE:
		return "Unknown finish filing date"
	}
	return "Unknown parsing error"
}

type OrderLaw int

const (
	// to correct template compilation it cannot starts from zero
	FZ44 OrderLaw = iota + 1
	FZ223
	FZ94
)

func ParseLow(str string) (OrderLaw, error) {
	switch {
	case strings.Contains(str, "44"):
		return FZ44, nil
	case strings.Contains(str, "223"):
		return FZ223, nil
	case strings.Contains(str, "94"):
		return FZ94, nil
	}
	return -1, _INVALID_LAW_ID
}

func (l OrderLaw) String() string {
	switch l {
	case FZ44:
		return "44-ФЗ"
	case FZ223:
		return "223-ФЗ"
	case FZ94:
		return "94-ФЗ"
	}
	return ""
}

type Price float64

func ParsePrice(str string) (Price, error) {
	price, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	// round order price with kopeika
	if price < 0 {
		price = -price
	}
	price = price * 100
	if math.Mod(price, 1) >= 0.5 {
		price = math.Ceil(price)
	} else {
		price = math.Floor(price)
	}
	return Price(price / 100), err
}

func (p Price) String() string {
	price := float64(p)
	kop := math.Mod(price*100, 100)
	output := fmt.Sprintf(",%02.f", kop)
	price -= kop / 100
	if price > 0 {
		for price > 0 {
			chunk := math.Mod(price, 1000)
			price -= chunk
			price /= 1000
			output = fmt.Sprintf(" %03.f%s", chunk, output)
		}
		return strings.TrimLeft(output, "0 ")
	}
	return "0" + output
}

type Order struct {
	LawId            OrderLaw // Номер ФЗ
	OrderId          string   // Реестровый номер закупки
	OrderType        string   // Способ размещения закупки
	OrderName        string   // Наименование закупки
	ExhibitionNumber int      // Номер лота
	ExhibitionName   string   // Наименование лота
	StartOrderPrice  Price    // Начальная (максимальная)
	CurrencyId       string   // Код валюты
	OKDP             string   // Классификация по ОКДП
	OKPD             string   // Классификация по ОКПД
	OrganisationName string   // Организация, размещающая заказ
	PubDate          string   // Дата публикации
	LastEventDate    string   // Дата последнего события
	OrderStage       string   // Этап закупки (размещения заказа)
	Features         string   // Особенности размещения заказа
	StartDilingDate  string   // Дата начала подачи заявок
	FinishDilingDate string   // Дата окончания подачи заявок
	Errors           []error  // Ошибки при анализе закупки
}

func NewOrder(row [17]string) (order *Order) {
	order = &Order{-1, "", row[2], row[3], 0, row[5], 0, row[7],
		row[8], row[9], row[10], row[11], row[12], row[13], row[14],
		"", "", nil}

	var err error
	order.LawId, err = ParseLow(row[0])
	if err != nil {
		order.PushError(_INVALID_LAW_ID)
		err = nil
	}
	order.OrderId = strings.TrimLeft(row[1], "№")
	if len(row[4]) > 0 {
		// Только для многолотовых закупок фз 223
		order.ExhibitionNumber, err = strconv.Atoi(row[4])
		if err != nil {
			order.PushError(_INVALID_EXHIBITION_NUMBER)
			err = nil
		}
	}
	order.StartOrderPrice, err = ParsePrice(row[6])
	if err != nil {
		order.PushError(_INVALID_START_ORDER_PRICE)
		err = nil
	}
	if len(order.CurrencyId) == 0 {
		order.PushError(_UNKNOWN_CURRENCY)
	}
	if len(row[15]) == 0 {
		// Когда по какой-то причине в csv файле отсутствует дата
		// начала приема заявок
		// назначаем дату последнего события и выводим ошибку
		order.StartDilingDate = row[12]
		order.PushError(_UNKNOWN_START_FILING_DATE)
	} else {
		order.StartDilingDate = row[15]
	}
	if len(row[16]) == 0 {
		// отсутствует дата окончания приема заявок
		order.PushError(_UNKNOWN_FINISH_FILING_DATE)
	} else {
		order.FinishDilingDate = row[16]
	}
	return
}

func (order *Order) PushError(err error) {
	if err != nil {
		order.Errors = append(order.Errors, err)
	}
}

func (order *Order) Title() (title string) {
	title = "№" + order.OrderId
	if order.ExhibitionNumber > 0 {
		title += " Лот " + strconv.Itoa(order.ExhibitionNumber)
	}
	return
}

func (order *Order) Link() string {
	return MakeLink(order.OrderId)
}

func (order *Order) ShortLink() string {
	host, _, _ := net.SplitHostPort(_LOCAL_ADDR)
	if len(host) == 0 {
		host = "localhost"
	}
	return fmt.Sprintf("http://%s/%s?order=%s", host,
		strings.TrimLeft(_PATH_TO_SHORT_LINKS, "/"), order.OrderId)
}

func (order *Order) Description() string {
	buff := bytes.NewBuffer(nil)
	err := tmpl.Execute(buff, order)
	if err != nil {
		log.Error.Println("template execution error:", err)
		return ""
	}
	return buff.String()
}

func (order *Order) PubDateRFC1123() string {
	chunks := strings.SplitN(order.PubDate, ".", 3)
	if len(chunks) != 3 {
		// Bad date format
		return order.PubDate
	}
	day, err := strconv.Atoi(chunks[0])
	if err != nil {
		return order.PubDate
	}
	month, err := strconv.Atoi(chunks[1])
	if err != nil {
		return order.PubDate
	}
	year, err := strconv.Atoi(chunks[2])
	if err != nil {
		return order.PubDate
	}
	return time.Date(year, time.Month(month), day,
		0, 0, 0, 0, time.Local).Format(time.RFC1123)
}

func MakeLink(order_id string) string {
	return fmt.Sprint("http://zakupki.gov.ru",
		"/epz/order/quicksearch/update.html",
		"?placeOfSearch=FZ_44&_placeOfSearch=on",  // ФЗ 44;
		"&placeOfSearch=FZ_223&_placeOfSearch=on", // ФЗ 223;
		"&placeOfSearch=FZ_94&_placeOfSearch=on",  // ФЗ 94;
		"&priceFrom=0&priceTo=200+000+000+000",    // любая НМЦК;
		"&publishDateFrom=&publishDateTo=",        // выкл. диапозоны
		"&updateDateFrom=&updateDateTo=",          // времени;
		"&orderStages=AF&_orderStages=on",         // подача заявок;
		"&orderStages=CA&_orderStages=on",         // работа комиссии;
		"&orderStages=PC&_orderStages=on",         // завершена;
		"&orderStages=PA&_orderStages=on",         // отменена;
		"&sortDirection=false&sortBy=UPDATE_DATE", // по убыванию...
		"&recordsPerPage=_10&pageNo=1",            // даты публикации;
		"&searchString=", order_id,                // поиск по ид;
		"&strictEqual=false&morphology=false",
		"&showLotsInfo=false&isPaging=false",
		"&isHeaderClick=&checkIds=")
}

var tmpl = template.Must(template.New("tmpl").Parse(`<!DOCTYPE html>
<html>
	<head>
		<title>{{.Title}}</title>
		<style>
			div, h1 {margin: 10px 0px;}
			ul {margin: 10px 15px;}
			h1 {font-size: 15pt;}
			a, s {text-decoration: none;}
			a {color: #000;}
			a:hover {background-color: #444; color: #fff;}
			s {color: #f00;}
			b {color: #999; margin-right: 8px;}
			i {color: #89f;}
		</style>
	</head>
	<body>
		<h1>
			<b>{{if .LawId}}{{.LawId}}{{else}}??-ФЗ{{end}}</b>
			{{.Title}}
		</h1>
		<div>
			<a href="{{.Link}}">
				{{if .OrderName}}{{.OrderName}}{{else}}unknown{{end}}
			</a>
		</div>
		{{if .OKDP}}
			<div><b>ОКДП:</b> {{.OKDP}}</div>
		{{end}}
		{{if .OKPD}}
			<div><b>ОКПД:</b> {{.OKPD}}</div>
		{{end}}
		<div>
			<b>Сроки подачи заявки:</b>
			с
			{{if .StartDilingDate}}{{.StartDilingDate}}
			{{else}}00.00.0000{{end}}
			по
			<s>
				{{if .FinishDilingDate}}{{.FinishDilingDate}}
				{{else}}00.00.0000{{end}}
			</s>
		</div>
		<div>
			<b>Начальная (максимальная) цена:</b>
			{{.StartOrderPrice}}
			{{if .CurrencyId}}{{.CurrencyId}}
			{{else}}unknown currency{{end}}
		</div>
		<hr />
		{{if .OrderType}}
			<div><b>Тип закупки:</b> {{.OrderType}}</div>
		{{end}}
		{{if .OrderStage}}
			<div><b>Этап закупки:</b> {{.OrderStage}}</div>
		{{end}}
		{{if .PubDate}}
			<div><b>Дата публикации извещения:</b> {{.PubDate}}</div>
		{{end}}
		{{if .OrganisationName}}
			<div><b>Организация:</b> {{.OrganisationName}}</div>
		{{end}}
		{{if .Features}}
			<div><i>{{.Features}}</i></div>
		{{end}}
		{{if .Errors}}
			<hr />
			<div>
				<s>Проверьте извещение</s>, были обнаружены ошибки:
			</div>
			<ul>
				{{range .Errors}}
					<li>{{.}}</li>
				{{end}}
			</ul>
		{{end}}
	</body>
</html>`))
