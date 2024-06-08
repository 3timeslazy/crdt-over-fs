//go:build wasm && js

package main

import (
	"fmt"
	"syscall/js"

	"github.com/3timeslazy/crdt-over-fs/sync"
	"github.com/3timeslazy/crdt-over-fs/sync/fs/s3"

	"github.com/aws/aws-sdk-go/aws"
	awscred "github.com/aws/aws-sdk-go/aws/credentials"
	awssess "github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
)

// NewSyncS3 creates new fSWrapper over S3.
//
// On js side:
//
//	const s3fs = newSyncS3({
//	    sync: {
//	        stateId: stateId,
//	        rootDir: rootDir,
//	    },
//	    crdt: {
//	        emptyState: () => {...},
//	        merge: () => {...}
//	    },
//	    s3: {
//	        keyId: keyId,
//	        keySecret: keySecret,
//	        endpoint: endpoint,
//	        region: region,
//	        bucket: bucket
//	    }
//	});
func NewSyncS3(this js.Value, inputs []js.Value) any {
	jsS3 := inputs[0].Get("s3")

	creds := awscred.NewStaticCredentials(
		jsS3.Get("keyId").String(),
		jsS3.Get("keySecret").String(),
		"",
	)
	s3conf := &aws.Config{
		Credentials:      creds,
		Endpoint:         aws.String(jsS3.Get("endpoint").String()),
		Region:           aws.String(jsS3.Get("region").String()),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := awssess.NewSession(s3conf)
	if err != nil {
		return err
	}
	s3fs := s3.NewFS(
		awss3.New(sess),
		jsS3.Get("bucket").String(),
	)

	jsCRDT := inputs[0].Get("crdt")
	jsSync := inputs[0].Get("sync")

	wrapper := sync.NewFSWrapper(
		s3fs,
		&JSCRDT{jsCRDT},
		jsSync.Get("stateId").String(),
		jsSync.Get("rootDir").String(),
	)

	return js.ValueOf(map[string]any{
		"loadOwnState": AsyncFn(func() (any, error) {
			return wrapper.LoadOwnState()
		}),
		// TODO: pass arguments to the async function
		// "saveOwnState": AsyncFn(func() (any, error) {
		// 	wrapper.SaveOwnState()
		// })
	})
}

func AsyncFn(gofn func() (any, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		asyncfn := js.FuncOf(func(this js.Value, args []js.Value) any {
			resolve := args[0]
			reject := args[1]

			go func() {
				v, err := gofn()
				if err != nil {
					jserr := js.Global().Get("Error")
					reject.Invoke(jserr.New(err.Error()))
					return
				}

				if b, ok := v.(sync.State); ok {
					arr := js.Global().Get("Uint8Array").New(len(b))
					fmt.Println(v, js.CopyBytesToJS(arr, b))
					v = arr
				}
				resolve.Invoke(v)
			}()

			return nil
		})

		pr := js.Global().Get("Promise")
		return pr.New(asyncfn)
	})
}

func main() {
	js.Global().Set("newSyncS3", js.FuncOf(NewSyncS3))
	<-make(chan bool)
}
