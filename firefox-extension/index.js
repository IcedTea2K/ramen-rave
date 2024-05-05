console.log("Extension entry")
let port = browser.runtime.connect({ name: "bg-comm"})

port.onMessage.addListener((msg) =>{
   console.log(msg)
})
