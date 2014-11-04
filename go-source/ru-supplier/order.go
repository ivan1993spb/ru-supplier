package main

import (
	"errors"
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
	return Price(price / 100), nil
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

const _CSV_FIELD_SEPARATOR = ';'

func ParseOrder(rowbyte []byte) (*Order, error) {
	if len(rowbyte) == 0 {
		return nil, errors.New("Passed empty rowbyte")
	}

	for i := 0; i < len(rowbyte); i++ {
		if rowbyte[i] == '\n' {
			rowbyte = rowbyte[:i]
		}
	}
	rowbyte = append(rowbyte, ';')

	var (
		row   [_ORDER_COLUMN_COUNT]string // array of order fields
		field int                         // index of current field

		buff   []byte // current block
		allow  = true // true if character is expected
		quoted bool   // true if current block is quoted
	)

	for i := 0; i < len(rowbyte) && field < _ORDER_COLUMN_COUNT; i++ {
		if rowbyte[i] == _CSV_FIELD_SEPARATOR {
			if quoted {
				// quoted separator is part of field
				buff = append(buff, rowbyte[i])
			} else {
				// separator is end of string
				// save
				row[field] = string(buff)
				field++
				// remove
				buff = buff[:0]
				allow = true
			}
		} else if allow {
			if rowbyte[i] == '"' {
				if quoted {
					if i+1 < len(rowbyte) && rowbyte[i+1] == '"' {
						// append only one double quote
						buff = append(buff, '"')
						// skip second quote
						i += 1
					} else {
						// closing quote
						quoted = false
						allow = false
					}
				} else if len(buff) == 0 {
					// new quoted field
					// opening quote
					quoted = true
				} else {
					// quote inside unquoted field
					return nil, errors.New("Unexpected quote")
				}
			} else {
				// append any char
				buff = append(buff, rowbyte[i])
			}
		} else {
			// block is closed
			return nil, errors.New("Unexpected character")
		}
	}

	if field < _ORDER_COLUMN_COUNT-1 {
		return nil, errors.New("Bad field count")
	}
	return NewOrder(row), nil
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
	}
	if len(row[_FIELD_EXHIBITION_NUMBER]) > 0 {
		// Только для многолотовых закупок по ФЗ 223
		order.ExhibitionNumber, err =
			strconv.Atoi(row[_FIELD_EXHIBITION_NUMBER])
		if err != nil {
			order.PushError(errors.New("Invalid exhibition number: " +
				err.Error()))
		}
	}
	order.StartOrderPrice, err =
		ParsePrice(row[_FIELD_START_ORDER_PRICE])
	if err != nil {
		order.PushError(errors.New("Invalid order price: " +
			err.Error()))
	}
	if len(row[_FIELD_CURRENCY_ID]) == 0 {
		order.PushError(errors.New("Unknown currency"))
	}
	order.PubDate, err = ParseRusFormatDate(row[_FIELD_PUB_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown publish date: " +
			err.Error()))
	}
	order.LastEventDate, err =
		ParseRusFormatDate(row[_FIELD_LAST_EVENT_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown last event date: " +
			err.Error()))
	}
	order.StartFilingDate, err =
		ParseRusFormatDate(row[_FIELD_START_FILING_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown start filing date: " +
			err.Error()))
	}
	order.FinishFilingDate, err =
		ParseRusFormatDate(row[_FIELD_FINISH_FILING_DATE])
	if err != nil {
		order.PushError(errors.New("Unknown finish filing date: " +
			err.Error()))
	}
	return
}

func (order *Order) PushError(err error) {
	if err != nil {
		order.Errors = append(order.Errors, err)
	}
}
