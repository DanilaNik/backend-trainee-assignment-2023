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
	Id int64 `json:"id" validate:"required"`
}

type Response struct {
	resp.Response
	Id int64 `json:"id,omitempty"`
}

type UserDeleter interface {
	DeleteUser(userId int64) error
}

func New(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.delete.New"

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

		err = userDeleter.DeleteUser(id)
		if err != nil {
			log.Error("failed to delete user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete user"))

			return
		}

		log.Info("user deleted")

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, userId int64) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       userId,
	})
}
