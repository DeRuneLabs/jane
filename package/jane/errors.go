package jane

var Errors = map[string]string{
	`file_not_jane`:            `this is not jane source file: `,
	`invalid_token`:            `undenfied code content`,
	`invalid_syntax`:           `invalid syntax`,
	`exist_name`:               `name is already exist`,
	`brace_not_closed`:         `brace is opened but not closed`,
	`function_body_not_exits`:  `function body is not declare`,
	`parameters_not_supported`: `function is not support parameters`,
	`not_support_expression`:   `expression is not supports yet`,
	`missing_return`:           `missing return at end of function`,
}
