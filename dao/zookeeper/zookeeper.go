package zookeeper

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/go-zookeeper/zk"
	"schedule/util/config"
	"schedule/util/consts"
	"schedule/util/log"
	"schedule/util/tool"
	"strconv"
	"time"
)

var zkConn *zk.Conn

func InitZookeeper() error {
	conn, _, err := zk.Connect(config.Viper.GetStringSlice("zookeeper.hosts"), config.Viper.GetDuration("zookeeper.timeout")*time.Second)
	if err != nil {
		log.ErrLogger.Printf("init zookeeper error:%+v", err)
		return err
	}

	zkConn = conn
	if exist, err := ExistNode("/go_schedule"); err != nil {
		log.ErrLogger.Printf("init zookeeper error:%+v", err)
		return err
	} else if !exist {
		if _, err := CreateNode("/go_schedule", nil); err != nil {
			log.ErrLogger.Printf("init zookeeper error:%+v", err)
			return err
		}
	}

	if exist, err := ExistNode(consts.ZKLockPath); err != nil {
		log.ErrLogger.Printf("init zookeeper error:%+v", err)
		return err
	} else if !exist {
		if _, err := CreateNode(consts.ZKLockPath, nil); err != nil {
			log.ErrLogger.Printf("init zookeeper error:%+v", err)
			return err
		}
	}

	if exist, err := ExistNode("/go_schedule/schedule"); err != nil {
		log.ErrLogger.Printf("init zookeeper error:%+v", err)
		return err
	} else if !exist {
		if _, err := CreateNode("/go_schedule/schedule", nil); err != nil {
			log.ErrLogger.Printf("init zookeeper error:%+v", err)
			return err
		}
	}

	hash := md5.New()
	if _, err := hash.Write([]byte(tool.IP)); err != nil {
		log.ErrLogger.Printf("init zookeeper error:%+v", err)
		return err
	}

	result := hash.Sum(nil)
	m := binary.BigEndian.Uint32(result)
	md5Str := strconv.FormatUint(uint64(m), 10)

	if _, err := CreateTemplateNode(fmt.Sprintf("/go_schedule/schedule/%s", tool.IP), []byte(md5Str)); err != nil {
		log.ErrLogger.Printf("init zookeeper error:%+v", err)
		return err
	}
	return nil
}

// CreateTemplateNode 创建临时节点
func CreateTemplateNode(path string, data []byte) (result string, err error) {
	for i := 0; i < 3; i++ {
		if result, err = zkConn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); err != nil {
			break
		}
	}
	if err != nil {
		log.ErrLogger.Printf("create template node fail, path:%s, data:%s, error:%+v", path, string(data), err)
		return "", err
	}

	return result, nil

}

// CreateNode 创建持久节点
func CreateNode(path string, data []byte) (string, error) {
	var result string
	var err error
	i := 0
	for i = 0; i < 3; i++ {
		if result, err = zkConn.Create(path, data, 0, zk.WorldACL(zk.PermAll)); err == nil {
			break
		}
	}
	if i >= 3 {
		log.ErrLogger.Printf("create node fail, path:%s, data:%s, error:%+v", path, string(data), err)
		return "", err
	}
	return result, nil
}

func SetData(path string, data []byte, version int32) error {
	var stat *zk.Stat
	var err error
	for i := 0; i < 3; i++ {
		if stat, err = zkConn.Set(path, data, version); err == nil {
			break
		}
	}
	if err != nil {
		log.ErrLogger.Printf("set data to zk fail, error:%+v, path:%s, data:%s, version:%d", err, path, data, version)
		return err
	}

	log.InfoLogger.Printf("set data status:%+v", *stat)
	return nil
}

// GetData 获取数据
func GetData(path string) ([]byte, error) {
	var data []byte
	var err error
	for i := 0; i < 3; i++ {
		if data, _, err = zkConn.Get(path); err == nil {
			break
		}
	}

	if err != nil {
		log.ErrLogger.Printf("get data from zk fail, error:%+v, path:%s", err, path)
		return nil, err
	}

	return data, nil
}

// DeleteNode 删除节点
func DeleteNode(path string) (err error) {
	for i := 0; i < 3; i++ {
		if err = zkConn.Delete(path, -1); err == nil {
			break
		}
	}
	if err != nil {
		log.ErrLogger.Printf("delete zk node fail, error:%+v", err)
	}

	return err
}

// ExistNode 判断节点是否存在
func ExistNode(path string) (exist bool, err error) {
	for i := 0; i < 3; i++ {
		if exist, _, err = zkConn.Exists(path); err == nil {
			break
		}
	}
	if err != nil {
		log.ErrLogger.Printf("request zk exist fail, path:%s, error:%+v", path, err)
	}
	return exist, err
}

// ChildrenNodes 获取子节点
func ChildrenNodes(path string) (list []string, err error) {
	for i := 0; i < 3; i++ {
		if list, _, err = zkConn.Children(path); err == nil {
			break
		}
	}
	if err != nil {
		log.ErrLogger.Printf("get children node fail, path:%s, error:%+v", path, err)
	}
	return list, err
}

// ChildrenWatch 监听子节点变化
func ChildrenWatch(path string) (list []string, wchann <-chan zk.Event, err error) {
	for i := 0; i < 3; i++ {
		if list, _, wchann, err = zkConn.ChildrenW(path); err == nil {
			break
		}
	}
	if err != nil {
		log.ErrLogger.Printf("children watch fail, path:%s, error:%+v", path, err)
	}
	return
}
