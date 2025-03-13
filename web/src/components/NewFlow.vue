<template>
    <div class="logic-flow-view">
      <h3 class="d-title">创建部署流程</h3>
      <!-- 辅助工具栏 -->
      <Control
        class="demo-control"
        v-if="lf"
        :lf="lf"
        @catData="catData"
      ></Control>
 
      <div id="LF-view" ref="container"></div>

      <div>
      <el-drawer
        title="设置节点属性"
        v-model="dialogVisible"  
        size="500px"
        :before-close="closeDialog"> 
        <PropertyDialog
          v-show="dialogVisible"
          :nodeData="clickNode"
          :lf="lf"
          :key="clickNode.value?.id"
          @setPropertiesFinish="closeDialog"
        >
        <p>这是测试打开抽屉显示的数据</p>
      </PropertyDialog>
      </el-drawer>
    </div>
      <!-- 数据查看面板 -->
      <el-dialog
        title="数据"
        width="50%"
        v-model="dataVisible"
        @close="dataVisible = false">
        <DataDialog v-if="dataVisible" :graphData="graphData"></DataDialog>
      </el-dialog>

    </div>
  </template>
<script setup>
  import LogicFlow from '@logicflow/core'
  // const LogicFlow = window.LogicFlow
  import { Menu, Snapshot, MiniMap,DndPanel } from '@logicflow/extension'
  import '@logicflow/core/lib/style/index.css'
  import '@logicflow/extension/lib/style/index.css'
  // import NodePanel from '../LFComponents/NodePanel.vue'
  // import AddPanel from '../LFComponents/AddPanel.vue'
  import Control from '@/LFComponents/Control.vue'  // 控制按钮
  import PropertyDialog from '@/PropertySetting/PropertyDialog.vue' // 点击节点后弹出的抽屉
  import DataDialog from '@/LFComponents/DataDialog.vue' // 数据弹框
  import { nodeList } from './config' // 预置的节点
  import { ref, onMounted, nextTick  } from 'vue'
  
  import {    registerStart,registerEnd, registerJava,registerWeb,registerJenkins,registerShell  } from '../LFComponents/nodes/'
  import { ElMessage } from 'element-plus'
  import edges from '../LFComponents/edges/index.js'

// const name = 'LF'
// 新增：存储当前选中的连线
const selectedEdge = ref(null)
// 新增：颜色映射
const colorMap = {
  'default': '#999999',
  'green': '#52c41a',
  'red': '#ff4d4f',
  'blue': '#1890ff'
}
const colorKeys = Object.keys(colorMap)

const lf = ref( null)
const container = ref(null)
const showAddPanel = ref(false)
const addPanelStyle = ref({
          top: 0,
          left: 0
        })
const nodeData = ref(null)
const addClickNode = ref(null)
const clickNode = ref(null)
const dialogVisible = ref(false)
const graphData = ref(null)
const dataVisible = ref(false)
const moveData = ref({})

 

function initLf () {
      // 画布配置
      LogicFlow.use(DndPanel);
      // LogicFlow.use(SelectionSelect);

      lf.value = new LogicFlow({

        background: {
            backgroundColor: '#f7f9ff',
        },
        grid: {
          size: 10,
          visible: true
        },
        keyboard: {
          enabled: true
        },
        edgeTextDraggable: true,
        hoverOutline: false,
      
          plugins: [
            Menu,
            MiniMap,
            DndPanel,
            Snapshot
          ],
          container: container.value,
      })
   
      lf.value.register(edges)
      lf.value.setDefaultEdgeType('myCurvedEdge')
      // 注册节点
      registerNode()

  }
  // 注册节点
  function registerNode () {
    const formattedNodeList = nodeList.map(node => {
      if (typeof node === 'string') {
        return { type: node, properties: {editable: true,
          refY:-100,
          textStyle: {
            textAnchor: 'middle',
            dominantBaseline: 'middle',
          },
        } }; // 确保是对象格式
      }
      return node; // 如果已经是对象，直接返回
    });
    // 注册节点
    registerStart(lf.value)
     registerEnd(lf.value)
    registerWeb(lf.value)
    registerJava(lf.value)
    registerJenkins(lf.value)
    registerShell(lf.value)
    // 注册节点到拖拽面板里
    lf.value.extension.dndPanel.setPatternItems(formattedNodeList)
    render()
  }
  // 渲染
  function render() {
    lf.value.render("")
     LfEvent()
  }
  // 动态更新边状态
  function updateEdgeStatus(edgeId, color) {
    const edgeModel = lf.value.getEdgeModelById(edgeId);
    edgeModel.setProperties({ color: color   });
    if (color === 'green') {
       lf.value.openEdgeAnimation(edgeId)
    }else{
      lf.value.closeEdgeAnimation(edgeId)
    }
    console.log(lf.value )
    console.log('lf.value.getEdgeModelById', lf.value.getEdgeModelById(edgeId));
    selectedEdge.value = edgeModel;
    }
  // 导出数据
  function getData () {
    const data = lf.value.getGraphData()
    console.log(JSON.stringify(data))
  }
  // 事件监听
  function LfEvent () {
    lf.value.on('node:click', ({data}) => {
      if (data.type === 'start' || data.type ==='end' ) {
        return
      }
      dialogVisible.value = false
      clickNode.value = data
      console.log('node:click', data,dialogVisible.value)
      nextTick(() => {
        dialogVisible.value = true
      })
      console.log('node:click--over', dialogVisible.value)
    })
 
    // 新增：监听连线点击事件
    lf.value.on('edge:click', ({data}) => {
      console.log('edge:click', data)
      selectedEdge.value = data
    })

    lf.value.on('element:click', () => {
      hideAddPanel()
    })
    lf.value.on('edge:add', ({data}) => {
      // 新增：更新边状态
      updateEdgeStatus(data.id, 'default')
      console.log('edge:add', data)
    })
    lf.value.on('node:mousemove', ({data}) => {
      console.log('node:mousemove')
      moveData.value = data
    })
    lf.value.on('blank:click', () => {
      hideAddPanel()
    })
    lf.value.on('connection:not-allowed', () => {
      const msg = '不允许连接'
      ElMessage.error(msg)
    })
    lf.value.on('node:mousemove', () => {
      console.log('on mousemove')
    })
  }
  // 点击加号事件
  function clickPlus (e, attributes) {
    e.stopPropagation()
    console.log('clickPlus', e, attributes)
    const { clientX, clientY } = e
    console.log(clientX, clientY)
    addPanelStyle.value.top = (clientY - 40) + 'px'
    addPanelStyle.value.left = clientX + 'px'
    showAddPanel.value = true
    addClickNode.value = attributes
  }
  // 鼠标按下事件
  function  mouseDownPlus (e, attributes) {
    e.stopPropagation()
    console.log('mouseDownPlus', e, attributes)
  }
  // 隐藏添加面板
  function  hideAddPanel () {
    showAddPanel.value = false
    addPanelStyle.value.top = 0
    addPanelStyle.value.left = 0
    addClickNode.value = null
  }
  // 关闭抽屉
  function closeDialog () {
    dialogVisible.value = false
  }
  /**
   * 获取LogicFlow图数据并显示数据对话框
   * 此函数用于获取LogicFlow实例的图数据，并将其赋值给graphData引用。
   * 同时，将dataVisible引用设置为true，以显示数据对话框。
   */
  function catData(){
    console.log('lf  cat data')
    graphData.value = lf.value.getGraphData();
    dataVisible.value = true;
    console.log('datavisable.value',dataVisible.value)
    console.log('graphData value',graphData.value)
  }
  function GetGraphData(){
    const data = lf.value.getGraphData()
    return data
  }

  defineExpose({
    GetGraphData
    })
  onMounted(()=> {
    initLf()
    })

</script>
<style>
  .logic-flow-view {
    height: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
  }
  .d-title{
    text-align: center;
    margin: 20px;
  }
  .demo-control{
    position: absolute;
    top: 50px;
    right: 50px;
    z-index: 2;
  }
  #LF-view{
    width: calc(100% - 100px);
    height: 80%;
    outline: none;
    margin-left: 50px;
  }
  .time-plus{
    cursor: pointer;
  }
  .add-panel {
    position: absolute;
    z-index: 11;
    background-color: white;
    padding: 10px 5px;
  }
  .el-drawer__body {
    height: 80%;
    overflow: auto;
    margin-top: -30px;
    z-index: 9999 !important;
  }
  .lf-dnd-shape {
  background-size: contain;
}
  @keyframes lf_animate_dash {
    to {
      stroke-dashoffset: 0;
    }
  }
  </style>  