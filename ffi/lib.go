package main

import "C"
import (
	"net/http"
	"runtime"
	"time"
	"unsafe"

	"github.com/BIQ-Cat/easyserver"
)

var pinner runtime.Pinner

const (
	// export OK
	OK = iota
	//export ERR_ABORT_HANDLER
	ERR_ABORT_HANDLER = 1 << (iota - 1)
	//export ERR_BODY_NOT_ALLOWED
	ERR_BODY_NOT_ALLOWED

	//export ERR_BODY_READ_AFTER_CLOSE
	ERR_BODY_READ_AFTER_CLOSE

	//export ERR_CONTENT_LENGTH
	ERR_CONTENT_LENGTH

	//export ERR_HANDLER_TIMEOUT
	ERR_HANDLER_TIMEOUT

	//export ERR_LINE_TOO_LONG
	ERR_LINE_TOO_LONG

	//export ERR_SERVER_CLOSED
	ERR_SERVER_CLOSED

	//export ERR_INTERNAL
	ERR_INTERNAL
)

var httpErrorCode = map[error]int{
	nil:                        OK,
	http.ErrAbortHandler:       ERR_ABORT_HANDLER,
	http.ErrBodyNotAllowed:     ERR_BODY_NOT_ALLOWED,
	http.ErrBodyReadAfterClose: ERR_BODY_READ_AFTER_CLOSE,
	http.ErrContentLength:      ERR_CONTENT_LENGTH,
	http.ErrHandlerTimeout:     ERR_HANDLER_TIMEOUT,
	http.ErrLineTooLong:        ERR_LINE_TOO_LONG,
	http.ErrServerClosed:       ERR_SERVER_CLOSED,
}

//export NewServer
func NewServer() (srv unsafe.Pointer) {
	server := new(easyserver.Router)
	server.Modules = make(map[string]easyserver.Module)
	srv = unsafe.Pointer(server)
	pinner.Pin(srv)
	return srv
}

//export RunServer
func RunServer(addr *C.char, writeTimeout C.int, readTimeout C.int, srv unsafe.Pointer) (errno int) {
	server := (*easyserver.Router)(srv)
	httpServer := &http.Server{
		Handler:      server,
		Addr:         C.GoString(addr),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}

	errno, ok := httpErrorCode[httpServer.ListenAndServe()]
	if !ok {
		errno = ERR_INTERNAL
	}
	return errno
}

// export RunServerTLS
func RunServerTLS(addr *C.char, writeTimeout C.int, readTimeout C.int, certFile, keyFile *C.char, srv unsafe.Pointer) (errno int) {
	server := (*easyserver.Router)(srv)
	httpServer := &http.Server{
		Handler:      server,
		Addr:         C.GoString(addr),
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}
	errno, ok := httpErrorCode[httpServer.ListenAndServeTLS(C.GoString(certFile), C.GoString(keyFile))]
	if !ok {
		errno = ERR_INTERNAL
	}
	return errno
}
