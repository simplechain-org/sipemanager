<template>
  <div>
    <div style="text-align: center;">
    <h3>添加节点的基本信息</h3>
    </div>
    <el-form ref="nodeForm" :model="form" label-width="100px" :rules="rules">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name"></el-input>
      </el-form-item>
      <el-form-item label="地址" prop="address">
        <el-input v-model="form.address"></el-input>
      </el-form-item>
      <el-form-item label="端口" prop="port">
        <el-input v-model="form.port"></el-input>
      </el-form-item>
      <el-form-item label="接入链" prop="chain_id">
        <el-select v-model="form.chain_id" placeholder="请选择目标链" style="width: 100%;">
          <el-option v-for="item in chains" :key="item.ID" :label="item.name" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div style="text-align: center;">
      <el-button type="primary" @click="onSubmit()">立即创建</el-button>
      <el-button @click="handleReset()" type="warning">重&nbsp;&nbsp;&nbsp;&nbsp;置</el-button>
    </div>
    <el-dialog title="提示" :visible.sync="centerDialogVisible" width="60%" center>
      <span>{{errMsg}}</span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="centerDialogVisible = false">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
  export default {
    name: 'NodeList',
    data() {
      return {
        chains: [],
        form: {
          name: '',
          address: '127.0.0.1',
          port: 8545,
          chain_id: 1
        },
        errMsg: '',
        centerDialogVisible: false,
        rules: {
          name: [{
            required: true,
            message: '名称不可为空',
            trigger: 'blur'
          }],
          address: [{
            required: true,
            message: '地址不可为空',
            trigger: 'blur'
          }],
          port: [{
            required: true,
            message: '端口不可为空',
            trigger: 'blur'
          }],
          chain_id: [{
            required: true,
            message: '链编号不可为空',
            trigger: 'blur'
          }]
        }
      }
    },
    created() {
      this.$http.get('/chain/list')
        .then(response => {
          if (response.data.code === 0) {
            this.chains = response.data.result
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
        this.form.name = ''
        this.form.address = ''
        this.form.port = ''
      },
      onSubmit() {
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            this.$http.post('/node', {
                name: this.form.name,
                address: this.form.address,
                port: parseInt(this.form.port),
                chain_id: parseInt(this.form.chain_id)
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '添加成功'
                  this.$router.push({
                    path: '/node/list'
                  })
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
