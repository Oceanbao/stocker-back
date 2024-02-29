package tracking

type Repository interface {
	GetTrackings() ([]Tracking, error)
	SetTracking(tracking Tracking) error
}
