<template>
  <div>
    <el-container>
      <svg version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink' xml:space='preserve'
        width='960' height='500'> </svg>
    </el-container>
    <el-button @click='add()' type='primary'>添加节点</el-button>
  </div>
</template>
<script>
  // https://blog.csdn.net/qq_34414916/article/details/80026029
  // https://www.npmjs.com/package/vue-native-websocket
  // npm install vue-native-websocket --save
  // vue-native-websocket
  import * as d3 from 'd3'
  // import {
  //   BigNumber
  // } from 'bignumber.js'
  export default {
    name: 'ChainView',
    data() {
      return {
        container: null,
        nodes: [],
        formatDifficulty: false,
        gs: null,
        links: null,
        linksText: null,
        forceSimulation: null,
        edges: [],
        g: null,
        circleTexts: [],
        circles: [],
        count: 8
      }
    },
    mounted() {
      var marge = {
        top: 60,
        bottom: 60,
        left: 60,
        right: 60
      }
      var svg = d3.select('svg')
      var width = svg.attr('width')
      var height = svg.attr('height')
      this.g = svg.append('g')
        .attr('transform', 'translate(' + marge.top + ',' + marge.left + ')')

      // 准备数据
      this.nodes = [ // 节点集
        {
          name: '湖南邵阳'
        },
        {
          name: '山东莱州'
        },
        {
          name: '广东阳江'
        },
        {
          name: '山东枣庄'
        },
        {
          name: '泽'
        },
        {
          name: '恒'
        },
        {
          name: '鑫'
        },
        {
          name: '明山'
        },
        {
          name: '班长'
        }
      ]

      this.edges = [ // 边集
        {
          source: 0,
          target: 1,
          relation: '籍贯',
          value: 1.3
        },
        {
          source: 1,
          target: 2,
          relation: '舍友',
          value: 1
        },
        {
          source: 2,
          target: 3,
          relation: '舍友',
          value: 1
        },
        {
          source: 3,
          target: 4,
          relation: '舍友',
          value: 1
        },
        {
          source: 4,
          target: 5,
          relation: '籍贯',
          value: 2
        },
        {
          source: 5,
          target: 6,
          relation: '籍贯',
          value: 0.9
        },
        {
          source: 6,
          target: 7,
          relation: '籍贯',
          value: 1
        },
        {
          source: 7,
          target: 8,
          relation: '同学',
          value: 1.6
        }
      ]

      var node = this.g.selectAll('g.node').data(this.nodes, function(d, i) {
        return d.id || (d.id = ++i)
      })
      var nodeEnter = node
        .enter()
        .append('g')
        .attr('class', 'node')
        .attr('id', d => d.id)
        .attr('name', d => d.name)

      // https://github.com/zhangzn3/D3-Es6
      // https://blog.csdn.net/qq_39141486/article/details/103024733
      // 设置一个color的颜色比例尺，为了让不同的扇形呈现不同的颜色
      var colorScale = d3.scaleOrdinal()
        .domain(d3.range(this.nodes.length))
        .range(d3.schemeCategory10)
      this.forceSimulation = d3.forceSimulation(this.nodes)
        .force('link', d3.forceLink(this.edges))
        .force('charge', d3.forceManyBody())
        .force('center', d3.forceCenter())

      this.forceSimulation.nodes(this.nodes)
        .on('tick', this.ticked)
      // 生成边数据
      this.forceSimulation.force('link')
        .links(this.edges)
        .distance(function(d) { // 每一边的长度
          return d.value * 100
        })
      // 设置图形的中心位置
      this.forceSimulation.force('center')
        .x(width / 2)
        .y(height / 2)

      this.links = nodeEnter.append('line')
        .data(this.edges)
        .attr('stroke', function(d, i) {
          return colorScale(i)
        })
        .attr('class', 'line')
        .attr('stroke-width', 1)

      this.gs = nodeEnter.append('g')
        .data(this.nodes)
        .attr('transform', function(d, i) {
          var cirX = d.x
          var cirY = d.y
          return 'translate(' + cirX + ',' + cirY + ')'
        })
        .call(d3.drag()
          .on('start', this.started)
          .on('drag', this.dragged)
          .on('end', this.ended)
        )
      this.circleTexts = this.gs.append('text')
        .data(this.nodes)
        .attr('x', -10)
        .attr('y', -20)
        .attr('dy', 10)
        .text(function(d) {
          return d.name
        })
      this.circles = this.gs.append('circle')
        .data(this.nodes)
        .attr('r', 10)
        .attr('fill', function(d, i) {
          return colorScale(i)
        })
      this.forceSimulation.restart()
    },
    methods: {
      add() {
        this.nodes.push({
          name: '新增节点' + this.count
        })
        this.edges.push({
          source: this.count,
          target: this.count + 1,
          relation: '新增',
          value: 4
        })
        this.count++

        var colorScale = d3.scaleOrdinal()
          .domain(d3.range(this.nodes.length))
          .range(d3.schemeCategory10)
        var node = this.g.selectAll('g.node').data(this.nodes, function(d, i) {
          return d.id || (d.id = ++i)
        })
        this.forceSimulation.nodes(this.nodes).on('tick', this.ticked)
        this.forceSimulation.force('link', d3.forceLink(this.edges).distance(function(d) {
          return d.value * 100
        }))

        var nodeEnter = node
          .enter()
          .append('g')
          .attr('class', 'node')
          .attr('id', d => d.id)
          .attr('name', d => d.name)

        this.links = nodeEnter.append('line')
          .data(this.edges)
          .attr('stroke', function(d, i) {
            alert(JSON.stringify(d))
            return colorScale(i)
          })
          .attr('stroke-width', 1)
          .merge(this.links)

        this.gs = nodeEnter.append('g')
          .data(this.nodes)
          .attr('transform', function(d, i) {
            var cirX = d.x
            var cirY = d.y
            return 'translate(' + cirX + ',' + cirY + ')'
          })
          .call(d3.drag()
            .on('start', this.started)
            .on('drag', this.dragged)
            .on('end', this.ended)
          )
          .merge(this.gs)
        this.circleTexts = this.gs.append('text')
          .data(this.nodes)
          .attr('x', -10)
          .attr('y', -20)
          .attr('dy', 10)
          .text(function(d) {
            return d.name
          })
          .merge(this.circleTexts)

        this.circles = this.gs.append('circle')
          .data(this.nodes)
          .attr('r', 10)
          .attr('fill', function(d, i) {
            return colorScale(i)
          })
          .merge(this.circles)
        this.forceSimulation.alphaDecay(0.001)
        this.forceSimulation.restart()
      },
      ticked() {
        this.links
          .attr('x1', function(d) {
            return d.source.x
          })
          .attr('y1', function(d) {
            return d.source.y
          })
          .attr('x2', function(d) {
            return d.target.x
          })
          .attr('y2', function(d) {
            return d.target.y
          })
        this.gs.attr('transform', function(d) {
          return 'translate(' + d.x + ',' + d.y + ')'
        })
      },
      started(d) {
        if (!d3.event.active) {
          this.forceSimulation.alphaTarget(0.8).restart()
        }
        d.fx = d.x
        d.fy = d.y
      },
      dragged(d) {
        d.fx = d3.event.x
        d.fy = d3.event.y
      },
      ended(d) {
        if (!d3.event.active) {
          this.forceSimulation.alphaTarget(0)
        }
        d.fx = null
        d.fy = null
      }
    }
  }
</script>
