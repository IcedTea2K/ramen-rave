import * as events from "../events.js";

let contentCommPort;
let queryRes = browser.tabs.query({
   currentWindow: true,
   active: true
});

let startPartyButton = document.querySelector("#start-party-btn");
queryRes.then(connectToTab, onErrorConnectToTab);
startPartyButton.addEventListener("click", () => messageContentScript("yummy", events.START_PARTY));

function connectToTab(tabs) {
   if (tabs.length > 0) {
      contentCommPort = browser.tabs.connect(tabs[0].id, {
         name: "popup-comm"
      });
   }
}

function onErrorConnectToTab() {
   console.log("Cannot connect to active tab");
}

function messageContentScript(msgStr, code) {
   let msg = messageMaker(msgStr, code);

   if (contentCommPort === undefined) {
      return;
   }
   contentCommPort.postMessage(msg)
}

function messageMaker(msg, code) {
   return {
      message: msg,
      event_code: code
   }
}
