<template>
  <div class="hello">
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="审批人">
        <el-input v-model="formInline.user" placeholder="审批人"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">查询</el-button>
      </el-form-item>
      <div>签名值：&nbsp;Server:<span style="color: red;font-weight: bold">{{ sign }}&nbsp;&nbsp;</span>
        JS:<span>{{ jsSign }}</span>
      </div>
      <div>msg：&nbsp;<span style="color: red;font-weight: bold">{{ msg }}</span></div>
    </el-form>
  </div>
</template>

<script>
import axios from 'axios';


const SendRequest = {
  Message: function (arg) {
  },
  getAxios: function () {
    const fetch = axios.create();
    fetch.interceptors.request.use((config) => {
      return config;
    }, error => Promise.reject(error));

    fetch.interceptors.response.use((response) => {
      // if (response.data.code == 403 || response.data.code == 401) {//如果是登录失败
      //   clearCookie('token');
      //   router.push({path: "/backend/login"});
      //   return Promise.reject({
      //     type: 'error',
      //     message: "权限异常" + response.data.msg + '(' + response.data.msg + ')',
      //     duration: 5000,
      //   });
      // }
      if (response.status !== 200) {
        this.Message({
          type: 'error',
          message: '接口错误',
          duration: 5000,
        });
        return Promise.reject(response.data);
      }
      return response.data;
    }, error => {
      return Promise.reject(error);
    });
    return fetch
  },
  makApiSign: function (params) {
    let code
    switch (params.method) {
      case "postJson":
        code = `post${params.uri}${params.timestamp}${params.secret}${JSON.stringify(params.form)}`.toLowerCase()
        break;
      case "put":
        code = `${params.method}${params.uri}${params.timestamp}${params.secret}${JSON.stringify(params.form)}`.toLowerCase()
        break;
      default:
        code = `${params.method}${params.uri}${params.timestamp}${params.secret}${this.serialize(params.form)}`.toLowerCase()
        console.log(code);
    }
    let base64 = require('js-base64').Base64;
    return this.encrypt(base64.encode(code), params.secret)
  },
  serialize: function (form) {
    const ordered = {};
    let s = ""
    Object.keys(form).sort().forEach(function (key) {
      ordered[key] = form[key];
      s += `${key}${form[key]}`
    });
    return s
  },
  send: function (params) {
    params["secret"] = "signxxx"
    params = this.orgCommonParameters(params);
    switch (params.method.toLowerCase()) {
      case "postJson":
        return this.getAxios().post(`${params.host}${params.uri}`, params.form, {headers: params.headers})
      case "post":
        return this.getAxios().post(`${params.host}${params.uri}`, params.form, {headers: params.headers})
      case "get":
        return this.getAxios().get(`${params.host}${params.uri}`, {params: params.form, headers: params.headers})
      case "delete":
        return this.getAxios().delete(`${params.host}${params.uri}`, {params: params.form, headers: params.headers})
      case "put":
        return this.getAxios().put(`${params.host}${params.uri}`, params.form, {headers: params.headers})
    }
  },
  orgCommonParameters: function (params) {
    const that = this;
    params.timestamp = (new Date()).getTime();
    params.uri = `/${params.app}/${params.uri}`
    let defaultHeaders = {
      "debug": true,
      "X-Sign": that.makApiSign(params),
      "X-Timestamp": params.timestamp,
    }
    // console.log(defaultHeaders["X-Sign"])
    params.headers = Object.assign(params.headers ? params.headers : {}, defaultHeaders)
    return params;
  },
  encrypt: function (base64Code, secret) {
    const crypto = require('crypto');
    return crypto.createHmac('sha1', secret)
        .update(base64Code)
        .digest()
        .toString('base64')
        .toLowerCase();
  },
}


export default {
  name: 'sign',
  data() {
    return {
      sign: "",
      jsSign: "",
      msg: "",
      formInline: {
        user: 'asdfasd',
        a: 'a',
        c: 'c',
        b: 'b',
      }
    }
  },
  methods: {
    onSubmit() {
      let arg = {
        host: "http://localhost:8192",
        uri: "page_sign/1",
        form: this.formInline,
        // method: "postJson",
        // method: "get",
        // method: "put",
        method: "delete",
        app: "base-wrapper",
      }
      SendRequest.send(arg).then(r => {
        this.sign = r.data;
        this.msg = r.message;
        console.log(r);
      }).catch(e => {
        console.log("err:", e);
      });
    },
  },


}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
