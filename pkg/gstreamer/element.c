#include "element.h"

GstFlowReturn gostreamer_element_push_buffer(GstElement *element, void *buffer, int len) {
    gpointer p = g_memdup(buffer, len);
    GstBuffer *data = gst_buffer_new_wrapped(p, len);
    
    return gst_app_src_push_buffer(GST_APP_SRC(element), data);
}

void gostreamer_pad_added_callback(GstElement* element, GstPad* pad, gpointer data){
    ElementUserData* d = (ElementUserData*)data;
    go_pad_added_callback(element, pad, d->callbackID);
}

void gostreamer_pad_removed_callback(GstElement* element, GstPad* pad, gpointer data){
    ElementUserData* d = (ElementUserData*)data;
    go_pad_removed_callback(element, pad, d->callbackID);
}

GstFlowReturn gostreamer_new_sample_callback(GstElement *element, gpointer data) {
    GstSample *sample = NULL;
    GstBuffer *buffer = NULL;
    gpointer copy = NULL;
    gsize copy_size = 0;
    ElementUserData* d = (ElementUserData*)data;

    g_signal_emit_by_name (element, "pull-sample", &sample);
    if (sample) {
        buffer = gst_sample_get_buffer(sample);
        if (buffer) {
            gst_buffer_extract_dup(buffer, 0, gst_buffer_get_size(buffer), &copy, &copy_size);
            go_new_sample_callback(element, copy, copy_size, GST_BUFFER_DURATION(buffer), d->callbackID);
        }
        gst_sample_unref (sample);
    }

    return GST_FLOW_OK;
}

gulong gostreamer_add_pad_added_signal(GstElement* element, guint64 callbackID){
    ElementUserData *data = calloc(1, sizeof(ElementUserData));
    data->callbackID = callbackID;

    return g_signal_connect(element, "pad-added", G_CALLBACK(gostreamer_pad_added_callback), data);
}

gulong gostreamer_add_pad_removed_signal(GstElement* element, guint64 callbackID){
    ElementUserData *data = calloc(1, sizeof(ElementUserData));
    data->callbackID = callbackID;

    return g_signal_connect(element, "pad-removed", G_CALLBACK(gostreamer_pad_removed_callback), data);
}

gulong gostreamer_add_sample_added_signal(GstElement* element, guint64 callbackID){
    ElementUserData *data = calloc(1, sizeof(ElementUserData));
    data->callbackID = callbackID;

    return g_signal_connect(element, "new-sample", G_CALLBACK(gostreamer_new_sample_callback), data);
}