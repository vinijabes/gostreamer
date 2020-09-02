#ifndef ELEMENT_H
#define ELEMENT_H

#include "pch.h"

typedef struct ElementUserData {
    guint64 callbackID;
} ElementUserData;

GstFlowReturn gostreamer_element_push_buffer(GstElement *element, void *buffer,int len);

extern void go_pad_added_callback(GstElement* element, GstPad* pad, guint64 callbackID);
void gostreamer_pad_added_callback(GstElement* element, GstPad* pad, gpointer data);

extern void go_pad_removed_callback(GstElement* element, GstPad* pad, guint64 callbackID);
void gostreamer_pad_removed_callback(GstElement* element, GstPad* pad, gpointer data);

extern void go_new_sample_callback(GstElement* element, void *buffer, int bufferLen, int samples, guint64 callbackID);
GstFlowReturn gostreamer_new_sample_callback(GstElement *object, gpointer user_data);

gulong gostreamer_add_pad_added_signal(GstElement* element, guint64 callbackID);
gulong gostreamer_add_pad_removed_signal(GstElement* element, guint64 callbackID);
gulong gostreamer_add_sample_added_signal(GstElement* element, guint64 callbackID);

#endif