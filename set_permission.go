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

type SetPermissionRequest struct {
	// Path of the object to get.
	//
	// Path is a required field
	Path *string `validate:"required"`

	// Name				owner
	// Description		The username who is the owner of a file/directory.
	// Type				String
	// Default Value	<empty> (means keeping it unchanged)
	// Valid Values		Any valid username.
	// Syntax			Any string.
	Owner *string

	// Name				group
	// Description		The name of a group.
	// Type				String
	// Default Value	<empty> (means keeping it unchanged)
	// Valid Values		Any valid group name.
	// Syntax			Any string.
	Group *string
}

type SetPermissionResponse struct {
	NameNode string `json:"-"`
	ErrorResponse
	HttpResponse `json:"-"`
}

func (req *SetPermissionRequest) RawPath() string {
	return aws.StringValue(req.Path)
}
func (req *SetPermissionRequest) RawQuery() string {
	v := url.Values{}
	v.Set("op", OpSetPermission)
	if req.Owner != nil {
		v.Set("owner", aws.StringValue(req.Owner))
	}
	if req.Group != nil {
		v.Set("group", aws.StringValue(req.Group))
	}
	return v.Encode()
}

func (resp *SetPermissionResponse) UnmarshalHTTP(httpResp *http.Response) error {
	resp.HttpResponse.UnmarshalHTTP(httpResp)
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

// Set Owner
// See: https://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/WebHDFS.html#Set_Owner
func (c *Client) SetPermission(req *SetPermissionRequest) (*SetPermissionResponse, error) {
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

		req, err := http.NewRequest(http.MethodPut, u.String(), nil)
		if err != nil {
			return nil, err
		}

		httpResp, err := c.httpClient.Do(req)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		var resp SetPermissionResponse
		resp.NameNode = addr

		if err := resp.UnmarshalHTTP(httpResp); err != nil {
			errs = append(errs, err)
			continue
		}

		return &resp, nil
	}
	return nil, errors.Multi(errs...)
}