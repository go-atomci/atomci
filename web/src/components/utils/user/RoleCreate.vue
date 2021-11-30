<style>
.createDialog .icon-question {
  position: absolute;
  left: -60px;
  top: 13px;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog" width='40%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules" v-if="isUser">
      <el-form-item :label="$t('bm.add.groupRole')" prop="rolesList">
        <el-select v-model="form.rolesList" :placeholder="$t('bm.add.selectGroupRole')" multiple filterable>
          <el-option v-for="(item, index) in groupRoleList" :key="index" :label="item.description" :value="item.role">
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <el-form :model="form" ref="ruleForm" :rules="rules" v-else>
      <el-form-item :label="$t('bm.add.roleName')" prop="role">
        <el-input v-model.trim="form.role" auto-complete="off" maxlength="64" :placeholder="$t('bm.add.inputRoleName')" :disabled="isEdit"></el-input>
      </el-form-item>
      <el-form-item :label="$t('bm.serviceM.description')" prop="description">
        <el-input v-model.trim="form.description" auto-complete="off" maxlength="64" :placeholder="$t('bm.add.inputDescInfo')" ></el-input>
      </el-form-item>
      <el-form-item :label="$t('bm.authorManage.resourceOper')" prop="perPolicy" v-if="!isEdit">
        <el-select v-model="form.perPolicy" :placeholder="$t('bm.add.selectOperation')" multiple filterable>
          <el-option v-for="(item, index) in operationsList" :key="index" :label="item.description" :value="item.id">
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
import backend from '@/api/backend';
import createTemplate from '@/common/createTemplate';
import validate from '@/common/validate';

const formData = {
  role: '',
  description: '',
  perPolicy: [],
  policies: [],
  rolesList: [],
};

export default {
  mixins: [createTemplate, validate],
  props: ['isUser', 'isTypeTitle'],
  data() {
    return {
      groupRoleList: [],
      operationsList: [],
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: this.$t('bm.add.newRole'),
      rules: {
        rolesList: [
          { required: true, message: '请选择角色', trigger: 'blur' },
        ],
        role: [
          { required: true, message: '请选择角色名称', trigger: 'blur' },
          { validator: this.validateResourceKeyValue, trigger: 'blur', validateKey: 'nameType' },
        ],
        description: [
          { required: true, message: '描述信息不能为空', trigger: 'blur' },
        ],
        perPolicy: [
          { required: true, message: '请选择权限策略', trigger: 'blur' },
        ],
      },
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  mounted() {
    if (this.$props.isUser) {
      backend.getRoleList((data) => {
        if (data) {
          this.groupRoleList = data;
        }
      });
    } else {
      backend.getResourcesOperation((data) => {
        if (data) {
          this.operationsList = data;
        }
      });
    }
  },
  methods: {
    doCreate(flag, item) {
      // this.form = Object.assign({}, flag ? item : formData);
      this.isEdit = flag;
      if (flag) {
        this.title = this.$t('bm.add.updateRole');
        const perPolicy = item.policies && item.policies.map((subItem) => {
          return subItem.policy_name;
        });
        this.form = {
          role: item.role || '',
          description: item.description || '',
          policies: perPolicy || [],
        };
      } else {
        this.title = this.isTypeTitle === 'new' ? this.$t('bm.add.newRole') : this.$t('bm.add.bindingRole');
        this.form = {
          role: '',
          description: '',
          perPolicy: [],
          rolesList: [],
        };
      }
      this.dialogFormVisible = true;
      this.isEdit = flag;
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            this.$emit('getlist');
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
          };
          if (this.$props.isUser) {
            const cl = {
              roles: this.form.rolesList,
            };
            // default group use system
            backend.addGroupUserRole("system", this.$route.params.user, cl, () => {
              successCallBack();
            });
          } else {
            const cl = {
              role: this.form.role,
              // use default group system
              group: 'system',
              description: this.form.description,
              operations: this.form.perPolicy,
            };
            if (this.isEdit) {
              // , JSON.stringify({ "role": this.form.role })
              backend.updateGroupRole(this.form.role, cl, () => {
                successCallBack();
              });
            } else {
              backend.addRole(cl, () => {
                successCallBack();
              });
            }
          }
        }
      });
    },
  },
};
</script>
