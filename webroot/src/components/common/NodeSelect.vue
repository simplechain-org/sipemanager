<template>
  <div>
    <el-form label-width="100px">
      <el-form-item label="连接节点" prop="nodeId">
        <el-select v-model="nodeId" placeholder="请选择" style="width: 100%;" @change="handleCommand">
          <el-option v-for="item in nodes" :key="item.ID" :label="item.description" :value="item.ID">
          </el-option>
        </el-select>
      </el-form-item>
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
    name: 'NodeSelect',
    data() {
      return {
        chain_id: '',
        custom: '',
        nodes: [],
        nodeId: '',
        errMsg: '',
        centerDialogVisible: false
      }
    },
    methods: {
      handleCommand(nodeId) {
        var userId = localStorage.getItem('user_id')
        this.$http.post('/node/change', {
            user_id: parseInt(userId),
            node_id: parseInt(nodeId)
          })
          .then(response => {
            if (response.data.code === 0) {
              this.nodeId = nodeId
              this.$emit('change')
            } else {
              this.centerDialogVisible = true
              this.errMsg = response.data.msg
            }
          })
          .catch(error => {
            console.log(error)
          })
      }
    },
    created() {
      var r1 = this.$http.get('/node/list')
      var r2 = this.$http.get('/node/current')
      this.$http.all([r1, r2])
        .then(this.$http.spread((res1, res2) => {
          this.nodes = res1.data.data
          var node = res2.data.data
          this.nodeId = node.ID
        }))
      // // 请求用户所有的节点
      // this.$http.get('/node/list')
      //   .then(response => {
      //     if (response.data.code === 0) {
      //       this.nodes = response.data.data
      //     }
      //   })
      //   .catch(error => {
      //     console.log(error)
      //   })
      // // 请求用户当前使用的节点
      // this.$http.get('/node/current')
      //   .then(response => {
      //     if (response.data.code === 0) {
      //       var node = response.data.data
      //       this.nodeId = node.ID
      //     }
      //   })
      //   .catch(error => {
      //     console.log(error)
      //   })
    }
  }
</script>
