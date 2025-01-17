#include "cxx.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

void GetControllers(char *controller[])
{
    controller[0] = "test";

    char *methods[1];
    GetControllerMethods("test", methods);
}

size_t LenControllers()
{
    return 1;
}

void GetControllerMethods(char *controller, char *methods[])
{
    if (!strcmp(controller, "test"))
    {
        methods[0] = "GET";
    }
}

size_t LenControllerMethods(char *controller)
{
    if (!strcmp(controller, "test"))
        return 1;

    return 0;
}

void GetControllerHeaderNames(char *controller, char *names[])
{
}

char *GetControllerHeader(char *controller, char *name)
{
    return NULL;
}

size_t LenControllerHeaders(char *controller)
{
    return 0;
}

void GetControllerSchemas(char *controller, char *schemas[])
{
}

size_t LenControllerSchemas(char *controller)
{
    return 0;
}

void CallController(char *controller, char *method, char *proto, char *body, char *urlScheme, char *urlHost, char *urlPort, char *urlPath, char **queryKeys, size_t queryLen, char ***queryValues, size_t *queryValueLen, char **headerNames, char **headerValues, size_t headerLen, Response *resp)
{
    resp->status = 200;
    resp->headerLen = 0;

    char *btn = "</p><br /><a href=\"/static/singin\">Sing in</a></body></html>";

    for (int i = 0; i < queryLen; i++)
    {
        if (!strcmp(*(queryKeys + i), "name"))
        {
            resp->data = malloc(strlen("<!DOCTYPE html><html><body><p>Hello, ") + strlen(queryValues[i][0]) + strlen(btn) + 1);
            strcpy(resp->data, "<!DOCTYPE html><html><body><p>Hello, ");
            strcat(resp->data, queryValues[i][0]);
            strcat(resp->data, btn);
        }
    }
}