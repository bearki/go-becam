#define UNICODE
#include <initguid.h>
#include <winnls.h>
#include <dshow.h>
#include <mmsystem.h>
#include <qedit.h>
#include <unistd.h>
#include <setupapi.h>
#include <devpkey.h>
#include <iostream>
#include "dshow.hpp"
#include "_cgo_export.h"

// printErr shows string representation of HRESULT.
// This is for debugging.
void printErr(HRESULT hr)
{
    char buf[128];
    AMGetErrorTextA(hr, buf, 128);
    fprintf(stderr, "%s\n", buf);
}

// utf16Decode returns UTF-8 string from UTF-16 string.
std::string utf16Decode(LPOLESTR olestr)
{
    std::wstring wstr(olestr);
    const int len = WideCharToMultiByte(
        CP_UTF8, 0,
        wstr.data(), (int)wstr.size(),
        nullptr, 0, nullptr, nullptr);
    std::string str(len, 0);
    WideCharToMultiByte(
        CP_UTF8, 0,
        wstr.data(), (int)wstr.size(),
        (LPSTR)str.data(), len, nullptr, nullptr);
    return str;
}

// getPin is a helper to get I/O pin of DirectShow filter.
IPin *getPin(IBaseFilter *filter, PIN_DIRECTION dir)
{
    IEnumPins *enumPins;
    if (FAILED(filter->EnumPins(&enumPins)))
        return nullptr;

    IPin *pin;
    while (enumPins->Next(1, &pin, nullptr) == S_OK)
    {
        PIN_DIRECTION d;
        pin->QueryDirection(&d);
        if (d == dir)
        {
            enumPins->Release();
            return pin;
        }
        pin->Release();
    }
    enumPins->Release();
    return nullptr;
}

// getPath returns path of the device.
// returned pointer must be released by free() after use.
char *getPath(IMoniker *moniker)
{
    LPOLESTR path;
    if (FAILED(moniker->GetDisplayName(nullptr, nullptr, &path)))
        return nullptr;

    std::string nameStr = utf16Decode(path);
    char *ret = new char[nameStr.size() + 1];
    memcpy(ret, nameStr.c_str(), nameStr.size() + 1);

    LPMALLOC comalloc;
    CoGetMalloc(1, &comalloc);
    comalloc->Free(path);

    return ret;
}

// getLocationInfo returns device path of the device.
// returned pointer must be released by free() after use.
char *getLocationInfo(IMoniker *moniker)
{
    IPropertyBag *pPropBag;
    if (FAILED(moniker->BindToStorage(0, 0, IID_IPropertyBag, (void **)&pPropBag)))
    {
        return nullptr;
    }

    VARIANT locationInfo;
    VariantInit(&locationInfo);
    if (SUCCEEDED(pPropBag->Read(L"DevicePath", &locationInfo, 0)))
    {
        HDEVINFO hDevInfo = SetupDiCreateDeviceInfoList(NULL, NULL);
        if (INVALID_HANDLE_VALUE == hDevInfo)
        {
            return nullptr;
        }

        BYTE buf[1024] = {0};
        TCHAR szTemp[MAX_PATH] = {0};
        SP_DEVICE_INTERFACE_DATA spdid = {0};
        spdid.cbSize = sizeof(SP_DEVICE_INTERFACE_DATA);
        PSP_DEVICE_INTERFACE_DETAIL_DATA pspdidd = (PSP_DEVICE_INTERFACE_DETAIL_DATA)buf;
        pspdidd->cbSize = sizeof(*pspdidd);
        SP_DEVINFO_DATA spdd = {0};
        spdd.cbSize = sizeof(spdd);
        DWORD dwSize = 0;

        do
        {
            if (!SetupDiOpenDeviceInterface(hDevInfo, locationInfo.bstrVal, 0, &spdid))
            {
                break;
            }

            dwSize = sizeof(buf);
            if (!SetupDiGetDeviceInterfaceDetail(hDevInfo, &spdid, pspdidd, dwSize, &dwSize, &spdd))
            {
                break;
            }

            DEVPROPTYPE prop_type = 0;
            DWORD rsize = 0;

            SetupDiGetDevicePropertyW(hDevInfo, &spdd, &DEVPKEY_Device_LocationInfo, &prop_type, (PBYTE)szTemp, MAX_PATH - 1, &rsize, 0);
        } while (0);

        SetupDiDestroyDeviceInfoList(hDevInfo);

        std::string nameStr = utf16Decode(szTemp);
        char *ret = new char[nameStr.size() + 1];
        memcpy(ret, nameStr.c_str(), nameStr.size() + 1);

        VariantClear(&locationInfo);
        safeRelease(&pPropBag);

        return ret;
    }

    safeRelease(&pPropBag);
    return nullptr;
}

// getName returns name of the device.
// returned pointer must be released by free() after use.
char *getName(IMoniker *moniker)
{
    IPropertyBag *pPropBag;
    if (FAILED(moniker->BindToStorage(0, 0, IID_IPropertyBag, (void **)&pPropBag)))
    {
        return nullptr;
    }

    VARIANT varName;
    VariantInit(&varName);
    if (SUCCEEDED(pPropBag->Read(L"FriendlyName", &varName, 0)))
    {
        std::string nameStr = utf16Decode(varName.bstrVal);
        char *ret = new char[nameStr.size() + 1];
        memcpy(ret, nameStr.c_str(), nameStr.size() + 1);

        VariantClear(&varName);
        safeRelease(&pPropBag);
        return ret;
    }

    safeRelease(&pPropBag);
    return nullptr;
}

// listCamera stores information of the devices to cameraList*.
int listCamera(cameraList *list, const char **errstr)
{
    ICreateDevEnum *sysDevEnum = nullptr;
    IEnumMoniker *enumMon = nullptr;

    if (FAILED(CoCreateInstance(
            CLSID_SystemDeviceEnum, nullptr, CLSCTX_INPROC,
            IID_ICreateDevEnum, (void **)&sysDevEnum)))
    {
        *errstr = errEnumDevice;
        goto fail;
    }

    if (FAILED(sysDevEnum->CreateClassEnumerator(
            CLSID_VideoInputDeviceCategory, &enumMon, 0)))
    {
        *errstr = errEnumDevice;
        goto fail;
    }

    safeRelease(&sysDevEnum);

    // 检查是否有枚举到设备
    if (enumMon == nullptr)
    {
        return 0;
    }

    {
        IMoniker *moniker;
        list->num = 0;
        while (enumMon->Next(1, &moniker, nullptr) == S_OK)
        {
            moniker->Release();
            list->num++;
        }

        enumMon->Reset();
        list->path = new char *[list->num];
        list->locationInfo = new char *[list->num];
        list->name = new char *[list->num];

        int i = 0;
        while (enumMon->Next(1, &moniker, nullptr) == S_OK)
        {
            list->path[i] = getPath(moniker);
            list->locationInfo[i] = getLocationInfo(moniker);
            list->name[i] = getName(moniker);
            moniker->Release();
            i++;
        }
    }

    safeRelease(&enumMon);
    return 0;

fail:
    safeRelease(&sysDevEnum);
    safeRelease(&enumMon);
    return 1;
}

// freeCameraList frees all resources stored in cameraList*.
int freeCameraList(cameraList *list, const char **errstr)
{
    if (list->path != nullptr)
    {
        for (int i = 0; i < list->num; ++i)
        {
            if (list->path[i])
            {
                delete list->path[i];
            }
        }
        delete list->path;
    }
    if (list->locationInfo != nullptr)
    {
        for (int i = 0; i < list->num; ++i)
        {
            if (list->locationInfo[i])
            {
                delete list->locationInfo[i];
            }
        }
        delete list->locationInfo;
    }
    if (list->name != nullptr)
    {
        for (int i = 0; i < list->num; ++i)
        {
            if (list->name[i])
            {
                delete list->name[i];
            }
        }
        delete list->name;
    }
    return 1;
}

// selectCamera stores pointer to the selected device IMoniker* according to the configs in camera*.
int selectCamera(camera *cam, IMoniker **monikerSelected, const char **errstr)
{
    ICreateDevEnum *sysDevEnum = nullptr;
    IEnumMoniker *enumMon = nullptr;

    if (FAILED(CoCreateInstance(
            CLSID_SystemDeviceEnum, nullptr, CLSCTX_INPROC,
            IID_ICreateDevEnum, (void **)&sysDevEnum)))
    {
        *errstr = errEnumDevice;
        goto fail;
    }

    if (FAILED(sysDevEnum->CreateClassEnumerator(
            CLSID_VideoInputDeviceCategory, &enumMon, 0)))
    {
        *errstr = errEnumDevice;
        goto fail;
    }

    safeRelease(&sysDevEnum);

    {
        IMoniker *moniker;
        while (enumMon->Next(1, &moniker, nullptr) == S_OK)
        {
            char *path = getPath(moniker);
            if (strcmp(cam->path, path) != 0)
            {
                free(path);
                safeRelease(&moniker);
                continue;
            }
            free(path);
            *monikerSelected = moniker;
            safeRelease(&enumMon);
            return 1;
        }
    }

    safeRelease(&enumMon);
    return 0;

fail:
    safeRelease(&sysDevEnum);
    safeRelease(&enumMon);
    return 1;
}

// getResolution stores list of the device to camera*.
int getResolution(camera *cam, const char **errstr)
{
    cam->props = nullptr;

    IMoniker *moniker = nullptr;
    IBaseFilter *captureFilter = nullptr;
    ICaptureGraphBuilder2 *captureGraph = nullptr;
    IAMStreamConfig *config = nullptr;
    IPin *src = nullptr;

    if (!selectCamera(cam, &moniker, errstr))
    {
        goto fail;
    }

    moniker->BindToObject(0, 0, IID_IBaseFilter, (void **)&captureFilter);
    safeRelease(&moniker);

    src = getPin(captureFilter, PINDIR_OUTPUT);
    if (src == nullptr)
    {
        *errstr = errGetConfig;
        goto fail;
    }

    // Getting IAMStreamConfig is stub on Wine. Requires real Windows.
    if (FAILED(src->QueryInterface(
            IID_IAMStreamConfig, (void **)&config)))
    {
        *errstr = errGetConfig;
        goto fail;
    }
    safeRelease(&src);

    {
        int count = 0, size = 0;
        if (FAILED(config->GetNumberOfCapabilities(&count, &size)))
        {
            *errstr = errGetConfig;
            goto fail;
        }
        cam->props = new imageProp[count];

        int iProp = 0;
        for (int i = 0; i < count; ++i)
        {
            VIDEO_STREAM_CONFIG_CAPS caps;
            AM_MEDIA_TYPE *mediaType;
            if (FAILED(config->GetStreamCaps(i, &mediaType, (BYTE *)&caps)))
                continue;

            if (mediaType->majortype != MEDIATYPE_Video ||
                mediaType->formattype != FORMAT_VideoInfo ||
                mediaType->pbFormat == nullptr)
                continue;

            VIDEOINFOHEADER *videoInfoHdr = (VIDEOINFOHEADER *)mediaType->pbFormat;
            cam->props[iProp].width = videoInfoHdr->bmiHeader.biWidth;
            cam->props[iProp].height = videoInfoHdr->bmiHeader.biHeight;
            cam->props[iProp].fps = 10000000 / videoInfoHdr->AvgTimePerFrame;
            cam->props[iProp].fcc = videoInfoHdr->bmiHeader.biCompression;
            iProp++;
        }
        cam->numProps = iProp;
    }
    safeRelease(&config);
    safeRelease(&captureGraph);
    safeRelease(&captureFilter);
    safeRelease(&moniker);
    return 0;

fail:
    safeRelease(&src);
    safeRelease(&config);
    safeRelease(&captureGraph);
    safeRelease(&captureFilter);
    safeRelease(&moniker);
    return 1;
}

// freeResolution frees all resources stored in props*.
void freeResolution(camera *cam)
{
    if (cam->props)
    {
        delete cam->props;
        cam->props = nullptr;
    }
}

// SampleCB is not used in this app.
HRESULT SampleGrabberCallback::SampleCB(double sampleTime, IMediaSample *sample)
{
    return S_OK;
}

// BufferCB receives image from DirectShow.
HRESULT SampleGrabberCallback::BufferCB(double sampleTime, BYTE *buf, LONG len)
{
    imageCallback((char *)buf, (size_t)len);
    return S_OK;
}

// openCamera opens a camera and stores interface handler to camera*.
// camera* should be freed by freeCamera() after use.
int openCamera(camera *cam, const char **errstr)
{
    cam->grabber = nullptr;
    cam->mediaControl = nullptr;
    cam->callback = nullptr;

    IMoniker *moniker = nullptr;
    IGraphBuilder *graphBuilder = nullptr;
    IBaseFilter *captureFilter = nullptr;
    IMediaControl *mediaControl = nullptr;
    IBaseFilter *grabberFilter = nullptr;
    ISampleGrabber *grabber = nullptr;
    IBaseFilter *nullFilter = nullptr;
    IPin *src = nullptr;
    IPin *dst = nullptr;
    IPin *end = nullptr;
    IPin *nul = nullptr;

    if (!selectCamera(cam, &moniker, errstr))
    {
        goto fail;
    }
    moniker->BindToObject(0, 0, IID_IBaseFilter, (void **)&captureFilter);
    safeRelease(&moniker);

    if (FAILED(CoCreateInstance(
            CLSID_FilterGraph, nullptr, CLSCTX_INPROC,
            IID_IGraphBuilder, (void **)&graphBuilder)))
    {
        *errstr = errGraphBuilder;
        goto fail;
    }

    if (FAILED(graphBuilder->QueryInterface(
            IID_IMediaControl, (void **)&mediaControl)))
    {
        *errstr = errNoControl;
        goto fail;
    }

    if (FAILED(graphBuilder->AddFilter(captureFilter, L"capture")))
    {
        *errstr = errAddFilter;
        goto fail;
    }

    if (FAILED(CoCreateInstance(
            CLSID_SampleGrabber, nullptr, CLSCTX_INPROC,
            IID_IBaseFilter, (void **)&grabberFilter)))
    {
        *errstr = errGrabber;
        goto fail;
    }

    if (FAILED(grabberFilter->QueryInterface(IID_ISampleGrabber, (void **)&grabber)))
    {
        *errstr = errGrabber;
        goto fail;
    }

    {
        AM_MEDIA_TYPE mediaType;
        memset(&mediaType, 0, sizeof(mediaType));
        mediaType.majortype = MEDIATYPE_Video;
        mediaType.subtype = MEDIASUBTYPE_MJPG;
        mediaType.formattype = FORMAT_VideoInfo;
        mediaType.bFixedSizeSamples = 1;
        mediaType.cbFormat = sizeof(VIDEOINFOHEADER);

        VIDEOINFOHEADER videoInfoHdr;
        memset(&videoInfoHdr, 0, sizeof(VIDEOINFOHEADER));
        videoInfoHdr.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
        videoInfoHdr.bmiHeader.biWidth = cam->width;
        videoInfoHdr.bmiHeader.biHeight = cam->height;
        videoInfoHdr.AvgTimePerFrame = 10000000 / cam->fps;
        // videoInfoHdr.bmiHeader.biPlanes = 1;
        // videoInfoHdr.bmiHeader.biBitCount = 16;
        videoInfoHdr.bmiHeader.biCompression = MAKEFOURCC('M', 'J', 'P', 'G');
        mediaType.pbFormat = (BYTE *)&videoInfoHdr;
        if (FAILED(grabber->SetMediaType(&mediaType)))
        {
            *errstr = errGrabber;
            goto fail;
        }
    }

    if (FAILED(graphBuilder->AddFilter(grabberFilter, L"grabber")))
    {
        *errstr = errAddFilter;
        goto fail;
    }

    if (FAILED(CoCreateInstance(
            CLSID_NullRenderer, nullptr, CLSCTX_INPROC,
            IID_IBaseFilter, (void **)&nullFilter)))
    {
        *errstr = errTerminator;
        goto fail;
    }

    if (FAILED(graphBuilder->AddFilter(nullFilter, L"bull")))
    {
        *errstr = errAddFilter;
        goto fail;
    }

    HRESULT hr;
    src = getPin(captureFilter, PINDIR_OUTPUT);
    dst = getPin(grabberFilter, PINDIR_INPUT);
    if (src == nullptr || dst == nullptr ||
        FAILED(hr = graphBuilder->Connect(src, dst)))
    {
        *errstr = errConnectFilters;
        goto fail;
    }

    safeRelease(&src);
    safeRelease(&dst);

    end = getPin(grabberFilter, PINDIR_OUTPUT);
    nul = getPin(nullFilter, PINDIR_INPUT);
    if (end == nullptr || nul == nullptr ||
        FAILED(hr = graphBuilder->Connect(end, nul)))
    {
        *errstr = errConnectFilters;
        goto fail;
    }

    safeRelease(&end);
    safeRelease(&nul);

    safeRelease(&nullFilter);
    safeRelease(&captureFilter);
    safeRelease(&grabberFilter);
    safeRelease(&graphBuilder);

    {
        SampleGrabberCallback *cb = new SampleGrabberCallback(cam);
        grabber->SetCallback(cb, 1);
        cam->grabber = (void *)grabber;
        cam->mediaControl = (void *)mediaControl;
        cam->callback = (void *)cb;

        grabber->SetBufferSamples(true);
        mediaControl->Run();
    }

    return 0;

fail:
    safeRelease(&src);
    safeRelease(&dst);
    safeRelease(&end);
    safeRelease(&nul);
    safeRelease(&nullFilter);
    safeRelease(&grabber);
    safeRelease(&grabberFilter);
    safeRelease(&mediaControl);
    safeRelease(&captureFilter);
    safeRelease(&graphBuilder);
    safeRelease(&moniker);
    return 1;
}

// freeCamera closes device and frees all resources allocated by openCamera().
void freeCamera(camera *cam)
{
    if (cam->mediaControl)
        ((IMediaControl *)cam->mediaControl)->Stop();

    safeRelease((ISampleGrabber **)&cam->grabber);
    safeRelease((IMediaControl **)&cam->mediaControl);

    if (cam->callback)
    {
        ((SampleGrabberCallback *)cam->callback)->Release();
        delete ((SampleGrabberCallback *)cam->callback);
        cam->callback = nullptr;
    }

    freeResolution(cam);
}
