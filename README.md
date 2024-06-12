## ⚠️ Everything in this repository is an unfinished experiment and not even close to a production-ready solution

## What

`crdt-over-fs` is an experimental library that uses CRDT and a filesystem for collaboration and synchronisation.

The library is inspired by [Collabs](https://github.com/mweidner037/fileshare-recipe-editor) and basically borrows the idea from the repo with the only difference that it moves the file system to the application layer. Moving the file system to the application layer makes it possible to use it in the browser and actually build an application on top of any storage that implements the file system interface, such as S3, local file system (with Dropbox/GDrive mounted), or Dropbox/GDrive/Any Drive API.

## Why

Mainstream local-first solutions have already solved (to some extent) many problems such as authorisation, synchronisation between edge devices and a remote server and many others. All this makes these libraries a perfect choice if you're building an application with a centralised infrastructure maintained by a company/team. These libraries can also be a good choice for applications with self-hosted capabilities, as your self-hosted users can deploy the infrastructure themselves.

But what if you're building a small, non-profit application for yourself or a small community of non-technical users? Obviously, it's not viable to have a centralised infrastructure because of the maintenance burden. So there has to be a solution that uses the infrastructure available for non-tech users. And it seems like google/dropbox/proton drives and S3 free-tier plans are a perfect fit for that. 

## Demo

The demo shows how a CLI and a web page can communicate via the S3 implementation of [File System](#fs).

https://github.com/3timeslazy/crdt-over-fs/assets/23486601/d2200f03-39c2-4faa-b521-f5b794bd6bd0

## Project structure

> ⚠️ If there is already a library like that please let me know!

> ⚠️ The library is on the very early stage. Its API most likely will change, it has no tests and lots of bugs and unsolved problems.

### Go library

Basically, the library requires two interfaces to be implemented for it to work:

#### FS

The `FS` interface abstracts a file system.

It's designed to be as simple and as minimalist as possible to make it easier to implement this interface for as many different types of storage as possible.

```go
type FS interface {
	MakeDir(name string) error
	ReadDir(name string) ([]DirEntry, error)
	WriteFile(name string, data []byte) error
	ReadFile(name string) ([]byte, error)
}

type DirEntry interface {
	Name() string
	IsDir() bool
}
```

#### CRDT

The `CRDT` abstracts a CRDT algorithm, making this library CRDT-agnostic. Again, the interface is kept as simple and small as possible.  

The small size of the interface puts certain constraints on what the library can do. Unfortunately, this interface can't help developers to separate the CRDT and application logic layers like [TinyBase](https://tinybase.org) does. 

```go
type CRDT interface {
	EmptyState() State
	Merge(s1, s2 State) (State, []Change, error)
}

type State []byte

type Change struct {
	Hash string
}
```

### How to use in JS/TS

Obviously, it doesn't make sense to make a local-first library that cannot be used by JS/TS for building web application. This problem is solved by WASM. To see how to use the compiled WASM library, please check out that [example](./app/todo-web).

In short, the exposed WASM API so far is as follows

```ts
declare global {
    interface Window {
        newSyncS3(opts: SyncS3Opts): Sync
    }
}

interface SyncS3Opts {
    sync: {
        stateId: string,
        rootDir: string,
    },
    crdt: Crdt,
    s3: {
        keyId: string,
        keySecret: string,
        endpoint: string,
        region: string,
        bucket: string
    }
}

interface Crdt {
    emptyState(): Uint8Array
    merge(s1: Uint8Array, s2: Uint8Array): { state: Uint8Array }
}

interface Sync {
    loadOwnState(): Promise<Uint8Array>
    saveOwnState(localState: Uint8Array): Promise<void>
    sync(localState: Uint8Array): Promise<{ state: Uint8Array }>
}
```

## Further plans and ideas

### Delete unused files.

Right now, the library does not delete unused files.

### Add local persistance.

Currently the library does only sync state over remote file system and does not store the state locally. 

### Make and publish the WASM package to npm.

So that users can just `npm install` it.

### Assets/Media files managment.

Some applications may need to store media files such as photos and videos somewhere. And since we've already got a file system, it makes sense to try and reuse that concept to store media files. 

I think it potentially might looks something like:

```ts
const sync = newSync({
	crdt:     automerge,
	localfs:  localStorageFs,
	remotefs: dropboxFs,
	assetsfs: s3Fs,
})
```

or even:

```ts
const sync = newSync({
	crdt:     automerge,
	localfs:  localStorageFs,
	remotefs: dropboxFs,
	assetsfs: {
		"png,jpeg,mp4": googlePhotosFs,
		"default":      dropboxFs,
	},
})
```

### Human-readable representation.

I think this is an important point in supporting user data ownership. Just because a user stores the files does not necessarily mean that the user owns those files. To a non-technical user, there's no difference between a binary automerge state file and magic, so it's important that the human-readable representation of the state is stored somewhere, so that the user can (at least manually) export their data and import it into another application.
