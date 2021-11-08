<style>
.createDialog .icon-question {
  position: absolute;
  left: -60px;
  top: 13px;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog" width='40%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item :label="$t('bm.add.perPolicy')" prop="polices">
        <el-select v-model="form.polices" :placeholder="$t('bm.add.selectPolicy')" filterable multiple>
          <el-option v-for="(item, index) in selOptions" :key="index" :label="item.label" :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="doCancelCreate">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit" :loading="loading">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '../../../api/backend';
import createTemplate from '../../../common/createTemplate';

const formData = {
  polices: [],
};

export default {
  mixins: [createTemplate],
  data() {
    return {
      selOptions: [],
      groupRoleList: [],
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: this.$t('bm.add.addResPolicy'),
      rules: {
        polices: [
          { required: true, message: '请选择权限策略', trigger: 'blur' },
        ],
      },
    };
  },
  created() {
    const options = [];
    backend.getPolicyList(this.$route.params.dept, (data) => {
      for (const i of data) {
        options.push({
          label: i.description,
          value: i.policy_name,
        });
      }
      this.selOptions = options;
    });
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  methods: {
    doCreate() {
      this.form = {
        polices: [],
      };
      this.dialogFormVisible = true;
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            this.$emit('getlist');
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
          };
          const cl = {
            policies: this.form.polices,
          };
          backend.addRolePolicies(this.$route.params.dept, this.$route.params.role, cl, () => {
            successCallBack();
          });
        }
      });
    },
  },
};
</script>
