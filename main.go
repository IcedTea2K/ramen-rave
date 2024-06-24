//go:build js && wasm

package main

import (
	"log"
   // "github.com/supabase-community/realtime-go/realtime"
)

func main()  {
   logger := log.Default()
   logger.Println("Starting the wasm")

   chatArea := createChatArea()

   // client := realtime.CreateRealtimeClient("fjveqmouznnqaigdtsxy", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZqdmVxbW91em5ucWFpZ2R0c3h5Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTUxNDEyOTIsImV4cCI6MjAzMDcxNzI5Mn0.ajUGIWg1vP5y4cR5X4OpTapGCTzdq0Oqv7fwWhoWAYQ")
   temp := make(chan struct{})
   select {
      case <-temp:
         break
   }

   chatArea.removeChatArea()
}
