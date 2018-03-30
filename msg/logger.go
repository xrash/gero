package msg

type Logger interface {
	Error(format string, params ...interface{}) []error
}
