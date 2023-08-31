package storage

import "errors"

var (
	ErrSegmentNotFound      = errors.New("Segment not found")
	ErrSegmentExists        = errors.New("Segment exists")
	ErrUserAlreadyInSegment = errors.New("User already in segment")
	ErrUserSegmentNotFound  = errors.New("User in segment not found")
	ErrUserNotFound         = errors.New("User not found")
)

type UserDTO struct {
	ID int64
}

type SegmentDTO struct {
	ID   int64
	Name string
}

type UserInSegmentDTO struct {
	UserID           int64
	AddedSegments    []string
	NotAddedSegments []string
	DeletedSegments  []string
}

type UserSegmentsDTO struct {
	UserId   int64
	Segments []SegmentDTO
}
