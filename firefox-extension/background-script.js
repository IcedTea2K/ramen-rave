import * as _ from "./wasm_exec.js";

console.log("Running background script")
// Fetch and run Go wasm
const go = new Go();

WebAssembly.instantiateStreaming(fetch("bin/main.wasm"), go.importObject).then((result) => {
   go.run(result.instance);
});

// Listen and talk to the content script
let contentCommPort;

browser.runtime.onConnect.addListener((port) => {
   contentCommPort = port;
   contentCommPort.postMessage({ message: "You stinky" })
})
