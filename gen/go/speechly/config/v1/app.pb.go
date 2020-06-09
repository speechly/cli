// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: speechly/config/v1/app.proto

package configv1

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type App_Status int32

const (
	App_STATUS_UNSPECIFIED App_Status = 0
	App_STATUS_NEW         App_Status = 1
	App_STATUS_TRAINING    App_Status = 2
	App_STATUS_TRAINED     App_Status = 3
	App_STATUS_FAILED      App_Status = 4
)

// Enum value maps for App_Status.
var (
	App_Status_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		1: "STATUS_NEW",
		2: "STATUS_TRAINING",
		3: "STATUS_TRAINED",
		4: "STATUS_FAILED",
	}
	App_Status_value = map[string]int32{
		"STATUS_UNSPECIFIED": 0,
		"STATUS_NEW":         1,
		"STATUS_TRAINING":    2,
		"STATUS_TRAINED":     3,
		"STATUS_FAILED":      4,
	}
)

func (x App_Status) Enum() *App_Status {
	p := new(App_Status)
	*p = x
	return p
}

func (x App_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (App_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_speechly_config_v1_app_proto_enumTypes[0].Descriptor()
}

func (App_Status) Type() protoreflect.EnumType {
	return &file_speechly_config_v1_app_proto_enumTypes[0]
}

func (x App_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use App_Status.Descriptor instead.
func (App_Status) EnumDescriptor() ([]byte, []int) {
	return file_speechly_config_v1_app_proto_rawDescGZIP(), []int{0, 0}
}

// App fields.
type App struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// IETF language tag
	Language string     `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
	Status   App_Status `protobuf:"varint,3,opt,name=status,proto3,enum=speechly.config.v1.App_Status" json:"status,omitempty"`
	// Name is for humans, if empty, id should be shown.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// Queue size when STATUS_NEW.
	QueueSize int32 `protobuf:"varint,5,opt,name=queue_size,json=queueSize,proto3" json:"queue_size,omitempty"`
	// Error when STATUS_FAILED, for more visibility.
	ErrorMsg string `protobuf:"bytes,6,opt,name=error_msg,json=errorMsg,proto3" json:"error_msg,omitempty"`
	// Estimated training time remaining. If 0 then no estimate is available.
	EstimatedRemainingSec int32 `protobuf:"varint,7,opt,name=estimated_remaining_sec,json=estimatedRemainingSec,proto3" json:"estimated_remaining_sec,omitempty"`
}

func (x *App) Reset() {
	*x = App{}
	if protoimpl.UnsafeEnabled {
		mi := &file_speechly_config_v1_app_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *App) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*App) ProtoMessage() {}

func (x *App) ProtoReflect() protoreflect.Message {
	mi := &file_speechly_config_v1_app_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use App.ProtoReflect.Descriptor instead.
func (*App) Descriptor() ([]byte, []int) {
	return file_speechly_config_v1_app_proto_rawDescGZIP(), []int{0}
}

func (x *App) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *App) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *App) GetStatus() App_Status {
	if x != nil {
		return x.Status
	}
	return App_STATUS_UNSPECIFIED
}

func (x *App) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *App) GetQueueSize() int32 {
	if x != nil {
		return x.QueueSize
	}
	return 0
}

func (x *App) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

func (x *App) GetEstimatedRemainingSec() int32 {
	if x != nil {
		return x.EstimatedRemainingSec
	}
	return 0
}

var File_speechly_config_v1_app_proto protoreflect.FileDescriptor

var file_speechly_config_v1_app_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x73, 0x70, 0x65, 0x65, 0x63, 0x68, 0x6c, 0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12,
	0x73, 0x70, 0x65, 0x65, 0x63, 0x68, 0x6c, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x22, 0xdf, 0x02, 0x0a, 0x03, 0x41, 0x70, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x73, 0x70, 0x65, 0x65, 0x63, 0x68, 0x6c,
	0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x36,
	0x0a, 0x17, 0x65, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x6d, 0x61,
	0x69, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x65, 0x63, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x15, 0x65, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x64, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e,
	0x69, 0x6e, 0x67, 0x53, 0x65, 0x63, 0x22, 0x6c, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x4e, 0x45, 0x57, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x54, 0x52, 0x41, 0x49, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x12, 0x0a,
	0x0e, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x54, 0x52, 0x41, 0x49, 0x4e, 0x45, 0x44, 0x10,
	0x03, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x04, 0x42, 0x5e, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x70, 0x65, 0x65,
	0x63, 0x68, 0x6c, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x08,
	0x41, 0x70, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x08, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x53, 0x43, 0x58, 0xaa, 0x02, 0x12, 0x53, 0x70, 0x65,
	0x65, 0x63, 0x68, 0x6c, 0x79, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x12, 0x53, 0x70, 0x65, 0x65, 0x63, 0x68, 0x6c, 0x79, 0x5c, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x5c, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_speechly_config_v1_app_proto_rawDescOnce sync.Once
	file_speechly_config_v1_app_proto_rawDescData = file_speechly_config_v1_app_proto_rawDesc
)

func file_speechly_config_v1_app_proto_rawDescGZIP() []byte {
	file_speechly_config_v1_app_proto_rawDescOnce.Do(func() {
		file_speechly_config_v1_app_proto_rawDescData = protoimpl.X.CompressGZIP(file_speechly_config_v1_app_proto_rawDescData)
	})
	return file_speechly_config_v1_app_proto_rawDescData
}

var file_speechly_config_v1_app_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_speechly_config_v1_app_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_speechly_config_v1_app_proto_goTypes = []interface{}{
	(App_Status)(0), // 0: speechly.config.v1.App.Status
	(*App)(nil),     // 1: speechly.config.v1.App
}
var file_speechly_config_v1_app_proto_depIdxs = []int32{
	0, // 0: speechly.config.v1.App.status:type_name -> speechly.config.v1.App.Status
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_speechly_config_v1_app_proto_init() }
func file_speechly_config_v1_app_proto_init() {
	if File_speechly_config_v1_app_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_speechly_config_v1_app_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*App); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_speechly_config_v1_app_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_speechly_config_v1_app_proto_goTypes,
		DependencyIndexes: file_speechly_config_v1_app_proto_depIdxs,
		EnumInfos:         file_speechly_config_v1_app_proto_enumTypes,
		MessageInfos:      file_speechly_config_v1_app_proto_msgTypes,
	}.Build()
	File_speechly_config_v1_app_proto = out.File
	file_speechly_config_v1_app_proto_rawDesc = nil
	file_speechly_config_v1_app_proto_goTypes = nil
	file_speechly_config_v1_app_proto_depIdxs = nil
}
