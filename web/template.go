package web

import (
	"html/template"
	"log"
)

var tmpl = template.New("")

func init() {
	var err error
	tmpl, err = tmpl.Parse(`
{{define "main"}}
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>HGRABER</title>
  </head>
  <body>
	<style>
		body{
			text-align: center;
		}
		table#main{
			border-spacing: 0px;
			display: inline-block;
		}
		#main tr{
			height: 75px;
		}
		#main *[t="red"]{
			color: red;
		}
		#main *[t="bred"]{
			background: pink;
		}
	</style>
	<form method="POST" action="/new">
		<input value="" name="url" placeholder="Загрузить новый тайтл">
		<input value="загрузить" name="submit" type="submit"><br/>
		<details>
			<summary>Пример</summary>
			<label>https://imhentai.xxx/gallery/653578/</label><br/>
			<label>https://manga-online.biz/rebirth_of_the_urban_immortal_cultivator/1/401/1/</label><br/>
		</details>
	</form>
	<form method="POST" action="/prepare" target="blank">
		<input value="" type="number" name="from" placeholder="С">
		<input value="" type="number" name="to" placeholder="По">
		<input value="подготовить архив" name="submit" type="submit">
	</form>
    <table id="main">
		<tbody>
		{{range .}}
			<tr t="{{if not .Loaded}}bred{{end}}">
				<td rowspan="2">
					{{if eq .Ext ""}}
					{{else}}
						<img src="/file/{{.ID}}/1.{{.Ext}}" style="max-width: 100px; max-height: 150px;">
					{{end}}
				</td>
				<td colspan="3" t="{{if not .Loaded}}red{{end}}">{{.Name}}</td>
			</tr>
			<tr t="{{if not .Loaded}}bred{{end}}">
				<td>#{{.ID}}</td>
				<td t="{{if not .ParsedPage}}red{{end}}">Страниц: {{.PageCount}}</td>
				<td t="{{if not .ParsedPage}}red{{end}}">Загружено: {{printf "%02.2f" .Avg}}%</td>
			</tr>
		{{end}}
		<tbody>
	</table>
  </body>
</html>
{{end}}
{{define "success"}}
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>HGRABER</title>
  </head>
  <body>
		<h1 style="color:green">Успешно: {{.}}</h1>
		<a href="/">главная</a>
  </body>
</html>
{{end}}
{{define "error"}}
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>HGRABER</title>
  </head>
  <body>
		<h1 style="color:red">Ошибка: {{.}}</h1>
		<a href="/">главная</a>
  </body>
</html>
{{end}}
`)
	if err != nil {
		log.Panicln(err)
	}
}
