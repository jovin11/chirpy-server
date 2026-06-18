package main

import "net/http"

/* 
To build any handler you want, you will manipulate these core steps inside the function:
1. Inspect the Request (*http.Request): The r parameter contains everything the client sent you. You can inspect it to make decisions
2. Write Headers (w.Header()): Headers must be set before you call w.WriteHeader or w.Write.
3. Write the Status Code (w.WriteHeader): Sends the HTTP status code (e.g., 200, 404, 500) to the client. 
4. Write the Body (w.Write): w.Write accepts a slice of bytes ([]byte). 
*/
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
    // 1. Set Headers
     w.Header().Add("Content-Type", "text/plain; charset=utf-8")

    // 2. Set Status Code
     w.WriteHeader(http.StatusOK)

    // 3. Write Body
     w.Write([]byte("OK"))
}