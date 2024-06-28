// Listen and talk to the content script
(async () => {
   const constantsURL = browser.runtime.getURL("events.js");
   const constants    = await import(constantsURL);
   let currActiveTab  = null;
   let contentPort = null;
   let popupPorts  = [] ;

   const handlePopupEvent = async (msg) => {
      switch (msg.event_code) {
         case constants.STOP_PARTY:
            if (contentPort !== null) {
               contentPort.postMessage(msg)
            }
            currActiveTab = null
            break;

         case constants.START_PARTY:
            if (currActiveTab != null)
               return
            const activeTabs = await browser.tabs.query({
               currentWindow: true,
               active: true
            });
            if (activeTabs.length <= 0) {
               console.error("Failed to query tabs")
               return
            }

            currActiveTab = activeTabs[0]

            try {
               await browser.scripting.executeScript({
                  target: {
                     tabId: currActiveTab.id,
                  },
                  files: ["../wasm_exec.js", "../index.js"],
               }) 
            } catch (error) {
               console.error(`Failed to load index.js: ${error}`)   
            }
            break;
         default:
            console.log(msg)
            break;
      }
   }

   const handleContentEvent = (msg) => {
      switch (msg.event_code) {
         case constants.PARTY_READY:
            if (contentPort === null){
               console.error("ASSERTION FAILED: contentPort must already set and connected")
               return
            }
            contentPort.postMessage({
               event_code: constants.START_PARTY,
               message: "Start the party!!!"
            })
            break;
         default:
            console.error(`Content script should not be sending event_code ${msg.event_code}`)
            break
      }
   }

   browser.runtime.onConnect.addListener((port) => {
      switch (port.name) {
         case constants.BG_PORT_NAME:
            popupPorts.push(port)

            port.onDisconnect.addListener(() => {
               popupPorts.splice(popupPorts.length - 1, 1)
            })
            port.onMessage.addListener(handlePopupEvent)
            break
         case constants.CONTENT_PORT_NAME:
            if (contentPort !== null) {
               console.error("Port for content script already existed")
               return
            }
            contentPort = port
            port.onDisconnect.addListener(() => { contentPort = null; console.log("REMOVING CONTENT PORT") })
            port.onMessage.addListener(handleContentEvent)
            break
         default:
            console.error(`Unrecognized port name: ${port.name}`)
            break;
      }
   })
})()
