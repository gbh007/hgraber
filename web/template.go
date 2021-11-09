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
		ul {
			padding: 0px;
			margin: 0px;
			list-style: none;
		}
	</style>
	<form method="POST" action="/new" style="display: inline">
		<input value="" name="url" placeholder="Загрузить новый тайтл">
		<details style="display: inline">
		<summary>Пример</summary>
			<ol style="text-align: left">
				<li>https://imhentai.xxx/gallery/653578/</li>
				<li>https://manga-online.biz/rebirth_of_the_urban_immortal_cultivator/1/401/1/</li>
			</ol>
		</details>
		<input value="загрузить" name="submit" type="submit">
	</form>
	<ul>
		<li>Всего <b>{{.Count}}</b> тайтлов</li>
		<li>Всего незагруженно <b>{{.UnloadCount}}</b> тайтлов</li>
		<li>Всего <b>{{.PageCount}}</b> страниц</li>
		<li>Всего незагруженно <b>{{.UnloadPageCount}}</b> страниц</li>
		<li><a href="/title/list?page=1">Список тайтлов</a></li>
	</ul>
  </body>
</html>
{{end}}
{{define "title-list"}}
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
		a.page {
			text-decoration: none;
			color: black;
			padding: 5px;
			border-radius: 5px;
			display: inline-block;
			border: 1px dashed black;
		}
		a.page[selected="true"] {
			background: lime;
		}
		a#title {
			text-decoration: none;
			color: black;
		}
		#title {
			display: inline-grid;
			grid-template-areas:
				"img name name name name"
				"img id pgc pgp dt"
				"img tag tag tag tag";
			grid-template-rows: none;
			grid-template-columns: 130px 1fr 1fr 1fr 1fr;
			border-spacing: 0px;
			max-width: 500px;
		}
		#title *[t="red"]{
			color: red;
		}
		#title *[t="bred"]{
			background: pink;
		}
		span.tag {
			border-radius: 3px;
			padding: 3px;
			margin: 2px;
			background: lightgrey;
			display: inline-block;
		}
	</style>
	<div>
		<form method="POST" action="/new" style="display: inline">
			<input value="" name="url" placeholder="Загрузить новый тайтл">
			<details style="display: inline">
			<summary>Пример</summary>
				<ol style="text-align: left">
					<li>https://imhentai.xxx/gallery/653578/</li>
					<li>https://manga-online.biz/rebirth_of_the_urban_immortal_cultivator/1/401/1/</li>
				</ol>
			</details>
			<input value="загрузить" name="submit" type="submit">
		</form>
		<form method="POST" action="/prepare" target="blank" style="display: inline">
			<input value="" type="number" name="from" placeholder="С">
			<input value="" type="number" name="to" placeholder="По">
			<input value="подготовить архив" name="submit" type="submit">
		</form>
	</div>
	<div style="padding: 10px">
		{{with $info := .}}
			{{range $info.Pages}}
				<a class="page" href="/title/list?page={{.}}" selected="{{if eq . $info.Page}}true{{end}}">{{.}}</a>
			{{end}}
		{{end}}
		<b>Всего {{.Count}} тайтлов</b>
	</div>
	{{range $ind, $e := .Titles}}
		{{template "title-short" $e}}
	{{end}}
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
{{define "title-short"}}
	<a href="/title/page?title={{.ID}}&page=1" id="title" t="{{if not .Loaded}}bred{{end}}">
		{{if eq .Ext ""}}
			<span style="grid-area: img;"></span>
		{{else}}
			<img src="/file/{{.ID}}/1.{{.Ext}}" style="max-width: 100%; max-height: 100%; grid-area: img;">
		{{end}}
		<span style="grid-area: name;" t="{{if not .Loaded}}red{{end}}">{{.Name}}</span>
		<span style="grid-area: id;">#{{.ID}}</span>
		<span style="grid-area: pgc;" t="{{if not .ParsedPage}}red{{end}}">Страниц: {{.PageCount}}</span>
		<span style="grid-area: pgp;" t="{{if ne .Avg 100.0}}red{{end}}">Загружено: {{printf "%02.2f" .Avg}}%</span>
		<span style="grid-area: dt;">{{.Created.Format "2006/01/02 15:04:05"}}</span>
		<span style="grid-area: tag;">
		{{range .Tags}}
			<span class="tag">{{.}}</span>
		{{end}}
		</span>
	</a>
{{end}}
{{define "title-page"}}
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>HGRABER</title>
  </head>
  <body>
    <script>
      document.addEventListener("keydown", function (event) {
        if (event.keyCode === 37) window.location.href="{{.Prev}}"
        if (event.keyCode === 39) window.location.href="{{.Next}}"
      });
    </script>
	<style>
		body {
		    text-align: center;
		}
		div.view {
			height: 90vh;
		}
		a.page {
			text-decoration: none;
			color: black;
		}
		h1.page {
			display: inline-block;
			writing-mode: vertical-lr;
			text-orientation: upright;
			text-decoration: none;
			color: black;
			height: 100%;
		    text-align: center;
			border: 2px dotted black;
			border-radius: 10px;
		}
	</style>
	<div>
		<a href="/">на главную</a>
		<details style="display: inline">
			<summary>перезагрузить изображение</summary>
			<form method="POST" action="/reload/page">
				<input value="{{.Page.TitleID}}" name="id" type="hidden">
				<input value="{{.Page.PageNumber}}" name="page" type="hidden">
				<input value="{{.Page.URL}}" name="url" placeholder="адрес">
				<input value="{{.Page.Ext}}" name="ext" placeholder="расширение">
				<input value="начать" name="submit" type="submit"><br/>
			</form>
		</details>
		Страница {{.Page.PageNumber}} из {{.Title.PageCount}}
	</div>
	<div class="view">
		<a class="page" href="{{.Prev}}"><h1 class="page">Назад</h1></a>
		<img src="{{.File}}" style="max-width: 100%; max-height: 100%;">
		<a class="page" href="{{.Next}}"><h1 class="page">Вперед</h1></a>
	<div>
  </body>
</html>
{{end}}
{{define "debug"}}
{{printf "%+v" .}}
{{end}}
`)
	if err != nil {
		log.Panicln(err)
	}
}
