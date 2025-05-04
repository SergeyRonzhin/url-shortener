package handlers

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func (h HTTPHandler) Ping(rw http.ResponseWriter, rq *http.Request) {

	db, err := sqlx.ConnectContext(context.Background(), "postgres", h.options.DatabaseDsn)

	defer func() {
		err = db.Close()

		if err != nil {
			h.logger.Error(err)
		}
	}()

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info("Connection to database successful!")
	rw.WriteHeader(http.StatusOK)
}
