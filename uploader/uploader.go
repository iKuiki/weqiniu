package uploader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ikuiki/wwdk"
	"time"
	"wegate/wechat"
	"weqiniu/conf"

	"github.com/qiniu/api.v7/storage"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"wegate/common"
	commontest "wegate/common/test"
)

// Uploader 上传者服务模块
type Uploader interface {
	Serve()
}

// NewUploader 创建新的上传服务模块
func NewUploader(conf conf.Conf) Uploader {
	u := &uploader{
		conf: conf,
	}
	return u
}

type uploader struct {
	conf conf.Conf
}

// Serve 运行
func (u *uploader) Serve() {
	w := u.prepareConnect()
	fileChan := make(chan wechat.MediaFile)
	w.On("upload", func(client MQTT.Client, msg MQTT.Message) {
		var file wechat.MediaFile
		err := json.Unmarshal(msg.Payload(), &file)
		if err != nil {
			u.conf.GetLogger().Error("msg Unmarshal fail")
			return
		}
		fileChan <- file
	})
	resp, _ := w.Request("Wechat/HD_Upload_RegisterMQTTUploader", []byte(fmt.Sprintf(
		`{"name":"%s","description":"%s","uploadListenerTopic":"%s"}`,
		"qiniuUploader",
		"七牛上传模块",
		"upload",
	)))
	if resp.Ret != common.RetCodeOK {
		u.conf.GetLogger().Fatalf("注册uploader失败: %s", resp.Msg)
	}
	token := resp.Msg
	for {
		select {
		case file := <-fileChan:
			u.conf.GetLogger().Infof("upload(%d)[%s]: %s\n", file.MediaType, file.QueueID, file.FileName)
			startAt := time.Now()
			// 上传文件
			putPolicy := storage.PutPolicy{
				Scope: u.conf.GetQiniuBucketName(),
			}
			upToken := putPolicy.UploadToken(u.conf.GetQiniuMac())
			ret := storage.PutRet{}
			putExtra := storage.PutExtra{}
			dataLen := int64(len(file.BinaryContent))
			key := file.FileName
			switch file.MediaType {
			case wwdk.MediaTypeUserHeadImg:
				key = "user/" + key
			case wwdk.MediaTypeContactHeadImg:
				key = "contact/" + key
			case wwdk.MediaTypeMemberHeadImg:
				key = "member/" + key
			case wwdk.MediaTypeMessageImage:
				key = "message_image/" + key
			case wwdk.MediaTypeMessageVoice:
				key = "message_voice/" + key
			case wwdk.MediaTypeMessageVideo:
				key = "message_video/" + key
			}
			err := u.conf.GetQiniuFormUploader().Put(context.Background(), &ret, upToken, key, bytes.NewReader(file.BinaryContent), dataLen, &putExtra)
			if err != nil && err.Error() != "file exists" {
				u.conf.GetLogger().Errorf("QiniuFormUploader.Put error: %+v", err)
				break
			}
			u.conf.GetLogger().Infof("%s上传完成，花费%v\n", key, time.Since(startAt))
			// 上传完成的回调
			resp, _ := w.Request("Wechat/HD_Upload_MQTTUploadFinish", []byte(fmt.Sprintf(
				`{"token":"%s","queueID":"%s","fileurl":"%s"}`,
				token,
				file.QueueID,
				"http://"+u.conf.GetQiniuBucketDomain()+"/"+key,
			)))
			if resp.Ret != common.RetCodeOK {
				u.conf.GetLogger().Errorf("通知服务器上传完毕失败: %s\n", resp.Msg)
			}
		}
	}
}

// 准备连接
func (u *uploader) prepareConnect() (w commontest.Work) {
	w = commontest.Work{}
	opts := w.GetDefaultOptions(u.conf.GetWegateURL())
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		u.conf.GetLogger().Info("ConnectionLost", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		u.conf.GetLogger().Info("OnConnectHandler")
	})
	err := w.Connect(opts)
	if err != nil {
		panic(err)
	}
	pass := u.conf.GetWegatePassword() + time.Now().Format(time.RFC822)
	resp, _ := w.Request("Login/HD_Login", []byte(`{"username":"uploader","password":"`+pass+`"}`))
	if resp.Ret != common.RetCodeOK {
		u.conf.GetLogger().Fatalf("登录失败: %s", resp.Msg)
	}
	return
}
