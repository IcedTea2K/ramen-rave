// Listen and talk to the content script
(async () => {
   const constantsURL = browser.runtime.getURL("events.js");
   const constants    = await import(constantsURL);

   let popupPorts  = [] ;
   let activeParties = new Map()

   const handlePopupEvent = async (msg) => {
      switch (msg.event_code) {
         case constants.STOP_PARTY:
            {
               const currTabId = await getCurrentTabId()
               if (currTabId === null || !activeParties.has(currTabId)) {
                  console.error("There is no active party on current tab")
                  return
               }
               activeParties.get(currTabId).postMessage({
                  event_code: constants.STOP_PARTY,
                  message: "Stop the party :("
               })
               activeParties.delete(currTabId)
               break;
            }

         case constants.START_PARTY:
            // Do nothing if already ran the program on this tab
            {
               const currTabId = await getCurrentTabId()
               if (currTabId === null || activeParties.has(currTabId))
                  return

               try {
                  await browser.scripting.executeScript({
                     target: {
                        tabId: currTabId,
                     },
                     files: ["../wasm_exec.js", "../index.js"],
                  }) 
               } catch (error) {
                  console.error(`Failed to load index.js: ${error}`)   
               }
               break;
            }
         default:
            console.log(msg)
            break;
      }
   }

   const handleContentEvent = async (port, msg) => {
      switch (msg.event_code) {
         case constants.PARTY_READY:
            const currTabId = await getCurrentTabId()
            if (currTabId === null || activeParties.has(currTabId))
               return
            activeParties.set(currTabId, port)

            port.postMessage({
               event_code: constants.START_PARTY,
               message: "Start the party!!!"
            })
            break;
         default:
            console.error(`Content script should not be sending event_code ${msg.event_code}`)
            break
      }
   }

   const getCurrentTabId = async () => {
      const queriedTabs = await browser.tabs.query({
         currentWindow: true,
         active: true
      });
      if (queriedTabs.length <= 0) {
         console.error("Failed to query tabs")
         return null
      }

      return queriedTabs[0].id
   }

   browser.runtime.onConnect.addListener(async (port) => {
      switch (port.name) {
         case constants.BG_PORT_NAME:
            popupPorts.push(port)

            port.onDisconnect.addListener(() => {
               popupPorts.splice(popupPorts.length - 1, 1)
            })
            port.onMessage.addListener(handlePopupEvent)
            break
         case constants.CONTENT_PORT_NAME:
            const currTabId = await getCurrentTabId()
            if (currTabId === null || activeParties.has(currTabId)) {
               console.error("Port for content script already existed")
               return
            }

            port.onDisconnect.addListener(async () => { 
               const currTabId = await getCurrentTabId()
               if (currTabId === null || activeParties.has(currTabId))
                  return
               activeParties.delete(currTabId)
            })
            port.onMessage.addListener((msg) => handleContentEvent(port, msg))
            break
         default:
            console.error(`Unrecognized port name: ${port.name}`)
            break;
      }
   })
})()
