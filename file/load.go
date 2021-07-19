package file

import (
	"app/config"
	"app/parser"
	"fmt"
	"log"
	"os"
	"strings"
)

func Load(id, page int, URL, ext string) error {
	// создаем папку с тайтлом
	err := os.MkdirAll(fmt.Sprintf("%s/%d", config.DefaultFilePath, id), 0777)
	if err != nil && !os.IsExist(err) {
		log.Println(err)
		return err
	}
	// скачиваем изображение
	data, err := parser.RequestBytes(URL)
	if err != nil {
		return err
	}
	// создаем файл и загружаем туда изображение
	f, err := os.Create(fmt.Sprintf("%s/%d/%d.%s", config.DefaultFilePath, id, page, ext))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return f.Close()
}

func escapeFileName(n string) string {
	const r = " "
	if len([]rune(n)) > 200 {
		n = string([]rune(n)[:200])
	}
	for _, e := range []string{`\`, `/`, `|`, `:`, `"`, `*`, `?`, `<`, `>`} {
		n = strings.ReplaceAll(n, e, r)
	}
	return n
}
