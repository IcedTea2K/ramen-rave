//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"
	"time"

	// "github.com/go-playground/validator/v10"
	"github.com/supabase-community/realtime-go/realtime"
)

type member struct {
   name string
   chat *chatArea
   targetVideo js.Value
   manipulateVideo bool

   rtClient  *realtime.RealtimeClient
   rtChannel *realtime.RealtimeChannel
}

type event struct {
   Sender string `json:"sender" validate:"required"`
   EventType string `json:"type" validate:"required"`

   MessageData message `json:"messageData"  validate:"required_if=EventType CHAT"`
   VideoData videoManipulation `json:"videoData" validate:"required_if=EventType VIDEO"`
}

type message struct {
   Payload string `json:"payload"   validate:"required"`
}

// Contains information about manipulating the video
type videoManipulation struct {
   ActionType  string  `json:"actionType" validate:"required"`
   CurrentTime float64 `json:"currentTime" validate:"required_if=ActionType SEEK"`
}

const (
   CHAT_ACTION  string = "CHAT"
   VIDEO_ACTION string = "VIDEO"

   SEEK_ACTION_TYPE  string = "SEEK"
   PLAY_ACTION_TYPE  string = "PLAY"
   PAUSE_ACTION_TYPE string = "PAUSE"

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

   return nil
}

func (me *member) postVideoEvent(action string, timeStamp float64) error {
   ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
   defer cancel()

   err := me.rtChannel.Send(realtime.CustomEvent{
      Event: NEW_EVENT, 
      Payload: event{
         EventType: VIDEO_ACTION,
         Sender: me.name,
         VideoData: videoManipulation{
            ActionType: action,
            CurrentTime: timeStamp,
         },
      },
      Type: "broadcast",
   }, ctx)
   if err != nil {
      return fmt.Errorf("Failed to post video event: %v", err)
   }

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

   // validate := validator.New(validator.WithRequiredStructEnabled())
   // err = validate.Struct(actualEvent)
   // if err != nil {
   //    log.Printf("Received event is invalid: %v", err)
   //    return
   // }

   // Ignore events that we post
   if actualEvent.Sender == me.name {
      return
   }
   
   switch actualEvent.EventType {
      case CHAT_ACTION: 
         me.handleIncomingMessage(actualEvent.Sender, actualEvent.MessageData)
         break
      case VIDEO_ACTION:
         me.handleIncomingVideoEvent(actualEvent.VideoData)
         break
      default:
         log.Printf("Unable to recognize event type: %v", actualEvent.EventType)
         break
   }
}

func (me *member) handleIncomingMessage(sender string, msg message) {
   me.chat.displayMsg(sender, msg.Payload, false)
}

func (me *member) handleIncomingVideoEvent(videoData videoManipulation) {
   if me.targetVideo.IsUndefined() {
      log.Println("Currently doesn't have a target video")
      return
   } 
   me.manipulateVideo = true
   switch videoData.ActionType {
      case PLAY_ACTION_TYPE:
         me.targetVideo.Call("play")
         break
      case PAUSE_ACTION_TYPE:
         me.targetVideo.Call("pause")
         break
      case SEEK_ACTION_TYPE:
         me.targetVideo.Set("currentTime", videoData.CurrentTime)
         break
   }
}

// Join party
func (me *member) joinParty() error {
   if me.chat == nil {
      return fmt.Errorf("Error: Need to associated the member with a chatArea first")
   }

   // Chatting Support
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

   // Video manipulating support
   document := js.Global().Get("document")
   videos := document.Call("getElementsByTagName", "video")
   if videos.Length() < 1 {
      return fmt.Errorf("Unable to find the target video")
   }
   // Assume that the first returned video is the target
   targetVideo := videos.Index(0)
   targetVideo.Call("pause") // Pause the video initially
   targetVideo.Set("onplay", js.FuncOf(func(this js.Value, args []js.Value) any {
      if me.manipulateVideo {
         me.manipulateVideo = false
         return nil
      }
      me.postVideoEvent(PLAY_ACTION_TYPE, 0)
      return nil
   }))
   targetVideo.Set("onpause", js.FuncOf(func(this js.Value, args []js.Value) any {
      if me.manipulateVideo {
         me.manipulateVideo = false
         return nil
      }
      me.postVideoEvent(PAUSE_ACTION_TYPE, 0)
      return nil
   }))
   targetVideo.Set("onseeked", js.FuncOf(func(this js.Value, args []js.Value) any {
      if len(args) != 1 {
         log.Println("There are more than 2 events passed into seek callback")
         return nil
      }
      if me.manipulateVideo {
         me.manipulateVideo = false
         return nil
      }
      me.postVideoEvent(SEEK_ACTION_TYPE, targetVideo.Get("currentTime").Float())
      return nil
   }))
   me.targetVideo = targetVideo

   return nil
}

// Exit party
func (me *member) exitParty() {
   // Unset event listener
   me.targetVideo.Set("onplay", js.Undefined())
   me.targetVideo.Set("onpause", js.Undefined())
   me.targetVideo.Set("onseeked", js.Undefined())
   me.rtChannel.Unsubscribe(context.Background())
}
