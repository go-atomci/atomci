<template>
  <el-dialog top='25vh' z-index="1100" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  :before-close="doCancelCreate">
    <el-form ref="ruleForm" :model="form" :rules="rules">
      <el-form-item :label="$t('bm.add.backStage')" prop="stage">
        <el-select v-model="form.stage" :placeholder="$t('bm.add.select')" filterable>
          <el-option v-for="(item, index) in tableData" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('bm.add.backReason')" prop="code">
        <el-select v-model="form.code" :placeholder="$t('bm.add.select')" filterable>
          <el-option v-for="(item, index) in codeList" :key="index" :label="item.key" :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('bm.serviceM.description')" v-if="form.code === 'other'" prop="message" style="padding-bottom:20px">
        <el-input v-model.trim="form.message" maxlength="64" :placeholder="$t('bm.add.inputDescInfo')" auto-complete="off"></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit('ruleForm')">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '@/api/backend';

const formData = {
  code: 'release_problem',
  message: '',
  stage: '',
};

export default {
  // inject: [getList],
  data() {
    return {
      form: formData,
      stage: '',
      tableData: [],
      // 是否属于编辑状态
      isEdit: false,
      // 是否属于查看状态
      isView: false,
      dialogFormVisible: false,
      showandhide: true,
      id: '',
      codeList: [
        {
          key: '代码冲突',
          value: 'code_conflict',
        },
        {
          key: '非最新代码',
          value: 'not_latest_code',
        },
        {
          key: '存在bug',
          value: 'bug',
        },
        {
          key: '发布问题',
          value: 'release_problem',
        },
        {
          key: '其他',
          value: 'other',
        },
      ],
      rules: {
        stage: [
          { required: true, message: this.$t('bm.add.inputBackStage'), trigger: 'blur' },
        ],
        code: [
          { required: false, message: this.$t('bm.add.inputBackReason'), trigger: 'blur' },
        ],
        message: [
          { required: true, message: '请输入描述信息', trigger: 'blur' },
        ],
      },
    };
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
    doSubmit(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          const body = {
            stage_id: this.form.stage,
            code: this.form.code,
            message: this.form.message,
          };
          const _this = this;
          backend.gocontinue(this.projectID, this.id, this.stageid, 'back-to', JSON.stringify(body), () => {
            Message.success(_this.$t('bm.add.optionSuc'));
            _this.$emit('getprojectReleaseList');
            _this.dialogFormVisible = false;
            _this.$emit('getprojectReleaseList');
            _this.form.stage = '';
            _this.form.message = '';
          });
        } else {
          return false;
        }
      });
    },
    handleClose() {
      this.dialogFormVisible = false;
    },
    show(id, stageid) {
      this.stageid = stageid;
      this.id = id;
      this.dialogFormVisible = true;
      this.form = Object.assign({}, formData);
      backend.goregression(this.projectID, id, stageid, 'back-to', (data) => {
        this.tableData = data;
      });
    },
    doCancelCreate() {
      this.dialogFormVisible = false;
      this.isEdit = false;
      this.isView = false;
      this.description = '';
    },
  },
};
</script>
