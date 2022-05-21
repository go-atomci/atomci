
<template>
  <el-dialog z-index="1100" top='15vh' v-if="dialogFormVisible" :close-on-click-modal="false" width='65%' :title="title"
    :visible.sync="dialogFormVisible" class="createDialog">
    <el-table border :data="dataList" @select-all="handleSelectAll" @select='handleSelect' ref="appDeploy">
      <el-table-column type="selection" width="50" disabled='true' :show-overflow-tooltip=true />
      <el-table-column prop="name" :label="$t('bm.deployCenter.name')" sortable min-width="12%" :show-overflow-tooltip=true />
      <el-table-column prop="type" :label="$t('bm.deployCenter.type')" sortable min-width="12%"
        :show-overflow-tooltip=true />
    </el-table>
    <div slot="footer" class="dialog-footer">
      <el-button @click="doClose" style="margin-top:20px">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>

<style scoped>
  .pv10 {
    padding-top: 10px;
    padding-bottom: 10px;
    height: 40px;
    position: relative;
  }
  .radioStyle {
    position: absolute;
    left: 100px;
    top: 12px;
  }
  .radioStyles {
    position: absolute;
    left: 0px;
    top: 12px;
  }
</style>
<script>
  import { mapGetters } from 'vuex';
  import { Message } from 'element-ui';
  import backend from '../../../api/backend';
  import Refresh from '../../../components/utils/Refresh';

  export default {
    data() {
      return {
        dialogFormVisible: false,
        title: '',
        dataList: [],
        publishId: '',
        stageId: '',
        selectList: [],
      };
    },
    components: {
      Refresh,
    },
    created() {

    },
    computed: {
      ...mapGetters({
        projectIDgetter: 'projectID',
    }),
    projectID() {
        if (this.projectIDgetter === 0 || this.projectIDgetter === undefined) {
          this.$store.dispatch('project/setProjectID', this.$route.params.projectID);
          return this.$route.params.projectID
        } else {
          return this.projectIDgetter
        }
    },
    },
    methods: {
      handleSelectAll(val) {
        this.selectList = val;
      },
      handleSelect(val) {
        this.selectList = val;
      },
      renderHeader(h, {column}) {
        return h(
          'div', [
            h('span', column.label),
            h('el-tooltip', {
              props: {
                effect: 'dark',
                content: '镜像版本由当前环境关联集群的镜像仓库根据应用及发布分支得到',
                placement: 'top',
              },
            },
            [h('i', {
              class: 'el-icon-question',
              style: 'font-size:16px;margin-left:8px;cursor:pointer;'
            })],
            )
          ]
        );
      },
      doShows(item) {
        this.title = item.step;
        this.publishId = item.id;
        this.stageId = item.stage_id;
        this.dataList = [];
        this.selectList = [];
        backend.getDeploy(this.projectID, item.id, item.stage_id, (data) => {
          this.dataList = data
          this.dialogFormVisible = true;
          this.$nextTick(() => {
            this.selectList = this.dataList;
            this.toggleSelection();
            this.selectList.forEach((item, index) => {
              this.toggleSelection([item]);
            });
          });

        });
      },
      toggleSelection(rows) {
        if (rows) {
          rows.forEach((row) => {
            this.$refs.appDeploy.toggleRowSelection(row, true);
          });
        } else {
          this.$refs.appDeploy.clearSelection();
        }
      },
      doClose() {
        this.toggleSelection();
        this.dialogFormVisible = false;
        this.selectList = [];
      },
      doSubmit() {
        if(this.selectList.length > 0) {
          let apps = [];
          this.selectList.forEach((item) => {
            const cl = {
              project_app_id: item.project_app_id,
            };
            apps.push(cl);
          });
        
            const params = {
              action_name: 'trigger',
              apps: apps
            };
            backend.setDeploy(this.projectID, this.publishId, this.stageId, params, (data)=> {
              Message.success('部署成功');
              this.$emit('getprojectReleaseList', true);
              this.dialogFormVisible = false;
            });
        } else {
          Message.error('请至少选择一条数据！');
        }
      }
    },
  };
</script>
