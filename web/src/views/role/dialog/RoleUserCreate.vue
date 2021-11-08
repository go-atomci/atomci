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
      <el-form-item :label="$t('bm.infrast.user')" prop="users">
        <el-select v-model="form.users" :placeholder="$t('bm.add.selectUser')" filterable multiple>
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
  users: [],
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
      title: this.$t('bm.add.bindUser'),
      rules: {
        users: [
          { required: true, message: '请选择用户', trigger: 'blur' },
        ],
      },
    };
  },
  created() {
    const options = [];
    backend.getBuUser(this.$route.params.dept, (data) => {
      if (data) {
        for (const i of data.users) {
          options.push({
            label: `${i.name} - ${i.user}`,
            value: i.user,
          });
        }
        this.selOptions = options;
      }
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
        users: [],
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
            users: this.form.users,
          };
          backend.addRoleBindUser(this.$route.params.dept, this.$route.params.role, cl, () => {
            successCallBack();
          });
        }
      });
    },
  },
};
</script>
