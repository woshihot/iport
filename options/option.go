package options

func RegisterOptions(o Options) {

}

func GetOptions() Options {
	if nil != &IportOptions {
		return IportOptions
	} else {
		return defaultOptions
	}
}

type Options struct {
	SendCloud bool
}

var IportOptions Options
var defaultOptions Options

func init() {
	defaultOptions = Options{
		SendCloud: true,
	}
}
