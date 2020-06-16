<template>
  <div>
    <div style="text-align: center;">
    <h3>添加钱包</h3>
    </div>
    <el-form ref="nodeForm" :model="form" label-width="150px" :rules="rules">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name"></el-input>
      </el-form-item>
      <el-form-item label="keystore文件内容" prop="content">
        <el-input v-model="form.content" type="textarea" placeholder="请输入内容">
        </el-input>
      </el-form-item>
    </el-form>
    <div style="text-align: center;">
      <el-button type="primary" @click="onSubmit()">立即导入</el-button>
      <el-button @click="handleReset()" type="warning">重&nbsp;&nbsp;&nbsp;&nbsp;置</el-button>
    </div>

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
    name: 'NodeList',
    data() {
      return {
        form: {
          name: '',
          content: ''
        },
        errMsg: '',
        centerDialogVisible: false,
        rules: {
          name: [{
            required: true,
            message: '名称不可为空',
            trigger: 'blur'
          }],
          content: [{
            required: true,
            message: '钱包keystore文件内容不可为空',
            trigger: 'blur'
          }]
        }
      }
    },
    methods: {
      handleReset() {
        this.form.name = ''
        this.form.content = ''
      },
      onSubmit() {
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            this.$http.post('/wallet', {
                name: this.form.name,
                content: this.form.content
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '添加成功'
                  this.$router.push({
                    path: '/wallet/list'
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
