//go:build js && wasm

package main

import (
	"log"
	"syscall/js"
)

type chatArea struct {
   htmlEl         js.Value
   messageAreaEl  js.Value
   isOpen bool
   currentMember *member
}

// Create the chat area with user's details
func createChatArea() *chatArea {
   var newChatArea chatArea
   document := js.Global().Get("document")

   // Inject html to the page
   chatAreaHtml := document.Call("createElement", "div")
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
      err := newChatArea.currentMember.postMsg(msg.String())
      if err != nil {
         log.Printf("Failed to send message: %v", err)
         return nil
      }
      newChatArea.displayMsg(newChatArea.currentMember.name, msg.String(), true)

      // reset the input area
      inputAreaHtml.Set("value", "")

      return false
   })
   inputAreaHtml.Call("addEventListener", "keypress", inputOnSubmitFunc)
   chatAreaHtml.Call("appendChild", inputAreaHtml)

   // Initialize the chat area
   newChatArea.isOpen = false
   newChatArea.htmlEl = chatAreaHtml
   newChatArea.messageAreaEl = messageAreaHtml

   return &newChatArea
}

func (ca *chatArea) addMember(newMember *member) {
   if ca.currentMember != nil {
      return
   }
   ca.currentMember = newMember
   ca.currentMember.addChatArea(ca)
}

func (ca *chatArea) injectChatArea() {
   document := js.Global().Get("document")
   body := document.Get("body")
   body.Call("insertAdjacentElement", "afterbegin", ca.htmlEl)
   
}

func (ca *chatArea) displayMsg(sender string, msg string, personal bool) {
   newMsgHtml := ca.createNewMsgHtml(sender, msg, personal) 
   ca.messageAreaEl.Call("appendChild", newMsgHtml) 
}

// Create a new message to be added to the chat
func (ca *chatArea) createNewMsgHtml(sender string, msg string, personal bool) js.Value {
   document := js.Global().Get("document")
   msgContainerHtml := document.Call("createElement", "div")
   msgContainerHtml.Get("style").Set("display", "flex")
   msgContainerHtml.Get("style").Set("width", "100%")
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

// Remove the chat area from DOM
func (ca *chatArea) removeChatArea() {
   ca.htmlEl.Call("remove")
}
