<template>
  <div class="jenkins-form">
    <el-form :model="form" label-width="100px">
      <el-form-item label="任务名称">
        <el-input 
          v-model="form.name" 
          placeholder="请输入Jenkins任务名称"
          clearable
        ></el-input>
      </el-form-item>

      <el-form-item label="构建任务">
        <el-select
          v-model="form.jobName"
          placeholder="请选择构建任务"
          filterable
          loading-text="加载任务列表中..."
          :loading="loading"
        >
          <el-option
            v-for="job in jobs"
            :key="job.value"
            :label="job.label"
            :value="job.value"
          />
        </el-select>
      </el-form-item>
 <!-- 新增节点选择表单项 -->
 <el-form-item label="构建节点">
        <el-select
          v-model="form.nodeName"
          placeholder="请选择构建节点"
          filterable
          :disabled="!form.jobName"
          loading-text="加载节点列表中..."
          :loading="loadingNodes"
        >
          <el-option
            v-for="node in nodes"
            :key="node.value"
            :label="node.label"
            :value="node.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">保存配置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

const props = defineProps(['nodeData', 'lf'])
const emit = defineEmits(['onClose'])

// 表单数据结构
 const form = reactive({
  name: '',
  jobName: '',
  nodeName: '' // 新增节点字段
})

// 任务列表相关状态
const jobs = ref([])
const loading = ref(false)

// 节点相关状态
const nodes = ref([])
const loadingNodes = ref(false)

// 获取Jenkins任务列表
async function fetchJobs() {
  try {
    loading.value = true
    const response = await fetch('/api/jenkins/jobs')
    const data = await response.json()
    jobs.value = data.map(job => ({
      label: job.fullName,
      value: job.name
    }))
  } finally {
    loading.value = false
  }
}

// 获取节点列表方法
async function fetchNodes(jobName) {
  try {
    loadingNodes.value = true
    const response = await fetch(`/api/jenkins/jobs/${jobName}/nodes`)
    const data = await response.json()
    nodes.value = data.map(node => ({
      label: node.name,
      value: node.name
    }))
  } finally {
    loadingNodes.value = false
  }
}

// 测试模仿数据
function fetchMockJobs() {
  jobs.value = [
    { label: 'Job 1', value: 'job1' },
    { label: 'Job 2', value: 'job2' },
    { label: 'Job 3', value: 'job3' }
  ]
  loading.value = false
}

// 测试模仿数据
function fetchMockNodes() {
  nodes.value = [
    { label: 'Node 1', value: 'node1' },
    { label: 'Node 2', value: 'node2' },
    { label: 'Node 3', value: 'node3' }
  ]
  loadingNodes.value = false
}

// 初始化加载任务列表
onMounted(async () => {
//   await fetchJobs()
  fetchMockJobs()
  fetchMockNodes()
  if (props.nodeData?.properties) {
    form.name = props.nodeData.properties.name || ''
    form.jobName = props.nodeData.properties.jobName || ''
    form.nodeName = props.nodeData.properties.nodeName || ''
  }
})

// // 监听构建任务变化
// watch(() => form.jobName, async (newJob) => {
//   if (newJob) {
//     await fetchNodes(newJob)
//   } else {
//     nodes.value = []
//     form.nodeName = ''
//   }
// })
// 提交处理
function onSubmit() {
  const nodeModel = props.lf.getNodeModelById(props.nodeData.id)
  nodeModel.text.value = form.name
  nodeModel.updateText(form.name)
  
  props.nodeData.properties = { ...form }
  props.lf.setProperties(props.nodeData.id, { ...form })
  emit('onClose')
}
</script>

<style scoped>
.jenkins-form {
  padding: 20px;
  max-width: 600px;
}
</style>
