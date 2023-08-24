package customErrors

// Illegal Path error
type IllegalPathError struct {
	Message string
}
type FileNotFoundError struct {
	Message string
}
type InvalidFormatError struct {
	Message string
}
type EncodingError struct {
	Message string
}
type IOError struct {
	Message string
}

func (err *IllegalPathError) Error() string {
	return err.Message
}
func (err *FileNotFoundError) Error() string {
	return err.Message
}
func (err *InvalidFormatError) Error() string {
	return err.Message
}
func (err *EncodingError) Error() string {
	return err.Message
}
func (err *IOError) Error() string {
	return err.Message
}
