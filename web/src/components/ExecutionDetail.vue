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
import { useRoute } from 'vue-router'
import LogicFlow from '@logicflow/core'
 
const route = useRoute()
const container = ref(null)
const lf = ref(null)
const logVisible = ref(false)
const selectedNode = ref(null)
const executionData = ref({})
const cors = {
      mode: 'cors', // 明确跨域模式
      headers: {
        'Content-Type': 'application/json'
      }
}
// 获取执行记录
async function fetchExecutionData() {
  try {
    // url参数为执行记录id
    const url = "http://localhost:8080/api/v1/deploy/"+`${route.params.id}`
    console.log(url)
    const response = await fetch(url,cors)
    console.log(response)
    executionData.value = await response.json()
    console.log(executionData.value)

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

  const edges = executionData.value.flowData?.edges || []
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
    edges: executionData.value.flowData?.edges || []
  })
}

// 轮询更新数据
let timer = null
onMounted(() => {
  fetchExecutionData()
  timer = setInterval(fetchExecutionData, 5000) // 5秒刷新
})

onBeforeUnmount(() => {
  clearInterval(timer)
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
