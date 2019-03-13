// This code is heavily inspired by:
// https://github.com/twitchtv/twirp-ruby/blob/master/protoc-gen-twirp_ruby/main.go
// which is licensed under the Apache License, Version 2.0.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	// "unicode"

	"./twirp/protobuf"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/twitchtv/protogen/typemap"
)

func main() {
	genReq := readGenRequest(os.Stdin)
	g := &generator{version: Version, genReq: genReq}
	genResp := g.Generate()
	writeGenResponse(os.Stdout, genResp)
}

type generator struct {
	version string
	genReq  *plugin.CodeGeneratorRequest
	reg     *typemap.Registry
}

func (g *generator) Generate() *plugin.CodeGeneratorResponse {
	resp := new(plugin.CodeGeneratorResponse)

	for _, f := range g.protoFilesToGenerate() {
		twirpFileName := packageFileName(filePath(f)) + "PB.hs"

		haskellCode := g.generateHaskellCode(f)
		respFile := &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(twirpFileName),
			Content: proto.String(haskellCode),
		}
		resp.File = append(resp.File, respFile)
	}

	return resp
}

func (g *generator) generateHaskellCode(file *descriptor.FileDescriptorProto) string {
	b := new(bytes.Buffer)
	print(b, "-- Code generated by protoc-gen-haskell %s, DO NOT EDIT.", g.version)
	print(b, "{-# LANGUAGE DerivingVia, DeriveAnyClass, DuplicateRecordFields #-}")
	print(b, "{-# OPTIONS_GHC -Wno-unused-imports -Wno-missing-export-lists #-}")

	moduleName := toModuleName(file)
	print(b, "module %sPB where", moduleName)
	print(b, "")

	print(b, "import           Control.DeepSeq")
	print(b, "import           Control.Monad (msum)")
	print(b, "import qualified Data.Aeson as A")
	print(b, "import qualified Data.Aeson.Encoding as E")
	print(b, "import           Data.ByteString (ByteString)")
	print(b, "import           Data.Int")
	print(b, "import           Data.Text (Text)")
	print(b, "import qualified Data.Text as T")
	print(b, "import           Data.Vector (Vector)")
	print(b, "import           Data.Word")
	print(b, "import           GHC.Generics")
	print(b, "import           Proto3.Suite")
	print(b, "import           Proto3.Suite.JSONPB as JSONPB")
	print(b, "import           Proto3.Wire (at, oneof)")



	ex, _ := proto.GetExtension(file.Options, haskell.E_Imports)
	if ex != nil {
		asString := *ex.(*string)
		imports := strings.Split(asString, ";")
		print(b, "");
		for _, val := range imports {
			print(b, "import qualified %s", val)
		}
	}


	for _, message := range file.MessageType {
		generateMessage(b, message)
	}
	for _, enum := range file.EnumType {
		generateEnum(b, enum)
	}

	return b.String()
}

func generateMessage(b *bytes.Buffer, message *descriptor.DescriptorProto) {
	oneofs := []string{}
	for _, oneof := range message.OneofDecl {
		generateOneof(b, message, oneof)
		oneofs = append(oneofs, oneof.GetName())
	}

	n := message.GetName()
	print(b, "")
	print(b, "data %s = %s", n, n)
	first := true
	for _, field := range message.Field {
		n := toHaskellFieldName(field.GetName())
		t := toType(field, "", "")
		if field.OneofIndex == nil {
			sep := ","
			if first {
				sep = "{"
			}
			print(b, "  %s %s :: %s", sep, n, t)
			first = false
		}
	}
	for _, n := range oneofs {
		t := fmt.Sprintf("%s%s", strings.Title(message.GetName()), pascalCase(n))
		sep := ","
		if first {
			sep = "{"
		}
		print(b, "  %s %s :: Maybe %s", sep, camelCase(n), t)
		first = false
	}
	print(b, "  } deriving stock (Eq, Ord, Show, Generic)")
	print(b, "    deriving anyclass (Named, NFData)")

	// Generate a FromJSONPB Instance
	print(b, "")
	print(b, "instance FromJSONPB %s where", n)
	print(b, "  parseJSONPB = A.withObject \"%s\" $ \\obj -> %s", n, n)
	for _, f := range fieldsForMessageInstance(message, "<$>", "<*>") {
		print(b, "    %s obj .: \"%s\"", f.sep, f.fieldName)
	}

	// Generate a ToJSONPB Instance
	print(b, "")
	print(b, "instance ToJSONPB %s where", n)
	print(b, "  toJSONPB %s{..} = object", n)
	print(b, "    [")
	for _, f := range fieldsForMessageInstance(message, " ", ",") {
		print(b, "    %s \"%s\" .= %s", f.sep, f.fieldName, f.fieldName)
	}
	print(b, "    ]")
	print(b, "  toEncodingPB %s{..} = pairs", n)
	print(b, "    [")
	for _, f := range fieldsForMessageInstance(message, " ", ",") {
		print(b, "    %s \"%s\" .= %s", f.sep, f.fieldName, f.fieldName)
	}
	print(b, "    ]")

	printToFromJSONInstances(b, n)

	print(b, "")
	print(b, "instance Message %s where", n)

	// encodeMessage impl
	print(b, "  encodeMessage _ %s{..} = mconcat", n)
	print(b, "    [")
	first = true
	for _, field := range message.Field {
		if field.OneofIndex == nil {
			fieldName := toHaskellFieldName(field.GetName())
			num := field.GetNumber()
			sep := ","
			if first {
				sep = " "
			}
			label := field.GetLabel()
			if label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					print(b, "    %s encodeMessageField %d (NestedVec %s)", sep, num, fieldName)
				} else {
					print(b, "    %s encodeMessageField %d (PackedVec %s)", sep, num, fieldName)
				}
			} else {
				if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					print(b, "    %s encodeMessageField %d (Nested %s)", sep, num, fieldName)
				} else {
					print(b, "    %s encodeMessageField %d %s", sep, num, fieldName)
				}
			}
			first = false
		}
	}
	for i, oneof := range message.OneofDecl {
		sep := ","
		if first {
			sep = " "
		}
		oneofName := toHaskellFieldName(oneof.GetName())
		print(b, "    %s case %s of", sep, oneofName)
		print(b, "         Nothing -> mempty")
		for _, field := range message.Field {
			if field.OneofIndex != nil && field.GetOneofIndex() == int32(i) {
				fieldName := toHaskellFieldName(field.GetName())
				n := pascalCase(field.GetName())
				num := field.GetNumber()
				print(b, "         Just (%s %s) -> encodeMessageField %d %s", n, fieldName, num, fieldName)
			}
		}

		first = false
	}
	print(b, "    ]")

	// decodeMessage impl
	print(b, "  decodeMessage _ = %s", n)
	first = true
	for _, field := range message.Field {
		if field.OneofIndex == nil {
			sep := "<*>"
			if first {
				sep = "<$>"
			}
			label := field.GetLabel()
			num := field.GetNumber()
			if label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					print(b, "    %s (nestedvec <$> at decodeMessageField %d)", sep, num)
				} else {
					print(b, "    %s (packedvec <$> at decodeMessageField %d)", sep, num)
				}
			} else {
				print(b, "    %s at decodeMessageField %d", sep, num)

			}
			first = false
		}
	}
	for i := range message.OneofDecl {
		sep := "<*>"
		if first {
			sep = "<$>"
		}
		print(b, "    %s oneof", sep)
		print(b, "         Nothing")
		print(b, "         [")
		firstSep2 := true
		for _, field := range message.Field {
			if field.OneofIndex != nil && field.GetOneofIndex() == int32(i) {
				sep2 := ","
				if firstSep2 {
					sep2 = " "
				}
				n := pascalCase(field.GetName())
				num := field.GetNumber()
				print(b, "         %s (%d, Just . %s <$> decodeMessageField)", sep2, num, n)
				firstSep2 = false
			}
		}
		print(b, "         ]")
	}

	// dotProto impl
	print(b, "  dotProto = undefined")

	for _, nested := range message.NestedType {
		generateMessage(b, nested)
	}
	for _, enum := range message.EnumType {
		generateEnum(b, enum)
	}
}

func generateOneof(b *bytes.Buffer, message *descriptor.DescriptorProto, oneof *descriptor.OneofDescriptorProto) {
	oneofName := oneof.GetName()
	n := fmt.Sprintf("%s%s", strings.Title(message.GetName()), pascalCase(oneofName))
	print(b, "")
	print(b, "data %s", n)
	first := true
	for _, field := range message.Field {
		n := pascalCase(field.GetName())
		t := toType(field, "(", ")")
		if field.OneofIndex != nil {
			sep := "|"
			if first {
				sep = "="
			}
			print(b, "  %s %s %s", sep, n, t)
			first = false
		}
	}
	print(b, "  deriving stock (Eq, Ord, Show, Generic)")
	print(b, "  deriving anyclass (Message, Named, NFData)")

	// Generate a FromJSONPB Instance
	print(b, "")
	print(b, "instance FromJSONPB %s where", n)
	print(b, "  parseJSONPB = A.withObject \"%s\" $ \\obj -> msum", n)
	print(b, "    [")
	for _, f := range fieldsForOneOfInstance(message, " ", ",") {
		print(b, "    %s %s <$> parseField obj \"%s\"", f.sep, f.fieldName, f.rawFieldName)
	}
	print(b, "    ]")

	// Generate a ToJSONPB Instance
	print(b, "")
	print(b, "instance ToJSONPB %s where", n)
	for _, f := range fieldsForOneOfInstance(message, " ", ",") {
		print(b, "  toJSONPB (%s x) = object [ \"%s\" .= x ]", f.fieldName, f.rawFieldName)
	}

	for _, f := range fieldsForOneOfInstance(message, " ", ",") {
		print(b, "  toEncodingPB (%s x) = pairs [ \"%s\" .= x ]", f.fieldName, f.rawFieldName)
	}

	printToFromJSONInstances(b, n)
}

func generateEnum(b *bytes.Buffer, enum *descriptor.EnumDescriptorProto) {
	n := enum.GetName()
	print(b, "")
	print(b, "data %s", n)
	first := true
	def := ""
	for _, value := range enum.Value {
		v := pascalCase(value.GetName())
		sep := "|"
		if first {
			sep = "="
			def = v
		}
		print(b, "  %s %s", sep, v)
		first = false
	}
	print(b, "  deriving stock (Eq, Ord, Show, Enum, Bounded, Generic)")
	print(b, "  deriving anyclass (Named, MessageField, NFData)")
	print(b, "  deriving Primitive via PrimitiveEnum %s", n)
	if def != "" {
		print(b, "")
		print(b, "instance HasDefault %s where def = %s", n, def)
	}

	// Generate a FromJSONPB Instance
	print(b, "")
	print(b, "instance FromJSONPB %s where", n)
	for _, value := range enum.Value {
		enum := strings.ToUpper(value.GetName())
		v := pascalCase(value.GetName())
		print(b, "  parseJSONPB (JSONPB.String \"%s\") = pure %s", enum, v)
	}
	print(b, "  parseJSONPB x = typeMismatch \"%s\" x", n)

	// Generate a ToJSONPB Instance
	print(b, "")
	print(b, "instance ToJSONPB %s where", n)
	print(b, "  toJSONPB x _ = A.String . T.toUpper . T.pack $ show x")
	print(b, "  toEncodingPB x _ = E.text . T.toUpper . T.pack  $ show x")

	printToFromJSONInstances(b, n)
}

type aField struct {
	sep          string
	fieldName    string
	rawFieldName string
}

func fieldsForOneOfInstance(message *descriptor.DescriptorProto, firstSep string, restSep string) []aField {
	fields := []aField{}
	first := true
	for _, field := range message.Field {
		fieldName := pascalCase(field.GetName())
		if field.OneofIndex != nil {
			sep := restSep
			if first {
				sep = firstSep
			}
			fields = append(fields, aField{sep: sep, fieldName: fieldName, rawFieldName: field.GetName()})
			first = false
		}
	}
	return fields
}

func fieldsForMessageInstance(message *descriptor.DescriptorProto, firstSep string, restSep string) []aField {
	fields := []aField{}

	first := true
	for _, field := range message.Field {
		if field.OneofIndex == nil {
			fieldName := toHaskellFieldName(field.GetName())
			sep := restSep
			if first {
				sep = firstSep
			}
			fields = append(fields, aField{sep: sep, fieldName: fieldName})
			first = false
		}
	}
	for _, oneof := range message.OneofDecl {
		fieldName := toHaskellFieldName(oneof.GetName())
		sep := restSep
		if first {
			sep = firstSep
		}
		fields = append(fields, aField{sep: sep, fieldName: fieldName})
		first = false
	}
	return fields
}

func printToFromJSONInstances(b *bytes.Buffer, n string) {
	print(b, "")
	print(b, "instance FromJSON %s where", n)
	print(b, "  parseJSON = parseJSONPB")

	print(b, "")
	print(b, "instance ToJSON %s where", n)
	print(b, "  toJSON = toAesonValue")
	print(b, "  toEncoding = toAesonEncoding")
}

// Reference: https://github.com/golang/protobuf/blob/c823c79ea1570fb5ff454033735a8e68575d1d0f/protoc-gen-go/descriptor/descriptor.proto#L136
func toType(field *descriptor.FieldDescriptorProto, prefix string, suffix string) string {
	ex, _ := proto.GetExtension(field.Options, haskell.E_Type)
	if ex != nil {
		return *ex.(*string)
	}

	label := field.GetLabel()
	res := ""
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		res = "Int32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		res = "Int64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		res = "Int32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		res = "Int64"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		res = fmt.Sprintf("%sSigned Int32%s", prefix, suffix)
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		res = fmt.Sprintf("%sSigned Int64%s", prefix, suffix)
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		res = "Word32"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		res = "Word64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		res = fmt.Sprintf("%sFixed Word32%s", prefix, suffix)
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		res = fmt.Sprintf("%sFixed Word64%s", prefix, suffix)
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		res = "Text"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		res = "ByteString"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		res = "Bool"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		res = "Float"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		res = "Double"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		res = toHaskellType(field.GetTypeName())
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		res = toHaskellType(field.GetTypeName())
	default:
		Fail(fmt.Sprintf("no mapping for type %s", field.GetType()))
	}

	if label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			res = fmt.Sprintf("Vector %s", res)
		} else {
			res = fmt.Sprintf("Vector %s", res)
		}
	} else if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		res = fmt.Sprintf("%sMaybe %s%s", prefix, res, suffix)
	}

	return res
}

// .foo.Message => Message
// google.protobuf.Empty => Google.Protobuf.Empty
func toHaskellType(s string) string {
	if len(s) > 1 && s[0:1] == "." {
		parts := strings.Split(s, ".")
		return parts[len(parts)-1]
	}

	parts := []string{}
	for _, x := range strings.Split(s, ".") {
		parts = append(parts, strings.Title(x))
	}
	return strings.Join(parts, ".")
}

// handle some names that are hard to deal with in Haskell like `id`.
func toHaskellFriendlyName(s string) string {
	switch s {
	case "id", "type":
		return s + "_"
	default:
		return s
	}
}

// snake_case to camelCase.
func toHaskellFieldName(s string) string {
	parts := []string{}
	for i, x := range strings.Split(s, "_") {
		if i == 0 {
			parts = append(parts, strings.ToLower(x))
		} else {
			parts = append(parts, strings.Title(strings.ToLower(x)))
		}
	}
	return toHaskellFriendlyName(strings.Join(parts, ""))
}

// protoFilesToGenerate selects descriptor proto files that were explicitly listed on the command-line.
func (g *generator) protoFilesToGenerate() []*descriptor.FileDescriptorProto {
	files := []*descriptor.FileDescriptorProto{}
	for _, name := range g.genReq.FileToGenerate { // explicitly listed on the command-line
		for _, f := range g.genReq.ProtoFile { // all files and everything they import
			if f.GetName() == name { // match
				files = append(files, f)
				continue
			}
		}
	}
	return files
}

func print(buf *bytes.Buffer, tpl string, args ...interface{}) {
	buf.WriteString(fmt.Sprintf(tpl, args...))
	buf.WriteByte('\n')
}

func filePath(f *descriptor.FileDescriptorProto) string {
	return *f.Name
}

// capitalize, with exceptions for common abbreviations
func capitalize(s string) string {
	return strings.Title(strings.ToLower(s))
}

func camelCase(s string) string {
	parts := []string{}
	for i, x := range strings.Split(s, "_") {
		if i == 0 {
			parts = append(parts, strings.ToLower(x))
		} else {
			parts = append(parts, capitalize(x))
		}
	}
	return strings.Join(parts, "")
}

func pascalCase(s string) string {
	parts := []string{}
	for _, x := range strings.Split(s, "_") {
		parts = append(parts, capitalize(x))
	}
	return strings.Join(parts, "")
}

func packageFileName(path string) string {
	ext := filepath.Ext(path)
	return pascalCase(strings.TrimSuffix(path, ext))
}

func packageType(path string) string {
	ext := filepath.Ext(path)
	path = strings.TrimSuffix(filepath.Base(path), ext)
	return pascalCase(path)
}

func toModuleName(file *descriptor.FileDescriptorProto) string {
	pkgName := file.GetPackage()
	parts := []string{}
	for _, p := range strings.Split(pkgName, ".") {
		parts = append(parts, capitalize(p))
	}

	apiName := packageType(filePath(file))
	parts = append(parts, apiName)

	return strings.Join(parts, ".")
}

func Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("error:", s)
	os.Exit(1)
}

func readGenRequest(r io.Reader) *plugin.CodeGeneratorRequest {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		Fail(err.Error(), "reading input")
	}

	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		Fail(err.Error(), "parsing input proto")
	}

	if len(req.FileToGenerate) == 0 {
		Fail("no files to generate")
	}

	return req
}

func writeGenResponse(w io.Writer, resp *plugin.CodeGeneratorResponse) {
	data, err := proto.Marshal(resp)
	if err != nil {
		Fail(err.Error(), "marshaling response")
	}
	_, err = w.Write(data)
	if err != nil {
		Fail(err.Error(), "writing response")
	}
}
