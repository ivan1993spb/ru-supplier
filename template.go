package main

import "html/template"

var tmpl = template.Must(template.New("tmpl").Parse(`
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
<h1>
	<b>{{.LawId}}</b>
	№{{.OrderId}}
	{{if .ExhibitionNumber}}
		Лот {{.ExhibitionNumber}}
	{{end}}
</h1>
<div><a href="{{.Link}}">{{.OrderName | "unknown name"}}</a></div>
{{if .OKDP}}
	<div><b>ОКДП:</b> {{.OKDP}}</div>
{{end}}
{{if .OKPD}}
	<div><b>ОКПД:</b> {{.OKPD}}</div>
{{end}}
<div>
	<b>Сроки подачи заявки:</b>
	с {{.StartDilingDate | "00.00.0000"}}
	по <s>{{.FinishDilingDate | "00.00.0000"}}</s>
</div>
<div>
	<b>Начальная (максимальная) цена:</b>
	{{.StartOrderPrice}}
	{{.CurrencyId | "unknown currency"}}
</div>
<hr />
<div><b>Тип закупки:</b> {{.OrderType}}</div>
<div><b>Этап закупки:</b> {{.OrderStage}}</div>
<div><b>Дата публикации извещения:</b> {{.PubDate}}</div>
<div><b>Организация:</b> {{.OrganisationName}}</div>
{{if .Features}}
	<div><i>{{.Features}}</i></div>
{{end}}
{{if .Errors}}
	<hr />
	<div>
		Ошибки при анализе закупки <s>(проверьте извещение)</s>:
	</div>
	<ul>
		{{range .Errors}}
			<li>{{.}}</li>
		{{end}}
	</ul>
{{end}}
`))