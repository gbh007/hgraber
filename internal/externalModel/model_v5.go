package externalModel

import (
	"app/internal/domain/hgraber"
	"app/pkg"
	"strconv"
	"time"
)

type V5Info struct {
	Version string `json:"version"`
	Meta    V5Meta `json:"meta"`
	Data    V5Book `json:"data"`
}

type V5Meta struct {
	Exported       time.Time `json:"exported"`
	ServiceVersion string    `json:"service_version,omitempty"`
	ServiceName    string    `json:"service_name,omitempty"`
}

type V5Book struct {
	Name             string        `json:"name"`
	OriginURL        string        `json:"origin_url,omitempty"`
	PageCount        int           `json:"page_count"`
	CreateAt         time.Time     `json:"create_at"`
	AttributesParsed bool          `json:"attributes_parsed"`
	Attributes       []V5Attribute `json:"attributes,omitempty"`
	Pages            []V5Page      `json:"pages,omitempty"`
	Labels           []V5Label     `json:"labels,omitempty"`
}

type V5Page struct {
	PageNumber int       `json:"page_number"`
	Ext        string    `json:"ext"`
	OriginURL  string    `json:"origin_url,omitempty"`
	CreateAt   time.Time `json:"create_at"`
	Downloaded bool      `json:"downloaded,omitempty"`
	LoadAt     time.Time `json:"load_at,omitempty"`
	Labels     []V5Label `json:"labels,omitempty"`
}

type V5Label struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	CreateAt time.Time `json:"create_at"`
}

type V5Attribute struct {
	Code   string   `json:"code"`
	Values []string `json:"values"`
}

func V5BookFromDomain(raw hgraber.Book) V5Book {
	bookLabels := make([]V5Label, 0, 2)

	bookLabels = append(bookLabels, V5Label{
		Name:     "hg4:id",
		Value:    strconv.Itoa(raw.ID),
		CreateAt: raw.Created,
	})

	if raw.Data.Rating > 0 {
		bookLabels = append(bookLabels, V5Label{
			Name:     "hg4:rating",
			Value:    strconv.Itoa(raw.Data.Rating),
			CreateAt: raw.Created,
		})
	}

	b := V5Book{
		Name:             raw.Data.Name,
		PageCount:        raw.PageCount(),
		CreateAt:         raw.Created,
		AttributesParsed: raw.AttributesParsed(),
		Labels:           bookLabels,
		OriginURL:        raw.URL,
		Attributes: pkg.MapToSlice(raw.Data.Attributes, func(code hgraber.Attribute, values []string) V5Attribute {
			return V5Attribute{
				Code:   string(code),
				Values: values,
			}
		}),
		Pages: pkg.Map(raw.Pages, func(p hgraber.Page) V5Page {
			labels := make([]V5Label, 0, 1)

			if p.Rating > 0 {
				labels = append(labels, V5Label{
					Name:     "hg4:rating",
					Value:    strconv.Itoa(p.Rating),
					CreateAt: raw.Created, // Данная информация отсутствует
				})
			}

			return V5Page{
				PageNumber: p.PageNumber,
				Ext:        "." + p.Ext,
				OriginURL:  p.URL,
				CreateAt:   raw.Created, // Данная информация отсутствует
				Downloaded: p.Success,
				LoadAt:     p.LoadedAt,
				Labels:     labels,
			}
		}),
	}

	return b
}

func V5Convert(raw hgraber.Book) V5Info {
	return V5Info{
		Version: "1.0.0",
		Meta: V5Meta{
			Exported:       time.Now().UTC(),
			ServiceName:    "hgraber",
			ServiceVersion: "v4.2.0", // FIXME: встраивать генерацией
		},
		Data: V5BookFromDomain(raw),
	}
}
