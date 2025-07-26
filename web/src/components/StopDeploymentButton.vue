<template>
    <el-button
      :type="variant"
      :size="size"
      :loading="isLoading"
      :disabled="isLoading"
      @click="handleStop"
      text
    >
      {{ isLoading ? 'stoping...' : 'stop' }}
    </el-button>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  
  const props = defineProps({
    deploymentId: {
      type: String,
      required: true
    },
    size: {
      type: String,
      default: 'large'
    },
    variant: {
      type: String,
      default: 'danger'
    }
  })
  
  const emit = defineEmits(['success'])
  
  const isLoading = ref(false)
  
  const handleStop = async () => {
    isLoading.value = true
  
    try {
      const response = await fetch(`/api/v1/deploy/stop/${props.deploymentId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      })
  
      const result = await response.json()
  
      if (response.ok) {
        ElMessage({
          message: result.message || '部署已成功停止',
          type: 'success',
          duration: 3000
        })
        emit('success')
      } else {
        ElMessage({
          message: result.error || '停止部署时发生错误',
          type: 'error',
          duration: 5000
        })
      }
    } catch (error) {
      ElMessage({
        message: '无法连接到服务器，请检查网络连接',
        type: 'error',
        duration: 5000
      })
    } finally {
      isLoading.value = false
    }
  }
  </script>