<template>
    <div>
      <el-form :model="form" label-width="80px">
        <el-form-item label="应用名称">
          <el-input v-model="form.appName" placeholder="请输入应用名"></el-input>
        </el-form-item>
        <el-form-item label="主机名称">
          <el-input v-model="form.host" placeholder="请输入主机名"></el-input>
        </el-form-item>
        <el-form-item label="jar包地址">
          <el-input type="textarea" v-model="form.packageSource" placeholder="jar包名称"></el-input>
        </el-form-item>
        <el-form-item label="部署目录">
            <el-input type="textarea" v-model="form.deployPath" placeholder="部署目录包含文件名"></el-input>
        </el-form-item>
       <el-form-item label="备份包名称">
            <el-input type="textarea" v-model="form.bakPath" placeholder="备份包的目录包含文件名"></el-input>
       </el-form-item>
       <el-form-item label="服务名称">
            <el-input type="textarea" v-model="form.serverName" placeholder="服务托管名称"></el-input>
       </el-form-item>
       <el-form-item label="健康检查端口">
            <el-input type="number" v-model="form.port" placeholder="请输入健康检查端口"></el-input>
       </el-form-item>
       <el-form-item label="健康检查路径">
            <el-input type="textarea" v-model="form.healthUri" placeholder="请输入健康检查路径"></el-input>
       </el-form-item>
       <el-form-item label="健康检查超时时间">
            <el-input type="textarea" v-model="form.healthCheckTimeout" placeholder="请输入健康检查超时时间"></el-input>
       </el-form-item>
        <el-form-item> 
            <el-button type="primary" @click="onSubmit">保存</el-button>
        </el-form-item>
      </el-form>
    </div>
  </template>
<script setup>
  import {   onMounted,watch, reactive } from 'vue'
  
  // 使用 `defineProps` 来定义 props
  const props = defineProps(['nodeData', 'lf'])

  // 使用 `reactive` 来定义响应式数据
  const form = reactive({
    appName: 'defaultAppName',
    host:  'test01',
    packageSource:   '/data/services/app/smallheart-app-entrance.jar',
    deployPath: '/data/services/app/smallheart-app-entrance.jar',
    bakPath: '/data/services/app/smallheart-app-entrance.jar.last',
    serverName:'app',
    port: 8080,
    healthUri: '/smallheart-mq-subscribe/api/mq/status',
    healthCheckTimeout: 180,
  })

  // 使用 `defineEmits` 来定义事件
  const emit = defineEmits(['onClose'])

  // onMounted 钩子
  onMounted(() => {
    if (props.nodeData && props.nodeData.properties) {

  }
    console.log('java组件已挂载');
    console.log('props.nodeData:', props.nodeData);
    console.log('form:', form);
  })

  watch(
  () => props.nodeData?.properties,
  (newProperties) => {
    if (newProperties) {
      form.appName = newProperties.appName || '';
      form.host = newProperties.host || '';
      form.jarSource = newProperties.jarSource || '';
      form.jarPath = newProperties.jarPath || '';
      form.bakPath = newProperties.bakPath || '';
      form.serverName = newProperties.serverName || '';
      form.port = newProperties.port || 0;
      form.healthUri = newProperties.healthUri || '';
      form.healthCheckTimeout = newProperties.healthCheckTimeout || 300
     }
  },
  { immediate: true }
); 
  // 提交处理方法
  function onSubmit() {
    console.log('submit!');
    console.log('form:', form);
    const nodeData = props.nodeData;
    const value = form.appName+"_"+form.host+"_"+form.serverName;
    console.log('lf:', props.lf);
    const nodeModel = props.lf.getNodeModelById(nodeData.id);
    nodeModel.text.value=value;
    nodeModel.updateText(value);

     nodeData.properties = { ...form }; // 更新 nodeData 的属性
    props.lf.setProperties(nodeData.id, { ...form }); // 确保传递的是 form 的最新值
    console.log('nodeData:', nodeData);
    emit('onClose');
   }
</script>
  
  