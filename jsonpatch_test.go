package jsonpatch

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRFC6901Encode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not need encode",
			args: args{
				s: "abcd_efg-hij",
			},
			want: "abcd_efg-hij",
		},
		{
			name: "need encode",
			args: args{
				s: "abc/def/hij~123",
			},
			want: "abc~1def~1hij~0123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RFC6901Encode(tt.args.s); got != tt.want {
				t.Errorf("RFC6901Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAddOperation(t *testing.T) {
	type args struct {
		path  string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Operation
	}{
		{
			// Add a new element to a positional array
			// kubectl patch sa default --type='json' -p='[{"op": "add", "path": "/secrets/1", "value": {"name": "whatever" } }]'
			name: "common",
			args: args{
				path:  "/secrets/1",
				value: map[string]string{"name": "whatever"},
			},
			want: &Operation{
				Op:    OpAdd,
				Path:  "/secrets/1",
				Value: map[string]string{"name": "whatever"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAddOperation(tt.args.path, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAddOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRemoveOperation(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *Operation
	}{
		{
			// Disable a deployment livenessProbe using a json patch with positional arrays
			// kubectl patch deployment valid-deployment  --type json   -p='[{"op": "remove", "path": "/spec/template/spec/containers/0/livenessProbe"}]'
			name: "common",
			args: args{
				path: "/spec/template/spec/containers/0/livenessProbe",
			},
			want: &Operation{
				Op:   OpRemove,
				Path: "/spec/template/spec/containers/0/livenessProbe",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRemoveOperation(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRemoveOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReplaceOperation(t *testing.T) {
	type args struct {
		path  string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Operation
	}{
		{
			// Update a container's image using a json patch with positional arrays
			// kubectl patch pod valid-pod --type='json' -p='[{"op": "replace", "path": "/spec/containers/0/image", "value":"new image"}]'
			name: "common",
			args: args{
				path:  "/spec/containers/0/image",
				value: "new image",
			},
			want: &Operation{
				Op:    OpReplace,
				Path:  "/spec/containers/0/image",
				Value: "new image",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReplaceOperation(tt.args.path, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReplaceOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAnnotationOperation(t *testing.T) {
	type args struct {
		op            Op
		annotationKey string
		value         interface{}
	}
	tests := []struct {
		name string
		args args
		want *Operation
	}{
		{
			name: "common",
			args: args{
				op:            OpReplace,
				annotationKey: "test/key",
				value:         "test-value",
			},
			want: &Operation{
				Op:    OpReplace,
				Path:  "/metadata/annotations/test~1key",
				Value: "test-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAnnotationOperation(tt.args.op, tt.args.annotationKey, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnnotationOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLabelOperation(t *testing.T) {
	type args struct {
		op       Op
		labelKey string
		value    interface{}
	}
	tests := []struct {
		name string
		args args
		want *Operation
	}{
		{
			name: "common",
			args: args{
				op:       OpReplace,
				labelKey: "test/key",
				value:    "test-value",
			},
			want: &Operation{
				Op:    OpReplace,
				Path:  "/metadata/labels/test~1key",
				Value: "test-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLabelOperation(tt.args.op, tt.args.labelKey, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLabelOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakePatchPath(t *testing.T) {
	type args struct {
		path    string
		newPart interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not need encode",
			args: args{
				path:    "/spec/labels/",
				newPart: "abcd_efg-hij",
			},
			want: "/spec/labels/abcd_efg-hij",
		},
		{
			name: "need encode 1",
			args: args{
				path:    "/spec/labels/",
				newPart: "abc/def/hij~123",
			},
			want: "/spec/labels/abc~1def~1hij~0123",
		},
		{
			name: "need encode 2",
			args: args{
				path:    "/spec/labels",
				newPart: "abc/def/hij~123",
			},
			want: "/spec/labels/abc~1def~1hij~0123",
		},
		{
			name: "empty path", // invalid case
			args: args{
				path:    "",
				newPart: "spec",
			},
			want: "/spec",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakePatchPath(tt.args.path, tt.args.newPart); got != tt.want {
				t.Errorf("MakePatchPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOperationData(t *testing.T) {
	type args struct {
		patchValues []*Operation
	}
	testPatch1 := NewRemoveOperation("/spec/template/spec/containers/0/livenessProbe")
	testPatch2 := NewReplaceOperation("/spec/containers/0/image", "new image")
	testBytes, _ := json.Marshal([]*Operation{testPatch1, testPatch2})
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				patchValues: []*Operation{testPatch1, testPatch2},
			},
			want:    testBytes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOperationData(tt.args.patchValues...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOperationData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOperationData() = %v, want %v", got, tt.want)
			}
		})
	}
}
