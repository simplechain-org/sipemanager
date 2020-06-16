<template>
  <div>
    <node-select v-on:change="onNodeChange"></node-select>
    <el-form ref="nodeForm" :model="form" label-width="150px" :rules="rules">
      <el-form-item label="目标链" prop="target_chain_id">
        <el-select v-model="form.target_chain_id" placeholder="请选择" style="width: 100%;">
          <el-option v-for="item in chains" :key="item.ID" :label="item.name" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="签名确认数" prop="sign_confirm_count">
        <el-input v-model="form.sign_confirm_count"></el-input>
      </el-form-item>
      <el-form-item v-for="(anchor, index) in form.anchors" :label="'锚定节点地址' + index" :key="anchor.key">
        <el-row>
          <el-col :span="22">
            <el-input v-model="form.anchors[index].value"></el-input>
          </el-col>
          <el-col :span="2" style="padding-left: 7px;">
            <el-button @click.prevent="removeDomain(anchor)" :disabled="form.anchors.length==1">删除</el-button>
          </el-col>
        </el-row>
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
      <el-button type="primary" @click="onSubmit('nodeForm')">立即注册</el-button>
      <el-button type="primary" @click="addAnchors()">增加锚定节点地址</el-button>
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
    name: 'RegisterChain',
    data() {
      return {
        wallets: [],
        chains: [],
        form: {
          password: '',
          target_chain_id: 1,
          sign_confirm_count: 1,
          anchors: [{
            value: '',
            key: Date.now()
          }],
          wallet_id: 0,
          button: 0
        },
        errMsg: '',
        centerDialogVisible: false,
        rules: {
          sign_confirm_count: [{
            required: true,
            message: '请输入最少确认数',
            trigger: 'blur'
          }],
          password: [{
            required: true,
            message: '请输入钱包密码',
            trigger: 'blur'
          }]
        }
      }
    },
    created() {
      this.getChain()
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
    },
    methods: {
      onNodeChange() {
        this.getChain()
      },
      getChain() {
        var r1 = this.$http.get('/chain/current')
        var r2 = this.$http.get('/chain/list')
        this.$http.all([r1, r2])
          .then(this.$http.spread((res1, res2) => {
            var node = res1.data.result
            var chains = res2.data.result
            this.chains = []
            for (let i = 0; i < chains.length; i++) {
              if (chains[i].ID !== node.chain_id) {
                this.chains.push(chains[i])
              }
            }
            if (this.chains.length > 0) {
              this.form.target_chain_id = this.chains[0].ID
            }
          }))
      },
      onSubmit() {
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            var anchors = []
            for (var i = 0; i < this.form.anchors.length; i++) {
              anchors.push(this.form.anchors[i].value)
            }
            this.$http.post('/contract/register', {
                target_chain_id: parseInt(this.form.target_chain_id),
                sign_confirm_count: parseInt(this.form.sign_confirm_count),
                anchor_addresses: anchors,
                wallet_id: parseInt(this.form.wallet_id),
                password: this.form.password
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '链信息注册成功 ' + response.data.result
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
      },
      removeDomain(item) {
        var index = this.form.anchors.indexOf(item)
        if (index !== -1 && this.form.anchors.length > 1) {
          this.form.anchors.splice(index, 1)
        }
      },
      addAnchors() {
        this.form.anchors.push({
          value: '',
          key: Date.now()
        })
      },
      handleReset() {

      }
    }
  }
</script>
