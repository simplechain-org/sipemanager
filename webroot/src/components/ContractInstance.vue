<template>
  <div>
    <el-table :data="instances" style="width: 100%" :border="true">
      <el-table-column prop="chain_id" label="上链编号" width="80">
      </el-table-column>
      <el-table-column prop="address" label="合约地址" width="360">
      </el-table-column>
      <el-table-column prop="tx_hash" label="交易哈希">
      </el-table-column>
    </el-table>
  </div>
</template>
<script>
  export default {
    name: 'ContractInstance',
    data() {
      return {
        instances: []
      }
    },
    created() {
      this.$http.get('/contract/instance/list')
        .then(response => {
          if (response.data.code === 0) {
            this.instances = response.data.result
          }
        })
        .catch(error => {
          console.log(error)
        })
    }
  }
</script>
