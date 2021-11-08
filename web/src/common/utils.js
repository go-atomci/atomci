export default {
  format(date, fmt) {
    const o = {
      'M+': date.getMonth() + 1,
      'd+': date.getDate(),
      'h+': date.getHours(),
      'm+': date.getMinutes(),
      's+': date.getSeconds(),
      'q+': Math.floor((date.getMonth() + 3) / 3),
      S: date.getMilliseconds(),
    };
    if (/(y+)/.test(fmt)) {
      fmt = fmt.replace(RegExp.$1, `${date.getFullYear()}`.substr(4 - RegExp.$1.length));
    }
    // eslint-disable-next-line
    for (const k in o) {
      if (new RegExp(`(${k})`).test(fmt)) {
        fmt = fmt.replace(
          RegExp.$1,
          RegExp.$1.length === 1 ? o[k] : `00${o[k]}`.substr(`${o[k]}`.length)
        );
      }
    }
    return fmt;
  },
  // 通过内存单位得到时间需要乘的倍数
  getMemoryCountByUnit(intervalUnit) {
    switch (intervalUnit) {
      case 'GB':
        return 1024 * 1024 * 1024;
      case 'MB':
        return 1024 * 1024;
      case 'KB':
        return 1024;
    }
    return 1;
  },
  // 通过时间单位得到时间需要乘的倍数
  getTimeCountByUnit(intervalUnit) {
    switch (intervalUnit) {
      case this.$t('bm.serviceM.day'):
        return 1000 * 1000 * 1000 * 60 * 60 * 24;
      case this.$t('bm.serviceM.hour'):
        return 1000 * 1000 * 1000 * 60 * 60;
      case this.$t('bm.serviceM.minute'):
        return 1000 * 1000 * 1000 * 60;
      case this.$t('bm.serviceM.second'):
        return 1000 * 1000 * 1000;
      case this.$t('bm.serviceM.Millis'):
        return 1000 * 1000;
      case this.$t('bm.serviceM.Micros'):
        return 1000;
    }
    return 1;
  },
  getMemoryTxt(val) {
    let unit = 'B';
    let middle = parseInt(val / 1024 / 1024 / 1024, 10);
    if (middle > 0) {
      unit = 'GB';
    } else {
      middle = parseInt(val / 1024 / 1024, 10);
      if (middle > 0) {
        unit = 'MB';
      } else {
        middle = parseInt(val / 1024, 10);
        if (middle > 0) {
          unit = 'KB';
        }
      }
    }
    return [middle, unit];
  },
  getIntervalTxt(val) {
    let unit = this.$t('bm.serviceM.nanos');
    let middle = parseInt(val / 1000 / 1000 / 1000 / 60 / 60 / 24, 10);
    if (middle > 0) {
      unit = this.$t('bm.serviceM.day');
    } else {
      middle = parseInt(val / 1000 / 1000 / 1000 / 60 / 60, 10);
      if (middle > 0) {
        unit = this.$t('bm.serviceM.hour');
      } else {
        middle = parseInt(val / 1000 / 1000 / 1000 / 60, 10);
        if (middle > 0) {
          unit = this.$t('bm.serviceM.minute');
        } else {
          middle = parseInt(val / 1000 / 1000 / 1000, 10);
          if (middle > 0) {
            unit = this.$t('bm.serviceM.second');
          } else {
            middle = parseInt(val / 1000 / 1000, 10);
            if (middle > 0) {
              unit = this.$t('bm.serviceM.Millis');
            } else {
              middle = parseInt(val / 1000, 10);
              if (middle > 0) {
                unit = this.$t('bm.serviceM.Micros');
              } else {
                middle = val;
              }
            }
          }
        }
      }
    }
    return middle + unit;
  },
  validateEmail(email) {
    const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(email);
  },
};
