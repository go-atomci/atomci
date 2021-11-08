<style>
.createDialog.hostCreate .lbsList.el-col-20 .el-input__inner {
  vertical-align: -1px;
}
</style>
<template>
  <el-dialog top='25vh' :title="textInfo.title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog hostCreate" width='50%' :before-close="doCancelCreate">
    <el-form :model="form" :rules="rules" ref="ruleForm">
      <el-form-item :label="textInfo.nameLabel" prop='nameType'>
        <el-input v-model.trim="form.nameType" :placeholder="textInfo.descPlace" maxlength="64" auto-complete="off" :disabled="isEdit"></el-input>
      </el-form-item>
      <el-form-item :label="$t('bm.serviceM.description')" prop='description'>
        <el-input v-model.trim="form.description" maxlength="64" :placeholder="$t('bm.add.inputDescInfo')" auto-complete="off"></el-input>
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
import validate from '../../../common/validate';

const formData = {
  nameType: '',
  description: '',
};

export default {
  mixins: [createTemplate, validate],
  props: ['label'],
  data() {
    return {
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      rules: {},
      textInfo: {
        title: '',
        nameLabel: '',
        descPlace: '',
      },
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  created() {
    switch (this.$props.label) {
      case 'type':
        this.textInfo = {
          title: this.$t('bm.add.addResType'),
          nameLabel: this.$t('bm.authorManage.resourceType'),
          descPlace: this.$t('bm.add.inputResType'),
        };
        this.rules = {
          nameType: [
            { required: true, message: this.$t('bm.add.inputResType'), trigger: 'blur' },
            { validator: this.validateResourceKeyValue, trigger: 'blur', validateKey: 'nameType' },
          ],
          description: [{ required: true, message: this.$t('bm.add.inputDescInfo'), trigger: 'blur' }],
        };
        break;
      case 'operation':
        this.textInfo = {
          title: this.$t('bm.add.addResOper'),
          nameLabel: this.$t('bm.authorManage.resourceOper'),
          descPlace: this.$t('bm.add.inputResOper'),
        };
        this.rules = {
          nameType: [
            { required: true, message: this.$t('bm.add.inputResOper'), trigger: 'blur' },
            { validator: this.validateResourceKeyValue, trigger: 'blur', validateKey: 'nameType' },
          ],
          description: [{ required: true, message: this.$t('bm.add.inputDescInfo'), trigger: 'blur' }],
        };
        break;
      case 'constraint':
        this.textInfo = {
          title: this.$t('bm..add.addResCons'),
          nameLabel: this.$t('bm.authorManage.resourceCons'),
          descPlace: this.$t('bm.add.inputResCons'),
        };
        this.rules = {
          nameType: [
            { required: true, message: this.$t('bm.add.inputResCons'), trigger: 'blur' },
            { required: true, validator: this.validateResourceKeyValue, trigger: 'blur', validateKey: 'nameType' },
          ],
          description: [{ required: true, message: this.$t('bm.add.inputDescInfo'), trigger: 'blur' }],
        };
        break;
    }
  },
  methods: {
    doCreate(flag, item) {
      this.dialogFormVisible = true;
      this.isEdit = flag;
      switch (this.$props.label) {
        case 'type':
          if (flag) {
            this.form.nameType = item.resource_type || '';
            this.form.description = item.description || '';
            this.textInfo.title = this.$t('bm.add.updateResType');
          } else {
            this.form = Object.assign({}, formData);
            this.textInfo.title = this.$t('bm.add.addResType');
          }
          break;
        case 'operation':
          if (flag) {
            this.form.nameType = item.resource_operation || '';
            this.form.description = item.description || '';
            this.textInfo.title = this.$t('bm.add.updateResOper');
          } else {
            this.form = Object.assign({}, formData);
            this.textInfo.title = this.$t('bm.add.addResOper');
          }
          break;
        case 'constraint':
          if (flag) {
            this.form.nameType = item.resource_constraint || '';
            this.form.description = item.description || '';
            this.textInfo.title = this.$t('bm.add.updateResCons');
          } else {
            this.form = Object.assign({}, formData);
            this.textInfo.title = this.$t('bm..add.addResCons');
          }
          break;
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
          switch (this.$props.label) {
            case 'type': {
              const paramsBody = {
                resource_type: this.form.nameType,
                description: this.form.description,
              };
              if (this.isEdit) {
                backend.updateResourceType(this.form.nameType, paramsBody, () => {
                  successCallBack();
                });
              } else {
                backend.addResourceType(paramsBody, () => {
                  successCallBack();
                });
              }
              break;
            }
            case 'operation': {
              const paramsBody = {
                resource_operation: this.form.nameType,
                description: this.form.description,
              };
              if (this.isEdit) {
                backend.updateResourceOperations(this.$route.params.resourceType, this.form.nameType, paramsBody, () => {
                  successCallBack();
                });
              } else {
                backend.addResourceOperations(this.$route.params.resourceType, paramsBody, () => {
                  successCallBack();
                });
              }
              break;
            }
            case 'constraint': {
              const paramsBody = {
                resource_constraint: this.form.nameType,
                description: this.form.description,
              };
              if (this.isEdit) {
                backend.updateResourceConstraints(this.$route.params.resourceType, this.form.nameType, paramsBody, () => {
                  successCallBack();
                });
              } else {
                backend.addResourceConstraints(this.$route.params.resourceType, paramsBody, () => {
                  successCallBack();
                });
              }
              break;
            }
          }
        }
      });
    },
  },
};
</script>
