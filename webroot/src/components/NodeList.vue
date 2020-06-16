<template>
  <div>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="name" label="名称">
      </el-table-column>
      <el-table-column prop="address" label="IP地址">
      </el-table-column>
      <el-table-column prop="port" label="端口">
      </el-table-column>
      <el-table-column prop="chain_name" label="接入链">
      </el-table-column>
    </el-table>
  </div>
</template>
<script>
  export default {
    name: 'NodeList',
    data() {
      return {
        tableData: [],
        number: 0
      }
    },
    created() {
      this.$http.get('/node/list')
        .then(response => {
          if (response.data.code === 0) {
            this.tableData = response.data.result
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    methods: {
      handleGoback(number) {
        this.$router.push({
          name: 'TransactionList',
          query: {
            number: number
          }
        })
      },
      handleAdd() {
        this.$router.push({
          path: '/node/add'
        })
      }
    }
  }
</script>
<style>
  .node-list{
    margin: 10px 5px;
    padding: 35px 35px 35px 35px;
    border-radius: 5px;
    -webkit-border-radius: 5px;
    -moz-border-radius: 5px;
    box-shadow: 0 0 25px #909399;
  }
</style>
