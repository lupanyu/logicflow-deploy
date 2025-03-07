<template>
  <div class="property-dialog">
    <Http v-if="nodeData.type === 'http'" :nodeData="nodeData" :lf="lf" @onClose="handleClose"/>
    <Java v-if="nodeData.type === 'java'" :nodeData="nodeData" :lf="lf" @onClose="handleClose"/>
    <Web v-if="nodeData.type === 'web'" :nodeData="nodeData" :lf="lf" @onClose="handleClose"/> 
    <Jenkins v-if="nodeData.type === 'jenkins'" :nodeData="nodeData" :lf="lf" @onClose="handleClose"/>
    <Shell v-if="nodeData.type === 'shell'" :nodeData="nodeData" :lf="lf" @onClose="handleClose"/>

   </div>
</template>
<script setup>
import {   onMounted } from 'vue';
// import CommonProperty from './CommonProperty.vue';
// import User from './User.vue';  
import Http from './Http.vue';  
import Java from './Java.vue';
import Web from './Web.vue';
import Jenkins from './Jenkins.vue';
import Shell from './Shell.vue';

const props = defineProps( {      
    nodeData: Object,
    lf: Object
  })
const emit = defineEmits(['setPropertiesFinish'])


function handleClose(){
  emit('setPropertiesFinish')
}
onMounted(() => {
  console.log('PropertyDialog props:', props.nodeData, props.lf)
})

</script>
<style>
.property-dialog{
  padding: 20px;
}
</style>
