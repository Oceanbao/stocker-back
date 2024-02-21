package screener

type Repository interface {
	GetScreenAll() ([]Screen, error)
	SetScreenAll([]Screen) error
}
