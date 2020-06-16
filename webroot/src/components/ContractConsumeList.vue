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
  <!-- <el-table :data="contractConsumeList" style="width: 100%" :border="true">
    <el-table-column prop="ID" label="编号" width="80">
    </el-table-column>
    <el-table-column prop="source_chain_id" label="本链">
    </el-table-column>
    <el-table-column prop="target_chain_id" label="目标链">
    </el-table-column>
    <el-table-column prop="taker" label="钱包地址">
    </el-table-column>
    <el-table-column prop="source_value" label="本链资产">
    </el-table-column>
    <el-table-column prop="target_value" label="目标资产">
    </el-table-column>
    <el-table-column prop="status_text" label="交易状态">
    </el-table-column>
    <el-table-column prop="tx_hash" label="交易哈希">
    </el-table-column>
    <el-table-column prop="ctx_id" label="跨链哈希">
    </el-table-column>
  </el-table> -->

   <el-card v-for="item in contractConsumeList" :key="item.ID" class="item_card">
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
      <el-col span="2">钱包地址</el-col>
       <el-col span="22">{{item.taker}}</el-col>
    </el-row>
    <el-row class="item_row">
      <el-col span="2">本链资产</el-col>
       <el-col span="22">{{item.source_value}}</el-col>
    </el-row>
    <el-row class="item_row">
      <el-col span="2">目标资产</el-col>
       <el-col span="22">{{item.target_value}}</el-col>
    </el-row>
    <el-row class="item_row">
      <el-col span="2">交易状态</el-col>
       <el-col span="22">{{item.status_text}}</el-col>
    </el-row>
    <el-row class="item_row">
      <el-col span="2">交易哈希</el-col>
       <el-col span="22">{{item.tx_hash}}</el-col>
    </el-row>
    <el-row class="item_row">
      <el-col span="2">跨链哈希</el-col>
       <el-col span="22">{{item.ctx_id}}</el-col>
    </el-row>
  </el-card>
  </div>
</template>
<script>
  export default {
    name: 'ContractConsumeList',
    data() {
      return {
        contractConsumeList: []
      }
    },
    created() {
      this.$http.get('/contract/consume/list')
        .then(response => {
          if (response.data.code === 0) {
            this.contractConsumeList = response.data.result
          }
        })
        .catch(error => {
          console.log(error)
        })
    }
  }
</script>
