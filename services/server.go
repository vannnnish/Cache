package services

const (
	APIVersion = "v1"
)

type Server interface {
	Run(address string) error
}
