<template>
  <div>
    <el-form label-width="80px" :model="formData">
      <el-form-item label="文案">
        <el-input v-model="formData.text"></el-input>
      </el-form-item>
      <el-form-item label="名称">
        <el-input v-model="formData.name"></el-input>
      </el-form-item>
      <el-form-item label="活动区域">
        <el-input v-model="formData.region"></el-input>
      </el-form-item>
      <el-form-item label="活动形式">
        <el-input v-model="formData.type"></el-input>
      </el-form-item>
      <el-form-item label="A">
        <el-input v-model="formData.a.a1"></el-input>
        <el-input v-model="formData.a.a2"></el-input>
      </el-form-item>
       <el-form-item>
        <el-button type="primary" @click="onSubmit">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'

// 通过 `defineProps` 来定义 props
const props = defineProps({
  nodeData: Object,
  lf: {
    type: [Object, String],
    required: true
  }
})

// 使用 `ref` 来定义响应式数据
const form = ref({
  name: '',
  region: '',
  date1: '',
  date2: '',
  delivery: false,
  type: [],
  resource: '',
  desc: ''
})

// 在 `onMounted` 中处理初始化逻辑
onMounted(() => {
  const { properties } = props.nodeData
  if (properties) {
    // 更新 `form` 数据
    form.value = { ...form.value, ...properties }
  }
})

// 使用 `defineEmits` 来定义事件
const emit = defineEmits(['onClose'])

// 提交处理方法
function onSubmit() {
  console.log('submit!');
  const nodeData = props.nodeData
  nodeData.properties = form.value
  props.lf.setProperties(nodeData.id, form.value)
  emit('onClose')
}
</script>

<style scoped>
</style>
