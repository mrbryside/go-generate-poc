package handlergen

type Report struct {
	PathToGenerateSuccess       []string
	PathToGenerateError         []PathToGenerateError
	BasePathOfJsonSpec          string
	SwagGenerateReport          SwagGenerateReport
	HandlerTemplateSuccessRoute []HandlerTemplateData
	HandlerTemplateErrorRoute   []HandlerTemplatedDataError
	Error                       error
}

type SwagGenerateReport struct {
	isSuccess bool
	Error     error
}

type PathToGenerateError struct {
	Path  string
	Error error
}