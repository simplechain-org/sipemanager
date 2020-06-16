<template>
  <div>
    <el-row>
      <el-col :span="24">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column prop="name" label="名称" width="150">
          </el-table-column>
          <el-table-column prop="address" label="钱包地址">
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>
  </div>
</template>
<script>
  export default {
    name: 'WalletList',
    data() {
      return {
        tableData: [],
        number: 0
      }
    },
    created() {
      this.$http.get('/wallet/list')
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
      handleAdd() {
        this.$router.push({
          path: '/wallet/add'
        })
      }
    }
  }
</script>
