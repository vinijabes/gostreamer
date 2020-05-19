#include "message.h"

GstMessageType gostreamer_get_message_type(GstMessage *message){
    return GST_MESSAGE_TYPE(message);
}