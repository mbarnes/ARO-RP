package vpc

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ListTrafficMirrorFilters invokes the vpc.ListTrafficMirrorFilters API synchronously
func (client *Client) ListTrafficMirrorFilters(request *ListTrafficMirrorFiltersRequest) (response *ListTrafficMirrorFiltersResponse, err error) {
	response = CreateListTrafficMirrorFiltersResponse()
	err = client.DoAction(request, response)
	return
}

// ListTrafficMirrorFiltersWithChan invokes the vpc.ListTrafficMirrorFilters API asynchronously
func (client *Client) ListTrafficMirrorFiltersWithChan(request *ListTrafficMirrorFiltersRequest) (<-chan *ListTrafficMirrorFiltersResponse, <-chan error) {
	responseChan := make(chan *ListTrafficMirrorFiltersResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListTrafficMirrorFilters(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ListTrafficMirrorFiltersWithCallback invokes the vpc.ListTrafficMirrorFilters API asynchronously
func (client *Client) ListTrafficMirrorFiltersWithCallback(request *ListTrafficMirrorFiltersRequest, callback func(response *ListTrafficMirrorFiltersResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListTrafficMirrorFiltersResponse
		var err error
		defer close(result)
		response, err = client.ListTrafficMirrorFilters(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ListTrafficMirrorFiltersRequest is the request struct for api ListTrafficMirrorFilters
type ListTrafficMirrorFiltersRequest struct {
	*requests.RpcRequest
	ResourceOwnerId         requests.Integer `position:"Query" name:"ResourceOwnerId"`
	TrafficMirrorFilterIds  *[]string        `position:"Query" name:"TrafficMirrorFilterIds"  type:"Repeated"`
	TrafficMirrorFilterName string           `position:"Query" name:"TrafficMirrorFilterName"`
	NextToken               string           `position:"Query" name:"NextToken"`
	ResourceOwnerAccount    string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount            string           `position:"Query" name:"OwnerAccount"`
	OwnerId                 requests.Integer `position:"Query" name:"OwnerId"`
	MaxResults              requests.Integer `position:"Query" name:"MaxResults"`
}

// ListTrafficMirrorFiltersResponse is the response struct for api ListTrafficMirrorFilters
type ListTrafficMirrorFiltersResponse struct {
	*responses.BaseResponse
	NextToken            string                `json:"NextToken" xml:"NextToken"`
	RequestId            string                `json:"RequestId" xml:"RequestId"`
	TotalCount           string                `json:"TotalCount" xml:"TotalCount"`
	TrafficMirrorFilters []TrafficMirrorFilter `json:"TrafficMirrorFilters" xml:"TrafficMirrorFilters"`
}

// CreateListTrafficMirrorFiltersRequest creates a request to invoke ListTrafficMirrorFilters API
func CreateListTrafficMirrorFiltersRequest() (request *ListTrafficMirrorFiltersRequest) {
	request = &ListTrafficMirrorFiltersRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ListTrafficMirrorFilters", "vpc", "openAPI")
	request.Method = requests.POST
	return
}

// CreateListTrafficMirrorFiltersResponse creates a response to parse from ListTrafficMirrorFilters response
func CreateListTrafficMirrorFiltersResponse() (response *ListTrafficMirrorFiltersResponse) {
	response = &ListTrafficMirrorFiltersResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}