package cli

import "regexp"

type ArgumentType struct {
	Type          string
	Aliases       []string
	DefaultFormat string
	IsValidFormat func(string) bool
	IsAny         bool
}

var KnwonArguments []*ArgumentType = []*ArgumentType{
	{
		Type:          "any",
		Aliases:       []string{"any", "unknown"},
		DefaultFormat: "v",
		IsValidFormat: func(s string) bool { return s == "v" },
		IsAny:         true,
	},
	{
		Type:          "bool",
		Aliases:       []string{"bool", "boolean"},
		DefaultFormat: "t",
		IsValidFormat: func(s string) bool { return s == "t" || s == "v" },
	},
	{
		Type:          "string",
		Aliases:       []string{"string", "str"},
		DefaultFormat: "s",
		IsValidFormat: func(s string) bool { return s == "s" || s == "v" },
	},
	{
		Type:          "int",
		Aliases:       []string{"int", "i"},
		DefaultFormat: "d",
		IsValidFormat: func(s string) bool { return s == "d" || s == "v" },
	},
	{
		Type:          "float64",
		Aliases:       []string{"float64", "f64", "double"},
		DefaultFormat: "f",
		IsValidFormat: func(s string) bool { return s == "v" || regexp.MustCompile(`^(\.\d+)?[feEgG]$`).MatchString(s) },
	},
}

func AnyKind() *ArgumentType {
	return KnwonArguments[0]
}

func IsKnownArgumentType(kind *ArgumentType) bool {
	for _, v := range KnwonArguments {
		if v == kind {
			return true
		}
	}
	return false
}

func FindArgumentType(kind string) *ArgumentType {
	for _, arg := range KnwonArguments {
		for _, alaias := range arg.Aliases {
			if alaias == kind {
				return arg
			}
		}
	}
	return nil
}
