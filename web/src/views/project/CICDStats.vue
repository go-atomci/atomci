<template>
  <div class="page-content">
    <div class="portlet-body p20">
      <el-form :model="form" ref="ruleForm">
        <el-row>
          <el-col class="census">
            <el-form-item prop="census" label="统计应用">
              <el-select v-model="form.cencus" filterable multiple v-on:change="cencusChange" placeholder="全部">
                <el-option v-for="(item, index) in censusList" :key="index" :label="item.name" :value="item.id">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col class="census">
            <el-form-item prop="time" label="发布时间">
              <el-date-picker style="width:288px" v-model="timeList" :clearable="false" type="daterange"
                v-on:change="dateChange" range-separator="~" format="yyyy-MM-dd" :picker-options="pickerOptions"
                start-placeholder="开始日期" end-placeholder="结束日期" />
            </el-form-item>
          </el-col>
          <el-col class="census">
            <el-form-item prop="env" label="发布环境">
              <el-select v-model="form.env" filterable multiple v-on:change="cencusChange" placeholder="全部">
                <el-option v-for="(item, index) in envList" :key="index" :label="item.name" :value="item.id">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>
    <div class="chartPanel">
      <el-row>
        <el-col :span="12" class="chartSpan">
          <div>
            <div class="chartTitle">持续构建统计</div>
            <div id="createChart" class="chartStyle"></div>
          </div>
        </el-col>
        <el-col :span="12" class="chartSpan">
          <div>
            <div class="chartTitle">持续部署统计</div>
            <div id="deployChart" class="chartStyle"></div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<style scoped>
  .p20 {
    padding: 20px;
    padding-bottom: 0;
  }

  .census {
    width: 400px;
  }

  .census .el-form-item__label {
    width: 80px;
    text-align: left;
  }

  .el-date-editor .el-range-input {
    width: 46% !important;
  }

  .censusChart {
    width: 900px;
    padding: 10px;
    border: 1px solid #ccc;
    overflow: hidden;
  }


  .el-select {
    width: 280px;
  }

  .chartStyle {
    width: 100%;
    height: 500px;
    overflow: hidden;
  }

  .chartPanel .chartSpan:first-child div {
    margin-right: 16px;
  }

  .chartSpan>div {
    background-color: #fff;
    padding-top: 15px;
  }

  .chartTitle {
    line-height: 45px;
    font-size: 20px;
    font-family: PingFangSC-Medium, PingFang SC;
    color: #333333;
    padding-left: 30px;
    margin-bottom: 15px;
  }
</style>
<script>
import { mapGetters } from 'vuex';
import Utils from '@/common/utils';
import backend from '@/api/backend';
const echarts = require('echarts');

  export default {
    data() {
      return {
        form: {
          cencus: [],
          env: []
        },
        pickerOptions: {
          disabledDate(time) {
            return time.getTime() > Date.now() - 8.64e6;
          },
        },
        censusList: [],
        timeList: '',
        envList: [],
        option: null,
        optionData: {},
        chartDataServer: {},
        chartDataCreate: {},
      };
    },
    computed: {
      ...mapGetters({
        loading: 'getLoading',
        projectID: 'projectID',
      }),
    },
    components: {},
    created() {
      this.getApps();
      this.getEnv();
      let end = (new Date()).getTime();
      // end = end - (1 * 24 * 3600 * 1000);
      let start = end - (7 * 24 * 3600 * 1000);
      end = Utils.format(new Date(end), 'yyyy-MM-dd');
      start = Utils.format(new Date(start), 'yyyy-MM-dd');
      this.timeList = [start, end];
    },
    mounted() {
      this.getData();
      window.addEventListener('resize', () => {
        this.$nextTick(() => {
          this.getDeploy('deployChart', this.chartDataServer);
          this.getDeploy('createChart', this.chartDataCreate);
        });
      });
    },
    methods: {
      getEnv() {
        backend.getProjectEnvsAll(this.projectID, (data) => {
          this.envList = data;
        });
      },
      getApps() {
        backend.getAppAll(this.projectID, (data) => {
          this.censusList = data;
        });
      },
      cencusChange() {
        this.getData();
      },
      dateChange(val) {
        if (val) {
          val.map((i, index) => {
            this.timeList[index] = Utils.format(new Date(i), 'yyyy-MM-dd');
          });
          this.getData();
        }
      },
      getData() {
        this.optionData = {};
        const params = {
          "start_time": this.timeList[0],
          "end_time": this.timeList[1],
          "app_ids": this.form.cencus,
          "env_ids": this.form.env
        };
        backend.getStats(this.projectID, params, (data) => {
          if (data) {
            this.optionData = {
              xData: [],
              createFailed: [],
              createSuccess: [],
              serverFailed: [],
              serverSuccess: [],
              totalFailed: [],
              totalSuccess: [],
              serverPercent: [],
              createPercent: []
            };
            data.map((i) => {
              this.optionData.xData.push(i.time);
              this.optionData.createFailed.push(i['build_failed']);
              this.optionData.createSuccess.push(i['build_success']);
              this.optionData.serverFailed.push(i['deploy_failed']);
              this.optionData.serverSuccess.push(i['deploy_success']);
              this.optionData.totalFailed.push(i['total_failed']);
              this.optionData.totalSuccess.push(i['total_success']);
              const pert = (i['deploy_success'] || i['deploy_failed']) ? Math.round((i['deploy_success']) * 100 / (i['deploy_failed'] + i['deploy_success'])) : 0;
              const perts = (i['build_success'] || i['build_failed']) ? Math.round((i['build_success']) * 100 / (i['build_failed'] + i['build_success'])) : 0;
              this.optionData.serverPercent.push(pert);
              this.optionData.createPercent.push(perts);
            });
            this.chartDataServer = {
              'data': ['部署成功', '部署失败', '成功率'],
              'success': {
                'name': '部署成功',
                'color': '#2DCA93',
                'count': this.optionData.serverSuccess
              },
              'failed': {
                'name': '部署失败',
                'color': '#FACF2A',
                'count': this.optionData.serverFailed
              },
              'percent': {
                'color': '#67C23A',
                'count': this.optionData.serverPercent
              }
            };
            this.chartDataCreate = {
              'data': ['构建成功', '构建失败', '成功率'],
              'success': {
                'name': '构建成功',
                'color': '#409EFF',
                'count': this.optionData.createSuccess
              },
              'failed': {
                'name': '构建失败',
                'color': '#FFBD74',
                'count': this.optionData.createFailed
              },
              'percent': {
                'color': '#67C23A',
                'count': this.optionData.createPercent
              }
            };
            this.getDeploy('deployChart', this.chartDataServer);
            this.getDeploy('createChart', this.chartDataCreate);
          }
        });
      },
      getDeploy(id, chartData) {
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
            trigger: 'axis',
            axisPointer: {
              type: 'cross',
              crossStyle: {
                color: 'red'
              }
            },
            formatter: function (params) {
              let text = params[0].name;
              params.map((i) => {
                i.value = i.seriesName === '成功率' ? i.value + '%' : i.value;
                text += '<br/>' + i.marker + i.seriesName + ' : ' + i.value;
              });
              return text;
            }
          },
          legend: {
            data: chartData.data
          },
          xAxis: [
            {
              type: 'category',
              data: this.optionData.xData,
              axisPointer: {
                type: 'shadow'
              }
            }
          ],
          yAxis: [
            {
              type: 'value',
              name: '数量',
              min: 0,
              axisLabel: {
                formatter: '{value}'
              }
            },
            {
              type: 'value',
              name: '成功率',
              min: 0,
              axisLabel: {
                formatter: '{value} %'
              }
            }
          ],
          series: [
            {
              name: chartData.success.name,
              type: 'bar',
              stack: '222',
              itemStyle: {
                normal: { color: chartData.success.color },
                width: 20
              },
              data: chartData.success.count
            },
            {
              name: chartData.failed.name,
              type: 'bar',
              stack: '222',
              itemStyle: {
                normal: { color: chartData.failed.color }
              },
              data: chartData.failed.count
            },
            {
              name: '成功率',
              type: 'line',
              itemStyle: {
                normal: { color: chartData.percent.color }
              },
              yAxisIndex: 1,
              data: chartData.percent.count
            }
          ]
        };
        if (option && typeof option === "object") {
          deployChart.setOption(option, true);
        }
      },
    }
  }
</script>
