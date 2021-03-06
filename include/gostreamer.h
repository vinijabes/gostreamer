#ifndef GOSTREAMER_H
#define GOSTREAMER_H

#include "pch.h"

extern void goPrint(char *str);
extern void goHandleBusMessage(GstMessage* msg, int pipelineID);
extern void goHandlePadAddedSignal(int elementId, GstPad* pad);

static gint
toGstMessageType(void *p) {
	return (GST_MESSAGE_TYPE(p));
}

static const char*
messageTypeName(void *p)
{
	return (GST_MESSAGE_TYPE_NAME(p));
}

static guint64
messageTimestamp(void *p)
{
	return (GST_MESSAGE_TIMESTAMP(p));
}

void gstreamer_init();

GstElement* gstreamer_element_factory_make(const char* factoryname,const char* name);
int gstreamer_element_link(GstElement* src, GstElement* dest);
void gstreamer_object_set(GstElement* object, const char* first_property_name, const char* arg);
void gstreamer_object_set_int(GstElement* object, const char* first_property_name, long arg);
void gstreamer_object_set_double(GstElement* object, const char* first_property_name, double arg);
void gstreamer_object_set_bool(GstElement* object, const char* first_property_name, gboolean arg);
void gstreamer_element_signal_connect(GstElement* element, const char* signal, int elementId);
void gstreamer_element_pad_added_signal_connect(GstElement* element, int elementId);
void gstreamer_set_caps(GstElement *element,const char *capstr); 
GstPadTemplate* gstreamer_get_pad_template(GstElement* element, const char* padName);
GstPad* gstreamer_element_request_pad(GstElement* element, GstPadTemplate* template);
void gstreamer_element_push_buffer(GstElement *element, void *buffer,int len);

GstPipeline* gstreamer_create_pipeline(const char* name); 
void gstreamer_pipeline_start(GstPipeline* pipeline);
void gstreamer_pipeline_pause(GstPipeline *pipeline);
void gstreamer_pipeline_stop(GstPipeline *pipeline);
void gstreamer_pipeline_sendeos(GstPipeline *pipeline);
void gstreamer_bin_add_element(GstPipeline *pipeline, GstElement* element);
void gstreamer_pipeline_bus_watch(GstPipeline* pipeline, int pipelineId);

gboolean gstreamer_pad_link(GstPad* src, GstPad* dest);
void gstreamer_pad_object_set(GstPad* object, const char* first_property_name, const char* arg);
void gstreamer_pad_object_set_int(GstPad* object, const char* first_property_name, long arg);
void gstreamer_pad_object_set_double(GstPad* object, const char* first_property_name, double arg);
void gstreamer_pad_object_set_bool(GstPad* object, const char* first_property_name, gboolean arg);


#endif