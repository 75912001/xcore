package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"tevat.nd.org/basecode/goost/async"
	"tevat.nd.org/basecode/goost/errors"
	"tevat.nd.org/toolchain/simulator/codec"
	"tevat.nd.org/toolchain/simulator/codec/keep"
	"tevat.nd.org/toolchain/simulator/codec/model"
	"tevat.nd.org/toolchain/simulator/codec/standard"
)

const (
	ProtoTypePROTOBUF ProtoType = 1
	ProtoTypeJSON     ProtoType = 2
)

const (
	CodecTypeStandard CodecType = 1
	CodecTypeKeep     CodecType = 2
)

type (
	ProtoType int32

	CodecType uint32

	ReqPack struct {
		id        uint32
		id2Method map[uint32]string
		lock      sync.RWMutex
	}

	Client struct {
		conn    net.Conn
		data    map[string]ApiData
		config  Config
		codec   codec.Coder
		onRec   func(res Res)
		res     Res
		cond    sync.Cond
		reqPack *ReqPack
	}

	Config struct {
		Addr               string
		ErrorProcessMethod string
		Codec              CodecType
		Proto              ProtoType
		ProtoFileConfig    ProtoFileConfig
		IgnoreMethod       []string

		ignoreMethodMap map[string]struct{}
	}

	ProtoFileConfig struct {
		Path  string
		Files []string
	}

	ApiData struct {
		ID       uint32
		Method   string
		Msg      map[string]interface{}
		Commands []string
	}

	Res struct {
		ID     uint32
		Result bool
		Method string
		Msg    map[string]interface{}
		ErrMsg string
	}
)

func (r *ReqPack) GetId() uint32 {
	return r.id
}

func (r *ReqPack) GetNewId(method string) uint32 {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.id += 1
	r.id2Method[r.id] = method
	return r.id
}

func (r *ReqPack) GetMethod(id uint32) string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if method, ok := r.id2Method[id]; ok {
		return method
	}
	return ""
}

func (c *Client) Connect() error {
	fmt.Println("start:", c.config.Addr)
	conn, err := net.Dial("tcp", c.config.Addr)
	if err != nil {
		fmt.Printf("conn Connect fail: %s \n", err.Error())
		fmt.Printf("Reconnect ...")
		time.Sleep(time.Second * 3)
		return c.Connect()
	}

	c.conn = conn

	if c.config.Codec == CodecTypeKeep {
		c.Send(
			model.MsgSyn{
				Version: keep.NetVersion,
				Token:   0,
				LastID:  0,
			})
	}

	return nil
}

func initFile(d any, filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("open file fail:%+v \n", err)
		panic(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			return
		}
	}(jsonFile)
	jsonDec := json.NewDecoder(jsonFile)
	jsonDec.UseNumber()
	err = jsonDec.Decode(d)
	if err != nil {
		fmt.Printf("Unmarshal file fail:%+v \n", err)
		panic(err)
	}
}

func (c *Client) Init(onRec func(res Res)) {
	if c.data != nil {
		return
	}

	c.reqPack = &ReqPack{
		id:        0,
		id2Method: make(map[uint32]string),
		lock:      sync.RWMutex{},
	}

	c.cond = *sync.NewCond(&sync.Mutex{})
	c.onRec = onRec
	initFile(&c.data, "bin/config/apiData.json")
	initFile(&c.config, "bin/config/config.json")

	if c.config.Codec == 0 {
		c.config.Codec = CodecTypeStandard
	}
	if c.config.Codec == CodecTypeKeep {
		c.codec = keep.NewCodec()
	} else {
		c.codec = standard.NewCodec()
	}

	c.config.ignoreMethodMap = make(map[string]struct{})
	for _, v := range c.config.IgnoreMethod {
		c.config.ignoreMethodMap[v] = struct{}{}
	}
	c.initMethodMap()
}

func (c *Client) ReloadApi() {
	c.data = nil
	initFile(&c.data, "bin/config/apiData.json")
	c.initMethodMap()
}

func (c *Client) GetApiDataByName(name string) ApiData {
	apiData := c.data[name]
	//if apiData.Metadata == nil {
	//	apiData.Metadata = make(map[string]interface{})
	//}
	if apiData.Msg == nil {
		apiData.Msg = make(map[string]interface{})
	}
	return apiData
}

func (c *Client) GetRequestFromJson(name string) (*model.Request, error) {
	// log.Printf("apiDataS:%+v \n", apiDataS)
	reqJson := c.GetApiDataByName(name)
	msg, _ := json.Marshal(reqJson.Msg)
	req := &model.Request{
		Method: []byte(reqJson.Method),
		Msg:    msg,
	}
	return req, nil
}

func (c *Client) ToString(data any) string {
	return fmt.Sprintf("%s", data)
}

func (c *Client) GetReqObject(name string, options ...func(map[string]interface{})) *model.Request {
	var err error
	req, err := c.GetRequestFromJson(name)
	if err != nil {
		log.Printf("GetRequestFromJson fail, err:%+v \n", errors.WithStack(err))
		return nil
	}
	msg := map[string]interface{}{}
	jsonDec := json.NewDecoder(bytes.NewReader(req.Msg))
	jsonDec.UseNumber()
	err = jsonDec.Decode(&msg)
	// fmt.Println("==============", req.Msg, msg, err, name)
	if err != nil {
		log.Printf("GetReqObject Unmarshal fail, err:%+v \n", errors.WithStack(err))
		return nil
	}
	for _, o := range options {
		o(msg)
	}

	if req.Method == nil || len(req.Method) == 0 {
		fmt.Printf("apiData.json 文件配置错误， 文件中没有找到指令:%s\n", name)
		log.Printf("apiData no finded:%+v \n", name)
		return nil
	}

	method := c.ToString(req.Method)
	msgBytes, err := c.MarshalMsg(method, msg)
	if err != nil {
		log.Printf("GetReqObject Marshal fail, method:%+v, msg:%+v, err:%+v \n", method, msg, errors.WithStack(err))
		return nil
	}
	req.Msg = msgBytes
	return req
}

func (c *Client) Send(reqAny any) (*model.Request, error) {
	var err error
	conn := c.conn
	if conn == nil {
		return nil, fmt.Errorf("conn is nil")
	}
	req := &model.Request{}
	switch v := reqAny.(type) {
	case model.MsgSyn:
		return req, c.codec.Encode(conn, v)
	case *model.Request:
		req = v
	case map[string]interface{}:
		id, ok := v["id"].(int64)
		if ok {
			req.ID = uint32(id)
		}
		method := v["method"].(string)
		req.Method = []byte(method)
		msg := v["msg"].(map[string]interface{})
		msgBytes, err := c.MarshalMsg(method, msg)
		if err != nil {
			log.Printf("send fail by json.Marshal(msg) msg:%+v \n", msg)
		}
		req.Msg = msgBytes
		if v["metadata"] != nil {
			metadata := v["metadata"].(map[string]interface{})
			req.Metadata = CreateMd(metadata)
		}
	default:
		err = errors.WithStack(fmt.Errorf("send fail by invalid req:%+v", reqAny))
		fmt.Println(err)
		return nil, err
	}

	// log.Printf("[%v]::req, req.ID:%+v, method:%v, msg:%+v \n", time.Now(), req.ID, m.ToString(req.Method), m.Parse(req.Msg))
	// log.Printf("%+v", req.ID)
	if req.Method == nil || len(req.Method) == 0 {
		return nil, errors.WithStack(fmt.Errorf("req is invalid, req: %+v \n", req))
	}
	return req, c.codec.Encode(conn, req)
}

func (c *Client) SendReq(name string) bool {
	req := c.GetReqObject(name)
	if req == nil {
		return false
	}
	id := c.reqPack.GetNewId(string(req.Method))
	req.ID = id
	_, err := c.Send(req)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return true
}

func (c *Client) WaitRes() Res {
	for {
		c.cond.L.Lock()
		c.cond.Wait()
		if c.res.ID == c.reqPack.GetId() {
			res := c.res
			c.cond.L.Unlock()
			return res
		} else {
			c.cond.L.Unlock()
			continue
		}
	}
}

func (c *Client) Decode(r io.Reader) (Res, error) {
	nilRes := Res{}
	if c.codec == nil {
		return nilRes, nil
	}
	res, err := c.codec.Decode(r)
	if err != nil {
		log.Printf("=====Rec err=====%+v \n", err)
		return nilRes, err
	}
	if v, ok := res.(model.Response); ok {
		return c.ParseRes(v), nil
	}
	if v, ok := res.(model.MsgSynAck); ok {
		return Res{Result: true, Method: "MsgSyn", Msg: map[string]interface{}{"msg": v}}, nil
	}
	return nilRes, nil
	//var res Res
	//var h uint32
	////fmt.Println("read h")
	//if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
	//	return res, err
	//}
	//
	////fmt.Println("read res", h)
	//resBase := model.Response{}
	//err := binary.Read(r, binary.LittleEndian, &resBase)
	////method := fmt.Sprintf("%s", resBase.Method)
	////fmt.Println("decode:, err:", resBase, err, method)
	////if method == "tevat.example.logic.Logic.WatchEventsEvents_PropsEvent" {
	////	msg := &logic.PropsEvent{}
	////	err = proto.Unmarshal(resBase.Msg, msg)
	////	fmt.Println("PropsEvent:", err, msg, msg.Props, msg.ProtoReflect(), msg.String())
	////	for i, v := range msg.GetProps() {
	////		fmt.Println("GetProps:", i, v)
	////	}
	////}
	////if method == "tevat.example.logic.Logic.WatchEventsEvents_MoneyEvent" {
	////	msg := &logic.MoneyEvent{}
	////	err = proto.Unmarshal(resBase.Msg, msg)
	////	fmt.Println("MoneyEvent:", err, msg, msg.Count, msg.MoneyType, msg.ProtoReflect(), msg.String())
	////}
	//res = c.ParseRes(resBase)
	//return res, err
}

func (c *Client) Rec() (Res, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("=====Rec panic=====%+v \n", r)
		}
	}()
	conn := c.conn
	if conn == nil {
		return Res{}, fmt.Errorf("conn is nil")
	}

	return c.Decode(conn)
	// log.Printf("[%v]:: m.Decode, v:%+v, err:%+v  \n", time.Now(), v, err)
}

func (c *Client) StartOnRec() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("=====StartOnRec panic=====%+v \n", r)
		}
	}()
	async.GoRaw(func() {
		for {
			res, err := c.Rec()
			if err != nil {
				fmt.Printf("[%v]::for onRec fail,err:%+v, res:%+v \n", time.Now(), err, res)
				fmt.Printf("Reconnect ...")
				time.Sleep(time.Second)
				c.Connect()
			}
			c.cond.L.Lock()
			c.res = res
			c.cond.L.Unlock()
			c.cond.Signal()
			if _, ok := c.config.ignoreMethodMap[res.Method]; ok {
				continue
			}
			if c.onRec != nil {
				c.onRec(res)
			}
		}
	})
}
