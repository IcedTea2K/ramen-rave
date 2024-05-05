import * as _ from "./wasm_exec.js";

console.log("Running background script")
const go = new Go();

WebAssembly.instantiateStreaming(fetch("bin/main.wasm"), go.importObject).then((result) => {
   go.run(result.instance);
});
