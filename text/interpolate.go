package text

import (
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"regexp"
)

// placeholdersRegex contains rules to find placeholders inside string
var placeholdersRegex = regexp.MustCompile(":[0-9a-zA-Z\\-_]+")

// Interpolate replaces all placeholders within source string using arguments bucket
// and string formatter
func Interpolate(source string, bucket xray.Bucket, format func(xray.Arg) string) string {
	if len(source) <= 1 || bucket == nil || bucket.Size() == 0 || format == nil {
		return source
	}

	return placeholdersRegex.ReplaceAllStringFunc(
		source,
		func(x string) string {
			arg := bucket.Get(x[1:])
			if arg == nil {
				arg = args.Nil(x[1:])
			}

			return format(arg)
		},
	)
}

// PlainInterpolator is argument to string converter, that returns only argument values
func PlainInterpolator(a xray.Arg) string {
	if a == nil {
		return ""
	}
	return a.Value()
}

// PlainInterpolatorBracketed is argument to string converter, that returns argument values in brackets
func PlainInterpolatorBracketed(a xray.Arg) string {
	if a == nil {
		return "<!>"
	} else if _, ok := a.(args.Nil); ok {
		return "<!" + a.Name() + "!>"
	}

	return "[" + a.Value() + "]"
}

// InterpolatePlainText performs plaintext interpolation
func InterpolatePlainText(source string, bucket xray.Bucket, brackets bool) string {
	if brackets {
		return Interpolate(source, bucket, PlainInterpolatorBracketed)
	}

	return Interpolate(source, bucket, PlainInterpolator)
}