<style scoped>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
.tag {
  float: left;
  min-width: 80px;
}
.tag .el-tag {
  margin-right: 5px;
}
.input200 {
  float: left;
  width: 200px;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="节点名称" prop="name">
        <el-select v-model="form.name"  filterable @change="nodeTypeChange">
          <el-option v-for="(item, index) in taskList" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="节点类型" prop="nodeType" v-show="isShow">
        <el-input v-model.trim="form.nodeType" :disabled="true"></el-input>
      </el-form-item>
      <el-form-item label="流转类型" prop="next">
        <el-select v-model="form.next"  filterable>
          <el-option v-for="(item, index) in nextList" :key="index" :label="item.label" :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <!-- <el-form-item label="参数" prop="params">
        <el-input v-model.trim="form.params" auto-complete="off" maxlength="10" placeholder="请输入参数" ></el-input>
      </el-form-item> -->
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

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      isEdit: true,
      // 是否属于编辑状态
      dialogFormVisible: false,
      form: {
        name: '',
        nodeType: '',
        next: 'auto',
      },
      title: '新增节点',
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
        next: [
          { required: true, message: '请选择流转类型', trigger: 'blur' },
        ],
      },
      pIndex: 0,
      thisIndex: 0,
      nextList: [
        {
          value: 'auto',
          label: this.$t('bm.add.autoFlow'),
        }, 
        {
          value: 'manual',
          label: this.$t('bm.add.manualFlow'),
        },
      ],
      isShow: true,
    };
  },
  props: {
    taskList: {
      type: Array,
      default: []
    }
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  created() {},
  methods: {
    doCreate(flag, parentIndex, index, item) {
      this.pIndex = parentIndex;
      this.thisIndex = index;
      this.isEdit = flag;
      if (flag) {
        this.title = '编辑节点';
        this.form = {
          name: item.step_id || '',
          nodeType: item.type || '',
          next: item.driver || '',
        };
        this.taskList.map((i) => {
          if(i.id === item.step_id){
            this.form.nodeType = item.type;
            this.form.sub_task = item.sub_task || '';
            return;
          }
        });
      } else {
        this.isShow = false;
        this.title = '新增节点';
        this.form = {
          name: '',
          next: 'auto',
        };
      }
      this.dialogFormVisible = true;
      this.isEdit = flag;
    },
    nodeTypeChange(item) {
      if(item) {
        this.taskList.map((i) => {
          if(i.id === item) {
            this.form.nodeType = i.type;
            this.form.sub_task = i.sub_task || [];
            return;
          }
        });
      }
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          let sname = '';
          this.taskList.map((i) => {
            if(i.id == this.form.name) {
              sname = i.name;
            }
          });
          const obj = {
            step_id: this.form.name || '',
            driver: this.form.next || '',
            type: this.form.nodeType || '',
            sub_task: this.form.sub_task || [],
            name: sname
          };
          this.$emit('updateNode',this.isEdit, this.pIndex, obj, this.thisIndex,);
        }
      });
    },
  },
};
</script>
