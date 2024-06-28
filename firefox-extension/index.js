// Register the function to handle messages from the background script
function registerFunction(func) {
   const port = browser.runtime.connect({ name: "content-script-comm" }) ;
   const constantsURL = browser.runtime.getURL("events.js");

   import(constantsURL)
      .then((constants) => {
         if (port === null) {
            console.error("Communication port with background hasn't been established yet")
            return
         } else if (port.onMessage.hasListener(func)) {
            console.error("The function has already been registered for port listening")
            return
         }
         port.onMessage.addListener(func)
         port.postMessage({
            event_code: constants.PARTY_READY, 
            message: "Party is ready to start"
         })
      })
      .catch((err) => {
         console.error(`Failed to load constants ${err}`)
      });

   // return the communication port for message posting and cleanup 
   return port
}

(async () => {
   console.log("Main script entry")
   const wasmResource   = browser.runtime.getURL("bin/main.wasm");

   // Fetch and run Go wasm
   const go = new Go();

   await WebAssembly.instantiateStreaming(fetch(wasmResource), go.importObject).then((result) => {
      go.run(result.instance);
   });
})();
