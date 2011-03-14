#include <JavaScriptCore/JSObjectRef.h>
#include <assert.h>
#include <stdlib.h>
#include "_cgo_export.h"
#include "callback.h"

static JSValueRef nativecallback_CallAsFunction(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef* exception)
{
	assert( exception );

	// Routine must set private to callback point in Go
	void* data = JSObjectGetPrivate( function );
	JSValueRef ret = nativecallback_CallAsFunction_go( data, (void*)ctx, (void*)function, (void*)thisObject, argumentCount, (void*)arguments, (void**)exception );
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

void*	new_nativeobject_data( void* typ, void* addr )
{
	nativeobject_data* ptr = (nativeobject_data*)malloc( sizeof(nativeobject_data) );
	if ( ptr ) {
		ptr->typ = typ;
		ptr->addr = addr;
	}
	return ptr;
}

static void nativeobject_Finalize(JSObjectRef object)
{
	void* data = JSObjectGetPrivate( object );
	free( data );
}

static JSValueRef nativeobject_GetProperty(JSContextRef ctx, JSObjectRef object, JSStringRef propertyName, JSValueRef* exception)
{
	assert( exception );

	// Routine must set private to callback point in Go
	void* data = JSObjectGetPrivate( object );
	JSValueRef ret = nativeobject_GetProperty_go( data, (void*)ctx, (void*)object, (void*)propertyName, (void**)exception );
	assert( *exception==NULL || (*exception && !ret) );
	return ret;
}

JSClassRef JSClassDefinition_NativeObject()
{
	static JSClassDefinition def = {
		0,
		kJSClassAttributeNone,
		"nativeobject",
		NULL,
        	NULL, // staticValues;
    		NULL, // staticFunctions;
		NULL, // initialize;
		nativeobject_Finalize, // finalize;
		NULL, // hasProperty;
		nativeobject_GetProperty, // getProperty;
		NULL, // setProperty;
		NULL, // deleteProperty;
		NULL, // getPropertyNames;
		NULL, // callAsFunction;
		NULL, // callAsConstructor;
		NULL, // hasInstance;
		NULL // convertToType;
	};

	return JSClassCreate( &def );
}

