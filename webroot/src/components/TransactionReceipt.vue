<template>
  <div>
    <el-row>
      <el-col :span="20">
        <node-select></node-select>
      </el-col>
    </el-row>
    <el-table :data="tableData" style="width: 100%" show-header="false">
      <el-table-column prop="name" width="200" label="字段">
      </el-table-column>
      <el-table-column prop="value" label="值">
      </el-table-column>
    </el-table>
  </div>
</template>
<script>
  export default {
    name: 'TransactionReceipt',
    data() {
      return {
        tableData: [],
        number: 0
      }
    },
    created() {
      this.number = this.$route.query.number
      this.$http.get('/transaction/' + this.$route.query.hash)
        .then(response => {
          console.log(response)
          if (response.data.code === 0) {
            for (var attr in response.data.data) {
              this.tableData.push({
                name: attr,
                value: response.data.data[attr]
              })
            }
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
      }
    }
  }
</script>
