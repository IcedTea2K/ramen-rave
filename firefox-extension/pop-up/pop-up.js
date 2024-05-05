let contentCommPort;

let queryRes = browser.tabs.query({
   currentWindow: true,
   active: true
});

queryRes.then(connectToTab, onErrorConnectToTab);

function connectToTab(tabs) {
   if (tabs.length > 0) {
      contentCommPort = browser.tabs.connect(tabs[0].id, {
         name: "popup-comm"
      });

      contentCommPort.postMessage({ message: "Just popped up" });
      console.log(tabs)
   }
}

function onErrorConnectToTab() {
   console.log("Cannot connect to active tab");
}
