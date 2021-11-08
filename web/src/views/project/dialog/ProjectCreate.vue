<template>
  <div>
    <el-dialog top='10vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog" :before-close="doCancelCreate">
      <el-form :model="form" :rules="rules" ref="ruleForm">
        <el-row>
          <el-col :span="24">
            <el-form-item label="项目名称" prop='name'>
              <el-input v-model.trim="form.name" placeholder="请输入项目名称(1-20字符)" auto-complete="off" maxlength="20"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <!-- <el-form-item label="项目编号" prop="project_code">
          <div>
            <el-radio v-model="radio" @change="changeone" label="1">{{$t('bm.add.defaultCreate')}}</el-radio>
            <el-radio class="radioStyle" v-model="radio" @change="changetwo" label="2">{{$t('bm.add.custom')}}</el-radio>
          </div>
          <el-input v-model="form.project_code" :disabled="showInput" placeholder="请输入项目编号，不能为中文" maxlength="20"></el-input>
        </el-form-item> -->
        <el-row>
          <el-col :span="24">
            <el-form-item :label="$t('bm.deployCenter.projectDesc')">
              <el-input v-model.trim="form.description" auto-complete="off" type="textarea" maxlength="256" placeholder="请输入项目描述(0-256个字)" :autosize="{ minRows: 7}"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="doCancelCreate">{{$t('bm.other.cancel')}}</el-button>
        <el-button type="primary" @click="doSubmit" :loading="createLoading">{{$t('bm.other.confirm')}}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
  .radioStyle {
    position: absolute;
    left: 100px;
    top: 12px;
  }
</style>
<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '../../../api/backend';
import keyTxts from '../../../common/validateKeyTxt';
import createTemplate from '../../../common/createTemplate';

const formData = {
  name: '',
  description: '',
  status: '',
  project_code: '',
};

export default {
  mixins: [createTemplate],
  data() {
    return {
      projId: null,
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      rules: {
        name: [
          { required: true, message: this.$t('bm.deployCenter.projectName'), trigger: 'blur' },
        ]
      },
      createLoading: false,
      showInput: false,
      radio: '2',
      title: '新建项目',
    };
  },
  created() {
    this.$store.dispatch('setNeedLoading', false);
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  methods: {
    changeone() {
      this.showInput = true;
      this.form.scope = '';
    },
    changetwo() {
      this.showInput = false;
    },
    doCreate(flag, item) {
      this.projId = '';
      if (!flag) {
        this.title = '新建项目';
        this.form = Object.assign({}, formData);
      } else {
        this.title = '编辑项目';
        this.form = {
          name: item.name || '',
          status: item.status || '1',
          project_code: item.project_code || '',
          description: item.description || ''
        };
        this.projId = item.id;
      }
      this.dialogFormVisible = true;
      this.isEdit = flag;
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
            this.$emit('getlist');
            this.createLoading = false;
          };
          const params = {
            name: this.form.name,
            //project_code: this.form.project_code,
            description: this.form.description
          };
          if (this.isEdit) {
            backend.updateNewProject(this.projId, params, (data) => {
              successCallBack();
              this.createLoading = false;
            }, (errCb) => {
              this.createLoading = false;
            });
          } else {
            //params.status = parseInt(this.radio);
            backend.getProjectCreate(params, (data) => {
              successCallBack();
            }, (errCb) => {
              this.createLoading = false;
            });
          }
        } else {
          return false;
        }
      });
    },
  },
};
</script>
