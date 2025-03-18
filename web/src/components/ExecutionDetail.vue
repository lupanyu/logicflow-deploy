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
import { registerStart, registerEnd, registerJava, registerWeb, registerJenkins, registerShell } from '../LFComponents/nodes/'
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
  }
  
  function registerNode () {
    const nodeConfig = {
      width: 60,
      height: 60,
      style: {
        background: {
          fill: '#FFF',
          stroke: '#DDD'
        }
      }
    };
    
    registerStart(lf.value,nodeConfig)
    registerEnd(lf.value,nodeConfig)
    registerWeb(lf.value,nodeConfig)
    registerJava(lf.value,nodeConfig)
    registerJenkins(lf.value,nodeConfig)
    registerShell(lf.value,nodeConfig)
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
    const nodes = executionData.value.flowData?.nodes?.map(node => ({
      ...node,
      properties: {
        ...node.properties,
        class: `status-${executionData.value.nodeResults[node.id]?.status || 'pending'}`
      }
    })) || []

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
      
      lf.value.on('render:finished', () => {
        nodes.forEach(node => {
          const element = document.querySelector(`[data-id="${node.id}"]`)
          if (element) {
            element.className.baseVal = node.properties.class
          }
        })
      })
    }
  
    lf.value.render({
      nodes,
      edges: executionData.value.flowData?.edges?.map(edge => ({
        ...edge,
        type: 'myCurvedEdge',
        text: edge.properties?.text,
        properties: edge.properties
      })) || []
    })
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

<style scoped>
.log-pre {
  line-height: 1.2;
  margin: 0;
  padding: 8px;
  font-size: 14px;
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

/* 根据状态设置节点颜色 */
:deep(.lf-node) {
  &.status-success rect { fill: #e1f3d8; }
  &.status-failed rect { fill: #fde2e2; }
  &.status-running rect { fill: #fdf6ec; }
}
</style>