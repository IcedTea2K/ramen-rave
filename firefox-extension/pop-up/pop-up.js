console.log("Running pop up")
let contentCommPort;

browser.runtime.onConnect.addListener((port) => {
   // Only listen to the pop-up port
   if (port.name != "popup-comm")
      return;
   contentCommPort = port;
   contentCommPort.postMessage({ message: "Just popped up" })
})
