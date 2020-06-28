<template>
  <div>
    <el-row>
      <el-col :span="20">
        <node-select></node-select>
      </el-col>
    </el-row>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="from" label="from" width="400">
      </el-table-column>
      <el-table-column prop="to" label="to" width="400">
      </el-table-column>
      <el-table-column prop="value" label="value">
      </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button @click="handleClick(scope.$index,tableData,number)" type="text" size="small">查看详情</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
<script>
  export default {
    name: 'TransactionList',
    data() {
      return {
        tableData: [],
        number: 0
      }
    },
    created() {
      this.number = this.$route.query.number
      this.$http.get('/block/transaction/' + this.$route.query.number)
        .then(response => {
          console.log(response)
          if (response.data.code === 0) {
            this.tableData = response.data.data
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    methods: {
      handleClick(row, tableData, number) {
        this.$router.push({
          path: '/transaction/receipt',
          query: {
            hash: tableData[row].hash,
            number: number
          }
        })
      }
    }
  }
</script>
