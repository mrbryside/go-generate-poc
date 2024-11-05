package myhttp

const (
	StatusContinue           = "100"
	StatusSwitchingProtocols = "101"
	StatusProcessing         = "102"
	StatusEarlyHints         = "103"

	StatusOK                   = "200"
	StatusCreated              = "201"
	StatusAccepted             = "202"
	StatusNonAuthoritativeInfo = "203"
	StatusNoContent            = "204"
	StatusResetContent         = "205"
	StatusPartialContent       = "206"
	StatusMultiStatus          = "207"
	StatusAlreadyReported      = "208"
	StatusIMUsed               = "226"

	StatusMultipleChoices   = "300"
	StatusMovedPermanently  = "301"
	StatusFound             = "302"
	StatusSeeOther          = "303"
	StatusNotModified       = "304"
	StatusUseProxy          = "305"
	_                       = "306"
	StatusTemporaryRedirect = "307"
	StatusPermanentRedirect = "308"

	StatusBadRequest                   = "400"
	StatusUnauthorized                 = "401"
	StatusPaymentRequired              = "402"
	StatusForbidden                    = "403"
	StatusNotFound                     = "404"
	StatusMethodNotAllowed             = "405"
	StatusNotAcceptable                = "406"
	StatusProxyAuthRequired            = "407"
	StatusRequestTimeout               = "408"
	StatusConflict                     = "409"
	StatusGone                         = "410"
	StatusLengthRequired               = "411"
	StatusPreconditionFailed           = "412"
	StatusRequestEntityTooLarge        = "413"
	StatusRequestURITooLong            = "414"
	StatusUnsupportedMediaType         = "415"
	StatusRequestedRangeNotSatisfiable = "416"
	StatusExpectationFailed            = "417"
	StatusTeapot                       = "418"
	StatusMisdirectedRequest           = "421"
	StatusUnprocessableEntity          = "422"
	StatusLocked                       = "423"
	StatusFailedDependency             = "424"
	StatusTooEarly                     = "425"
	StatusUpgradeRequired              = "426"
	StatusPreconditionRequired         = "428"
	StatusTooManyRequests              = "429"
	StatusRequestHeaderFieldsTooLarge  = "431"
	StatusUnavailableForLegalReasons   = "451"

	StatusInternalServerError           = "500"
	StatusNotImplemented                = "501"
	StatusBadGateway                    = "502"
	StatusServiceUnavailable            = "503"
	StatusGatewayTimeout                = "504"
	StatusHTTPVersionNotSupported       = "505"
	StatusVariantAlsoNegotiates         = "506"
	StatusInsufficientStorage           = "507"
	StatusLoopDetected                  = "508"
	StatusNotExtended                   = "510"
	StatusNetworkAuthenticationRequired = "511"
)

var StatusCodeMap = map[string]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "SwitchingProtocols",
	StatusProcessing:         "Processing",
	StatusEarlyHints:         "EarlyHints",

	StatusOK:                   "OK",
	StatusCreated:              "Created",
	StatusAccepted:             "Accepted",
	StatusNonAuthoritativeInfo: "NonAuthoritativeInfo",
	StatusNoContent:            "NoContent",
	StatusResetContent:         "ResetContent",
	StatusPartialContent:       "PartialContent",
	StatusMultiStatus:          "MultiStatus",
	StatusAlreadyReported:      "AlreadyReported",
	StatusIMUsed:               "IMUsed",

	StatusMultipleChoices:   "MultipleChoices",
	StatusMovedPermanently:  "MovedPermanently",
	StatusFound:             "Found",
	StatusSeeOther:          "SeeOther",
	StatusNotModified:       "NotModified",
	StatusUseProxy:          "UseProxy",
	StatusTemporaryRedirect: "TemporaryRedirect",
	StatusPermanentRedirect: "PermanentRedirect",

	StatusBadRequest:                   "BadRequest",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "PaymentRequired",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "NotFound",
	StatusMethodNotAllowed:             "MethodNotAllowed",
	StatusNotAcceptable:                "NotAcceptable",
	StatusProxyAuthRequired:            "ProxyAuthRequired",
	StatusRequestTimeout:               "RequestTimeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "LengthRequired",
	StatusPreconditionFailed:           "PreconditionFailed",
	StatusRequestEntityTooLarge:        "RequestEntityTooLarge",
	StatusRequestURITooLong:            "RequestURITooLong",
	StatusUnsupportedMediaType:         "UnsupportedMediaType",
	StatusRequestedRangeNotSatisfiable: "RequestedRangeNotSatisfiable",
	StatusExpectationFailed:            "ExpectationFailed",
	StatusTeapot:                       "Teapot",
	StatusMisdirectedRequest:           "MisdirectedRequest",
	StatusUnprocessableEntity:          "UnprocessableEntity",
	StatusLocked:                       "Locked",
	StatusFailedDependency:             "FailedDependency",
	StatusTooEarly:                     "TooEarly",
	StatusUpgradeRequired:              "UpgradeRequired",
	StatusPreconditionRequired:         "PreconditionRequired",
	StatusTooManyRequests:              "TooManyRequests",
	StatusRequestHeaderFieldsTooLarge:  "RequestHeaderFieldsTooLarge",
	StatusUnavailableForLegalReasons:   "UnavailableForLegalReasons",

	StatusInternalServerError:           "InternalServerError",
	StatusNotImplemented:                "NotImplemented",
	StatusBadGateway:                    "BadGateway",
	StatusServiceUnavailable:            "ServiceUnavailable",
	StatusGatewayTimeout:                "GatewayTimeout",
	StatusHTTPVersionNotSupported:       "HTTPVersionNotSupported",
	StatusVariantAlsoNegotiates:         "VariantAlsoNegotiates",
	StatusInsufficientStorage:           "InsufficientStorage",
	StatusLoopDetected:                  "LoopDetected",
	StatusNotExtended:                   "NotExtended",
	StatusNetworkAuthenticationRequired: "NetworkAuthenticationRequired",
}

func statusCodeMapReversed() map[string]string {
	result := make(map[string]string)
	for key, value := range StatusCodeMap {
		result[value] = key
	}
	return result
}

func IsStatusCodeMap(statusCodeMap string) bool {
	for k := range statusCodeMapReversed() {
		if k == statusCodeMap {
			return true
		}
	}
	return false
}

func IsKeyInStatusMapping(key string) bool {
	for k := range StatusCodeMap {
		if k == key {
			return true
		}
	}
	return false
}
