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
<template>
  <el-dialog top='15vh' v-if="dialogFormVisible" :close-on-click-modal="false" width='65%' :title="setname" :visible.sync="dialogFormVisible" class="createDialog">
    <el-form ref="ruleForm" :model="form" :rules="rules" label-width="52px">
      <div>
        <el-row>
          <el-col :span="18">
            <el-form-item :label="$t('bm.deployCenter.pipelineName')" prop="image_version">
              <el-input :disabled="details" v-model="form.image_version" :placeholder="$t('bm.add.verNameNo16Node')"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="18">
            <el-form-item :label="$t('bm.deployCenter.pipelineDesc')" prop="name">
              <el-input v-model="form.name" :disabled="details"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="18">
            <el-form-item label="绑定流程" prop="pipe">
              <el-select v-model="form.pipe" filterable :disabled="details">
                <el-option v-for="(item, index) in pipeArray" :key="index" :label="item.name" :value="item.id"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </div>
    </el-form>
    <el-table v-show="!details" style="margin-top:2%" border :data="tableList" @select-all="handleSelectAll" @select='handleSelect' ref="multipleTable">
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
    <el-table v-show="details" style="margin-top:2%" border :data="tableListdetails" ref="appmodule">
      <el-table-column prop="name" :label="$t('bm.deployCenter.name')" sortable min-width="12%" :show-overflow-tooltip=true />
      <el-table-column prop="type" :label="$t('bm.deployCenter.type')" sortable min-width="12%" :show-overflow-tooltip=true />
      <el-table-column prop="language" :label="$t('bm.deployCenter.language')" sortable min-width="10%" :show-overflow-tooltip=true />
      <el-table-column prop="branch_name" :label="$t('bm.deployCenter.releaseBran')" min-width="11%" :show-overflow-tooltip=true />
      <el-table-column prop="image_version" :label="$t('bm.add.version')"  sortable min-width="10%" :show-overflow-tooltip=true />
    </el-table>
    <div slot="footer" class="dialog-footer">
      <el-button v-show="!details" @click="handleClose" style="margin-top:20px">{{$t('bm.other.cancel')}}</el-button>
      <el-button v-show="!details" type="primary" @click="doSubmit('ruleForm')">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
  import { mapGetters } from 'vuex';
  import { Message } from 'element-ui';
  import backend from '@/api/backend';
  import Refresh from '@/components/utils/Refresh';

  const formData = {
    name: '',
    image_version: '',
    pipe: 0,
  };
  export default {
    props: ['listData', 'cpData'],
    data() {
      return {
        setname: this.$t('bm.add.createPipeline'),
        tableListdetails: [],
        details: false,
        form: JSON.parse(JSON.stringify(formData)),
        pipeArray: [],
        dialogFormVisible: false,
        branchList: [],
        selectList: [],
        version: [],
        rules: {
          name: [
            { required: true, message: '请输入流水线描述', trigger: 'blur' },
          ],
          pipe: [
            { required: true, message: '请选择绑定流程', trigger: 'blur' },
          ],
          image_version: [
            { required: true, message: '请输入流水线名称', trigger: 'blur' },
            { min: 1, max: 64, message: this.$t('bm.add.long0T16Char'), trigger: 'blur' },
          ],
        },
        tableList: [],
        dataList: [],
      };
    },
    computed: {
      ...mapGetters({
        projectID: 'projectID',
      }),
    },
    components: {
      Refresh,
    },
    methods: {
      golist(rows) {
        if (rows) {
          rows.forEach((row) => {
          });
        }
      },
      handleSelectAll(val) {
        this.selectList = val;
      },
      handleSelect(val) {
        this.selectList = val;
      },
      doSubmit(formName) {
        this.$refs[formName].validate((valid) => {
          if (valid) {
            const apps = [];
            if(this.form.pipe === 0) {
              Message.error("请您先通过'项目设置－项目流程'添加项目流程后，再操作'创建流水线'")
              return;
            }
            if(this.selectList.length === 0) {
              Message.error('请先选择应用！');
              return;
            }
            for(const i in this.selectList) {
              const st = {
                'app_id': this.selectList[i].id,
                'branch_name': this.selectList[i].branch_name,
                'compile_command': this.selectList[i].compile_command,
              };
              apps.push(st);
            }
            const params = {
              version_no: this.form.image_version,
              name: this.form.name,
              bind_pipeline_id: this.form.pipe,
              apps: apps,
            };
            const that = this;
            backend.addProjectCI(this.projectID, params, (data) => {
              Message.success(this.$t('bm.add.optionSuc'));
              that.$emit('getprojectReleaseList');
              this.handleClose();
            }, () => {
            });
          } else {
            return false;
          }
        });
      },
      toggleSelection(rows) {
        if (rows) {
          rows.forEach((row) => {
            this.$refs.multipleTable.toggleRowSelection(row, true);
          });
        }
      },

      doShow(details, publishid) {
        this.setname = this.$t('bm.add.createPipeline');
        this.dialogFormVisible = true;
        this.details = details;
        if (this.details === true) {
          this.setname = this.$t('bm.add.publishBillDetail');
          backend.getListdetail(this.projectID, publishid, (data) => {
            this.form.name = data.name;
            this.form.pipe = data.id;
            this.tableListdetails = data.apps;
          });
          return false;
        }
        backend.getProjectPipe(this.projectID, (data) => {
          if (data) {
            this.pipeArray = data;
            this.form.pipe = this.pipeArray[0].id;
          }          
        });
        this.form = Object.assign({}, formData);
        backend.getProjectApp(this.projectID, (data) => {
          this.tableList = data;
          this.dataList = data;
          this.$nextTick(() => {
            for (const i of this.tableList) {
              this.toggleSelection([i]);
            }
            this.selectList = Object.assign({},data);
          });
        });

        this.$refs.ruleForm && this.$refs.ruleForm.clearValidate();
      },
      handleClose() {
        this.dialogFormVisible = false;
      },
    },
  };
</script>
