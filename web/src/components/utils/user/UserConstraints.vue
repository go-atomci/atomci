<style>
.createDialog .icon-question {
  position: absolute;
  left: -60px;
  top: 13px;
}
.constraint-info {
  line-height: 50px;
}
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog" width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item v-if="!isEdit" :label="$t('bm.authorManage.resourceType')" prop="resType">
        <el-select v-model="form.resType" :placeholder="$t('bm.add.selectResType')" filterable @change="resourceChange">
          <el-option v-for="(item, index) in resourceTypeList" :key="index" :label="item.description" :value="item.resource_type">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="约束" prop="constraint">
        <el-select v-model="form.constraint" placeholder="请选择约束" :disabled="isEdit" filterable>
          <el-option v-for="(item, index) in resourceConsList"
            :key="index"
            :label="item.resource_constraint"
            :value="item.resource_constraint"
            :disabled="isEditConstrint(item.resource_constraint)">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item v-for="(citem, cindex) in constraintList" :key="cindex" :label="'值' + (cindex+1)">
        <div class="constraint-info" >
          <el-col :span="8">
            <el-input v-model="citem.value" auto-complete="off" placeholder="value"></el-input>
          </el-col>
          <el-col :span="3" style="text-align:center">
            <el-button @click="removeConstraint(cindex)">删除</el-button>
          </el-col>
        </div>
      </el-form-item>
      <el-form-item label=" ">
        <el-button style="width:80px" size="small" type="primary" @click="addResourceConstanst">{{$t('bm.authorManage.addCons')}}</el-button>
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
import { Message, Submenu } from 'element-ui';
import backend from '../../../api/backend';
import createTemplate from '../../../common/createTemplate';

const formData = {
  constraint: '',
  resType: '',
};
/*eslint-disable*/
export default {
  mixins: [createTemplate],
  // isUser： true是个人权限策略界面，false 是组权限策略界面
  props: ['isUser'],
  data() {
    return {
      // selOptions: [],
      resourceTypeList: [],
      resourceOperationList: [],
      resourceConsList: [],
      rules: {
        resType: [
          { required: true, message: "请选择资源类型", trigger: 'blur' },
        ],
        constraint: [
          { required: true, message: "请选择约束", trigger: 'blur' },
        ],
      },
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: formData,
      title: this.$t('bm.add.addResOperation'),
      constraintList: [{
        key: '',
        value: '',
      }],
      adminEditList: ['resourceType', 'cluster', 'resourceConstraint', 'resourceOperation'],
      noEditList: ['resourceType', 'cluster', 'resourceConstraint', 'resourceOperation', 'group', 'namespace'],
    };
  },
  created() {
    this.queryResoreceType();
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
      detailInfo: 'getUserInfo',
    }),
  },
  methods: {
    isEditConstrint(value) {
      if (!this.isEdit || value === '' || value === undefined || value === null) {
        return false;
      }
      if (this.detailInfo.admin === 1 && this.adminEditList.includes(value)) {
        return false
      }
      return this.noEditList.includes(value);
    },
    removeConstraint(index) {
      this.constraintList.splice(index, 1);
    },
    queryResoreceType() {
      backend.getResourceTypeList((data) => {
        if (data) {
          this.resourceTypeList = data;
        }
      });
    },
    resourceChange() {
      this.resourceOperationList = [];
      this.resourceConsList = [];
      this.constraintList = [{
        key: '',
        value: '',
      }];
      this.form.resOption = [];
      this.form.resourceCons = [];
      this.queryResourceDetail();
    },
    queryResourceDetail() {
     if (this.form.resType) {
        backend.getResourceTypeDetail(this.form.resType, (data) => {
          if (data.resource_operations) {
            this.resourceOperationList = data.resource_operations;
          }
          if (data.resource_constraints) {
            this.resourceConsList = data.resource_constraints;
          }
        });
      }
    },
    doCreate(flag, item) {
      this.constraintList = [{
        value: '',
      }];
      this.isEdit = flag;
      this.form = Object.assign({}, formData);
      this.$refs.ruleForm && this.$refs.ruleForm.clearValidate();
      if (flag) {
        let list = [];
        this.title = "编辑约束";
        item.value.forEach((subVal) => {
          const constraint = {
            value: subVal,
          }
          list.push(constraint);
        });
        this.constraintList = list;
      } else {
        this.title = "添加约束";
      }
      this.form = {
        constraint: item && item.key || '',
      };
      this.dialogFormVisible = true;
      this.isEdit = flag;
      this.queryResourceDetail();
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            this.$emit('getlist');
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
          };
          const constraint = {};
          let dataValue = [];
          this.constraintList && this.constraintList.forEach((item) => {
            dataValue.push(item.value);
          });
          if (this.isEdit) {
            // TODO: use default group system 
            backend.putUserConstraintsItem("system", this.$route.params.user, this.form.constraint, dataValue, () => {
              successCallBack();
            });
          } else {
            // TODO: use default group system 
            backend.postUserConstraints("system", this.$route.params.user, this.form.constraint, dataValue, () => {
              successCallBack();
            });
          }
        }
      })
    },
    addResourceConstanst() {
      // 添加资源约束
      this.constraintList.push({
        value: '',
      });
    },
  },
};
</script>
