package screener

type Repository interface {
	GetScreens() ([]Screen, error)
	SetScreens([]Screen) error
}
