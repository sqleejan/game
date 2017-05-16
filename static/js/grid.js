// register the grid component
Vue.component('demo-grid', {
  template: '#grid-template',
  replace: true,
  props: {
    data: Array,
    columns: Array,
    filterKey: String
  },
  data: function () {
    var sortOrders = {}
    this.columns.forEach(function (key) {
      sortOrders[key] = 1
    })
    return {
      msg: sharedata,
      showModal: false,
      sortKey: '',
      sortOrders: sortOrders,
      postDict: {"UserId":"申请用户ID","nicname":"申请用户名","create_at":"创建时间","duration":"房间时长","active":"是否激活"}
    }
  },
  computed: {
    filteredData: function () {
      	
      var sortKey = this.sortKey
      var filterKey = this.filterKey && this.filterKey.toLowerCase()
      var order = this.sortOrders[sortKey] || 1
      var data = this.data
      console.log("filteredData")
      if (filterKey) {
        data = data.filter(function (row) {
          return Object.keys(row).some(function (key) {
            return String(row[key]).toLowerCase().indexOf(filterKey) > -1
          })
        })
      }
      if (sortKey) {
        data = data.slice().sort(function (a, b) {
          a = a[sortKey]
          b = b[sortKey]
          return (a === b ? 0 : a > b ? 1 : -1) * order
        })
      }
      return data
    }
  },
  filters: {
    
    capitalize: function (str) {
      console.log("capitalize")
      return str.charAt(0).toUpperCase() + str.slice(1)
    }
  },
  methods: {
  
    sortBy: function (key) {
        console.log("sortBy")
      this.sortKey = key
      this.sortOrders[key] = this.sortOrders[key] * -1
    }
  }
})


// bootstrap the demo
var demo = new Vue({
  el: '#demo',
  data: {
    searchQuery: '',
    gridColumns: ['UserId', 'nicname','create_at','duration',"active"],
    gridData: [],
    apiUrl: '/v1/room/request/list',
    sdata: sharedata
  },
  ready: function() {
    	console.log("ttttt")
					this.getCustomers()
				},
	methods: {
			getCustomers: function() {
			this.$http.get(this.apiUrl+"?token="+this.sdata.token)
					.then((response) => {
              bd=response.json()
              
							this.$set('gridData', bd.data)
					})
					.catch(function(response) {
							console.log(response)
					})
			}
	}
})
