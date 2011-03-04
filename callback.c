#include <JavaScriptCore/JSObjectRef.h>
#include <assert.h>
#include <stdio.h>
#include "_cgo_export.h"
#include "callback.h"

static JSValueRef nativecallback_CallAsFunction(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception)
{
	assert( exception );

	// Routine must set private to callback point in Go
	void* data = JSObjectGetPrivate( function );
	JSValueRef ret = nativecallback_CallAsFunction_go( data, (void*)ctx, (void*)function, (void*)thisObject, argumentCount, arguments, (void**)exception );
	assert( *exception==NULL || (*exception && !ret) );
	return ret;
}

static JSValueRef nativecallback_ObjectConvertToType(JSContextRef ctx, JSObjectRef object, JSType type, JSValueRef* exception)
{
	if ( type == kJSTypeString ) {
		JSStringRef str = JSStringCreateWithUTF8CString( "nativecallback" );
		JSValueRef ret = JSValueMakeString( ctx, str );
		JSStringRelease( str );
		return ret;
	}

	return 0;
}

JSClassRef JSClassDefinition_NativeCallback()
{
	static JSClassDefinition def = {
		0,
		kJSClassAttributeNone,
		"nativecallback",
		NULL,
        	NULL, // staticValues;
    		NULL, // staticFunctions;
		NULL, // initialize;
		NULL, // finalize;
		NULL, // hasProperty;
		NULL, // getProperty;
		NULL, // setProperty;
		NULL, // deleteProperty;
		NULL, // getPropertyNames;
		nativecallback_CallAsFunction, // callAsFunction;
		NULL, // callAsConstructor;
		NULL, // hasInstance;
		nativecallback_ObjectConvertToType // convertToType;
	};

	return JSClassCreate( &def );
}

