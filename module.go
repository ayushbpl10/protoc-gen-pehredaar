package main

import (
	"encoding/json"
	"fmt"

	//"github.com/davecgh/go-spew/spew"

	"github.com/ayushbpl10/protoc-gen-pehredaar/pehredaar"

	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

var IMPORTS = []string{
	"github.com/golang/protobuf/ptypes/empty",
	"google/protobuf/fieldmask.proto",
}

type RightsGen struct {
	pgs.ModuleBase
	pgsgo.Context
}

func (*RightsGen) Name() string {
	return "pehredaar"
}

func (m *RightsGen) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (m *RightsGen) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {

	for _, f := range targets {

		//modulePath := m.Context.OutputPath(f).BaseName()
		name := m.Context.OutputPath(f).SetExt(".rights.go").String()
		fm := fileModel{PackageName: m.Context.PackageName(f).String()}

		for _, i := range IMPORTS {
			fm.Imports = append(fm.Imports, i)
		}

		for _, im := range f.Imports() {

			found := false
			for _, i := range IMPORTS {
				if i == im.Descriptor().Options.GetGoPackage() {
					found = true
				}
			}
			if !found {
				x := im.Descriptor().Options.GetGoPackage()
				if strings.Contains(x, ";") {
					x = strings.Split(x, ";")[0]
				}
				fm.Imports = append(fm.Imports, x)
			}
		}

		//fm.Imports = append(fm.Imports, modulePath+f.Descriptor().Options.GetGoPackage())

		serviceLength := 0
		serviceName := ""

		for _, srv := range f.Services() {

			if srv.Name().UpperCamelCase().String() == "ParentService" {
				continue
			}
			serviceLength++
			serviceName = srv.Name().String()

			service := serviceModel{}
			service.ServiceName = srv.Name().String()
			service.PackageName = m.Context.PackageName(f).String()

			// Code to handle service pattern option in pehredaar
			//if srv.Descriptor() == nil || srv.Descriptor().Options == nil {
			//	continue
			//}
			//
			//missingS := false
			//optS := srv.Descriptor().GetOptions()
			//optionS, err := proto.GetExtension(optS, pehredaar.E_Pattern)
			//if err != nil {
			//	if err == proto.ErrMissingExtension {
			//		missingS = true
			//	} else {
			//		panic(err)
			//	}
			//}
			//pattern := pehredaar.ServicePattern{}
			//
			//if !missingS {

			//data, err := json.Marshal(optionS)
			//if err != nil {
			//	panic(err)
			//}
			//err = json.Unmarshal(data, &pattern)
			//if err != nil {
			//	panic(err)
			//}
			//
			//service.PatternName = pattern.Pattern

			for _, rpc := range srv.Methods() {

				inputMsgMap := getMsgMap(rpc.Input(), "")

				right := pehredaar.MyRights{}
				rpcModel := rpcModel{RPCName: rpc.Name().UpperCamelCase().String(), Input: rpc.Input().Name().UpperCamelCase().String(), Output: rpc.Output().Name().UpperCamelCase().String(), Option: right, PackageName: m.Context.PackageName(f).String(), Missing: true}

				if rpc.Descriptor() == nil || rpc.Descriptor().Options == nil {
					service.Rpcs = append(service.Rpcs, rpcModel)
					continue
				}

				missing := false
				opt := rpc.Descriptor().GetOptions()
				option, err := proto.GetExtension(opt, pehredaar.E_Paths)
				if err != nil {
					if err == proto.ErrMissingExtension {
						missing = true
					} else {
						panic(err)
					}
				}

				//rpcModel := rpcModel{RPCName: rpc.Name().UpperCamelCase().String(), Input: rpc.Input().Name().UpperCamelCase().String(), Output: rpc.Output().Name().UpperCamelCase().String(), Option: right, PackageName: m.Context.PackageName(f).String(), Missing: missing}

				if !missing {
					byteData, err := json.Marshal(option)
					if err != nil {
						panic(err)
					}

					err = json.Unmarshal(byteData, &right)
					if err != nil {
						panic(err)
					}

					rpcModel.DefaultStaticRight = right.ResourceStaticDefault
					rpcModel.OnlyAttributeBased = right.OnlyAttributeBased
					rpcModel.AllowParent = right.AllowParent
					rpcModel.AllowStaff = right.AllowStaff

					if len(right.Resource) == 0 && len(right.ResourceAnd) == 0 && len(right.ResourceStatic) == 0 && right.ResourceStaticDefault == false && right.OnlyAttributeBased == false {
						rpcModel.Missing = true
					} else {

						rpcModel.Missing = false
						// Removing the validation field of resource and marking has validation in object in template
						for _, fieldRight := range right.Resource {
							if fieldRight == ".Validate" {
								rpcModel.HasValidation = true
								service.HasValidation = true
							}
						}

						// to distinguish between resourceOR and resourceAND
						resourceAndMap := map[string]bool{}

						//re := regexp.MustCompile("{([^{]*)}")
						//[right] : [with{} , without{}]
						fieldsInResource := make(map[string][][]string, 0)
						for _, fieldRight := range right.Resource { // ResourceOR
							elementCurlyWithoutCurlyBrackets := []string{fmt.Sprintf("{%s}", fieldRight), fieldRight}
							fieldsInResource[fieldRight] = append(fieldsInResource[fieldRight], elementCurlyWithoutCurlyBrackets)
						}
						for _, fieldRight := range right.ResourceAnd { // ResourceAnd
							elementCurlyWithoutCurlyBrackets := []string{fmt.Sprintf("{%s}", fieldRight), fieldRight}
							fieldsInResource[fieldRight] = append(fieldsInResource[fieldRight], elementCurlyWithoutCurlyBrackets)
							resourceAndMap[fieldRight] = true
						}

						for _, fieldRight := range right.ResourceStatic {
							fieldsInResource[fieldRight] = [][]string{}
						}

						if right.ResourceStaticDefault {
							fieldsInResource[fmt.Sprintf("/%s/.%s", srv.Name().UpperCamelCase().String(), rpcModel.RPCName)] = [][]string{}
						}

						//[right] : {[without {}] : seperatedFields}
						fieldVsfieldSeperatedRightMap := make(map[string]map[string][]string, 0)

						//Track keys to maintan map order
						fieldVsfieldSeperatedRightMapTrack := make(map[string][]string)

						ToBEreplacedByPlaceHolder := make([]string, 0)

						for fieldRight, ArrayOfCurlyBraces := range fieldsInResource {
							fieldVsfieldSeperated := make(map[string][]string, 0)
							for _, fr := range ArrayOfCurlyBraces {
								//%s place holders
								ToBEreplacedByPlaceHolder = append(ToBEreplacedByPlaceHolder, fr[0])
								//splitting the dot operator
								fieldVsfieldSeperated[fr[1]] = strings.Split(fr[1], ".")

								fieldVsfieldSeperatedRightMapTrack[fieldRight] = append(fieldVsfieldSeperatedRightMapTrack[fieldRight], fr[1])
								fieldVsfieldSeperatedRightMap[fieldRight] = fieldVsfieldSeperated
							}

							if len(ArrayOfCurlyBraces) == 0 {
								fieldVsfieldSeperatedRightMap[fieldRight] = fieldVsfieldSeperated
							}
						}

						//[right] : isrepeated
						RightRepeatedMap := make(map[string]bool)

						// validation of keys specified in right
						for rightVal, fieldVsfieldSeperatedElement := range fieldVsfieldSeperatedRightMap {

							for _, inOrder := range fieldVsfieldSeperatedRightMapTrack[rightVal] {

								dotSeperatedkeys, _ := fieldVsfieldSeperatedElement[inOrder]

								for _, r := range dotSeperatedkeys {
									//m.Log(r)
									//checking the key is a field in message
									found := false
									IsRepeated := false
									for _, msg := range f.AllMessages() {

										for _, field := range msg.Fields() {
											if field.Name().String() == r {
												found = true
												if _, ok := inputMsgMap[msg.Name().String()]; ok {
													IsRepeated = field.Type().IsRepeated()
												}
											}
										}
									}
									if !found {
										m.Log("key is not a primitive type : " + r)
									}

									if IsRepeated {
										RightRepeatedMap[rightVal] = true
									}
								}
							}

						}

						//m.Log(spew.Sdump(fieldVsfieldSeperatedRightMap))

						//if GlobalIsRepeated {
						for rightVal, fieldVsfieldSeperatedElement := range fieldVsfieldSeperatedRightMap {

							if len(fieldVsfieldSeperatedElement) == 0 && rpcModel.DefaultStaticRight {
								rpcModel.StaticRight = fmt.Sprintf("/%s/.%s", srv.Name().UpperCamelCase().String(), rpcModel.RPCName)
							} else if len(fieldVsfieldSeperatedElement) == 0 && !rpcModel.DefaultStaticRight {
								rpcModel.StaticRight = fmt.Sprintf("/%s/%s/.%s", srv.Name().UpperCamelCase().String(), rightVal, rpcModel.RPCName)
							}

							fieldVsGetString := make(map[string]string, 0)
							resource := Resource{And: resourceAndMap[rightVal]}
							//if repeated
							if RightRepeatedMap[rightVal] {
								for _, inOrder := range fieldVsfieldSeperatedRightMapTrack[rightVal] {
									fr := inOrder
									dotSeperatedkeys, _ := fieldVsfieldSeperatedElement[inOrder]

									//m.Log(fr)
									//[without curly braces] : Current For Loop Dotted Access
									ForLoopMapWithGet := make(map[string]string, 0)
									Level := 0
									ForLoopMapWithOutGet := make(map[string]string, 0)
									for _, r := range dotSeperatedkeys {

										//checking field is repeated
										IsRepeated := false
										IsNotPrimitive := false
										IsRepeatedPrimitive := false
										for _, msg := range f.AllMessages() {

											if _, ok := inputMsgMap[msg.Name().String()]; !ok {
												continue
											}

											for _, field := range msg.Fields() {
												if field.Name().String() == r {
													IsRepeated = field.Type().IsRepeated()
													IsNotPrimitive = field.Type().IsEmbed()
													if field.Type().Element() != nil {
														IsRepeatedPrimitive = field.Type().Element().IsEmbed()
													}

												}
											}
										}

										if IsRepeated {

											if _, ok := ForLoopMapWithGet[fr]; !ok {
												//m.Log(ForLoopMapWithGet[fr])
												ForLoopMapWithGet[fr] = "Get" + toCamelInitCase(r, true)
												ForLoopMapWithOutGet[fr] = toCamelInitCase(r, true)

												//for loops already contains the previous repeated value
												for _, forLoopResourceExists := range resource.Loops {
													if strings.Contains(forLoopResourceExists.RangeKey, ForLoopMapWithOutGet[fr]) {
														resource.Loops = resource.Loops[0 : len(resource.Loops)-1]
													}
												}

											} else {

												if strings.Contains(fr, ForLoopMapWithOutGet[fr]) {
													resource.Loops = resource.Loops[0 : len(resource.Loops)-1]
												}

												ForLoopMapWithGet[fr] = ForLoopMapWithOutGet[fr] + "." + "Get" + toCamelInitCase(r, true)

												//Level to add the ForLoopMapWithOutGet[fr] of the previous loop to next
												if Level >= 0 {
													ForLoopMapWithOutGet[fr] = toCamelInitCase(r, true)
												}
												Level++
											}

											forLoop := ForLoop{}
											//m.Log(ForLoopMapWithGet[fr])
											if _, ok := fieldVsGetString[fr]; ok {
												forLoop.RangeKey = fieldVsGetString[fr] + "." + "Get" + toCamelInitCase(r, true)
											} else {
												forLoop.RangeKey = ForLoopMapWithGet[fr]
											}
											forLoop.ValueKey = toCamelInitCase(r, true)

											resource.Loops = append(resource.Loops, forLoop)

											//starting the getString from inner for loop value
											fieldVsGetString[fr] = toCamelInitCase(r, true)

										} else {
											if _, ok := fieldVsGetString[fr]; ok {
												fieldVsGetString[fr] = fieldVsGetString[fr] + ".Get" + toCamelInitCase(r, true) + "()"
											} else {
												fieldVsGetString[fr] = "Get" + toCamelInitCase(r, true) + "()"
											}
										}

										resource.IsRepeated = RightRepeatedMap[rightVal]

										if !IsRepeated {
											if !IsNotPrimitive {
												found := false
												for _, strMap := range resource.GetStrings {
													if _, ok := strMap[fieldVsGetString[fr]]; ok {
														found = true
													}
												}
												if !found {
													mapGetString := make(map[string]bool, 0)
													mapGetString[fieldVsGetString[fr]] = false
													for _, forLoop := range resource.Loops {
														if strings.Contains(fieldVsGetString[fr], forLoop.ValueKey) {
															//m.Log(fr,fieldVsGetString[fr])
															mapGetString[fieldVsGetString[fr]] = true
														}
													}

													resource.GetStrings = append(resource.GetStrings, mapGetString)
												}
											}
										} else {
											if !IsRepeatedPrimitive {
												mapGetString := make(map[string]bool, 0)
												mapGetString[fieldVsGetString[fr]] = false
												for _, forLoop := range resource.Loops {
													if strings.Contains(fieldVsGetString[fr], forLoop.ValueKey) {
														//m.Log(fr,fieldVsGetString[fr])
														mapGetString[fieldVsGetString[fr]] = true
													}
												}

												resource.GetStrings = append(resource.GetStrings, mapGetString)
											}

										}

										resource.ResourceStringWithCurlyBraces = "/" + srv.Name().UpperCamelCase().String() + fmt.Sprintf("/{%s}**/.%s", rightVal, rpcModel.RPCName)
									}

								}
								//preparing formatted string
								if len(ToBEreplacedByPlaceHolder) > 0 {
									resource.ResourceStringWithFormatter = fmt.Sprintf("/%s/%s**/.%s", serviceName, "{"+rightVal+"}", rpcModel.RPCName)
								} else if rpcModel.DefaultStaticRight {
									resource.ResourceStringWithFormatter = fmt.Sprintf("/%s/.%s", serviceName, rpcModel.RPCName)
								} else {
									resource.ResourceStringWithFormatter = fmt.Sprintf("/%s/%s/.%s", serviceName, "{"+rightVal+"}", rpcModel.RPCName)
								}

								for _, p := range ToBEreplacedByPlaceHolder {
									resource.ResourceStringWithFormatter = strings.Replace(resource.ResourceStringWithFormatter, p, "%s", -1)
								}

								resource.ResourceName = rpc.Name().String()
								resource.ResourceComment = strings.Replace(rpc.SourceCodeInfo().LeadingComments(), "\n", "\\n", -1)

								if resource.ResourceStringWithCurlyBraces == "" && !rpcModel.DefaultStaticRight {
									resource.Constant = strings.ToUpper(srv.Name().SnakeCase().String() + "_" + rpc.Name().SnakeCase().String())
									resource.ConstantValue = "/" + srv.Name().String() + resource.ResourceStringWithFormatter
								} else if resource.ResourceStringWithCurlyBraces == "" && rpcModel.DefaultStaticRight {
									resource.Constant = strings.ToUpper(srv.Name().SnakeCase().String() + "_" + rpc.Name().SnakeCase().String())
									resource.ConstantValue = resource.ResourceStringWithFormatter

									// NOTE: Assigning just for module role patterns
									resource.ResourceStringWithCurlyBraces = resource.ResourceStringWithFormatter
								} else {
									resource.Constant = strings.ToUpper(srv.Name().SnakeCase().String() + "_" + rpc.Name().SnakeCase().String() + "_" + getResourceConstant(resource.ResourceStringWithCurlyBraces))
									resource.ConstantValue = resource.ResourceStringWithCurlyBraces
								}

								rpcModel.Resources = append(rpcModel.Resources, resource)

							} else {
								for fr, dotSeperatedkeys := range fieldVsfieldSeperatedElement {

									for _, r := range dotSeperatedkeys {

										//checking field is repeated
										IsNotPrimitive := false
										for _, msg := range f.AllMessages() {

											if _, ok := inputMsgMap[msg.Name().String()]; !ok {
												continue
											}

											for _, field := range msg.Fields() {
												if field.Name().String() == r {
													IsNotPrimitive = field.Type().IsEmbed()
													if field.InOneOf() {
														oneOf := field.OneOf()
														resource.IsOneOf = true
														resource.OneOfInput = inputMsgMap[msg.Name().String()] + oneOf.Name().UpperCamelCase().String()
														resource.OneOfName = msg.Name().UpperCamelCase().String()
														resource.OneOfField = field.Name().UpperCamelCase().String()
													}
												}
											}
										}

										if _, ok := fieldVsGetString[fr]; ok {
											fieldVsGetString[fr] = fieldVsGetString[fr] + ".Get" + toCamelInitCase(r, true) + "()"
										} else {
											fieldVsGetString[fr] = "Get" + toCamelInitCase(r, true) + "()"
										}

										resource.IsRepeated = RightRepeatedMap[rightVal]

										if !IsNotPrimitive {
											mapGetString := make(map[string]bool, 0)
											mapGetString[fieldVsGetString[fr]] = false
											resource.GetStrings = append(resource.GetStrings, mapGetString)
										}

										resource.ResourceStringWithCurlyBraces = "/" + srv.Name().UpperCamelCase().String() + fmt.Sprintf("/{%s}**/.%s", rightVal, rpcModel.RPCName)
									}

								}

								//preparing formatted string
								if len(ToBEreplacedByPlaceHolder) > 0 {
									resource.ResourceStringWithFormatter = fmt.Sprintf("/%s/%s**/.%s", serviceName, "{"+rightVal+"}", rpcModel.RPCName)
								} else if rpcModel.DefaultStaticRight {
									resource.ResourceStringWithFormatter = fmt.Sprintf("/%s/.%s", serviceName, rpcModel.RPCName)
								} else {
									resource.ResourceStringWithFormatter = fmt.Sprintf("/%s/%s/.%s", serviceName, "{"+rightVal+"}", rpcModel.RPCName)
								}
								for _, p := range ToBEreplacedByPlaceHolder {
									resource.ResourceStringWithFormatter = strings.Replace(resource.ResourceStringWithFormatter, p, "%s", -1)
								}

								resource.ResourceName = rpc.Name().String()
								resource.ResourceComment = strings.Replace(rpc.SourceCodeInfo().LeadingComments(), "\n", "\\n", -1)

								if resource.ResourceStringWithCurlyBraces == "" && !rpcModel.DefaultStaticRight {
									resource.Constant = strings.ToUpper(srv.Name().SnakeCase().String() + "_" + rpc.Name().SnakeCase().String())
									resource.ConstantValue = resource.ResourceStringWithFormatter
								} else if resource.ResourceStringWithCurlyBraces == "" && rpcModel.DefaultStaticRight {
									resource.Constant = strings.ToUpper(srv.Name().SnakeCase().String() + "_" + rpc.Name().SnakeCase().String())
									resource.ConstantValue = resource.ResourceStringWithFormatter

									// NOTE: Assigning just for module role patterns
									resource.ResourceStringWithCurlyBraces = resource.ResourceStringWithFormatter
								} else {
									resource.Constant = strings.ToUpper(srv.Name().SnakeCase().String() + "_" + rpc.Name().SnakeCase().String() + "_" + getResourceConstant(resource.ResourceStringWithCurlyBraces))
									resource.ConstantValue = resource.ResourceStringWithCurlyBraces
								}
								rpcModel.Resources = append(rpcModel.Resources, resource)
							}
						}

						for _, f := range rpc.Input().Fields() {

							if f.Type().ProtoType() == pgs.BoolT || f.Type().ProtoType() == pgs.MessageT {
								fName := ""

								m.MakeValidatorResource(f, rpc.Input(), &service, rpc, srv, &rpcModel, "", fName)
							}
						}
					}
				}

				service.Rpcs = append(service.Rpcs, rpcModel)

			}

			//} // end-if service-pattern ...

			fm.Services = append(fm.Services, service)
		}

		// Module Role options
		roles := getFileOptions(f)
		fm.ModuleRoles = []moduleRoles{}
		if roles != nil {

			for _, role := range roles.ModuleRole {

				patterns := []string{}

				var priority = role.Priority
				if priority == 0 {
					switch role.ModuleRoleName {
					case "Admin":
						priority = 3
					case "Editor":
						priority = 2
					case "Viewer":
						priority = 1
					}
				}

				if role.ServiceName == "" {

					if role.SkipServiceName {

						if len(fm.Services) > 1 {
							panic("Multiple Services in file, ServiceName file option required in Module Role, cannot skip it")
						}

						for _, p := range role.Pattern {
							if p != "" {
								patterns = append(patterns, fmt.Sprintf("%s", p))
							}
						}

						fm.ModuleRoles = append(fm.ModuleRoles, moduleRoles{
							Name:            role.ModuleRoleName,
							DisplayName:     role.DisplayName,
							Patterns:        patterns,
							UniqueForModule: !role.GroupingAllowed,
							Priority:        priority,
							ServiceName:     "",
							SkipServiceName: true,
							Description:     role.Description,
							AppName:         role.AppName,
							External:        role.External,
						})

					} else {

						if serviceLength != 1 {
							panic("Multiple Services in file, ServiceName file option required in Module Role")
						}

						if len(role.Rpc) > 0 {

							for _, rpcInProto := range role.Rpc {
								found := false
								for _, srv := range fm.Services {
									for _, rpc := range srv.Rpcs {
										if rpcInProto == rpc.RPCName {
											found = true
											for _, p := range rpc.Resources {
												if p.ResourceStringWithCurlyBraces != "" {
													patterns = append(patterns, p.ResourceStringWithCurlyBraces)
												}
											}
										}
									}
								}
								if !found {
									m.Log("RPC mentioned in module role not found in any service : " + rpcInProto + " , module role: " + role.ModuleRoleName)
								}
							}

						}

						for _, p := range role.Pattern {
							if p != "" {
								patterns = append(patterns, fmt.Sprintf("/%s/%s", serviceName, p))
							}
						}

						fm.ModuleRoles = append(fm.ModuleRoles, moduleRoles{
							Name:            serviceName + role.ModuleRoleName,
							DisplayName:     serviceName + " " + role.DisplayName,
							Patterns:        patterns,
							UniqueForModule: !role.GroupingAllowed,
							Priority:        priority,
							ServiceName:     serviceName,
							Description:     role.Description,
							AppName:         role.AppName,
							External:        role.External,
						})

					}

				} else {

					if len(role.Rpc) > 0 {

						for _, rpcInProto := range role.Rpc {
							found := false
							for _, srv := range fm.Services {
								for _, rpc := range srv.Rpcs {
									if rpcInProto == rpc.RPCName {
										found = true
										for _, p := range rpc.Resources {
											if p.ResourceStringWithCurlyBraces != "" {
												patterns = append(patterns, p.ResourceStringWithCurlyBraces)
											}
										}
									}
								}
							}
							if !found {
								m.Log("RPC mentioned in module role not found in any service : " + rpcInProto + " , module role: " + role.ModuleRoleName)
							}
						}

					}

					for _, p := range role.Pattern {
						if p != "" {
							patterns = append(patterns, fmt.Sprintf("/%s/%s", role.ServiceName, p))
						}
					}

					fm.ModuleRoles = append(fm.ModuleRoles, moduleRoles{
						Name:            role.ServiceName + role.ModuleRoleName,
						DisplayName:     role.ServiceName + " " + role.DisplayName,
						Patterns:        patterns,
						ServiceName:     role.ServiceName,
						UniqueForModule: !role.GroupingAllowed,
						Priority:        priority,
						Description:     role.Description,
						AppName:         role.AppName,
						External:        role.External,
					})
				}

			}
		}

		m.OverwriteGeneratorTemplateFile(
			name,
			T.Lookup("File"),
			&fm,
		)
	}

	return m.Artifacts()
}

func getMsgMap(msg pgs.Message, protoName string) map[string]string {

	msgMap := map[string]string{
		msg.Name().String(): protoName,
	}

	for _, field := range msg.Fields() {
		if field.Type().IsEmbed() {
			emb := field.Type().Embed()
			if emb.FullyQualifiedName() == msg.FullyQualifiedName() {
				continue
			}
			m := getMsgMap(emb, protoName+field.Name().UpperCamelCase().String()+".")

			for k, v := range m {
				msgMap[k] = v
			}
		}
		if field.Type().IsRepeated() && field.Type().Element().IsEmbed() {
			emb := field.Type().Element().Embed()
			if emb.FullyQualifiedName() == msg.FullyQualifiedName() {
				continue
			}
			m := getMsgMap(emb, protoName+field.Name().UpperCamelCase().String()+".")

			for k, v := range m {
				msgMap[k] = v
			}
		}
	}

	return msgMap
}

func getFileOptions(file pgs.File) *pehredaar.ModuleRoles {

	if file.Descriptor() == nil || file.Descriptor().GetOptions() == nil {
		return nil
	}

	opt := file.Descriptor().GetOptions()
	option, err := proto.GetExtension(opt, pehredaar.E_ModuleRoles)
	if err != nil {
		if err == proto.ErrMissingExtension {
			return nil
		} else {
			panic(err)
		}
	}

	byteData, err := json.Marshal(option)
	if err != nil {
		panic(err)
	}

	roles := pehredaar.ModuleRoles{}
	err = json.Unmarshal(byteData, &roles)
	if err != nil {
		panic(err)
	}

	return &roles
}

func getFieldOptionsAttribute(field pgs.Field) *pehredaar.Attribute {

	if field.Descriptor() == nil || field.Descriptor().GetOptions() == nil {
		return nil
	}

	opt := field.Descriptor().GetOptions()
	option, err := proto.GetExtension(opt, pehredaar.E_Attribute)
	if err != nil {
		if err == proto.ErrMissingExtension {
			return nil
		} else {
			panic(err)
		}
	}

	byteData, err := json.Marshal(option)
	if err != nil {
		panic(err)
	}

	attribute := pehredaar.Attribute{}
	err = json.Unmarshal(byteData, &attribute)
	if err != nil {
		panic(err)
	}

	return &attribute
}

func getResourceConstant(resource string) string {

	re := regexp.MustCompile("{([^{]*)}")

	arr := re.FindAllStringSubmatch(resource, -1)

	if len(arr) == 1 {
		return strings.Join(strings.Split(arr[0][1], "."), "_")
	}

	return ""
}

func (m *RightsGen) MakeValidatorResource(f pgs.Field, par pgs.Message, service *serviceModel, rpc pgs.Method, srv pgs.Service, rpcModel *rpcModel, getStr string, fName string) {

	if f.Type().ProtoType() == pgs.MessageT {
		if msg := f.Type().Embed(); msg != nil {
			for _, field := range msg.Fields() {

				fNameNew := fName
				if f.Type().ProtoType() == pgs.MessageT {
					if fNameNew != "" {
						fNameNew = fmt.Sprintf("%v.%v", fName, f.Name().LowerSnakeCase().String())
					} else {
						fNameNew = fmt.Sprintf("%v", f.Name().LowerSnakeCase().String())
					}

				}

				if par.FullyQualifiedName() != msg.FullyQualifiedName() {
					m.MakeValidatorResource(field, msg, service, rpc, srv, rpcModel, fmt.Sprintf("%v.Get%v()", getStr, f.Name().UpperCamelCase().String()), fNameNew)
				}
			}
		}

	}

	if f.Type().ProtoType() != pgs.BoolT {
		return
	}

	attribute := getFieldOptionsAttribute(f)
	if attribute != nil {

		service.HasValidation = true

		attrValue := fmt.Sprintf("%v.%v", fName, f.Name().LowerSnakeCase().String())
		if fName == "" {
			attrValue = fmt.Sprintf("%v", f.Name().LowerSnakeCase().String())
		}

		v := ValidatorResource{
			Name:          rpc.Input().Name().UpperCamelCase().String() + " " + f.Name().UpperCamelCase().String(),
			ValidatorName: fmt.Sprintf("rightsvar%v.Get%v()", getStr, f.Name().UpperCamelCase().String()),
			Description:   strings.Replace(f.SourceCodeInfo().LeadingComments(), "\n", "\\n", -1),
			Resource:      "/" + srv.Name().UpperCamelCase().String() + "/" + attrValue + "." + rpc.Name().UpperCamelCase().String(),
		}

		rpcModel.Validators = append(rpcModel.Validators, v)
	}
}
