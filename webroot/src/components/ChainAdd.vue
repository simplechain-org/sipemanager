<template>
  <div>
    <div style="text-align: center;">
    <h3>添加链的基本信息</h3>
    </div>
    <el-form ref="nodeForm" :model="form" label-width="150px" :rules="rules">
      <el-form-item label="链的名称" prop="name">
        <el-input v-model="form.name"></el-input>
      </el-form-item>
      <el-form-item label="网络编号" prop="network_id">
        <el-input v-model="form.network_id"></el-input>
      </el-form-item>
      <el-form-item label="币名" prop="coin_name">
        <el-input v-model="form.coin_name"></el-input>
      </el-form-item>
      <el-form-item label="符号" prop="symbol">
        <el-input v-model="form.symbol"></el-input>
      </el-form-item>
    </el-form>
    <div style="text-align: center;">
      <el-button @click="handleAdd" type="primary">立即创建</el-button>
    </div>
    <el-dialog title="提示" :visible.sync="centerDialogVisible" width="50%" center>
      <span>{{errMsg}}</span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="ok()">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
  export default {
    name: 'ChainAdd',
    data() {
      return {
        centerDialogVisible: false,
        errMsg: '',
        form: {
          name: '',
          network_id: '',
          coin_name: '',
          symbol: ''
        },
        rules: {
          name: [{
            required: true,
            message: '名称不可为空',
            trigger: 'blur'
          }],
          coin_name: [{
            required: true,
            message: '币的名称不可为空',
            trigger: 'blur'
          }],
          symbol: [{
            required: true,
            message: '币的符号不可为空',
            trigger: 'blur'
          }],
          network_id: [{
            required: true,
            message: '网络编号不可为空',
            trigger: 'blur'
          }]
        }
      }
    },
    methods: {
      ok() {
        this.centerDialogVisible = false
        this.$router.push('/chain/list')
      },
      handleAdd() {
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            this.$http.post('/chain/create', {
                name: this.form.name,
                network_id: parseInt(this.form.network_id),
                coin_name: this.form.coin_name,
                symbol: this.form.symbol
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '链信息添加成功 '
                  this.form.name = ''
                  this.form.network_id = ''
                  this.form.coin_name = ''
                  this.form.symbol = ''
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
