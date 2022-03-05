package file

import (
	"app/db"
	"app/parser"
	"app/system"
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func Load(ctx system.Context, id, page int, URL, ext string) error {
	// создаем папку с тайтлом
	err := os.MkdirAll(fmt.Sprintf("%s/%d", system.GetFileStoragePath(ctx), id), 0666)
	if err != nil && !os.IsExist(err) {
		system.Error(ctx, err)
		return err
	}
	// скачиваем изображение
	data, err := parser.RequestBytes(ctx, URL)
	if err != nil {
		return err
	}
	// создаем файл и загружаем туда изображение
	f, err := os.Create(fmt.Sprintf("%s/%d/%d.%s", system.GetFileStoragePath(ctx), id, page, ext))
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	return f.Close()
}

// LoadToZip сохраняет тайтлы на диск zip архивом
func LoadToZip(ctx system.Context, id int) error {

	titleInfo, err := db.SelectTitleByID(ctx, id)
	if err != nil {
		return err
	}
	pages := db.SelectPagesByTitleID(ctx, id)

	buff := &bytes.Buffer{}
	zw := zip.NewWriter(buff)

	for _, p := range pages {
		f, err := os.Open(fmt.Sprintf("%s/%d/%d.%s", system.GetFileStoragePath(ctx), id, p.PageNumber, p.Ext))
		defer system.IfErrFunc(ctx, f.Close)
		if err != nil {
			system.Error(ctx, err)
			return err
		}

		tmpBuff := &bytes.Buffer{}

		_, err = tmpBuff.ReadFrom(f)
		if err != nil {
			system.Error(ctx, err)
			return err
		}

		w, err := zw.Create(fmt.Sprintf("%d.%s", p.PageNumber, p.Ext))
		if err != nil {
			system.Error(ctx, err)
			return err
		}

		_, err = w.Write(tmpBuff.Bytes())
		if err != nil {
			system.Error(ctx, err)
			return err
		}
	}

	w, err := zw.Create("info.txt")
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	_, err = fmt.Fprintf(
		w,
		"URL:%s\nNAME:%s\nPAGE-COUNT:%d\nINNER-ID:%d",
		titleInfo.URL,
		titleInfo.Name,
		titleInfo.PageCount,
		titleInfo.ID,
	)
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	err = zw.Close()
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	f, err := os.Create(fmt.Sprintf(
		"%s/%d)_%s.zip",
		system.GetFileStoragePath(ctx),
		id,
		escapeFileName(titleInfo.Name),
	))

	defer system.IfErrFunc(ctx, f.Close)

	if err != nil {
		system.Error(ctx, err)
		return err
	}
	_, err = buff.WriteTo(f)
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	return nil
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
