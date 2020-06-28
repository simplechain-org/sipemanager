<template>
  <div>
    <div style="padding: 20px;">
      <el-table :data="chains" style="width: 100%">
        <el-table-column prop="name" label="链的名称">
        </el-table-column>
        <el-table-column prop="network_id" label="网络编号">
        </el-table-column>
        <el-table-column prop="coin_name" label="币名">
        </el-table-column>
        <el-table-column prop="symbol" label="符号">
        </el-table-column>
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button @click="handleRemove(scope.$index)" type="text" size="small">删除</el-button>
            <el-button @click="handleContract(scope.$index)" type="text" size="small">设定跨链合约</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-container>
      <p style="color: red;">修改了跨链合约地址以后，务必修改锚定节点命令行参数--contract.main和--contract.sub的值，一定要匹配上当前设置的地址</p>
    </el-container>
    <el-dialog title="选择跨链合约地址" :visible.sync="dialogTableVisible">
      <el-table :data="contracts">
        <el-table-column prop="ID" label="选择" width="50px">
          <template slot-scope="scope">
            <el-radio v-model="contract_instance_id" :label="contracts[scope.$index].ID"></el-radio>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="地址"></el-table-column>
      </el-table>
      <el-button @click="modifyAddress()" type="primary" style="margin-top: 15px">确定使用该地址</el-button>
    </el-dialog>
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
    name: 'ChainList',
    data() {
      return {
        chain_id: 0,
        chains: [],
        contracts: [],
        contract_instance_id: 0,
        dialogTableVisible: false,
        centerDialogVisible: false,
        errMsg: ''
      }
    },
    created() {
      this.loadChainList()
    },
    methods: {
      loadChainList() {
        this.$http.get('/chain/list')
          .then(response => {
            if (response.data.code === 0) {
              this.chains = response.data.data
            }
          })
          .catch(error => {
            console.log(error)
          })
      },
      handleContract(index) {
        this.contract_id = 0
        this.chain_id = this.chains[index].ID
        this.$http.get('/contract/chain?chain_id=' + this.chains[index].ID)
          .then(response => {
            if (response.data.code === 0) {
              this.dialogTableVisible = true
              this.contracts = response.data.data
            }
          })
          .catch(error => {
            console.log(error)
          })
      },
      modifyAddress() {
        this.dialogTableVisible = false
        this.$http.post('/chain/address', {
            contract_instance_id: parseInt(this.contract_instance_id),
            chain_id: parseInt(this.chain_id)
          })
          .then(response => {
            if (response.data.code === 0) {
              this.centerDialogVisible = true
              this.errMsg = '跨链合约地址设定成功'
            } else {
              this.centerDialogVisible = true
              this.errMsg = response.data.msg
            }
          })
          .catch(error => {
            console.log(error)
          })
      },
      handleModify(index) {

      },
      // 链的删除
      handleRemove(index) {
        this.$confirm(
          '提示', {
            title: '提示',
            message: '确定删除该链吗?',
            showCancelButton: true,
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).then(() => {
          let chainId = this.chains[index].ID
          this.$http.delete('/chain/' + chainId)
            .then(response => {
              console.log(response.data)
              if (response.data.code === 0) {
                this.centerDialogVisible = true
                this.errMsg = response.data.data
                this.loadChainList()
              } else {
                this.centerDialogVisible = true
                this.errMsg = response.data.msg
              }
            })
            .catch(error => {
              console.log(error)
            })
        })
      },
      handleAdd() {
        this.$router.push({
          path: '/chain/add'
        })
      }
    }
  }
</script>
