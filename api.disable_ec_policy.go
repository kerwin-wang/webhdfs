package webhdfs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	strings_ "github.com/searKing/golang/go/strings"

	"github.com/searKing/golang/go/errors"
)

type DisableECPolicyRequest struct {
	Authentication
	ProxyUser
	CSRF

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
	if req.Authentication.Delegation != nil {
		v.Set("delegation", aws.StringValue(req.Authentication.Delegation))
	}
	if req.ProxyUser.Username != nil {
		v.Set("user.name", aws.StringValue(req.ProxyUser.Username))
	}
	if req.ProxyUser.DoAs != nil {
		v.Set("doas", aws.StringValue(req.ProxyUser.DoAs))
	}

	if req.ECPolicy != nil {
		v.Set("ecpolicy", aws.StringValue(req.ECPolicy))
	}
	return v.Encode()
}

func (resp *DisableECPolicyResponse) UnmarshalHTTP(httpResp *http.Response) error {
	resp.HttpResponse.UnmarshalHTTP(httpResp)
	if isSuccessHttpCode(httpResp.StatusCode) {
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return nil
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return fmt.Errorf("parse %s: %w", strings_.Truncate(string(body), MaxHTTPBodyLengthDumped), err)
	}

	if err := resp.Exception(); err != nil {
		return err
	}
	return nil
}

// Disable EC Policy
// See: https://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/WebHDFS.html#Disable_EC_Policy
func (c *Client) DisableECPolicy(req *DisableECPolicyRequest) (*DisableECPolicyResponse, error) {
	return c.disableECPolicy(nil, req)
}
func (c *Client) DisableECPolicyWithContext(ctx context.Context, req *DisableECPolicyRequest) (*DisableECPolicyResponse, error) {
	if ctx == nil {
		panic("nil context")
	}
	return c.disableECPolicy(ctx, req)
}
func (c *Client) disableECPolicy(ctx context.Context, req *DisableECPolicyRequest) (*DisableECPolicyResponse, error) {
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

		httpReq, err := http.NewRequest(http.MethodPut, u.String(), nil)
		if err != nil {
			return nil, err
		}
		if req.CSRF.XXsrfHeader != nil {
			httpReq.Header.Set("X-XSRF-HEADER", aws.StringValue(req.CSRF.XXsrfHeader))
		}

		if ctx != nil {
			httpReq = httpReq.WithContext(ctx)
		}
		httpResp, err := c.httpClient.Do(httpReq)
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