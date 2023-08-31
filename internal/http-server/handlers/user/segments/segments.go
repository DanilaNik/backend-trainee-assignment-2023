package segments

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
	Id int64 `json:"id" validate:"required"`
}

type Response struct {
	resp.Response
	Segments storage.UserSegmentsDTO `json:"segments,omitempty"`
}

type UserSegmentsGetter interface {
	GetUserSegments(userId int64) (*storage.UserSegmentsDTO, error)
}

func New(log *slog.Logger, userSegmentsGetter UserSegmentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.segments.New"

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

		id := req.Id

		userSegments, err := userSegmentsGetter.GetUserSegments(id)
		if err != nil {
			log.Error("failed to get user segments", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get user segments"))

			return
		}

		log.Info("get user segments", slog.Int64("id", id))

		responseOK(w, r, userSegments)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, userSegments *storage.UserSegmentsDTO) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Segments: *userSegments,
	})
}
