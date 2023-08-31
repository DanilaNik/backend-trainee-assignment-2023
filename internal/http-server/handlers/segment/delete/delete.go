package delete

import (
	"errors"
	resp "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/api/response"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type Request struct {
	Name string `json:"Name" validate:"required"`
}

type Response struct {
	resp.Response
	Name string `json:"name,omitempty"`
}

type SegmentDeleter interface {
	DeleteSegment(name string) error
}

func New(log *slog.Logger, segmentDeleter SegmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		reqName := req.Name

		err = segmentDeleter.DeleteSegment(reqName)
		if err != nil {
			log.Error("failed to delete segment", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete segment"))

			return
		}

		log.Info("segment deleted", slog.String("name", reqName))

		responseOK(w, r, reqName)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, segmentName string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Name:     segmentName,
	})
}
