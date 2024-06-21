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

func main()  {
   logger := log.Default()
   logger.Println("Starting the wasm")

   _ = injectChatArea()

   // document = js.Global().Get("document")
   // body := document.Get("body")
   // style := body.Get("style")
   // style.Set("border", "5px solid green")
   // fmt.Println("Done changing doc")

   // browser := js.Global().Get("browser")
}

func injectChatArea() *chatArea {
   document := js.Global().Get("document")
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
   // Positioning
   chatAreaHtml.Get("style").Set("position", "fixed")
   chatAreaHtml.Get("style").Set("z-index", "100")
   chatAreaHtml.Get("style").Set("height", "30rem")
   chatAreaHtml.Get("style").Set("width", "27.5rem")
   chatAreaHtml.Get("style").Set("bottom", "0px")
   chatAreaHtml.Get("style").Set("right", "0px")

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
   chatAreaHtml.Call("appendChild", inputAreaHtml)

   body.Call("insertAdjacentElement", "afterbegin", chatAreaHtml)
   
   return &chatArea{
      isOpen: false,
      htmlEl: chatAreaHtml,
   }
}
