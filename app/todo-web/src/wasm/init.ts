export let module: WebAssembly.Module;
export let instance: WebAssembly.Instance;
import mainWasm from '../main.wasm?url'

// TODO: how to bundle it? 
// TODO: better UI

// @ts-ignore
export const go = new Go();

await WebAssembly.instantiateStreaming(fetch(mainWasm), go.importObject).then( async (result) => {
    module = result.module;
    instance = result.instance;
    go.run(instance)
    console.log("WASM initialized")
});

// declare functions from the wasm module
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

export const crdtOverFs = {
    newSyncS3: window.newSyncS3,
}
