<style>
.createDialog.hostCreate .lbsList.el-col-20 .el-input__inner {
  vertical-align: -1px;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog hostCreate" width='50%' :before-close="doCancelCreate">

    <el-form :model="form" :rules="rules" ref="ruleForm">
      <el-form-item :label="$t('bm.operCenter.account')" prop='user'>
        <el-input v-model.trim="form.user" :placeholder="$t('bm.add.inputAcount')" maxlength="64" auto-complete="on" :disabled="isEdit"></el-input>
      </el-form-item>
      <el-form-item :label="$t('bm.authorManage.userName')" prop='name'>
        <el-input v-model.trim="form.name" :placeholder="$t('bm.add.inputUsername')" maxlength="64" auto-complete="off"></el-input>
      </el-form-item>
      <el-form-item :label="$t('bm.infrast.email')" prop='email'>
        <el-input v-model.trim="form.email" :placeholder="$t('bm.add.inputEmail')" maxlength="64" auto-complete="off"></el-input>
      </el-form-item>

      <el-form-item label="密码" prop='password'>
        <el-input v-model.trim="form.password" :placeholder="$t('bm.add.inputPassword')" maxlength="64" type="password" auto-complete="off" ></el-input>
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
import backend from '@/api/backend';
import createTemplate from '@/common/createTemplate';
import validate from '@/common/validate';

const formData = {
  user: '',
  name: '',
  email: '',
  password: '',
};
export default {
  mixins: [createTemplate, validate],
  props: ['user'],
  data() {
    return {
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      title:  this.$t('bm.authorManage.createUser'),
      form: JSON.parse(JSON.stringify(formData)),
      rolesList: [],
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
    rules() {
      return {
        user: [
          { required: true, message: this.$t('bm.add.inputAcount'), trigger: 'blur' },
          { validator: this.validateResourceKeyValue, trigger: 'blur' },
        ],
        name: [{ required: true, message: this.$t('bm.add.inputUsername'), trigger: 'blur' }],
        email: [
          { required: true, message: this.$t('bm.add.inputEmail'), trigger: 'blur' },
          { validator: this.validateEmail, trigger: 'blur' },
        ],
        password: [{ required: !this.isEdit, message: this.$t('bm.add.inputPassword'), trigger: 'blur' }],
      }
    }
  },
  methods: {
    doCreate(flag, item) {
      this.dialogFormVisible = true;
      this.isEdit = flag;
      if (flag) {
        this.title = "编辑用户";
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
          const cl = {
              user: this.form.user,
              name: this.form.name,
              email: this.form.email,
          };

        if (this.form.password != undefined && this.form.password != '') {
          cl['password'] = this.form.password
        }
          if (this.isEdit) {
              backend.updateUser(this.form.user, cl, () => {
                successCallBack();
              });
          } else {
            backend.addUser(JSON.stringify(this.form), () => {
              successCallBack();
            });            
          }
        }
      });
    },
  },
};
</script>
