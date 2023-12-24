package hgraberagent

import (
	"app/internal/domain/agent"
	"net/http"
	"strconv"
)

func (c *Controller) pageUpload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bookID, err := strconv.Atoi(r.Header.Get(agent.HeaderBookID))
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		page, err := strconv.Atoi(r.Header.Get(agent.HeaderPageNumber))
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		ext := r.Header.Get(agent.HeaderPageExtension)
		url := r.Header.Get(agent.HeaderPageUrl)

		err = c.useCases.UploadPage(
			ctx,
			agent.PageInfoToUpload{
				BookID:     bookID,
				PageNumber: page,
				URL:        url,
				Ext:        ext,
			},
			r.Body,
		)
		if err != nil {
			c.webtool.WritePlain(ctx, w, http.StatusBadRequest, err.Error())

			return
		}

		c.webtool.WriteNoContent(ctx, w)
	})
}
