package pehredaar

import "go.uber.org/fx"

var Options = []fx.Option{}

func GetRightsModule() fx.Option {
	return fx.Options(Options...)
}
