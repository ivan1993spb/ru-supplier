package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

const (
	_DEFAULT_TITLE    = "Закупки"
	_FEED_DESCRIPTION = "Лента закупок с гибким механизмом фильтрации"
	_FEED_LINK        = "http://zakupki.gov.ru"
)

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
			{{if .StartFilingDate}}{{.StartFilingDate}}
			{{else}}00.00.0000{{end}}
			по
			<s>
				{{if .FinishFilingDate}}{{.FinishFilingDate}}
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

func MakeTitle(order *Order) (title string) {
	title = "№" + order.OrderId
	if order.ExhibitionNumber > 0 {
		title += " Лот " + strconv.Itoa(order.ExhibitionNumber)
	}
	return
}

func MakeDescription(order *Order) string {
	buff := bytes.NewBuffer(nil)
	err := tmpl.Execute(buff, map[string]interface{}{
		"Title":            MakeTitle(order),
		"LawId":            LawIdToString(order.LawId),
		"Link":             MakeLink(order.OrderId),
		"OrderName":        order.OrderName,
		"OKDP":             order.OKDP,
		"OKPD":             order.OKPD,
		"StartFilingDate":  RusFormatDate(order.StartFilingDate),
		"FinishFilingDate": RusFormatDate(order.FinishFilingDate),
		"StartOrderPrice":  FormatPrice(order.StartOrderPrice),
		"CurrencyId":       order.CurrencyId,
		"OrderType":        order.OrderType,
		"OrderStage":       order.OrderStage,
		"PubDate":          RusFormatDate(order.PubDate),
		"OrganisationName": order.OrganisationName,
		"Features":         order.Features,
		"Errors":           order.Errors,
	})
	if err != nil {
		log.Println("template execution error:", err)
	}
	return buff.String()
}

func MakeLink(id string) string {
	return fmt.Sprint("http://zakupki.gov.ru",
		"/epz/order/quicksearch/update.html",
		"?placeOfSearch=FZ_44&_placeOfSearch=on",  // ФЗ 44
		"&placeOfSearch=FZ_223&_placeOfSearch=on", // ФЗ 223
		"&placeOfSearch=FZ_94&_placeOfSearch=on",  // ФЗ 94
		"&priceFrom=0&priceTo=200+000+000+000",    // любая НМЦК
		"&publishDateFrom=&publishDateTo=",        // выкл. диапозоны
		"&updateDateFrom=&updateDateTo=",          // времени
		"&orderStages=AF&_orderStages=on",         // подача заявок
		"&orderStages=CA&_orderStages=on",         // работа комиссии
		"&orderStages=PC&_orderStages=on",         // завершена
		"&orderStages=PA&_orderStages=on",         // отменена
		"&sortDirection=false&sortBy=UPDATE_DATE", // по убыванию даты обновления
		"&recordsPerPage=_10&pageNo=1",            // без этого не работает
		"&searchString=", id,                      // поиск по ид
		"&strictEqual=false&morphology=false",
		"&showLotsInfo=false&isPaging=false",
		"&isHeaderClick=&checkIds=")
}

func MakeShortLink(id, host string) string {
	return fmt.Sprintf("http://%s/%s?order=%s", host,
		strings.TrimLeft(_PATH_TO_SHORT_LINKS, "/"), id)
}

func LawIdToString(law OrderLaw) string {
	switch law {
	case FZ44:
		return "44-ФЗ"
	case FZ223:
		return "223-ФЗ"
	case FZ94:
		return "94-ФЗ"
	}
	return ""
}

func FormatPrice(p Price) string {
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

type Render struct {
	config *Config
	feed   *feeds.RssFeed
}

func NewRender(config *Config) *Render {
	return &Render{
		config,
		&feeds.RssFeed{
			Link:        _FEED_LINK,
			Description: _FEED_DESCRIPTION,
		},
	}
}

func (r *Render) Compose(title string, orders []*Order) {
	if len(title) == 0 {
		r.feed.Title = _DEFAULT_TITLE
	} else {
		r.feed.Title = title
	}
	r.feed.Items = make([]*feeds.RssItem, len(orders))
	for i, order := range orders {
		r.feed.Items[i] = &feeds.RssItem{
			Title: MakeTitle(order),
			Link: MakeShortLink(
				order.OrderId,
				r.config.HTTPHost(),
			),
			Description: MakeDescription(order),
			Author:      order.OrganisationName,
			PubDate:     order.PubDate.Format(time.RFC1123),
		}
	}
}

func (r *Render) WriteTo(w io.Writer) error {
	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}
	// clear feed
	defer func() {
		r.feed.Title = _DEFAULT_TITLE
		r.feed.Items = nil
	}()
	// write data
	return xml.NewEncoder(w).Encode(r.feed.FeedXml())
}
