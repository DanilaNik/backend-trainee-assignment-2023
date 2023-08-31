package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/config"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/storage"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Storage) (*Storage, error) {
	const op = "storage.postgresql.New"

	dataSource := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s  sslmode=%s",
		cfg.Addr, cfg.Port, cfg.User, cfg.DB, cfg.Password, cfg.Sslmode,
	)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser() (*storage.UserDTO, error) {
	const op = "storage.postgresql.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(id) VALUES(DEFAULT) RETURNING *")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user storage.UserDTO
	err = stmt.QueryRow().Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get last insert id %w", op, err)
	}
	return &user, nil
}

func (s *Storage) DeleteUser(userId int64) error {
	const op = "storage.postgresql.DeleteUser"

	stmt, err := s.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(userId)
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}
	return nil
}

func (s *Storage) SaveSegment(name string) (*storage.SegmentDTO, error) {
	const op = "storage.postgresql.SaveSegment"

	stmt, err := s.db.Prepare("INSERT INTO segments(name) VALUES($1) RETURNING *")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var segment storage.SegmentDTO
	err = stmt.QueryRow(name).Scan(&segment.ID, &segment.Name)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, storage.ErrSegmentExists
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &segment, nil
}

func (s *Storage) DeleteSegment(name string) error {
	const op = "storage.postgresql.DeleteSegment"

	stmt, err := s.db.Prepare("DELETE FROM segments WHERE name = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(name)

	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}

	return nil
}

func (s *Storage) AddUserToSegments(segmentsToSave []string, segmentsToDelete []string, userId int64) (*storage.UserInSegmentDTO, error) {
	const op = "storage.postgresql.AddUserToSegments"

	addSegments := make([]string, 0)
	notAddSegments := make([]string, 0)
	for _, segment := range segmentsToSave {
		err := s.AddUserSegment(segment, userId)
		if err != nil {
			notAddSegments = append(notAddSegments, segment)
		} else {
			addSegments = append(addSegments, segment)
		}
	}
	for _, segment := range segmentsToDelete {
		s.DeleteUserSegment(segment, userId)
	}
	var userInSegment storage.UserInSegmentDTO
	userInSegment.UserID = userId
	userInSegment.AddedSegments = addSegments
	userInSegment.NotAddedSegments = notAddSegments
	userInSegment.DeletedSegments = segmentsToDelete

	return &userInSegment, nil
}

func (s *Storage) AddUserSegment(name string, id int64) error {
	const op = "storage.postgresql.AddUserSegment"

	segmentId, err := s.GetSegmentId(name)
	if err != nil {
		return storage.ErrSegmentNotFound
	}
	err = s.GetUserId(id)
	if err != nil {
		return storage.ErrSegmentNotFound
	}
	stmt, err := s.db.Prepare("INSERT INTO user_segments(user_id, segment_id) VALUES($1,$2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(id, segmentId)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return storage.ErrUserAlreadyInSegment
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetSegmentId(name string) (int64, error) {
	const op = "storage.postgresql.GetSegmentId"

	stmt, err := s.db.Prepare("SELECT id FROM segments WHERE name = ?")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var segmentId int64
	err = stmt.QueryRow(name).Scan(&segmentId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return segmentId, nil
}

func (s *Storage) GetUserId(id int64) error {
	const op = "storage.postgresql.GetUserId"

	stmt, err := s.db.Prepare("SELECT id FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	var userId int64
	err = stmt.QueryRow(id).Scan(&userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserSegment(name string, id int64) {
	const op = "storage.postgresql.DeleteUserSegment"

	segmentId, err := s.GetSegmentId(name)
	if err != nil {
		return
	}
	err = s.GetUserId(id)
	if err != nil {
		return
	}

	stmt, err := s.db.Prepare("DELETE FROM user_segments WHERE user_id = $1 AND segment_id = $2")
	if err != nil {
		return
	}

	_, _ = stmt.Exec(id, segmentId)

	return
}

func (s *Storage) GetUserSegments(userId int64) (*storage.UserSegmentsDTO, error) {
	const op = "storage.postgresql.GetUserSegments"

	stmt, err := s.db.Prepare("SELECT segment_id, segments.name FROM user_segments JOIN segments ON user_segments.segment_id = segments.id WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var userSegments storage.UserSegmentsDTO
	userSegments.UserId = userId
	for rows.Next() {
		var segment storage.SegmentDTO
		err := rows.Scan(&segment.ID, &segment.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		userSegments.Segments = append(userSegments.Segments, segment)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &userSegments, nil
}
