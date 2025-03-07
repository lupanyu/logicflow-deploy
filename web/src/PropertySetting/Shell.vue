<template>
    <div>
      <el-form :model="form" label-width="100px">
        <el-form-item label="应用名称">
          <el-input v-model="form.appName" placeholder="请输入应用名"></el-input>
        </el-form-item>
        <el-form-item label="主机名称">
          <el-input v-model="form.host" placeholder="请输入主机名"></el-input>
        </el-form-item>
        
        <el-form-item label="前置脚本内容">
          <el-input type="textarea" v-model="form.preScript" 
            :autosize="{ minRows: 4 }" 
            placeholder="请输入前置脚本内容（如环境准备）"/>
        </el-form-item>
        <el-form-item label="部署脚本内容">
          <el-input type="textarea" v-model="form.deployScript"
            :autosize="{ minRows: 6 }"
            placeholder="请输入部署脚本内容（核心部署逻辑）"/>
        </el-form-item>
        <el-form-item label="后置脚本内容">
          <el-input type="textarea" v-model="form.postScript"
            :autosize="{ minRows: 4 }"
            placeholder="请输入后置脚本内容（如清理操作）"/>
        </el-form-item>

        <el-form-item>
            <el-button type="primary" @click="onSubmit">保存</el-button>
        </el-form-item>
      </el-form>
    </div>
</template>

<script setup>
  import { onMounted, watch, reactive } from 'vue'
  
  const props = defineProps(['nodeData', 'lf'])
  const emit = defineEmits(['onClose'])

  const form = reactive({
    appName: '',
    host: '',
    preScript: '#!/bin/bash\necho pre  script',
    deployScript: '#!/bin/bash\necho deploy script',
    postScript: '#!/bin/bash\necho deploy script'
  })

  onMounted(() => {
    if (props.nodeData?.properties) {
      Object.assign(form, props.nodeData.properties)
    }
  })

  watch(
    () => props.nodeData?.properties,
    (newProperties) => {
      if (newProperties) {
        form.appName = newProperties.appName || ''
        form.host = newProperties.host || ''
        form.preScript = newProperties.preScript || ''
        form.deployScript = newProperties.deployScript || ''
        form.postScript = newProperties.postScript || ''
      }
    },
    { immediate: true }
  )

  function onSubmit() {
    const nodeModel = props.lf.getNodeModelById(props.nodeData.id)
    const value = `${form.appName}_${form.host}`
    
    nodeModel.text.value = value
    nodeModel.updateText(value)
    
    props.nodeData.properties = { ...form }
    props.lf.setProperties(props.nodeData.id, { ...form })
    emit('onClose')
  }
</script>
