/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/12/3
   Description :
-------------------------------------------------
*/

package mysql_binlog

import (
	"errors"
	"math/rand"
	"time"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/service"
)

type RegistryMysqlBinlogHandlerFunc = func (c core.IComponent) IEventHandler

func init() {
	service.RegisterCreator(core.MysqlBinlogService, new(mysqlBinlogCreator))
}

type mysqlBinlogCreator struct{}

func (*mysqlBinlogCreator) Create(app core.IApp) core.IService {
	return NewMysqlBinlogService(app)
}

type MysqlBinlogService struct {
	app core.IApp

	canal               *canal.Canal
	analyzer            *analyzer
	oldSchema, oldTable string

	handler IEventHandler
	err     error
}

func NewMysqlBinlogService(app core.IApp) core.IService {
	conf := app.GetConfig().Config().MysqlBinlogService
	cfg := &canal.Config{
		Addr:                  conf.Host,
		User:                  conf.UserName,
		Password:              conf.Password,
		Charset:               "utf8mb4",
		ServerID:              uint32(rand.New(rand.NewSource(time.Now().Unix())).Intn(1000)) + 1001,
		Flavor:                "mysql",
		DiscardNoMetaRowEvent: conf.DiscardNoMetaRowEvent,
		Dump: canal.DumpConfig{
			ExecutionPath:  conf.DumpExecutionPath,
			DiscardErr:     true,
			SkipMasterData: false,
		},
	}
	if conf.Charset != nil {
		cfg.Charset = *conf.Charset
	}
	if len(conf.IncludeTableRegex) > 0 {
		cfg.IncludeTableRegex = append([]string{}, conf.IncludeTableRegex...)
	}
	if len(conf.ExcludeTableRegex) > 0 {
		cfg.ExcludeTableRegex = append([]string{}, conf.ExcludeTableRegex...)
	}

	ca, err := canal.NewCanal(cfg)

	return &MysqlBinlogService{
		app:      app,
		canal:    ca,
		analyzer: newAnalyzer(app, conf.IgnoreWKBDataParseError),
		err:      err,
	}
}

func (m *MysqlBinlogService) Inject(a ...interface{}) {
	if m.handler != nil {
		m.app.Fatal("mysql-binlog服务重复注入")
	}

	if len(a) != 1 {
		m.app.Fatal("mysql-binlog服务注入数量必须为1个")
	}

	fn, ok := a[0].(RegistryMysqlBinlogHandlerFunc)
	if !ok {
		m.app.Fatal("mysql-binlog服务注入类型错误, 它必须能转为 mysql_binlog.RegistryMysqlBinlogHandlerFunc")
	}

	m.handler = fn(m.app.GetComponent())
}

func (m *MysqlBinlogService) Start() error {
	if m.err != nil {
		return m.err
	}
	if m.handler == nil {
		return errors.New("未注入handler")
	}

	m.canal.SetEventHandler(m)

	binlogName, pos, err := m.handler.GetStartPos()
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func(errChan chan error) {
		switch binlogName {
		case OldestPos:
			if err := m.canal.Run(); err != nil {
				errChan <- err
			}
		case LatestPos:
			pos, err := m.canal.GetMasterPos()
			if err == nil {
				err = m.canal.RunFrom(pos)
			}
			if err != nil {
				errChan <- err
			}
		default:
			if err := m.canal.RunFrom(mysql.Position{Name: binlogName, Pos: pos}); err != nil {
				errChan <- err
			}
		}
	}(errChan)

	select {
	case <-time.After(time.Second):
	case err := <-errChan:
		return err
	}

	m.app.Debug("mysql-bing服务启动成功")
	return nil
}

func (m *MysqlBinlogService) Close() error {
	m.canal.Close()
	return nil
}
