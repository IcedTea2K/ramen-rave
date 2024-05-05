// Listen and talk to the content script
let contentCommPort;

browser.runtime.onConnect.addListener((port) => {
   // Only listen to the background port
   if (port.name != "bg-comm")
      return;
   contentCommPort = port;
   contentCommPort.postMessage({ message: "You stinky in the background" })
})
