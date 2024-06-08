export let module: WebAssembly.Module;
export let instance: WebAssembly.Instance;
import mainWasm from '../main.wasm?url'

// @ts-ignore
export const go = new Go();

WebAssembly.instantiateStreaming(fetch(mainWasm), go.importObject).then( async (result) => {
    module = result.module;
    instance = result.instance;
    go.run(instance)
    console.log("WASM initialized")
});