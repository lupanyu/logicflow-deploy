<template>
  <div class="container">
     <div id="lf-container" ref="container"></div>
    
    <!-- 日志弹窗 -->
    <el-dialog 
      v-model="logVisible"
      :title="`${selectedNode?.id} 执行日志`"
      width="70%"
    >
      <pre class="log-pre">{{ selectedNode?.logs || '无日志' }}</pre>
      <div v-if="selectedNode?.error" class="error-message">
        错误信息：{{ selectedNode.error }}
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute,useRouter } from 'vue-router'
import LogicFlow from '@logicflow/core'
   import '@logicflow/core/lib/style/index.css'
  import '@logicflow/extension/lib/style/index.css'
  import {    registerStart,registerEnd, registerJava,registerWeb,registerJenkins,registerShell  } from '../LFComponents/nodes/'
  // import demoData from  './demo-data.json'  
  import { ElMessage } from 'element-plus'
  import edges from '../LFComponents/edges/index.js' 
 
const route = useRoute()
const router = useRouter()
const container = ref(null)
const lf = ref(null)
const logVisible = ref(false)
const selectedNode = ref(null)
const executionData = ref({})
const execResult = ref({})
const cors = {
      mode: 'cors', // 明确跨域模式
      headers: {
        'Content-Type': 'application/json'
      }
}


async function initLf () {
       // 画布配置
       // LogicFlow.use(SelectionSelect);

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
        usePassiveEvent: true,// 被动事件
        edgeTextDraggable: true,
        hoverOutline: false,
        isSilentMode: true, // 禁用所有交互
        stopScrollGraph: true, // 禁止画布滚动
        stopZoomGraph: true,  // 禁止缩放
        adjustNodePosition: false, // 禁止节点拖拽
        disabledPlugins: ['control'], // 禁用控制点件
 
          container: container.value,
      })
   
      lf.value.register(edges)
      lf.value.setDefaultEdgeType('myCurvedEdge')
      // 添加节点点击事件监听
lf.value.on('node:click', ({ data }) => {
  selectedNode.value = {
    ...executionData.value.nodeResults[data.id],
    id: data.id
  }
  logVisible.value = true
})
      initEvent()
      // 注册节点
      registerNode()

  }
  // 注册节点
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
    // 注册节点
    registerStart(lf.value,nodeConfig)
     registerEnd(lf.value,nodeConfig)
    registerWeb(lf.value,nodeConfig)
    registerJava(lf.value,nodeConfig)
    registerJenkins(lf.value,nodeConfig)
    registerShell(lf.value,nodeConfig)
    // 注册节点到拖拽面板里
    //lf.value.extension.dndPanel.setPatternItems(nodeList)
  }
// 在fetchExecutionData之后添加以下方法
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

 
   // 事件监听
   function initEvent () {
    lf.value.on('blank:click', () => {
    logVisible.value = false
  })
  } 
// 获取执行记录
async function fetchExecutionData(flowId) {
  try {
    // url参数为执行记录id
    const url = "/api/v1/deploy/"+flowId
    console.log(url)
    const response = await fetch(url,cors)
    console.log(response)
    if (response.status === 404) {
      ElMessage.warning('请求的执行记录不存在...')
      // 跳转到首页

      setTimeout(() => {router.push('/')}, 2000)
      return
    }
    if (!response.ok) {
      ElMessage.error('Oops, this is a error message.')
    }
    executionData.value = await response.json()
    console.log(executionData.value)
    updateEdgeStatus() // 新增调用

    renderFlow()
  } catch (e) {
    console.log('获取执行记录失败', e)
  }
}

// 渲染流程图
function renderFlow() {
    // 确保每次渲染都使用最新数据
    const nodes = executionData.value.flowData?.nodes?.map(node => ({
    ...node,
    properties: {
      ...node.properties,
      // 添加状态类名用于CSS匹配
      class: `status-${executionData.value.nodeResults[node.id]?.status || 'pending'}`
    }
  })) || []

  const edges = executionData.value.flowData?.edges?.map(edge => ({
    ...edge,
    type: 'myCurvedEdge',
    text: edge.properties?.text,
    properties: {
      ...edge.properties,
      // 强制应用样式
      style: edge.properties?.style 
    }
  })) || []
    if (!lf.value) {
    lf.value = new LogicFlow({
      container: container.value,
      grid: true,
      isSilentMode: true, // 禁用交互
      disabledPlugins: ['control', 'dnd-panel'], // 禁用不需要的插件
     })
    
    // 自定义节点点击事件
    lf.value.on('node:click', ({ data }) => {
      selectedNode.value = executionData.value.nodeResults[data.id]
      logVisible.value = true
    })
        // 添加样式更新逻辑
        lf.value.on('render:finished', () => {
      nodes.forEach(node => {
        const element = document.querySelector(`[data-id="${node.id}"]`)
        if (element) {
          element.className.baseVal = node.properties.class
        }
      })})
  }
  

  
  lf.value.render({
    nodes,
    edges: executionData.value.flowData?.edges?.map(edge => ({
      ...edge,
      type: 'myCurvedEdge', // 确保使用自定义边类型
      text: edge.properties?.text // 传递状态文本
    })) || []
})
}

// 轮询更新数据
let timer = null
onMounted(() => {
  initLf()
  const flowId = props.flowId
  console.log(flowId)
  fetchExecutionData(flowId)
  // fetchExecutionData()
  timer = setInterval(()=> {
    if (execResult.value.status === 'running') {
      fetchExecutionData(props.execResult.flowId)
    }else{
      // 如果不是running 就 停止轮询
      clearInterval(timer)
     }
      },5000)
})

onBeforeUnmount(() => {
  clearInterval(timer)
})

defineExpose({
    executionData
    })  
</script>

<style scoped>
.log-pre {
  white-space: pre-wrap;
  max-height: 60vh;
  overflow: auto;
}

.error-message {
  color: #f56c6c;
  margin-top: 15px;
}

/* 根据状态设置节点颜色 */
:deep(.lf-node) {
  &.status-success rect { fill: #e1f3d8; }
  &.status-failed rect { fill: #fde2e2; }
  &.status-running rect { fill: #fdf6ec; }
}
</style>
