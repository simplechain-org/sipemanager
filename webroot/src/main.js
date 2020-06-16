// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import axios from 'axios'
import NodeSelect from '@/components/common/NodeSelect'
import 'font-awesome/css/font-awesome.min.css'
import * as d3 from 'd3'
Vue.prototype.$d3 = d3
Vue.component('node-select', NodeSelect)
Vue.use(ElementUI)
Vue.config.productionTip = false
Vue.prototype.$http = axios
axios.defaults.baseURL = '/api/v1'
// 请求拦截器：在发送请求前拦截
axios.interceptors.request.use(config => {
  if (localStorage.getItem('accessToken')) {
    config.headers.Authorization = localStorage.getItem('accessToken')
  }
  return config
}, error => {
  return Promise.reject(error)
})

// 响应拦截器：在请求响应之后拦截
axios.interceptors.response.use(response => {
  if (response.data.code === 401) {
    router.push({
      path: '/login'
    })
  }
  return response
}, error => {
  return Promise.reject(error)
})

let vm = new Vue({
  router,
  el: '#app',
  render: h => h(App)
})

Vue.use({
  vm
})
