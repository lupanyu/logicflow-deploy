<template>
    <div>
      <el-form :model="form" label-width="80px">
        <el-form-item label="请求 URL">
          <el-input v-model="form.url" placeholder="请输入请求 URL"></el-input>
        </el-form-item>
        <el-form-item label="请求方法">
          <el-select v-model="form.method" placeholder="请选择请求方法">
            <el-option label="GET" value="GET"></el-option>
            <el-option label="POST" value="POST"></el-option>
            <el-option label="PUT" value="PUT"></el-option>
            <el-option label="DELETE" value="DELETE"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="请求体">
          <el-input type="textarea" v-model="form.body" placeholder="请输入请求体"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">保存</el-button>
        </el-form-item>
      </el-form>
    </div>
  </template>
<script setup>
  import { ref,  onMounted,watch, reactive } from 'vue'
  
  // 使用 `defineProps` 来定义 props
  const props = defineProps(['nodeData', 'lf'])

  // 使用 `reactive` 来定义响应式数据
  const form = reactive({
    url: '222',
    method:  'GET',
    body:   '222'
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
    console.log('组件已挂载');
    console.log('props.nodeData:', props.nodeData);
    console.log('form:', form);
  })

  watch(
  () => props.nodeData?.properties,
  (newProperties) => {
    if (newProperties) {
      form.url = newProperties.url || '';
      form.method = newProperties.method || 'GET';
      form.body = newProperties.body || '';
    }
  },
  { immediate: true }
);
  // 提交处理方法
  function onSubmit() {
    console.log('submit!');
    console.log('form:', form);
    const nodeData = props.nodeData;
    nodeData.properties = { ...form }; // 更新 nodeData 的属性
    props.lf.setProperties(nodeData.id, { ...form }); // 确保传递的是 form 的最新值
    emit('onClose');
  }
</script>
  
  