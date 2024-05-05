console.log("Extension entry")
let popupPort;
let bgPort = browser.runtime.connect({ name: "bg-comm" });

bgPort.onMessage.addListener((msg) =>{
   console.log(msg);
});

browser.runtime.onConnect.addListener((port) => {
   if (port.name != "popup-comm")
      return;

   popupPort = port;
   popupPort.onMessage.addListener((msg) => {
      console.log("From popup script");
      console.log(msg);
   });
});
