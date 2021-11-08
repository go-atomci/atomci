<style scoped>
  .projectPubModuleList .el-table {
    min-height: 32vh;
  }

  .projectPubModuleList .dialog-footer {
    margin-top: 0px;
  }

  .projectPubModuleList .dialog-footer .el-checkbox__label {
    color: #f56c6c;
  }

  .projectPubModuleList .el-dialog__footer {
    overflow: hidden;
    width: 100%;
    margin-top: 0px;
  }

  .projectPubModuleList .el-checkbox {
    margin-right: 5px;
  }

  #color:hover {
    color: #409EFF;
  }
  .reset-el-select {
    width: 150px;
  }

</style>
<style>
  .dialogTable th .el-checkbox__input {
    vertical-align: -2px;
    margin-right: 5px;
  }
</style>
<template>
  <el-dialog top='15vh' v-if="dialogFormVisible" :close-on-click-modal="false" :show-close="false" width='70%' :title="username"
    :visible.sync="dialogFormVisible" class="commonDialog projectPubModuleList">
    <div>
      <i id="color" class="el-icon-close" @click="handleClose" style="cursor:pointer;position:absolute;right:15px;top:15px;"></i>
    </div>
    <el-form ref="ruleForm" :model="form" :rules="rules">
      <div class="deploy-mirror-wrap">
        <el-row>
          <el-col :span="18">
            <el-form-item :label="$t('bm.deployCenter.pipelineName')" prop="version_no">
              <el-input v-model="form.version_no" :placeholder="$t('bm.add.verNameNo16Node')" disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="18">
            <el-form-item :label="$t('bm.deployCenter.pipelineDesc')" prop="name">
              <el-input v-model="form.name" disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </div>
    </el-form>

    <el-table  border :data="tableList" @select-all="handleSelectAll" @select='handleSelect' ref="appmodule" >
      <el-table-column type="selection" min-width="7%" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="name" :label="$t('bm.deployCenter.name')" sortable min-width="12%" :show-overflow-tooltip="true" />
      <el-table-column prop="language" :label="$t('bm.deployCenter.language')" sortable min-width="8%" :show-overflow-tooltip="true" />
      <el-table-column prop="build_path" label="构建目录" sortable min-width="10%" :show-overflow-tooltip="true" />
      <el-table-column :label="$t('bm.deployCenter.releaseBran')" min-width="15%">
        <template slot-scope="scope">
          <el-select v-model.trim="scope.row.branch_name" filterable :placeholder="$t('bm.add.selectSubmitBra')">
            <el-option v-for="(item, index) in scope.row.branch_history_list" :key="index" :label="item" :value="item">
            </el-option>
          </el-select>
        </template>
      </el-table-column>
      <el-table-column :label="$t('bm.add.customBuild')" sortable min-width="40%">
        <template slot-scope="scope">
          <el-input v-model="scope.row.compile_command" :placeholder="$t('bm.add.customBuildCom')"> </el-input>
        </template>
      </el-table-column>
    </el-table>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose" style="margin-top:20px">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit" :loading="subloading">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
  import {
    Message
  } from 'element-ui';
  import {
    mapGetters
  } from 'vuex';
  import backend from '@/api/backend';


  export default {
    props: ['listData', 'pubType', 'cpData'],
    data() {
      return {
        // 是否启用切换分支
        enable_switch_branch: true,
        username: '',
        envStageID: '',
        publishid: '',
        unit_test: false,
        form: {
          name: '',
          image_version: '',
          version_no: ''
        },
        premerge: true,
        show: '1',
        radio: '1',
        dialogFormVisible: false,
        searchList: [{
            key: 'app_name',
            txt: this.$t('bm.deployCenter.name')
          },
          {
            key: 'app_type',
            txt: this.$t('bm.deployCenter.type')
          },
          {
            key: 'app_language',
            txt: this.$t('bm.deployCenter.language')
          },
        ],
        branchList: [],
        preBranchList: [],
        good: true,
        subloading: true,
        // 直接选中项
        selectList: [],
        rules: {
          name: [{
            required: true
          }, ],
          version_no: [{
            required: true
          }, ]
        },
        tableList: [],
        cpList: [],
      };
    },
    computed: {
      ...mapGetters({
        projectID: 'projectID',
      })
    },
    methods: {
      golist(rows) {
        if (rows) {
          rows.forEach((row) => {
            // this.goselect(row, row.branch_name)
          });
        }
      },
      handleSelectAll(val) {
        this.selectList = val;
      },
      handleSelect(val) {
        this.selectList = val;
      },
      doSubmit() {
        const app = [];
        for (const a of this.selectList) {
          const ats = {
            branch_name: a.branch_name,
            project_app_id: a.project_app_id,
            compile_command: a.compile_command,
          };
          app.push(ats);
        }
        if (app.length === 0) {
          Message.error('请至少选择一条数据！');
          return;
        }
        // TODO: add env vars
        const params = {
          apps: app,
          action_name: 'trigger'
        };
        const that = this;
        backend.setBuildMerge(this.projectID, this.publishid, this.envStageID, params, (data) => {
          Message.success(this.$t('bm.add.optionSuc'));
          this.dialogFormVisible = false;
          that.$emit('getprojectReleaseList');
        }, () => {
          that.$emit('getprojectReleaseList');
        });
      },
      toggleSelection(rows) {
        if (rows) {
          rows.forEach((row) => {
            this.$refs.appmodule.toggleRowSelection(row, true);
          });
        } else {
          this.$refs.appmodule.clearSelection();
        }
      },
      doShows(publishid, envStageID, step) {
        this.selectList = [];
        this.enable_switch_branch = true;
        this.tableList = [];
        this.cpList = [];
        this.username = step;
        backend.getBuildMerge(this.projectID, publishid, envStageID, (data) => {
          this.subloading = false;
          this.tableList = data.apps.map((i, index) => {
            // 
            i.branch_history_list = i.branch_history_list
            i.name = i.app_name;
            return i;
          });
          this.form.version_no = data.versionNo || '';
          this.form.name = data.versionName || '';
          if (data.apps[0].image_version === '') {
            this.radio = '1';
            this.show = '1';
          } else {
            this.radio = '2';
            this.show = '2';
          }
          this.selectList = data.apps;
          this.toggleSelection();
          this.$nextTick(() => {
            this.selectList.forEach((item, index) => {
              this.toggleSelection([item]);
            });
          });
        });
        this.publishid = publishid;
        this.envStageID = envStageID;
        this.form.name = '';
        this.premerge = true;
        this.form.image_version = '';
        this.show = '1';
        this.radio = '1';
        this.dialogFormVisible = true;
      },
      handleClose() {
        this.toggleSelection();
        this.dialogFormVisible = false;
        this.selectList = [];
      },
    },
  };
</script>
