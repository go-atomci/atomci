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
  <el-dialog top='15vh' v-if="dialogFormVisible" :close-on-click-modal="true" :show-close="false" width='46%' title="编辑"
    :visible.sync="dialogFormVisible" class="commonDialog projectPubModuleList">
    <div>
      <i id="color" class="el-icon-close" @click="handleClose" style="cursor:pointer;position:absolute;right:15px;top:15px;"></i>
    </div>
    <el-form ref="ruleForm" :model="form" :rules="rules">
      <div class="deploy-mirror-wrap">
        <el-row>
          <el-col :span="18">
            <el-form-item :label="$t('bm.deployCenter.pipelineName')" prop="version_no">
              <el-input v-model="form.version_no" :placeholder="$t('bm.add.verNameNo16Node')"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="18">
            <el-form-item :label="$t('bm.deployCenter.pipelineDesc')" prop="name">
              <el-input v-model="form.name"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </div>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose" style="margin-top:20px">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit" >{{$t('bm.other.confirm')}}</el-button>
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
  import backend from '../../../api/backend';


  export default {
    data() {
      return {
        // 是否启用切换分支
        form: {
          name: '',
          id: undefined,
          version_no: ''
        },
        dialogFormVisible: false,
        rules: {
          name: [{
            required: true, message: '请输入版本描述', trigger: 'blur',
          }, ],
          version_no: [{
            required: true, message: '请输入版本号', trigger: 'blur'
          }, ]
        },
      };
    },
    components: {
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
        const params = {
          name: this.form.name,
          version_no: this.form.version_no,
          id: this.form.id
        };
        const that = this;
        backend.updateProjectCI(this.projectID, this.form.id, params, (data) => {
          Message.success(this.$t('bm.add.optionSuc'));
          this.dialogFormVisible = false;
          that.$emit('getPublishBaseInfo');
        }, () => {
          that.$emit('getPublishBaseInfo');
        });
      },
      doShows(publishInfo) {
        this.form.name = publishInfo.name
        this.form.id = publishInfo.id
        this.form.version_no = publishInfo.version_no
        this.dialogFormVisible = true
      },
      handleClose() {
        this.dialogFormVisible = false;
      },
    },
  };
</script>
