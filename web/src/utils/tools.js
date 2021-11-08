export const getParameterByName = (name, url) => {
  if (!url) {
    url = window.location.href;
  }
  name = name.replace(/[[]]/g, '\\$&');
  const regex = new RegExp(`[?&]${name}(=([^&#]*)|&|#|$)`);
  const results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, ' '));
};

export const addParameterByName = (name, key, url) => {
  if (!url) {
    url = window.location.href;
  }
  const array = url.split('#');
  array[0] = `${array[0]}&${name}=${key}`.replace(/[&?]{1,2}/, '?');
  const result = array.join('#');

  return result;
};


  /**
   *协议检查
   *
   * @returns
   */
export function  protocolCheck () {
    let protocol = ''
    if (window.location.protocol === 'https:') {
      protocol = 'https'
    } else if (window.location.protocol === 'http:') {
      protocol = 'http'
    }
    return protocol
}