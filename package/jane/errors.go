package jane

var Errors = map[string]string{
	`file_not_jane`:            `This is not jane source file: `,
	`invalid_token`:            `Undefined code content`,
	`invalid_syntax`:           `Invalid syntax`,
	`function_body`:            `Function body is not declared`,
	`no_entry_point`:           `main function is not defined`,
	`exist_name`:               `Name is already exists`,
	`brace_not_closed`:         `Brace is opened but not closed`,
	`function_body_not_exist`:  `Function body is not declared`,
	`parameters_not_supported`: `Function is not supported parameters`,
	`not_support_expression`:   `expression is not support yet`,
	`missing_return`:           `missing return at end of function`,
}
