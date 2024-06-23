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
   inputOnSubmitFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
      if len(args) == 0 || !args[0].Truthy() {
         return nil
      }

      event := args[0]
      if !event.Get("key").Equal(js.ValueOf("Enter")) || 
         event.Get("shiftKey").Truthy() {
         return nil
      }

      msg := this.Get("value")
      if !msg.Truthy() {
         return nil
      }
      newMsgHtml := createNewMsgHtml(msg.String())
      messageAreaHtml.Call("appendChild", newMsgHtml)

      // reset the input area
      inputAreaHtml.Set("value", "")

      return nil
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
func createNewMsgHtml(msg string) js.Value {
   msgHtml := document.Call("createElement", "div")
   msgContentHtml := document.Call("createElement", "p")
   msgContentHtml.Set("textContent", msg)
   msgHtml.Call("appendChild", msgContentHtml)
   return msgHtml
}
