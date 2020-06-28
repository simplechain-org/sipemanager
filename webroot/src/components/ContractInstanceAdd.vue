<template>
  <div>
    <!--要使用校验框架，1 必须在这里定义 ref和:rules-->
    <el-form ref="nodeForm" :model="form" label-width="100px" :rules="rules">
      <el-form-item label="所在链" prop="chain_id">
        <el-select v-model="form.chain_id" placeholder="请选择" style="width: 100%;">
          <el-option v-for="item in chains" :key="item.ID" :label="item.name" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="选择合约" prop="contract_id">
        <el-select v-model="form.contract_id" placeholder="请选择合约" style="width:100%">
          <el-option v-for="item in contracts" :key="item.ID" :label="item.description" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <!--要使用校验框架，2 在el-form-item中定义prop且名字要和对象的属性一致-->
      <el-form-item label="合约地址" prop="address">
        <el-input v-model="form.address" type="text" placeholder="请输入内容">
        </el-input>
      </el-form-item>
      <el-form-item label="交易哈希" prop="tx_hash">
        <el-input v-model="form.tx_hash" type="text" placeholder="请输入内容">
        </el-input>
      </el-form-item>
    </el-form>
    <div style="text-align: center;">
      <el-button type="primary" @click="onSubmit('nodeForm')">提交</el-button>
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
    name: 'ContractInstanceAdd',
    data() {
      return {
        contracts: [],
        form: {
          chain_id: 0,
          contract_id: '',
          address: '',
          tx_hash: ''
        },
        errMsg: '',
        centerDialogVisible: false,
        // 3、定义校验规则
        rules: {
          address: [{
            required: true,
            message: '请输入合约的地址',
            trigger: 'blur'
          }],
          tx_hash: [{
            required: true,
            message: '生成合约实例的交易哈希',
            trigger: 'blur'
          }]
        }
      }
    },
    created() {
      this.$http.get('contract/list')
        .then(response => {
          if (response.data.code === 0) {
            this.contracts = response.data.data
            if (this.contracts.length > 0) {
              this.form.contract_id = this.contracts[0].ID
            }
          }
        })
        .catch(error => {
          console.log(error)
        })

      this.$http.get('/chain/list').then(response => {
          if (response.data.code === 0) {
            this.chains = response.data.data
            if (this.chains.length > 0) {
              this.form.chain_id = this.chains[0].ID
            }
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    methods: {
      handleReset() {
        this.form.address = ''
        this.form.tx_hash = ''
      },
      onSubmit(formName) {
        // 在el-form中定义的ref
        // 4、在提交方法中检查用户行为
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            this.$http.post('/contract/instance/add', {
                chain_id: this.form.chain_id,
                contract_id: this.form.contract_id,
                tx_hash: this.form.tx_hash,
                address: this.form.address
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '合约实例添加成功' + response.data.data
                  this.handleReset()
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
