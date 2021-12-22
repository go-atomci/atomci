<template>
  <div class="page-content">
    <div class="portlet-body min-height150">
      <template>
        <div v-show="proShow">
          <el-row>
            <el-col :span="24">
              <div class="setTitle">
                <div class="f-r">
                  <div class="summaryUser">
                    {{projectInfo.members}}<span class="fontSmall"> 人</span>
                    <p class="fontSmall">项目人员</p>
                  </div>
                  <div class="summaryUser">
                    {{projectInfo.code_repos}}<span class="fontSmall"> 个</span>
                    <p class="fontSmall">关联代码仓库</p>
                  </div>
                </div>
                {{projectInfo.name}}<span class="setDescription">{{projectInfo.description}}</span>
              </div>
            </el-col>
          </el-row>
          <el-row>
            <el-col class="setText">
              管理员：{{projectInfo.owner}}
            </el-col>
            <el-col class="setText">
              项目状态：{{projectInfo.status === 1 ? '运行中' : '已结束'}}
            </el-col>
          </el-row>
          <el-row>
            <el-col class="setText">
              开始时间：{{projectInfo.start_at}}
            </el-col>
            <el-col class="setText">结束时间：{{projectInfo.end_at}}</el-col>
          </el-row>
          <el-row>
            <el-col class="setText">项目成员：{{projectInfo.membersName}}</el-col>
          </el-row>
        </div>
      </template>
    </div>
    <div class="portlet-body min-height150" style="float: left">
      <div class="setTitle" style="float: left">流水线统计</div>
      <ul class="version-list mt10 clearfix">
        <li v-for="(item,index) in projectInfo.releases">
          <span class="span-count">{{item.count}}</span> 个
          <p>{{item.env}}</p>
        </li>
      </ul>
    </div>
    <!-- <div class="portlet-body p-chart min-height150">
      <div class="setTitle pl16">已部署应用</div>
      <ul class="chart-list mt10 clearfix">
        <li v-for="(item,index) in optionData">
          <div>
            <div class="chart-title">{{item.cluster}}</div>
            <div class="chart-style" :id="'chart'+ index"></div>
          </div>
        </li>
      </ul>
    </div> -->
    <div class="portlet-body" v-show="false">
      <div class="setTitle">操作日志</div>
      <div class="mt10">
        <el-table border :data="logList">
          <el-table-column prop="name" label="人员" min-width="8%" />
          <el-table-column prop="content" label="内容" min-width="20%" />
          <el-table-column prop="create_at" label="时间" min-width="8%" />
        </el-table>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .min-height150 {
    min-height: 150px;
  }

  .pl16 {
    padding-left: 16px;
  }

  .p-chart {
    padding: 10px 16px 10px 0;
  }

  .member-btn {
    width: 550px;
    text-align: right;
  }

  .f-r {
    float: right;
  }

  .pv20 {
    padding: 20px;
  }

  .memberRow .el-select {
    width: 100%;
  }

  .containerMember {
    width: 550px;
    border: 1px solid #ccc;
    padding: 10px;
  }

  .el-tag {
    margin-right: 5px;
    margin-bottom: 3px;
  }

  .mb15 {
    margin-bottom: 15px;
  }

  .setTitle {
    font-size: 18px;
    color: #333;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: bold;
    line-height: 40px;
  }

  .setDescription {
    color: #606266;
    font-size: 14px;
    margin-left: 10px;
    font-weight: 400;
  }

  .summaryUser {
    line-height: 25px;
    display: inline-block;
    padding: 0 20px;
    border-right: 1px solid #DCDFE6;
  }

  .setTitle .summaryUser:last-child {
    border-right: 0;
  }

  .fontSmall {
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    line-height: 20px;
    color: #909399;
  }

  .setText {
    width: 400px;
    line-height: 20px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-right: 10px;
    margin-bottom: 15px;
    color: #333;
    font-weight: 400;
  }

  .form-number {
    width: 320px;
    float: right;
    margin-right: 5px;
  }

  .form-col {
    width: 150px;
    margin-right: 10px;
  }

  .font-title {
    font-size: 16px;
    color: #409EFF;
    display: inline-block;
    line-height: 50px;
  }

  .el-form-item {
    margin-bottom: 5px;
  }

  .proj-edit {
    font-size: 20px;
    color: #409EFF;
    vertical-align: middle;
    cursor: pointer;
  }

  .version-list li {
    float: left;
    width: 250px;
    border-radius: 6px;
    background-color: #F5F7FA;
    text-align: center;
    padding: 20px 0;
    margin-right: 20px;
    margin-bottom: 20px;
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: rgba(144, 147, 153, 1);
  }

  .span-count {
    font-size: 32px;
    font-family: PingFangSC-Medium, PingFang SC;
    font-weight: 500;
    color: rgba(51, 51, 51, 1);
    line-height: 45px;
  }

  .chart-list li {
    float: left;
    width: 33.3%;
    padding-left: 16px;
    margin-bottom: 16px;
  }

  .chart-list li>div {
    height: 274px;
    padding: 10px;
    background: rgba(255, 255, 255, 1);
    box-shadow: 0px 2px 4px 0px rgba(64, 158, 255, 0.12);
    border: 1px solid rgba(220, 223, 230, 1);
  }

  .chart-title {
    font-size: 16px;
    font-family: PingFangSC-Medium, PingFang SC;
    font-weight: 500;
    color: rgba(51, 51, 51, 1);
    line-height: 22px;
  }

  .chart-style {
    width: 100%;
    height: 100%;
  }
</style>
<style>
  .pv30 .el-form-item__error {
    left: 3px !important;
  }
</style>
<script>
  import { mapGetters } from 'vuex';
  import backend from '@/api/backend';
  import UtilsFn from '@/common/utils';
  const echarts = require('echarts');
  export default {
    data() {
      return {
        proShow: false,
        projectInfo: {},
        versionList: [],
        optionData: [],
        logList: [],
      };
    },
    computed: {
      ...mapGetters({
        projectID: 'projectID',
      }),
    },
    components: {},
    activated() {
      this.getProjectInfo();
      // this.getChartData();
    },
    created() {},
    mounted() {
      window.addEventListener('resize', () => {
        this.$nextTick(() => {
          this.getData();
        });
      });
    },
    methods: {
      getChartData() {
        backend.getSummaryChart(this.projectID, (data) => {
          this.optionData = data;
          this.$nextTick(() => {
            this.getData();
          });
        });
      },
      getProjectInfo() {
        // TODO: project id refacotr to router params
        if (this.$route.query.id !== undefined) {
          this.projectID = this.$route.query.id
          this.$store.dispatch('project/setProjectID', this.$route.query.id);
        }
        backend.getProjectDetail(this.projectID, (data) => {
          if (data) {
            data.start_at = data.start_at ? UtilsFn.format(new Date(data.start_at), 'yyyy-MM-dd hh:mm:ss') : '';
            data.end_at = data.end_at ? UtilsFn.format(new Date(data.end_at), 'yyyy-MM-dd hh:mm:ss') : '';
            data.membersName = data.membersName && data.membersName.length > 0 ? data.membersName.join('; ') : '';
            this.projectInfo = Object.assign({}, data);
            this.proShow = true;
          }
        });
      },
      getData() {
        this.optionData.map((i, index) => {
          this.getChart('chart' + index, i);
        });
      },
      getChart(id, chartData) {
        const dom = document.getElementById(id);
        if(!dom) return;
        const parentNodeStyle = window.getComputedStyle(dom.parentNode);
        // 用于使chart自适应高度和宽度,通过窗体高宽计算容器高宽
        dom.style.width = `${parentNodeStyle.width}px`;
        dom.style.height = `${parentNodeStyle.height}px`;
        const deployChart = echarts.init(dom);
        deployChart.resize();
        const option = {
          tooltip: {
            trigger: 'item',
            formatter: "{a} <br/>{b}: {c} ({d}%)"
          },
          legend: {
            orient: 'vertical',
            data: [`正常：${chartData.runningNum}`, `警告：${chartData.pendingNum}`, `严重：${chartData.stoppedNum}`],
            bottom: 50,
            right: 0
          },
          series: [{
            name: '部署',
            type: 'pie',
            radius: ['50%', '70%'],
            avoidLabelOverlap: false,
            label: {
              normal: {
                show: false,
                position: 'center'
              },
              emphasis: {
                show: true,
                textStyle: {
                  fontSize: '30',
                  fontWeight: 'bold'
                }
              }
            },
            labelLine: {
              normal: {
                show: false
              }
            },
            data: [{
                value: chartData.runningNum,
                name: `正常：${chartData.runningNum}`,
                itemStyle: {
                  normal: {
                    'color': '#34BD85'
                  }
                }
              },
              {
                value: chartData.pendingNum,
                name: `警告：${chartData.pendingNum}`,
                itemStyle: {
                  normal: {
                    'color': '#FDC942'
                  }
                }
              },
              {
                value: chartData.stoppedNum,
                name: `严重：${chartData.stoppedNum}`,
                itemStyle: {
                  normal: {
                    'color': '#F91628'
                  }
                }
              }
            ]
          }]
        };
        if (option && typeof option === "object") {
          deployChart.setOption(option, true);
        }
      },

    }
  }
</script>
