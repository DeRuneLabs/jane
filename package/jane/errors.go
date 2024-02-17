package jane

var Errors = map[string]string{
	`file_not_jane`:  `This is not jane source file: `,
	`invalid_token`:  `Undefined code content`,
	`invalid_syntax`: `Invalid syntax`,
	`function_body`:  `Function body is not declared`,
}
