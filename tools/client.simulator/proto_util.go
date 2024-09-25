package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	_ "unsafe"

	"tevat.nd.org/basecode/goost/encoding/binary"
	"tevat.nd.org/toolchain/simulator/codec/model"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	pb "tevat.nd.org/framework/proxy/proto"
	pbparserPb "tevat.nd.org/pkg/pbparser/proto"
)

var (
	methodMap map[string]protoreflect.MethodDescriptor
	notifyMap map[string]protoreflect.FieldDescriptor
)

func (c *Client) GetMd(method string) (protoreflect.MethodDescriptor, bool) {
	md, ok := methodMap[method]
	return md, ok
}

func (c *Client) GetField(notify string) protoreflect.FieldDescriptor {
	return notifyMap[notify]
}

var initMapLock sync.Mutex

func (c *Client) initMethodMap() {
	methodMap = make(map[string]protoreflect.MethodDescriptor)
	notifyMap = make(map[string]protoreflect.FieldDescriptor)

	files := c.config.ProtoFileConfig.Files
	for _, v := range files {
		err := c.GetMdFromFile(v)
		if err != nil {
			fmt.Println("initMethodMap GetMdFromFile err:", v)
		}
	}
}

func (c *Client) GetFieldFromMd(md protoreflect.MethodDescriptor) {
	if !md.IsStreamingClient() && !md.IsStreamingServer() {
		return
	}
	out := md.Output()
	if out.Oneofs().Len() > 0 {
		for i := 0; i < out.Oneofs().Len(); i++ {
			fields := out.Oneofs().Get(i).Fields()
			for j := 0; j < fields.Len(); j++ {
				fd := fields.Get(j)
				opt, ok := fd.Options().(*descriptorpb.FieldOptions)
				var notifyName string
				if ok && opt != nil {
					notifyName = proto.GetExtension(opt, pbparserPb.E_Method).(string)
				} else {
					notifyName = string(fd.FullName())
				}
				fmt.Println("::notifyName", notifyName)
				notifyMap[notifyName] = fd
			}
		}
	} else {
		fields := out.Fields()
		for j := 0; j < fields.Len(); j++ {
			fd := fields.Get(j)
			opt, ok := fd.Options().(*descriptorpb.FieldOptions)
			var notifyName string
			if ok && opt != nil {
				notifyName = proto.GetExtension(opt, pbparserPb.E_Method).(string)
			} else {
				notifyName = string(fd.FullName())
			}
			fmt.Println("::notifyName", notifyName)
			notifyMap[notifyName] = fd
		}
	}
}

func (c *Client) GetMdFromFile(filename string) error {
	initMapLock.Lock()
	defer initMapLock.Unlock()
	protoPath := c.config.ProtoFileConfig.Path
	registry, err := createProtoRegistry(protoPath, filename)
	if err != nil {
		return err
	}

	// var md protoreflect.MethodDescriptor
	registry.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		sds := fd.Services()
		for i := 0; i < sds.Len(); i++ {
			sd := sds.Get(i)
			sdname := sd.FullName()

			mds := sd.Methods()
			for j := 0; j < mds.Len(); j++ {
				mdi := mds.Get(j)
				c.GetFieldFromMd(mdi)
				methodName := fmt.Sprintf("/%s/%s", sdname, mdi.Name())
				fmt.Println("::methodName", methodName)
				methodMap[methodName] = mdi
			}
		}
		return true
	})
	return nil
}

func CreateMd(m map[string]interface{}) binary.BytesWithUint16Len {
	metadata := make(map[string]*pb.Metadata_Value)
	for k, v := range m {
		metadata[k] = &pb.Metadata_Value{
			Values: []string{fmt.Sprintf("%v", v)},
		}
	}
	md := &pb.Metadata{
		Metadata: metadata,
	}
	b, _ := proto.Marshal(md)
	return b
}

func createProtoRegistry(srcDir, filename string) (*protoregistry.Files, error) {
	// Create descriptors using the protoc binary.
	// Imported dependencies are included so that the descriptors are self-contained.
	fns := strings.Split(filename, "/")
	tmpFile := fns[len(fns)-1] + "-tmp.pb"
	cmd := exec.Command("protoc",
		"--include_imports",
		"--descriptor_set_out="+tmpFile,
		"-I",
		srcDir,
		path.Join(srcDir, filename))
	// cmd = exec.Command("protoc", "--include_imports --descriptor_set_out=auth.proto-tmp.pb -I proto auth.proto")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("cmd err:", err)
		return nil, err
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile)

	marshalledDescriptorSet, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, err
	}
	descriptorSet := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(marshalledDescriptorSet, &descriptorSet)
	if err != nil {
		return nil, err
	}

	files, err := protodesc.NewFiles(&descriptorSet)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (c *Client) ParseRes(res model.Response) Res {
	r := Res{}
	r.ID = res.ID
	r.Result = res.Result
	r.Method = fmt.Sprintf("%s", res.Method)

	if r.Method == "" { // 升级 0.6 后好像这边 res.Method 没数据了
		r.Method = c.reqPack.GetMethod(r.ID)
	}

	var resp *dynamicpb.Message
	if !r.Result {
		// notify 不会有 Result = false 的情况
		errMehod := "/lobbyplatform.logic.LogicServer/TestGetErr"
		md, _ := c.GetMd(errMehod)
		resp = dynamicpb.NewMessage(md.Output())

		err := error(nil)
		if c.config.Proto != ProtoTypePROTOBUF {
			err = protojson.Unmarshal(res.Msg, resp)
		} else {
			err = proto.Unmarshal(res.Msg, resp)
		}
		if err != nil {
		} else {
			msg := make(map[string]interface{})
			marshaler := protojson.MarshalOptions{EmitUnpopulated: true}
			raw, _ := marshaler.Marshal(resp)
			jsonDec := json.NewDecoder(bytes.NewReader(raw))
			jsonDec.UseNumber()
			_ = jsonDec.Decode(&msg)
			r.Msg = msg
		}
		r.ErrMsg = string(res.Msg)

		return r
	} else {
		if r.ID == 0 {
			field := c.GetField(r.Method)
			resp = dynamicpb.NewMessage(field.Message())
		} else {
			md, _ := c.GetMd(r.Method)
			resp = dynamicpb.NewMessage(md.Output())
		}
		err := error(nil)
		if c.config.Proto != ProtoTypePROTOBUF {
			err = protojson.Unmarshal(res.Msg, resp)
		} else {
			err = proto.Unmarshal(res.Msg, resp)
		}

		if err != nil {
			r.Msg = map[string]interface{}{"msg": res.Msg}
		} else {
			msg := make(map[string]interface{})
			marshaler := protojson.MarshalOptions{EmitUnpopulated: true}
			raw, _ := marshaler.Marshal(resp)
			jsonDec := json.NewDecoder(bytes.NewReader(raw))
			jsonDec.UseNumber()
			_ = jsonDec.Decode(&msg)
			r.Msg = msg
		}

		return r
	}
}

func (c *Client) MarshalMsg(method string, msgMap map[string]interface{}) ([]byte, error) {
	if c.config.Proto == ProtoTypePROTOBUF {
		md, ok := c.GetMd(method)
		if !ok {
			fmt.Printf("apiData.json 文件配置错误， proto文件中没有找到对应的method:%s\n", method)
			return nil, fmt.Errorf("method (%v) no find", method)
		}
		req := dynamicpb.NewMessage(md.Input())

		msgBytes, err := json.Marshal(msgMap)
		if err != nil {
			return nil, err
		}

		if err := protojson.Unmarshal(msgBytes, req); err != nil {
			fmt.Printf("apiData.json 文件配置错误， %s方法中的req message与proto文件中的message不匹配，err:=%s \n", method, err)
			return nil, fmt.Errorf("unable to serialise request object to protocol buffer: %w", err)
		}
		b, err := proto.Marshal(req)
		if err != nil {
			return nil, err
		}
		// fmt.Println("reqmarshal", req, b)
		return b, nil
	}
	return json.Marshal(msgMap)
}
