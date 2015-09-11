package uhost

import (
	"github.com/xiaohui/goucloud/ucloud"
)

// CreateUHostInstance will create instances
type CreateUHostInstanceParams struct {
	ucloud.CommonRequest

	Region    string
	ImageId   string
	LoginMode string
	Password  string
	KeyPair   string
	CPU       int
	Memory    int
	DiskSpace int
	Name      string
	NetworkId string

	SecurityGroupId string
	ChargeType      string
	Quantity        int
	Count           int
	UHostType       string
	NetCapability   string
	Tag             string
	CouponId        string
}

type CreateUHostInstanceResponse struct {
	ucloud.CommonResponse
	HostIds []string
}

func (u *UHost) CreateUHostInstance(params *CreateUHostInstanceParams) (*CreateUHostInstanceResponse, error) {
	response := CreateUHostInstanceResponse{}
	err := u.DoRequest("CreateUHostInstance", params, response)

	return &response, err
}

type DescribeImageParams struct {
	Region    string
	ImageType string
	OsType    string
	ImageId   string
	Offset    int
	Limit     int
}

type ImageSet struct {
	ImageId   string
	ImageName string
	OsType    string
	OsName    string
	State     string

	ImageDescription string
	CreateTime       string
}

type ImageSetArray []ImageSet

type DescribeImageResponse struct {
	ucloud.CommonResponse

	TotalCount int
	ImageSet   ImageSetArray
}

func (u *UHost) DescribeImage(params *DescribeImageParams) (*DescribeImageResponse, error) {
	response := DescribeImageResponse{}
	err := u.DoRequest("DescribeImage", params, response)

	return &response, err
}
