package webhdfs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/searKing/golang/go/errors"
)

type OpenRequest struct {
	// Path of the object to get.
	//
	// Path is a required field
	Path *string `validate:"required"`

	// Name				offset
	// Description		The starting byte position.
	// Type				long
	// Default Value	0
	// Valid Values		>= 0
	// Syntax			Any integer.
	Offset *int64
	// Name				length
	// Description		The number of bytes to be processed.
	// Type				long
	// Default Value	null (means the entire file)
	// Valid Values		>= 0 or null
	// Syntax			Any integer.
	Length *int64
	// Name				buffersize
	// Description		The size of the buffer used in transferring data.
	// Type				int
	// Default Value	Specified in the configuration.
	// Valid Values		> 0
	// Syntax			Any integer.
	BufferSize *int32
	// Name				nodirect
	// Description		Disable automatically redirected.
	// Type				bool
	// Default Value	false
	// Valid Values		true|false
	// Syntax			Any Bool.
	NoDirect *bool
}

type OpenResponse struct {
	NameNode string `json:"-"`
	ErrorResponse
	HttpResponse `json:"-"`
}

func (req *OpenRequest) RawPath() string {
	return aws.StringValue(req.Path)
}
func (req *OpenRequest) RawQuery() string {
	v := url.Values{}
	v.Set("op", OpOpen)
	if req.Offset != nil {
		v.Set("offset", fmt.Sprintf("%d", aws.Int64Value(req.Offset)))
	}
	if req.Length != nil {
		v.Set("length", fmt.Sprintf("%d", aws.Int64Value(req.Length)))
	}
	if req.BufferSize != nil {
		v.Set("buffersize", fmt.Sprintf("%d", aws.Int32Value(req.BufferSize)))
	}
	if req.NoDirect != nil {
		v.Set("noredirect", fmt.Sprintf("%t", aws.BoolValue(req.NoDirect)))
	}
	return v.Encode()
}

func (resp *OpenResponse) UnmarshalHTTP(httpResp *http.Response) error {
	resp.HttpResponse.UnmarshalHTTP(httpResp)

	if isSuccessHttpCode(httpResp.StatusCode) {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if err := resp.Exception(); err != nil {
		return err
	}
	return nil
}

// Open and Read a File
// See: https://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/WebHDFS.html#Open_and_Read_a_File
func (c *Client) Open(req *OpenRequest) (*OpenResponse, error) {
	err := c.opts.Validator.Struct(req)
	if err != nil {
		return nil, err
	}

	nameNodes := c.opts.Addresses
	if nameNodes == nil {
		return nil, fmt.Errorf("missing namenode addresses")
	}
	var u = c.HttpUrl(req)

	var errs []error
	for _, addr := range nameNodes {
		u.Host = addr
		httpResp, err := c.httpClient.Get(u.String())
		if err != nil {
			errs = append(errs, err)
			continue
		}

		var resp OpenResponse
		resp.NameNode = addr

		if err := resp.UnmarshalHTTP(httpResp); err != nil {
			errs = append(errs, err)
			continue
		}

		return &resp, nil
	}
	return nil, errors.Multi(errs...)
}