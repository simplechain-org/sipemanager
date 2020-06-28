<template>
    <el-table :data="contracts" style="width: 100%" :border="true">
      <el-table-column prop="ID" label="编号" width="80">
      </el-table-column>
      <el-table-column prop="description" label="合约描述">
      </el-table-column>
    </el-table>
</template>
<script>
  export default {
    name: 'ContractList',
    data() {
      return {
        contracts: []
      }
    },
    created() {
      this.$http.get('/contract/list')
        .then(response => {
          if (response.data.code === 0) {
            this.contracts = response.data.data
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    methods: {
      handleRemove(index) {
          alert(this.contracts[index].address)
      }
    }
  }
</script>
