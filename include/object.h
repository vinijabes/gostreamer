#ifndef OBJECT_H
#define OBJECT_H

#include "pch.h"

void gostreamer_object_set_string(GstObject* object, const char* first_property_name, gchar* arg);
void gostreamer_object_set_int(GstObject* object, const char* first_property_name, gint arg);
void gostreamer_object_set_uint(GstObject* object, const char* first_property_name, guint arg);
void gostreamer_object_set_bool(GstObject* object, const char* first_property_name, gboolean arg);
void gostreamer_object_set_caps(GstObject* object, const char* first_property_name, const GstCaps *arg);

gchar*      gostreamer_object_get_string(GstObject* object, const char* first_property_name);
gint        gostreamer_object_get_int(GstObject* object, const char* first_property_name);
guint       gostreamer_object_get_uint(GstObject* object, const char* first_property_name);
gboolean    gostreamer_object_get_bool(GstObject* object, const char* first_property_name);

#endif