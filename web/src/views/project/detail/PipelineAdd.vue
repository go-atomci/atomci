<template>
    <div class="portlet-body">
      <div class="clearfix pv20">
        <div class="pipeBegin" :style="{'min-height': initHeight + 'px'}">
          <div class="pipeCircle"></div>
          <div class="pipeLine"></div>
          <div class="pipeCircle pipeCircleR"></div>
          <div class="cicleDashed" @click="$refs.createEnv.doCreate(-1)">+</div>
        </div>
        <template v-if="list">
          <template v-for="(item,pindex) in list">
            <template v-if="item">
              <div class="f-l" :style="{'min-height': initHeight + 'px'}">
                <div class="pipeContent">
                  <div class="pipeClose"><i class="el-icon-error" @click="delPipe(pindex)"></i></div>
                  <div class="pipeTitle">{{item.name}}</div>
                  <template v-if="item.steps">
                    <template v-for="(term,index) in item.steps">
                      <div class="clearfix">
                        <i class="el-close-fr el-icon-close" @click="delSteps(pindex, index)"></i>
                        <div class="pipeTask" @click="$refs.createNode.doCreate(true, pindex, index, term)">{{term.name}}</div>
                      </div>
                    </template>
                  </template>
                  <div class="pipeTask pipeAdd" @click="$refs.createNode.doCreate(false, pindex)">+ 添加节点</div>
                </div>
                <div class="pipeBegin">
                  <div class="pipeCircle"></div>
                  <div class="pipeLine"></div>
                  <div class="pipeCircle pipeCircleR"></div>
                  <div class="cicleDashed" @click="$refs.createEnv.doCreate(pindex)">+</div>
                </div>
              </div>
            </template>
          </template>
        </template>
        <div class="fb-layout">
          <el-button plain icon='el-icon-back' @click="goback">{{$t('bm.other.cancel')}}</el-button>
          <el-button type="primary" class="fb-ly-rbtn" icon="el-icon-edit" @click="doSubmit">{{$t('bm.add.saveStart')}}</el-button>
          <el-tooltip class="item" effect="dark" content="此次变更会应用于引用该流程的新建、回退的流水线！" placement="top-start">
            <i class="el-icon-info" style="font-size: 18px; margin-left:15px;"></i>
          </el-tooltip>
        </div>
      </div>
      <add-env ref="createEnv" v-on:listAdd='listAdd' :envList="envList"></add-env>
      <add-node ref="createNode" v-on:updateNode='updateNode' :taskList="taskList"></add-node>
    </div>
</template>

<style scoped>
  .pv20 {
    padding: 40px 20px;
  }
  .f-l {
    float: left;
    margin-bottom: 25px;
  }
  .pipeBegin {
    width: 80px;
    height: 80px;
    position: relative;
    float: left;
    cursor: pointer;
  }
  .cicleDashed {
    position: absolute;
    top: 43px;
    left: 28px;
    width:24px;
    height:24px;
    line-height: 24px;
    border-radius: 24px;
    background-color: #fff;
    text-align: center;
    font-size: 30px;
    border: 1px dashed #ccc;
    color: #555;
    display: none;
  }
  .pipeBegin:hover .cicleDashed {
    display: block;
  }
  .pipeCircle {
    position: absolute;
    width: 10px;
    height: 10px;
    border-radius: 10px;
    overflow: hidden;
    background-color: rgb(82, 155, 200);
    left: 15px;
    top: 50px;
  }
  .pipeCircleR {
    right: 15px;
    left: auto;
  }
  .pipeLine {
    width: 40px;
    margin: 0 auto;
    height: 1px;
    overflow: hidden;
    background-color: #ccc;
    margin-top: 55px;
  }
  .pipeContent {
    position: relative;
    float: left;
    width: 210px;
    min-height: 100px;
    box-shadow: 1px 1px 3px #ccc;
    border: 1px solid #fafafa;
    border-radius: 6px;
    padding-top: 35px;
  }
  .pipeTitle {
    position: absolute;
    left: 15%;
    top: -20px;
    z-index: 3;
    width: 70%;
    height: 40px;
    line-height: 40px;
    overflow: hidden;
    color: #fff;
    background: #2a2a2b;
    border-radius: 25px;
    text-align: center;
  }
  .pipeTask {
    width: 82%;
    height: 36px;
    line-height: 36px;
    margin-left: 5%;
    background-color: rgb(250, 251, 253);
    border-radius: 4px;
    margin-bottom: 6px;
    padding: 0 6px;
    overflow: hidden;
  }
  .el-close-fr {
    float: right;
    width: 9%;
    margin-top: 12px;
    margin-right: 4px;
    font-size: 16px;
    cursor: pointer;
  }
  .pipeTask:hover {
    cursor: pointer;
  }
  .pipeAdd {
    width:90%;
    background: transparent;
    border: 1px dashed #ccc;
    cursor: pointer;
  }
  .pipeClose {
    position: absolute;
    right: -10px;
    top: -5px;
  }
  .pipeClose>i {
    font-size: 24px;
    color:#ddd;
  }
  .pipeClose>i:hover {
    color: #555;
    cursor: pointer;
  }
</style>
<script>
import { mapGetters } from 'vuex';
  import { Message, MessageBox } from 'element-ui';
  import backend from '@/api/backend';
  import addEnv from '../components/PipeAddEnv';
  import addNode from '../components/PipeAddNode';

export default {
  data() {
    return {
      col: {},
      list: [],
      listData: {},
      initHeight: 100,
      envList:[],
      taskList:[]
    };
  },
  components: {
    addEnv,
    addNode,
  },
  computed: {
    ...mapGetters({
      projectID: 'projectID',
    }),
  },
  mounted() {},
  created() {
    backend.getProjectEnvsAll(this.projectID, (data) => {
      this.envList = data
    }),
    backend.getStepAll((data) => {
      if(data){
        this.taskList = data;
      }
    });
    this.getList();
  },
  // watch: {
  //   '$route': 'getList',  // 监听router值改变时，改变导航菜单激活项
  // },
  methods: {
    getList() {
      if(this.$route.params.pipeId) {
        backend.getProjectPipeDetail(this.projectID, this.$route.params.pipeId,(data) => {
          this.list = data.config;
          this.listData = data;
          this.initNodeHeight();
        });
      }
    },
    initNodeHeight() {
      let array = [];
      let max = 0;
      if ( this.list === undefined ){
        this.list =[] 
      } 
      this.list.map((i) => {
        if(i.steps && i.steps.length) array.push(i.steps.length);
      });
      if(array) max = Math.max.apply(Math,array);
      this.initHeight = (max + 1)*42 + 60;
    },
    listAdd(index, obj) {
      let checks = false;
      this.list.map((i) => {
        if(i.stage_id == obj.stage_id){
          checks = obj.name;
          return;
        }
      });
      if(checks) {
        Message.error(`存在重复的环境-【${checks}】`);
      } else {
        this.list.splice(index, 0, obj);
        Message.success('环境添加成功');
        this.$refs.createEnv.dialogFormVisible = false;
      }
    },
    delPipe(index) {
      MessageBox.confirm('是否要删除当前环境', '提示', { type: 'warning' })
        .then(() => {
          this.list.splice(index, 1);
          Message.success('成功删除环境');
          this.initNodeHeight();
        })
        .catch(() => {});
    },
    updateNode(flag, pIndex, obj, index) {
      if(flag) {
        this.list[pIndex].steps.splice(index, 1, obj);
        //Message.success('更新节点成功');
      } else {
        this.list[pIndex].steps.push(obj);
        //Message.success('添加节点成功');
        this.initNodeHeight();
      }
      this.$refs.createNode.dialogFormVisible = false;
    },
    delSteps(pIndex, index) {
      MessageBox.confirm('是否要删除当前节点', '提示', { type: 'warning' })
        .then(() => {
          this.list[pIndex].steps.splice(index, 1);
          //Message.success('成功删除节点');
          this.initNodeHeight();
        })
        .catch(() => {});
    },
    goback() {
      this.$router.push({
        name: 'projectPipeline',
        params: {'projectID': this.projectID}
      });
    },
    doSubmit() {
      this.list.map((item, index) => {
        item.index = index + 1;
        const steps = item.steps;
        if(steps && steps.length > 0) {
          steps.map((i, sindex) => {
            i.index = sindex + 1;
          });
        }
      });
      this.listData.config = this.list;
      MessageBox.confirm('是否确定更新当前配置', '提示', {type: 'warning'})
        .then(() => {
          backend.editProjectPipe(this.projectID, this.$route.params.pipeId, this.listData, (data) => {
            Message.success('配置更新成功');
            this.$router.push({
              name: 'projectPipeline',
              params: {'projectID': this.projectID}
            });
          });
        })
        .catch(() => {});
    },
  },
};
</script>