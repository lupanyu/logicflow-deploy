<template>
    <div class="flow-list-container">
      <div class="header">
 
      </div>
  
      <el-table :data="flowList" v-loading="loading">
        <el-table-column prop="name" label="流程名称" width="200" />
        <el-table-column label="最后更新时间" width="220">
          <template #default="{row}">
            {{ formatTime(row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{row}">
            <el-button 
              type="primary" 
              link
              @click="viewFlowDetail(row.name)"
            >
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import dayjs from 'dayjs'
  import { ElMessage } from 'element-plus'
  
  const router = useRouter()
  const flowList = ref([])
  const loading = ref(false)
  
  // 获取流程列表
  async function fetchFlows() {
    try {
      loading.value = true
      const response = await fetch('/api/v1/flow', {
      method: 'GET',
       headers: {
        'Content-Type': 'application/json'
      }
    })
      if (!response.ok) throw new Error('获取数据失败')
      const result = await response.json()
    console.log(result)
      flowList.value =  result.data || []
    } catch (e) {
      ElMessage.error(e.message)
    } finally {
      loading.value = false
    }
  }
  
  // 时间格式化
  function formatTime(timestamp) {
    return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
  }
  
  // 查看详情
  function viewFlowDetail(name) {
    router.push(`/flow/${name}`)
  }
  
  // 创建新流程
  function createNewFlow() {
    router.push('/flow/new')
  }
  
  onMounted(() => {
    fetchFlows() 
  }
  )
  </script>
  
  <style scoped>
  .flow-list-container {
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .header {
    margin-bottom: 20px;
    display: flex;
    justify-content: flex-end;
  }
  </style>