#ifndef ELEMENT_H
#define ELEMENT_H

#include "pch.h"

GstFlowReturn gostreamer_element_push_buffer(GstElement *element, void *buffer,int len);

#endif