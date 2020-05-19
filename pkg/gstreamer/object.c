#include "object.h"

void gostreamer_object_set_string(GstObject* object,const char* first_property_name, gchar* arg){
    g_object_set(object, first_property_name, arg, NULL);
}

void gostreamer_object_set_int(GstObject* object,const char* first_property_name, gint arg){
    gint64 garg = arg;
    g_object_set(object, first_property_name, garg, NULL);
}

void gostreamer_object_set_uint(GstObject* object,const char* first_property_name, guint arg){
    g_object_set(object, first_property_name, arg, NULL);
}

void gostreamer_object_set_bool(GstObject* object, const char* first_property_name, gboolean arg){
    g_object_set(object, first_property_name, arg, NULL);
}

void gostreamer_object_set_double(GstObject* object, const char* first_property_name, gdouble arg){
    g_object_set(object, first_property_name, arg, NULL);
}

void gostreamer_object_set_caps(GstObject* object, const char* first_property_name,const GstCaps *arg){
    g_object_set(object, first_property_name, arg, NULL);
}

gchar* gostreamer_object_get_string(GstObject* object, const char* first_property_name) {
    gchar *data;
    g_object_get (object, first_property_name, &data, NULL);
    return data;
}

gint gostreamer_object_get_int(GstObject* object, const char* first_property_name) {
    gint data;
    g_object_get (object, first_property_name, &data, NULL);
    return data;
}

guint gostreamer_object_get_uint(GstObject* object, const char* first_property_name){
    guint data;
    g_object_get (object, first_property_name, &data, NULL);
    return data;
}

gboolean gostreamer_object_get_bool(GstObject* object, const char* first_property_name){
    gboolean data;
    g_object_get (object, first_property_name, &data, NULL);
    return data;
}

gdouble gostreamer_object_get_double(GstObject* object, const char* first_property_name){
    gdouble data;
    g_object_get (object, first_property_name, &data, NULL);
    return data;
}
