(async () => {
   console.log("Main script entry")
   const wasmResource   = browser.runtime.getURL("bin/main.wasm");

   // Fetch and run Go wasm
   const go = new Go();

   WebAssembly.instantiateStreaming(fetch(wasmResource), go.importObject).then((result) => {
      go.run(result.instance);
   });
})();
