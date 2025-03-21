<template>
  <div class="container">
     <div id="lf-container" ref="container"></div>
    
    <!-- 日志弹窗 -->
    <el-dialog 
      v-model="logVisible"
      :title="`${selectedNode?.id} 执行日志`"
      width="70%"
    >
    <pre class="log-pre">
      <template v-if="selectedNode?.logs">
        <template v-for="(line, index) in formatLogs(selectedNode.logs)" :key="index">
          <span v-if="line.isTimestamp" class="timestamp">{{ line.content }}</span>
          <span v-else class="log-line">{{ line.content }}</span>
        </template>
      </template>
      <span v-else>无日志</span>
    </pre>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import LogicFlow from '@logicflow/core'
import '@logicflow/core/lib/style/index.css'
import '@logicflow/extension/lib/style/index.css'
import {  RegisterNodes,StartNode  } from '../LFComponents/nodes/'
import { ElMessage } from 'element-plus'
import edges from '../LFComponents/edges/index.js' 
 
const router = useRouter()
const container = ref(null)
const lf = ref(null)
const logVisible = ref(false)
const selectedNode = ref(null)
const executionData = ref({})
const executionDetailResult = ref({})
const cors = {
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json'
      }
}

const props = defineProps({
  getSelectedDeployment: Function,
  setSelectedDeployment: Function,
})

function formatLogs(logs) {
  const logLines = logs.split('\n').filter(line => line.trim() !== '');
  const timestampReg = /^(\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\])(.*)/;
  const result = [];

  logLines.forEach(line => {
    const match = line.match(timestampReg);
    if (match) {
      const timestamp = match[1];
      const content = match[2].trim();
      if (timestamp) {
        result.push({ isTimestamp: true, content: timestamp });
      }
      if (content) {
        result.push({ isTimestamp: false, content });
      }
    } else if (line.trim()) {
      result.push({ isTimestamp: false, content: line });
    }
  });
  return result;
}

// 自定义节点类型，添加动画效果
function registerCustomNodes() {
   console.log("在executionDetail里注册自定义节点")
  // 定义状态颜色
  const statusColors = {
    pending: { fill: '#F0F4FF', stroke: '#B3C0DB', strokeWidth: 1 },  // 浅蓝灰
  running: { fill: '#FFF7E6', stroke: '#FFA940', strokeWidth: 3 },  // 更明亮的橙色
  success: { fill: '#F6FFED', stroke: '#73D13D', strokeWidth: 2 }, // 鲜绿色
  failed: { fill: '#FFF1F0', stroke: '#FF4D4F', strokeWidth: 2 },    // 红色
    skipped: { fill: '#f4f4f5', stroke: '#909399', strokeWidth: 1 },
    error: { fill: '#fde2e2', stroke: '#f56c6c', strokeWidth: 2 },
    timeout: { fill: '#fde2e2', stroke: '#f56c6c', strokeWidth: 2 },
    rollbacked: { fill: '#fde2e2', stroke: '#f56c6c', strokeWidth: 2 }
  }

  // 为每种节点类型创建状态变体
  const nodeTypes = ['start', 'end', 'java', 'web', 'jenkins', 'shell']
  
  nodeTypes.forEach(baseType => {
    console.log(baseType,lf.value)
    // 获取原始节点类型
    const BaseNodeModel = lf.value.graphModel.getModel(baseType)
    if (!BaseNodeModel){
      console.log("未找到节点类型",baseType)
      return}

    
    // 扩展节点模型，添加状态支持
    class StatusNodeModel extends BaseNodeModel {
      constructor(data, graphModel) {
        super(data, graphModel)
        this.status = data.properties?.status || 'pending'
      }
      
      getNodeStyle() {
        const style = super.getNodeStyle()
        const statusStyle = statusColors[this.status] || statusColors.pending
        
        return {
          ...style,
          fill: statusStyle.fill,
          stroke: statusStyle.stroke,
          strokeWidth: statusStyle.strokeWidth,
          // 如果是运行中状态，添加动画类
          className: this.status === 'running' ? 'running-animation' : ''
        }
      }
      
      // 更新节点状态
      updateStatus(status) {
        if (this.status !== status) {
          this.status = status
          // 强制更新节点属性
          this.setAttributes({
            fill: statusColors[status].fill,
            stroke: statusColors[status].stroke,
            strokeWidth: statusColors[status].strokeWidth
          })
          // 触发双重更新确保视图刷新
          this.graphModel.eventCenter.emit('element:update', {
            element: this,
            type: 'node'
          })
          this.graphModel.executeRender() // 添加强制渲染
        }
      } 
    }
    
    // 注册扩展后的节点类型
    lf.value.register({
      type: baseType,
      model: StatusNodeModel,
      // 保留原始视图
      view: lf.value.getView(baseType)
    })
    // 更新节点 
    console.log("注册节点",baseType,lf.value.graphModel.getModel(baseType))
  })
}

// 在 updateEdgeStatus 方法后添加新方法
function updateEdgeAnimations(runningNodeId) {
  // 设置边显示动画
  const activeEdges = executionData.value.flowData?.edges?.filter(edge => 
    edge.sourceNodeId === runningNodeId
  ) || []

  activeEdges.forEach(edge => {
    const edgeModel = lf.value.getEdgeModelById(edge.id)
    if (edgeModel) {
      edgeModel.setProperties({
        ...edge.properties,
        style: {
          ...edge.properties?.style,
          stroke: '#FFA940',
          strokeDasharray: '5 5',
          animation: 'flow 2s linear infinite'
        }
      })
       lf.value.openEdgeAnimation(edge.id)
    }
  })
  把非活动的边关闭动画
  executionData.value.flowData?.edges?.forEach(edge => {
    if (!activeEdges.some(e => e.id === edge.id)) {
      lf.value.closeEdgeAnimation(edge.id)
    }
  })
}

async function initLf () {
      lf.value = new LogicFlow({
        width: 1000,
        height: 600,
        background: {
            backgroundColor: '#f7f9ff',
        },
        grid: {
          size: 10,
          visible: true
        },
        keyboard: {
          enabled: true
        },
        usePassiveEvent: true,
        edgeTextDraggable: true,
        hoverOutline: false,
        isSilentMode: true,
        stopScrollGraph: true,
        stopZoomGraph: true,
        adjustNodePosition: false,
        disabledPlugins: ['control'],
        container: container.value,
      })
   
      lf.value.register(edges)
      lf.value.setDefaultEdgeType('myCurvedEdge')
      
      lf.value.on('node:click', ({ data }) => {
        selectedNode.value = {
          ...executionData.value.nodeResults[data.id],
          id: data.id
        }
        logVisible.value = true
      })
      
      initEvent()
      registerNode()
      registerCustomNodes()
  }
  
  function registerNode () {

    
    RegisterNodes(lf.value)
    // lf.value.setTheme({
    //  rect:{
    //   fill: 'transparent',
    //   // 恢复默认边框设置
    //   stroke: '#DDD',       // 恢复默认边框颜色
    //   strokeWidth: 1        // 恢复默认边框宽度
    //  } 
    // })
  }

function updateEdgeStatus() {
  const statusColors = {
    pending: '#CCCCCC',
    running: '#E6A23C',
    success: '#67C23A',
    failed: '#F56C6C',
    skipped: '#909399',
    error: '#F56C6C',
    timeout: '#F56C6C',
    rollbacked: '#F56C6C'
  }

  executionData.value.flowData?.edges?.forEach(edge => {
    const sourceNode = executionData.value.nodeResults[edge.sourceNodeId]
    if (sourceNode?.status) {
      edge.properties = {
        ...edge.properties,
        style: {
          stroke: statusColors[sourceNode.status],
          strokeWidth: 2
        },
        text: {
          content: sourceNode.status.toUpperCase(),
          style: {
            fill: statusColors[sourceNode.status],
            fontSize: 12
          }
        }
      }
    }
  })
}

function initEvent () {
  lf.value.on('blank:click', () => {
    logVisible.value = false
  })
} 

// 更新节点状态的新方法 - 使用LogicFlow API
function updateNodeVisualStates() {
  if (!lf.value) return
  let runningNodeIds = []  // 改为数组存储多个运行节点

  
  executionData.value.flowData?.nodes?.forEach(node => {
    const nodeStatus = executionData.value.nodeResults[node.id]?.status  
    const nodeModel = lf.value.getNodeModelById(node.id)

    if (nodeModel && typeof nodeModel.updateStatus === 'function') {
      nodeModel.updateStatus(nodeStatus)
      if (nodeStatus === 'running') {
        runningNodeIds.push(node.id)  // 收集所有运行节点ID
      }
    }
  })
  // 更新所有运行节点的边动画
  runningNodeIds.forEach(id => {
    updateEdgeAnimations(id)
  })
}

async function fetchExecutionData(flowId) {
  try {
    const url = "/api/v1/deploy/"+flowId
    const response = await fetch(url,cors)
    
    if (response.status === 404) {
      ElMessage.warning('请求的执行记录不存在...')
      setTimeout(() => {router.push('/')}, 2000)
      return
    }
    
    if (!response.ok) {
      ElMessage.error('Oops, this is a error message.')
    }
    
    executionData.value = await response.json()
    props.setSelectedDeployment(executionData.value)
    updateEdgeStatus()
    renderFlow()
  } catch (e) {
    console.log('获取执行记录失败', e)
  }
}

function renderFlow() {
    const nodes = executionData.value.flowData?.nodes?.map(node => {
      const nodeStatus = executionData.value.nodeResults[node.id]?.status || 'pending'
      return {
        ...node,
        properties: {
          ...node.properties,
          status: nodeStatus
        }
      }
    }) || []

    const edges = executionData.value.flowData?.edges?.map(edge => ({
      ...edge,
      type: 'myCurvedEdge',
      text: edge.properties?.text,
      properties: {
        ...edge.properties,
        style: edge.properties?.style 
      }
    })) || []
    
    if (!lf.value) {
      lf.value = new LogicFlow({
        container: container.value,
        grid: true,
        isSilentMode: true,
        disabledPlugins: ['control', 'dnd-panel'],
      })
      
      lf.value.on('node:click', ({ data }) => {
        selectedNode.value = executionData.value.nodeResults[data.id]
        logVisible.value = true
      })
      

    }
  
    lf.value.render({
      nodes,
      edges
    })
    
    // 渲染完成后更新节点状态
    setTimeout(() => {
      updateNodeVisualStates()
    }, 100)
}

let timer = null
onMounted(() => {
  initLf()
  executionData.value = props.getSelectedDeployment()
  
  if (!executionData.value.flowId) {
    console.log('未获取到flowId')
    return 
  }
  
  const flowId = executionData.value.flowId
  fetchExecutionData(flowId)
  
  timer = setInterval(() => {
    if (executionData.value.status === 'running') {
      fetchExecutionData(flowId)
    } else {
      clearInterval(timer)
    }
  }, 5000)
})

onBeforeUnmount(() => {
  clearInterval(timer)
})
</script>

<style>
/* 全局样式，确保动画效果可以应用到LogicFlow节点 */
@keyframes pulse {
  0% {
    stroke-width: 2px;
    stroke-opacity: 1;
  }
  50% {
    stroke-width: 4px;
    stroke-opacity: 0.8;
  }
  100% {
    stroke-width: 2px;
    stroke-opacity: 1;
  }
}
/* .running-edge {
  stroke: #FFA940;
  stroke-width: 2;
  stroke-dasharray: 5 5;
  animation: flow 2s linear infinite;
} */

@keyframes flow {
  0% {
    stroke-dashoffset: 10;
  }
  100% {
    stroke-dashoffset: 0;
  }
}
@keyframes glow {
  0% {
    filter: drop-shadow(0 0 2px rgba(230, 162, 60, 0.5));
  }
  50% {
    filter: drop-shadow(0 0 6px rgba(230, 162, 60, 0.8));
  }
  100% {
    filter: drop-shadow(0 0 2px rgba(230, 162, 60, 0.5));
  }
}

.running-animation {
  animation: pulse 1.5s infinite ease-in-out, glow 1.5s infinite ease-in-out;
}
</style>

<style scoped>
.log-pre {
  line-height: 1.2;
  margin: 0;
  padding: 8px;
  font-size: 12px;
  font-family: monospace;
  white-space: pre-wrap;
  overflow-y: auto;
  max-height: 60vh;
  display: flex;
  flex-direction: column;
}

.log-line, .timestamp {
  display: block;
  margin: 0;
  padding: 0;
  line-height: 1.2;
}

.timestamp {
  color: #67C23A;
  font-weight: 500;
  margin-top: 4px;
}

.log-line {
  padding-left: 20px;
}

/* 在 <style> 块中添加 */
.lf-node-content {
  background-color: inherit !important;
}

.lf-node-shape {
  fill: inherit !important;
  stroke: inherit !important;
  stroke-width: inherit !important;
}
</style>