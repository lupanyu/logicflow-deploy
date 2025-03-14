<template>
  <div class="jenkins-form">
    <el-form :model="form" label-width="100px">
      <el-form-item label="任务名称">
        <el-input 
          v-model="form.name" 
          placeholder="请输入Jenkins任务名称(标识显示用)"
          clearable
        ></el-input>
      </el-form-item>

      <el-form-item label="构建节点选择">
        <el-select
          v-model="form.nodeName"
          placeholder="请选择构建节点"
          filterable
          loading-text="加载任务列表中..."
          :loading="loading"
        >
          <el-option
            v-for="(jobs,node) in nodes"
            :key="node"
            :label="node"
            :value="node"
          />
        </el-select>
      </el-form-item>
 <!-- 新增节点选择表单项 -->
 <el-form-item label="构建任务选择">
  <el-select
    v-model="form.jobName"
    placeholder="请选择构建任务"
    filterable
    :disabled="!form.nodeName"
    loading-text="加载任务列表中..."
    :loading="loadingNodes"
  >
    <el-option 
      v-for="item in (nodes[form.nodeName] || []).filter(Boolean)"
      :key="String(item)"
      :label="String(item)"
      :value="String(item)"
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
const loading = ref(false)
const nodes=ref({})

 
// 获取jenkins任务相关数据 
async function fetchJenkinsData() {
  try {
    loading.value = true
    const response = await fetch('/api/v1/jenkins')
    if (!response.ok) {
      throw new Error('Network response was not ok')
    }   
    const data = await response.json()
    // 转换数据结构，确保所有节点值都是数组，并过滤无效值
    nodes.value = Object.fromEntries(
      Object.entries(data).map(([key, val]) => [
        key, 
        Array.isArray(val) ? val.filter(Boolean).map(String) : []
      ])
    )
    console.log(nodes.value)
    loading.value = false
  }catch (error) {
    console.error('There was a problem with the fetch operation:', error)
  }
}

// 初始化加载任务列表
onMounted(async () => {
//   await fetchJobs()
  await fetchJenkinsData()
  if (props.nodeData?.properties) {
    form.name = props.nodeData.properties.name || ''
    form.jobName = props.nodeData.properties.jobName || ''
    form.nodeName = props.nodeData.properties.nodeName || ''
  }
})


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
