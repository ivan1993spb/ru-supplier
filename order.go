package main

import (
	"bytes"
	"fmt"
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
	FZ44 OrderLaw = iota
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
	return "??-ФЗ"
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

func NewOrder(law_id, order_id, order_type, order_name,
	exhibition_number, exhibition_name, start_order_price,
	currency_id, okdp, okpd, organisation_name, pub_date,
	last_event_date, order_stage, features, start_diling_date,
	finish_diling_date string) (order *Order) {

	order = &Order{-1, "", order_type, order_name,
		0, exhibition_name, 0, currency_id, okdp, okpd,
		organisation_name, pub_date, last_event_date,
		order_stage, features, "", "", nil}

	var err error
	order.LawId, err = ParseLow(law_id)
	if err != nil {
		order.PushError(_INVALID_LAW_ID)
		err = nil
	}
	order.OrderId = strings.TrimLeft(order_id, "№")
	if len(exhibition_number) > 0 {
		// Только для многолотовых закупок
		order.ExhibitionNumber, err = strconv.Atoi(exhibition_number)
		if err != nil {
			order.PushError(_INVALID_EXHIBITION_NUMBER)
			err = nil
		}
	}
	order.StartOrderPrice, err = ParsePrice(start_order_price)
	if err != nil {
		order.PushError(_INVALID_START_ORDER_PRICE)
		err = nil
	}
	if len(order.CurrencyId) == 0 {
		order.PushError(_UNKNOWN_CURRENCY)
	}
	if len(start_diling_date) == 0 {
		// Когда по какой-то причине в csv файле отсутствует дата
		// начала приема заявок
		// назначаем дату последнего события и выводим ошибку
		order.StartDilingDate = last_event_date
		order.PushError(_UNKNOWN_START_FILING_DATE)
	} else {
		order.StartDilingDate = start_diling_date
	}
	if len(finish_diling_date) == 0 {
		// отсутствует дата окончания приема заявок
		order.PushError(_UNKNOWN_FINISH_FILING_DATE)
	} else {
		order.FinishDilingDate = finish_diling_date
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
		"?placeOfSearch=FZ_44&_placeOfSearch=on",
		"&placeOfSearch=FZ_223&_placeOfSearch=on",
		"&placeOfSearch=FZ_94&_placeOfSearch=on",
		"&priceFrom=0&priceTo=200+000+000+000",
		"&publishDateFrom=&publishDateTo=",
		"&updateDateFrom=&updateDateTo=",
		"&orderStages=AF&_orderStages=on",
		"&orderStages=CA&_orderStages=on",
		"&orderStages=PC&_orderStages=on",
		"&orderStages=PA&_orderStages=on",
		"&sortDirection=false&sortBy=UPDATE_DATE",
		"&recordsPerPage=_10&pageNo=1",
		"&searchString=", order_id,
		"&strictEqual=false&morphology=false",
		"&showLotsInfo=false&isPaging=false",
		"&isHeaderClick=&checkIds=")
}
