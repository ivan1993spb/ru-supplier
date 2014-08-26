package main

import (
	"bytes"
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

func ParseOrder(rowBytes []byte) (*Order, error) {
	if len(rowBytes) == 0 {
		return nil, errors.New("Passed empty bytes row")
	}

	var (
		// Array of order fields
		row   [_ORDER_COLUMN_COUNT]string
		field int // index of current order field

		block  []byte // block
		quoted bool   // true if current block is quoted

		// Tokenes:
		//     i is index of `;`
		//     j is index of `"`
		//     k is index of `""`
		i, j, k int

		l int // parser position in rowBytes
	)

	for {
		i = bytes.IndexByte(rowBytes[l:], ';')
		j = bytes.IndexByte(rowBytes[l:], '"')
		k = bytes.Index(rowBytes[l:], []byte{'"', '"'})

		if (i > -1) && (i < j || j < 0) && (i < k || k < 0) {
			// next token is `;`
			if quoted {
				// `;` is part of block
				block = append(block, rowBytes[l:i+1]...)
			} else {
				// `;` is separator
				block = append(block, rowBytes[l:i]...)
				// write block
				row[field] = string(block)
				field++
			}
			l += i + 1

		} else if (j > -1) && (j < i || i < 0) && (j < k || k < 0) {
			// next token is `"`
			if quoted {
				block = append(block, rowBytes[l:j]...)
				row[field] = string(block)
				field++
			}
			l += j + 1
			quoted = !quoted

		} else if quoted && (k > -1) && (k < i || i < 0) && k == j {
			// next token is `""`
			block = append(block, rowBytes[l:k]...)
			// append only one double quote
			block = append(block, '"')
			l += k + 2
			row[field] = string(block)
			field++

		} else if i < 0 && j < 0 && k < 0 {
			block = append(block, rowBytes[l:])

		}
	}

	return NewOrder(row), nil
}

func ParseOrder2(rowBytes []byte) {
	if len(rowBytes) == 0 {
		fmt.Println("empty")
	}

	var (
		// Array of order fields
		// row   [_ORDER_COLUMN_COUNT]string
		// field int // index of current order field

		block  []byte // block
		quoted bool   // true if current block is quoted
		inside = true

		l int // parser cursor

		// Tokenes:
		//     l+i is current index of `;`
		//     l+j is current index of `"`
		//     l+k is current index of `""`
		i, j, k int
	)

	for l < len(rowBytes) {
		i = bytes.IndexByte(rowBytes[l:], ';')
		j = bytes.IndexByte(rowBytes[l:], '"')
		k = bytes.Index(rowBytes[l:], []byte{'"', '"'})

		if (i > -1) && (i < j || j < 0) && (i < k || k < 0) {
			// next token is `;`
			if quoted {
				// `;` is part of block
				block = append(block, rowBytes[l:l+i+1]...)
			} else {
				// `;` is separator
				block = append(block, rowBytes[l:l+i]...)
				// close block
				inside = false
			}
			l += i + 1

		} else if (j > -1) && (j < i || i < 0) && (j < k || k < 0 || !quoted) {
			// next token is `"`
			// if field are quoted
			if quoted {
				block = append(block, rowBytes[l:l+j]...)
				// ignore all bytes before next seporator `;`
				if i > j {
					j = i
				}
				// close block
				inside = false // ????????
			}
			l += j + 1
			quoted = !quoted

		} else if quoted && (k > -1) && (k < i || i < 0) && k == j {
			// next token is `""`
			// append only one double quote
			block = append(block, rowBytes[l:l+k+1]...)
			// but skip two
			l += k + 2
			continue

		} else if i < 0 && j < 0 && k < 0 {
			if quoted {
				fmt.Println("error")
				break
			}
			if inside {
				block = append(block, rowBytes[l:]...)
				inside = false
			} else {
				fmt.Println("err 2")
			}
			//fmt.Println("4:", string())
			// block = append(block, rowBytes[l:])
			l = len(rowBytes)
		}

		if !inside {
			// fmt.Println(inside)
			fmt.Printf("%q\n", block)
			block = block[:0]
			inside = true
		}
		// fmt.Println(inside)
	}
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
