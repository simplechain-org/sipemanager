<template>
  <div>
        <!--要使用校验框架，1 必须在这里定义 ref和:rules-->
        <el-form ref="nodeForm" :model="form" label-width="100px" :rules="rules">
          <el-form-item label="合约描述" prop="description">
            <el-input v-model="form.description" type="text" placeholder="请输入合约描述">
            </el-input>
          </el-form-item>
          <!--要使用校验框架，2 在el-form-item中定义prop且名字要和对象的属性一致-->
          <el-form-item label="合约源码" prop="sol">
            <el-input v-model="form.sol" type="textarea" placeholder="请输入合约sol文件内容">
            </el-input>
          </el-form-item>
          <el-form-item label="合约abi" prop="abi">
            <el-input v-model="form.abi" type="textarea" placeholder="请输入合约abi文件内容">
            </el-input>
          </el-form-item>
          <el-form-item label="合约bin" prop="bin">
            <el-input v-model="form.bin" type="textarea" placeholder="请输入合约bin文件内容">
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
    name: 'ContractAdd',
    data() {
      return {
        form: {
          sol: '',
          abi: '',
          bin: '',
          description: ''
        },
        errMsg: '',
        centerDialogVisible: false,
        // 3、定义校验规则
        rules: {
          sol: [{
            required: true,
            message: '请输入合约的源码',
            trigger: 'blur'
          }],
          bin: [{
            required: true,
            message: '请输入合约的二进制内容',
            trigger: 'blur'
          }],
          abi: [{
            required: true,
            message: '请输入合约的abi',
            trigger: 'blur'
          }],
          description: [{
            required: true,
            message: '请输入合约描述',
            trigger: 'blur'
          }]
        }
      }
    },
    methods: {
      handleReset() {
        this.form.sol = ''
        this.form.bin = ''
        this.form.abi = ''
        this.form.description = ''
      },
      onSubmit(formName) {
        // 在el-form中定义的ref
        // 4、在提交方法中检查用户行为
        this.$refs.nodeForm.validate((valid) => {
          if (valid) {
            this.$http.post('/contract', {
                description: this.form.description,
                sol: this.form.sol,
                abi: this.form.abi,
                bin: this.form.bin
              })
              .then(response => {
                if (response.data.code === 0) {
                  this.centerDialogVisible = true
                  this.errMsg = '合约添加成功' + response.data.result
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
