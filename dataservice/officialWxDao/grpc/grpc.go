package grpc

import (
	"context"
	"fmt"

	"github.com/Naist4869/awesomeProject/api"
	"github.com/Naist4869/awesomeProject/model/officialmodel"
	"google.golang.org/grpc"
)

type OfficialWxClient struct {
	FileSystemClient api.FileSystemClient
	TBKClient        api.TBKClient
}

func NewOfficialWxClient(fileSystemClient api.FileSystemClient, tbkClient api.TBKClient) *OfficialWxClient {
	return &OfficialWxClient{
		FileSystemClient: fileSystemClient,
		TBKClient:        tbkClient,
	}
}

func (oc OfficialWxClient) MediaIDGet(ctx context.Context, req officialmodel.MediaIDReq, args ...interface{}) (resp officialmodel.MediaIDResp, err error) {
	option := make([]grpc.CallOption, len(args))
	for i := range args {
		if o, ok := args[i].(grpc.CallOption); ok {
			option[i] = o
		}
	}
	get := new(api.MediaIDResp)
	get, err = oc.FileSystemClient.MediaIDGet(ctx, &api.MediaIDReq{
		FakeID:    req.FakeID,
		Timestamp: req.Timestamp,
	}, option...)
	if err != nil || get == nil {
		fmt.Println("---------------------------------------------------")
		err = fmt.Errorf("MediaIDGet%w", err)
		return
	}
	resp.MediaID = get.MediaID
	return
}
func (oc OfficialWxClient) TitleConvertTBKey(ctx context.Context, req officialmodel.TitleConvertTBKeyReq, args ...interface{}) (resp officialmodel.TitleConvertTBKeyResp, err error) {
	option := make([]grpc.CallOption, len(args))
	for i := range args {
		if o, ok := args[i].(grpc.CallOption); ok {
			option[i] = o
		}
	}
	get := new(api.TitleConvertTBKeyResp)
	get, err = oc.TBKClient.TitleConvertTBKey(ctx, &api.TitleConvertTBKeyReq{
		Title: req.Title,
	}, option...)
	if err != nil || get == nil {
		fmt.Println("---------------------------------------------------")
		err = fmt.Errorf("TitleConvertTBKey%w", err)
		return
	}
	resp.TBKey = get.TBKey
	return
}
func (oc OfficialWxClient) KeyConvertKey(ctx context.Context, req officialmodel.KeyConvertKeyReq, args ...interface{}) (resp officialmodel.KeyConvertKeyResp, err error) {
	option := make([]grpc.CallOption, len(args))
	for i := range args {
		if o, ok := args[i].(grpc.CallOption); ok {
			option[i] = o
		}
	}
	get := new(api.KeyConvertKeyResp)
	get, err = oc.TBKClient.KeyConvertKey(ctx, &api.KeyConvertKeyReq{
		FromKey: req.FromKey,
		UserID:  req.UserID,
	}, option...)
	if err != nil || get == nil {
		fmt.Println("---------------------------------------------------")
		err = fmt.Errorf("KeyConvertKey%w", err)
		return
	}
	resp.ToKey = get.ToKey
	resp.Price = get.Price
	resp.Rebate = get.Rebate
	resp.Title = get.Title
	resp.PicURL = get.PicURL
	return
}
