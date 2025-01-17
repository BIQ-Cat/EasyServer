#ifndef CXX_MODULE_H
#define CXX_MODULE_H

#include <stddef.h>

typedef struct
{
    int status;
    char **headerNames;
    char **headerValues;
    size_t headerLen;
    char *data;
} Response;

void GetControllers(char *controller[]);
size_t LenControllers();

void GetControllerMethods(char *controller, char *methods[]);
size_t LenControllerMethods(char *controller);

void GetControllerHeaderNames(char *controller, char *names[]);
char *GetControllerHeader(char *controller, char *name);
size_t LenControllerHeaders(char *controller);

void GetControllerSchemas(char *controller, char *schemas[]);
size_t LenControllerSchemas(char *controller);

void CallController(char *controller, char *method, char *proto, char *body, char *urlScheme, char *urlHost, char *urlPort, char *urlPath, char **queryKeys, size_t queryLen, char ***queryValues, size_t *queryValueLen, char **headerNames, char **headerValues, size_t headerLen, Response *resp);

#endif // CXX_MODULE_H