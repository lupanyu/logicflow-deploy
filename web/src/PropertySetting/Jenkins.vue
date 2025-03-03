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
  jobName: ''
})

// 任务列表相关状态
const jobs = ref([])
const loading = ref(false)

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

// 测试模仿数据
function fetchMockJobs() {
  jobs.value = [
    { label: 'Job 1', value: 'job1' },
    { label: 'Job 2', value: 'job2' },
    { label: 'Job 3', value: 'job3' }
  ]
  loading.value = false
}

// 初始化加载任务列表
onMounted(async () => {
//   await fetchJobs()
  fetchMockJobs()
  if (props.nodeData?.properties) {
    form.name = props.nodeData.properties.name || ''
    form.jobName = props.nodeData.properties.jobName || ''
  }
})

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
