// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: proto/order.proto

package order

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OrderItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     int64                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity      int32                  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderItem) Reset() {
	*x = OrderItem{}
	mi := &file_proto_order_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderItem) ProtoMessage() {}

func (x *OrderItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderItem.ProtoReflect.Descriptor instead.
func (*OrderItem) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{0}
}

func (x *OrderItem) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *OrderItem) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type OrderRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Items         []*OrderItem           `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderRequest) Reset() {
	*x = OrderRequest{}
	mi := &file_proto_order_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderRequest) ProtoMessage() {}

func (x *OrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderRequest.ProtoReflect.Descriptor instead.
func (*OrderRequest) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{1}
}

func (x *OrderRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *OrderRequest) GetItems() []*OrderItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type OrderResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Status        string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Items         []*OrderItem           `protobuf:"bytes,5,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderResponse) Reset() {
	*x = OrderResponse{}
	mi := &file_proto_order_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderResponse) ProtoMessage() {}

func (x *OrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderResponse.ProtoReflect.Descriptor instead.
func (*OrderResponse) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{2}
}

func (x *OrderResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *OrderResponse) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *OrderResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *OrderResponse) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *OrderResponse) GetItems() []*OrderItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type OrderID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderID) Reset() {
	*x = OrderID{}
	mi := &file_proto_order_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderID) ProtoMessage() {}

func (x *OrderID) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderID.ProtoReflect.Descriptor instead.
func (*OrderID) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{3}
}

func (x *OrderID) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type OrderList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Orders        []*OrderResponse       `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderList) Reset() {
	*x = OrderList{}
	mi := &file_proto_order_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderList) ProtoMessage() {}

func (x *OrderList) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderList.ProtoReflect.Descriptor instead.
func (*OrderList) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{4}
}

func (x *OrderList) GetOrders() []*OrderResponse {
	if x != nil {
		return x.Orders
	}
	return nil
}

type StatusUpdate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StatusUpdate) Reset() {
	*x = StatusUpdate{}
	mi := &file_proto_order_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatusUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusUpdate) ProtoMessage() {}

func (x *StatusUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusUpdate.ProtoReflect.Descriptor instead.
func (*StatusUpdate) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{5}
}

func (x *StatusUpdate) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *StatusUpdate) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type UserOrdersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserOrdersRequest) Reset() {
	*x = UserOrdersRequest{}
	mi := &file_proto_order_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserOrdersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserOrdersRequest) ProtoMessage() {}

func (x *UserOrdersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_order_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserOrdersRequest.ProtoReflect.Descriptor instead.
func (*UserOrdersRequest) Descriptor() ([]byte, []int) {
	return file_proto_order_proto_rawDescGZIP(), []int{6}
}

func (x *UserOrdersRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

var File_proto_order_proto protoreflect.FileDescriptor

const file_proto_order_proto_rawDesc = "" +
	"\n" +
	"\x11proto/order.proto\x12\x05order\"F\n" +
	"\tOrderItem\x12\x1d\n" +
	"\n" +
	"product_id\x18\x01 \x01(\x03R\tproductId\x12\x1a\n" +
	"\bquantity\x18\x02 \x01(\x05R\bquantity\"O\n" +
	"\fOrderRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId\x12&\n" +
	"\x05items\x18\x02 \x03(\v2\x10.order.OrderItemR\x05items\"\x97\x01\n" +
	"\rOrderResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\x12\x16\n" +
	"\x06status\x18\x03 \x01(\tR\x06status\x12\x1d\n" +
	"\n" +
	"created_at\x18\x04 \x01(\tR\tcreatedAt\x12&\n" +
	"\x05items\x18\x05 \x03(\v2\x10.order.OrderItemR\x05items\"\x19\n" +
	"\aOrderID\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\"9\n" +
	"\tOrderList\x12,\n" +
	"\x06orders\x18\x01 \x03(\v2\x14.order.OrderResponseR\x06orders\"6\n" +
	"\fStatusUpdate\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\",\n" +
	"\x11UserOrdersRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\x03R\x06userId2\xf4\x01\n" +
	"\fOrderService\x128\n" +
	"\vCreateOrder\x12\x13.order.OrderRequest\x1a\x14.order.OrderResponse\x120\n" +
	"\bGetOrder\x12\x0e.order.OrderID\x1a\x14.order.OrderResponse\x12>\n" +
	"\x11UpdateOrderStatus\x12\x13.order.StatusUpdate\x1a\x14.order.OrderResponse\x128\n" +
	"\n" +
	"ListOrders\x12\x18.order.UserOrdersRequest\x1a\x10.order.OrderListB\x1eZ\x1corder-service/pb/order;orderb\x06proto3"

var (
	file_proto_order_proto_rawDescOnce sync.Once
	file_proto_order_proto_rawDescData []byte
)

func file_proto_order_proto_rawDescGZIP() []byte {
	file_proto_order_proto_rawDescOnce.Do(func() {
		file_proto_order_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_order_proto_rawDesc), len(file_proto_order_proto_rawDesc)))
	})
	return file_proto_order_proto_rawDescData
}

var file_proto_order_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_order_proto_goTypes = []any{
	(*OrderItem)(nil),         // 0: order.OrderItem
	(*OrderRequest)(nil),      // 1: order.OrderRequest
	(*OrderResponse)(nil),     // 2: order.OrderResponse
	(*OrderID)(nil),           // 3: order.OrderID
	(*OrderList)(nil),         // 4: order.OrderList
	(*StatusUpdate)(nil),      // 5: order.StatusUpdate
	(*UserOrdersRequest)(nil), // 6: order.UserOrdersRequest
}
var file_proto_order_proto_depIdxs = []int32{
	0, // 0: order.OrderRequest.items:type_name -> order.OrderItem
	0, // 1: order.OrderResponse.items:type_name -> order.OrderItem
	2, // 2: order.OrderList.orders:type_name -> order.OrderResponse
	1, // 3: order.OrderService.CreateOrder:input_type -> order.OrderRequest
	3, // 4: order.OrderService.GetOrder:input_type -> order.OrderID
	5, // 5: order.OrderService.UpdateOrderStatus:input_type -> order.StatusUpdate
	6, // 6: order.OrderService.ListOrders:input_type -> order.UserOrdersRequest
	2, // 7: order.OrderService.CreateOrder:output_type -> order.OrderResponse
	2, // 8: order.OrderService.GetOrder:output_type -> order.OrderResponse
	2, // 9: order.OrderService.UpdateOrderStatus:output_type -> order.OrderResponse
	4, // 10: order.OrderService.ListOrders:output_type -> order.OrderList
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_order_proto_init() }
func file_proto_order_proto_init() {
	if File_proto_order_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_order_proto_rawDesc), len(file_proto_order_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_order_proto_goTypes,
		DependencyIndexes: file_proto_order_proto_depIdxs,
		MessageInfos:      file_proto_order_proto_msgTypes,
	}.Build()
	File_proto_order_proto = out.File
	file_proto_order_proto_goTypes = nil
	file_proto_order_proto_depIdxs = nil
}
