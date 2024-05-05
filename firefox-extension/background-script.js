import * as _ from "./wasm_exec.js";

// Fetch and run Go wasm
const go = new Go();

WebAssembly.instantiateStreaming(fetch("bin/main.wasm"), go.importObject).then((result) => {
   go.run(result.instance);
});

// Listen and talk to the content script
let contentCommPort;

browser.runtime.onConnect.addListener((port) => {
   // Only listen to the background port
   if (port.name != "bg-comm")
      return;
   contentCommPort = port;
   contentCommPort.postMessage({ message: "You stinky in the background" })
})
