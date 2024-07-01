//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/supabase-community/realtime-go/realtime"
   "github.com/go-playground/validator/v10"
)

type member struct {
   name string
   chat *chatArea

   rtClient  *realtime.RealtimeClient
   rtChannel *realtime.RealtimeChannel
}

type message struct {
   Sender  string `json:"sender"    validate:"required"`
   Payload string `json:"payload"   validate:"required"`
}

// Create a member and join the realtime channel
func createMember(name string, partyCode string) (*member, error) {
   // Hard-coded PUBLIC key (it's okay to do so)
   client := realtime.CreateRealtimeClient("fjveqmouznnqaigdtsxy", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZqdmVxbW91em5ucWFpZ2R0c3h5Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTUxNDEyOTIsImV4cCI6MjAzMDcxNzI5Mn0.ajUGIWg1vP5y4cR5X4OpTapGCTzdq0Oqv7fwWhoWAYQ")
   channel, err := client.Channel(partyCode)
   if err != nil {
      return nil, fmt.Errorf("Unable to create channel: %v", err)
   }

   return &member{
      name: name,
      rtClient: client,
      rtChannel: channel,
   }, nil
}

func (me *member) addChatArea(newChatArea *chatArea) {
   if me.chat != nil {
      return
   }
   me.chat = newChatArea
   me.chat.addMember(me)
   err := me.rtChannel.On("broadcast", map[string]string{
      "event" : "message",
   }, me.handleIncomingMessage)
   if err != nil {
      fmt.Printf("AHHHHHHHHHHHHHHHHHHH %v", err)
   }
   fmt.Println("Done adding new chat area")
}

func (me *member) handleIncomingMessage(msg any) {
   log.Println("Handling some messages")
   encodedMsg, err := json.Marshal(msg)
   if err != nil {
      log.Printf("Failed to handle incoming message: %v", err)
      return
   }

   var actualMsg message
   err = json.Unmarshal(encodedMsg, &actualMsg)
   if err != nil {
      log.Printf("Failed to handle incoming message: %v", err)
      return
   }
   log.Printf("Received: %+v", actualMsg)

   validate := validator.New(validator.WithRequiredStructEnabled())
   err = validate.Struct(actualMsg)
   if err != nil {
      log.Printf("Received message is invalid: %v", err)
      return
   }
}

// Join party
func (me *member) joinParty() error {
   log.Println("Joining party")
   ctx, _ := context.WithTimeout(context.Background(), time.Second * 60)
   err := me.rtChannel.Subscribe(ctx)
   if err != nil {
      fmt.Println("Error Joining Party")
      return fmt.Errorf("Unable to join the party: %v", err)
   }
   fmt.Println("Joined Party")
   return nil
}

// Exit party
func (me *member) exitParty() {

}
