package main

const (
	PackageRegexp = `^package\s[a-zA-Z_][0-9a-zA-Z_]*$`

	LineCommentRegexp         = `\/\/`
	MutLineCommentStartRegexp = `\/\*`
	MutLineCommentEndRegexp   = `\*\/`

	InitFunctionRegexp = `^func\sinit\(\)\s\{$`

	HessianImportRegexp = `"github.com/apache/dubbo-go-hessian2"`

	HessianPOJORegexp     = `\*[0-9a-zA-Z_]+\)\sJavaClassName\(\)\sstring\s\{$`
	HessianPOJONameRegexp = `\*[0-9a-zA-Z_]+\)`
)

const (
	newLine byte = '\n'
	funcEnd byte = '}'

	targetFileSuffix = ".go"
)
