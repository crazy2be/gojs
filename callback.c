#include <JavaScriptCore/JSObjectRef.h>
#include <assert.h>
#include <stdio.h>
#include "_cgo_export.h"
#include "callback.h"

static JSValueRef JSObjectCallAsFunctionCallback_trampoline(
	JSContextRef ctx, JSObjectRef function, 
	JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], 
	JSValueRef* exception)
{
	assert( ctx );
	assert( function );

	JSObjectCallAsFunctionCallback_go( ctx, function, thisObject, argumentCount, arguments, exception );
	return 0;
}

JSObjectRef JSObjectMakeFunctionWithCallback_wka( JSContextRef ctx, JSStringRef name )
{
	assert( ctx );
	assert( name );

	JSObjectRef ref = JSObjectMakeFunctionWithCallback( ctx, name, JSObjectCallAsFunctionCallback_trampoline );
	return ref;
}

