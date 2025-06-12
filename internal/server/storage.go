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
	feKeys     []string
	maxNum     int
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
func NewMemoryStorage(maxNum int, jsonFile string) Storage {
	ms := &MemoryStorage{
		feKeys:     make([]string, 0, maxNum+1),
		maxNum:     maxNum,
		executions: make(map[string]*schema.FlowExecution),
	}
	// 新增文件加载逻辑
	if file, err := os.Open(jsonFile); err == nil {
		defer file.Close()
		var data map[string]*schema.FlowExecution
		if err := json.NewDecoder(file).Decode(&data); err == nil {
			ms.lock.Lock()
			for flowID, execution := range data {
				ms.executions[flowID] = execution
				ms.feKeys = append(ms.feKeys, flowID)
				// 保持最大数量限制
				if len(ms.feKeys) > maxNum {
					oldest := ms.feKeys[0]
					delete(ms.executions, oldest)
					ms.feKeys = ms.feKeys[1:]
				}
			}
			ms.lock.Unlock()
			log.Printf("Loaded %d executions from storage file", len(data))
		}
	}

	return ms
}

// Save 方法用于保存流程执行状态
func (ms *MemoryStorage) Save(execution *schema.FlowExecution) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.executions[execution.FlowID] = execution
	isExist := false
	// 如果不存在当前key,则添加到切片中
	for _, key := range ms.feKeys {
		if key == execution.FlowID {
			isExist = true
			break
		}
	}
	if !isExist {
		ms.feKeys = append(ms.feKeys, execution.FlowID)
	}
	// 清理过期处理器（修复条件判断）
	if len(ms.feKeys) > ms.maxNum { // 改为>=30确保31个时触发清理
		oldest := ms.feKeys[0]
		if _, exists := ms.executions[oldest]; exists {
			delete(ms.executions, oldest)
		}
		ms.feKeys = ms.feKeys[1:] // 使用切片操作保持顺序
	}
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

// 新增内存存储转文件存储的方法
func SaveMemStorageToFile(ms *MemoryStorage, path string) error {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ms.executions)
}
