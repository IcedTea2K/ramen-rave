import * as constants from "../events.js";


(async () => {
   let bgPort = null;
   let startParty = async () => {
      let activeTabs = await browser.tabs.query({
         currentWindow: true,
         active: true
      });
      if (activeTabs.length > 0) {
         browser.tabs.onRemoved.addListener((tabId, _) => {
            if (tabId == activeTabs[0].id) {
               stopParty()
            }
         })
      }

      const mainScript = {
         id: "main-script",
         js: ["../wasm_exec.js", "../index.js"],
         matches: ["<all_urls>"],
      }

      try {
         await browser.scripting.registerContentScripts([mainScript]); 
      } catch (error) {
         console.error(`Failed to start the party: ${error}`) 
      }
   }

   let stopParty = () => {
      if (bgPort === null)
      return

      bgPort.postMessage({
         message: "",
         event_code: constants.STOP_PARTY
      })
   }

   let startPartyButton = document.getElementById("start-party-btn");
   startPartyButton.addEventListener("click", () => {
      console.log("trying to start party")
      bgPort = browser.runtime.connect({ name: constants.BG_PORT_NAME });

      // Initial confirmation with the background script
      bgPort.onMessage.addListener((msg) => {
         switch (msg.event_code) {
            case constants.PORT_REGISTERED:
               startParty() 
               break;
            case constants.PORT_EXISTS:
            default:
               bgPort.disconnect()
         }
      })
   });

   let stopPartyButton = document.getElementById("stop-party-btn");
   stopPartyButton.addEventListener("click", stopParty);
})()
