package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"server/zinx/utils"
	"server/zinx/ziface"
)

// DataPack 封包拆包的具体模块
type DataPack struct {
}

// NewDataPack 拆包封包实例的初始化方法{ }
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取报的头部长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	//dataLen uint32(4 bytes) + ID uint32(4 bytes)
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放字节的缓存
	var test []byte
	dataBuffer := bytes.NewBuffer(test)
	//将DataLen写进DataBuffer中，将msg id写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

// Unpack 拆包方法,只需要将包的Head信息读出来，之后再根据head信息里面的data的长度
func (dp *DataPack) Unpack(dataBytes []byte) (ziface.IMessage, error) {
	// 创建一个IO.read对象
	dataBuffer := bytes.NewReader(dataBytes)
	msg := &Message{}
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读dataLen
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//读取数据操作，判断dataLen是否已经超出了最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large Msg data recv")
	}

	//此时还没有将data的内容，只是解压了head信息
	return msg, nil
}
