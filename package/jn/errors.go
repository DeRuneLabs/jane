package jn

import "fmt"

var Errors = map[string]string{
	`no_stdlib`:                                `standard library directory not found`,
	`file_not_useable`:                         `file is not useable for this operating system or architecture`,
	`file_not_jn`:                              `this is not jn source file: %s`,
	`invalid_token`:                            `undefined code content: %c`,
	`invalid_syntax`:                           `invalid syntax`,
	`no_entry_point`:                           `entry point (main) function is not defined`,
	`exist_id`:                                 `identifier is already exist: %s`,
	`extra_closed_parentheses`:                 `extra closed parentheses`,
	`extra_closed_braces`:                      `extra closed braces`,
	`extra_closed_brackets`:                    `extra closed brackets`,
	`wait_close_parentheses`:                   `parentheses waiting to close`,
	`wait_close_brace`:                         `brace waiting to close`,
	`wait_close_bracket`:                       `bracket are waiting to close`,
	`expected_parentheses_close`:               `was expected parentheses close`,
	`expected_brace_close`:                     `was expected brace close`,
	`expected_bracket_close`:                   `was expected bracket close`,
	`body_not_exist`:                           `body is not exist`,
	`operator_overflow`:                        `operator overflow`,
	`incompatible_datatype`:                    `%s and %s data-types are not compatible`,
	`operator_notfor_jntype`:                   `%s operator is not defined for %s type`,
	`operator_notfor_float`:                    `%s operator is not defined for float type(s)`,
	`operator_notfor_int`:                      `%s operator is not defined for integer type(s)`,
	`operator_notfor_uint`:                     `%s operator is not defined for unsigned integer type(s)`,
	`id_noexist`:                               `identifier is not exist: %s`,
	`not_function_call`:                        `value is not function`,
	`argument_overflow`:                        `argument overflow`,
	`func_have_return`:                         `%s function cannot have return type`,
	`func_have_parameters`:                     `%s function cannot have parameter(s)`,
	`func_have_attributes`:                     `%s function cannot have attribute(s)`,
	`require_return_value`:                     `return statements of non-void functions should have return value`,
	`void_function_return_value`:               `void functions is cannot returns any value`,
	`bitshift_must_unsigned`:                   `bit shifting value is must be unsigned`,
	`logical_not_bool`:                         `logical expression is have only boolean type values`,
	`assign_const`:                             `constants is can't assign`,
	`assign_nonlvalue`:                         `lvalue required assignment`,
	`assign_type_not_support_value`:            `type is not support assignment`,
	`invalid_type`:                             `invalid data-type`,
	`invalid_attribute`:                        `invalid attribute for type`,
	`invalid_numeric_range`:                    `arithmetic value overflow`,
	"invalid_operator":                         "invalid operator",
	"invalid_type_unary_operator":              "invalid data-type for unary %s operator",
	`invalid_escape_sequence`:                  `invalid escape sequence`,
	`invalid_type_source`:                      `invalid data-type source`,
	`invalid_preprocessor`:                     `invalid preprocessor directive`,
	`invalid_pragma_directive`:                 `invalid pragma directive`,
	`invalid_type_for_const`:                   `%s is invalid data-type for constant`,
	`invalid_value_for_key`:                    `"%s" is invalid value for the "%s" key`,
	`invalid_expr`:                             `invalid expression`,
	`invalid_header_ext`:                       `invalid header extension: %s`,
	`missing_autotype_value`:                   `auto-type declarations should have a initializer`,
	`missing_type`:                             `data-type missing`,
	`missing_expr`:                             `expression missing`,
	`missing_block_comment`:                    `missing block comment close`,
	`missing_rune_end`:                         `rune is not finished`,
	`missing_ret`:                              `missing return at end of function`,
	`missing_string_end`:                       `string is not finished`,
	`missing_const_value`:                      `constants must have value specification`,
	`missing_multi_return`:                     `missing return values for multi return`,
	`missing_multiassign_identifiers`:          `missing identifier(s) for multiple assignment`,
	`missing_use_path`:                         `missing path of use statement`,
	`missing_pragma_directive`:                 `missing pragma directive`,
	`missing_goto_label`:                       `missing label identifier for goto statement`,
	`missing_expr_for`:                         `missing expression for %s`,
	`missing_generics`:                         `missing generics`,
	`expr_not_const`:                           `expressions is not constant expression`,
	`nil_for_autotype`:                         `nil is cannot use with auto-type definitions`,
	`void_for_autotype`:                        `void data is cannot use for auto-type definitions`,
	`rune_empty`:                               `rune is cannot empty`,
	`rune_overflow`:                            `rune is should be single`,
	`not_supports_indexing`:                    `%s data type is not support indexing`,
	`not_supports_slicing`:                     `%s data type is not support slicing`,
	`undefined_attribute`:                      `undefined attribute`,
	`attribute_repeat`:                         `this attribute is already given`,
	`already_constant`:                         `this define is already constant`,
	`already_variadic`:                         `this define is already variadic`,
	`already_reference`:                        `this define is already reference`,
	`already_uses`:                             `this path is already uses`,
	`ignore_id`:                                `ignore operator cannot use as identifier`,
	`overflow_multiassign_identifiers`:         `overflow multi assignment identifers`,
	`overflow_return`:                          `overflow return expressions`,
	`break_at_outiter`:                         `break keyword is cannot used at out of iter block`,
	`continue_at_outiter`:                      `continue keyword is cannot used at out of iter block`,
	`iter_while_notbool_expr`:                  `while iterations must be have boolean expression`,
	`iter_foreach_nonenumerable_expr`:          `foreach iterations must be have enumerable expression`,
	`much_foreach_vars`:                        `foreach variables can be maximum two`,
	`if_notbool_expr`:                          `if conditions must be have boolean expression`,
	`else_have_expr`:                           `else's cannot have any expression`,
	`variadic_parameter_notlast`:               `variadic parameter can only be last parameter`,
	`variadic_with_nonvariadicable`:            `%s data-type is not variadicable`,
	`more_args_with_variadiced`:                `variadic argument can't use with more argument`,
	`type_notsupports_casting`:                 `%s data-type not supports casting`,
	`type_notsupports_casting_to`:              `%s data-type not supports casting to %s data-type`,
	`notallow_declares`:                        `declare not allowed`,
	`notallow_multiple_assign`:                 `multiple assignments not allowed`,
	`attribute_not_supports`:                   `attribute is not supports by define`,
	`generics_not_supports`:                    `generics is not supports by define`,
	`use_at_content`:                           `use declaration must be start of source code`,
	`use_not_found`:                            `used directory path not found/access: %s`,
	`use_has_errors`:                           `used package has errors`,
	`def_not_support_pub`:                      `define is not supports pub modifier`,
	`obj_not_support_sub_fields`:               `object is not supports sub fields: %s`,
	`obj_have_not_id`:                          `object is not have sub field in this identifier: %s`,
	`doc_couldnt_generated`:                    `%s: documentation could not generated because X source code has an errors`,
	`declared_but_not_used`:                    `%s declared but not used`,
	`expr_not_func_call`:                       `statement must have function call expression`,
	`label_exist`:                              `label is already exist in this identifier: %s`,
	`label_not_exist`:                          `not exist any label in this identifier: %s`,
	`goto_jumps_declarations`:                  `goto %s jumps over declaration(s)`,
	`function_not_has_parameter`:               `function is not has parameter in this identifier: %s`,
	`already_has_expr`:                         `%s already has expression`,
	`argument_must_target_to_parameter`:        `argument must target to parameter`,
	`namespace_not_exist`:                      `namespace is not exist in this identifier: %s`,
	`overflow_limits`:                          `overflow the limit of data-type`,
	`generics_overflow`:                        `overflow generics`,
	`has_generics`:                             `define has generics`,
	`not_has_generics`:                         `define not has generics`,
	`not_lvalue_for_reference_param`:           `require lvalue for reference parameter`,
	`variadic_reference_param`:                 `referencing cannot combined with variadic parameters`,
	`func_must_have_generics_if_has_attribute`: `function is must be have minimum one generic type if has @%s attribute`,
	`func_cant_have_params_if_has_attribute`:   `function is cannot have parameter(s) if has @%s attribute`,
	`divide_by_zero`:                           `divide by zero`,
	`trait_hasnt_id`:                           `%s trait is not have this identifier: %s`,
	`notimpl_trait_def`:                        `not implemented %s trait's %s define`,
	`dynamic_generic_annotation_failed`:        `dynamic generic type annotation failed`,
	`fallthrough_wrong_use`:                    `fallthrough keyword can only useable at end of the case scopes`,
	`fallthrough_into_final_case`:              `fallthrough cannot useable at final case`,
}

func GetError(key string, args ...any) string {
	return fmt.Sprintf(Errors[key], args...)
}
