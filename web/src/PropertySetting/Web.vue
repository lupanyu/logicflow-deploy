<template>
  <div class="app-config-form">
    <el-form :model="form" label-width="120px">
      <!-- 核心配置项 -->
      <el-form-item label="应用名称" prop="appName">
        <el-input 
          v-model="form.appName" 
          placeholder="请输入应用名称（英文）"
          clearable>
        </el-input>
      </el-form-item>

      <el-form-item label="主机名称" prop="host">
        <el-input
          v-model="form.host"
          placeholder="请输入主机名称"
          clearable>
        </el-input>
      </el-form-item>

      <el-form-item label="代码包地址" prop="packageSource">
        <el-input
          type="textarea"
          v-model="form.packageSource"
          placeholder="示例：/data/pkg/app-service-1.0.0.jar"
          :autosize="{ minRows: 2 }">
        </el-input>
      </el-form-item>

      <el-form-item label="代码部署位置" prop="deployPath">
        <el-input
          type="textarea"
          v-model="form.deployPath"
          placeholder="示例：/data/services/app-service/"
          :autosize="{ minRows: 2 }">
        </el-input>
      </el-form-item>

      <el-form-item label="备份路径" prop="bakPath">
        <el-input
          type="textarea"
          v-model="form.bakPath"
          placeholder="示例：/data/backup/app-service-{{timestamp}}.jar"
          :autosize="{ minRows: 2 }">
        </el-input>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="onSubmit">保存配置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
  import {   onMounted,watch, reactive } from 'vue'
  const props = defineProps(['nodeData', 'lf'])

const form = reactive({
  appName: '',
  host: '',
  packageSource: '',
  deployPath: '',
  bakPath: ''
})

 // 使用 `defineEmits` 来定义事件
 const emit = defineEmits(['onClose'])

// onMounted 钩子
onMounted(() => {
  if (props.nodeData && props.nodeData.properties) {
  form.url = props.nodeData.properties.url || '';
  form.method = props.nodeData.properties.method || 'GET';
  form.body = props.nodeData.properties.body || '';
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
  const value = form.appName+"_"+form.host;
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

<style scoped>
.app-config-form {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}
</style>
