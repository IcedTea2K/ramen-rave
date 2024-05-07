//go:build js && wasm
package main

import (
   "fmt"
   "syscall/js"
)

func main()  {
   fmt.Println("Oh Yeahhhh, webassembly")

   document := js.Global().Get("document")
   body := document.Get("body")
   style := body.Get("style")
   style.Set("border", "5px solid green")
   fmt.Println("Done changing doc")
}
