package server

import (
	"encoding/json"
	"fmt"
	"log"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"os"
	"sync"
)

type Storage interface {
	Save(execution *schema.FlowExecution)
	Get(flowID string) (*schema.FlowExecution, bool)
	GetAll() []*schema.FlowExecution
}

// MemoryStorage 结构体用于内存存储
type MemoryStorage struct {
	executions map[string]*schema.FlowExecution
	// map读写锁
	lock sync.RWMutex
}

// GetAll 方法用于获取所有的流程执行状态
func (ms *MemoryStorage) GetAll() []*schema.FlowExecution {
	ms.lock.RLock()
	defer ms.lock.RUnlock()
	var executions []*schema.FlowExecution
	for _, execution := range ms.executions {
		executions = append(executions, execution)
	}
	return executions
}

// NewMemoryStorage 创建一个新的 MemoryStorage 实例
func NewMemoryStorage() Storage {
	return &MemoryStorage{
		executions: make(map[string]*schema.FlowExecution),
	}
}

// Save 方法用于保存流程执行状态
func (ms *MemoryStorage) Save(execution *schema.FlowExecution) {
	ms.lock.Lock()
	ms.executions[execution.FlowID] = execution
	ms.lock.Unlock()
	log.Printf(" [%s]Saved flow execution with ID: %s", utils.GetCallerInfo(), execution.FlowID)
}

// Get 方法用于获取指定 ID 的流程执行状态
func (ms *MemoryStorage) Get(flowID string) (*schema.FlowExecution, bool) {
	ms.lock.RLock()
	execution, exists := ms.executions[flowID]
	ms.lock.RUnlock()
	return execution, exists
}

type FileStorage struct {
	executions map[string]schema.FlowExecution
	files      map[string][]byte // 存储文件内容的映射
	path       string            // 存储文件的路径
}

var defaultPath = "./storage"

func NewFileStorage(path string) FileStorage {
	if path == "" {
		path = defaultPath
	}
	// 创建存储文件的目录,如果不存在就创建
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
		return FileStorage{path: path, executions: make(map[string]schema.FlowExecution), files: make(map[string][]byte)}
	} else {
		// 加载目录中的已执行过的流程
		executions, err := loadExecutionsFromDirectory(path)
		if err != nil {
			panic(err)
		}
		return FileStorage{path: path, executions: executions, files: make(map[string][]byte)}
	}
}

func (fs FileStorage) Save(execution schema.FlowExecution) {
	fs.executions[execution.FlowID] = execution
	// 保存内容到文件
	file, err := os.OpenFile("./"+fs.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	// 把执行结果序列化
	content, err := json.Marshal(execution)
	if err != nil {
		log.Println("Error marshalling flow execution:", file.Name(), err.Error())
		return
	}
	_, err = file.WriteString(string(content))
	if err != nil {
		log.Println("Error write flow execution:", file.Name(), err.Error())
		return
	}
}

func (fs FileStorage) Get(flowID string) (schema.FlowExecution, bool) {
	execution, exists := fs.executions[flowID]
	return execution, exists
}

func loadExecutionsFromDirectory(path string) (map[string]schema.FlowExecution, error) {
	executions := make(map[string]schema.FlowExecution)
	// 读取目录中的文件
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		var s schema.FlowExecution
		// 读取文件内容
		content, err := os.ReadFile(path + "/" + file.Name())
		if err != nil {
			fmt.Println("加载文件【%s】内容出错:%v", file.Name(), err)
			continue
		}
		// 把json文件的内容反序列化
		if err := json.Unmarshal(content, &s); err != nil {
			fmt.Println("解析文件【%s】内容出错:%v", file.Name(), err)
			continue
		}
		executions[s.FlowID] = s
	}
	return executions, nil
}
