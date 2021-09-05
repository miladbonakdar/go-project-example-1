package logger

func Trace(message string) {
	newUnitLogger().Trace(message)
}
func Debug(message string) {
	newUnitLogger().Debug(message)
}
func Print(message string) {
	newUnitLogger().Print(message)
}
func Info(message string) {
	newUnitLogger().Info(message)
}
func Warn(message string) {
	newUnitLogger().Warn(message)
}
func Warning(message string) {
	newUnitLogger().Warning(message)
}
func Error(message string) {
	newUnitLogger().Error(message)
}
func Fatal(message string) {
	newUnitLogger().Fatal(message)
}
func Panic(message string) {
	newUnitLogger().Panic(message)
}

func PanicException(exception error, message string) {
	newUnitLogger().PanicException(exception, message)
}
func WarnException(exception error, message string) {
	newUnitLogger().WarnException(exception, message)
}
func WarningException(exception error, message string) {
	newUnitLogger().WarningException(exception, message)
}
func ErrorException(exception error, message string) {
	newUnitLogger().ErrorException(exception, message)
}
func FatalException(exception error, message string) {
	newUnitLogger().FatalException(exception, message)
}

func WithException(exception error) *EntryLog {
	return newUnitLogger().WithException(exception)
}
func WithDevMessage(devMessage string) *EntryLog {
	return newUnitLogger().WithDevMessage(devMessage)
}
func WithData(data interface{}) *EntryLog {
	return newUnitLogger().WithData(data)
}
