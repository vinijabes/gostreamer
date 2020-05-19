#include "element.h"

GstFlowReturn gostreamer_element_push_buffer(GstElement *element, void *buffer, int len) {
    gpointer p = g_memdup(buffer, len);
    GstBuffer *data = gst_buffer_new_wrapped(p, len);
    
    return gst_app_src_push_buffer(GST_APP_SRC(element), data);
}