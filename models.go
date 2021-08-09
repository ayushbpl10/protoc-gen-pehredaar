package main

import "github.com/ayushbpl10/protoc-gen-pehredaar/pehredaar"

type rpcModel struct {
	Missing            bool
	PackageName        string
	RPCName            string
	Input              string
	Output             string
	Option             pehredaar.MyRights
	Resources          []Resource
	HasValidation      bool
	Validators         []ValidatorResource
	StaticRight        string
	DefaultStaticRight bool
	OnlyAttributeBased bool
	AllowParent        bool
	AllowStaff         bool
}

type ValidatorResource struct {
	Name, Description, Resource string
	ValidatorName               string
	OnlyValidator               bool
}

type Resource struct {
	IsRepeated                    bool
	GetStrings                    []map[string]bool
	ResourceStringWithCurlyBraces string
	ResourceStringWithFormatter   string
	Loops                         []ForLoop
	ResourceComment               string
	ResourceName                  string
	IsOneOf                       bool
	OneOfName                     string
	OneOfField                    string
	OneOfInput                    string
	And                           bool
	Constant                      string
	ConstantValue                 string
}

type moduleRoles struct {
	Name            string
	DisplayName     string
	Patterns        []string
	ServiceName     string
	SkipServiceName bool
	Priority        int32
	UniqueForModule bool
	Description     string
	External        bool
	AppName         string
}

type ForLoop struct {
	RangeKey string
	ValueKey string
	Level    int
}

type serviceModel struct {
	ServiceName   string
	PatternName   string
	PackageName   string
	Rpcs          []rpcModel
	HasValidation bool
}

type fileModel struct {
	PackageName string
	Imports     []string
	Services    []serviceModel
	ModuleRoles []moduleRoles
}
