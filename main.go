//go:build js && wasm

package main

import (
	"log"
	"syscall/js"
)

type chatArea struct {
   htmlEl js.Value
   isOpen bool
}

const (
   chatAreaTag = "chat-area"
)

var document js.Value

func main()  {
   logger := log.Default()
   logger.Println("Starting the wasm")

   document = js.Global().Get("document")
   _ = injectChatArea()

   // document = js.Global().Get("document")
   // body := document.Get("body")
   // style := body.Get("style")
   // style.Set("border", "5px solid green")
   // fmt.Println("Done changing doc")

   // browser := js.Global().Get("browser")
   temp := make(chan struct{})
   select {
      case <-temp:
         break
   }
}

func injectChatArea() *chatArea {
   body := document.Get("body")

   // Inject html to the page
   chatAreaHtml := document.Call("createElement", "div")
   chatAreaHtml.Get("classList").Call("add", chatAreaTag)
   // Style
   chatAreaHtml.Get("style").Set("background-color", "rgba(109, 72, 72, 0.40)")
   chatAreaHtml.Get("style").Set("backdrop-filter", "blur(16.25px)")
   chatAreaHtml.Get("style").Set("display", "flex")
   chatAreaHtml.Get("style").Set("flex-direction", "column")
   chatAreaHtml.Get("style").Set("align-items", "center")
   chatAreaHtml.Get("style").Set("font-family", "Droid Sans")
   chatAreaHtml.Get("style").Set("padding", "0.8rem 0.7rem ")
   // Positioning
   chatAreaHtml.Get("style").Set("position", "fixed")
   chatAreaHtml.Get("style").Set("z-index", "100")
   chatAreaHtml.Get("style").Set("height", "30rem")
   chatAreaHtml.Get("style").Set("width", "27.5rem")
   chatAreaHtml.Get("style").Set("bottom", "0px")
   chatAreaHtml.Get("style").Set("right", "0px")

   // Adding messages area
   messageAreaHtml := document.Call("createElement", "div")
   messageAreaHtml.Get("style").Set("height", "100%")
   messageAreaHtml.Get("style").Set("width",  "100%")
   messageAreaHtml.Get("style").Set("overflow",  "scroll")
   chatAreaHtml.Call("appendChild", messageAreaHtml)

   // Adding input area
   inputAreaHtml := document.Call("createElement", "textarea")
   inputAreaHtml.Set("placeholder", "What's on your mind?.....")
   inputAreaHtml.Get("style").Set("resize", "none")
   inputAreaHtml.Get("style").Set("width", "calc(100% - 3rem)")
   inputAreaHtml.Get("style").Set("height", "2.6rem")
   inputAreaHtml.Get("style").Set("border-radius", "0.75rem")
   inputAreaHtml.Get("style").Set("background", "rgba(99, 99, 99, 0.38)")
   inputAreaHtml.Get("style").Set("box-shadow", "none")
   inputAreaHtml.Get("style").Set("outline", "none")
   inputAreaHtml.Get("style").Set("border", "none")
   inputAreaHtml.Get("style").Set("padding", "0.6rem 0.8rem")
   inputAreaHtml.Get("style").Set("color", "#9B9595")
   inputAreaHtml.Get("style").Set("font-size", "0.9375rem")
   inputAreaHtml.Set("rows", 4)
   // TODO: REMOVE LATER
   temp := true
   inputOnSubmitFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
      if len(args) == 0 || !args[0].Truthy() {
         return true
      }

      event := args[0]
      if !event.Get("key").Equal(js.ValueOf("Enter")) || 
         event.Get("shiftKey").Truthy() {
         return true
      }

      msg := this.Get("value")
      if !msg.Truthy() {
         return true
      }
      newMsgHtml := createNewMsgHtml("Dawg", msg.String(), temp)
      temp = !temp
      messageAreaHtml.Call("appendChild", newMsgHtml)

      // reset the input area
      inputAreaHtml.Set("value", "")

      return false
   })
   inputAreaHtml.Call("addEventListener", "keypress", inputOnSubmitFunc)
   chatAreaHtml.Call("appendChild", inputAreaHtml)

   body.Call("insertAdjacentElement", "afterbegin", chatAreaHtml)
   
   return &chatArea{
      isOpen: false,
      htmlEl: chatAreaHtml,
   }
}

// Create a new message to be added to the chat
func createNewMsgHtml(sender string, msg string, personal bool) js.Value {
   msgContainerHtml := document.Call("createElement", "div")
   msgContainerHtml.Get("style").Set("display", "flex")
   msgContainerHtml.Get("style").Set("width", "100%")
   msgContainerHtml.Get("style").Set("height", "3.2rem")
   msgContainerHtml.Get("style").Set("margin", "0 0 1.9rem 0")

   msgHtml := document.Call("createElement", "div")
   msgHtml.Get("style").Set("display", "flex")
   msgHtml.Get("style").Set("flex-direction", "column")
   msgHtml.Get("style").Set("width", "75%")
   msgHtml.Get("style").Set("justify-content", "center")
   msgHtml.Get("style").Set("gap", "0.5rem")
   msgHtml.Get("style").Set("border-radius", "0.5625rem")
   msgHtml.Get("style").Set("padding", "0.3rem 0.5rem")
   msgHtml.Get("style").Set("font-size", "0.7rem")

   msgSenderHtml := document.Call("createElement", "h3")
   msgSenderHtml.Set("textContent", sender)
   msgSenderHtml.Get("style").Set("margin", "0")
   msgSenderHtml.Get("style").Set("color", "rgba(0, 0, 0, 0.53)")
   msgHtml.Call("appendChild", msgSenderHtml)

   msgContentHtml := document.Call("createElement", "p")
   msgContentHtml.Set("textContent", msg)
   msgContentHtml.Get("style").Set("margin", "0")
   msgContentHtml.Get("style").Set("color", "#9B9595")
   msgHtml.Call("appendChild", msgContentHtml)

   // Change style between user's message vs other people's mesaage
   if personal {
      msgContainerHtml.Get("style").Set("justify-content", "flex-end")
      msgHtml.Get("style").Set("background", "rgba(201, 49, 94, 0.29)")
   } else {
      msgContainerHtml.Get("style").Set("justify-content", "flex-start")
      msgHtml.Get("style").Set("background", "rgba(255, 168, 168, 0.29)")
   }

   msgContainerHtml.Call("appendChild", msgHtml)

   return msgContainerHtml
}
