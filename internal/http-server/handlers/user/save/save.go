package save

import (
	resp "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/api/response"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/logger/sl"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Id int64 `json:"id,omitempty"`
}

type UserSaver interface {
	SaveUser() (*storage.UserDTO, error)
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		res, err := userSaver.SaveUser()
		if err != nil {
			log.Error("failed to save user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save user"))

			return
		}
		log.Info("user added", slog.Int64("id", res.ID))

		responseOK(w, r, res.ID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, userId int64) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       userId,
	})
}
