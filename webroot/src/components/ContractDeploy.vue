<template>
  <div>
   <!-- <el-row>
      <el-col :span="20"> -->
        <node-select></node-select>
   <!--   </el-col>
    </el-row> -->
    <!--要使用校验框架，1 必须在这里定义 ref和:rules-->
    <el-form ref="nodeForm" :model="form" label-width="100px" :rules="rules">
      <!--要使用校验框架，2 在el-form-item中定义prop且名字要和对象的属性一致-->
      <el-form-item label="选择合约" prop="contract_id">
        <el-select v-model="form.contract_id" placeholder="请选择合约" style="width:100%">
          <el-option v-for="item in contracts" :key="item.ID" :label="item.description" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="选择钱包" prop="wallet_id">
        <el-select v-model="form.wallet_id" placeholder="请选择钱包" style="width:100%">
          <el-option v-for="item in wallets" :key="item.id" :label="item.address" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="钱包密码" prop="password">
        <el-input v-model="form.password" show-password>></el-input>
      </el-form-item>
    </el-form>
    <div style="text-align: center;">
      <el-button type="primary" @click="onSubmit('nodeForm')">部署合约</el-button>
      <el-button @click="handleReset()" type="warning">重&nbsp;&nbsp;&nbsp;&nbsp;置</el-button>
    </div>

    <el-dialog title="提示" :visible.sync="centerDialogVisible" width="50%" center>
      <span>{{errMsg}}</span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="centerDialogVisible = false">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
  export default {
    name: 'ContractDeploy',
    data() {
      return {
        contracts: [],
        wallets: [],
        form: {
          password: '',
          wallet_id: 0,
          contract_id: 0
        },
        errMsg: '',
        centerDialogVisible: false,
        // 3、定义校验规则
        rules: {
          password: [{
            required: true,
            message: '请输入钱包密码',
            trigger: 'blur'
          }]
        }
      }
    },
    created() {
      this.$http.get('/wallet/list')
        .then(response => {
          if (response.data.code === 0) {
            this.wallets = response.data.result
            if (this.wallets.length > 0) {
              this.form.wallet_id = this.wallets[0].ID
            }
          }
        })
        .catch(error => {
          console.log(error)
        })
      this.$http.get('contract/list')
        .then(response => {
          if (response.data.code === 0) {
            this.contracts = response.data.result
            if (this.contracts.length > 0) {
              this.form.contract_id = this.contracts[0].ID
            }
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    methods: {
      handleReset() {
        this.form.password = ''
      },
      onSubmit(formName) {
        // 在el-form中定义的ref
        // 4、在提交方法中检查用户行为
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            this.$http.post('/contract/instance', {
                password: this.form.password,
                contract_id: this.form.contract_id,
                wallet_id: parseInt(this.form.wallet_id)
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '合约部署成功' + response.data.result
                  this.handleReset()
                  this.$router.push('/contract/instance')
                } else {
                  this.centerDialogVisible = true
                  this.errMsg = response.data.msg
                }
              })
              .catch(error => {
                console.log(error)
              })
          } // 校验失败的话，会在界面上有浮动提示，不需要额外做工作
        })
      } // onSubmit
    } // methods
  }
</script>
