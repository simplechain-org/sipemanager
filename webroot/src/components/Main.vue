<template>
  <div>
    <p style="color: green;">{{localTime}}</p>
    <el-steps :active="8" align-center finish-status="success" direction="vertical">
      <el-step title="链的添加" description="添加链的基本信息,包括链的名称,networkid等"></el-step>
      <el-step title="添加节点" description="添加节点的基本信息,包括地址和端口,用于rpc连接"></el-step>
       <el-step title="添加钱包" description="添加钱包,为发起交易准备"></el-step>
      <el-step title="添加合约" description="保存合约到系统中"></el-step>
      <el-step title="部署合约" description="将合约部署到链上"></el-step>
      <el-step title="设定跨链合约实例" description="设定当前链所使用的跨链合约实例"></el-step>
      <el-step title="链的注册" description="调用合约注册跨链信息"></el-step>
      <el-step title="发起跨链交易" description="调用合约发起跨链交易"></el-step>
    </el-steps>
   <!-- <el-calendar v-model="value">
    </el-calendar> -->
  </div>
</template>
<script>
  function showLocale(objD) {
    var str
    var yy = objD.getYear()
    if (yy < 1900) yy = yy + 1900
    var MM = objD.getMonth() + 1
    if (MM < 10) MM = '0' + MM
    var dd = objD.getDate()
    if (dd < 10) dd = '0' + dd
    var hh = objD.getHours()
    if (hh < 10) hh = '0' + hh
    var mm = objD.getMinutes()
    if (mm < 10) mm = '0' + mm
    var ss = objD.getSeconds()
    if (ss < 10) ss = '0' + ss
    var ww = objD.getDay()
    if (ww === 0) ww = '星期日'
    if (ww === 1) ww = '星期一'
    if (ww === 2) ww = '星期二'
    if (ww === 3) ww = '星期三'
    if (ww === 4) ww = '星期四'
    if (ww === 5) ww = '星期五'
    if (ww === 6) ww = '星期六'
    str = yy + '年' + MM + '月' + dd + '日' + hh + ':' + mm + ':' + ss + ' ' + ww
    return str
  }
  export default {
    data() {
      return {
        value: new Date(),
        localTime: ''
      }
    },
    created() {
      this.tick(this)
    },
    methods: {
      tick(o) {
        var today
        today = new Date()
        o.localTime = showLocale(today)
        setTimeout(o.tick, 1000, o)
      }
    }
  }
</script>
<style>
</style>
