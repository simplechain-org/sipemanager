<template>
<div>
    <node-select v-on:change="onNodeChange"></node-select>
    <el-card  v-for="(obj,index) in tableData" class="item_card" :key="index">
      <el-row class="item_row">
         <el-col :span="2">跨链哈希</el-col>
        <el-col :span="22">{{obj.ctx_id}}</el-col>
      </el-row>
      <el-row class="item_row">
        <el-col :span="2">交易哈希</el-col>
        <el-col :span="22">{{obj.tx_hash}}</el-col>
      </el-row>
      <el-row class="item_row">
        <el-col :span="2">源资产</el-col>
        <el-col :span="22">{{obj.value}}</el-col>
      </el-row>
      <el-row class="item_row">
        <el-col :span="2">目标资产</el-col>
        <el-col :span="22">{{obj.destination_value}}</el-col>
      </el-row>
      <el-row>
        <el-col :span="2" :offset="22">
          <el-button @click="handleClick(obj)" type="primary" size="medium">买入</el-button>
        </el-col>
      </el-row>
    </el-card>

    <el-dialog title="达成交易" :visible.sync="dialogFormVisible">
      <el-form :model="form">
        <el-form-item label="交易哈希" :label-width="formLabelWidth">
          {{form.ctx_id}}
        </el-form-item>
        <el-form-item label="选择钱包" :label-width="formLabelWidth">
          <el-select v-model="form.wallet_id" placeholder="请选择钱包" style="width: 100%;">
            <el-option v-for="item in wallets" :key="item.id" :label="item.address" :value="item.ID"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="钱包密码" :label-width="formLabelWidth">
          <el-input v-model="form.password" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleConfirm()">确定买入</el-button>
      </div>
    </el-dialog>

    <el-dialog title="提示" :visible.sync="centerDialogVisible" width="50%" center>
      <span>{{errMsg}}</span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="handleOk()">确 定</el-button>
      </span>
    </el-dialog>

  </div>
</template>
<script>
  export default {
    name: 'ContractTransaction',
    data() {
      return {
        wallets: [],
        centerDialogVisible: false,
        errMsg: '',
        formLabelWidth: '68px',
        dialogFormVisible: false,
        tableData: [],
        number: 0,
        form: {
          ctx_id: '',
          password: '',
          wallet_id: 0
        }
      }
    },
    created() {
      this.loadTransaction()
      this.$http.get('/wallet/list')
        .then(response => {
          if (response.data.code === 0) {
            this.wallets = response.data.data
            if (this.wallets.length > 0) {
              this.form.wallet_id = this.wallets[0].ID
            }
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    methods: {
      handleOk() {
        this.centerDialogVisible = false
        this.$router.push('/contract/consume/list')
      },
      onNodeChange() {
        this.loadTransaction()
      },
      loadTransaction() {
        this.tableData = []
        this.$http.get('/contract/transaction')
          .then(response => {
            if (response.data.code === 0) {
              this.tableData = response.data.data
            }
          })
          .catch(error => {
            console.log(error)
          })
      },
      handleClick(obj) {
        this.form.ctx_id = obj.ctx_id
        this.dialogFormVisible = true
      },
      handleConfirm() {
        // todo 数据校验
        this.dialogFormVisible = false
        this.$http.post('/contract/consume', {
            ctx_id: this.form.ctx_id,
            wallet_id: parseInt(this.form.wallet_id),
            password: this.form.password
          })
          .then(response => {
            if (response.data.code === 0) {
              this.centerDialogVisible = true
              this.errMsg = '买入成功' + response.data.data
            } else {
              this.centerDialogVisible = true
              this.errMsg = response.data.msg
            }
          })
          .catch(error => {
            console.log(error)
          })
      }
    }
  }
</script>
