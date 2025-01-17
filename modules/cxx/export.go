package cxx

// #include <stddef.h>
// #include <stdlib.h>
// #include "cxx.h"
import "C"
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"unsafe"

	"github.com/BIQ-Cat/easyserver"
	"github.com/BIQ-Cat/easyserver/internal/router"
)

var pinner runtime.Pinner

func init() {
	var module easyserver.Module

	controllersLen := C.LenControllers()
	if controllersLen == 0 {
		return
	}

	module.Route = make(easyserver.Route)

	controllerNames := make([]*C.char, controllersLen)
	C.GetControllers((**C.char)(unsafe.Pointer(&controllerNames[0])))
	for _, controllerName := range controllerNames {
		var controller easyserver.Controller

		methodsLen := C.LenControllerMethods(controllerName)
		if methodsLen != 0 {
			methods := make([]*C.char, methodsLen)
			controller.Methods = make([]string, methodsLen)
			C.GetControllerMethods(controllerName, (**C.char)(unsafe.Pointer(&methods[0])))
			for i, method := range methods {
				controller.Methods[i] = C.GoString(method)
			}
		}

		headersLen := C.LenControllerHeaders(controllerName)
		if headersLen != 0 {
			headerNames := make([]*C.char, headersLen)
			controller.Headers = make(map[string]string)
			C.GetControllerHeaderNames(controllerName, (**C.char)(unsafe.Pointer(&headerNames[0])))
			for _, headerName := range headerNames {
				controller.Headers[C.GoString(headerName)] = C.GoString(C.GetControllerHeader(controllerName, headerName))
			}
		}

		schemasLen := C.LenControllerSchemas(controllerName)
		if schemasLen != 0 {
			schemas := make([]*C.char, schemasLen)
			controller.Schemas = make([]string, schemasLen)
			C.GetControllerSchemas(controllerName, (**C.char)(unsafe.Pointer(&schemas[0])))
			for i, schema := range schemas {
				controller.Schemas[i] = C.GoString(schema)
			}
		}

		controller.Handler = makeHandler(C.GoString(controllerName))

		fmt.Println(controller)

		module.Route[C.GoString(controllerName)] = controller
	}

	router.DefaultRouter.Modules["cxx"] = module
}

func makeHandler(controllerName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		queryLen := len(r.URL.Query())
		queryKeys := make([]*C.char, queryLen)
		goQueryValues := make([][]*C.char, queryLen)
		queryValues := make([]**C.char, queryLen)
		queryValueLens := make([]C.size_t, queryLen)
		i := 0
		for key, values := range r.URL.Query() {
			valueLen := len(values)
			queryValueLens[i] = C.size_t(valueLen)
			queryKeys[i] = C.CString(key)

			goQueryValues[i] = make([]*C.char, valueLen)
			for j, value := range values {
				goQueryValues[i][j] = C.CString(value)
			}

			pinner.Pin(&goQueryValues[i][0])

			queryValues[i] = (**C.char)(unsafe.Pointer(&goQueryValues[i][0]))
			i++
		}

		var cQueryKeys **C.char
		var cQueryValues ***C.char
		var cQueryValuesLen *C.size_t
		if queryLen > 0 {
			cQueryKeys = (**C.char)(unsafe.Pointer(&queryKeys[0]))
			cQueryValues = (***C.char)(unsafe.Pointer(&queryValues[0]))
			cQueryValuesLen = (*C.size_t)(unsafe.Pointer(&queryValueLens[0]))
		}

		headerLen := len(r.Header)
		headerNames := make([]*C.char, headerLen)
		headerValues := make([]*C.char, headerLen)
		i = 0
		for name := range r.Header {
			headerNames[i] = C.CString(name)
			headerValues[i] = C.CString(r.Header.Get(name))
		}

		var cHeaderNames, cHeaderValues **C.char
		if headerLen > 0 {
			cHeaderNames = (**C.char)(unsafe.Pointer(&headerNames[0]))
			cHeaderValues = (**C.char)(unsafe.Pointer(&headerValues[0]))
		}

		var resp C.Response
		C.CallController(
			C.CString(controllerName),

			C.CString(r.Method),
			C.CString(r.Proto),
			C.CString(string(body)),

			C.CString(r.URL.Scheme),
			C.CString(r.URL.Hostname()),
			C.CString(r.URL.Port()),
			C.CString(r.URL.Path),

			cQueryKeys,
			C.size_t(queryLen),
			cQueryValues,
			cQueryValuesLen,

			cHeaderNames,
			cHeaderValues,
			C.size_t(headerLen),

			(*C.Response)(unsafe.Pointer(&resp)),
		)

		respHeaderNames := unsafe.Slice(resp.headerNames, int(resp.headerLen))
		respHeaderValues := unsafe.Slice(resp.headerValues, int(resp.headerLen))
		for i, name := range respHeaderNames {
			value := C.GoString(respHeaderValues[i])
			w.Header().Add(C.GoString(name), value)
		}

		w.WriteHeader(int(resp.status))

		io.WriteString(w, C.GoString(resp.data))
		C.free(unsafe.Pointer(resp.data))
	})
}
