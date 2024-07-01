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

type event struct {
   Sender string `json:"sender" validate:"required"`
   EventType string `json:"type" validate:"required"`

   MessageData message `json:"messageData"  validate:"required_if=EventType CHAT"`
   VideoData videoManipulation `json:"videoData" validate:"required_unless=EventType CHAT"`
}

type message struct {
   Payload string `json:"payload"   validate:"required"`
}

// Contains information about manipulating the video
type videoManipulation struct {

}

const (
   CHAT_ACTION string = "CHAT"
   SEEK_ACTION string = "SEEK"
   PLAY_ACTION string = "PLAY"
   PAUSE_ACTION string = "SEEK"

   NEW_EVENT string = "NEW EVENT"
)

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
}

// Post the message to the realtime connections
func (me *member) postMsg(msg string) error {
   log.Println("Posting new messages")
   ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
   defer cancel()

   err := me.rtChannel.Send(realtime.CustomEvent{
      Event: NEW_EVENT, 
      Payload: event{
         EventType: CHAT_ACTION,
         Sender: me.name,
         MessageData: message{
            Payload: msg,
         }, 
      },
      Type: "broadcast",
   }, ctx)
   if err != nil {
      return fmt.Errorf("Failed to post message: %v", err)
   }
   log.Println("Done posting")

   return nil
}

func (me *member) handleIncomingEvent(newEvent any) {
   encodedEvent, err := json.Marshal(newEvent)
   if err != nil {
      log.Printf("Failed to handle incoming event: %v", err)
      return
   }

   actualEvent := &event{}
   err = json.Unmarshal(encodedEvent, actualEvent)
   if err != nil {
      log.Printf("Failed to handle incoming event: %v", err)
      return
   }

   validate := validator.New(validator.WithRequiredStructEnabled())
   err = validate.Struct(actualEvent)
   if err != nil {
      log.Printf("Received event is invalid: %v", err)
      return
   }

   // Ignore events that we post
   if actualEvent.Sender == me.name {
      return
   }
   
   switch actualEvent.EventType {
      case CHAT_ACTION: 
         me.handleIncomingMessage(actualEvent.Sender, actualEvent.MessageData)
         break
      default:
         log.Printf("Unable to recognize event type: %v", actualEvent.EventType)
         break
   }
}

func (me *member) handleIncomingMessage(sender string, msg message) {
   me.chat.displayMsg(sender, msg.Payload, false)
}

// Join party
func (me *member) joinParty() error {
   log.Println("Joining party")
   if me.chat == nil {
      return fmt.Errorf("Error: Need to associated the member with a chatArea first")
   }

   err := me.rtChannel.On("broadcast", map[string]string{
      "event" : NEW_EVENT,
   }, me.handleIncomingEvent)
   if err != nil {
      return fmt.Errorf("Unable to join the party: %v", err)
   }

   ctx, cancel := context.WithTimeout(context.Background(), time.Second * 60)
   defer cancel()
   err = me.rtChannel.Subscribe(ctx)
   if err != nil {
      return fmt.Errorf("Unable to join the party: %v", err)
   }
   fmt.Println("Joined Party")

   return nil
}

// Exit party
func (me *member) exitParty() {

}
