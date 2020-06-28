<template>
  <div>
    <el-form ref="nodeForm" :model="form" label-width="150px" :rules="rules">
      <el-form-item label="源链" prop="source_chain_id">
        <el-select v-model="form.source_chain_id" placeholder="请选择" style="width: 100%;">
          <el-option v-for="item in chains" :key="item.ID" :label="item.name" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="目标链" prop="target_chain_id">
        <el-select v-model="form.target_chain_id" placeholder="请选择" style="width: 100%;">
          <el-option v-for="item in chains" :key="item.ID+2" :label="item.name" :value="item.ID">
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
      <el-form-item label="交易哈希" prop="tx_hash">
        <el-input v-model="form.tx_hash"></el-input>
      </el-form-item>
    </el-form>

    <div style="text-align: center;">
      <el-button type="primary" @click="onSubmit('nodeForm')">提交</el-button>
      <el-button type="primary" @click="addAnchors()">增加锚定节点地址</el-button>
    </div>

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
    name: 'RegisterChainAdd',
    data() {
      return {
        chains: [],
        form: {
          source_chain_id: '',
          target_chain_id: '',
          sign_confirm_count: '',
          anchors: [{
            value: '',
            key: Date.now()
          }],
          tx_hash: '',
          button: ''
        },
        errMsg: '',
        centerDialogVisible: false,
        rules: {
          sign_confirm_count: [{
            required: true,
            message: '请输入最小签名确认数',
            trigger: 'blur'
          }]
        }
      }
    },
    created() {
      this.getChain()
    },
    methods: {
      handleOk() {
        this.centerDialogVisible = false
        this.$router.push('/chain/register/list')
      },
      getChain() {
        this.$http.get('/chain/list')
          .then(response => {
            this.chains = response.data.data
          })
      },
      onSubmit() {
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            var anchors = []
            for (var i = 0; i < this.form.anchors.length; i++) {
              anchors.push(this.form.anchors[i].value)
            }
            if (this.form.source_chain_id === this.form.target_chain_id) {
              this.centerDialogVisible = true
              this.errMsg = '源链和目标链不能相同'
            }
            this.$http.post('/contract/register/add', {
                source_chain_id: parseInt(this.form.source_chain_id),
                target_chain_id: parseInt(this.form.target_chain_id),
                sign_confirm_count: parseInt(this.form.sign_confirm_count),
                anchor_addresses: anchors,
                tx_hash: this.form.tx_hash
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '链信息注册成功 ' + response.data.data
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
      }
    }
  }
</script>
