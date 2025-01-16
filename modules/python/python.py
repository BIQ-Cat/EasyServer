import time
import lib


def test(req: lib.Request):
    res = lib.Response(200)
    res.data = b"Hello from Python!\n"
    return res


def some_fun():
    time.sleep(2)
    print("Fun")


testController = lib.Controller(test, [b"GET"])

EASYSERVER_MODULE = lib.Module({
    "test": testController,
})
