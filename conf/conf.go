package conf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/kataras/golog"

	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

// Conf 配置
type Conf interface {
	// LoadJSON 从json文件载入配置
	LoadJSON(jsonFilepath string) (err error)
	GetLogger() *golog.Logger
	GetWegateURL() string
	GetWegatePassword() string
	GetQiniuMac() *qbox.Mac
	GetQiniuFormUploader() *storage.FormUploader
	GetQiniuBucketName() string
	GetQiniuBucketDomain() string
}

// NewConfig 创建新的配置实例
func NewConfig() Conf {
	return new(conf)
}

// conf 配置文件实现
type conf struct {
	// WegateURL wegate微信网关的url地址
	WegateURL string
	// WegatePassword wegate微信网关的接入密码
	WegatePassword string
	// QiniuAccessID 七牛的accessID
	QiniuAccessID string
	// QiniuAccessSecret 七牛的accessSecret
	QiniuAccessSecret string
	// QiniuBucketName 七牛空间名称
	QiniuBucketName string
	// QiniuBucketDomain 七牛空间域名
	QiniuBucketDomain string
	// 七牛上传组件
	qiniuMac          *qbox.Mac
	qiniuFormUploader *storage.FormUploader
	logger            *golog.Logger
}

// LoadJSON 从json文件载入配置
func (c *conf) LoadJSON(jsonFilepath string) (err error) {
	var data []byte
	buf := new(bytes.Buffer)
	f, err := os.Open(jsonFilepath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			if len(line) > 0 {
				buf.Write(line)
			}
			break
		}
		if !strings.HasPrefix(strings.TrimLeft(string(line), "\t "), "//") {
			buf.Write(line)
		}
	}
	data = buf.Bytes()
	err = json.Unmarshal(data, c)
	if err != nil {
		return errors.WithStack(err)
	}
	err = c.Init()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Init 根据载入的参数初始化
func (c *conf) Init() (err error) {
	c.qiniuMac = qbox.NewMac(c.QiniuAccessID, c.QiniuAccessSecret)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 空间对应的机房
		UseHTTPS:      false,               // 是否使用https域名
		UseCdnDomains: false,               // 上传是否使用CDN上传加速
	}
	c.qiniuFormUploader = storage.NewFormUploader(&cfg)
	c.logger = golog.New()
	return nil
}

func (c *conf) GetLogger() *golog.Logger {
	return c.logger
}

func (c *conf) GetWegateURL() string {
	return c.WegateURL
}

func (c *conf) GetWegatePassword() string {
	return c.WegatePassword
}

func (c *conf) GetQiniuMac() *qbox.Mac {
	return c.qiniuMac
}

func (c *conf) GetQiniuFormUploader() *storage.FormUploader {
	return c.qiniuFormUploader
}

func (c *conf) GetQiniuBucketName() string {
	return c.QiniuBucketName
}
func (c *conf) GetQiniuBucketDomain() string {
	return c.QiniuBucketDomain
}
