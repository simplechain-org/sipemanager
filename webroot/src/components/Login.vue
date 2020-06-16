<template>
  <div>
    <el-form ref="loginForm" :model="form" :rules="rules" label-width="50px" class="login-box">
      <h3 class="login-title">欢迎登录</h3>
      <el-form-item label="账号" prop="username">
        <el-input type="text" placeholder="请输入账号" v-model="form.username" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input type="password" placeholder="请输入密码" v-model="form.password" />
      </el-form-item>
      <div style="text-align: center;">
        <el-button type="primary" v-on:click="onSubmit()">登&nbsp;&nbsp;&nbsp;&nbsp;录</el-button>
        <el-button type="primary" v-on:click="onRegister()">注&nbsp;&nbsp;&nbsp;&nbsp;册</el-button>
      </div>
    </el-form>
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
    name: 'Login',
    data() {
      return {
        centerDialogVisible: false,
        errMsg: '',
        form: {
          username: '',
          password: ''
        },
        // 表单验证，需要在 el-form-item 元素中增加 prop 属性
        rules: {
          username: [{
            required: true,
            message: '账号不可为空',
            trigger: 'blur'
          }],
          password: [{
            required: true,
            message: '密码不可为空',
            trigger: 'blur'
          }]
        }
      }
    },
    methods: {
      onRegister() {
        this.$router.push({
          path: '/register'
        })
      },
      onSubmit() {
        // 为表单绑定验证功能
        this.$refs.loginForm.validate((valid) => {
          if (valid) {
            this.$http.post('/user/login', {
                username: this.form.username,
                password: this.form.password
              })
              .then(response => {
                if (response.data.code === 0) {
                  localStorage.setItem('accessToken', 'Bearer ' + response.data.result.token)
                  localStorage.setItem('user_id', response.data.result.user_id)
                  this.$router.push({
                    path: '/'
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

<style lang="css" scoped>
  .login-box {
    border: 1px solid #DCDFE6;
    width: 360px;
    margin: 180px auto;
    padding: 35px 35px 35px 35px;
    border-radius: 5px;
    -webkit-border-radius: 5px;
    -moz-border-radius: 5px;
    box-shadow: 0 0 25px #909399;
  }

  .login-title {
    text-align: center;
    margin: 0 auto 40px auto;
    color: #303133;
  }
</style>
