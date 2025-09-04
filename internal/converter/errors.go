package converter

type URLError struct {
	Message string
}

func (e *URLError) Error() string {
	return e.Message
}

type ConversionError struct {
	Message string
}

func (e *ConversionError) Error() string {
	return e.Message
}

type FileError struct {
	Message string
}

func (e *FileError) Error() string {
	return e.Message
}

type BrowserError struct {
	Message string
}

func (e *BrowserError) Error() string {
	return e.Message
}