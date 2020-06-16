<template>
  <div>
    <el-row>
      <el-col :span="20">
        <node-select></node-select>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="12" :offset="4">
        <el-row>
          <el-col :span="20">
            <el-input v-model="hash" placeholder="请输入交易哈希值"></el-input>
          </el-col>
          <el-col :span="4" style="padding-left: 20px;">
            <el-button type="primary" @click="search()">查询</el-button>
          </el-col>
        </el-row>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="20" :offset="2">
        <el-table :data="tableData" style="width: 100%" show-header="false" v-if="show">
          <el-table-column prop="name" width="200" label="字段">
          </el-table-column>
          <el-table-column prop="value" label="值">
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>
  </div>
</template>
<script>
  export default {
    name: 'SearchTransaction',
    data() {
      return {
        tableData: [],
        number: 0,
        show: false,
        hash: ''
      }
    },
    methods: {
      search() {
        this.$http.get('/transaction/' + this.hash)
          .then(response => {
            console.log(response)
            if (response.data.code === 0) {
              for (var attr in response.data.result) {
                this.tableData.push({
                  name: attr,
                  value: response.data.result[attr]
                })
              }
              this.show = true
            }
          })
          .catch(error => {
            console.log(error)
          })
      }
    }
  }
</script>
<style>
  .el-row {
    margin-bottom: 15px;
  }
</style>
