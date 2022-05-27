package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
)

type consoleColorModeValue int

type RequestLabelMappingFn func(c *gin.Context) string

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var consoleColorMode = autoColor

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the webserve returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the webserve cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("%s %3d %s| %13v | %15s |%s %-7s %s %#v %s",
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// disableConsoleColor disables color output in the consoleserve.
func disableConsoleColor() {
	consoleColorMode = disableColor
}

// forceConsoleColor force color output in the consoleserve.
func forceConsoleColor() {
	consoleColorMode = forceColor
}

// ErrorLogger returns a handler func for any error type.
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT returns a handler func for a given error type.
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errors := c.Errors.ByType(typ)
		if len(errors) > 0 {
			c.JSON(-1, errors)
		}
	}
}

// New instances a Logger middleware that will write the logs to gin.DefaultWriter. By default gin.DefaultWriter = os.Stdout.
func New(opts ...Option) gin.HandlerFunc {
	cfg := &config{
		consoleColor: true,
		endpointLabelMappingFn: func(c *gin.Context) string {
			return c.Request.URL.Path
		}}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.formatter == nil {
		cfg.formatter = defaultLogFormatter
	}
	if cfg.consoleColor {
		forceConsoleColor()
	} else {
		disableConsoleColor()
	}
	isTerm := true
	//gin.DefaultWriter = &writeLogger{pool: buffer.Pool{}}
	return func(c *gin.Context) {
		if cfg.logger == nil {
			return
		}
		// Start timer
		start := time.Now()
		method := c.Request.Method
		endpoint := cfg.endpointLabelMappingFn(c)
		isOk := cfg.checkLabel(fmt.Sprintf("%d", c.Writer.Status()), cfg.excludeRegexStatus) && cfg.checkLabel(endpoint, cfg.excludeRegexEndpoint) && cfg.checkLabel(method, cfg.excludeRegexMethod)
		if !isOk {
			return
		}

		// Process request
		c.Next()
		raw := c.Request.URL.RawQuery
		param := LogFormatterParams{
			Request: c.Request,
			isTerm:  isTerm,
			Keys:    c.Keys,
		}
		// Stop timer
		param.ClientIP = c.ClientIP()
		param.Method = method
		param.StatusCode = c.Writer.Status()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			endpoint = endpoint + "?" + raw
		}
		param.Path = endpoint
		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		cfg.logger.Info(cfg.formatter(param))

	}
}

// checkLabel returns the match result of labels.
// Return true if regex-pattern compiles failed.
func (c *config) checkLabel(label string, patterns []string) bool {
	if len(patterns) <= 0 {
		return true
	}
	for _, pattern := range patterns {
		if pattern == "" {
			return true
		}
		matched, err := regexp.MatchString(pattern, label)
		if err != nil {
			return true
		}
		if matched {
			return false
		}
	}
	return true
}
