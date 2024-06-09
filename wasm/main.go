//go:build wasm && js

package main

import (
	"syscall/js"

	"github.com/3timeslazy/crdt-over-fs/sync"
	"github.com/3timeslazy/crdt-over-fs/sync/fs/s3"

	"github.com/aws/aws-sdk-go/aws"
	awscred "github.com/aws/aws-sdk-go/aws/credentials"
	awssess "github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
)

// TODO: more args validation

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

	// TODO: add initRootDir
	return js.ValueOf(map[string]any{
		"loadOwnState": js.FuncOf(func(this js.Value, args []js.Value) any {
			return Promise(func() (any, error) {
				return wrapper.LoadOwnState()
			})
		}),
		"saveOwnState": js.FuncOf(func(this js.Value, args []js.Value) any {
			return Promise(func() (any, error) {
				state := BytesFromJS(args[0])
				return nil, wrapper.SaveOwnState(state)
			})
		}),
		"sync": js.FuncOf(func(this js.Value, args []js.Value) any {
			localState := BytesFromJS(args[0])

			return Promise(func() (any, error) {
				// TODO: handle changes
				newState, _, err := wrapper.Sync(localState)
				if err != nil {
					return nil, err
				}

				return map[string]any{
					"state": BytesToJS(newState),
				}, nil
			})
		}),
	})
}

func main() {
	js.Global().Set("newSyncS3", js.FuncOf(NewSyncS3))
	<-make(chan bool)
}
