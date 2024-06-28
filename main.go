//go:build js && wasm

package main

import (
	"log"
	"syscall/js"
)

const (
   START_PARTY int = iota
   STOP_PARTY
)

func main()  {
   logger := log.Default()
   logger.Println("Starting the wasm")

   chatArea := createChatArea()

   msgChan  := make(chan []js.Value)
   commPort := js.Global().Call("registerFunction", js.FuncOf(func(this js.Value, args []js.Value) any {
      msgChan <- args
      return nil
   }))
   defer func() {
      commPort.Call("disconnect")
   }()

MAIN_LOOP:
   for {
      msg := <- msgChan
      if len(msg) != 2 {
         logger.Fatal("TYPE ASSERTION FAILED: received more than just data and sender")
      } else if msg[0].Type() != js.TypeObject {
         logger.Fatal("TYPE ASSERTION FAILED: expecting an object for message")
      } else if msg[0].Get("event_code").Type() != js.TypeNumber {
         logger.Fatal("TYPE ASSERTION FAILED: failed to retrieve the event_code")
      }

      switch eventCode := msg[0].Get("event_code").Int(); eventCode {
         case START_PARTY:
            logger.Println("STARTING THE PARTY")
            chatArea.injectChatArea()
            break
         case STOP_PARTY:
            chatArea.removeChatArea()
            break MAIN_LOOP
         default:
            logger.Printf("Ignoring event: %v", eventCode)
            break
      }
   }
   logger.Println("Stopping the wasm")
}
