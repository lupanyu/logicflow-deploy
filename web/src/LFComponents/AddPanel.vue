<template>
  <el-tabs tab-position="left">
    <el-tab-pane label="添加动作">
      <div v-for="item in nodeList" :key="item.type">
        <el-button class="add-node-btn" type="primary" size="mini" @click="$_addNode(item)">{{item.label}}</el-button>
      </div>
    </el-tab-pane>
    <el-tab-pane label="添加组">
      <el-button class="add-node-btn" type="primary" size="mini" @click="$_addTempalte">模板</el-button>
    </el-tab-pane>
  </el-tabs>
</template>
<script>
import { defineComponent, ref, toRefs } from 'vue';

export default defineComponent({
  name: 'AddPanel',
  props: {
    nodeData: Object,
    lf: [Object, String] // 修正了 props 的类型定义
  },
  setup(props) {
    const { nodeData } = toRefs(props); // 使用 toRefs 解构 props
    const nodeList = ref([
      {
        type: 'user',
        label: '用户'
      },
      {
        type: 'push',
        label: '推送'
      }
    ]);

    const $_addNode = (item) => {
      const { lf } = props;
      const { id, x, y } = nodeData.value; // 使用 value 访问 ref 的值
      const nextNode = lf.addNode({
        type: item.type,
        x: x + 150,
        y: y + 150
      });
      const nextId = nextNode.id;
      lf.addEdge({ sourceNodeId: id, targetNodeId: nextId });
      emit('addNodeFinish'); // 使用 emit 触发事件
    };

    const $_addTempalte = () => {
      const { lf } = props;
      const { id, x, y } = nodeData.value;
      // ... 其余代码保持不变 ...
    };

    return {
      nodeList,
      $_addNode,
      $_addTempalte
    };
  }
});
</script>
<style scoped>
.add-node-btn{
  margin-bottom: 10px;
  margin-right: 20px;
}
</style>
