#include <JavaScriptCore/JSObjectRef.h>

struct nativeobject_data_tag { 
	void* typ;
	void* addr;
};
typedef struct nativeobject_data_tag nativeobject_data;

void*	new_nativeobject_data( void*, void * );

JSClassRef JSClassDefinition_NativeCallback();
JSClassRef JSClassDefinition_NativeFunction();
JSClassRef JSClassDefinition_NativeObject();

