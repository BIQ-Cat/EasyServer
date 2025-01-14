import lib


def test(req: lib.Request):
    res = lib.Response(200)
    res.data.write(b"Hello from Python!\n")
    return res


testController = lib.Controller(test, [b"GET"])

EASYSERVER_MODULE = lib.Module({
    "test": testController,
})
