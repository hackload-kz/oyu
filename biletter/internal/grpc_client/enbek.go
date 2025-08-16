package grpc_client

//
//import (
//	"biletter/internal/domain/entity"
//	pb "biletter/internal/pb/router"
//	"context"
//)

//func (g *GrpcClient) SendSoapRequestByOperation(ctx context.Context, cachedData entity.SignXMLEnbek, xmlData string) (
//	err error,
//) {
//
//	rpcRequest := &pb.SignedXmlAndContextRequest{
//		Context: cachedData.RpcContext,
//		SignedXmlData: &pb.SignedXmlData{
//			Xml:       xmlData,
//			Operation: cachedData.Operation,
//			DocId:     int32(cachedData.ID),
//			Type:      cachedData.Type,
//			Code:      cachedData.Code,
//		},
//	}
//
//	_, err = g.gateway.SendSoapRequestByOperation(ctx, rpcRequest)
//	if err != nil {
//		g.logger.Error(err)
//		return err
//	}
//
//	return nil
//}
