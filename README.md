# protoc-gen-pehredaar

Generates code for rights from proto messages for rights implememntation .

protoc -I /usr/local/include -I ./  --go_out=plugins=grpc:./  ./src/rights/rights.proto

protoc -I /usr/local/include -I ./  --go_out=plugins=grpc:./  ./src/pehredaar/pehredaar.proto

protoc -I /usr/local/include -I ./ -I ./src --go_out=plugins=grpc:./  ./src/example/example.proto

packr && go build && protoc -I /usr/local/include -I  ./ -I ./src --plugin=protoc-gen-pehredaar=protoc-gen-pehredaar  --pehredaar_out=:./  ./src/example/example.proto && goimports -w ./src

# protoc-gen-pehredaar - Build a Rights Module to secure the API/GraphQL endpoints

protoc-gen-pehredaar is a protoc plugin which is used to generate pehredaar middleware for securing endpoints. This middleware uses a glob notation to build specific rights pattern for each module.

# Getting Started

Let's get started with a quick example.

For information regarding protobuf :
Protocol Buffer Basics: You can refer [here](https://developers.google.com/protocol-buffers/docs/gotutorial) 

Now we are good to go:

The rpc of the proto looks like :

```
    rpc SinglePrimitive (SinglePrimitiveReq) returns (.google.protobuf.Empty){
        option (pehredaar.paths) = {
            resource: "id"
        };
    }
    
    message SinglePrimitiveReq {
        string id = 1;
    }
    
```

Corresponding to this rpc the generated code will be as folllows:

```
func (s *RightsRightsSamplesServer) SinglePrimitive(ctx context.Context, rightsvar *SinglePrimitiveReq) (*empty.Empty, error) {

	ResourcePathOR := make([]string, 0)
	ResourcePathAND := make([]string, 0)

	ResourcePathOR = append(ResourcePathOR,

		fmt.Sprintf("/RightsSamples/%s**/.SinglePrimitive",

			rightsvar.GetId(),
		),
	)

	validations := map[string]bool{}

	res, err := s.rightsCli.IsValid(ctx, &rightspb.IsValidRequest{
		ResourcePathOr:       ResourcePathOR,
		ResourcePathAnd:      ResourcePathAND,
		UserId:               userinfo.FromContext(ctx).Id,
		ModuleName:           "RightsSamples",
		AttributeValidations: validations,
	})
	if err != nil {
		return nil, err
	}

	if !res.IsValid {
		return nil, status.Errorf(codes.PermissionDenied, res.Reason)
	}

	return s.RightsSamplesCli.SinglePrimitive(ctx, rightsvar)
}

```

As seen in the generated code : The pattern is formed using the input value of id and then the pattern is appended in an array of
resource paths.

The user id is taken from context , which is added in a previous layer of middlewares of authorization when the user login in the system.
And then context is passed in this middleware.

The module name defines the module for which the rights needs to be validated . 
This provides modularity for only checking in specific module rights assigned to the user and optimizes the functionality.

Case 2 : When objects are present in request:
```
    rpc NestedObject (NestedObjectReq) returns (.google.protobuf.Empty){
        option (pehredaar.paths) = {
            resource: "object.id"
        };
    }
    
    message NestedObjectReq {
        Primitive object = 1;
    }
    
```

Case 3 : When multi level object is present

```
    rpc NestedObjects (NestedObjectsReq) returns (.google.protobuf.Empty){
        option (pehredaar.paths) = {
            resource: "object.string.id"
        };
    }
    
    message NestedObjectsReq {
        Object object = 1;
    }
    
```
Case 4 : When one of is present with primitive types

```
    rpc OneOfPrimitive (OneOfPrimitiveReq) returns (.google.protobuf.Empty){
        option (pehredaar.paths) = {
            resource: "id1"
            resource: "id2"
        };
    }

    message OneOfPrimitiveReq {
        oneof data {
            string id1 = 1;
            string id2 = 2;
        }
    }
```

Case 5 : When one of has objects 
    
```

    rpc OneOfObject (OneOfObjectReq) returns (.google.protobuf.Empty){
        option (pehredaar.paths) = {
            resource: "object1.id"
            resource: "object2.ids"
        };
    }
    
    
    message OneOfObjectReq {
        oneof data {
            Primitive object1 = 1;
            Primitive object2 = 2;
        }
    }

```

Case 6: When one of has both primitive and object

```

    rpc OneOfPrimitiveAndObject (OneOfPrimitiveAndObjectReq) returns (.google.protobuf.Empty){
        option (pehredaar.paths) = {
            resource: "object1.id"
            resource: "object2.ids"
        };
    }
    
    message OneOfPrimitiveAndObjectReq {
        oneof data {
            Primitive object1 = 1;
            StringArray object2 = 2;
        }
    }
    
```

Case 7 : When nested objects has primitives

```
    rpc NestedObjectPrimitive (NestedObjectPrimitiveReq) returns (.google.protobuf.Empty) {
        option (pehredaar.paths) = {
            resource: "object.string_array.ids"
        };
    }

    message NestedObjectPrimitiveReq {
        ObjectWithArray object = 1;
    }

```

Case 8 : When nested objects has repeated values 

```
    rpc NestedObjectRepeated (NestedObjectRepeatedReq) returns (.google.protobuf.Empty) {
        option (pehredaar.paths) = {
            resource: "object.object_array.ids"
        };
    }

    message NestedObjectRepeatedReq {
        repeated ObjectArray object = 1;
    }
```

Additionally , 

You can also define attribute based rights that are dependent on the value of a boolean attribute in request.

```
    rpc AttributeBasedRight (AttributeBasedRightReq) returns (.google.protobuf.Empty) {
        option (pehredaar.paths) = {
            resource: "id"
        };
    }
    
    
    message AttributeBasedRightReq {
        string id = 1 ;
        bool skip_input_validation = 2 [(pehredaar.attribute).skip = true];
    }
    
```

When Only Attribute Based validation is needed

```
    rpc OnlyAttributeBasedRight (OnlyAttributeBasedRightReq) returns (.google.protobuf.Empty) {
        option (pehredaar.paths) = {
            only_attribute_based: true
        };
    }

    message OnlyAttributeBasedRightReq {
        bool attr_obj1 = 1 [(pehredaar.attribute).skip = true];
        NestedObjectOnlyAttributeBasedRightReq1 obj1 = 2;
    }
    
    message NestedObjectOnlyAttributeBasedRightReq1 {
        bool attr_obj2 = 1 [(pehredaar.attribute).skip = true];
        NestedObjectOnlyAttributeBasedRightReq2 obj2 = 2;
    }
    
    message NestedObjectOnlyAttributeBasedRightReq2 {
        bool attr_obj3 = 1 [(pehredaar.attribute).skip = true];
    }

```

Also , There is also the functionality of generating default static rights if there is no dynamic value from input in the pattern formation.
   
```

    rpc StaticRights (StaticRightsReq) returns (.google.protobuf.Empty) {
        option (pehredaar.paths) = {
            resource_static: "StaticValue"
        };
    }

    rpc DefaultStaticRights (StaticRightsReq) returns (.google.protobuf.Empty) {
        option (pehredaar.paths) = {
            resource_static_default: true
        };
    }
```

Dummy Messages used for reference:

```
    message Primitive {
        string id = 1;
    }
    
    message Object {
        Primitive string = 1;
    }
    
    message StringArray {
        repeated string ids = 1;
    }
    
    message ObjectWithArray {
        StringArray string_array = 1;
    }
    
    message ObjectArray {
        repeated StringArray object_array = 1;
    }
    
    
    message StaticRightsReq {
        string id = 1;
    }

```

There is also the functionality of defining module roles that generates a function for inserting patterns in database.

```
option (pehredaar.module_roles).module_role = {
    module_role_name: "RPC and Pattern Test case Editor"
    display_name: "Editor"
    rpc: "SinglePrimitive"
    pattern: "{parent}/.BatchCheckAvailability"
    pattern : "skip_input_validation.AttributeBasedRight"
    pattern : "obj_1.obj_2.attr_obj3.OnlyAttributeBasedRight"
};

```

protoc-gen-pehredaar uses the options provided in the rpc to identify which input message property is dynamic and needs to be considered while glob pattern matching from input.


# Installing

The installation of protoc-gen-pehredaar can be done directly by running go get.

```
go get go.saastack.io/pehredaar
```

# Usage

The project uses [packr](https://github.com/gobuffalo/packr)

For a proto file sample.proto, the corresponding code is generated in sample.pb.rights.go file.
The command can be used for generating the code file.

````
packr && \
go build \
-o protoc-gen-pehredaar && \
mv ./protoc-gen-pehredaar $GOPATH/bin && \
protoc -I ./src -I ./ \
--go_out=:. --pehredaar_out=:. \
./src/sample/sample.proto && && goimports -w .
````

protoc-gen-pehredaar generates the code for the middleware that calls a rights validator client for verifying the glob pattern corresponding to the module.

The Isvalid function uses the roles assigned to the user to build rights corresponding to user and verify them with respect to ResourcePathOr and ResourcePathAND.

The IsValid function of rights validator verifies that if the user has any of the rights mentioned in ResourcePathOr as well as all the tights mentioned in ResourcePathAND.

ResourcePathOr code can be understood from this:
```

for _,permission := range ResourcePathOr {

		for _, module := range userRights.Resources {

			for i := 0; i < len(module.Allowed); i++ {

				g = glob.MustCompile(module.Allowed[i])

				if g.Match(permission) {
					ORPatternAllowed = true
					isAllowed = true
				}
			}
		}

	}

```
ResourcePathAND code can be understood from this:

```
for _,permission := range ResourcePathAND {

			isAllowed = false

			for _, module := range userRights.Resources {

				for i := 0; i < len(module.Allowed); i++ {

					g = glob.MustCompile(module.Allowed[i])

					if g.Match(permission) {
						isAllowed = true
						break
					}
				}
			}

			//Checking for AND
			if isAllowed == false {
				break
			}
		}

```

# Available Options

The behaviour of protoc-gen-pehredaar can be modified using the following options:

# File Option : (pehredaar.module_roles).module_role

*  This option is used to generate a default patterns function for inserting in db 
   which are then stored against user and later used in verifying.
         
    Properties of options:                                       
    * module_role_name : Name of module
    * display_name : The disply name of module
    * pattern : Manually define patterns
    * service_name : Name of the Module/Service . NOTE: Only if multiple services in a single proto
    * skip_service_name : Boolean defining whether to skip adding of service name in front of patterns (only valid for pattern not rpc)
    * rpc :Name of rpc (to add the patterns of rpc mentioned) . NOTE : exact name of rpc required

# Method Option : option (pehredaar.paths)

*   This option is used to define the generation of rights middleware code that will form a pattern using the input value.

    Properties of options:   
    * resource : To define the input property to use while checking glob pattern
    * resource_and = To define the input property to use while checking glob pattern for all repeated values of input.
    * resource_static = To define a static pattern for the module.
    * resource_static_default : Boolean value to generate a default right.
    * only_attribute_based : Boolean value to generate only attribute based right. When input level rights are not needed. NOTE: This skips the validation if no attribute is sent true in request input.
    * allow_parent : Boolean value to represent the additional checking that will be done corresponding to parent in the id.

# Field Options : (pehredaar.attribute)

* skip : This option is used to mark the field for attribute based validation.

# Authors:

* Ayush Gupta 
