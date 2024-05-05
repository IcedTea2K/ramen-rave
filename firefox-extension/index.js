console.log("Extension entry")
let bgPort = browser.runtime.connect({ name: "bg-comm" })
let popupPort = browser.runtime.connect({ name: "popup-comm" })

bgPort.onMessage.addListener((msg) =>{
   console.log("From background script")
   console.log(msg)
})

popupPort.onMessage.addListener((msg) => {
   console.log("From popup script")
   console.log(msg)
})
