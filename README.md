# jsonpatch

A simple JSON Patch Utils Lib for Kubenetes

reference: [https://github.com/mattbaird/jsonpatch](https://github.com/mattbaird/jsonpatch)

## Using

```shell
go get -u github.com/chinaran/jsonpatch
```

## Example

controller client using `client.RawPatch(types.JSONPatchType, patchData)`

dynamic client, client set, see the following example:

```go
import (
	"context"

	"github.com/chinaran/jsonpatch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func JSONPatchResource(dyclient dynamic.NamespaceableResourceInterface) (*unstructured.Unstructured, error) {
	patchValues := []*jsonpatch.Operation{
		// patch annotations
		jsonpatch.NewAnnotationOperation(jsonpatch.OpReplace, "key-1", "value-1"),
		// patch labels
		jsonpatch.NewLabelOperation(jsonpatch.OpReplace, "key-2", "value-2"),
		// add patch
		jsonpatch.NewAddOperation("/spec/arrayx/0", "value-3"),
		// replace patch
		jsonpatch.NewReplaceOperation("/spec/arrayx", []string{"value-3", "value-4"}),
		/// remove patch
		jsonpatch.NewRemoveOperation("/spec/arrayx/0"),
	}

	patchBytes, err := jsonpatch.NewOperationData(patchValues...)
	if err != nil {
		return nil, err
	}
	return dyclient.Namespace("test-ns").Patch(context.TODO(), "test-name", types.JSONPatchType, patchBytes, metav1.PatchOptions{})
}
```