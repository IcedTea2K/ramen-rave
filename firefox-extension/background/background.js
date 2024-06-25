// Listen and talk to the content script
(async () => {
   const constantsURL = browser.runtime.getURL("events.js");
   const constants    = await import(constantsURL);

   let ports = new Map()

   browser.runtime.onConnect.addListener((port) => {
      console.log(port.sender)
      // Only listen to the background port
      if (port.name != "bg-comm" || ports.has(port.sender.contextId)) {
         port.postMessage({ 
            event_code: constants.PORT_EXISTS,
            message: "Port already exists: " + port.sender.contextId
         })
         return;
      }

      ports.set(port.sender.contextId, port)
      port.postMessage({ 
         event_code: constants.PORT_REGISTERED,
         message: "Communication with background has been established with " + port.sender.contextId
      })

      port.onMessage.addListener(() => {
         switch (msg.event_code) {
            case constants.STOP_PARTY:
               ports.delete(port.sender.contextId)
               break;
            default:
               break;
         }
      })
   })
})()
