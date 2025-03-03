<template>
  <div>
     <el-button-group>
      <el-button type="default" size="small" @click="zoomIn">放大</el-button>
      <el-button type="default" size="small" @click="zoomOut">缩小</el-button>
      <el-button type="default" size="small" @click="zoomReset">大小适应</el-button>
      <el-button type="default" size="small" @click="translateRest">定位还原</el-button>
      <!-- <el-button type="default" size="small" @click="reset">还原(大小&定位)</el-button>
      <el-button type="default" size="small" @click="undo" :disabled="undoDisable">上一步(ctrl+z)</el-button>
      <el-button type="default" size="small" @click="redo" :disabled="redoDisable">下一步(ctrl+y)</el-button> -->
      <!-- <el-button type="default" size="small" @click="download">下载图片</el-button> -->
      <el-button type="default" size="small" @click="catData">查看数据</el-button>
      <!-- <el-button v-if="catTurboData" type="default" size="small" @click="catTurboData">查看turbo数据</el-button> -->
      <!-- <el-button type="default" size="small" @click="showMiniMap">查看缩略图</el-button> -->
    </el-button-group>
  </div>
</template>

<script>
import { defineComponent, ref, onMounted, onBeforeUnmount } from 'vue';

export default defineComponent({
  name: 'Control',
  props: {
    lf: {
      type: Object,
      required: true
    },
    catTurboData: {
      type: Boolean,
      default: false
    }
  },
  setup(props, { emit }) {
    const undoDisable = ref(true);
    const redoDisable = ref(true);

    const historyChangeHandler = ({ data: { undoAble, redoAble } }) => {
      console.log("undoAble:", !undoAble, "redoAble", !redoAble);
      undoDisable.value = !undoAble;
      redoDisable.value = !redoAble;
    };

    onMounted(() => {
      props.lf.on('history:change', historyChangeHandler);
    });

    onBeforeUnmount(() => {
      props.lf.off('history:change', historyChangeHandler);
    });

    const zoomIn = () => {
      props.lf.zoom(true);
    };
    const zoomOut = () => {
      props.lf.zoom(false);
    };
    const zoomReset = () => {
      console.log(props)
      props.lf.resetZoom();
    };
    const translateRest = () => {
      props.lf.resetTranslate();
    };
    const reset = () => {
      props.lf.resetZoom();
      props.lf.resetTranslate();
    };
    const undo = () => {
      props.lf.undo();
    };
    const redo = () => {
      props.lf.redo();
    };
    const download = () => {
      props.lf.getSnapshot();
    };
    const catData = () => {
      console.log('Control.vue catData')
      emit('catData');
    };
    const catTurboData = () => {
      emit('catTurboData');  // 原样返回给父组件来处理
    };
    const showMiniMap = () => {
      const { lf } = props;
      lf.extension.miniMap.show(lf.graphModel.width - 150, 40);
    };

    return {
      undoDisable,
      redoDisable,
      zoomIn,
      zoomOut,
      zoomReset,
      translateRest,
      reset,
      undo,
      redo,
      download,
      catData,
      catTurboData,
      showMiniMap
    };
  }
});
</script>
<style scoped>
</style>
