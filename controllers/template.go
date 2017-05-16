package controllers

const (
	requestRoomTemp = `
	<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Vue.js grid component example</title>
    <link rel="stylesheet" href="/static/css/style_grid.css">
    <link rel="stylesheet" href="/static/css/style_modal.css">
       
    <!-- Delete ".min" for console warnings in development -->
    <script src="/static/js/vue_ok.js"></script>
    <script src="/static/js/vue-resource.js"></script>
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
