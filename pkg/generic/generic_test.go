/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generic

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
)

func TestBinaryThriftGeneric(t *testing.T) {
	g := BinaryThriftGeneric()
	defer g.Close()

	test.Assert(t, g.Framed() == true)
	test.Assert(t, g.PayloadCodec().Name() == "RawThriftBinary")
	test.Assert(t, g.PayloadCodecType() == serviceinfo.Thrift)

	method, err := g.GetMethod(nil, "Test")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "Test")
}

func TestMapThriftGeneric(t *testing.T) {
	p, err := NewThriftFileProvider("./map_test/idl/mock.thrift")
	test.Assert(t, err == nil)

	// new
	g, err := MapThriftGeneric(p)
	test.Assert(t, err == nil)
	defer g.Close()

	mg, ok := g.(*mapThriftGeneric)
	test.Assert(t, ok)

	test.Assert(t, g.PayloadCodec().Name() == "MapThrift")

	err = SetBinaryWithBase64(g, true)
	test.Assert(t, err == nil)
	test.Assert(t, mg.codec.binaryWithBase64)
	test.Assert(t, g.Framed() == false)

	err = SetBinaryWithByteSlice(g, true)
	test.Assert(t, err == nil)
	test.Assert(t, mg.codec.binaryWithByteSlice)

	test.Assert(t, g.PayloadCodecType() == serviceinfo.Thrift)

	method, err := g.GetMethod(nil, "Test")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "Test")
}

func TestMapThriftGenericForJSON(t *testing.T) {
	p, err := NewThriftFileProvider("./map_test/idl/mock.thrift")
	test.Assert(t, err == nil)

	// new
	g, err := MapThriftGenericForJSON(p)
	test.Assert(t, err == nil)
	defer g.Close()

	mg, ok := g.(*mapThriftGeneric)
	test.Assert(t, ok)

	test.Assert(t, g.PayloadCodec().Name() == "MapThrift")

	err = SetBinaryWithBase64(g, true)
	test.Assert(t, err == nil)
	test.Assert(t, mg.codec.binaryWithBase64)
	test.Assert(t, g.Framed() == false)

	test.Assert(t, g.PayloadCodecType() == serviceinfo.Thrift)

	method, err := g.GetMethod(nil, "Test")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "Test")
}

func TestHTTPThriftGeneric(t *testing.T) {
	p, err := NewThriftFileProvider("./http_test/idl/binary_echo.thrift")
	test.Assert(t, err == nil)

	var opts []Option
	opts = append(opts, UseRawBodyForHTTPResp(true))
	g, err := HTTPThriftGeneric(p, opts...)
	test.Assert(t, err == nil)
	defer g.Close()

	hg, ok := g.(*httpThriftGeneric)
	test.Assert(t, ok)

	test.Assert(t, g.PayloadCodec().Name() == "HttpThrift")

	test.Assert(t, !hg.codec.dynamicgoEnabled)
	test.Assert(t, hg.codec.useRawBodyForHTTPResp)
	test.Assert(t, !hg.codec.binaryWithBase64)
	test.Assert(t, !hg.codec.convOpts.NoBase64Binary)
	test.Assert(t, !hg.codec.convOptsWithThriftBase.NoBase64Binary)

	err = SetBinaryWithBase64(g, true)
	test.Assert(t, err == nil)
	test.Assert(t, hg.codec.binaryWithBase64)
	test.Assert(t, !hg.codec.convOpts.NoBase64Binary)
	test.Assert(t, !hg.codec.convOptsWithThriftBase.NoBase64Binary)

	test.Assert(t, g.Framed() == false)

	test.Assert(t, g.PayloadCodecType() == serviceinfo.Thrift)

	url := "https://example.com/BinaryEcho"
	body := map[string]interface{}{
		"msg": "Test",
	}
	data, err := customJson.Marshal(body)
	test.Assert(t, err == nil)
	// NewRequest
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(data))
	test.Assert(t, err == nil)
	customReq, err := FromHTTPRequest(req)
	test.Assert(t, err == nil)

	method, err := g.GetMethod(customReq, "Test")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "BinaryEcho")
}

func TestHTTPThriftGenericWithDynamicGo(t *testing.T) {
	p, err := NewThriftFileProviderWithDynamicGo("./http_test/idl/binary_echo.thrift")
	test.Assert(t, err == nil)

	g, err := HTTPThriftGeneric(p)
	test.Assert(t, err == nil)
	defer g.Close()

	hg, ok := g.(*httpThriftGeneric)
	test.Assert(t, ok)

	test.Assert(t, g.PayloadCodec().Name() == "HttpThrift")

	test.Assert(t, hg.codec.dynamicgoEnabled)
	test.Assert(t, !hg.codec.useRawBodyForHTTPResp)
	test.Assert(t, !hg.codec.binaryWithBase64)
	test.Assert(t, hg.codec.convOpts.NoBase64Binary)
	test.Assert(t, hg.codec.convOptsWithThriftBase.NoBase64Binary)

	err = SetBinaryWithBase64(g, true)
	test.Assert(t, err == nil)
	test.Assert(t, hg.codec.binaryWithBase64)
	test.Assert(t, !hg.codec.convOpts.NoBase64Binary)
	test.Assert(t, !hg.codec.convOptsWithThriftBase.NoBase64Binary)

	test.Assert(t, g.Framed() == true)

	test.Assert(t, g.PayloadCodecType() == serviceinfo.Thrift)

	url := "https://example.com/BinaryEcho"
	body := map[string]interface{}{
		"msg": "Test",
	}
	data, err := customJson.Marshal(body)
	test.Assert(t, err == nil)
	// NewRequest
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(data))
	test.Assert(t, err == nil)
	customReq, err := FromHTTPRequest(req)
	test.Assert(t, err == nil)

	method, err := g.GetMethod(customReq, "Test")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "BinaryEcho")
}

func TestJSONThriftGeneric(t *testing.T) {
	p, err := NewThriftFileProvider("./json_test/idl/mock.thrift")
	test.Assert(t, err == nil)

	g, err := JSONThriftGeneric(p)
	test.Assert(t, err == nil)
	defer g.Close()

	jg, ok := g.(*jsonThriftGeneric)
	test.Assert(t, ok)

	test.Assert(t, g.PayloadCodec().Name() == "JSONThrift")

	test.Assert(t, !jg.codec.dynamicgoEnabled)
	test.Assert(t, jg.codec.binaryWithBase64)
	test.Assert(t, !jg.codec.convOpts.NoBase64Binary)
	test.Assert(t, !jg.codec.convOptsWithThriftBase.NoBase64Binary)

	err = SetBinaryWithBase64(g, false)
	test.Assert(t, err == nil)
	test.Assert(t, !jg.codec.binaryWithBase64)
	test.Assert(t, !jg.codec.convOpts.NoBase64Binary)
	test.Assert(t, !jg.codec.convOptsWithThriftBase.NoBase64Binary)

	test.Assert(t, g.Framed() == false)

	test.Assert(t, g.PayloadCodecType() == serviceinfo.Thrift)

	method, err := g.GetMethod(nil, "Test")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "Test")
}

func TestJSONThriftGenericWithDynamicGo(t *testing.T) {
	p, err := NewThriftFileProviderWithDynamicGo("./json_test/idl/mock.thrift")
	test.Assert(t, err == nil)

	g, err := JSONThriftGeneric(p)
	test.Assert(t, err == nil)
	defer g.Close()

	jg, ok := g.(*jsonThriftGeneric)
	test.Assert(t, ok)

	test.Assert(t, g.PayloadCodec().Name() == "JSONThrift")

	test.Assert(t, jg.codec.dynamicgoEnabled)
	test.Assert(t, jg.codec.binaryWithBase64)
	test.Assert(t, !jg.codec.convOpts.NoBase64Binary)
	test.Assert(t, !jg.codec.convOptsWithThriftBase.NoBase64Binary)

	err = SetBinaryWithBase64(g, false)
	test.Assert(t, err == nil)
	test.Assert(t, !jg.codec.binaryWithBase64)
	test.Assert(t, jg.codec.convOpts.NoBase64Binary)
	test.Assert(t, jg.codec.convOptsWithThriftBase.NoBase64Binary)

	test.Assert(t, g.Framed() == true)

	test.Assert(t, g.PayloadCodecType() == serviceinfo.Thrift)

	method, err := g.GetMethod(nil, "Test")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "Test")
}

func TestJSONPbGeneric(t *testing.T) {
	path := "./jsonpb_test/idl/echo.proto"
	content := `syntax = "proto3";
	package pbapi;
	// The greeting service definition.
	option go_package = "pbapi";
	
	message Request {
	  string message = 1;
	}
	
	message Response {
	  string message = 1;
	}
	
	service Echo {
	  rpc Echo (Request) returns (Response) {}
	}`

	includes := map[string]string{
		path: content,
	}

	// get pb generic
	p, err := NewPbContentProvider(path, includes)
	test.Assert(t, err == nil)

	g, err := JSONPbGeneric(p)
	test.Assert(t, err == nil)
	defer g.Close()

	test.Assert(t, g.PayloadCodec().Name() == "JSONPb")

	// err = SetBinaryWithBase64(g, true)
	// test.Assert(t, err == nil)
	test.Assert(t, g.Framed() == false)

	test.Assert(t, g.PayloadCodecType() == serviceinfo.Protobuf)

	method, err := g.GetMethod(nil, "Echo")
	test.Assert(t, err == nil)
	test.Assert(t, method.Name == "Echo")
}
