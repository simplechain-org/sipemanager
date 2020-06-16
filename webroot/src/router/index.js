import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Block from '@/components/Block'
import TransactionList from '@/components/TransactionList'
import TransactionReceipt from '@/components/TransactionReceipt'
import Login from '@/components/Login'
import NodeList from '@/components/NodeList'
import NodeAdd from '@/components/NodeAdd'
import WalletAdd from '@/components/WalletAdd'
import WalletList from '@/components/WalletList'
import ContractDeploy from '@/components/ContractDeploy'
import SearchTransaction from '@/components/SearchTransaction'
import RegisterChain from '@/components/RegisterChain'
import ContractProduce from '@/components/ContractProduce'
import ContractTransaction from '@/components/ContractTransaction'
import ChainList from '@/components/ChainList'
import ChainAdd from '@/components/ChainAdd'
import ContractList from '@/components/ContractList'
import Register from '@/components/Register'
import ContractAdd from '@/components/ContractAdd'
import ContractInstance from '@/components/ContractInstance'
import ContractInstanceAdd from '@/components/ContractInstanceAdd'
import Main from '@/components/Main'
import ContractProduceList from '@/components/ContractProduceList'
import ContractConsumeList from '@/components/ContractConsumeList'
import RegisterChainList from '@/components/RegisterChainList'
import RegisterChainAdd from '@/components/RegisterChainAdd'
Vue.use(Router)

export default new Router({
  routes: [{
      path: '/login',
      name: 'Login',
      component: Login,
      hidden: true
    },
    {
      path: '/register',
      name: 'Register',
      component: Register,
      hidden: true
    },
    {
      path: '/',
      component: Home,
      name: '链的管理',
      iconCls: 'el-icon-message',
      children: [{
        path: '/',
        component: Main,
        hidden: true
      }, {
        path: '/chain/add',
        component: ChainAdd,
        name: '链的添加'
      }, {
        path: '/chain/list',
        component: ChainList,
        name: '区块链表'
      }]
    },
    {
      path: '/',
      component: Home,
      name: '节点管理',
      iconCls: 'el-icon-message',
      children: [{
          path: '/node/list',
          name: '节点列表',
          component: NodeList
        },
        {
          path: '/node/add',
          name: '添加节点',
          component: NodeAdd
        }
      ]
    },
    {
      path: '/',
      component: Home,
      name: '合约管理',
      iconCls: 'el-icon-message',
      children: [{
          path: '/contract/add',
          name: '添加合约',
          component: ContractAdd
        }, {
          path: '/contract/list',
          name: '合约列表',
          component: ContractList
        },
        {
          path: '/contract/deploy',
          name: '合约部署',
          component: ContractDeploy
        },
        {
          path: '/contract/instance',
          name: '合约实例',
          component: ContractInstance
        },

        {
          path: '/contract/instance/add',
          name: '添加实例',
          component: ContractInstanceAdd
        }
      ]
    },
    {
      path: '/',
      component: Home,
      name: '钱包管理',
      iconCls: 'el-icon-message',
      children: [{
          path: '/wallet/add',
          name: '钱包导入',
          component: WalletAdd
        },
        {
          path: '/wallet/list',
          name: '钱包列表',
          component: WalletList
        }
      ]
    },
    {
      path: '/',
      component: Home,
      name: '跨链交易',
      iconCls: 'el-icon-message',
      children: [{
          path: '/chain/register',
          name: '链的注册',
          component: RegisterChain
        },
        {
          path: '/chain/register/list',
          name: '注册日志',
          component: RegisterChainList
        },
        {
          path: '/chain/register/add',
          name: '增加已有注册日志',
          component: RegisterChainAdd,
          hidden: true
        },
        {
          path: '/contract/produce',
          name: '发起跨链交易',
          component: ContractProduce
        },
        {
          path: '/contract/produce/list',
          name: '发单日志',
          component: ContractProduceList
        },
        {
          path: '/contract/transaction',
          name: '跨链交易列表',
          component: ContractTransaction
        },
        {
          path: '/contract/consume/list',
          name: '接单日志',
          component: ContractConsumeList
        }
      ]
    },
    {
      path: '/',
      component: Home,
      name: '区块浏览器',
      iconCls: 'el-icon-message',
      children: [{
          path: '/block',
          name: '区块列表',
          component: Block
        },
        {
          path: '/transaction/list',
          name: '交易列表',
          component: TransactionList,
          hidden: true
        },
        {
          path: '/transaction/receipt',
          name: '交易凭证',
          component: TransactionReceipt,
          hidden: true
        },
        {
          path: '/transaction/search',
          name: '交易查询',
          component: SearchTransaction
        }
      ]
    }
  ]
})
