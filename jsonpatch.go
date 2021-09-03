// patch example: https://kubernetes.io/docs/reference/kubectl/cheatsheet/#patching-resources

package jsonpatch

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Op string

const (
	OpAdd     Op = "add"
	OpRemove  Op = "remove"
	OpReplace Op = "replace"
)

type Operation struct {
	Op    Op          `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

var rfc6901Encoder = strings.NewReplacer("~", "~0", "/", "~1")

// RFC6901Encode escapes as per RFC6901.
// Because the characters '~' (%x7E) and '/' (%x2F) have special
//   meanings in JSON Pointer, '~' needs to be encoded as '~0' and '/'
//   needs to be encoded as '~1' when these characters appear in a
//   reference token.
func RFC6901Encode(s string) string {
	return rfc6901Encoder.Replace(s)
}

func NewAddOperation(path string, value interface{}) *Operation {
	return &Operation{
		Op:    OpAdd,
		Path:  path,
		Value: value,
	}
}

func NewRemoveOperation(path string) *Operation {
	return &Operation{
		Op:   OpRemove,
		Path: path,
	}
}

func NewReplaceOperation(path string, value interface{}) *Operation {
	return &Operation{
		Op:    OpReplace,
		Path:  path,
		Value: value,
	}
}

func NewAnnotationOperation(op Op, annotationKey string, value interface{}) *Operation {
	path := MakePatchPath("/metadata/annotations/", annotationKey)
	return &Operation{
		Op:    op,
		Path:  path,
		Value: value,
	}
}

func NewLabelOperation(op Op, labelKey string, value interface{}) *Operation {
	path := MakePatchPath("/metadata/labels/", labelKey)
	return &Operation{
		Op:    op,
		Path:  path,
		Value: value,
	}
}

func MakePatchPath(path string, newPart interface{}) string {
	key := rfc6901Encoder.Replace(fmt.Sprintf("%v", newPart))
	if path == "" {
		return "/" + key
	}
	if strings.HasSuffix(path, "/") {
		return path + key
	}
	return path + "/" + key
}

func NewOperationData(patchValues ...*Operation) ([]byte, error) {
	return json.Marshal(patchValues)
}

func MustNewOperationData(patchValues ...*Operation) []byte {
	data, _ := json.Marshal(patchValues)
	return data
}
