import keyTxts from './validateKeyTxt';

export default {
  methods: {
    // 验证标识名
    validateName(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case value.trim().length < 1:
          callback(new Error(this.$t('bm.add.lengthMinOne', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case !/^[0-9a-zA-Z-]+$/.test(value):
          callback(new Error(this.$t('bm.add.letter_num', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        default:
          callback();
      }
    },
    // 验证中文名
    validateCnName(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case value.trim().length < 3:
          callback(new Error(this.$t('bm.add.lengthMinThree', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case !/^[\u4E00-\u9FA5\uF900-\uFA2D]+$/.test(value):
          callback(new Error(this.$t('bm.add.shouldChinese', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        default:
          callback();
      }
    },
    // 验证邮箱
    validateEmail(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error('邮箱地址不能为空'));
          break;
        case !/^(\w-*\.*)+@(\w-?)+(\.\w{2,})+$/.test(value):
          callback(new Error('邮箱地址不正确'));
          break;
        default:
          callback();
      }
    },
    // 验证手机
    validatePhone(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error('手机号码不能为空'));
          break;
        case !/^1\d{10}$/.test(value):
          callback(new Error('手机号码格式不正确'));
          break;
        default:
          callback();
      }
    },
    // 验证应用名称
    validateAppName(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case value.trim().length > 64:
          callback(new Error(this.$t('bm.add.lengthNo64Node', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case !/^[a-z0-9]+(?:[-][a-z0-9]+)+$/.test(value):
          callback(new Error(this.$t('bm.add.appNameMatchRule', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        default:
          callback();
      }
    },
    // 验证IP地址
    validateIP(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        default:
          callback();
      }
    },
    validateIP2(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case !/^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(value):
          callback(new Error(this.$t('bm.releaseSetting.hostIPAdrForamtErr', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        default:
          callback();
      }
    },
    // 验证键值对，中间用冒号分隔
    validateColonKeyValue(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case !/([0-9a-zA-Z-]+):([0-9a-zA-Z-]+)/.test(value):
          callback(this.$t('bm.add.keyValueColonPart', { info: this.$t(keyTxts[rule.validateKey]) }));
          break;
        default:
          callback();
      }
    },
    // 验证键值对，中间用等号分隔
    validateEqualKeyValue(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case !/([0-9a-zA-Z-]+)=([0-9a-zA-Z]+)/.test(value):
          callback(this.$t('bm.add.keyValueEqualPart', { info: this.$t(keyTxts[rule.validateKey]) }));
          break;
        default:
          callback();
      }
    },
    // 验证键值对，中间用 -, 分隔
    validateResourceKeyValue(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.add.inputInfo', { info: this.$t(keyTxts[rule.validateKey]) })));
          break;
        case !/^[0-9a-zA-Z-.]+$/.test(value):
          callback(this.$t('bm.add.keyValueResoucrePart', { info: this.$t(keyTxts[rule.validateKey]) }));
          break;
        default:
          callback();
      }
    },
    // 验证键值对，中间用 - 分隔  可为空
    validateResourceValue(rule, value, callback) {
      switch (true) {
        case value !== '' && !/^[0-9a-zA-Z-]+$/.test(value):
          callback(this.$t('bm.add.keyValueResoucrePart', { info: this.$t(keyTxts[rule.validateKey]) }));
          break;
        default:
          callback();
      }
    },
    validateDuration(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.alarms.intervalNoEmpty')));
          break;
        case Number(value) < 0 || /\./.test(value):
          callback(new Error(this.$t('bm.alarms.errMsgDuration')));
          break;
        default:
          callback();
      }
    },
    validateInterval(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.alarms.intervalNoEmpty')));
          break;
        case (Number(value) <= 0 || /\./.test(value)):
          callback(new Error(this.$t('bm.alarms.errMsgInterval')));
          break;
        default:
          callback();
      }
    },
    validateThreshold(rule, value, callback) {
      switch (true) {
        case value === '':
          callback(new Error(this.$t('bm.alarms.intervalNoEmpty')));
          break;
        case Number(value) < 0:
          callback(new Error(this.$t('bm.alarms.errMsgThreshold')));
          break;
        case Number(value) < 0.000001:
          callback(new Error(this.$t('bm.alarms.errMsgThresholdMin')));
          break;
        default:
          callback();
      }
    },
  },
};
