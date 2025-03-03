<template>
    <div>
      <el-form :model="form" label-width="80px">
        <el-form-item label="服务">
            <el-select v-model="form.service" placeholder="请选择服务">
                <el-option v-for="service in serviceList" :key="service.id" :label="service.name" :value="service.id"></el-option>
            </el-select>

        </el-form-item>
        <el-form-item label="参数配置">
          <el-input v-model="form.params" placeholder="请输入参数配置"></el-input>
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
    service: '',
    params:  ''
   })

  // 使用 `defineEmits` 来定义事件
  const emit = defineEmits(['onClose'])

  // 通过api来获取服务类型的数组
  const serviceList = ref([]);

  // 定义一个函数来获取服务类型列表
  function fetchServiceList() {
    // 使用全局变量来拼接URL
    const url = apiUrl + '/service/list';
    try {
        // 使用 fetch 来获取数据
            fetch(url)
            .then(response => response.json())
            .then(data => {
                // 假设数据是一个数组，你可以在这里处理它
                serviceList.value = data;
            })
      // 假设你有一个名为 fetchServiceList 的函数来获取服务类型列表
      fetchServiceList().then((data) => {
        serviceList.value = data;
      });
    }
    catch (error) {
      console.error('Error fetching service list:', error); 
    }
  }

  // onMounted 钩子
  onMounted(() => {
    fetchServiceList();

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
  
  