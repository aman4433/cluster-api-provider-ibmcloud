// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package tlog

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/sigstore/rekor/pkg/generated/models"
)

// GetLogInfoReader is a Reader for the GetLogInfo structure.
type GetLogInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetLogInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetLogInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetLogInfoDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetLogInfoOK creates a GetLogInfoOK with default headers values
func NewGetLogInfoOK() *GetLogInfoOK {
	return &GetLogInfoOK{}
}

/*
GetLogInfoOK describes a response with status code 200, with default header values.

A JSON object with the root hash and tree size as properties
*/
type GetLogInfoOK struct {
	Payload *models.LogInfo
}

// IsSuccess returns true when this get log info o k response has a 2xx status code
func (o *GetLogInfoOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get log info o k response has a 3xx status code
func (o *GetLogInfoOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get log info o k response has a 4xx status code
func (o *GetLogInfoOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get log info o k response has a 5xx status code
func (o *GetLogInfoOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get log info o k response a status code equal to that given
func (o *GetLogInfoOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get log info o k response
func (o *GetLogInfoOK) Code() int {
	return 200
}

func (o *GetLogInfoOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/log][%d] getLogInfoOK  %+v", 200, o.Payload)
}

func (o *GetLogInfoOK) String() string {
	return fmt.Sprintf("[GET /api/v1/log][%d] getLogInfoOK  %+v", 200, o.Payload)
}

func (o *GetLogInfoOK) GetPayload() *models.LogInfo {
	return o.Payload
}

func (o *GetLogInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.LogInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetLogInfoDefault creates a GetLogInfoDefault with default headers values
func NewGetLogInfoDefault(code int) *GetLogInfoDefault {
	return &GetLogInfoDefault{
		_statusCode: code,
	}
}

/*
GetLogInfoDefault describes a response with status code -1, with default header values.

There was an internal error in the server while processing the request
*/
type GetLogInfoDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get log info default response has a 2xx status code
func (o *GetLogInfoDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get log info default response has a 3xx status code
func (o *GetLogInfoDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get log info default response has a 4xx status code
func (o *GetLogInfoDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get log info default response has a 5xx status code
func (o *GetLogInfoDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get log info default response a status code equal to that given
func (o *GetLogInfoDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get log info default response
func (o *GetLogInfoDefault) Code() int {
	return o._statusCode
}

func (o *GetLogInfoDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1/log][%d] getLogInfo default  %+v", o._statusCode, o.Payload)
}

func (o *GetLogInfoDefault) String() string {
	return fmt.Sprintf("[GET /api/v1/log][%d] getLogInfo default  %+v", o._statusCode, o.Payload)
}

func (o *GetLogInfoDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetLogInfoDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}