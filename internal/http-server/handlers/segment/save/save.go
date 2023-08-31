package save

import (
	"errors"
	resp "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/api/response"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/logger/sl"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/storage"
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
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SegmentSaver interface {
	SaveSegment(name string) (*storage.SegmentDTO, error)
}

func New(log *slog.Logger, segmentSaver SegmentSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.save.New"

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

		segment, err := segmentSaver.SaveSegment(reqName)
		if err != nil {
			log.Error("failed to save segment", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save segment"))

			return
		}

		log.Info("segment saved", slog.String("name", reqName))

		responseOK(w, r, segment)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, segment *storage.SegmentDTO) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       segment.ID,
		Name:     segment.Name,
	})
}
