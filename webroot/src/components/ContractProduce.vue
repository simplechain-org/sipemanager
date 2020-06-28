<template>
  <div>
    <node-select v-on:change="onNodeChange"></node-select>
    <el-form ref="nodeForm" :model="form" label-width="100px" :rules="rules">
      <el-form-item label="目标链" prop="chain_id">
        <el-select v-model="form.chain_id" placeholder="请选择目标链" style="width: 100%;" @change="chainChange()">
          <el-option v-for="item in chains" :key="item.ID" :label="item.name" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="source_coin+'的数量'" prop="source_value">
        <el-input v-model="form.source_value"></el-input>
      </el-form-item>
      <el-form-item :label="target_coin+'的数量'" prop="target_value">
        <el-input v-model="form.target_value"></el-input>
      </el-form-item>
      <el-form-item label="选择钱包" prop="wallet_id">
        <el-select v-model="form.wallet_id" placeholder="请选择" style="width: 100%;">
          <el-option v-for="item in wallets" :key="item.ID" :label="item.address" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="钱包密码" prop="password">
        <el-input v-model="form.password" show-password></el-input>
      </el-form-item>

    </el-form>
    <div style="text-align: center;">
      <el-button type="primary" @click="onSubmit()">立即创建</el-button>
      <el-button @click="handleReset()" type="warning">重&nbsp;&nbsp;&nbsp;&nbsp;置</el-button>
    </div>
    <el-dialog title="提示" :visible.sync="centerDialogVisible" width="60%" center>
      <span>{{errMsg}}</span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="handleOk()">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
  export default {
    name: 'ContractProduce',
    data() {
      return {
        source_chain_id: 0,
        target_chain_id: 0,
        source_coin: '',
        target_coin: '',
        wallets: [],
        chains: [],
        form: {
          password: '',
          source_value: '',
          target_value: '',
          extra: '',
          chain_id: 1,
          wallet_id: 0
        },
        errMsg: '',
        centerDialogVisible: false,
        rules: {
          password: [{
            required: true,
            message: '请输入钱包密码',
            trigger: 'blur'
          }],
          source_value: [{
            required: true,
            message: '请输入源值',
            trigger: 'blur'
          }],
          target_value: [{
            required: true,
            message: '请输入目标值',
            trigger: 'blur'
          }]
        }
      }
    },
    created() {
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
      this.getChain()
    },
    methods: {
      handleOk() {
        this.centerDialogVisible = false
        this.$router.push('/contract/produce/list')
      },
      chainChange(data) {
        this.target_chain_id = this.form.chain_id
        this.getTargetCoin()
      },
      getSourceCoin() {
        this.$http.get('/chain/info/' + this.source_chain_id)
          .then(response => {
            if (response.data.code === 0) {
              this.source_coin = response.data.data.symbol
            }
          })
          .catch(error => {
            console.log(error)
          })
      },
      getTargetCoin() {
        this.$http.get('/chain/info/' + this.target_chain_id)
          .then(response => {
            if (response.data.code === 0) {
              this.target_coin = response.data.data.symbol
            }
          })
          .catch(error => {
            console.log(error)
          })
      },
      onNodeChange() {
        this.getChain()
      },
      getChain() {
        var r1 = this.$http.get('/chain/current')
        var r2 = this.$http.get('/chain/list')
        this.$http.all([r1, r2])
          .then(this.$http.spread((res1, res2) => {
            var node = res1.data.data
            this.source_chain_id = node.chain_id
            this.getSourceCoin()
            var chains = res2.data.data
            this.chains = []
            for (let i = 0; i < chains.length; i++) {
              if (chains[i].ID !== node.chain_id) {
                this.chains.push(chains[i])
              }
            }
            if (this.chains.length > 0) {
              this.form.chain_id = this.chains[0].ID
              this.target_chain_id = this.form.chain_id
              this.getTargetCoin()
            }
          }))
      },
      handleReset() {
        this.form.password = ''
        this.form.source_value = ''
        this.form.target_value = ''
        this.form.extra = ''
      },
      onSubmit() {
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            this.$http.post('/contract/produce', {
                chain_id: parseInt(this.form.chain_id),
                source_value: parseInt(this.form.source_value),
                target_value: parseInt(this.form.target_value),
                wallet_id: parseInt(this.form.wallet_id),
                password: this.form.password,
                extra: this.form.extra
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '跨链交易创建成功 ' + response.data.data
                  this.handleReset()
                } else {
                  this.centerDialogVisible = true
                  this.errMsg = response.data.msg
                }
              })
              .catch(error => {
                console.log(error)
              })
          }
        })
      }
    }
  }
</script>
