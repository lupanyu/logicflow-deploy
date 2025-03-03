<template>
  <div class="common-layout">
    <el-container>
      <el-header>Header</el-header>
      <el-container>
        <el-aside width="200px">Aside</el-aside>
        <el-main>
          <div id="app">
            <router-view/>
          </div>      
        </el-main>  
      </el-container>
    </el-container>
  </div>
</template>


<script>
export default {
  name: 'App',
}
</script>

<style>
</style>
