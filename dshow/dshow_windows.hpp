#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C"
{
#endif

    typedef struct
    {
        int width;
        int height;
        int fps;
        uint32_t fcc;
    } imageProp;

    typedef struct
    {
        int width;
        int height;
        int fps;
        size_t buf; // uintptr

        char *path;
        int numProps;
        imageProp *props;

        void *grabber;
        void *mediaControl;
        void *callback;
    } camera;

    typedef struct
    {
        int num;
        char **path;
        char **locationInfo;
        char **name;
    } cameraList;

    static imageProp *getProp(camera *cam, int i) { return &cam->props[i]; }
    static char *getPath(cameraList *list, int i) { return list->path[i]; }
    static char *getLocationInfo(cameraList *list, int i) { return list->locationInfo[i]; }
    static char *getName(cameraList *list, int i) { return list->name[i]; }

    int listCamera(cameraList *list, const char **errstr);
    int freeCameraList(cameraList *list, const char **errstr);
    int getResolution(camera *cam, const char **errstr);
    void freeResolution(camera *cam);
    int openCamera(camera *cam, const char **errstr);
    void freeCamera(camera *cam);

#ifdef __cplusplus
}
#endif

#ifdef __cplusplus
#include <dshow.h>
#include <windows.h>

#include <string>
std::string utf16Decode(LPOLESTR olestr);
IPin *getPin(IBaseFilter *filter, PIN_DIRECTION dir);
char *getPath(IMoniker *monier);
char *getName(IMoniker *monier);

template <class T>
void safeRelease(T **p)
{
    if (*p)
    {
        (*p)->Release();
        *p = nullptr;
    }
}

class SampleGrabberCallback : public ISampleGrabberCB
{
private:
    camera *cam_;

public:
    inline SampleGrabberCallback(camera *cam) : cam_(cam) {}

    HRESULT SampleCB(double sampleTime, IMediaSample *sample) final;
    HRESULT BufferCB(double sampleTime, BYTE *buffer, LONG bufferLen) final;

    inline ULONG STDMETHODCALLTYPE AddRef() final
    {
        return 2;
    }
    inline ULONG STDMETHODCALLTYPE Release() final
    {
        return 1;
    }
    inline HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, void **ppv) final
    {
        if (riid == IID_ISampleGrabberCB || riid == IID_IUnknown)
        {
            *ppv = (void *)static_cast<ISampleGrabberCB *>(this);
            return NOERROR;
        }
        return E_NOINTERFACE;
    }
};

EXTERN_C const CLSID CLSID_NullRenderer;
EXTERN_C const CLSID CLSID_SampleGrabber;

const static char *errAddFilter = "failed to add filter";
const static char *errConnectFilters = "failed to connect filters";
const static char *errEnumDevice = "failed to enum device";
const static char *errGrabber = "failed to create grabber";
const static char *errGraphBuilder = "failed to build graph";
const static char *errNoControl = "failed to control media";
const static char *errTerminator = "failed to create graph terminator";
const static char *errGetConfig = "failed to get config";

#endif