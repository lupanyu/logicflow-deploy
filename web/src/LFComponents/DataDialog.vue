<template>
  <div>
    <vue-json-pretty :path="'res'" :data="graphData"> </vue-json-pretty>
  </div>
</template>
<script setup>
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
 
const props = defineProps({
  graphData: Object
});
</script>
<style scoped>
</style>
