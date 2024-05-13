package validation

type Input interface {
	Validate() map[string]string
}
