<style>
  .item_row{
    padding: 5px;
  }
  .item_card{
    margin-bottom: 10px;
  }
</style>
<template>
  <div>
    <el-card v-for="item in registerChainList" :key="item.ID" class="item_card">
      <el-row class="item_row">
        <el-col span="2">编号</el-col>
         <el-col span="22">{{item.ID}}</el-col>
      </el-row>
      <el-row class="item_row">
        <el-col span="2">本链</el-col>
         <el-col span="22">{{item.source_chain_id}}</el-col>
      </el-row>
      <el-row class="item_row">
        <el-col span="2">目标链</el-col>
         <el-col span="22">{{item.target_chain_id}}</el-col>
      </el-row>
      <el-row class="item_row">
        <el-col span="2">最少签名数</el-col>
         <el-col span="22">{{item.confirm}}</el-col>
      </el-row>

      <el-row class="item_row" v-for="(address,index) in item.anchor_addresses.split(',')" :key="index">
        <el-col span="2">锚定地址{{index}}</el-col>
         <el-col span="22">{{address}}</el-col>
      </el-row>

      <el-row class="item_row">
        <el-col span="2">交易状态</el-col>
         <el-col span="22">{{item.status_text}}</el-col>
      </el-row>
      <el-row class="item_row">
        <el-col span="2">交易哈希</el-col>
         <el-col span="22">{{item.tx_hash}}</el-col>
      </el-row>
    </el-card>
    <router-link to="/chain/register/add">
      <el-button @click="handleClick(obj)" type="primary" size="medium">导入现有注册日志</el-button>
    </router-link>
  </div>

<!--  <div>
 <el-table :data="registerChainList" style="width: 100%" :border="true">
    <el-table-column prop="ID" label="编号" width="80">
    </el-table-column>
    <el-table-column prop="source_chain_id" label="本链">
    </el-table-column>
    <el-table-column prop="target_chain_id" label="目标链">
    </el-table-column>
    <el-table-column prop="confirm" label="最少签名数">
    </el-table-column>
    <el-table-column prop="anchor_addresses" label="锚定节点地址">
    </el-table-column>
    <el-table-column prop="status_text" label="交易状态">
    </el-table-column>
    <el-table-column prop="tx_hash" label="交易哈希">
    </el-table-column>
  </el-table>
  <router-link to="/chain/register/add">
    导入现有注册日志
  </router-link>
  </div> -->
</template>
<script>
  export default {
    name: 'RegisterChainList',
    data() {
      return {
        registerChainList: []
      }
    },
    created() {
      this.$http.get('/contract/register/list')
        .then(response => {
          if (response.data.code === 0) {
            this.registerChainList = response.data.result
          }
        })
        .catch(error => {
          console.log(error)
        })
    }
  }
</script>
