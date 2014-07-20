package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

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
	return -1, errors.New("Invalid or unknown law id")
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
		return 0, errors.New("Invalid order price")
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
	return Price(price / 100), nil
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

const (
	_FIELD_LAW_ID int = iota
	_FIELD_ORDER_ID
	_FIELD_ORDER_TYPE
	_FIELD_ORDER_NAME
	_FIELD_EXHIBITION_NUMBER
	_FIELD_EXHIBITION_NAME
	_FIELD_START_ORDER_PRICE
	_FIELD_CURRENCY_ID
	_FIELD_OKDP
	_FIELD_OKPD
	_FIELD_ORGANISATION_NAME
	_FIELD_PUB_DATE
	_FIELD_LAST_EVENT_DATE
	_FIELD_ORDER_STAGE
	_FIELD_FEATURES
	_FIELD_START_FILING_DATE
	_FIELD_FINISH_FILING_DATE
	_ORDER_COLUMN_COUNT // result column count
)

type Order struct {
	LawId            OrderLaw  // Номер ФЗ
	OrderId          string    // Реестровый номер закупки
	OrderType        string    // Способ размещения закупки
	OrderName        string    // Наименование закупки
	ExhibitionNumber int       // Номер лота
	ExhibitionName   string    // Наименование лота
	StartOrderPrice  Price     // Начальная (максимальная) цена
	CurrencyId       string    // Код валюты
	OKDP             string    // Классификация по ОКДП
	OKPD             string    // Классификация по ОКПД
	OrganisationName string    // Организация, размещающая заказ
	PubDate          time.Time // Дата публикации
	LastEventDate    time.Time // Дата последнего события
	OrderStage       string    // Этап закупки (размещения заказа)
	Features         string    // Особенности размещения заказа
	StartFilingDate  time.Time // Дата начала подачи заявок
	FinishFilingDate time.Time // Дата окончания подачи заявок
	Errors           []error   // Ошибки при анализе закупки
}

func NewOrder(row [_ORDER_COLUMN_COUNT]string) (order *Order) {
	order = &Order{
		OrderId:          strings.TrimLeft(row[_FIELD_ORDER_ID], "№"),
		OrderType:        row[_FIELD_ORDER_TYPE],
		OrderName:        row[_FIELD_ORDER_NAME],
		ExhibitionName:   row[_FIELD_EXHIBITION_NAME],
		CurrencyId:       row[_FIELD_CURRENCY_ID],
		OKDP:             row[_FIELD_OKDP],
		OKPD:             row[_FIELD_OKPD],
		OrganisationName: row[_FIELD_ORGANISATION_NAME],
		OrderStage:       row[_FIELD_ORDER_STAGE],
		Features:         row[_FIELD_FEATURES],
	}
	var err error
	order.LawId, err = ParseLow(row[_FIELD_LAW_ID])
	if err != nil {
		order.PushError(err)
		err = nil
	}
	if len(row[_FIELD_EXHIBITION_NUMBER]) > 0 {
		// Только для многолотовых закупок по ФЗ 223
		order.ExhibitionNumber, err = strconv.Atoi(row[_FIELD_EXHIBITION_NUMBER])
		if err != nil {
			order.PushError(errors.New("Invalid exhibition number"))
			err = nil
		}
	}
	order.StartOrderPrice, err = ParsePrice(row[_FIELD_START_ORDER_PRICE])
	if err != nil {
		order.PushError(err)
		err = nil
	}
	if len(row[_FIELD_CURRENCY_ID]) == 0 {
		order.PushError(errors.New("Unknown currency"))
	}
	order.PubDate, err = ParseRusFormatDate(row[_FIELD_PUB_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown publish date"))
		err = nil
	}
	order.LastEventDate, err = ParseRusFormatDate(row[_FIELD_LAST_EVENT_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown last event date"))
		err = nil
	}
	order.StartFilingDate, err = ParseRusFormatDate(row[_FIELD_START_FILING_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown start filing date"))
		err = nil
	}
	order.FinishFilingDate, err = ParseRusFormatDate(row[_FIELD_FINISH_FILING_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown finish filing date"))
		err = nil
	}
	return
}

func (order *Order) PushError(err error) {
	if err != nil {
		order.Errors = append(order.Errors, err)
	}
}
