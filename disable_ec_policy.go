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

type DisableECPolicyRequest struct {
	// Name				ecpolicy, Erasure Coding Policy
	// Description		The name of the erasure coding policy.
	// Type				String
	// Default Value	<empty>
	// Valid Values		Any valid erasure coding policy name;
	// Syntax			Any string.
	ECPolicy *string `validate:"required"`
}

type DisableECPolicyResponse struct {
	NameNode string `json:"-"`
	ErrorResponse
	HttpResponse `json:"-"`
}

func (req *DisableECPolicyRequest) RawPath() string {
	return ""
}
func (req *DisableECPolicyRequest) RawQuery() string {
	v := url.Values{}
	v.Set("op", OpDisableECPolicy)
	if req.ECPolicy != nil {
		v.Set("ecpolicy", aws.StringValue(req.ECPolicy))
	}
	return v.Encode()
}

func (resp *DisableECPolicyResponse) UnmarshalHTTP(httpResp *http.Response) error {
	resp.HttpResponse.UnmarshalHTTP(httpResp)
	defer resp.Body.Close()
	if isSuccessHttpCode(httpResp.StatusCode) {
		return nil
	}

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

// Disable EC Policy
// See: https://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/WebHDFS.html#Disable_EC_Policy
func (c *Client) DisableECPolicy(req *DisableECPolicyRequest) (*DisableECPolicyResponse, error) {
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

		var resp DisableECPolicyResponse
		resp.NameNode = addr

		if err := resp.UnmarshalHTTP(httpResp); err != nil {
			errs = append(errs, err)
			continue
		}

		return &resp, nil
	}
	return nil, errors.Multi(errs...)
}
