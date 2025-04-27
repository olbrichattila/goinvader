package api

type APIClient interface {
	AddScore(name string, score int) error
	Top10() (UserScores, error)
}
