import ctypes
import os
import json
import platform

dll_path = "./easyserver"
if platform.system() == 'Windows':
    dll_path += '.dll'
else:
    dll_path += '.so'

class GoSideError(Exception):
    pass


class GoString(ctypes.Structure):
    _fields_ = [('p', ctypes.c_char_p), ('n', ctypes.c_longlong)]


class GetConfiguration_return(ctypes.Structure):
    _fields_ = [('r0', ctypes.c_char_p), ('r1', ctypes.c_bool)]


def GetDefaultModuleConfiguration(module: str) -> tuple[bytes, bool]:
    lib = ctypes.cdll.LoadLibrary(os.path.abspath(dll_path))
    lib.GetDefaultModuleConfiguration.argtypes = [GoString]
    lib.GetDefaultModuleConfiguration.restype = GetConfiguration_return
    moduleBytes = module.encode()
    moduleGo = GoString(moduleBytes, len(moduleBytes))
    res = lib.GetDefaultModuleConfiguration(moduleGo)
    return (res.r0, res.r1)


def GetEnvironmentConfiguration() -> dict[str, int | float | str | None]:
    lib = ctypes.cdll.LoadLibrary(os.path.abspath(dll_path))
    lib.GetEnvironmentConfiguration.restype = GetConfiguration_return

    res = lib.GetEnvironmentConfiguration()
    if not res.r1:
        raise GoSideError

    return json.loads(res.r0)
