package controllers

const (
	requestRoomTemp = `
	<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Vue.js grid component example</title>
    <link rel="stylesheet" href="/bg/static/css/style_grid.css">
    <link rel="stylesheet" href="/bg/static/css/style_modal.css">
       
    <!-- Delete ".min" for console warnings in development -->
    <script src="/bg/static/js/vue_ok.js"></script>
    <script src="/bg/static/js/vue-resource.js"></script>
    </head>
  <body>
    <script>
      var sharedata = {
        msg: "",
        openid:"",
        token:"{{{.Token}}}"
      }
    </script>
   
    <nav class="navbar navbar-inverse navbar-fixed-top" role="navigation">
    <!-- component template -->
    <script type="text/x-template" id="grid-template">
      
      <!-- use the modal component, pass in the prop -->
      <modal  v-if="showModal" @close="showModal = false">
        <!--
          you can use custom content here to overwrite
          default content
        -->
        <h3 slot="header">操作确认</h3>
      </modal>
      <table v-if="filteredData.length">
        <thead>
          <tr>
            <th v-for="key in columns"
              @click="sortBy(key)"
              :class="{ active: sortKey == key }">
              {{ postDict[key] | capitalize }}
              <span class="arrow" :class="sortOrders[key] > 0 ? 'asc' : 'dsc'">
              </span>
            </th>
             <th>
               操作
             </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="entry in filteredData">
            <td v-for="key in columns">
              {{entry[key]}}
            </td>
            <td >
              <button id="show-modal" @click="showModal = true, msg.msg='是否要拒绝申请', msg.openid=entry['UserId']">拒绝</button>
              <button id="show-modal" @click="showModal = true, msg.msg='是否同意申请', msg.openid=entry['UserId']">同意</button>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-else>No matches found.</p>
    </script>

    <!-- demo root element -->
    <div id="demo">
      <form id="search">
        Search <input name="query" v-model="searchQuery">
      </form>
      <demo-grid
        :data="gridData"
        :columns="gridColumns"
        :filter-key="searchQuery">
      </demo-grid>
    </div>
    </nav>
    <script src="/static/js/grid.js"></script>
     <script type="text/x-template" id="modal-template">
      <transition name="modal">
        <div class="modal-mask">
          <div class="modal-wrapper">
            <div class="modal-container">

              <div class="modal-header">
                <slot name="header">
                  default header
                </slot>
              </div>

              <div class="modal-body">
                <slot name="body">
                  {{ body.msg }}
                </slot>
              </div>

              <div class="modal-footer">
                <slot name="footer">
                  <button class="modal-left-button" @click="$emit('close')">
                    取消
                  </button>
                  <button class="modal-default-button" @click="accept()">
                    确定
                  </button>
                </slot>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </script>
    <script>
      // register modal component
      Vue.component('modal', {
        template: '#modal-template',
        data:function () {
          var customActions = {
            del: {method: 'DELETE', url: '/v1/room/request{/id}?token={token}'},
            put: {method: 'POST', url: '/v1/room/request{/id}?token={token}'}
          }
          var resource = this.$resource('/v1/room/request{/id}', {}, customActions);
          return{
            body:sharedata,
            client: resource
          }
        },
      methods: {
			  accept: function() {
          if (this.body.msg=="是否同意申请"){
            this.client.put({id: this.body.openid,token: this.body.token})
					.then((response) => {
              console.log(response)
					})
					.catch(function(response) {
							console.log(response)
					})
          }else{
             this.client.del({id: this.body.openid,token: this.body.token})
					.then((response) => {
              console.log(response)
					})
					.catch(function(response) {
							console.log(response)
					})
          }
			  
          this.$emit('close')
			   }
         
    	  }
      })

      

    </script>

  </body>
</html>`
)

const bakhtml = `<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8">
  <title>欢乐派送管理后台</title>
  <link href="http://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
  <link rel="stylesheet" href="/bg/static/assets/css/materialize.css">
  <link rel="stylesheet" href="/bg/static/assets/css/toastr.min.css" type="text/css" />
  <link rel="stylesheet" href="/bg/static/css/main.css" media="screen" title="no title" charset="utf-8">
</head>
<body>

      <nav class="top red darken-4">
        <div class="nav-wrapper">
          <a href=" " class="brand-logo">后台管理系统</a >
          <ul id="nav-mobile" class="right hide-on-med-and-down">
            <li><a href="sass.html">Sass</a ></li>
            <li><a href="badges.html">admin</a ></li>
            <li><a href="collapsible.html">推出</a ></li>
          </ul>
        </div>
      </nav>
      <div class="row main-content">
        <div class="col s2">
          <ul class="list z-depth-5 blue-grey darken-4 clearfix">
            <li class="blue-grey waves-effect waves-light btn" id="roomlist">房间列表</li>
            <li class="blue-grey waves-effect waves-light btn" id="roombug">调度服务</li>
            <li class="blue-grey waves-effect waves-light btn" id="admintele">设置联系方式</li>
          </ul>
        </div>
        <div class="col s10 template">
          <div class="frombox">
            <div class="search left">
              <div class="left">关键字：</div>
              <div class="search-wrapper left">
                <input id="search">
              </div>
            </div>
            <a class='dropdown-button btn red left dropdown1' href='#' data-activates='dropdown1'>房间状态</a >

            <ul id='dropdown1' class='dropdown-content'>
              <li class="timeover"><a href="#!">申请中</a ></li>
              <li class="showactive"><a href="#!">激活</a ></li>
              <li class="useactive"><a href="#!">使用中</a ></li>
              <li class="timeovering"><a href="#!">快到期</a ></li>
              <li class="timeover"><a href="#!">已过期</a ></li>
            </ul>
            <a class='dropdown-button btn red left dropdown2' href='#' data-activates='dropdown2'>房间状态</a >

            <ul id='dropdown2' class='dropdown-content'>
              <li class="bugopen"><a href="#!">关闭</a ></li>
              <li class="bugclose"><a href="#!">开启</a ></li>
            </ul>
            <table>
              <thead class="roomlistheader">
                <tr>
                    <th>房间号</th>
                    <th>房间名称</th>
                    <th>申请时间</th>
                    <th>激活时间</th>
                    <th>启用时间</th>
                    <th>到期时间</th>
                    <th>购买时长</th>
                    <th>房间状态</th>
                    <th>操作</th>
                </tr>
              </thead>
              <thead class="roombugheader">
                <tr>
                    <th>房间号</th>
                    <th>房间名称</th>
                    <th>房间状态</th>
                </tr>
              </thead>

              <tbody class="roomlist">
                <tr>
                  <td>Alvin</td>
                  <td>Eclair</td>
                  <td>$0.87</td>
                  <td>Alvin</td>
                  <td>Eclair</td>
                  <td>$0.87</td>
                  <td>Alvin</td>
                  <td>
                    <ul class="btnlist">
                      <li><a class="waves-effect waves-light btn actived">激活</a ></li>
                      <li><a class="waves-effect waves-light btn close">注销</a ></li>
                      <li><a class="waves-effect waves-light btn banker">账单</a ></li>
                      <li><a class="waves-effect waves-light btn timer">设定时长</a ></li>
                      <li><a class="waves-effect waves-light btn delet">删除</a ></li>
                    </ul>
                  </td>
                </tr>
              </tbody>
              <tbody class="roombug">
                <tr>
                  <td>Alvin</td>
                  <td>Eclair</td>
                  <td>$0.87</td>
                  <td>Alvin</td>
                  <td>Eclair</td>
                  <td>$0.87</td>
                  <td>Alvin</td>
                  <td>
                    <ul class="btnlist">
                      <li><a class="waves-effect waves-light btn actived">激活</a ></li>
                      <li><a class="waves-effect waves-light btn close">注销</a ></li>
                      <li><a class="waves-effect waves-light btn banker">账单</a ></li>
                      <li><a class="waves-effect waves-light btn timer">设定时长</a ></li>
                      <li><a class="waves-effect waves-light btn delet">删除</a ></li>
                    </ul>
                  </td>
                </tr>
              </tbody>
            </table>

          </div>
          <ul class="pagination">
            <li class="disabled pevr"><a href="#!"><i class="material-icons">chevron_left</i></a ></li>
            <li class="waves-effect page"><a href="#!"></a ></li>
            <li class="waves-effect next"><a href="#!"><i class="material-icons">chevron_right</i></a ></li>
          </ul>
          <ul class="admintele">
            <li>管理员联系方式：</li>
            <li><input type="text" class="adminteleinput"></li>
            <li><a class="waves-effect waves-light btn" id="pushtele">确认</a ></li>
          </ul>
        </div>
      </div>
      <footer class="page-footer red darken-4">
        <div class="footer-copyright">
          <div class="container">
          © 2017 Copyright Text
          <a class="grey-text text-lighten-4 right" href="#!">More Links</a >
          </div>
        </div>
      </footer>
      <div class="mask">
      </div>
      <div class="mask-time pupbox">
        <ul>
          <li>
            <div class="roomtimeinput">
              <input type="text" name="" value="" placeholder="输入购买时长（小时）" class="timeinput">
            </div>
          </li>
          <li>
            <a class="waves-effect waves-light btn" id="updatatime">button</a >
          </li>
        </ul>
      </div>
      <div class="mask-room pupbox">
        <div class="roomBillInfbox">

        </div>
        <div class="preloader-wrapper loading">
          <div class="spinner-layer spinner-red-only">
            <div class="circle-clipper left">
              <div class="circle"></div>
            </div><div class="gap-patch">
              <div class="circle"></div>
            </div><div class="circle-clipper right">
              <div class="circle"></div>
            </div>
          </div>
        </div>

        <ul class="pagination">
          <li class="waves-effect page2"><a href="#!"></a ></li>
          <li class="waves-effect next2 btn"><a href="#!">加载更多</a ></li>
        </ul>

      </div>
      <div class="preloader-wrapper big active tips">
        <div class="spinner-layer spinner-red">
          <div class="circle-clipper left">
            <div class="circle"></div>
          </div><div class="gap-patch">
            <div class="circle"></div>
          </div><div class="circle-clipper right">
            <div class="circle"></div>
          </div>
        </div>
      </div>
      

<script src="/bg/static/assets/js/jq.js" charset="utf-8"></script>
<script src="/bg/static/assets/js/toastr.js"></script>
<script src="/bg/static/assets/js/materialize.min.js"></script>
<script type="text/javascript">
$().ready(function(){
  var token = {{.Token}}
  sessionStorage.admin = token
})
</script>
<script src="/bg/static/js/main.js" charset="utf-8"></script>

</body>
</html>`
