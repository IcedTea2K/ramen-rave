(async () => {
   console.log("Extension entry")
   const eventsResource = browser.runtime.getURL("events.js");
   const wasmResource   = browser.runtime.getURL("bin/main.wasm");
   const events         = await import(eventsResource);

   let popupPort;
   let bgPort = browser.runtime.connect({ name: "bg-comm" });

   // Listen to background scripts
   bgPort.onMessage.addListener((msg) =>{
      console.log(msg);
   });

   // Listen to popup
   browser.runtime.onConnect.addListener((port) => {
      if (port.name != "popup-comm")
         return;

      popupPort = port;
      popupPort.onMessage.addListener(processPopUpMsg);
   });

   async function startParty() {
      // Fetch and run Go wasm
      const go = new Go();

      WebAssembly.instantiateStreaming(fetch(wasmResource), go.importObject).then((result) => {
         go.run(result.instance);
      });
   }

   function processPopUpMsg(msg) {
      console.log("PopUp Message: ", msg.message);
      switch (msg.event_code) {
         case events.START_PARTY:
            startParty();
            break;
         case events.STOP_PARTY:
            break;
         default:
            console.log("Uknown event");
      }
   }
})();
