package api

// #cgo pkg-config: python3
// #cgo LDFLAGS: -lpython3.13
// #include <Python.h>
import "C"
import (
	"io"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/BIQ-Cat/easyserver"
)

func StartPython() {
	C.Py_Initialize()
}

var mu sync.Mutex

func EndPython() {
	C.Py_Finalize()
}

func PythonImportLib() *C.PyObject {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	return PythonImport(path, "lib")
}

func PythonImport(path string, module string) *C.PyObject {
	defer C.PyErr_Print()
	sys := C.PyImport_ImportModule(C.CString("sys"))

	sys_path := C.PyObject_GetAttrString(sys, C.CString("path"))
	defer C.Py_DECREF(sys)

	folder_path := C.PyUnicode_FromString(C.CString(path))
	defer C.Py_DECREF(sys_path)

	C.PyList_Append(sys_path, folder_path)
	defer C.Py_DECREF(folder_path)

	pName := C.PyUnicode_FromString(C.CString(module))
	if pName == nil {
		return nil
	}

	pModule := C.PyImport_Import(pName)
	if pModule == nil {
		return nil
	}
	defer C.Py_DECREF(pName)

	pDict := C.PyModule_GetDict(pModule)
	defer C.Py_DECREF(pModule)

	return pDict
}

func CreateModule(modulePath, moduleName string) (module easyserver.Module, ok bool) {
	defer C.PyErr_Print()

	PythonImportLib()
	pModule := PythonImport(modulePath, moduleName)

	pModuleObject := C.PyDict_GetItemString(pModule, C.CString("EASYSERVER_MODULE"))
	if pModuleObject == nil {
		return
	}

	pRoute := C.PyObject_GetAttrString(pModuleObject, C.CString("route"))
	if pRoute == nil {
		return
	}
	defer C.Py_DECREF(pRoute)

	var pKey, pController *C.PyObject
	var pos C.Py_ssize_t = 0
	module.Route = make(easyserver.Route)
	for C.PyDict_Next(pRoute, &pos, &pKey, &pController) != 0 {
		if pController == nil || pKey == nil {
			return
		}

		pByteKey := C.PyUnicode_EncodeLocale(pKey, nil)
		if pByteKey == nil {
			return
		}
		defer C.Py_DECREF(pByteKey)

		key := C.PyBytes_AsString(pByteKey)
		if key == nil {
			return
		}

		pControllerCopy := &(*pController)

		controller := easyserver.Controller{}

		pHandler := C.PyObject_GetAttrString(pControllerCopy, C.CString("handler"))
		if pHandler == nil {
			return
		}

		controller.Handler = handlePython(moduleName, modulePath, C.GoString(key))

		pMethods := C.PyObject_GetAttrString(pController, C.CString("methods"))
		if pMethods == nil {
			return
		}
		defer C.Py_DECREF(pMethods)

		methods, good := makeStrSlice(pMethods)
		if !good {
			return
		}
		controller.Methods = methods

		pSchemas := C.PyObject_GetAttrString(pController, C.CString("schemas"))
		if pSchemas == nil {
			return
		}
		defer C.Py_DECREF(pSchemas)

		schemas, good := makeStrSlice(pSchemas)
		if !good {
			return
		}
		controller.Schemas = schemas

		pHeaders := C.PyObject_GetAttrString(pController, C.CString("headers"))
		if pHeaders == nil {
			return
		}
		defer C.Py_DECREF(pHeaders)

		headers, good := makeHeadersMap(pHeaders)
		if !good {
			return
		}
		controller.Headers = headers

		module.Route[C.GoString(key)] = controller
	}

	ok = true
	return
}

func makeStrSlice(pList *C.PyObject) ([]string, bool) {
	if pList == C.Py_None {
		return nil, true
	}
	length := int(C.PyList_Size(pList))
	if length == -1 {
		return nil, false
	}
	res := make([]string, length)
	for i := 0; i < length; i++ {
		item := C.PyList_GetItem(pList, C.Py_ssize_t(i))
		if item == nil {
			return nil, false
		}
		defer C.Py_DECREF(item)

		str := C.PyBytes_AsString(item)
		if str == nil {
			return nil, false
		}
		res[i] = C.GoString(str)
	}

	return res, true
}

func makeHeadersMap(pDict *C.PyObject) (map[string]string, bool) {
	if pDict == C.Py_None {
		return nil, true
	}

	res := make(map[string]string)
	var pKey, pValue *C.PyObject
	var pos C.Py_ssize_t = 0
	for C.PyDict_Next(pDict, &pos, &pKey, &pValue) != 0 {
		if pKey == nil || pValue == nil {
			return nil, false
		}

		pByteKey := C.PyUnicode_EncodeLocale(pKey, nil)
		if pByteKey == nil {
			return nil, false
		}
		defer C.Py_DECREF(pByteKey)

		key := C.PyBytes_AsString(pByteKey)
		if key == nil {
			return nil, false
		}

		value := C.PyBytes_AsString(pValue)
		if value == nil {
			return nil, false
		}

		res[C.GoString(key)] = C.GoString(value)
	}
	return res, true
}

func handlePython(moduleName string, modulePath string, controllerName string) http.Handler {
	defer C.PyErr_Print()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers, status, data, ok := wrapPythonHandler(moduleName, modulePath, controllerName, r)
		if !ok {
			w.WriteHeader(500)
			return
		}

		for key, header := range headers {
			w.Header().Set(key, header)
		}

		w.WriteHeader(status)

		io.WriteString(w, data)
	})
}

func wrapPythonHandler(moduleName string, modulePath string, controllerName string, r *http.Request) (headers map[string]string, status int, data string, ok bool) {
	mu.Lock()
	defer mu.Unlock()

	StartPython()
	defer EndPython()

	defer C.PyErr_Print()

	pLib := PythonImportLib()
	pModule := PythonImport(modulePath, moduleName)

	pModuleObject := C.PyDict_GetItemString(pModule, C.CString("EASYSERVER_MODULE"))
	if pModuleObject == nil {
		return
	}

	pRoute := C.PyObject_GetAttrString(pModuleObject, C.CString("route"))
	if pRoute == nil {
		return
	}
	defer C.Py_DECREF(pRoute)

	pController := C.PyDict_GetItemString(pRoute, C.CString(controllerName))
	if pController == nil {
		return
	}

	pHandler := C.PyObject_GetAttrString(pController, C.CString("handler"))
	if pHandler == nil {
		return
	}
	defer C.Py_DECREF(pHandler)

	pClassRequest := C.PyDict_GetItemString(pLib, C.CString("Request"))
	if pClassRequest == nil {
		return
	}

	pClassURL := C.PyDict_GetItemString(pLib, C.CString("URL"))
	if pClassURL == nil {
		return
	}

	pRequestArgs := C.PyTuple_New(7)
	if pRequestArgs == nil {
		return
	}
	defer C.Py_DECREF(pRequestArgs)

	pMethod := C.PyUnicode_FromString(C.CString(r.Method))
	if pMethod == nil {
		return
	}
	defer C.Py_DECREF(pMethod)

	if C.PyTuple_SetItem(pRequestArgs, 0, pMethod) == -1 {
		return
	}

	pURLArgs := C.PyTuple_New(5)
	if pURLArgs == nil {
		return
	}
	defer C.Py_DECREF(pURLArgs)

	pScheme := C.PyUnicode_FromString(C.CString(r.URL.Scheme))
	if pScheme == nil {
		return
	}
	defer C.Py_DECREF(pScheme)

	if C.PyTuple_SetItem(pURLArgs, 0, pScheme) == -1 {
		return
	}

	pHost := C.PyUnicode_FromString(C.CString(r.URL.Hostname()))
	if pHost == nil {
		return
	}
	defer C.Py_DECREF(pHost)

	if C.PyTuple_SetItem(pURLArgs, 1, pHost) == -1 {
		return
	}

	pPort := C.PyUnicode_FromString(C.CString(r.URL.Port()))
	if pPort == nil {
		return
	}
	defer C.Py_DECREF(pPort)

	if C.PyTuple_SetItem(pURLArgs, 2, pPort) == -1 {
		return
	}

	pPath := C.PyUnicode_FromString(C.CString(r.URL.Path))
	if pPath == nil {
		return
	}
	defer C.Py_DECREF(pPath)

	if C.PyTuple_SetItem(pURLArgs, 3, pPath) == -1 {
		return
	}

	pQuery := C.PyDict_New()
	for key, values := range r.URL.Query() {
		pList := C.PyList_New(0)
		for _, value := range values {
			pStr := C.PyUnicode_FromString(C.CString(value))
			if pStr == nil {
				return
			}

			C.PyList_Append(pList, pStr)
		}

		C.PyDict_SetItemString(pQuery, C.CString(key), pList)
	}

	if C.PyTuple_SetItem(pURLArgs, 4, pQuery) == -1 {
		return
	}

	pURL := C.PyObject_CallObject(pClassURL, pURLArgs)
	if pURL == nil {
		return
	}
	defer C.Py_DECREF(pURL)

	if C.PyTuple_SetItem(pRequestArgs, 1, pURL) == -1 {
		return
	}

	pProtocol := C.PyUnicode_FromString(C.CString(r.Proto))
	if pProtocol == nil {
		return
	}
	defer C.Py_DECREF(pProtocol)

	if C.PyTuple_SetItem(pRequestArgs, 2, pProtocol) == -1 {
		return
	}

	pHeaders := C.PyDict_New()
	for key := range r.Header {
		pHeader := C.PyUnicode_FromString(C.CString(r.Header.Get(key)))
		if pHeader == nil {
			return
		}
		defer C.Py_DECREF(pHeader)

		if C.PyDict_SetItemString(pHeaders, C.CString(key), pHeader) == -1 {
			return
		}
	}

	if C.PyTuple_SetItem(pRequestArgs, 3, pHeaders) == -1 {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	pBody := C.PyBytes_FromString(C.CString(string(body)))
	if pBody == nil {
		return
	}
	defer C.Py_DECREF(pBody)

	if C.PyTuple_SetItem(pRequestArgs, 4, pBody) == -1 {
		return
	}

	pContext := C.PyDict_New()
	if pContext == nil {
		return
	}
	defer C.Py_DECREF(pContext)

	if C.PyTuple_SetItem(pRequestArgs, 5, pContext) == -1 {
		return
	}

	pCookies := C.PyDict_New()
	if pCookies == nil {
		return
	}
	defer C.Py_DECREF(pCookies)

	if C.PyTuple_SetItem(pRequestArgs, 6, pCookies) == -1 {
		return
	}

	pRequest := C.PyObject_CallObject(pClassRequest, pRequestArgs)
	if pRequest == nil {
		return
	}
	defer C.Py_DECREF(pRequest)

	pHandlerArgs := C.PyTuple_New(1)
	if pHandlerArgs == nil {
		return
	}
	defer C.Py_DECREF(pHandlerArgs)

	if C.PyTuple_SetItem(pHandlerArgs, 0, pRequest) == -1 {
		return
	}

	pResponse := C.PyObject_CallObject(pHandler, pHandlerArgs)
	if pResponse == nil {
		return
	}
	defer C.Py_DECREF(pResponse)

	pStatus := C.PyObject_GetAttrString(pResponse, C.CString("status"))
	if pStatus == nil {
		return
	}
	defer C.Py_DECREF(pStatus)

	status = int(C.PyLong_AsInt(pStatus))

	pRespHeaders := C.PyObject_GetAttrString(pResponse, C.CString("headers"))
	if pRespHeaders == nil {
		return
	}
	defer C.Py_DECREF(pRespHeaders)

	headers, ok = makeHeadersMap(pRespHeaders)
	if !ok {
		return
	}

	ok = false

	pData := C.PyObject_GetAttrString(pResponse, C.CString("data"))
	if pData == nil {
		return
	}
	defer C.Py_DECREF(pData)

	cData := C.PyBytes_AsString(pData)
	if cData == nil {
		return
	}
	data = C.GoString(cData)

	ok = true
	return
}
