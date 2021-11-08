<style>
.createDialog.hostCreate .lbsList.el-col-20 .el-input__inner {
  vertical-align: -1px;
}
</style>
<template>
  <el-dialog top='25vh' :title="$t('bm.authorManage.createUser')" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog hostCreate" width='50%' :before-close="doCancelCreate">
    <el-form :model="form" :rules="rules" ref="ruleForm">
      <el-form-item :label="$t('bm.operCenter.account')" prop='user'>
        <el-input v-model.trim="form.user" :placeholder="$t('bm.add.inputAcount')" maxlength="64" auto-complete="on" :disabled="isEdit"></el-input>
      </el-form-item>
      <el-form-item :label="$t('bm.authorManage.userName')" prop='name'>
        <el-input v-model.trim="form.name" :placeholder="$t('bm.add.inputUsername')" maxlength="64" auto-complete="off" :disabled="isEdit"></el-input>
      </el-form-item>
      <el-form-item :label="$t('bm.infrast.email')" prop='email'>
        <el-input v-model.trim="form.email" :placeholder="$t('bm.add.inputEmail')" maxlength="64" auto-complete="off" :disabled="isEdit"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop='password'>
        <el-input v-model.trim="form.password" :placeholder="$t('bm.add.inputPassword')" maxlength="64" type="password" auto-complete="off" :disabled="isEdit"></el-input>
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
import backend from '../../api/backend';
import createTemplate from '../../common/createTemplate';
import validate from '../../common/validate';
import keyTxts from '../../common/validateKeyTxt';

const formData = {
  user: '',
  name: '',
  email: '',
};
export default {
  mixins: [createTemplate, validate],
  props: ['user'],
  data() {
    return {
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      rolesList: [],
      rules: {
        user: [
          { required: true, message: this.$t('bm.add.inputAcount'), trigger: 'blur' },
          { validator: this.validateResourceKeyValue, trigger: 'blur' },
        ],
        name: [{ required: true, message: this.$t('bm.add.inputUsername'), trigger: 'blur' }],
        email: [
          { required: true, message: this.$t('bm.add.inputEmail'), trigger: 'blur' },
          { validator: this.validateEmail, trigger: 'blur' },
        ],
        password: [{ required: true, message: this.$t('bm.add.inputPassword'), trigger: 'blur' }],
      },
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  methods: {
    doCreate(flag, item) {
      this.dialogFormVisible = true;
      this.isEdit = flag;
      if (flag) {
        this.form.user = item.user;
        this.form.name = item.name;
        this.form.email = item.email;
      } else {
        this.form = Object.assign({}, formData);
      }
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            this.$emit('getlist');
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
          };
          backend.addUser(JSON.stringify(this.form), () => {
            successCallBack();
          });
          return false;
        }
      });
    },
  },
};
</script>
