import * as constants from "../events.js";


(async () => {
   let bgPort = browser.runtime.connect({ name: constants.BG_PORT_NAME });
   // Start button will trigger port registration with the background
   let startPartyButton = document.getElementById("start-party-btn");
   startPartyButton.addEventListener("click", () => {
      bgPort.postMessage({
         event_code: constants.START_PARTY,
         message: "",
      })
   });

   // Stop button will trigger port deregistration with the background
   let stopPartyButton = document.getElementById("stop-party-btn");
   stopPartyButton.addEventListener("click", () => {
      bgPort.postMessage({
         event_code: constants.STOP_PARTY,
         message: "",
      })
   });

   const activeTabs = await browser.tabs.query({
      currentWindow: true,
      active: true
   });
   if (activeTabs.length <= 0) {
      throw "Cannot query active tabs"
   }
})()
