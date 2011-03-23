#include <JavaScriptCore/JSObjectRef.h>

struct nativeobject_data_tag { 
	void* typ;
	void* addr;
};
typedef struct nativeobject_data_tag nativeobject_data;

struct nativemethod_data_tag { 
	void* typ;
	void* addr;
	unsigned method;
};
typedef struct nativemethod_data_tag nativemethod_data;

void*	new_nativeobject_data( void*, void * );
void*	new_nativemethod_data( void*, void *, unsigned );

JSClassRef JSClassDefinition_NativeCallback();
JSClassRef JSClassDefinition_NativeFunction();
JSClassRef JSClassDefinition_NativeObject();
JSClassRef JSClassDefinition_NativeMethod();

