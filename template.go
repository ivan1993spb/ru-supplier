package main

import "html/template"

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
			№{{.OrderId}}
			{{if .ExhibitionNumber}}
				Лот {{.ExhibitionNumber}}
			{{end}}
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
