<template>
  <div>
    <el-row style="height: 30px;">
      <el-col :span="20">
        <node-select v-on:change="onNodeChange"></node-select>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="24">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column prop="number" label="区块高度">
          </el-table-column>
          <el-table-column prop="miner" label="矿工地址" width="356">
          </el-table-column>
          <el-table-column prop="gasLimit" label="gasLimit">
          </el-table-column>
          <el-table-column prop="gasUsed" label="gasUsed">
          </el-table-column>
          <el-table-column prop="timestamp" label="出块时间">
          </el-table-column>
          <el-table-column prop="difficulty" label="区块难度">
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button @click="handleClick(scope.$index,tableData)" type="text" size="small" v-if="tableData[scope.$index].transactions>0">查看交易</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>
    <el-row style="margin-top: 10px;">
      <el-col :span="24">
          <el-pagination background layout="prev, pager, next" :total="total" :page-size="10" @current-change="handleChangePage"
            v-if="paginationVisible">
          </el-pagination>
      </el-col>
    </el-row>
    <el-dialog title="提示" :visible.sync="centerDialogVisible" width="30%" center>
      <span>{{errMsg}}</span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="centerDialogVisible = false">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
  export default {
    name: 'Block',
    data() {
      return {
        tableData: [],
        total: 100,
        errMsg: '',
        centerDialogVisible: false,
        paginationVisible: false
      }
    },
    created() {
      this.$http.get('/block/list')
        .then(response => {
          if (response.data.code === 0) {
            this.tableData = response.data.data.data
            this.total = response.data.data.total
            this.paginationVisible = true
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    methods: {
      handleClick(row, tableData) {
        this.$router.push({
          path: '/transaction/list',
          query: {
            number: tableData[row].number
          }
        })
      },
      handleChangePage(current) {
        this.$http.get('/block/list', {
            params: {
              currentPage: current
            }
          })
          .then(response => {
            if (response.data.code === 0) {
              this.tableData = response.data.data.data
              this.total = response.data.data.total
            }
          })
          .catch(error => {
            console.log(error)
          })
      },
      onNodeChange() {
        this.$http.get('/block/list')
          .then(response => {
            if (response.data.code === 0) {
              this.tableData = response.data.data.data
              this.total = response.data.data.total
              this.paginationVisible = true
            } else {
              this.errMsg = response.data.msg
              this.centerDialogVisible = true
              this.tableData = []
              this.paginationVisible = false
            }
          })
          .catch(error => {
            console.log(error)
          })
      }
    }
  }
</script>
